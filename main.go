package main

import (
	"calculationengine/constants"
	"calculationengine/logging"
	"calculationengine/router"
	// "calculationengine/service/evaluator"
	// "calculationengine/service/parser"
	"calculationengine/store"
	// "github.com/gin-contrib/cors"
	// "encoding/json"
	// "fmt"
	// "time"
	// "strings"
	// "text/scanner"
)

func main() {
	logger := logging.Init()
	constants.Load()
	logger.Info("application configuration loaded")

	storage.Connect()
	logger.Info("database connection established")

	if err := storage.AutoMigrate(); err != nil {
		logger.Error("database automigrate failed", "error", err)
		panic(err)
	}
	logger.Info("database migrations completed")

	//^Manual Test golang scanner only
	// var s scanner.Scanner
	// s.Init(strings.NewReader(`IF map, mrp, ^ <> "Hello"`))
	// for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
	// 	fmt.Println(scanner.TokenString(tok))
	// }
	//^Manual Test golang scanner only

	//^ Manual Test Parser
	// start := time.Now()
	// // lexer := parser.NewLexer("100 + 102 + 103 * 9999 + 23987 - 876876 / 3876786")
	// lexer := parser.NewLexer("2 - TRUE")
	// nparser := parser.NewParser(lexer)
	// program := nparser.ParseProgram()
	// eval := evaluator.Eval(program.Statements[0])
	// fmt.Println(eval)
	//^ Manula Test Parser

	router.Api()
	logger.Info("http router initialized", "address", "0.0.0.0:3000")
	router.Router.Run("0.0.0.0:3000")

}
