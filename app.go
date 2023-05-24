package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
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
	app.handleRoutes()
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

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload) // serialize payload into json format
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode) // sets HTTP status code of the response
	w.Write(response)         // writes to the httpResponseWriter (sends response to client)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	errorMessage := map[string]string{"error": err}
	sendResponse(w, statusCode, errorMessage)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, products)
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product := Product{}
	result := app.DB.First(&product, id)

	if result.Error != nil {
		sendError(w, http.StatusNotFound, "Product not found")
		return
	}

	sendResponse(w, http.StatusOK, product)
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")

}

func getProducts(db *gorm.DB) ([]Product, error) {
	products := make([]Product, 0)
	result := db.Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}
