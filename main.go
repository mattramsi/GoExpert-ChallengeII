package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"goexpert-challenge-ii/internal/cep"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <CEP>")
		fmt.Println("Exemplo: go run main.go 01153000")
		os.Exit(1)
	}

	cepInput := strings.TrimSpace(os.Args[1])

	cepInput = strings.ReplaceAll(cepInput, "-", "")
	cepInput = strings.ReplaceAll(cepInput, ".", "")
	cepInput = strings.ReplaceAll(cepInput, " ", "")

	if len(cepInput) != 8 {
		log.Fatalf("CEP inválido: deve conter 8 dígitos. Recebido: %s", cepInput)
	}

	client := cep.NewClient()
	
	fmt.Printf("Buscando CEP: %s\n", cepInput)
	fmt.Println("Aguardando resposta das APIs (timeout: 1 segundo)...")
	fmt.Println()

	start := time.Now()
	address, err := client.SearchCEP(cepInput)
	duration := time.Since(start)

	if err != nil {
		log.Fatalf("Erro ao buscar CEP: %v", err)
	}

	fmt.Println("=== Resultado ===")
	fmt.Printf("CEP: %s\n", address.CEP)
	if address.Street != "" {
		fmt.Printf("Rua/Logradouro: %s\n", address.Street)
	}
	if address.Neighborhood != "" {
		fmt.Printf("Bairro: %s\n", address.Neighborhood)
	}
	if address.City != "" {
		fmt.Printf("Cidade: %s\n", address.City)
	}
	if address.State != "" {
		fmt.Printf("Estado: %s\n", address.State)
	}
	fmt.Printf("\nFonte: %s\n", address.Source)
	fmt.Printf("Tempo de resposta: %v\n", duration)
}
