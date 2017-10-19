package controllers

import (
	"fmt"
	"github.com/michaelbui/AWSSG-1710/entities"
	"github.com/michaelbui/AWSSG-1710/models"
	"github.com/michaelbui/AWSSG-1710/types"
	"mime/multipart"
)

type (
	FileController struct {
		configs *types.AppConfigs
	}
)

func NewFileController(configs *types.AppConfigs) *FileController {
	return &FileController{
		configs: configs,
	}
}

func (c *FileController) Post(file *multipart.FileHeader) (uint, error) {
	model := models.NewFile()
	dbConn := getDatabaseConnector(c.configs.DB)
	id, err := model.SaveInfoToDB(dbConn)
	if err != nil {
		return id, err
	}

	cloudId := fmt.Sprint(id)
	cloudConn := getCloudConnector(c.configs.Cloud)
	if err := model.UploadVideoToCloud(cloudConn, file, cloudId); err != nil {
		return id, err
	}
	if err := model.ConvertVideoOnCloud(cloudConn, file, cloudId); err != nil {
		return id, err
	}

	return id, nil
}

func (c *FileController) List() ([]entities.File, error) {
	model := models.NewFile()
	dbConn := getDatabaseConnector(c.configs.DB)
	return model.Find(dbConn)
}
