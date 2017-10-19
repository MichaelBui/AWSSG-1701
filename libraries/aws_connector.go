package libraries

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/michaelbui/AWSSG-1710/types"
	"io"
)

type (
	AwsConnector struct {
		configs types.Configs
	}
)

func NewAwsConnector(configs types.Configs) *AwsConnector {
	return &AwsConnector{
		configs: configs,
	}
}

func (conn *AwsConnector) newSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	}))
}

func (conn *AwsConnector) UploadVideo(file io.Reader, id string, ext string) error {
	uploader := s3manager.NewUploader(conn.newSession())
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:   aws.String(conn.configs["s3"].(types.AwsS3Config).Bucket),
		Key:      aws.String(conn.generateS3Key(id, ext)),
		Body:     file,
		Metadata: map[string]*string{"Id": aws.String(id)},
	})
	return err
}

func (conn *AwsConnector) generateS3Key(id string, ext string) string {
	return id + "/original" + ext
}

func (conn *AwsConnector) ConvertVideo(id string, ext string) error {
	configs := conn.configs["et"].(types.AwsETConfig)
	transcoderClient := elastictranscoder.New(conn.newSession())
	_, err := transcoderClient.CreateJob(&elastictranscoder.CreateJobInput{
		PipelineId: aws.String(configs.Pipeline),
		Input: &elastictranscoder.JobInput{
			Key: aws.String(conn.generateS3Key(id, ext)),
		},
		Output: &elastictranscoder.CreateJobOutput{
			Key:      aws.String(id + "/converted.mp4"),
			PresetId: aws.String(configs.Preset),
		},
		UserMetadata: map[string]*string{"Id": aws.String(id)},
	})
	return err
}
