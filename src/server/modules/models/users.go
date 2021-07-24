package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"ivankprod.ru/src/server/modules/db"
	"ivankprod.ru/src/server/modules/utils"
)

// User struct
type User struct {
	ID             int64  `db:"user_id"           json:"userID"`
	Group          int64  `db:"user_group"        json:"userGroup"`
	SocialID       string `db:"user_social_id"    json:"userSocialID"`
	NameFirst      string `db:"user_name_first"   json:"userNameFirst"`
	NameLast       string `db:"user_name_last"    json:"userNameLast"`
	AvatarPath     string `db:"user_avatar_path"  json:"userAvatarPath"`
	Email          string `db:"user_email"        json:"userEmail"`
	AccessToken    string `db:"user_access_token" json:"-"`
	LastAccessTime string `db:"user_last_access"  json:"userLastAccessTime"`
	Role           int    `db:"user_role"         json:"userRole"`
	RoleDesc       string `db:"user_role_desc"    json:"userRoleDesc"`
	Type           int    `db:"user_type"         json:"userType"`
	TypeDesc       string `db:"user_type_desc"    json:"userTypeDesc"`
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

// Get users conditions by type to map
func (users *Users) GetCondsByType(aType int) fiber.Map {
	result := make(fiber.Map)

	for _, v := range *users {
		switch t := v.Type; t {
		case 0:
			result["vk"] = aType == 0 || true
		case 1:
			result["ok"] = aType == 1 || true
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
	Role   int
	Page   int
}

// Add new user
func AddUser(user *User) (int64, error) {
	var (
		tNow = utils.TimeMSK_ToString()

		query = "INSERT INTO users (user_social_id, user_group, user_role, user_access_token, " +
			"user_avatar_path, user_email, user_name_first, user_name_last, user_last_access, user_type) " +
			"VALUES (:user_social_id, :user_group, :user_role, :user_access_token, :user_avatar_path, " +
			":user_email, :user_name_first, :user_name_last, :user_last_access, :user_type)"
	)

	(*user).LastAccessTime = tNow

	if (*user).Role == 0 {
		(*user).Role = 2
	}

	db, err := db.Connect()
	if err != nil {
		return 0, err
	}

	res, err := db.NamedExec(query, user)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if (*user).ID == 0 && (*user).Group == 0 {
		setUserGroup(id, id)
	}

	return id, nil
}

// Update user group if it's new user
func setUserGroup(uID int64, uGroup int64) error {
	type PQuery struct {
		ID    int64
		Group int64
	}

	var (
		query = "UPDATE users SET user_group = :group WHERE user_id = :id"

		pqs = &PQuery{ID: uID, Group: uGroup}
	)

	db, err := db.Connect()
	if err != nil {
		return err
	}

	_, err = db.NamedExec(query, pqs)
	if err != nil {
		return err
	}

	return nil
}

// Get user
func GetUser(uID int64) (*User, error) {
	type PQuery struct {
		ID int64
	}

	var (
		query = "SELECT users.*, users_roles.role AS user_role_desc, users_types.type AS user_type_desc FROM users " +
			"INNER JOIN users_roles INNER JOIN users_types ON " +
			"users.user_role = users_roles.id AND users.user_type = users_types.id WHERE users.user_id = :id LIMIT 1"

		pqs = &PQuery{ID: uID}

		result = &User{}
	)

	db, err := db.Connect()
	if err != nil {
		return result, err
	}

	row, err := db.NamedQuery(query, pqs)
	if err != nil {
		return result, err
	}

	defer row.Close()

	for row.Next() {
		if err := row.StructScan(result); err != nil {
			return result, err
		}
	}

	(*result).AccessToken = "<restricted>"

	return result, nil
}

// Get user credentials
func getUserCredentials(uID int64) (*User, error) {
	type PQuery struct {
		ID int64
	}

	var (
		query = "SELECT users.*, users_roles.role AS user_role_desc, users_types.type AS user_type_desc FROM users " +
			"INNER JOIN users_roles INNER JOIN users_types ON " +
			"users.user_role = users_roles.id AND users.user_type = users_types.id WHERE users.user_id = :id LIMIT 1"

		pqs = &PQuery{ID: uID}

		result = &User{}
	)

	db, err := db.Connect()
	if err != nil {
		return result, err
	}

	row, err := db.NamedQuery(query, pqs)
	if err != nil {
		return result, err
	}

	defer row.Close()

	for row.Next() {
		if err := row.StructScan(result); err != nil {
			return result, err
		}
	}

	return result, nil
}

// Get users by specified group
func GetUsersGroup(uGroup int64) (*Users, error) {
	type PQuery struct {
		Group int64
	}

	var (
		query = "SELECT users.*, users_roles.role AS user_role_desc, users_types.type AS user_type_desc FROM users " +
			"INNER JOIN users_roles INNER JOIN users_types ON " +
			"users.user_role = users_roles.id AND users.user_type = users_types.id WHERE users.user_group = :group LIMIT 4"

		pqs = &PQuery{Group: uGroup}

		result = &Users{}
	)

	db, err := db.Connect()
	if err != nil {
		return result, err
	}

	rows, err := db.NamedQuery(query, pqs)
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

// Check if user exists
func ExistsUser(uSocialID string, uSocialType int) (int64, int64, int, error) {
	type PQuery struct {
		ID   string
		Type int
	}

	var (
		query = "SELECT user_id, user_group, user_role FROM users WHERE user_social_id = :id AND user_type = :type LIMIT 1"

		pqs = &PQuery{ID: uSocialID, Type: uSocialType}

		result = &User{}
	)

	db, err := db.Connect()
	if err != nil {
		return 0, 0, 0, err
	}

	row, err := db.NamedQuery(query, pqs)
	if err != nil {
		return 0, 0, 0, err
	}

	defer row.Close()

	for row.Next() {
		if err := row.StructScan(result); err != nil {
			return 0, 0, 0, err
		}
	}

	return (*result).ID, (*result).Group, (*result).Role, nil
}

// Get all users by args
func GetUsers(args *ArgsGetUsers) (*Users, error) {
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

	db, err := db.Connect()
	if err != nil {
		return result, err
	}

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
func SignInUser(u *User) error {
	var (
		tNow = utils.TimeMSK_ToString()

		query = "UPDATE users SET user_access_token = :user_access_token, user_avatar_path = :user_avatar_path, user_email = :user_email, " +
			"user_name_first = :user_name_first, user_name_last = :user_name_last, user_last_access = :user_last_access " +
			"WHERE user_social_id = :user_social_id AND user_type = :user_type"
	)

	(*u).LastAccessTime = tNow

	db, err := db.Connect()
	if err != nil {
		return err
	}

	_, err = db.NamedExec(query, u)
	if err != nil {
		return err
	}

	return nil
}

// User update access time
func UpdateUserAccessTime(uID int64) error {
	type PQuery struct {
		ID   int64
		Time string
	}

	var (
		tNow = utils.TimeMSK_ToString()

		query = "UPDATE users SET user_last_access = :time WHERE user_id = :id"

		pqs = &PQuery{ID: uID, Time: tNow}
	)

	db, err := db.Connect()
	if err != nil {
		return err
	}

	_, err = db.NamedExec(query, pqs)
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
func IsAuthenticated(uAuth string, uAgent string) (*User, error) {
	uAuthParsed, err := userAuthParse(uAuth)
	if err != nil {
		return nil, err
	}
	if uAuthParsed == nil {
		return nil, nil
	}

	result, err := getUserCredentials((*uAuthParsed).ID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	/*resultArr, err := getUsersGroup((*result).Group, (*result).ID)
	if err != nil {
		return nil, err
	}*/

	uSessionHash := utils.HashSHA512(strconv.FormatInt(((*result).ID), 10) + (*result).SocialID + (*result).AccessToken + uAgent)
	if uSessionHash == (*uAuthParsed).Hash {
		(*result).AccessToken = "<restricted>"

		/*r := &Users{}
		(*r) = append([]User{}, *result)

		if resultArr != nil {
			(*r) = append((*r), (*resultArr)...)
		}*/

		return result, nil
	}

	return nil, nil
}
