package repositories

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tarantool/go-tarantool"

	"ivankprod.ru/src/server/internal/domain"
	"ivankprod.ru/src/server/pkg/utils"
)

type AX []interface{}

type UserRepository interface {
	Create(uDTO *domain.UserCreateDTO) (*domain.User, error)
	FindOne(uID uint64) (*domain.User, error)
	FindOneBySocialID(uDTO *domain.UserFindOneBySocialIDDTO) (*domain.User, error)
	FindGroup(uGroup uint64) (*domain.Users, error)
	FindAll(uDTO *domain.UserFindAllDTO) (*domain.Users, error)
	Update(uID uint64, uDTO *domain.UserUpdateDTO) (*domain.User, error)
}

type userRepository struct {
	db *tarantool.Connection
}

func NewUserRepository(dbc *tarantool.Connection) UserRepository {
	return &userRepository{
		db: dbc,
	}
}

func (r *userRepository) Create(uDTO *domain.UserCreateDTO) (*domain.User, error) {
	var (
		tuplesRoles domain.UserRoles
		tuplesTypes domain.UserTypes
		tuplesUsers domain.Users

		err error
	)

	if uDTO == nil {
		return nil, fmt.Errorf("UserRepository error: UserCreateDTO can't be empty")
	}

	if uDTO.Role == 0 {
		uDTO.Role = domain.USER_ROLE_GUEST
	}

	err = r.db.InsertTyped("users", AX{
		nil, uDTO.Group, uDTO.SocialID, uDTO.AccessToken,
		uDTO.AvatarPath, uDTO.Email, uDTO.NameFirst, uDTO.NameLast, utils.TimeMSK_ToLocaleString(),
		uDTO.Role, uDTO.Type}, &tuplesUsers)
	if err != nil {
		return nil, err
	}

	if len(tuplesUsers) == 0 {
		return nil, nil
	}

	id := tuplesUsers[0].ID

	if uDTO.Group == 0 {
		return r.Update(id, &domain.UserUpdateDTO{
			Group: &id,
		})
	}

	err = r.db.SelectTyped("users_roles", "primary_id", 0, 1, tarantool.IterEq, AX{tuplesUsers[0].Role}, &tuplesRoles)
	if err != nil {
		return nil, err
	}

	err = r.db.SelectTyped("users_types", "primary_id", 0, 1, tarantool.IterEq, AX{tuplesUsers[0].Type}, &tuplesTypes)
	if err != nil {
		return nil, err
	}

	tuplesUsers[0].RoleDesc = tuplesRoles[0].Role
	tuplesUsers[0].TypeDesc = tuplesTypes[0].Type

	return &tuplesUsers[0], nil
}

func (r *userRepository) FindOne(uID uint64) (*domain.User, error) {
	var (
		tuplesRoles domain.UserRoles
		tuplesTypes domain.UserTypes
		tuplesUsers domain.Users

		err error
	)

	err = r.db.SelectTyped("users", "primary_id", 0, 1, tarantool.IterEq, AX{uID}, &tuplesUsers)
	if err != nil {
		return nil, err
	}

	if len(tuplesUsers) == 0 {
		return nil, nil
	}

	err = r.db.SelectTyped("users_roles", "primary_id", 0, 1, tarantool.IterEq, AX{tuplesUsers[0].Role}, &tuplesRoles)
	if err != nil {
		return nil, err
	}

	err = r.db.SelectTyped("users_types", "primary_id", 0, 1, tarantool.IterEq, AX{tuplesUsers[0].Type}, &tuplesTypes)
	if err != nil {
		return nil, err
	}

	tuplesUsers[0].RoleDesc = tuplesRoles[0].Role
	tuplesUsers[0].TypeDesc = tuplesTypes[0].Type

	return &tuplesUsers[0], nil
}

func (r *userRepository) FindOneBySocialID(uDTO *domain.UserFindOneBySocialIDDTO) (*domain.User, error) {
	var (
		tuplesRoles domain.UserRoles
		tuplesTypes domain.UserTypes
		tuplesUsers domain.Users

		err error
	)

	if uDTO == nil {
		return nil, fmt.Errorf("UserRepository error: UserFindBySocialIDDTO can't be empty")
	}

	err = r.db.SelectTyped("users", "secondary_socialid_type", 0, 1, tarantool.IterEq, AX{uDTO.SocialID, uDTO.Type}, &tuplesUsers)
	if err != nil {
		return nil, err
	}

	if len(tuplesUsers) == 0 {
		return nil, nil
	}

	err = r.db.SelectTyped("users_roles", "primary_id", 0, 1, tarantool.IterEq, AX{tuplesUsers[0].Role}, &tuplesRoles)
	if err != nil {
		return nil, err
	}

	err = r.db.SelectTyped("users_types", "primary_id", 0, 1, tarantool.IterEq, AX{tuplesUsers[0].Type}, &tuplesTypes)
	if err != nil {
		return nil, err
	}

	tuplesUsers[0].RoleDesc = tuplesRoles[0].Role
	tuplesUsers[0].TypeDesc = tuplesTypes[0].Type

	return &tuplesUsers[0], nil
}

func (r *userRepository) FindGroup(uGroup uint64) (*domain.Users, error) {
	var (
		tuplesRoles domain.UserRoles
		tuplesTypes domain.UserTypes
		tuplesUsers domain.Users

		err error
	)

	err = r.db.SelectTyped("users", "secondary_group", 0, 100, tarantool.IterEq, AX{uGroup}, &tuplesUsers)
	if err != nil {
		return nil, err
	}

	if len(tuplesUsers) == 0 {
		return nil, nil
	}

	err = r.db.SelectTyped("users_roles", "primary_id", 0, 100, tarantool.IterEq, AX{}, &tuplesRoles)
	if err != nil {
		return nil, err
	}

	err = r.db.SelectTyped("users_types", "primary_id", 0, 100, tarantool.IterEq, AX{}, &tuplesTypes)
	if err != nil {
		return nil, err
	}

	for i, v := range tuplesUsers {
		tuplesUsers[i].RoleDesc = tuplesRoles[v.Role-1].Role
		tuplesUsers[i].TypeDesc = tuplesTypes[v.Type].Type
	}

	return &tuplesUsers, nil
}

func (r *userRepository) FindAll(uDTO *domain.UserFindAllDTO) (*domain.Users, error) {
	if uDTO == nil {
		return nil, fmt.Errorf("UserRepository error: UserFindBySocialIDDTO can't be empty")
	}

	var (
		query = "SELECT \"users\".*, \"users_roles\".\"role\" AS \"user_role_desc\", \"users_types\".\"type\" " +
			"AS \"user_type_desc\" FROM \"users\" " +
			"INNER JOIN \"users_roles\" INNER JOIN \"users_types\" ON " +
			"\"users\".\"user_role\" = \"users_roles\".\"id\" AND \"users\".\"user_type\" = \"users_types\".\"id\" "

		where = ""
		limit = ""

		tuplesUsers domain.Users
	)

	if uDTO.Search != nil || uDTO.Role != nil {
		where += "WHERE "

		if uDTO.Search != nil {
			where += "(\"users\".\"user_email\" LIKE '%" + *uDTO.Search + "%' OR (\"users\".\"user_name_first\" || ' ' || \"users\".\"user_name_last\") LIKE '%" + *uDTO.Search + "%')"

			if uDTO.Role != nil {
				where += " AND "
			}
		}

		if uDTO.Role != nil {
			where += "(\"users\".\"user_role\" = " + strconv.FormatUint(*uDTO.Role, 10) + ")"
		}
	}

	if uDTO.Page != nil {
		limit += " LIMIT " + strconv.FormatUint((*uDTO.Page-1)*10, 10) + ", " + "10"
	} else {
		limit += " LIMIT 10"
	}

	query += where + " ORDER BY \"users\".\"user_role\" DESC" + limit

	resp, err := r.db.Call("box.execute", AX{query})
	if err != nil {
		return nil, err
	}

	respData := resp.Data

	if len(respData) > 1 {
		respError, _ := respData[1].([]interface{})[0].(string)
		if respError != "" {
			return nil, fmt.Errorf("UserRepository error: SQL error: %s", respError)
		}
	}

	respParsed, ok := respData[0].([]interface{})[0].(map[interface{}]interface{})["rows"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("UserRepository error during parsing SQL response")
	}

	for _, v := range respParsed {
		data := v.([]interface{})

		tuplesUsers = append(tuplesUsers, domain.User{
			ID:             data[0].(uint64),
			Group:          data[1].(uint64),
			SocialID:       data[2].(string),
			NameFirst:      data[6].(string),
			NameLast:       data[7].(string),
			AvatarPath:     data[4].(string),
			Email:          data[5].(string),
			AccessToken:    data[3].(string),
			LastAccessTime: data[8].(string),
			Role:           data[9].(uint64),
			RoleDesc:       data[11].(string),
			Type:           data[10].(uint64),
			TypeDesc:       data[12].(string),
		})
	}

	return &tuplesUsers, nil
}

func (r *userRepository) Update(uID uint64, uDTO *domain.UserUpdateDTO) (*domain.User, error) {
	var (
		tuplesRoles domain.UserRoles
		tuplesTypes domain.UserTypes
		tuplesUsers domain.Users

		err error
	)

	if uDTO == nil {
		return nil, fmt.Errorf("UserRepository error: UserUpdateDTO can't be empty")
	}

	set := make([]interface{}, 0)
	re := regexp.MustCompile(`[A-Z][a-z0-9]*`)

	utils.IterateStruct(*uDTO, func(field string, value interface{}) {
		f := strings.Join(re.FindAllString(field, -1), "_")
		set = append(set, AX{"=", "user_" + strings.ToLower(f), value})
	})

	err = r.db.UpdateTyped("users", "primary_id", AX{uID}, set, &tuplesUsers)
	if err != nil {
		return nil, err
	}

	if len(tuplesUsers) == 0 {
		return nil, nil
	}

	err = r.db.SelectTyped("users_roles", "primary_id", 0, 1, tarantool.IterEq, AX{tuplesUsers[0].Role}, &tuplesRoles)
	if err != nil {
		return nil, err
	}

	err = r.db.SelectTyped("users_types", "primary_id", 0, 1, tarantool.IterEq, AX{tuplesUsers[0].Type}, &tuplesTypes)
	if err != nil {
		return nil, err
	}

	tuplesUsers[0].RoleDesc = tuplesRoles[0].Role
	tuplesUsers[0].TypeDesc = tuplesTypes[0].Type

	return &tuplesUsers[0], nil
}
