package videocontroller

import (
	"errors"
	"net/http"
	"streaming/services/videoservice"
	"streaming/util"
	"strings"
)

const videoDir = "media/video"

// Controller Controls paths for /video path.
type Controller struct {
	Head         string
	Handler      http.Handler
	VideoService *videoservice.VideoService
}

// client for testing video stream
// https://flowplayer.com/developers/tools/stream-tester
// sample url for testing
// http://localhost:8080/video/playback/sample-mp4-file/sample-mp4-file.m3u8

// New Creates http request handlers for paths under /video.
func New(res http.ResponseWriter, req *http.Request) *Controller {
	vc := new(Controller)
	vc.Head, req.URL.Path = util.ShiftPath(req.URL.Path)
	vc.VideoService = videoservice.New()
	switch vc.Head {
	case "playback":
		vc.Handler = http.FileServer(http.Dir(videoDir))
	case "upload":
		vc.Handler = &uploadHandler{}
	}
	return vc
}

type uploadHandler struct{}

// ServeHTTP Serves HTTP requests for /playback/upload.
// Implements http.Handler interface.
func (h *uploadHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(20)
	file, fileHeader, err := req.FormFile("fileToUpload")
	if err != nil {
		util.CreateHTTPResponse(res, "Error retrieving file from request.\n", http.StatusBadRequest, err)
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if strings.Contains(contentType, "video") {
		vs := videoservice.New()
		vs.Add(file, fileHeader)
		if err != nil {
			util.CreateHTTPResponse(res, "Failed to upload video file.\n", http.StatusInternalServerError, err)
		} else {
			util.CreateHTTPResponse(res, "Upload succeeded.\n", http.StatusOK, nil)
		}
	} else {
		err = errors.New("Invalid content-type")
		util.CreateHTTPResponse(res, "Incorrect media type. Please make sure you are uploading a video file.\n", http.StatusBadRequest, err)
	}
}
