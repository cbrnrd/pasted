package backends

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Backend struct {
	ctx context.Context

	// The function used to generate paths/keys.
	pathGenFunc PathGenFunc

	// The S3 bucket to use.
	bucket string

	// The S3 client to use.
	client *s3.Client
}

var _ Backend = (*S3Backend)(nil)

func NewS3Backend(ctx context.Context, pgf PathGenFunc, bucket string, client *s3.Client) *S3Backend {
	return &S3Backend{ctx: ctx, pathGenFunc: pgf, bucket: bucket, client: client}
}

// Get writes the contents of the file at key to w.
// If the key does not exist, Get should return nil.
//
// Note that the value is read into memory before being written to w.
func (b *S3Backend) Get(key string, w io.Writer) error {
	input := &s3.GetObjectInput{
		Bucket: &b.bucket,
		Key:    &key,
	}

	resp, err := b.client.GetObject(b.ctx, input)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, resp.Body)
	return err
}

// Put stores the contents of r in memory and returns the key.
func (b *S3Backend) Put(r io.Reader) (string, error) {

	key := b.pathGenFunc()
	input := &s3.PutObjectInput{
		Bucket: &b.bucket,
		Key:    &key,
		Body:   r,
	}

	_, err := b.client.PutObject(b.ctx, input)
	if err != nil {
		return "", err
	}

	return key, nil
}
