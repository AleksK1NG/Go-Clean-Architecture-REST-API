package utils

import (
	"errors"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/google/uuid"

	"github.com/AleksK1NG/api-mc/pkg/httpErrors"
)

var allowedImagesContentType = map[string]string{
	"image/png":  "png",
	"image/jpg":  "jpg",
	"image/jpeg": "jpeg",
}

func determineFileContentType(fileHeader textproto.MIMEHeader) (string, error) {
	contentTypes := fileHeader["Content-Type"]
	if len(contentTypes) < 1 {
		return "", httpErrors.NotAllowedImageHeader
	}
	return contentTypes[0], nil
}

func CheckImageContentType(image *multipart.FileHeader) error {
	// Check content type from header
	if !IsAllowedImageHeader(image) {
		return httpErrors.NotAllowedImageHeader
	}

	// Check real content type
	imageFile, err := image.Open()
	if err != nil {
		return httpErrors.BadRequest
	}
	defer imageFile.Close()

	fileHeader := make([]byte, 512)
	if _, err = imageFile.Read(fileHeader); err != nil {
		return httpErrors.BadRequest
	}

	if !IsAllowedImageContentType(fileHeader) {
		return httpErrors.NotAllowedImageHeader
	}
	return nil
}

func IsAllowedImageHeader(image *multipart.FileHeader) bool {
	contentType, err := determineFileContentType(image.Header)
	if err != nil {
		return false
	}
	_, allowed := allowedImagesContentType[contentType]
	return allowed
}

func GetImageExtension(image *multipart.FileHeader) (string, error) {
	contentType, err := determineFileContentType(image.Header)
	if err != nil {
		return "", err
	}

	extension, has := allowedImagesContentType[contentType]
	if !has {
		return "", errors.New("prohibited image extension")
	}
	return extension, nil
}

func GetImageContentType(image []byte) (string, bool) {
	contentType := http.DetectContentType(image)
	extension, allowed := allowedImagesContentType[contentType]
	return extension, allowed
}

func IsAllowedImageContentType(image []byte) bool {
	_, allowed := GetImageContentType(image)
	return allowed
}

func GetUniqFileName(userID string, fileExtension string) string {
	randString := uuid.New().String()
	return "userid_" + userID + "_" + randString + "." + fileExtension
}
