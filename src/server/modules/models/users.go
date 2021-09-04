package models

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tarantool/go-tarantool"
	"gopkg.in/vmihailenco/msgpack.v2"

	"ivankprod.ru/src/server/modules/utils"
)

// User struct
type User struct {
	ID             int64  `json:"userID"`
	Group          int64  `json:"userGroup"`
	SocialID       string `json:"userSocialID"`
	NameFirst      string `json:"userNameFirst"`
	NameLast       string `json:"userNameLast"`
	AvatarPath     string `json:"userAvatarPath"`
	Email          string `json:"userEmail"`
	AccessToken    string `json:"userAccessToken"`
	LastAccessTime string `json:"userLastAccessTime"`
	Role           int64  `json:"userRole"`
	RoleDesc       string `json:"userRoleDesc"`
	Type           int64  `json:"userType"`
	TypeDesc       string `json:"userTypeDesc"`
}

// User struct: msgpack encoder
func (u *User) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeSliceLen(11); err != nil {
		return err
	}

	if err := e.EncodeInt64(u.ID); err != nil {
		return err
	}

	if err := e.EncodeInt64(u.Group); err != nil {
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

	if err := e.EncodeInt64(u.Role); err != nil {
		return err
	}

	if err := e.EncodeInt64(u.Type); err != nil {
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

	if u.ID, err = d.DecodeInt64(); err != nil {
		return err
	}

	if u.Group, err = d.DecodeInt64(); err != nil {
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

	if u.Role, err = d.DecodeInt64(); err != nil {
		return err
	}

	if u.Type, err = d.DecodeInt64(); err != nil {
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
	ID   int64  `json:"roleID"`
	Role string `json:"roleDesc"`
	Sort int64  `json:"roleSort"`
}

// UserRole struct: msgpack encoder
func (r *UserRole) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeSliceLen(3); err != nil {
		return err
	}

	if err := e.EncodeInt64(r.ID); err != nil {
		return err
	}

	if err := e.EncodeString(r.Role); err != nil {
		return err
	}

	if err := e.EncodeInt64(r.Sort); err != nil {
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

	if r.ID, err = d.DecodeInt64(); err != nil {
		return err
	}

	if r.Role, err = d.DecodeString(); err != nil {
		return err
	}

	if r.Sort, err = d.DecodeInt64(); err != nil {
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
	ID   int64  `json:"typeID"`
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

	if err := e.EncodeInt64(t.ID); err != nil {
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

	if t.ID, err = d.DecodeInt64(); err != nil {
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
func (users *Users) GetCondsByType(aType int) fiber.Map {
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
	ID   int64
	Hash string
}

// Args struct for GetUsers function
type ArgsGetUsers struct {
	Search string
	Role   int64
	Page   int64
}

// Add new user
func AddUser(db *tarantool.Connection, user *User) (int64, error) {
	var (
		tuples Users
		tNow   = utils.TimeMSK_ToString()
	)

	user.LastAccessTime = tNow

	if user.Role == 0 {
		user.Role = 2
	}

	err := db.InsertTyped("users", []interface{}{
		nil, user.Group, user.SocialID, user.AccessToken,
		user.AvatarPath, user.Email, user.NameFirst, user.NameLast, user.LastAccessTime,
		user.Role, user.Type}, &tuples)
	if err != nil {
		return 0, err
	}

	id := tuples[0].ID

	if user.ID == 0 && user.Group == 0 {
		setUserGroup(db, id, id)
	}

	return id, nil
}

// Update user group if it's new user
func setUserGroup(db *tarantool.Connection, uID int64, uGroup int64) error {
	_, err := db.Update("users", "primary_id", []interface{}{uID}, []interface{}{[]interface{}{"=", "user_group", uGroup}})
	if err != nil {
		return err
	}

	return nil
}

// Get user
func GetUser(db *tarantool.Connection, uID int64) (*User, error) {
	var (
		tuplesRoles UserRoles
		tuplesTypes UserTypes
		tuplesUsers Users

		err error
	)

	err = db.SelectTyped("users_roles", "primary_id", 0, 4, tarantool.IterEq, []interface{}{}, &tuplesRoles)
	if err != nil {
		return nil, err
	}

	err = db.SelectTyped("users_types", "primary_id", 0, 4, tarantool.IterEq, []interface{}{}, &tuplesTypes)
	if err != nil {
		return nil, err
	}

	err = db.SelectTyped("users", "primary_id", 0, 1, tarantool.IterEq, []interface{}{uID}, &tuplesUsers)
	if err != nil {
		return nil, err
	}

	tuplesUsers[0].RoleDesc = tuplesRoles[tuplesUsers[0].Role-1].Role
	tuplesUsers[0].TypeDesc = tuplesTypes[tuplesUsers[0].Type].Type
	tuplesUsers[0].AccessToken = "<restricted>"

	return &tuplesUsers[0], nil
}

// Get user credentials
func getUserCredentials(db *tarantool.Connection, uID int64) (*User, error) {
	var (
		tuplesRoles UserRoles
		tuplesTypes UserTypes
		tuplesUsers Users

		err error
	)

	err = db.SelectTyped("users_roles", "primary_id", 0, 4, tarantool.IterEq, []interface{}{}, &tuplesRoles)
	if err != nil {
		return nil, err
	}

	err = db.SelectTyped("users_types", "primary_id", 0, 4, tarantool.IterEq, []interface{}{}, &tuplesTypes)
	if err != nil {
		return nil, err
	}

	err = db.SelectTyped("users", "primary_id", 0, 1, tarantool.IterEq, []interface{}{uID}, &tuplesUsers)
	if err != nil {
		return nil, err
	}

	tuplesUsers[0].RoleDesc = tuplesRoles[tuplesUsers[0].Role-1].Role
	tuplesUsers[0].TypeDesc = tuplesTypes[tuplesUsers[0].Type].Type

	return &tuplesUsers[0], nil
}

// Get users by specified group
func GetUsersGroup(db *tarantool.Connection, uGroup int64) (*Users, error) {
	var (
		tuplesRoles UserRoles
		tuplesTypes UserTypes
		tuplesUsers Users

		err error
	)

	err = db.SelectTyped("users_roles", "primary_id", 0, 4, tarantool.IterEq, []interface{}{}, &tuplesRoles)
	if err != nil {
		return nil, err
	}

	err = db.SelectTyped("users_types", "primary_id", 0, 4, tarantool.IterEq, []interface{}{}, &tuplesTypes)
	if err != nil {
		return nil, err
	}

	err = db.SelectTyped("users", "secondary_group", 0, 4, tarantool.IterEq, []interface{}{uGroup}, &tuplesUsers)
	if err != nil {
		return nil, err
	}

	for i, v := range tuplesUsers {
		tuplesUsers[i].RoleDesc = tuplesRoles[v.Role-1].Role
		tuplesUsers[i].TypeDesc = tuplesTypes[v.Type].Type
	}

	return &tuplesUsers, nil
}

// Check if user exists
func ExistsUser(db *tarantool.Connection, uSocialID string, uSocialType int) (int64, int64, int64, error) {
	var (
		tuplesUsers Users

		err error
	)

	err = db.SelectTyped("users", "secondary_socialid_type", 0, 1, tarantool.IterEq, []interface{}{uSocialID, uSocialType}, &tuplesUsers)
	if err != nil {
		return 0, 0, 0, err
	}

	tuplesUsers[0].AccessToken = "<restricted>"

	return tuplesUsers[0].ID, tuplesUsers[0].Group, tuplesUsers[0].Role, nil
}

/*
// Get all users by args
func GetUsers(db *tarantool.Connection, args *ArgsGetUsers) (*Users, error) {
	type PQuery struct {
		Search string
		Role   int
	}

	var (
		query = "SELECT users.*, users_roles.role AS user_role_desc, users_types.type AS user_type_desc FROM users " +
			"INNER JOIN users_roles INNER JOIN users_types ON " +
			"users.user_role = users_roles.id AND users.user_type = users_types.id "
		search = (*args).Search
		where  = ""
		limit  = ""

		role = (*args).Role
		page = (*args).Page

		result = &Users{}
	)

	if search != "" || role != 0 {
		where += "WHERE "

		if search != "" {
			where += "(users.user_email LIKE :search OR concat(users.user_name_first, ' ', users.user_name_last) LIKE :search)"

			if role != 0 {
				where += " AND "
			}
		}

		if role != 0 {
			where += "(users.user_role = :role)"
		}
	}

	if page != 0 {
		limit += " LIMIT " + strconv.Itoa((page-1)*10) + ", " + "10"
	}

	pqs := &PQuery{Search: "%" + search + "%", Role: role}

	rows, err := db.NamedQuery(query+where+" ORDER BY users.user_role DESC"+limit, pqs)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		user := &User{}

		if err := rows.StructScan(user); err != nil {
			return result, err
		}

		(*user).AccessToken = "<restricted>"

		(*result) = append((*result), (*user))
	}

	return result, nil
}

// User sign in
func SignInUser(db *tarantool.Connection, u *User) error {
	var (
		tNow = utils.TimeMSK_ToString()

		query = "UPDATE users SET user_access_token = :user_access_token, user_avatar_path = :user_avatar_path, user_email = :user_email, " +
			"user_name_first = :user_name_first, user_name_last = :user_name_last, user_last_access = :user_last_access " +
			"WHERE user_social_id = :user_social_id AND user_type = :user_type"
	)

	(*u).LastAccessTime = tNow

	_, err := db.NamedExec(query, u)
	if err != nil {
		return err
	}

	return nil
}

// User update access time
func UpdateUserAccessTime(db *tarantool.Connection, uID int64) error {
	type PQuery struct {
		ID   int64
		Time string
	}

	var (
		tNow = utils.TimeMSK_ToString()

		query = "UPDATE users SET user_last_access = :time WHERE user_id = :id"

		pqs = &PQuery{ID: uID, Time: tNow}
	)

	_, err := db.NamedExec(query, pqs)
	if err != nil {
		return err
	}

	return nil
}

// User auth string parser
// (format: userID:userGroup:hash)
func userAuthParse(str string) (*UserAuth, error) {
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

	(*result).ID, err = strconv.ParseInt(strarr[0], 10, 64)
	if err != nil {
		return result, err
	}

	(*result).Hash = strarr[1]

	return result, nil
}

// Check for login
func IsAuthenticated(db *tarantool.Connection, uAuth string, uAgent string) (*User, error) {
	uAuthParsed, err := userAuthParse(uAuth)
	if err != nil {
		return nil, err
	}
	if uAuthParsed == nil {
		return nil, nil
	}

	result, err := getUserCredentials(db, (*uAuthParsed).ID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	uSessionHash := utils.HashSHA512(strconv.FormatInt(((*result).ID), 10) + (*result).SocialID + (*result).AccessToken + uAgent)
	if uSessionHash == (*uAuthParsed).Hash {
		(*result).AccessToken = "<restricted>"

		return result, nil
	}

	return nil, nil
}*/
