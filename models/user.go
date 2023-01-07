package models

import "fmt"

type User struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

func (u User) Saludar() {
	fmt.Println(u.Nombre)
}
