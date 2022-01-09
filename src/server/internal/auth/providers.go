package auth

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"ivankprod.ru/src/server/internal/domain"
	"ivankprod.ru/src/server/internal/services"
	"ivankprod.ru/src/server/pkg/utils"
)

//  VK authentication
func authVK(service services.UserService, c *fiber.Ctx, userExisting *domain.User) error {
	query := &utils.URLParams{
		"client_id":     os.Getenv("AUTH_VK_CLIENT_ID"),
		"client_secret": os.Getenv("AUTH_VK_CLIENT_SECRET"),
		"redirect_uri":  "https://" + os.Getenv("SERVER_HOST") + "/auth/",
		"code":          c.Query("code"),
	}

	req := &fiber.Client{NoDefaultUserAgentHeader: false}

	_, res, errs := req.Get("https://oauth.vk.com/access_token" + query.ToString(true)).Bytes()

	if len(errs) > 0 {
		for _, v := range errs {
			if v != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "ВКонтакте OAuth - "+v.Error())
			}
		}
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
		query := &utils.URLParams{
			"v":            "5.131",
			"uids":         userID,
			"access_token": (*ress1)["access_token"].(string),
			"fields":       "photo_400_orig",
		}

		_, res, errs := req.Get("https://api.vk.com/method/users.get" + query.ToString(true)).Bytes()

		if len(errs) > 0 {
			for _, v := range errs {
				if v != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "ВКонтакте OAuth - "+v.Error())
				}
			}
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
		user := &domain.User{
			SocialID:    userID,
			NameFirst:   (resp["first_name"]).(string),
			NameLast:    (resp["last_name"]).(string),
			AvatarPath:  (resp["photo_400_orig"]).(string),
			Email:       ((*ress1)["email"]).(string),
			AccessToken: ((*ress1)["access_token"]).(string),
			Type:        0,
		}

		ue, err := service.FindOneBySocialID(&domain.UserFindOneBySocialIDDTO{
			SocialID: user.SocialID,
			Type:     0,
		})
		if err != nil {
			return err
		}

		id := ue.ID

		if id > 0 {
			if _, err := service.SignIn(id, &domain.UserSignInDTO{
				NameFirst:   user.NameFirst,
				NameLast:    user.NameLast,
				AvatarPath:  user.AvatarPath,
				Email:       user.Email,
				AccessToken: user.AccessToken,
			}); err != nil {
				return err
			}
		} else {
			if userExisting != nil {
				user.Group = userExisting.Group
				user.Role = userExisting.Role
			}

			ue, err = service.Create(&domain.UserCreateDTO{
				Group:       user.Group,
				SocialID:    user.SocialID,
				NameFirst:   user.NameFirst,
				NameLast:    user.NameLast,
				AvatarPath:  user.AvatarPath,
				Email:       user.Email,
				AccessToken: user.AccessToken,
				Role:        user.Role,
				Type:        user.Type,
			})
			if err != nil {
				return err
			}

			id = ue.ID
		}

		if userExisting == nil {
			c.Cookie(&fiber.Cookie{
				Name:     "session",
				Value:    strconv.FormatUint(id, 10) + ":" + utils.HashSHA512(strconv.FormatUint(id, 10)+user.SocialID+user.AccessToken+c.Get("user-agent")),
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
func authFacebook(service services.UserService, c *fiber.Ctx, userExisting *domain.User) error {
	query := &utils.URLParams{
		"client_id":     os.Getenv("AUTH_FB_CLIENT_ID"),
		"client_secret": os.Getenv("AUTH_FB_CLIENT_SECRET"),
		"redirect_uri":  "https://" + os.Getenv("SERVER_HOST") + "/auth/",
		"auth_type":     "rerequest",
		"code":          c.Query("code"),
	}

	req := &fiber.Client{NoDefaultUserAgentHeader: false}
	_, res, errs := req.Get("https://graph.facebook.com/v11.0/oauth/access_token" + query.ToString(true)).Bytes()

	if len(errs) > 0 {
		for _, v := range errs {
			if v != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Facebook OAuth - "+v.Error())
			}
		}
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
		query := &utils.URLParams{
			"access_token": (*ress1)["access_token"].(string),
			"fields":       "id,email,first_name,last_name,picture.width(400)",
		}

		_, res, errs := req.Get("https://graph.facebook.com/me" + query.ToString(true)).Bytes()

		if len(errs) > 0 {
			for _, v := range errs {
				if v != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "Facebook OAuth - "+v.Error())
				}
			}
		}

		ress := new(map[string]interface{})
		err := json.Unmarshal(res, ress)
		if err != nil {
			return err
		}
		if (*ress)["error"] != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Facebook OAuth - "+(((*ress)["error"]).(map[string]interface{})["message"]).(string))
		}

		user := &domain.User{
			SocialID:    ((*ress)["id"]).(string),
			NameFirst:   ((*ress)["first_name"]).(string),
			NameLast:    ((*ress)["last_name"]).(string),
			AvatarPath:  ((*ress)["picture"]).(map[string]interface{})["data"].(map[string]interface{})["url"].(string),
			Email:       ((*ress)["email"]).(string),
			AccessToken: ((*ress1)["access_token"]).(string),
			Type:        2,
		}

		ue, err := service.FindOneBySocialID(&domain.UserFindOneBySocialIDDTO{
			SocialID: user.SocialID,
			Type:     2,
		})
		if err != nil {
			return err
		}

		id := ue.ID

		if id > 0 {
			if _, err := service.SignIn(id, &domain.UserSignInDTO{
				NameFirst:   user.NameFirst,
				NameLast:    user.NameLast,
				AvatarPath:  user.AvatarPath,
				Email:       user.Email,
				AccessToken: user.AccessToken,
			}); err != nil {
				return err
			}
		} else {
			if userExisting != nil {
				user.Group = userExisting.Group
				user.Role = userExisting.Role
			}

			ue, err = service.Create(&domain.UserCreateDTO{
				Group:       user.Group,
				SocialID:    user.SocialID,
				NameFirst:   user.NameFirst,
				NameLast:    user.NameLast,
				AvatarPath:  user.AvatarPath,
				Email:       user.Email,
				AccessToken: user.AccessToken,
				Role:        user.Role,
				Type:        user.Type,
			})
			if err != nil {
				return err
			}

			id = ue.ID
		}

		if userExisting == nil {
			c.Cookie(&fiber.Cookie{
				Name:     "session",
				Value:    strconv.FormatUint(id, 10) + ":" + utils.HashSHA512(strconv.FormatUint(id, 10)+user.SocialID+user.AccessToken+c.Get("user-agent")),
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
func authYandex(service services.UserService, c *fiber.Ctx, userExisting *domain.User) error {
	query := &utils.URLParams{
		"client_id":     os.Getenv("AUTH_YA_CLIENT_ID"),
		"client_secret": os.Getenv("AUTH_YA_CLIENT_SECRET"),
		"grant_type":    "authorization_code",
		"code":          c.Query("code"),
	}

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

	_, res, errs := a.Bytes()

	if len(errs) > 0 {
		for _, v := range errs {
			if v != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Яндекс OAuth - "+v.Error())
			}
		}
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

		_, res, errs := a.Bytes()

		if len(errs) > 0 {
			for _, v := range errs {
				if v != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "Яндекс OAuth - "+v.Error())
				}
			}
		}

		ress := new(map[string]interface{})
		if err := json.Unmarshal([]byte(res), ress); err != nil {
			return err
		}
		if (*ress)["error"] != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Яндекс OAuth - "+((*ress)["error_description"]).(string))
		}

		user := &domain.User{
			SocialID:    ((*ress)["id"]).(string),
			NameFirst:   ((*ress)["first_name"]).(string),
			NameLast:    ((*ress)["last_name"]).(string),
			AvatarPath:  "https://avatars.yandex.net/get-yapic/" + ((*ress)["default_avatar_id"]).(string) + "/islands-200",
			Email:       ((*ress)["default_email"]).(string),
			AccessToken: ((*ress1)["access_token"]).(string),
			Type:        1,
		}

		ue, err := service.FindOneBySocialID(&domain.UserFindOneBySocialIDDTO{
			SocialID: user.SocialID,
			Type:     1,
		})
		if err != nil {
			return err
		}

		id := ue.ID

		if id > 0 {
			if _, err := service.SignIn(id, &domain.UserSignInDTO{
				NameFirst:   user.NameFirst,
				NameLast:    user.NameLast,
				AvatarPath:  user.AvatarPath,
				Email:       user.Email,
				AccessToken: user.AccessToken,
			}); err != nil {
				return err
			}
		} else {
			if userExisting != nil {
				user.Group = userExisting.Group
				user.Role = userExisting.Role
			}

			ue, err = service.Create(&domain.UserCreateDTO{
				Group:       user.Group,
				SocialID:    user.SocialID,
				NameFirst:   user.NameFirst,
				NameLast:    user.NameLast,
				AvatarPath:  user.AvatarPath,
				Email:       user.Email,
				AccessToken: user.AccessToken,
				Role:        user.Role,
				Type:        user.Type,
			})
			if err != nil {
				return err
			}

			id = ue.ID
		}

		if userExisting == nil {
			c.Cookie(&fiber.Cookie{
				Name:     "session",
				Value:    strconv.FormatUint(id, 10) + ":" + utils.HashSHA512(strconv.FormatUint(id, 10)+user.SocialID+user.AccessToken+c.Get("user-agent")),
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
func authGoogle(service services.UserService, c *fiber.Ctx, userExisting *domain.User) error {
	query := &utils.URLParams{
		"client_id":     os.Getenv("AUTH_GL_CLIENT_ID"),
		"client_secret": os.Getenv("AUTH_GL_CLIENT_SECRET"),
		"grant_type":    "authorization_code",
		"redirect_uri":  "https://" + os.Getenv("SERVER_HOST") + "/auth/",
		"code":          c.Query("code"),
	}

	req := &fiber.Client{NoDefaultUserAgentHeader: false}
	_, res, errs := req.Post("https://oauth2.googleapis.com/token" + query.ToString(true)).Bytes()

	if len(errs) > 0 {
		for _, v := range errs {
			if v != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Google OAuth - "+v.Error())
			}
		}
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
		a := req.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
		a.Request().Header.Set("authorization", "Bearer "+(*ress1)["access_token"].(string))

		_, res, errs := a.Bytes()

		if len(errs) > 0 {
			for _, v := range errs {
				if v != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "Google OAuth - "+v.Error())
				}
			}
		}

		ress := new(map[string]interface{})
		err := json.Unmarshal(res, ress)
		if err != nil {
			return err
		}

		if (*ress)["error"] != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Google OAuth - "+(((*ress)["error"]).(map[string]interface{})["message"]).(string))
		}

		user := &domain.User{
			SocialID:    ((*ress)["id"]).(string),
			NameFirst:   ((*ress)["given_name"]).(string),
			NameLast:    ((*ress)["family_name"]).(string),
			AvatarPath:  strings.ReplaceAll(((*ress)["picture"]).(string), "=s96-c", "=s400-c"),
			Email:       ((*ress)["email"]).(string),
			AccessToken: ((*ress1)["access_token"]).(string),
			Type:        3,
		}

		ue, err := service.FindOneBySocialID(&domain.UserFindOneBySocialIDDTO{
			SocialID: user.SocialID,
			Type:     3,
		})
		if err != nil {
			return err
		}

		id := ue.ID

		if id > 0 {
			if _, err := service.SignIn(id, &domain.UserSignInDTO{
				NameFirst:   user.NameFirst,
				NameLast:    user.NameLast,
				AvatarPath:  user.AvatarPath,
				Email:       user.Email,
				AccessToken: user.AccessToken,
			}); err != nil {
				return err
			}
		} else {
			if userExisting != nil {
				user.Group = userExisting.Group
				user.Role = userExisting.Role
			}

			ue, err = service.Create(&domain.UserCreateDTO{
				Group:       user.Group,
				SocialID:    user.SocialID,
				NameFirst:   user.NameFirst,
				NameLast:    user.NameLast,
				AvatarPath:  user.AvatarPath,
				Email:       user.Email,
				AccessToken: user.AccessToken,
				Role:        user.Role,
				Type:        user.Type,
			})
			if err != nil {
				return err
			}

			id = ue.ID
		}

		if userExisting == nil {
			c.Cookie(&fiber.Cookie{
				Name:     "session",
				Value:    strconv.FormatUint(id, 10) + ":" + utils.HashSHA512(strconv.FormatUint(id, 10)+user.SocialID+user.AccessToken+c.Get("user-agent")),
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
