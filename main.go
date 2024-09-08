package main

import (
	"github.com/vivekmv23/go-web-frameworks/database"
	wfgorillamux "github.com/vivekmv23/go-web-frameworks/wf-gorilla-mux"
	wfstandardlib "github.com/vivekmv23/go-web-frameworks/wf-standard-lib"
)

func main() {
	StartGorillaMuxServer()
}

func StartStdLibServer() {
	standardLibWebServer := wfstandardlib.StandardLibWebServer{}
	standardLibWebServer.Start(8080)
}

func StartGorillaMuxServer() {
	d := database.NewDatabase()
	gorillamux := wfgorillamux.NewGorillaMuxWebServer(d)
	gorillamux.Start(8080)
}
