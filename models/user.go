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
	Celular    string             `bson:"celular,omitempty" json:"celular"`
	Rol        string             `bson:"rol" json:"rol"`
}

func (u *User) getNameCollection() string {
	return "usuarios"
}

func (u *User) Registrar() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	col := DB.Collection(u.getNameCollection())

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

func (u *User) Login() (User, error) {

	// Selecionando la colección
	col := DB.Collection(u.getNameCollection())

	// Recuperando el user de la DB
	var userDB User
	err := col.FindOne(context.TODO(), bson.D{
		{Key: "correo", Value: u.Correo},
	}).Decode(&userDB)

	if err != nil {
		return userDB, newUserError("No existe el usuario", false)
	}

	// comparamos el hash y la contraseña
	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(u.Password))
	if err != nil {
		return userDB, newUserError("La contraseña es incorrecta", true)
	}

	return userDB, nil
}

func (u *User) ObtenerUsuarios(selector bson.M) (usuarios []User, err error) {

	colUsers := DB.Collection(u.getNameCollection())
	cursor, err := colUsers.Find(context.TODO(), selector)
	if err != nil {
		return usuarios, err
	}

	// Iterando el cursor
	for cursor.Next(context.TODO()) {
		var usuario User
		err := cursor.Decode(&usuario)
		if err != nil {
			return usuarios, err
		}

		usuarios = append(usuarios, usuario)
	}

	if err := cursor.Err(); err != nil {
		return usuarios, err
	}

	return usuarios, nil
}
