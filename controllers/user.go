package controllers

import "github.com/gofiber/fiber/v2"

type User struct{}

func NewControllerUser() *User {
	return &User{}
}

func (user *User) ConfigPath(router *fiber.App) *fiber.App {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("todos los usuarios")
	})

	return router
}
