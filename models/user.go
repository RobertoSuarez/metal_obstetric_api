package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nombres    string             `bson:"nombres,omitempty" json:"nombres"`
	Apellidos  string             `bson:"apellidos,omitempty" json:"apellidos"`
	Cedula     string             `bson:"cedula,omitempty" json:"cedula"`
	Correo     string             `bson:"correo,omitempty" json:"correo"`
	Password   string             `bson:"password,omitempty" json:"password"`
	Genero     string             `bson:"genero,omitempty" json:"genero"`
	Nacimiento time.Time          `bson:"nacimiento,omitempty" json:"nacimiento"`
}

func (u *User) Registrar() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	col := DB.Collection("usuarios")

	// verificar que el usuario no exita+

	var userExiste User
	col.FindOne(context.TODO(), bson.D{{"correo", u.Correo}}).Decode(&userExiste)

	if len(userExiste.Correo) > 0 {
		return errors.New("el usuario ya existe")
	}

	err := u.encryptPassword()
	if err != nil {
		return err
	}

	result, err := col.InsertOne(ctx, u)
	if err != nil {
		return err
	}

	objID, _ := result.InsertedID.(primitive.ObjectID)
	u.ID = objID
	return nil
}

// Encriptar la contraseña del usuario
func (u *User) encryptPassword() error {
	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), cost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}
