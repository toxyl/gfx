package main

import (
	_ "embed"
	"image/png"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/toxyl/gfx/parser"
)

//go:embed index.html
var index string

//go:embed style.css
var style string

//go:embed script.js
var script string

//go:embed gfxs-mode.js
var gfxsmode string

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		c.Context().SetContentType("text/html")
		return c.SendString(index)
	})
	app.Get("/styles.css", func(c *fiber.Ctx) error {
		c.Context().SetContentType("text/css")
		return c.SendString(style)
	})
	app.Get("/script.js", func(c *fiber.Ctx) error {
		c.Context().SetContentType("text/javascript")
		return c.SendString(script)
	})
	app.Get("/gfxs-mode.js", func(c *fiber.Ctx) error {
		c.Context().SetContentType("text/javascript")
		return c.SendString(gfxsmode)
	})
	app.Post("/render", func(c *fiber.Ctx) error {
		gfxs := c.FormValue("gfxs")
		if gfxs == "" {
			return c.Status(400).SendString("GFXS script is required")
		}
		comp, err := parser.ParseComposition(gfxs)
		if err != nil {
			return c.Status(400).SendString("GFXS script could not be parsed: " + err.Error())
		}
		return png.Encode(c.Response().BodyWriter(), comp.Render().Get())
	})

	app.Listen(":8080")
}
