package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientOptions = options.Client().ApplyURI("mongodb+srv://root:facil2020@cluster0.23imj.mongodb.net/?retryWrites=true&w=majority")

func ConectarDB() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return client
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
		return client
	}
	log.Println("Conexión exitosa")
	return client
}

func GetDB() *mongo.Database {
	client := ConectarDB()
	return client.Database("app_metal_db")
}
