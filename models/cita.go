package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Cita struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Paciente      User               `bson:"paciente" json:"paciente"`
	Doctor        User               `bson:"doctor" json:"doctor"`
	Descripcion   string             `bson:"descripcion" json:"descripcion"`
	Fecha         time.Time          `bson:"fecha" json:"fecha"`
	Recordatorios int                `bson:"recordatorios" json:"recordatorios"`

	// Los estaos que podria tener serian
	// Por atender: La cita está pendiente de ser atendida por el profesional médico.
	// Atendida: La cita ha sido atendida y completada.
	// Cancelada: La cita ha sido cancelada por el paciente o por el profesional médico.
	// Reagendada: La cita ha sido reagendada para una fecha y hora diferentes.
	// No asistida: La cita no ha sido atendida por el paciente.
	Estado             string     `bson:"estado" json:"estado"`
	DetallesEmbarazada Embarazada `bson:"detalles_embarazada" json:"detalles_embarazada"`
	DetallesFeto       Feto       `bson:"detalles_feto" json:"detalles_feto"`
}

type Embarazada struct {
	Peso   string `bson:"peso" json:"peso"`
	Altura string `bson:"altura" json:"altura"`
}

type Feto struct {
	Fotos  []string `bson:"fotos" json:"fotos"`
	Tamaño string   `bson:"tamaño" json:"tamaño"`
}

func (c *Cita) RegistrarCita() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	citas := DB.Collection("citas")
	usuario := DB.Collection("usuarios")

	// rellena de datos en usuario
	err := usuario.FindOne(context.TODO(), bson.D{{"_id", c.Paciente.ID}}).Decode(&c.Paciente)
	if err != nil {
		return err
	}
	c.Paciente.Password = ""
	err = usuario.FindOne(context.TODO(), bson.D{{"_id", c.Doctor.ID}}).Decode(&c.Doctor)
	if err != nil {
		return err
	}
	c.Doctor.Password = ""

	// establecemos el estado de la cita
	c.Estado = "Por atender"

	result, err := citas.InsertOne(ctx, c)
	if err != nil {
		return err
	}

	objID, _ := result.InsertedID.(primitive.ObjectID)
	c.ID = objID
	return nil
}

func (c *Cita) ObtenerCitas(selector bson.M) ([]Cita, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	col := DB.Collection("citas")
	var citas []Cita

	opciones := options.Find().SetSort(bson.M{"fecha": 1})

	cur, err := col.Find(ctx, selector, opciones)
	if err != nil {
		return citas, err
	}

	for cur.Next(context.TODO()) {
		var result Cita
		err := cur.Decode(&result)
		if err != nil {
			return citas, err
		}

		citas = append(citas, result)
	}

	if err := cur.Err(); err != nil {
		return citas, err
	}

	return citas, err
}

// Actualizar cita/consulta
func (c *Cita) ActualizarCita() error {

	citas := DB.Collection("citas")
	fmt.Println(c)
	update := bson.D{{
		Key: "$set", Value: c,
	}}
	result, err := citas.UpdateByID(context.TODO(), c.ID, update)
	if err != nil {
		return err
	}

	fmt.Println(err, result)

	return nil
}

// obtener una cita por id
func (c *Cita) ObtenerCitaPorID() error {

	citas := DB.Collection("citas")
	return citas.FindOne(context.TODO(), bson.D{{"_id", c.ID}}).Decode(&c)
}

// Incrementar el recordatorio
func (c *Cita) IncrementarRecordatorio() error {
	citas := DB.Collection("citas")

	result, err := citas.UpdateOne(
		context.TODO(),
		bson.M{"_id": c.ID},
		bson.M{"$inc": bson.M{"recordatorio": 1}},
	)
	fmt.Println("Documentos actualizados: ", result.UpsertedCount)

	return err
}
