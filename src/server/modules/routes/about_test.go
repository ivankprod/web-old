package routes

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars"
	"ivankprod.ru/src/server/modules/models"
)

func TestRouteAboutIndex(t *testing.T) {
	engine := handlebars.New("../../views", ".hbs")

	configCommon := fiber.Config{
		Prefork:       false,
		Views:         engine,
		StrictRouting: true,
	}

	config404 := fiber.Config{
		Prefork:       false,
		StrictRouting: true,
	}

	type args struct {
		method     string
		route      string
		config     fiber.Config
		handler    fiber.Handler
		middleware fiber.Handler
	}

	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "About route should return code 200 with logger",
			args: args{
				method:  "GET",
				route:   "/about/",
				config:  configCommon,
				handler: RouteAboutIndex,
				middleware: func(c *fiber.Ctx) error {
					os.Setenv("STAGE_MODE", "dev")

					return c.Next()
				},
			},
			wantCode: 200,
		},
		{
			name: "About route should return code 200 with locals",
			args: args{
				method:  "GET",
				route:   "/about/",
				config:  configCommon,
				handler: RouteAboutIndex,
				middleware: func(c *fiber.Ctx) error {
					c.Locals("user_auth", &models.User{
						ID:             0,
						Group:          0,
						SocialID:       "",
						NameFirst:      "",
						NameLast:       "",
						AvatarPath:     "",
						Email:          "",
						AccessToken:    "",
						LastAccessTime: "",
						Role:           0,
						RoleDesc:       "",
						Type:           0,
						TypeDesc:       "",
					})

					return c.Next()
				},
			},
			wantCode: 200,
		},
		{
			name: "About route should return code 404",
			args: args{
				method:     "GET",
				route:      "/about/",
				config:     config404,
				handler:    RouteAboutIndex,
				middleware: func(c *fiber.Ctx) error { return c.Next() },
			},
			wantCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(tt.args.config)

			app.Add(tt.args.method, tt.args.route, tt.args.middleware, tt.args.handler)

			req := httptest.NewRequest(tt.args.method, tt.args.route, nil)
			resp, err := app.Test(req)

			if err != nil {
				t.Errorf("RouteAboutIndex() error = %v, want no errors", err)
			}

			if resp.StatusCode != tt.wantCode {
				t.Errorf("RouteAboutIndex() status code = %v, wantCode %v", resp.StatusCode, tt.wantCode)
			}
		})
	}
}
