package videoservice

import (
	"mime/multipart"
	"streaming/models/videometamodel"
	videometafs "streaming/repos/videometafs"
	videometasql "streaming/repos/videometasql"
)

// VideoService Handles business logic for videometadata.
type VideoService struct {
	SQLRepo *videometasql.SQLRepo
	FSRepo  *videometafs.FSRepo
}

// New Constructor for videometaservice.Model.
func New() *VideoService {
	vs := new(VideoService)
	vs.SQLRepo = videometasql.New()
	vs.FSRepo = videometafs.New()
	return vs
}

// Add Adds videometamodel.Model to DB.
func (s *VideoService) Add(file multipart.File, fileHeader *multipart.FileHeader) error {
	tempFilePath, err := s.FSRepo.CacheVideo(file, fileHeader)
	if err != nil {
		return err
	}
	video := videometamodel.New(tempFilePath)
	s.FSRepo.ChunkAndSaveVideo(video)
	s.SQLRepo.Insert(video)
	return nil
}
