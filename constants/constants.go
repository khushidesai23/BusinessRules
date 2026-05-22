package constants

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"slices"
)

const (
	SUCCESS = "success"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBPort     string
	DBName     string
	SSLMode    string
}

var AppConfig *Config

var AcceptedTokens []string
var ParanthesisTokens []string
var OperatorTokens []string
var AttributesAndConstantTokens []string

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = &Config{
		DBHost:     os.Getenv("PGHOST"),
		DBUser:     os.Getenv("PGUSER"),
		DBPassword: os.Getenv("PGPASSWORD"),
		DBName:     os.Getenv("PGDATABASE"),
		SSLMode:    os.Getenv("PGSSLMODE"),
		DBPort:     os.Getenv("PGPORT"),
	}

	ParanthesisTokens = []string{"(", ")"}
	OperatorTokens = []string{"+", "/", "-", "*"}
	AttributesAndConstantTokens = []string{"Ident", "Int", "Float", "String"}
	AcceptedTokens = slices.Concat(ParanthesisTokens, OperatorTokens, AttributesAndConstantTokens)
}
