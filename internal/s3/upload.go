package s3

import (
	"fmt"
	"io"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/kmjayadeep/totpm/internal/config"
)

func UploadIcon(body io.ReadSeeker) (string, error) {
	conf := config.Get()
	s3conf := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.Get().S3Key, config.Get().S3Secret, ""),
		Endpoint:         aws.String(conf.S3Endpoint),
		Region:           aws.String(conf.S3Region),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession := session.New(s3conf)

	s3Client := s3.New(newSession)
	key := fmt.Sprintf("icons/%s", uuid.New().String())

	obj := &s3.PutObjectInput{
		Key:    aws.String(key),
		Body:   body,
		Bucket: aws.String(config.Get().S3Bucket),
	}

	_, err := s3Client.PutObject(obj)
	if err != nil {
		return "", err
	}

	publicUrl := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.%s", *obj.Bucket, *s3conf.Endpoint),
		Path:   key,
	}

	return publicUrl.String(), nil
}
