package videometafs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"streaming/models/videometamodel"
)

const videoDir = "media/video"
const tempVideoStore = "tempVideoStore"

// FSRepo Handles all logic for FS and FFMPEG.
type FSRepo struct{}

// New Constructor for videometafs.FSRepo.
func New() *FSRepo {
	return new(FSRepo)
}

// CacheVideo Saves file in temporary location on file system and returns path to that file.
func (fs *FSRepo) CacheVideo(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	if _, err := os.Stat(tempVideoStore); os.IsNotExist(err) {
		err := os.Mkdir(tempVideoStore+"/", 0755)
		if err != nil {
			return "", err
		}
	}

	tempFile, err := ioutil.TempFile(tempVideoStore, fmt.Sprintf("%+v", fileHeader.Filename))
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	tempFile.Write(fileBytes)
	name := tempFile.Name()
	return name, nil
}

// ChunkAndSaveVideo Uses FFMPEG to chunk videos and save to file system.
func (fs *FSRepo) ChunkAndSaveVideo(video videometamodel.Model) error {
	if _, err := os.Stat(videoDir + "/" + video.Filename); os.IsNotExist(err) {
		err := os.Mkdir(videoDir+"/"+video.Filename, 0755)
		if err != nil {
			return err
		}
	} else {
		return errors.New("video already exists")
	}

	// ffmpeg -i bunnyVideo.mp4 -codec: copy -start_number 0 -hls_time 1 -hls_list_size 0 -f hls bunnyVideo.m3u8
	ffmpeg := "ffmpeg"
	chunkSize := "1"
	args := []string{
		"-i",
		video.Tempfilepath,
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
		video.Filepath,
	}
	cmd := exec.Command(ffmpeg, args...)

	err := cmd.Run()
	if err != nil {
		return err
	}

	err = os.Remove(video.Tempfilepath)
	if err != nil {
		return err
	}

	return nil
}
