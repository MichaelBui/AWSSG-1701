package types

import "io"

type (
	CloudConnector interface {
		UploadVideo(file io.Reader, id string, ext string) error
		ConvertVideo(id string, ext string) error
	}
)
