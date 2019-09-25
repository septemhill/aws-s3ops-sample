package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadFile(sess *session.Session, filepath string) error {
	uploader := s3manager.NewUploader(sess)
	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filepath, err)
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("exampleBucket"),
		Key:    aws.String("file1"),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))

	return nil
}

func DownloadFile(sess *session.Session) error {
	downloader := s3manager.NewDownloader(sess)
	f, err := os.Create("qqq")
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", "qqq", err)
	}
	defer f.Close()

	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String("exampleBucket"),
		Key:    aws.String("file1"),
	})
	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}

	fmt.Printf("file downloaded, %d bytes\n", n)

	return nil
}

func DeleteObject(svc *s3.S3, bucket, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := svc.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}

	fmt.Println(result)
	return nil
}

func GetObject(svc *s3.S3, bucket, key string) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())

			default:
				fmt.Println(aerr.Error())
			}
			return aerr
		} else {
			return err
		}
	}

	f, err := os.Create(key)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(result.Body); err != nil {
		return err
	}

	if _, err := buf.WriteTo(f); err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}

func PutObject(svc *s3.S3, bucket, key string) error {
	f, err := os.Open("./docker-compose.yml")
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Body:    f,
		Bucket:  aws.String(bucket),
		Key:     aws.String(key),
		Tagging: aws.String("key1=value1&key2=value2"),
	}

	result, err := svc.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
				return aerr
			}
		} else {
			return err
		}
	}

	fmt.Println(result)
	return nil
}

func CreateBucket(svc *s3.S3) {
	input := &s3.CreateBucketInput{
		Bucket: aws.String("exampleBucket"),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String("us-east-1"),
		},
	}

	result, err := svc.CreateBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func ListBuckets(svc *s3.S3) {
	input := &s3.ListBucketsInput{}

	result, err := svc.ListBuckets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4572"),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	svc := s3.New(sess)
	_ = svc

	bucket, key := "exampleBucket", "file1"
	if err := PutObject(svc, bucket, key); err != nil {
		fmt.Println(err)
		return
	}
	if err := GetObject(svc, bucket, key); err != nil {
		fmt.Println(err)
		return
	}
	if err := DeleteObject(svc, bucket, key); err != nil {
		fmt.Println(err)
		return
	}
}
