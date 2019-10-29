package simple_uploader

import (
	"fmt"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/pborman/uuid"
	"os"
)

func generateFileName(extension string) string {
	return fmt.Sprintf("%s.%s", uuid.New(), extension)
}

func UploadFile(file *os.File, bucket string) (string, error) {
	session := session.New(&aws.Config{Region: aws.String("eu-central-1")})
	uploader := s3manager.NewUploader(session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		ContentType: aws.String("image/jpeg"),
		Body:        file,
		//Bucket:      aws.String("trendee-selfies"),
		Bucket: aws.String(bucket),
		Key:    aws.String(generateFileName("jpg")),
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		return "", err
	}
	return result.Location, err
}
