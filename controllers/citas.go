package controllers

import (
	"fmt"

	"github.com/RobertoSuarez/api_metal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cita struct{}

func NewControllerCitas() *Cita {
	return &Cita{}
}

func (cita *Cita) ConfigPath(router *fiber.App) *fiber.App {
	router.Get("/", cita.HandlerObtenerCitas)
	router.Post("/", cita.HandlerRegistrarCitas)
	router.Put("/", cita.HandlerActualizarCita)

	router.Get("/:idcita", cita.HandlerObtenerCitaPorID)

	return router
}

func (citaController *Cita) HandlerObtenerCitas(c *fiber.Ctx) error {
	selector := bson.M{}
	var cita models.Cita

	//para buscar todas las citas por la cédula que coincida
	cedula := c.Query("cedula")
	if len(cedula) > 0 {
		selector["paciente.cedula"] = cedula
	}

	idDoctor := c.Query("id-doctor")
	fmt.Println("ID Doctor: ", idDoctor)
	if len(idDoctor) > 0 {
		objId, _ := primitive.ObjectIDFromHex(idDoctor)
		selector["doctor._id"] = objId
	}

	//cuando ya esté armado el filtro
	citas, err := cita.ObtenerCitas(selector)
	if err != nil {
		return c.JSON(err)
	}

	return c.JSON(citas)
}

func (citaController *Cita) HandlerRegistrarCitas(c *fiber.Ctx) error {
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

func (citaController *Cita) HandlerActualizarCita(c *fiber.Ctx) error {
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

// se va a obtener una cita por el id que se pase por query
func (citaController *Cita) HandlerObtenerCitaPorID(c *fiber.Ctx) error {
	var cita models.Cita

	idCita := c.Params("idcita")
	if len(idCita) < 1 {
		return c.SendString("falta el id del doctor")
	}

	cita.ID, _ = primitive.ObjectIDFromHex(idCita)

	err := cita.ObtenerCitaPorID()
	if err != nil {
		return c.JSON("Error: " + err.Error())
	}

	return c.JSON(cita)
}
