package main

import (
	"fmt"
	"log"
	"net/http"
	"streaming/controllers/audiocontroller"
	"streaming/controllers/videocontroller"
	"streaming/util"
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
		audioHandler := audiocontroller.HandleAudio(res, req)
		audioHandler.ServeHTTP(res, req)
	case "video":
		videoHandler := videocontroller.HandleVideo(res, req)
		videoHandler.ServeHTTP(res, req)
	default:
		http.Error(res, "Not found.\n", http.StatusNotFound)
	}
}
