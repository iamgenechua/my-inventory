package main

func main() {
	app := App{}
	app.Initialize()
	app.PopulateDatabase()
	app.Run("localhost:10000")
}
