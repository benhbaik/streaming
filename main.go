package main

import (
	"fmt"
	"log"
	"net/http"
	"streaming/controllers/audiocontroller"
	"streaming/controllers/videocontroller"
	"streaming/util"

	_ "github.com/go-sql-driver/mysql"
)

const port = ":8080"

func main() {
	fmt.Println("Starting streaming app...")
	a := &streamingApp{}

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, a))
}

type streamingApp struct{}

func (streamApp *streamingApp) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = util.ShiftPath(req.URL.Path)

	util.AddCORSHeaders(res)

	switch head {
	case "audio":
		audioController := audiocontroller.NewAudioController(res, req)
		audioController.ServeHTTP(res, req)
	case "video":
		vc := videocontroller.New(res, req)
		vc.Handler.ServeHTTP(res, req)
	default:
		http.Error(res, "Not found.\n", http.StatusNotFound)
	}
}
