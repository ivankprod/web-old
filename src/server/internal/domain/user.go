package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/vmihailenco/msgpack.v2"
)

// User Role constants
const (
	USER_ROLE_ADMINISTRATOR uint64 = 4
	USER_ROLE_WEBMASTER     uint64 = 3
	USER_ROLE_GUEST         uint64 = 2
	USER_ROLE_BANNED        uint64 = 1
)

// User struct
type User struct {
	ID             uint64 `json:"userID"`
	Group          uint64 `json:"userGroup"`
	SocialID       string `json:"userSocialID"`
	NameFirst      string `json:"userNameFirst"`
	NameLast       string `json:"userNameLast"`
	AvatarPath     string `json:"userAvatarPath"`
	Email          string `json:"userEmail"`
	AccessToken    string `json:"-"`
	LastAccessTime string `json:"userLastAccessTime"`
	Role           uint64 `json:"userRole"`
	RoleDesc       string `json:"userRoleDesc"`
	Type           uint64 `json:"userType"`
	TypeDesc       string `json:"userTypeDesc"`
}

// UserUpdateDTO struct
type UserCreateDTO struct {
	Group       uint64
	SocialID    string
	NameFirst   string
	NameLast    string
	AvatarPath  string
	Email       string
	AccessToken string
	Role        uint64
	Type        uint64
}

// UserFindOneBySocialIDDTO struct
type UserFindOneBySocialIDDTO struct {
	SocialID string
	Type     uint64
}

// UserFindAllDTO struct
type UserFindAllDTO struct {
	Search *string
	Role   *uint64
	Page   *uint64
}

// UserUpdateDTO struct
type UserUpdateDTO struct {
	Group       *uint64
	SocialID    *string
	NameFirst   *string
	NameLast    *string
	AvatarPath  *string
	Email       *string
	AccessToken *string
	LastAccess  *string
	Role        *uint64
	Type        *uint64
}

// UserSignInDTO struct
type UserSignInDTO struct {
	NameFirst   string
	NameLast    string
	AvatarPath  string
	Email       string
	AccessToken string
}

// User struct: msgpack encoder
func (u *User) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeSliceLen(11); err != nil {
		return err
	}

	if err := e.EncodeUint64(u.ID); err != nil {
		return err
	}

	if err := e.EncodeUint64(u.Group); err != nil {
		return err
	}

	if err := e.EncodeString(u.SocialID); err != nil {
		return err
	}

	if err := e.EncodeString(u.AccessToken); err != nil {
		return err
	}

	if err := e.EncodeString(u.AvatarPath); err != nil {
		return err
	}

	if err := e.EncodeString(u.Email); err != nil {
		return err
	}

	if err := e.EncodeString(u.NameFirst); err != nil {
		return err
	}

	if err := e.EncodeString(u.NameLast); err != nil {
		return err
	}

	if err := e.EncodeString(u.LastAccessTime); err != nil {
		return err
	}

	if err := e.EncodeUint64(u.Role); err != nil {
		return err
	}

	if err := e.EncodeUint64(u.Type); err != nil {
		return err
	}

	return nil
}

// User struct: msgpack decoder
func (u *User) DecodeMsgpack(d *msgpack.Decoder) error {
	var (
		l   int
		err error
	)

	if l, err = d.DecodeSliceLen(); err != nil {
		return err
	}

	if l != 11 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}

	if u.ID, err = d.DecodeUint64(); err != nil {
		return err
	}

	if u.Group, err = d.DecodeUint64(); err != nil {
		return err
	}

	if u.SocialID, err = d.DecodeString(); err != nil {
		return err
	}

	if u.AccessToken, err = d.DecodeString(); err != nil {
		return err
	}

	if u.AvatarPath, err = d.DecodeString(); err != nil {
		return err
	}

	if u.Email, err = d.DecodeString(); err != nil {
		return err
	}

	if u.NameFirst, err = d.DecodeString(); err != nil {
		return err
	}

	if u.NameLast, err = d.DecodeString(); err != nil {
		return err
	}

	if u.LastAccessTime, err = d.DecodeString(); err != nil {
		return err
	}

	if u.Role, err = d.DecodeUint64(); err != nil {
		return err
	}

	if u.Type, err = d.DecodeUint64(); err != nil {
		return err
	}

	return nil
}

// Stringify user struct
func (user *User) ToJSON() string {
	var (
		result []byte
		err    error
	)

	if result, err = json.Marshal(user); err != nil {
		return err.Error()
	}

	return string(result)
}

// Users struct
type Users []User

// Stringify users struct
func (users *Users) ToJSON() string {
	var (
		result []byte
		err    error
	)

	if result, err = json.Marshal(users); err != nil {
		return err.Error()
	}

	return string(result)
}

// UserRole struct
type UserRole struct {
	ID   uint64 `json:"roleID"`
	Role string `json:"roleDesc"`
	Sort uint64 `json:"roleSort"`
}

// UserRole struct: msgpack encoder
func (r *UserRole) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeSliceLen(3); err != nil {
		return err
	}

	if err := e.EncodeUint64(r.ID); err != nil {
		return err
	}

	if err := e.EncodeString(r.Role); err != nil {
		return err
	}

	if err := e.EncodeUint64(r.Sort); err != nil {
		return err
	}

	return nil
}

// UserRole struct: msgpack decoder
func (r *UserRole) DecodeMsgpack(d *msgpack.Decoder) error {
	var (
		l   int
		err error
	)

	if l, err = d.DecodeSliceLen(); err != nil {
		return err
	}

	if l != 3 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}

	if r.ID, err = d.DecodeUint64(); err != nil {
		return err
	}

	if r.Role, err = d.DecodeString(); err != nil {
		return err
	}

	if r.Sort, err = d.DecodeUint64(); err != nil {
		return err
	}

	return nil
}

// Stringify UserRole struct
func (r *UserRole) ToJSON() string {
	var (
		result []byte
		err    error
	)

	if result, err = json.Marshal(r); err != nil {
		return err.Error()
	}

	return string(result)
}

// UserRoles struct
type UserRoles []UserRole

// Stringify UserRoles struct
func (r *UserRoles) ToJSON() string {
	var (
		result []byte
		err    error
	)

	if result, err = json.Marshal(r); err != nil {
		return err.Error()
	}

	return string(result)
}

// UserType struct
type UserType struct {
	ID   uint64 `json:"typeID"`
	Type string `json:"typeDesc"`
}

// Stringify UserType struct
func (t *UserType) ToJSON() string {
	var (
		result []byte
		err    error
	)

	if result, err = json.Marshal(t); err != nil {
		return err.Error()
	}

	return string(result)
}

// UserType struct: msgpack encoder
func (t *UserType) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeSliceLen(2); err != nil {
		return err
	}

	if err := e.EncodeUint64(t.ID); err != nil {
		return err
	}

	if err := e.EncodeString(t.Type); err != nil {
		return err
	}

	return nil
}

// UserType struct: msgpack decoder
func (t *UserType) DecodeMsgpack(d *msgpack.Decoder) error {
	var (
		l   int
		err error
	)

	if l, err = d.DecodeSliceLen(); err != nil {
		return err
	}

	if l != 2 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}

	if t.ID, err = d.DecodeUint64(); err != nil {
		return err
	}

	if t.Type, err = d.DecodeString(); err != nil {
		return err
	}

	return nil
}

// UserTypes struct
type UserTypes []UserType

// Stringify UserTypes struct
func (t *UserTypes) ToJSON() string {
	var (
		result []byte
		err    error
	)

	if result, err = json.Marshal(t); err != nil {
		return err.Error()
	}

	return string(result)
}

// Get users conditions by type to map
func (users *Users) GetCondsByType(aType uint64) fiber.Map {
	result := make(fiber.Map)

	for _, v := range *users {
		switch t := v.Type; t {
		case 0:
			result["vk"] = aType == 0 || true
		case 1:
			result["ya"] = aType == 1 || true
		case 2:
			result["fb"] = aType == 2 || true
		case 3:
			result["gl"] = aType == 3 || true
		}
	}

	return result
}

// User auth struct
type UserAuth struct {
	ID   uint64
	Hash string
}

// User auth string parser
// (format: userID:userGroup:hash)
func NewUserAuthFromString(str string) (*UserAuth, error) {
	var err error

	if str == "" {
		return nil, nil
	}

	result := &UserAuth{}
	strarr := strings.Split(str, ":")

	if len(strarr) < 2 {
		return nil, nil
	}
	if strarr[0] == "" || strarr[1] == "" {
		return nil, nil
	}

	result.ID, err = strconv.ParseUint(strarr[0], 10, 64)
	if err != nil {
		return result, err
	}

	result.Hash = strarr[1]

	return result, nil
}
