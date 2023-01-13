package controllers

import (
	"net/http"

	"github.com/RobertoSuarez/api_metal/models"
	"github.com/gofiber/fiber/v2"
)

type User struct {
}

func NewControllerUser() *User {
	return &User{}
}

func (user *User) ConfigPath(router *fiber.App) *fiber.App {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("todos los usuarios")
	})

	router.Post("/registrar", user.registrarUsurio)
	return router
}

func (user *User) registrarUsurio(c *fiber.Ctx) error {
	var usuarioData models.User
	err := c.BodyParser(&usuarioData)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err = usuarioData.Registrar()

	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	usuarioData.Password = ""

	return c.JSON(usuarioData)

}
