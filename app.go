package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//type App struct {
//	Router *mux.Router
//	DB     *sql.DB
//}

func openDatabase() *gorm.DB {
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	return db
}

func populateDatabase(db *gorm.DB) {

}

func closeDatabase(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("failed to get DB instance")
	}
	dbSQL.Close()
}

func main() {

	db := openDatabase()

	// Use the db variable to perform database operations

	// Close the database connection when done
	closeDatabase(db)
}
