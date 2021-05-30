package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

type URLParams map[string]string

func (p *URLParams) ToString(escaped bool) string {
	object := *p
	result := ""

	for key, value := range object {
		if escaped {
			value = url.QueryEscape(value)
		}

		result += "&" + key + "=" + value
	}

	if len(object) > 0 {
		result = "?" + result[1:]
	}

	return result
}

func GetAuthLinks() fiber.Map {
	query_vk := &URLParams{}
	(*query_vk)["client_id"] = os.Getenv("AUTH_VK_CLIENT_ID")
	(*query_vk)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
	(*query_vk)["scope"] = "email"
	(*query_vk)["response_type"] = "code"
	(*query_vk)["state"] = "vk"

	query_gl := &URLParams{}
	(*query_gl)["client_id"] = os.Getenv("AUTH_GL_CLIENT_ID")
	(*query_gl)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
	(*query_gl)["scope"] = "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email"
	(*query_gl)["access_type"] = "online"
	(*query_gl)["include_granted_scopes"] = "false"
	(*query_gl)["response_type"] = "code"
	(*query_gl)["state"] = "google"

	return fiber.Map{
		"vk": "https://oauth.vk.com/authorize" + (*query_vk).ToString(true),
		"gl": "https://accounts.google.com/o/oauth2/v2/auth" + (*query_gl).ToString(true),
	}
}

func IsEmptyStruct(object interface{}) bool {
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

func TimeMSK_ToLocaleString() string {
	return TimeMSK_ToTime().Format("02.01.2006 15:04:05")
}

func DateMSK_ToLocaleString() string {
	return TimeMSK_ToTime().Format("02.01.2006")
}

func DateMSK_ToLocaleSepString() string {
	return TimeMSK_ToTime().Format("02-01-2006")
}

// SHA512 hash
func HashSHA512(str string) string {
	hash := sha512.New()

	hash.Write([]byte(str))

	return hex.EncodeToString(hash.Sum(nil))
}

// Logger
func Logger(uri string, ip string, status int) {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	f, err := os.OpenFile("./logs/"+DateMSK_ToLocaleSepString()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Printf("Error opening file: %v", err)
	} else {
		log.SetOutput(f)
		log.Printf("\nREQUEST (%s): %s\nFROM: %s\nSTATUS: %d\nMEMORY USAGE (KiB): Alloc = %v; TotalAlloc = %v; Sys = %v; NumGC = %v;\n\n",
			TimeMSK_ToLocaleString(), uri, ip, status,
			memStats.Alloc/1024, memStats.TotalAlloc/1024, memStats.Sys/1024, memStats.NumGC)

		defer f.Close()
	}
}
