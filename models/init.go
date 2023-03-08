package models

import (
	"github.com/RobertoSuarez/api_metal/db"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Database
var logger *logrus.Logger

func InitDataBase() {
	DB = db.GetDB()

	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
}
