package main

import (
	"fmt"
	"log"

	"github.com/RobertoSuarez/api_metal/config"
	"github.com/RobertoSuarez/api_metal/controllers"
	"github.com/RobertoSuarez/api_metal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	configvar := viper.New()

	configvar.AddConfigPath(".")
	configvar.SetConfigName("app")
	configvar.SetConfigType("env")

	configvar.AutomaticEnv()

	if err := configvar.ReadInConfig(); err != nil {
		fmt.Println("Error al leer las variables de configuración")
		log.Println(err)
	} else {
		fmt.Println("Las variables se establecierón correctamente")
	}

	models.InitDataBase()

	fmt.Println(configvar.GetString("configapp"))
	fmt.Println(configvar.GetString("MONGODB"))

	api := app.Group("/api/v1")

	config.UseMount("/users", api, controllers.NewControllerUser())
	config.UseMount("/citas", api, controllers.NewControllerCitas())

	app.Listen(":" + configvar.GetString("port"))
}
