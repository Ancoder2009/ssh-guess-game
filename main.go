package main

import "os"

func main() {
	port, success := os.LookupEnv("PORT")
	if !success {
		port = "12345"
	}
	app := NewApp("0.0.0.0", port)
	app.Start()
	// p := tea.NewProgram(NewModel("test"))
	// p.Run()
}
