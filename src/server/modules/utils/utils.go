package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Error statuses
var ErrorStatus = map[string]string{
	"400": "Некорректный запрос",
	"401": "Необходима авторизация",
	"403": "Запрещено",
	"404": "Не найдено",
	"408": "Истекло время ожидания",
	"413": "Слишком большая нагрузка",
	"414": "Слишком длинный URI",
	"429": "Слишком много запросов",
	"500": "Внутреняя ошибка сервера",
	"502": "Неверный шлюз",
	"503": "Сервис недоступен",
	"505": "Версия HTTP не поддерживается",
	"520": "Неизвестная ошибка",
}

// Sitemap: SitemapPath type
type SitemapPath struct {
	ID       int64         `json:"id"`
	ParentID int64         `json:"parent_id"`
	Title    string        `json:"title"`
	Path     string        `json:"path"`
	Priority float64       `json:"priority"`
	Children []SitemapPath `json:"children"`
}

// Sitemap: add child to path
func (p *SitemapPath) addChild(parentID int64, child *SitemapPath) bool {
	for i, v := range (*p).Children {
		if v.ID == parentID {
			(*p).Children[i].Children = append((*p).Children[i].Children, *child)

			return true
		} else if len(v.Children) > 0 {
			return (*p).Children[i].addChild(parentID, child)
		}
	}

	return false
}

// Sitemap: type
type Sitemap []SitemapPath

// Sitemap: add child to path (top level)
func (p *Sitemap) addChild(parentID int64, child *SitemapPath) bool {
	for i, v := range *p {
		if v.ID == parentID {
			(*p)[i].Children = append((*p)[i].Children, *child)

			return true
		} else if len(v.Children) > 0 {
			return (*p)[i].addChild(parentID, child)
		}
	}

	return false
}

// Sitemap: remove path by index
func (p *Sitemap) removePath(index int) *Sitemap {
	*p = append((*p)[:index], (*p)[index+1:]...)

	return p
}

// Sitemap: nest sitemap
func (p *Sitemap) Nest() *Sitemap {
	var added []int64

	// insert childs
	for _, v := range *p {
		if v.ParentID != 0 {
			(*p).addChild(v.ParentID, &v)

			added = append(added, v.ID)
		}
	}

	// remove added childs from top level paths
	for _, id := range added {
		for i, v := range *p {
			if v.ID == id {
				p.removePath(i)

				break
			}
		}
	}

	return p
}

// Sitemap: returns path's child html
func childLookup(item *SitemapPath) string {
	html := "\n<li><a href=\"" + item.Path + "\" class=\"spa\">" + item.Title + "</a>"

	if len(item.Children) > 0 {
		for _, v := range item.Children {
			html += "<ul>" + childLookup(&v) + "</ul>"
		}
	}

	return html + "</li>"
}

// Sitemap: convert sitemap to HTML
func (p *Sitemap) ToHTMLString() *string {
	output := "<ul>"

	for i := range *p {
		output += childLookup(&(*p)[i])
	}

	output += "</ul>"
	return &output
}

// URL params: type
type URLParams map[string]interface{}

// URL params: toString
func (p *URLParams) ToString(escaped bool) string {
	result := ""
	keys := make([]string, 0, len(*p))

	for k := range *p {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		tv := fmt.Sprintf("%v", (*p)[key])

		if escaped {
			tv = url.QueryEscape(tv)
		}

		result += "&" + key + "=" + tv
	}

	if len(*p) > 0 {
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

	if reflect.TypeOf(object).Kind() == reflect.Ptr {
		object = reflect.ValueOf(object).Elem().Interface()
	}

	switch reflect.TypeOf(object).Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Chan:
		if reflect.ValueOf(object).Len() == 0 {
			return true
		}
	}

	empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()

	if reflect.DeepEqual(object, empty) {
		return true
	}

	return false
}

// Time functions: to time
func TimeMSK_ToTime(mock ...time.Time) time.Time {
	now := time.Now()

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}

	if len(mock) > 1 {
		now = mock[0]
	}

	return now.In(loc).Add(time.Hour * time.Duration(3))
}

// Time functions: to string
func TimeMSK_ToString(mock ...time.Time) string {
	return TimeMSK_ToTime(mock...).Format("2006-01-02 15:04:05")
}

// Time functions: to locale string
func TimeMSK_ToLocaleString(mock ...time.Time) string {
	return TimeMSK_ToTime(mock...).Format("02.01.2006 15:04:05")
}

// Time functions: date to locale string
func DateMSK_ToLocaleString(mock ...time.Time) string {
	return TimeMSK_ToTime(mock...).Format("02.01.2006")
}

// Time functions: date to locale string with - separator
func DateMSK_ToLocaleSepString(mock ...time.Time) string {
	return TimeMSK_ToTime(mock...).Format("02-01-2006")
}

// SHA512 hash
func HashSHA512(str string) string {
	hash := sha512.New()

	hash.Write([]byte(str))

	return hex.EncodeToString(hash.Sum(nil))
}

// Logger (for dev only)
func DevLogger(uri string, ip string, status int) bool {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	f, err := os.OpenFile("./logs/"+DateMSK_ToLocaleSepString()+"_dev.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Printf("Error opening devlog file: %v\n", err)

		return false
	} else {
		log.SetFlags(0)
		log.SetOutput(f)

		log.Printf("\nREQUEST (%s): %s\nFROM: %s\nSTATUS: %d\nMEMORY USAGE (KiB): Alloc = %v; TotalAlloc = %v; Sys = %v; NumGC = %v;\n\n",
			TimeMSK_ToLocaleString(), uri, ip, status,
			memStats.Alloc/1024, memStats.TotalAlloc/1024, memStats.Sys/1024, memStats.NumGC)

		defer f.Close()
	}

	return true
}
