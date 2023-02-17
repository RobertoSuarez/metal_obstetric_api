package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

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

	router.Put("/:idcita/recordatorio", cita.HandlerRecordarCita)

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

// generar el recordatorio para la paciente
func (citaController *Cita) HandlerRecordarCita(c *fiber.Ctx) error {
	var cita models.Cita

	idCita := c.Params("idcita")
	if len(idCita) < 1 {
		return c.SendString("falta el id de la cita")
	}

	cita.ID, _ = primitive.ObjectIDFromHex(idCita)

	cita.ObtenerCitaPorID()

	// configuración del correo
	from := "electrosonix12@gmail.com"
	password := "ntcrqetvrfbtsxvq"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Asusto del correo
	subject := "Correo electrónico de prueba"

	// Destinatario del correo electrónico
	to := []string{cita.Paciente.Correo}

	// Cargar la plantilla HTML
	tmpl, err := template.ParseFiles("plantillas/plantilla.html")
	if err != nil {
		return c.SendString("Error: " + err.Error())
	}

	// Generar el contenido del correo electrónico a partir de la plantilla HTML
	var body bytes.Buffer
	err = tmpl.Execute(&body, cita)
	if err != nil {
		return c.SendString("Error: " + err.Error())
	}

	// Autenticarse en el servidor SMTP y enviar el correo electrónico
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Mensaje del correo con el html incrustado
	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html\r\n\r\n%s", to[0], subject, body.String()))

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return c.SendString("Error: " + err.Error())
	}

	// Incrementar el contador de los recordatorios en la base de datos.
	err = cita.IncrementarRecordatorio()
	if err != nil {
		return c.SendString("Error: " + err.Error())
	}

	fmt.Println("Correo electrónico enviado correctamente.")

	return c.SendString("Se envio el correo electronico")
}
