package routes

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"ivankprod.ru/src/backend/modules/models"
	"ivankprod.ru/src/backend/modules/utils"
)

func RouteAuthIndex(c *fiber.Ctx) error {
	uAuth := c.Locals("user_auth")
	data := make(fiber.Map)

	if c.Query("code") != "" {
		if uAuth == nil {
			// auth from vk
			query := &utils.URLParams{}

			(*query)["client_id"] = os.Getenv("AUTH_VK_CLIENT_ID")
			(*query)["client_secret"] = os.Getenv("AUTH_VK_CLIENT_SECRET")
			(*query)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
			(*query)["code"] = c.Query("code")

			req := &fiber.Client{NoDefaultUserAgentHeader: false}
			_, res, errs := (*req).Get("https://oauth.vk.com/access_token" + (*query).ToString()).String()
			for _, v := range errs {
				if v != nil {
					return v
				}
			}

			ress1 := new(map[string]interface{})
			err := json.Unmarshal([]byte(res), ress1)
			if err != nil {
				return err
			}
			if (*ress1)["error"] != nil {
				return fiber.NewError(fiber.StatusInternalServerError, ((*ress1)["error_description"]).(string))
			}

			userID := strconv.FormatFloat((*ress1)["user_id"].(float64), 'f', 0, 64)
			userIDint, err := strconv.Atoi(userID)
			if err != nil {
				return err
			}

			if (*ress1)["access_token"] != nil {
				query := &utils.URLParams{}

				(*query)["v"] = "5.52"
				(*query)["uids"] = userID
				(*query)["access_token"] = (*ress1)["access_token"].(string)
				(*query)["fields"] = "photo_big"

				_, res, errs := (*req).Get("https://api.vk.com/method/users.get" + (*query).ToString()).String()
				for _, v := range errs {
					if v != nil {
						data["error"] = v.Error()
					}
				}

				ress := new(map[string]interface{})
				err := json.Unmarshal([]byte(res), ress)
				if err != nil {
					data["error"] = err.Error()
				}
				if (*ress)["error"] != nil {
					data["error"] = (*ress)["error_description"]
				}

				resp := (((*ress)["response"]).([]interface{})[0]).(map[string]interface{})
				user := &models.User{
					ID:          userIDint,
					NameFirst:   (resp["first_name"]).(string),
					NameLast:    (resp["last_name"]).(string),
					AvatarPath:  (resp["photo_big"]).(string),
					Email:       ((*ress1)["email"]).(string),
					AccessToken: ((*ress1)["access_token"]).(string),
					Type:        0,
				}

				b, err := models.ExistsUser((*user).ID)
				if err != nil {
					return err
				}

				if b {
					if err := models.SignInUser(user); err != nil {
						return err
					}
				} else {
					if err := models.AddUser(user); err != nil {
						return err
					}
				}

				c.Cookie(&fiber.Cookie{
					Name:     "session",
					Value:    userID + ":" + utils.HashSHA512(userID+(*user).AccessToken+c.Get("user-agent")),
					Path:     "/",
					MaxAge:   86400 * 7,
					Expires:  time.Now().Add(time.Hour * 168),
					Secure:   true,
					HTTPOnly: true,
					SameSite: "strict",
				})

				c.Redirect("/auth/", 303)
			}
		}
	} else {
		if uAuth == nil {
			data = utils.GetAuthLinks()
		} else {
			data = fiber.Map{"user": uAuth}
		}
	}

	err := c.Render("auth", fiber.Map{
		"urlCanonical": c.BaseURL() + c.Path(),
		"pageTitle":    "Авторизация - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":     os.Getenv("INFO_DESC_BASE"),
		"ogTags": fiber.Map{
			"title": os.Getenv("INFO_TITLE_BASE"),
		},
		"data": data,
	})
	if err == nil {
		return nil
	}

	return fiber.NewError(fiber.StatusNotFound, "Страница не найдена!")
}
