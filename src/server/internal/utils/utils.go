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
	for i, v := range *p {
		if v.ParentID != 0 {
			(*p).addChild(v.ParentID, &(*p)[i])

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
		for i := range item.Children {
			html += "<ul>" + childLookup(&item.Children[i]) + "</ul>"
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
	queryVK := &URLParams{
		"client_id":     os.Getenv("AUTH_VK_CLIENT_ID"),
		"redirect_uri":  "https://" + os.Getenv("SERVER_HOST") + "/auth/",
		"scope":         "email",
		"response_type": "code",
		"state":         "vk",
	}

	queryFB := &URLParams{
		"client_id":     os.Getenv("AUTH_FB_CLIENT_ID"),
		"redirect_uri":  "https://" + os.Getenv("SERVER_HOST") + "/auth/",
		"scope":         "email",
		"response_type": "code",
		"state":         "facebook",
	}

	queryYA := &URLParams{
		"client_id":     os.Getenv("AUTH_YA_CLIENT_ID"),
		"redirect_uri":  "https://" + os.Getenv("SERVER_HOST") + "/auth/",
		"response_type": "code",
		"state":         "yandex",
	}

	queryGL := &URLParams{
		"client_id":              os.Getenv("AUTH_GL_CLIENT_ID"),
		"redirect_uri":           "https://" + os.Getenv("SERVER_HOST") + "/auth/",
		"scope":                  "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
		"access_type":            "online",
		"include_granted_scopes": "false",
		"response_type":          "code",
		"state":                  "google",
	}

	return fiber.Map{
		"vk": "https://oauth.vk.com/authorize" + queryVK.ToString(true),
		"fb": "https://www.facebook.com/v11.0/dialog/oauth" + queryFB.ToString(true),
		"ya": "https://oauth.yandex.ru/authorize" + queryYA.ToString(true),
		"gl": "https://accounts.google.com/o/oauth2/v2/auth" + queryGL.ToString(true),
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

	return reflect.DeepEqual(object, empty)
}

// Contains for uint64[]
func Contains(val uint64, slice ...uint64) bool {
	for _, v := range slice {
		if val == v {
			return true
		}
	}

	return false
}

// Time functions: to time
func TimeMSK_ToTime(mock ...time.Time) time.Time {
	now := time.Now()

	loc, _ := time.LoadLocation("UTC")

	if len(mock) > 0 {
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

	f, err := os.OpenFile("./logs/"+DateMSK_ToLocaleSepString()+"_dev.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)

	if err != nil {
		log.SetPrefix(TimeMSK_ToLocaleString() + "\n")
		log.Printf("Error opening devlog file: %v\n", err)

		return false
	} else {
		log.SetFlags(0)
		log.SetPrefix(TimeMSK_ToLocaleString() + "\n")
		log.SetOutput(f)

		log.Printf("REQUEST (%s): %s\nFROM: %s\nSTATUS: %d\nMEMORY USAGE (KiB): Alloc = %v; TotalAlloc = %v; Sys = %v; NumGC = %v;\n\n",
			TimeMSK_ToLocaleString(), uri, ip, status,
			memStats.Alloc/1024, memStats.TotalAlloc/1024, memStats.Sys/1024, memStats.NumGC)

		defer func(f *os.File) {
			_ = f.Close()
		}(f)
	}

	return true
}
