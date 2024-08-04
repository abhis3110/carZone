package main

import (
	"database/sql"
	"fmt"
	"github.com/abhis3110/carZone/driver"
	carHandler "github.com/abhis3110/carZone/handler/car"
	engineHandler "github.com/abhis3110/carZone/handler/engine"
	"github.com/gorilla/mux"
	"net/http"
	"os"

	carService "github.com/abhis3110/carZone/service/car"
	engineService "github.com/abhis3110/carZone/service/engine"
	carStore "github.com/abhis3110/carZone/store/car"
	engineStore "github.com/abhis3110/carZone/store/engine"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	driver.InitDB()
	defer driver.CloseDB()

	db := driver.GetDB()
	carstore := carStore.New(db)
	carservice := carService.NewCarService(carstore)

	enginestore := engineStore.New(db)
	engineservice := engineService.NewEngineService(enginestore)

	carhandler := carHandler.NewCarHandler(carservice)
	enginehandler := engineHandler.NewEngineHandler(engineservice)

	router := mux.NewRouter()

	schemaFile := "store/schema.sql"
	if err := executeSchemaFile(db, schemaFile); err != nil {
		log.Fatal("Error while executing schema file ", err)
	}

	router.HandleFunc("/cars/{id}", carhandler.GetCarByID).Methods("GET")
	router.HandleFunc("/cars", carhandler.GetCarByBrand).Methods("GET")
	router.HandleFunc("/cars", carhandler.GetCarByID).Methods("POST") // need to correct func call
	router.HandleFunc("/cars/{id}", carhandler.UpdateCar).Methods("PUT")
	router.HandleFunc("/cars/{id}", carhandler.DeleteCar).Methods("DELETE")

	router.HandleFunc("/engine/{id}", enginehandler.GetEngineByID).Methods("GET")
	router.HandleFunc("/engine", enginehandler.CreateEngine).Methods("POST")
	router.HandleFunc("/engine/{id}", enginehandler.UpdateEngine).Methods("PUT")
	router.HandleFunc("/engine/{id}", enginehandler.DeleteEngine).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(addr, router))

}

func executeSchemaFile(db *sql.DB, schemaFile string) error {
	sqlFile, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(sqlFile)); err != nil {
		return err
	}
	return nil
}
