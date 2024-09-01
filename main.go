package main

import (
	wfstandardlib "github.com/vivekmv23/go-web-frameworks/wf-standard-lib"
)

func main() {
	standardLibWebServer := wfstandardlib.StandardLibWebServer{}
	standardLibWebServer.Start(8080)
}
