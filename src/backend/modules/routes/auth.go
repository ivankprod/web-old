package routes

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"ivankprod.ru/src/backend/modules/models"
	"ivankprod.ru/src/backend/modules/utils"
)

//  VK authentication
func authVK(c *fiber.Ctx, userExisting *models.User) error {
	query := &utils.URLParams{}

	(*query)["client_id"] = os.Getenv("AUTH_VK_CLIENT_ID")
	(*query)["client_secret"] = os.Getenv("AUTH_VK_CLIENT_SECRET")
	(*query)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
	(*query)["code"] = c.Query("code")

	req := &fiber.Client{NoDefaultUserAgentHeader: false}
	code, res, errs := (*req).Get("https://oauth.vk.com/access_token" + (*query).ToString(true)).String()
	if code != 200 {
		return fiber.NewError(fiber.StatusInternalServerError, "аутентификация не выполнена!")
	}
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

	if (*ress1)["access_token"] != nil {
		query := &utils.URLParams{}

		(*query)["v"] = "5.52"
		(*query)["uids"] = userID
		(*query)["access_token"] = (*ress1)["access_token"].(string)
		(*query)["fields"] = "photo_400_orig"

		_, res, errs := (*req).Get("https://api.vk.com/method/users.get" + (*query).ToString(true)).String()
		for _, v := range errs {
			if v != nil {
				return v
			}
		}

		ress := new(map[string]interface{})
		err := json.Unmarshal([]byte(res), ress)
		if err != nil {
			return err
		}
		if (*ress)["error"] != nil {
			return fiber.NewError(fiber.StatusInternalServerError, ((*ress)["error_description"]).(string))
		}

		resp := (((*ress)["response"]).([]interface{})[0]).(map[string]interface{})
		user := &models.User{
			SocialID:    userID,
			NameFirst:   (resp["first_name"]).(string),
			NameLast:    (resp["last_name"]).(string),
			AvatarPath:  (resp["photo_400_orig"]).(string),
			Email:       ((*ress1)["email"]).(string),
			AccessToken: ((*ress1)["access_token"]).(string),
			Type:        0,
		}

		id, _, _, err := models.ExistsUser((*user).SocialID, 0)
		if err != nil {
			return err
		}

		if id > 0 {
			if err := models.SignInUser(user); err != nil {
				return err
			}
		} else {
			if userExisting != nil {
				(*user).Group = (*userExisting).Group
				(*user).Role = (*userExisting).Role
			}

			id, err = models.AddUser(user)
			if err != nil {
				return err
			}
		}

		if userExisting == nil {
			c.Cookie(&fiber.Cookie{
				Name:     "session",
				Value:    strconv.FormatInt(id, 10) + ":" + utils.HashSHA512(strconv.FormatInt(id, 10)+userID+(*user).AccessToken+c.Get("user-agent")),
				Path:     "/",
				MaxAge:   86400 * 7,
				Expires:  time.Now().Add(time.Hour * 168),
				Secure:   true,
				HTTPOnly: true,
				SameSite: "Lax",
			})
		}
	}

	return nil
}

//  Google authentication
func authGoogle(c *fiber.Ctx, userExisting *models.User) error {
	query := &utils.URLParams{}

	(*query)["client_id"] = os.Getenv("AUTH_GL_CLIENT_ID")
	(*query)["client_secret"] = os.Getenv("AUTH_GL_CLIENT_SECRET")
	(*query)["grant_type"] = "authorization_code"
	(*query)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
	(*query)["code"] = c.Query("code")

	req := &fiber.Client{NoDefaultUserAgentHeader: false}
	code, res, errs := (*req).Post("https://oauth2.googleapis.com/token" + (*query).ToString(true)).String()
	if code != 200 {
		return fiber.NewError(fiber.StatusInternalServerError, "Аутентификация не выполнена!")
	}
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
		return fiber.NewError(fiber.StatusInternalServerError, "Аутентификация не выполнена!")
	}

	if (*ress1)["access_token"] != nil {
		query := &utils.URLParams{}

		(*query)["access_token"] = (*ress1)["access_token"].(string)
		(*query)["alt"] = "json"

		code, res, errs := (*req).Get("https://www.googleapis.com/oauth2/v1/userinfo" + (*query).ToString(true)).String()
		if code != 200 {
			return fiber.NewError(fiber.StatusInternalServerError, "Невозможно получить данные!")
		}
		for _, v := range errs {
			if v != nil {
				return v
			}
		}

		ress := new(map[string]interface{})
		err := json.Unmarshal([]byte(res), ress)
		if err != nil {
			return err
		}
		if (*ress)["error"] != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Невозможно получить данные!")
		}

		user := &models.User{
			SocialID:    ((*ress)["id"]).(string),
			NameFirst:   ((*ress)["given_name"]).(string),
			NameLast:    ((*ress)["family_name"]).(string),
			AvatarPath:  strings.ReplaceAll(((*ress)["picture"]).(string), "=s96-c", "=s400-c"),
			Email:       ((*ress)["email"]).(string),
			AccessToken: ((*ress1)["access_token"]).(string),
			Type:        3,
		}

		id, _, _, err := models.ExistsUser((*user).SocialID, 3)
		if err != nil {
			return err
		}

		if id > 0 {
			if err := models.SignInUser(user); err != nil {
				return err
			}
		} else {
			if userExisting != nil {
				(*user).Group = (*userExisting).Group
				(*user).Role = (*userExisting).Role
			}

			id, err = models.AddUser(user)
			if err != nil {
				return err
			}
		}

		if userExisting == nil {
			c.Cookie(&fiber.Cookie{
				Name:     "session",
				Value:    strconv.FormatInt(id, 10) + ":" + utils.HashSHA512(strconv.FormatInt(id, 10)+(*user).SocialID+(*user).AccessToken+c.Get("user-agent")),
				Path:     "/",
				MaxAge:   86400 * 7,
				Expires:  time.Now().Add(time.Hour * 168),
				Secure:   true,
				HTTPOnly: true,
				SameSite: "Lax",
			})
		}
	}

	return nil
}

func RouteAuthIndex(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*models.User)
	if !ok {
		uAuth = nil
	}

	data := make(fiber.Map)
	title := "Авторизация"

	if c.Query("code") != "" && c.Query("state") != "" {
		//if uAuth == nil {
		if c.Query("state") == "vk" {
			if err := authVK(c, uAuth); err != nil {
				return err
			}
		} else if c.Query("state") == "google" {
			if err := authGoogle(c, uAuth); err != nil {
				return err
			}
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString("<!DOCTYPE html><html><head><script>window.location.href=\"/auth/\"</script></head><body></body></html>")
		//}
	} else {
		if uAuth == nil {
			data["links"] = utils.GetAuthLinks()
		} else {
			data["user"] = *uAuth
			data["links"] = utils.GetAuthLinks()
			title = "Личный кабинет"

			userAccounts, err := models.GetUsersGroup((*uAuth).Group, (*uAuth).ID)
			if err != nil {
				return err
			}

			if userAccounts != nil {
				data["user_accounts"] = *userAccounts
			}
		}
	}

	err := c.Render("auth", fiber.Map{
		"urlBase":      c.BaseURL(),
		"urlCanonical": c.BaseURL() + c.Path(),
		"pageTitle":    title + " - " + os.Getenv("INFO_TITLE_BASE"),
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

func RouteAuthLogout(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*models.User)
	if !ok || uAuth == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Вы не авторизованы!")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   0,
		Expires:  time.Now().Add(-(time.Hour * 1)),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "strict",
	})

	return c.Redirect("/auth/", 303)
}
