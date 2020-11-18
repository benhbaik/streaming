package videometamodel

import (
	"path/filepath"
	"strings"
)

const videoDir = "media/video"
const tempVideoStore = "tempVideoStore"

// Model data about video.
type Model struct {
	Id           string
	Filename     string
	Fullfilename string
	Extension    string
	Directory    string
	Filepath     string
	Tempfilepath string
}

// New Creates an instance of videometamodel.Model.
func New(tempFileFullPath string) Model {
	v := new(Model)
	v.Tempfilepath = tempFileFullPath
	v.Fullfilename = strings.TrimPrefix(tempFileFullPath, tempVideoStore+"/")
	v.Extension = filepath.Ext(tempFileFullPath)
	v.Filename = v.Fullfilename[0 : len(v.Fullfilename)-len(v.Extension)]
	v.Filepath = videoDir + "/" + v.Filename + "/" + v.Filename + ".m3u8"
	return *v
}
