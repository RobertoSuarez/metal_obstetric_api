package models

import (
	"github.com/RobertoSuarez/api_metal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Database

func InitDataBase() {
	DB = db.GetDB()
}
