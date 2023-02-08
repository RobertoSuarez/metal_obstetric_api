package controllers

import (
	"github.com/RobertoSuarez/api_metal/models"
	"github.com/gofiber/fiber/v2"
)

type Cita struct{}

func NewControllerCitas() *Cita {
	return &Cita{}
}

func (cita *Cita) ConfigPath(router *fiber.App) *fiber.App {
	router.Get("/", HandlerObtenerCitas)
	router.Post("/", HandlerRegistrarCitas)

	return router
}

func HandlerObtenerCitas(c *fiber.Ctx) error {
	var cita models.Cita

	citas, err := cita.ObtenerCitas()
	if err != nil {
		return c.JSON(err)
	}

	return c.JSON(citas)
}

func HandlerRegistrarCitas(c *fiber.Ctx) error {
	var citaData models.Cita
	// Parseamos los datos
	err := c.BodyParser(&citaData)
	if err != nil {
		return c.JSON(err)
	}

	// Registramos en MongoDB
	err = citaData.Registrar()
	if err != nil {
		return c.JSON(err)
	}

	// Retornamos los datos al cliente
	return c.JSON(citaData)
}
