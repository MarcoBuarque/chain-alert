package main

import (
	"fmt"

	"github.com/MarcoBuarque/chain-alert/chain-alert/internal/listener"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if err := listener.Run(); err != nil {
		fmt.Println("OOOOPS ;-;", err)
	}
}
