package config

import "github.com/gofiber/fiber/v2"

type ConfigMicroServicio interface {
	ConfigPath(app *fiber.App) *fiber.App
}

func UseMount(prefix string, r fiber.Router, microServicio ConfigMicroServicio) {
	r.Mount(prefix, microServicio.ConfigPath(fiber.New()))
}
