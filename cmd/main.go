package main

import (
	"api-employees-and-departaments/config"
	"fmt"
)

func main() {
	fmt.Println("Iniciando projeto")

	config.GlobalConfig.LoadVariables()
}
