package main

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Downloader struct {
	*s3manager.Downloader
	wg        sync.WaitGroup
	s3path    *S3Path
	localpath string
	parallel  int
	downloadC chan string
}

func NewDownloader(s3path *S3Path, localpath string, parallel int, manager *s3manager.Downloader) *Downloader {

	d := &Downloader{
		s3path:     s3path,
		localpath:  localpath,
		parallel:   parallel,
		downloadC:  make(chan string, 1),
		Downloader: manager,
	}

	for i := 0; i < d.parallel; i++ {
		d.wg.Add(1)
		go d.downloadTo()
	}

	return d
}

func (d *Downloader) eachPage(page *s3.ListObjectsOutput, more bool) bool {
	for _, obj := range page.Contents {
		d.downloadC <- *obj.Key
	}

	if more == false {
		close(d.downloadC)
		return false
	}

	return true
}

func (d *Downloader) downloadTo() {
L:
	for {
		key, ok := <-d.downloadC
		if !ok {
			break L
		}
		d.downloadToFile(key)
	}
	d.wg.Done()
}

func (d *Downloader) Wait() {
	d.wg.Wait()
}

func (d *Downloader) downloadToFile(key string) {

	// Create the directories in the path
	file := filepath.Join(d.localpath, key)
	if err := os.MkdirAll(filepath.Dir(file), 0775); err != nil {
		log.Fatalf("Error:%v", err)
	}

	// Setup the local file
	fd, err := os.Create(file)
	if err != nil {
		log.Fatalf("Error:%v", err)
	}
	defer fd.Close()

	// Download the file using the AWS SDK
	log.Printf("Downloading s3://%s/%s to %s...\n", d.s3path.Bucket, key, file)
	params := &s3.GetObjectInput{Bucket: &d.s3path.Bucket, Key: &key}
	d.Download(fd, params)
}
