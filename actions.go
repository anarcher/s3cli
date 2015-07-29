package main

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/codegangsta/cli"
	"log"
)

func GetAction(c *cli.Context) {
	if len(c.Args()) != 2 {
		log.Fatal("get s3path localpath")
	}

	s3path, err := NewS3Path(c.Args().Get(0))
	if err != nil {
		log.Fatal(err)
	}

	localpath := c.Args().Get(1)
	if localpath == "" {
		log.Fatal("get s3path localpath")
	}

	log.Printf("s3path Bucket:%v Prefix:%v", s3path.Bucket, s3path.Prefix)
	parallel := c.Int("parallel")

	manager := s3manager.NewDownloader(nil)
	d := NewDownloader(s3path, localpath, parallel, manager)

	client := s3.New(nil)
	params := &s3.ListObjectsInput{Bucket: &s3path.Bucket, Prefix: &s3path.Prefix}
	err = client.ListObjectsPages(params, d.eachPage)
	if err != nil {
		log.Fatal(err)
	}
	d.Wait()
}
