package main

import (
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
)

var runLocal bool

func main() {
	rl, ok := os.LookupEnv("RUN_LOCAL")
	if !ok {
		rl = "false"
	}
	runLocal, err := strconv.ParseBool(rl)
	if err != nil {
		log.Fatal(err)
	}
	if runLocal {
		handler()
	} else {
		lambda.Start(handler)
	}

}
