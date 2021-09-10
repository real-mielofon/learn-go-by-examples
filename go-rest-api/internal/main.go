package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/real-mielofon/learn-go-by-examples/go-rest-api/pkg/swagger/server/restapi"
	"github.com/real-mielofon/learn-go-by-examples/go-rest-api/pkg/swagger/server/restapi/operations"
)

func main() {

	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err)
	}

	api := operations.NewHelloAPIAPI(swaggerSpec)
	server := restapi.NewServer(api)

	defer func() {
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatal(err)
		}
	}()

	server.Port = 8080
	api.CheckHealthHandler = operations.CheckHealthHandlerFunc(Health)
	api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(GetHelloUser)
	api.GetGopherHandler = operations.GetGopherHandlerFunc(GetGopherByName)

	// Start server with listening
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}

// Health route returns Ok
func Health(operations.CheckHealthParams) middleware.Responder {
	return operations.NewCheckHealthOK().WithPayload("OK")
}

//GetHelloUser returns Hello + your name
func GetHelloUser(user operations.GetHelloUserParams) middleware.Responder {
	return operations.NewGetHelloUserOK().WithPayload("Hello " + user.User + "!")
}

// GetGopherByName return a gopher in png
func GetGopherByName(gopher operations.GetGopherParams) middleware.Responder {
	var URL string
	if gopher.Name != nil {
		URL = "https://github.com/scraly/gophers/raw/main/" + gopher.Name* + ".png"
	} else {
		//by default we return dr who gopher
		URL = "https://github.com/scraly/gophers/raw/main/dr-who.png"
	}
	response, err := http.Get(URL)
	if err != nil {
		fmt.Println("error")
	}

	return operations.NewGetGopherNameOK().WithPayload(response.Body)
}
