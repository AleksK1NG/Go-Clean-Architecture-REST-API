package models

import "io"

// AWS Upload Input
type UploadInput struct {
	File        io.Reader
	Name        string
	Size        int64
	ContentType string
	BucketName  string
}
