package storage

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	bucket   string
	uploader *s3manager.Uploader
}

func NewS3(region, bucket string) *S3 {
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(region)},
	}))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	return &S3{
		bucket:   bucket,
		uploader: uploader,
	}
}

func (s3 *S3) Copy(localfile, remotefile string) error {
	f, err := os.Open(localfile)
	if err != nil {
		return err
	}

	_, err = s3.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3.bucket),
		Key:    aws.String(remotefile),
		Body:   f,
	})
	if err != nil {
		return err
	}

	return nil
}
