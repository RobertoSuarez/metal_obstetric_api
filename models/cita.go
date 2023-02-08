package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cita struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Usuario     User               `bson:"usuario" json:"usuario"`
	Descripcion string             `bson:"descripcion" json:"descripcion"`
	Fecha       time.Time          `bson:"fecha" json:"fecha"`
}

func (c *Cita) Registrar() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	col := DB.Collection("citas")
	usuario := DB.Collection("usuarios")

	// rellena de datos en usuario
	usuario.FindOne(context.TODO(), bson.D{{"_id", c.Usuario.ID}}).Decode(&c.Usuario)
	c.Usuario.Password = ""

	result, err := col.InsertOne(ctx, c)
	if err != nil {
		return err
	}

	objID, _ := result.InsertedID.(primitive.ObjectID)
	c.ID = objID
	return nil
}

func (c *Cita) ObtenerCitas() ([]Cita, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	col := DB.Collection("citas")
	var citas []Cita
	filtro := bson.D{}
	cur, err := col.Find(ctx, filtro)
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
