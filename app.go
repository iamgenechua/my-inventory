package main

import (
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (app *App) Initialize() error {
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=postgres port=5432 sslmode=disable"
	var err error
	app.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)

	return nil
}

func (app *App) PopulateDatabase() error {
	err := app.DB.AutoMigrate(&Product{})
	if err != nil {
		panic("failed to create Product table")
	}

	// clear Product table first
	result := app.DB.Where("1=1").Delete(&Product{})
	if result.Error != nil {
		panic("failed to clear the table")
		return result.Error
	}

	// Add products into Product table
	products := []*Product{
		{ID: 1, Name: "chair", Quantity: 100, Price: 200},
		{ID: 2, Name: "desk", Quantity: 800, Price: 600.00},
	}

	app.DB.Create(&products) // pass pointer of data to create

	return nil
}

func (app *App) CloseDatabase() error {
	dbSQL, err := app.DB.DB()
	if err != nil {
		panic("failed to get DB instance")
		return err
	}
	dbSQL.Close()

	return nil
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router)) // log errors, if any
}

//func (app *App) handleRoutes() {
//	app.Router.HandleFunc("/products", getProducts).Methods("GET")
//}
