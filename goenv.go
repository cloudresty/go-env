package goenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Var(key string) string {

	//
	// Local .env File
	//

	/*

		Checking if there's any '.env' file stored locally.
		This should not be present in a Production environment,
		it should be used only within a local / development environment.

		Make sure that .gitignore contains ".env" so we won't
		push this file to the code repository

	*/

	//
	// Local .env File
	//

	localEnvFile := ".env"
	godotenv.Load(localEnvFile)

	//
	// System Environment Variables
	//

	value, found := os.LookupEnv(key)

	// Check if the environment variable has been set up
	if !found {

		log.Printf("ERROR | Couldn't find '%s' environment variable", key)
		os.Exit(1)

	}

	// Check if the environment variable is empty or not
	if len(key) < 1 {

		log.Printf("ERROR | Environment variable '%s' is empty, can't find any value attached to it", key)
		os.Exit(1)

	}

	return value

}
