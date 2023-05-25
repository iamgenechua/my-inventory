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

	// Add products into Product table
	products := []*Product{
		{Name: "chair", Quantity: 100, Price: 200},
		{Name: "desk", Quantity: 800, Price: 600.00},
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

	// clear Product table
	result := app.DB.Where("1=1").Delete(&Product{})
	if result.Error != nil {
		panic("failed to clear the table")
		return result.Error
	}

	dbSQL.Close()

	log.Println("closed")

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

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	product := Product{}
	err := json.NewDecoder(r.Body).Decode(&product) // decode request body into product struct

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = product.createProduct(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusCreated, product)
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}

func getProducts(db *gorm.DB) ([]Product, error) {
	products := make([]Product, 0)
	result := db.Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (p *Product) createProduct(db *gorm.DB) error {
	result := db.Select("Name", "Quantity", "Price").Create(&p)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (app *App) updateProduct(writer http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(writer, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product := Product{}
	result := app.DB.First(&product, id)

	if result.Error != nil {
		sendError(writer, http.StatusNotFound, "Product not found")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&product) // decode request body into product struct

	if err != nil {
		sendError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}

	result = app.DB.Save(&product)

	if result.Error != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(writer, http.StatusOK, product)
}

func (app *App) deleteProduct(writer http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(writer, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product := Product{}
	result := app.DB.First(&product, id)

	if result.Error != nil {
		sendError(writer, http.StatusNotFound, "Product not found")
		return
	}

	result = app.DB.Delete(&product)

	if result.Error != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(writer, http.StatusOK, product)
}
