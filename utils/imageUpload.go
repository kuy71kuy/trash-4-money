package utils

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"mime/multipart"
	"os"
)

func generateRandomUUID() string {
	uuidObj := uuid.New()
	return uuidObj.String()
}
func UploadImage(file multipart.File) string {
	cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUD_NAME"), os.Getenv("CLOUDINARY_KEY"), os.Getenv("CLOUDINARY_SECRET"))
	ctx := context.Background()
	resp, _ := cld.Upload.Upload(ctx, file, uploader.UploadParams{PublicID: generateRandomUUID()})
	return resp.SecureURL
}
