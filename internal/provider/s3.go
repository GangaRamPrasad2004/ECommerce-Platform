package providers

import (
	"context"
	"log"
	"mime/multipart"
	"strings"

	appconfig "github.com/GangaRamPrasad2004/ECommerce-Platform/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	//nolint:staticcheck // LocalStack compatibility requires the legacy uploader API in this project.
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Provider stores uploaded files in an Amazon S3-compatible service.
type S3Provider struct {
	client *s3.Client
	//nolint:staticcheck // LocalStack compatibility requires the legacy uploader API in this project.
	uploader *manager.Uploader
	bucket   string
	endpoint string
}

// NewS3Provider creates an S3 provider from the application configuration.
func NewS3Provider(cfg *appconfig.Config) *S3Provider {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		panic("failed to create AWS config " + err.Error())
	}
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.AWS.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.AWS.S3Endpoint)
			o.UsePathStyle = true

		}
	})
	return &S3Provider{
		client: client,
		//nolint:staticcheck // LocalStack compatibility requires the legacy uploader API in this project.
		uploader: manager.NewUploader(client),
		bucket:   cfg.AWS.S3Bucket,
		endpoint: cfg.AWS.S3Endpoint,
	}
}

// UploadFile uploads a multipart file to the requested object path.
func (p *S3Provider) UploadFile(file *multipart.FileHeader, path string) (string, error) {

	log.Printf("Uploading file %s using s3", path)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = src.Close() }()
	//nolint:staticcheck // LocalStack compatibility requires the legacy uploader API in this project.
	result, err := p.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
		Body:   src,
	})
	if err != nil {
		return "", err
	}
	return *result.Key, nil
}

// DeleteFile removes a previously uploaded S3 object.
func (p *S3Provider) DeleteFile(path string) error {
	_, err := p.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(strings.TrimPrefix(path, "/")),
	})
	return err
}
