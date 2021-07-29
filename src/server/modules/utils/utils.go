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

// Sitemap: SitemapPath type
type SitemapPath struct {
	ID       int64         `json:"id"`
	ParentID int64         `json:"parent_id"`
	Title    string        `json:"title"`
	Path     string        `json:"path"`
	Priority float64       `json:"priority"`
	Children []SitemapPath `json:"children"`
}

// Sitemap: type
type Sitemap []SitemapPath

// Sitemap: get path index by ID
func (p *Sitemap) GetPathIndexByID(id int64) int {
	for i, v := range *p {
		if v.ID == id {
			return i
		}
	}

	return -1
}

// Sitemap: remove path by index
func (p *Sitemap) RemovePath(index int) {
	*p = append((*p)[:index], (*p)[index+1:]...)
}

// Sitemap: nest sitemap
func (p *Sitemap) Nest() {
	for i, v := range *p {
		if v.ParentID != 0 {
			parentIndex := p.GetPathIndexByID(v.ParentID)
			(*p)[parentIndex].Children = append((*p)[parentIndex].Children, v)
			p.RemovePath(i)
		}
	}
}

// URL params: type
type URLParams map[string]string

// URL params: toString
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

// OAuth links
func GetAuthLinks() fiber.Map {
	query_vk := &URLParams{}
	(*query_vk)["client_id"] = os.Getenv("AUTH_VK_CLIENT_ID")
	(*query_vk)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
	(*query_vk)["scope"] = "email"
	(*query_vk)["response_type"] = "code"
	(*query_vk)["state"] = "vk"

	query_fb := &URLParams{}
	(*query_fb)["client_id"] = os.Getenv("AUTH_FB_CLIENT_ID")
	(*query_fb)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
	(*query_fb)["scope"] = "email"
	(*query_fb)["response_type"] = "code"
	(*query_fb)["state"] = "facebook"

	query_ya := &URLParams{}
	(*query_ya)["client_id"] = os.Getenv("AUTH_YA_CLIENT_ID")
	(*query_ya)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
	(*query_ya)["response_type"] = "code"
	(*query_ya)["state"] = "yandex"

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
		"fb": "https://www.facebook.com/v11.0/dialog/oauth" + (*query_fb).ToString(true),
		"ya": "https://oauth.yandex.ru/authorize" + (*query_ya).ToString(true),
		"gl": "https://accounts.google.com/o/oauth2/v2/auth" + (*query_gl).ToString(true),
	}
}

// Check for empty
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

// Time functions: to time
func TimeMSK_ToTime() time.Time {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}

	return time.Now().In(loc)
}

// Time functions: to string
func TimeMSK_ToString() string {
	return TimeMSK_ToTime().Format("2006-01-02 15:04:05")
}

// Time functions: to locale string
func TimeMSK_ToLocaleString() string {
	return TimeMSK_ToTime().Format("02.01.2006 15:04:05")
}

// Time functions: date to locale string
func DateMSK_ToLocaleString() string {
	return TimeMSK_ToTime().Format("02.01.2006")
}

// Time functions: date to locale string with - separator
func DateMSK_ToLocaleSepString() string {
	return TimeMSK_ToTime().Format("02-01-2006")
}

// SHA512 hash
func HashSHA512(str string) string {
	hash := sha512.New()

	hash.Write([]byte(str))

	return hex.EncodeToString(hash.Sum(nil))
}

// Logger (for dev only)
func DevLogger(uri string, ip string, status int) {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	f, err := os.OpenFile("./logs/"+DateMSK_ToLocaleSepString()+"_dev.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Printf("Error opening devlog file: %v", err)
	} else {
		log.SetOutput(f)
		log.Printf("\nREQUEST (%s): %s\nFROM: %s\nSTATUS: %d\nMEMORY USAGE (KiB): Alloc = %v; TotalAlloc = %v; Sys = %v; NumGC = %v;\n\n",
			TimeMSK_ToLocaleString(), uri, ip, status,
			memStats.Alloc/1024, memStats.TotalAlloc/1024, memStats.Sys/1024, memStats.NumGC)

		defer f.Close()
	}
}
