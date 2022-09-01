package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ivmello/fakehub/pkg/crawler"
)

func Alive(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"data": "alive",
	})
}

type LupaPayload struct {
	Query string `json:"query" xml:"query" form:"query"`
}

// ShowAccount godoc
// @Summary      Agencia Lupa crawler
// @Description  search on Agencia Lupa query param sent on body
// @Tags         crawl/lupa
// @Accept       json
// @Produce      json
// @Param        query body LupaPayload  true  "Query"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} error
// @Failure      404 {object} error
// @Failure      500 {object} error
// @Router       /crawl/lupa [post]
func Lupa(c *fiber.Ctx) error {
	payload := new(LupaPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	result, err := crawler.CrawlLupa(payload.Query)

	if err != nil {
		return c.JSON(map[string]interface{}{
			"error": err,
		})
	}

	return c.JSON(result)
}
