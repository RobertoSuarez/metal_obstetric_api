package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RobertoSuarez/api_metal/jwt"
	"github.com/RobertoSuarez/api_metal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
}

func NewControllerUser() *User {
	return &User{}
}

func (user *User) ConfigPath(router *fiber.App) *fiber.App {
	router.Get("/", user.HandlerObtenerUsuarios)

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

// Query rol
func (user *User) HandlerObtenerUsuarios(c *fiber.Ctx) error {

	selector := bson.M{}
	userData := models.User{}

	// el cliente puede filtrar por el rol del usuario
	rol := c.Query("rol")
	if len(rol) > 0 {
		selector["rol"] = rol
	}

	// el cliente puede filtrar por algun termino los documentos
	termino := c.Query("termino")
	if len(termino) > 0 {
		selector["$or"] = []bson.M{
			bson.M{"nombres": bson.M{"$regex": fmt.Sprintf(".*%s.*", termino), "$options": "i"}},
			bson.M{"apellidos": bson.M{"$regex": fmt.Sprintf(".*%s.*", termino), "$options": "i"}},
			bson.M{"cedula": bson.M{"$regex": fmt.Sprintf(".*%s.*", termino), "$options": "i"}},
			bson.M{"correo": bson.M{"$regex": fmt.Sprintf(".*%s.*", termino), "$options": "i"}},
		}
	}

	limiteStrimg := c.Query("limite")
	limite := 0
	if len(limiteStrimg) > 0 {
		limite, _ = strconv.Atoi(limiteStrimg)
	}

	usuarios, err := userData.ObtenerUsuarios(selector, limite)
	if err != nil {
		return err
	}
	return c.JSON(usuarios)
}
