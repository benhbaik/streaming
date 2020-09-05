package videocontroller

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"streaming/util"
	"strings"
)

const videoDir = "media/video"
const tempVideoStore = "tempVideoStore"

// client for testing video stream
// https://flowplayer.com/developers/tools/stream-tester

// ffmpeg command
// ffmpeg -i bunnyVideo.mp4 -codec: copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls bunnyVideo.m3u8

// HandleVideo handles http requests and returns http handler for /video
func HandleVideo(res http.ResponseWriter, req *http.Request) http.Handler {
	var head string
	var handler http.Handler

	head, req.URL.Path = util.ShiftPath(req.URL.Path)

	switch head {
	// sample url for testing
	// http://localhost:8080/video/playback/upload-630824529-bunnyVideo/upload-630824529-bunnyVideo.m3u8
	case "playback":
		handler = handlePlayback(res, req)
	case "upload":
		handler = &uploadHandler{}
	}

	return handler
}

func handlePlayback(res http.ResponseWriter, req *http.Request) http.Handler {
	return http.FileServer(http.Dir(videoDir))
}

type uploadHandler struct{}

// ServeHTTP serves HTTP requests for /playback/upload
func (h *uploadHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(20)
	file, fileHeader, err := req.FormFile("fileToUpload")
	if err != nil {
		util.CreateHTTPResponse(res, "Error retrieving file from request.\n", http.StatusBadRequest, err)
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if strings.Contains(contentType, "video") {
		err := uploadVideo(file, fileHeader)
		if err != nil {
			util.CreateHTTPResponse(res, "Failed to upload video file.\n", http.StatusInternalServerError, err)
		} else {
			util.CreateHTTPResponse(res, "Upload succeeded.\n", http.StatusOK, nil)
		}
	} else {
		util.CreateHTTPResponse(res, "Incorrect media type. Please make sure you are uploading a video file.\n", http.StatusBadRequest, nil)
	}
}

func uploadVideo(file multipart.File, fileHeader *multipart.FileHeader) error {
	tempFile, err := ioutil.TempFile(tempVideoStore, fmt.Sprintf("upload-*-%+v", fileHeader.Filename))
	if err != nil {
		return err
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	tempFile.Write(fileBytes)

	err = createVideoChunks(tempFile.Name())
	if err != nil {
		return err
	}

	err = os.Remove(tempFile.Name())
	if err != nil {
		return err
	}

	return nil
}

func createVideoChunks(fullPathAndExt string) error {
	fileNameAndExt := strings.TrimPrefix(fullPathAndExt, tempVideoStore+"/")
	extension := filepath.Ext(fullPathAndExt)
	fileName := fileNameAndExt[0 : len(fileNameAndExt)-len(extension)]

	err := os.Mkdir(videoDir+"/"+fileName, 0755)
	if err != nil {
		return err
	}

	ffmpeg := "ffmpeg"
	chunkSize := "10"
	args := []string{
		"-i",
		fullPathAndExt,
		"-codec:",
		"copy",
		"-start_number",
		"0",
		"-hls_time",
		chunkSize,
		"-hls_list_size",
		"0",
		"-f",
		"hls",
		videoDir + "/" + fileName + "/" + fileName + ".m3u8",
	}
	cmd := exec.Command(ffmpeg, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
