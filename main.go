package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		user := User{
			ID:       1,
			Nombre:   "Roberto",
			Apellido: "Suárez",
		}
		return c.JSON(user)
	})

	app.Listen(":3000")
}

type User struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}
