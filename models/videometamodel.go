package videometamodel

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const videoDir = "media/video"
const tempVideoStore = "tempVideoStore"

// Model data about video.
type Model struct {
	id           string
	filename     string
	fullfilename string
	extension    string
	directory    string
	filepath     string
}

// New Creates an instance of videometamodel.Model.
func New(tempFileFullPath string) Model {
	v := new(Model)
	v.fullfilename = strings.TrimPrefix(tempFileFullPath, tempVideoStore+"/")
	v.extension = filepath.Ext(tempFileFullPath)
	v.filename = v.fullfilename[0 : len(v.fullfilename)-len(v.extension)]
	v.filepath = videoDir + "/" + v.filename + "/" + v.filename + ".m3u8"
	return *v
}

// Insert saves video to FS and inserts metadata into DB.
// move this method to a class responsible for s3 bucket
func Insert(file multipart.File, fileHeader *multipart.FileHeader) error {
	if _, err := os.Stat(tempVideoStore); os.IsNotExist(err) {
		err := os.Mkdir(tempVideoStore+"/", 0755)
		if err != nil {
			return err
		}
	}

	tempFile, err := ioutil.TempFile(tempVideoStore, fmt.Sprintf("%+v", fileHeader.Filename))
	if err != nil {
		return err
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	tempFile.Write(fileBytes)

	err = saveVideoToFS(tempFile.Name())
	if err != nil {
		return err
	}

	return nil
}

func saveVideoToFS(tempFilePath string) error {
	v := New(tempFilePath)

	if _, err := os.Stat(videoDir + "/" + v.filename); os.IsNotExist(err) {
		err := os.Mkdir(videoDir+"/"+v.filename, 0755)
		if err != nil {
			return err
		}
	}

	// ffmpeg -i bunnyVideo.mp4 -codec: copy -start_number 0 -hls_time 1 -hls_list_size 0 -f hls bunnyVideo.m3u8
	ffmpeg := "ffmpeg"
	chunkSize := "1"
	args := []string{
		"-i",
		tempFilePath,
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
		v.filepath,
	}
	cmd := exec.Command(ffmpeg, args...)

	err := cmd.Run()
	if err != nil {
		return err
	}

	err = os.Remove(tempFilePath)
	if err != nil {
		return err
	}

	return nil
}
