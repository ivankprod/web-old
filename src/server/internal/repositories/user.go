package repositories

import (
	"fmt"
	"strings"

	"github.com/tarantool/go-tarantool"

	"ivankprod.ru/src/server/internal/domain"
	"ivankprod.ru/src/server/pkg/utils"
)

type AX []interface{}

type UserRepository interface {
	Create(user *domain.User) (uint64, error)
	FindOne(uID uint64) (*domain.User, error)
	Update(uID uint64, uDTO *domain.UserDTO) error
}

type userRepository struct {
	db *tarantool.Connection
}

func NewUserRepository(dbc *tarantool.Connection) UserRepository {
	return &userRepository{
		db: dbc,
	}
}

func (r *userRepository) Create(user *domain.User) (uint64, error) {
	var tuplesUsers domain.Users

	if user.Role == 0 {
		user.Role = domain.USER_ROLE_GUEST
	}

	err := r.db.InsertTyped("users", AX{
		nil, user.Group, user.SocialID, user.AccessToken,
		user.AvatarPath, user.Email, user.NameFirst, user.NameLast, utils.TimeMSK_ToLocaleString(),
		user.Role, user.Type}, &tuplesUsers)
	if err != nil {
		return 0, err
	}

	if len(tuplesUsers) == 0 {
		return 0, nil
	}

	id := tuplesUsers[0].ID

	if user.ID == 0 && user.Group == 0 {
		if err := r.setUserGroup(id, id); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (r *userRepository) FindOne(uID uint64) (*domain.User, error) {
	var (
		tuplesRoles domain.UserRoles
		tuplesTypes domain.UserTypes
		tuplesUsers domain.Users

		err error
	)

	err = r.db.SelectTyped("users_roles", "primary_id", 0, 4, tarantool.IterEq, AX{}, &tuplesRoles)
	if err != nil {
		return nil, err
	}

	err = r.db.SelectTyped("users_types", "primary_id", 0, 4, tarantool.IterEq, AX{}, &tuplesTypes)
	if err != nil {
		return nil, err
	}

	err = r.db.SelectTyped("users", "primary_id", 0, 1, tarantool.IterEq, AX{uID}, &tuplesUsers)
	if err != nil {
		return nil, err
	}

	if len(tuplesUsers) == 0 {
		return nil, nil
	}

	tuplesUsers[0].RoleDesc = tuplesRoles[tuplesUsers[0].Role-1].Role
	tuplesUsers[0].TypeDesc = tuplesTypes[tuplesUsers[0].Type].Type

	return &tuplesUsers[0], nil
}

func (r *userRepository) Update(uID uint64, uDTO *domain.UserDTO) error {
	if uDTO == nil {
		return fmt.Errorf("UserService error: UserDTO can't be empty")
	}

	set := make([]interface{}, 0)

	utils.IterateStruct(uDTO, func(field string, value interface{}) {
		set = append(set, AX{"=", "user_" + strings.ToLower(field), value})
	})

	_, err := r.db.Update("users", "primary_id", AX{uID}, set)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) setUserGroup(uID uint64, uGroup uint64) error {
	_, err := r.db.Update("users", "primary_id", AX{uID}, AX{AX{"=", "user_group", uGroup}})
	if err != nil {
		return err
	}

	return nil
}
