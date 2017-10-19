package models

import (
	"github.com/michaelbui/AWSSG-1710/entities"
	"github.com/michaelbui/AWSSG-1710/types"
	"mime/multipart"
	"path"
)

type (
	File struct{}
)

func NewFile() *File {
	return &File{}
}

func (f *File) Find(db types.DatabaseConnector) ([]entities.File, error) {
	l := &[]entities.File{}
	err := db.Find(l)
	return *l, err
}

func (f *File) SaveInfoToDB(db types.DatabaseConnector) (uint, error) {
	entity := &entities.File{}
	err := db.Save(entity)
	if err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (f *File) UploadVideoToCloud(c types.CloudConnector, video *multipart.FileHeader, cloudId string) error {
	file, err := video.Open()
	if err != nil {
		return err
	}
	return c.UploadVideo(file, cloudId, path.Ext(video.Filename))
}

func (f *File) ConvertVideoOnCloud(c types.CloudConnector, video *multipart.FileHeader, cloudId string) error {
	return c.ConvertVideo(cloudId, path.Ext(video.Filename))
}
