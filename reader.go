package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

var bucketName = "waldo-recruiting"

type PhotoReader struct {
	svc    *s3.S3
	bucket *string
}

func NewReader() (*PhotoReader, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)
	return &PhotoReader{svc: svc, bucket: &bucketName}, nil
}

func (r *PhotoReader) List() ([]*s3.Object, error) {
	//SDK Ref: http://docs.aws.amazon.com/sdk-for-go/api/service/s3/#ListObjectsOutput
	params := &s3.ListObjectsInput{Bucket: r.bucket}
	resp, err := r.svc.ListObjects(params)

	return resp.Contents, err
}

type Walker struct {
	tags map[string]string
}

func (w Walker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	val, _ := tag.StringVal()
	if val != "" {
		w.tags[string(name)] = val
	}
	return nil
}

type PhotoResult struct {
	tags map[string]string
	key  string
}

func (r *PhotoReader) worker(id int, jobs chan string, out chan *PhotoResult) {
	for k := range jobs {
		fmt.Println("worker", id, "reading photo", k)
		p, _ := r.FetchEXIF(k)
		out <- p
	}
}

//Fetch Photo data and parse EXIF
func (r *PhotoReader) FetchEXIF(key string) (*PhotoResult, error) {
	//Read
	params := &s3.GetObjectInput{
		Bucket: r.bucket, // Required
		Key:    &key,     // Required
	}
	resp, e := r.svc.GetObject(params)

	if e != nil {
		fmt.Println("Error photo ", key, "Details: ", e.Error())
		return nil, e
	}
	//Parse EXIF Data
	x, err := exif.Decode(resp.Body)
	if err != nil {
		fmt.Println("Error photo ", key, "Details: ", err.Error())
		return nil, err
	}
	w := &Walker{tags: make(map[string]string)}
	x.Walk(w) 
	return &PhotoResult{tags: w.tags, key: key}, err
}
