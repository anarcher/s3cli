package main

import (
	"strings"
)

type S3Path struct {
	Bucket string
	Prefix string
}

func NewS3Path(path string) (*S3Path, error) {
	path = strings.Replace(path, "s3://", "", 1)
	parts := strings.SplitN(path, "/", 2)
	bucket := parts[0]
	prefix := parts[1]

	p := &S3Path{
		Bucket: bucket,
		Prefix: prefix,
	}
	return p, nil
}
