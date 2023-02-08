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
	router.Put("/", HandlerActualizarCita)

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

	if len(citaData.Doctor.ID) < 1 {
		return c.SendString("Falta el doctor")
	}

	if len(citaData.Paciente.ID) < 1 {
		return c.SendString("Falta el paciente")
	}

	// Registramos en MongoDB
	err = citaData.RegistrarCita()
	if err != nil {
		return c.JSON(err)
	}

	// Retornamos los datos al cliente
	return c.JSON(citaData)
}

func HandlerActualizarCita(c *fiber.Ctx) error {
	var cita models.Cita
	err := c.BodyParser(&cita)
	if err != nil {
		return c.JSON("Error: " + err.Error())
	}

	err = cita.ActualizarCita()
	if err != nil {
		return c.JSON("Error: " + err.Error())
	}

	return c.JSON(cita)
}
