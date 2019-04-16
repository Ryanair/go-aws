package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ryanair/goaws"
)

type Client struct {
	s3 *s3.S3
}

func NewClient(cfg *goaws.Config, options ...func(*s3.S3)) *Client {
	cli := s3.New(cfg.Provider)
	for _, opt := range options {
		opt(cli)
	}

	return &Client{s3: cli}
}

func Endpoint(endpoint string) func(*s3.S3) {
	return func(s3 *s3.S3) {
		s3.Endpoint = endpoint
	}
}

func (c *Client) GeneratePutURL(bucket, key, contentType string, expire time.Duration) (string, error) {
	req, _ := c.s3.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		ContentType: &contentType,
	})

	url, err := req.Presign(expire)
	if err != nil {
		return "", newErr(err, SigningURLErrCode, "signing url failed")
	}

	return url, nil
}

func (c *Client) DeleteObject(bucket, key string) error {
	if _, err := c.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return newOpsErr(err, "delete object failed")
	}

	return nil
}

func (c *Client) GetObject(bucket, key string) (io.ReadCloser, error) {
	out, err := c.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, newOpsErr(err, "get object failed")
	}

	return out.Body, nil
}

