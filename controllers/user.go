package controllers

import (
	"net/http"

	"github.com/RobertoSuarez/api_metal/jwt"
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
	router.Post("/login", user.login)
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

func (user *User) login(c *fiber.Ctx) error {
	var userData models.User
	err := c.BodyParser(&userData)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.UserError{
			IsUser:  false,
			Message: "No se pudo obtener los datos",
		})
	}

	userDB, err := userData.Login()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.(models.UserError))
	}

	userDB.Password = ""

	token, err := jwt.GenerateJWT(userDB)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.UserError{
			IsUser:  true,
			Message: err.Error(),
		})
	}

	return c.JSON(struct {
		User  models.User `json:"user"`
		Token string      `json:"token"`
	}{
		User:  userDB,
		Token: token,
	})
}
