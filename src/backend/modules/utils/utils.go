package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"os"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
)

type URLParams map[string]string

func (p *URLParams) ToString() string {
	object := *p
	result := ""

	for key, value := range object {
		result += "&" + key + "=" + value
	}

	if len(object) > 0 {
		result = "?" + result[1:]
	}

	return result
}

func GetAuthLinks() fiber.Map {
	query := &URLParams{}

	(*query)["client_id"] = os.Getenv("AUTH_VK_CLIENT_ID")
	(*query)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
	(*query)["scope"] = "email"
	(*query)["response_type"] = "code"

	return fiber.Map{"links": fiber.Map{
		"vk": "https://oauth.vk.com/authorize" + (*query).ToString(),
	}}
}

func IsEmptySctruct(object interface{}) bool {
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	if reflect.ValueOf(object).Kind() == reflect.Struct {
		empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()

		if reflect.DeepEqual(object, empty) {
			return true
		}
	}

	return false
}

func TimeMSK_ToTime() time.Time {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}

	return time.Now().In(loc)
}

func TimeMSK_ToString() string {
	return TimeMSK_ToTime().Format("2006-01-02 15:04:05")
}

func HashSHA512(str string) string {
	hash := sha512.New()

	hash.Write([]byte(str))

	return hex.EncodeToString(hash.Sum(nil))
}
