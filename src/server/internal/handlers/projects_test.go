package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars"

	"ivankprod.ru/src/server/internal/domain"
)

func TestHandlerProjectsIndexView(t *testing.T) {
	middlewareAuth := func(c *fiber.Ctx) error {
		c.Locals("user_auth", &domain.User{
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
	}

	type args struct {
		method     string
		route      string
		routePath  string
		handler    fiber.Handler
		middleware fiber.Handler
	}

	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "Projects route should return code 200",
			args: args{
				method:  "GET",
				route:   "/projects/",
				handler: HandlerProjectsIndex,
			},
			wantCode: 200,
		},
		{
			name: "Projects route should return code 200 with locals",
			args: args{
				method:     "GET",
				route:      "/projects/",
				handler:    HandlerProjectsIndex,
				middleware: middlewareAuth,
			},
			wantCode: 200,
		},
		{
			name: "Projects route should return code 404",
			args: args{
				method:  "GET",
				route:   "/projectsss/",
				handler: HandlerProjectsIndex,
			},
			wantCode: 404,
		},
		{
			name: "Projects view route should return code 200",

			args: args{
				method:    "GET",
				route:     "/projects/it/",
				routePath: ":type/",
				handler:   HandlerProjectsView,
			},
			wantCode: 200,
		},
		{
			name: "Projects view route should return code 200 with locals",

			args: args{
				method:     "GET",
				route:      "/projects/it/",
				routePath:  ":type/",
				handler:    HandlerProjectsView,
				middleware: middlewareAuth,
			},
			wantCode: 200,
		},
		{
			name: "Projects view route should return code 404",

			args: args{
				method:    "GET",
				route:     "/projects/itt/",
				routePath: ":type/",
				handler:   HandlerProjectsView,
			},
			wantCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				Prefork:       true,
				Views:         handlebars.New("../../views", ".hbs"),
				StrictRouting: true,
			})

			if tt.args.middleware != nil {
				app.Add(tt.args.method, "/projects/"+tt.args.routePath, tt.args.middleware, tt.args.handler)
			} else {
				app.Add(tt.args.method, "/projects/"+tt.args.routePath, tt.args.handler)
			}

			req := httptest.NewRequest(tt.args.method, tt.args.route, nil)
			resp, err := app.Test(req)

			if err != nil {
				t.Errorf("HandlerProjectsIndexView() error = %v, want no errors", err)
			}

			if resp.StatusCode != tt.wantCode {
				t.Errorf("HandlerProjectsIndexView() status code = %v, wantCode %v", resp.StatusCode, tt.wantCode)
			}
		})
	}
}
