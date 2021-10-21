package routes

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tarantool/go-tarantool"

	"ivankprod.ru/src/server/modules/models"
	"ivankprod.ru/src/server/modules/utils"
)

//  VK authentication
func authVK(c *fiber.Ctx, db *tarantool.Connection, userExisting *models.User) error {
	query := &utils.URLParams{}

	(*query)["client_id"] = os.Getenv("AUTH_VK_CLIENT_ID")
	(*query)["client_secret"] = os.Getenv("AUTH_VK_CLIENT_SECRET")
	(*query)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
	(*query)["code"] = c.Query("code")

	req := &fiber.Client{NoDefaultUserAgentHeader: false}

	code, res, errs := req.Get("https://oauth.vk.com/access_token" + query.ToString(true)).Bytes()

	if (code != 200 && code != 400) || len(errs) > 0 {
		for _, v := range errs {
			if v != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "ВКонтакте OAuth - "+v.Error())
			}
		}

		return fiber.NewError(fiber.StatusInternalServerError, "ВКонтакте OAuth - "+"код ошибки: "+strconv.Itoa(code))
	}

	ress1 := new(map[string]interface{})
	err := json.Unmarshal(res, ress1)
	if err != nil {
		return err
	}
	if (*ress1)["error"] != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ВКонтакте OAuth - "+((*ress1)["error_description"]).(string))
	}

	userID := strconv.FormatFloat((*ress1)["user_id"].(float64), 'f', 0, 64)

	if (*ress1)["access_token"] != nil {
		query := &utils.URLParams{}

		(*query)["v"] = "5.131"
		(*query)["uids"] = userID
		(*query)["access_token"] = (*ress1)["access_token"].(string)
		(*query)["fields"] = "photo_400_orig"

		code, res, errs := (*req).Get("https://api.vk.com/method/users.get" + (*query).ToString(true)).Bytes()

		if (code != 200 && code != 400) || len(errs) > 0 {
			for _, v := range errs {
				if v != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "ВКонтакте OAuth - "+v.Error())
				}
			}

			return fiber.NewError(fiber.StatusInternalServerError, "ВКонтакте OAuth - "+"код ошибки: "+strconv.Itoa(code))
		}

		ress := new(map[string]interface{})
		err := json.Unmarshal(res, ress)
		if err != nil {
			return err
		}
		if (*ress)["error"] != nil {
			return fiber.NewError(fiber.StatusBadRequest, "ВКонтакте OAuth - "+((*ress)["error_description"]).(string))
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

		id, _, _, err := models.ExistsUser(db, (*user).SocialID, 0)
		if err != nil {
			return err
		}

		if id > 0 {
			if err := models.SignInUser(db, user, id); err != nil {
				return err
			}
		} else {
			if userExisting != nil {
				(*user).Group = (*userExisting).Group
				(*user).Role = (*userExisting).Role
			}

			id, err = models.AddUser(db, user)
			if err != nil {
				return err
			}
		}

		if userExisting == nil {
			c.Cookie(&fiber.Cookie{
				Name:     "session",
				Value:    strconv.FormatUint(id, 10) + ":" + utils.HashSHA512(strconv.FormatUint(id, 10)+userID+(*user).AccessToken+c.Get("user-agent")),
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

//  Facebook authentication
func authFacebook(c *fiber.Ctx, db *tarantool.Connection, userExisting *models.User) error {
	query := &utils.URLParams{}

	(*query)["client_id"] = os.Getenv("AUTH_FB_CLIENT_ID")
	(*query)["client_secret"] = os.Getenv("AUTH_FB_CLIENT_SECRET")
	(*query)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
	(*query)["auth_type"] = "rerequest"
	(*query)["code"] = c.Query("code")

	req := &fiber.Client{NoDefaultUserAgentHeader: false}
	code, res, errs := (*req).Get("https://graph.facebook.com/v11.0/oauth/access_token" + (*query).ToString(true)).Bytes()

	if (code != 200 && code != 400) || len(errs) > 0 {
		for _, v := range errs {
			if v != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Facebook OAuth - "+v.Error())
			}
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Facebook OAuth - "+"код ошибки: "+strconv.Itoa(code))
	}

	ress1 := new(map[string]interface{})
	err := json.Unmarshal(res, ress1)
	if err != nil {
		return err
	}
	if (*ress1)["error"] != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Facebook OAuth - "+(((*ress1)["error"]).(map[string]interface{})["message"]).(string))
	}

	if (*ress1)["access_token"] != nil {
		query := &utils.URLParams{}

		(*query)["access_token"] = (*ress1)["access_token"].(string)
		(*query)["fields"] = "id,email,first_name,last_name,picture.width(400)"

		code, res, errs := (*req).Get("https://graph.facebook.com/me" + (*query).ToString(true)).Bytes()

		if (code != 200 && code != 400) || len(errs) > 0 {
			for _, v := range errs {
				if v != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "Facebook OAuth - "+v.Error())
				}
			}

			return fiber.NewError(fiber.StatusInternalServerError, "Facebook OAuth - "+"код ошибки: "+strconv.Itoa(code))
		}

		ress := new(map[string]interface{})
		err := json.Unmarshal(res, ress)
		if err != nil {
			return err
		}
		if (*ress)["error"] != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Facebook OAuth - "+(((*ress)["error"]).(map[string]interface{})["message"]).(string))
		}

		user := &models.User{
			SocialID:    ((*ress)["id"]).(string),
			NameFirst:   ((*ress)["first_name"]).(string),
			NameLast:    ((*ress)["last_name"]).(string),
			AvatarPath:  ((*ress)["picture"]).(map[string]interface{})["data"].(map[string]interface{})["url"].(string),
			Email:       ((*ress)["email"]).(string),
			AccessToken: ((*ress1)["access_token"]).(string),
			Type:        2,
		}

		id, _, _, err := models.ExistsUser(db, (*user).SocialID, 2)
		if err != nil {
			return err
		}

		if id > 0 {
			if err := models.SignInUser(db, user, id); err != nil {
				return err
			}
		} else {
			if userExisting != nil {
				(*user).Group = (*userExisting).Group
				(*user).Role = (*userExisting).Role
			}

			id, err = models.AddUser(db, user)
			if err != nil {
				return err
			}
		}

		if userExisting == nil {
			c.Cookie(&fiber.Cookie{
				Name:     "session",
				Value:    strconv.FormatUint(id, 10) + ":" + utils.HashSHA512(strconv.FormatUint(id, 10)+(*user).SocialID+(*user).AccessToken+c.Get("user-agent")),
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

//  Yandex authentication
func authYandex(c *fiber.Ctx, db *tarantool.Connection, userExisting *models.User) error {
	query := &utils.URLParams{}

	(*query)["client_id"] = os.Getenv("AUTH_YA_CLIENT_ID")
	(*query)["client_secret"] = os.Getenv("AUTH_YA_CLIENT_SECRET")
	(*query)["grant_type"] = "authorization_code"
	(*query)["code"] = c.Query("code")

	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("accept", "application/json")
	req.SetBodyString((*query).ToString(true)[1:])
	req.SetRequestURI("https://oauth.yandex.ru/token")

	if err := a.Parse(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Яндекс OAuth - "+err.Error())
	}

	code, res, errs := a.Bytes()

	if (code != 200 && code != 400) || len(errs) > 0 {
		for _, v := range errs {
			if v != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Яндекс OAuth - "+v.Error())
			}
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Яндекс OAuth - "+"код ошибки: "+strconv.Itoa(code))
	}

	ress1 := new(map[string]interface{})
	err := json.Unmarshal(res, ress1)
	if err != nil {
		return err
	}

	if (*ress1)["error"] != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Яндекс OAuth - "+(*ress1)["error_description"].(string))
	}

	if (*ress1)["access_token"] != nil {
		a = fiber.AcquireAgent()
		req = a.Request()
		req.Header.SetMethod(fiber.MethodPost)
		req.Header.Set("accept", "application/json")
		req.Header.Set("authorization", "OAuth "+(*ress1)["access_token"].(string))
		req.SetRequestURI("https://login.yandex.ru/info?format=json")

		if err := a.Parse(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Яндекс OAuth - "+err.Error())
		}

		code, res, errs := a.Bytes()

		if (code != 200 && code != 400) || len(errs) > 0 {
			for _, v := range errs {
				if v != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "Яндекс OAuth - "+v.Error())
				}
			}

			return fiber.NewError(fiber.StatusInternalServerError, "Яндекс OAuth - "+"код ошибки: "+strconv.Itoa(code))
		}

		ress := new(map[string]interface{})
		if err := json.Unmarshal([]byte(res), ress); err != nil {
			return err
		}
		if (*ress)["error"] != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Яндекс OAuth - "+((*ress)["error_description"]).(string))
		}

		user := &models.User{
			SocialID:    ((*ress)["id"]).(string),
			NameFirst:   ((*ress)["first_name"]).(string),
			NameLast:    ((*ress)["last_name"]).(string),
			AvatarPath:  "https://avatars.yandex.net/get-yapic/" + ((*ress)["default_avatar_id"]).(string) + "/islands-200",
			Email:       ((*ress)["default_email"]).(string),
			AccessToken: ((*ress1)["access_token"]).(string),
			Type:        1,
		}

		id, _, _, err := models.ExistsUser(db, (*user).SocialID, 1)
		if err != nil {
			return err
		}

		if id > 0 {
			if err := models.SignInUser(db, user, id); err != nil {
				return err
			}
		} else {
			if userExisting != nil {
				(*user).Group = (*userExisting).Group
				(*user).Role = (*userExisting).Role
			}

			id, err = models.AddUser(db, user)
			if err != nil {
				return err
			}
		}

		if userExisting == nil {
			c.Cookie(&fiber.Cookie{
				Name:     "session",
				Value:    strconv.FormatUint(id, 10) + ":" + utils.HashSHA512(strconv.FormatUint(id, 10)+(*user).SocialID+(*user).AccessToken+c.Get("user-agent")),
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
func authGoogle(c *fiber.Ctx, db *tarantool.Connection, userExisting *models.User) error {
	query := &utils.URLParams{}

	(*query)["client_id"] = os.Getenv("AUTH_GL_CLIENT_ID")
	(*query)["client_secret"] = os.Getenv("AUTH_GL_CLIENT_SECRET")
	(*query)["grant_type"] = "authorization_code"
	(*query)["redirect_uri"] = "https://" + os.Getenv("SERVER_HOST") + "/auth/"
	(*query)["code"] = c.Query("code")

	req := &fiber.Client{NoDefaultUserAgentHeader: false}
	code, res, errs := (*req).Post("https://oauth2.googleapis.com/token" + (*query).ToString(true)).Bytes()

	if (code != 200 && code != 400) || len(errs) > 0 {
		for _, v := range errs {
			if v != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Google OAuth - "+v.Error())
			}
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Google OAuth - "+"код ошибки: "+strconv.Itoa(code))
	}

	ress1 := new(map[string]interface{})
	err := json.Unmarshal(res, ress1)
	if err != nil {
		return err
	}

	if (*ress1)["error"] != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Google OAuth - "+(*ress1)["error_description"].(string))
	}

	if (*ress1)["access_token"] != nil {
		a := (*req).Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
		a.Request().Header.Set("authorization", "Bearer "+(*ress1)["access_token"].(string))

		code, res, errs := a.Bytes()

		if (code != 200 && code != 400) || len(errs) > 0 {
			for _, v := range errs {
				if v != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "Google OAuth - "+v.Error())
				}
			}

			return fiber.NewError(fiber.StatusInternalServerError, "Google OAuth - "+"код ошибки: "+strconv.Itoa(code))
		}

		ress := new(map[string]interface{})
		err := json.Unmarshal(res, ress)
		if err != nil {
			return err
		}

		if (*ress)["error"] != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Google OAuth - "+(((*ress)["error"]).(map[string]interface{})["message"]).(string))
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

		id, _, _, err := models.ExistsUser(db, (*user).SocialID, 3)
		if err != nil {
			return err
		}

		if id > 0 {
			if err := models.SignInUser(db, user, id); err != nil {
				return err
			}
		} else {
			if userExisting != nil {
				(*user).Group = (*userExisting).Group
				(*user).Role = (*userExisting).Role
			}

			id, err = models.AddUser(db, user)
			if err != nil {
				return err
			}
		}

		if userExisting == nil {
			c.Cookie(&fiber.Cookie{
				Name:     "session",
				Value:    strconv.FormatUint(id, 10) + ":" + utils.HashSHA512(strconv.FormatUint(id, 10)+(*user).SocialID+(*user).AccessToken+c.Get("user-agent")),
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

func RouteAuthIndex(db *tarantool.Connection) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uAuth, ok := c.Locals("user_auth").(*models.User)
		if !ok {
			uAuth = nil
		}

		data := make(fiber.Map)
		title := "Авторизация"

		if c.Query("code") != "" && c.Query("state") != "" {
			if c.Query("state") == "vk" {
				if err := authVK(c, db, uAuth); err != nil {
					return err
				}
			} else if c.Query("state") == "facebook" {
				if err := authFacebook(c, db, uAuth); err != nil {
					return err
				}
			} else if c.Query("state") == "yandex" {
				if err := authYandex(c, db, uAuth); err != nil {
					return err
				}
			} else if c.Query("state") == "google" {
				if err := authGoogle(c, db, uAuth); err != nil {
					return err
				}
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
			return c.SendString("<!DOCTYPE html><html><head><script>window.location.href=\"/auth/\"</script></head><body></body></html>")
		} else {
			if uAuth == nil {
				data["links"] = utils.GetAuthLinks()
			} else {
				data["user"] = *uAuth
				data["links"] = utils.GetAuthLinks()
				title = "Личный кабинет"

				userAccounts, err := models.GetUsersGroup(db, uAuth.Group)
				if err != nil {
					return err
				}

				if userAccounts != nil {
					data["user_accounts"] = (*userAccounts).GetCondsByType(uAuth.Type)
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
				"type":  "website",
			},
			"data": data,
		})
		if err == nil {
			if os.Getenv("STAGE_MODE") == "dev" {
				go utils.DevLogger(c.Request().URI().String(), c.IP(), 200)
			}

			return nil
		}

		return fiber.NewError(fiber.StatusNotFound, "Запрашиваемая страница не найдена либо ещё не создана")
	}
}

func RouteAuthLogout(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*models.User)
	if !ok || uAuth == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Для того, чтобы выйти из системы, Вы должны быть авторизованы")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   0,
		Expires:  time.Now().Add(-(time.Hour * 1)),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	if os.Getenv("STAGE_MODE") == "dev" {
		go utils.DevLogger(c.Request().URI().String(), c.IP(), 200)
	}

	return c.Redirect("/auth/", 303)
}
