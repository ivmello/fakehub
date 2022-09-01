package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/ivmello/fakehub/docs"
)

func RegisterRoutes(app fiber.Router) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/docs/*", swagger.HandlerDefault)
	v1.Get("/docs/*", swagger.New(swagger.Config{
		URL:          "http://localhost:3000/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		OAuth2RedirectUrl: "http://localhost:3000/docs/oauth2-redirect.html",
	}))

	v1.Get("/alive", Alive)
	v1.Post("/crawl/lupa", Lupa)
}
