package main

import (
	"log"

	"vintage-vision-api/infrastructure/db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	db.Connect()

	log.Println("🟢 Servidor iniciado")

}
