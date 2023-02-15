package controllers

import (
	"github.com/RobertoSuarez/api_metal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type Cita struct{}

func NewControllerCitas() *Cita {
	return &Cita{}
}

func (cita *Cita) ConfigPath(router *fiber.App) *fiber.App {
	router.Get("/", cita.HandlerObtenerCitas)
	router.Post("/", HandlerRegistrarCitas)
	router.Put("/", HandlerActualizarCita)

	return router
}

func (citaa *Cita) HandlerObtenerCitas(c *fiber.Ctx) error {
	selector := bson.M{}
	var cita models.Cita

	//para buscar todas las citas por la cédula que coincida
	cedula := c.Query("cedula")
	if len(cedula) > 0 {
		selector["paciente.cedula"] = cedula
	}

	//cuando ya esté armado el filtro
	citas, err := cita.ObtenerCitas(selector)
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
