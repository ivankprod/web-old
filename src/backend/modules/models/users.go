package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"ivankprod.ru/src/backend/modules/db"
	"ivankprod.ru/src/backend/modules/utils"
)

// User struct
type User struct {
	ID             int    `db:"user_id"           json:"userID"`
	NameFirst      string `db:"user_name_first"   json:"userNameFirst"`
	NameLast       string `db:"user_name_last"    json:"userNameLast"`
	ProfileLink    string `db:"user_profile_link" json:"userProfileLink"`
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
type Users struct {
	Users []User `db:"users" json:"users"`
}

// Stringify users struct
func (users *Users) ToJSON() string {
	var (
		result []byte
		err    error
	)

	if result, err = json.Marshal(users.Users); err != nil {
		return err.Error()
	}

	return string(result)
}

// User auth struct
type UserAuth struct {
	ID   int
	Hash string
}

// Args struct for GetUsers function
type ArgsGetUsers struct {
	Search string
	Role   int
	Page   int
}

// Add new user
func AddUser(user *User) error {
	var (
		tNow = utils.TimeMSK_ToString()

		query = "INSERT INTO users (user_id, user_role, user_access_token, user_profile_link, " +
			"user_avatar_path, user_email, user_name_first, user_name_last, user_last_access, user_type) " +
			"VALUES (:user_id, :user_role, :user_access_token, :user_profile_link, :user_avatar_path, " +
			":user_email, :user_name_first, :user_name_last, :user_last_access, :user_type)"
	)

	(*user).LastAccessTime = tNow
	(*user).ProfileLink = "/user/" + strconv.Itoa((*user).ID) + "/"
	(*user).Role = 1

	db, err := db.Connect()
	if err != nil {
		return err
	}

	_, err = db.NamedExec(query, user)
	if err != nil {
		return err
	}

	return nil
}

// Get user
func GetUser(uID int) (*User, error) {
	type PQuery struct {
		ID int
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
func getUserCredentials(uID int) (*User, error) {
	type PQuery struct {
		ID int
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

// Check if user exists
func ExistsUser(uID int) (bool, error) {
	type PQuery struct {
		ID int
	}

	var (
		query = "SELECT count(*) FROM users WHERE users.user_id = :id"

		pqs = &PQuery{ID: uID}

		result = 0
	)

	db, err := db.Connect()
	if err != nil {
		return false, err
	}

	rows, err := db.NamedQuery(query, pqs)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&result); err != nil {
			return result > 0, err
		}
	}

	return result > 0, nil
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

		(*result).Users = append((*result).Users, (*user))
	}

	return result, nil
}

// User sign in
func SignInUser(u *User) error {
	var (
		tNow = utils.TimeMSK_ToString()

		query = "UPDATE users SET user_access_token = :user_access_token, user_avatar_path = :user_avatar_path, user_email = :user_email, " +
			"user_name_first = :user_name_first, user_name_last = :user_name_last, user_last_access = :user_last_access " +
			"WHERE user_id = :user_id"
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

// User auth string parser
// (format: userID:hash)
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

	(*result).ID, err = strconv.Atoi(strarr[0])
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

	uSessionHash := utils.HashSHA512(strconv.Itoa((*result).ID) + (*result).AccessToken + uAgent)
	if uSessionHash == (*uAuthParsed).Hash {
		(*result).AccessToken = "<restricted>"
		return result, nil
	}

	return nil, nil
}
