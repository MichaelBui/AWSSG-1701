package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type (
	File struct {
		gorm.Model
		Path   string `gorm:"size:255"`
		Status string `gorm:"size:16;default:'PENDING'"`
	}
)

func main() {
	db, err := gorm.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		panic("Unable to access database!")
	}
	defer db.Close()
	db.DropTableIfExists(&File{})
	db.CreateTable(&File{})

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/files", func(c echo.Context) error {
		files := []File{}
		db.Find(&files)
		return c.JSON(http.StatusOK, files)
	})

	e.POST("/files", func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		fileRow := File{Path: file.Filename}
		db.Create(&fileRow)

		awsSession := session.Must(session.NewSession(&aws.Config{
			Region: aws.String("ap-southeast-1"),
		}))
		s3uploader := s3manager.NewUploader(awsSession)
		result, err := s3uploader.Upload(&s3manager.UploadInput{
			Bucket:   aws.String("awssg-1710"),
			Key:      aws.String(file.Filename),
			Body:     src,
			Metadata: map[string]*string{"Id": aws.String(strconv.FormatInt(int64(fileRow.ID), 10))},
		})
		if err != nil {
			return err
		}

		s3Client := s3.New(awsSession)
		s3Res, err := s3Client.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String("awssg-1710"),
		})
		if err != nil {
			return err
		}

		transcoderClient := elastictranscoder.New(awsSession)
		etRes, err := transcoderClient.CreateJob(&elastictranscoder.CreateJobInput{
			PipelineId: aws.String("1508403132974-cj4bdz"),
			Input: &elastictranscoder.JobInput{
				Key: aws.String(file.Filename),
			},
			Output: &elastictranscoder.CreateJobOutput{
				Key:      aws.String(fmt.Sprintf("%v.mp4", time.Now().Unix())),
				PresetId: aws.String("1351620000001-000061"),
			},
			UserMetadata: map[string]*string{"Id": aws.String(strconv.FormatInt(int64(fileRow.ID), 10))},
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, struct {
			File     File
			S3Result *s3manager.UploadOutput
			ETResult *elastictranscoder.Job
			S3Keys   []*s3.Object
		}{
			File:     fileRow,
			S3Result: result,
			ETResult: etRes.Job,
			S3Keys:   s3Res.Contents,
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
