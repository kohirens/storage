package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type BucketStorage struct {
	Name                  string
	S3                    *s3.Client
	Prefix                string
	requestListParameters *RequestListParameters
}

var _ Storage = (*BucketStorage)(nil)

// Exist Verify the object is in the bucket.
func (s *BucketStorage) Exist(key string) bool {
	fullKey := s.Location(key)

	Log.Infof(stdout.LoadKey, fullKey)

	_, e1 := s.S3.HeadObject(
		context.Background(),
		&s3.HeadObjectInput{
			Bucket:       &s.Name,
			Key:          &fullKey,
			ChecksumMode: types.ChecksumModeEnabled,
		},
	)
	if e1 != nil {
		return false
	}

	return true
}

// List Files in a location in storage. It is not recursive, it only lists
// files in the specified directory.
//
//	There is special care that must be taken with S3. It requires a call to
//	ListObjectsV2, which requires parameters to be set. So there is a
//	prerequisite to call BucketStorage.SetRequestParams before calling
//	List on S3 storage. If you need to call BucketStorage.List multiple
//	times, care MUST be taken to call BucketStorage.SetRequestParams should
//	the parameters need to change. For that reason, once BucketStorage.List
//	is called, the request parameters are reset to nil.
func (s *BucketStorage) List(location string) ([]string, error) {
	requestParameter := s.requestListParameters
	if requestParameter == nil {
		return []string{}, fmt.Errorf("%v", stderr.RequestListParameters)
	}

	filePath := s.Location(location)

	Log.Dbugf(stdout.Load, filePath)

	lo, e1 := s.S3.ListObjectsV2(
		context.Background(),
		&s3.ListObjectsV2Input{
			Bucket: &s.Name,
			Prefix: &filePath,
		},
	)

	if e1 != nil {
		return []string{}, fmt.Errorf(stderr.ListFiles, e1.Error())
	}

	files := make([]string, 0)
	for _, v := range lo.Contents {
		files = append(files, strings.Replace(*v.Key, filePath+"/", "", 1))
	}

	return files, nil
}

// Load data from S3. It may be best to use a prefix, like the site domain,
// to prevent key name collision in the bucket. See an example at
// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/s3-checksums.html#use-service-S3-checksum-download
func (s *BucketStorage) Load(key string) ([]byte, error) {
	fullKey := s.Location(key)

	Log.Infof(stdout.LoadKey, fullKey)

	obj, e1 := s.S3.GetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket:       &s.Name,
			Key:          &fullKey,
			ChecksumMode: types.ChecksumModeEnabled,
		},
	)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.LoadKey, key, s.Name, e1.Error())
	}

	content, e2 := io.ReadAll(obj.Body)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.ReadObject, key)
	}

	return content, nil
}

// Location Mainly for internal use, this allows a prefix while ensuring all
// methods use a consistent location.
func (s *BucketStorage) Location(key string) string {
	return s.Prefix + key
}

// Save Uploads an object to S3, validating the checksum on success.
// For an example, see
// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/s3-checksums.html#use-service-S3-checksum-upload
func (s *BucketStorage) Save(key string, content []byte) error {
	fullKey := s.Location(key)

	Log.Infof(stdout.SaveKey, fullKey)

	_, e1 := s.S3.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Bucket:               &s.Name,
			Key:                  &fullKey,
			Body:                 bytes.NewReader(content),
			ChecksumAlgorithm:    types.ChecksumAlgorithmCrc32,
			ServerSideEncryption: "AES256",
		},
	)
	if e1 != nil {
		return fmt.Errorf(stderr.PutObject, e1.Error())
	}

	return nil
}

// Remove Delete an object from S3.
func (s *BucketStorage) Remove(key string) error {
	fullKey := s.Location(key)

	Log.Infof(stdout.SaveKey, fullKey)

	_, e1 := s.S3.DeleteObject(
		context.Background(),
		&s3.DeleteObjectInput{
			Bucket: &s.Name,
			Key:    &fullKey,
		},
	)
	if e1 != nil {
		return fmt.Errorf(stderr.DeleteObject, e1.Error())
	}

	return nil
}

// NewBucketStorage Initializes an S3 client to use as storage.
// Credentials are expected to be configured in the environment to be picked up
// by the AWS SDK. Panics on failure.
func NewBucketStorage(bucket string, ctx context.Context) (*BucketStorage, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf(stderr.AwsConfig, err)
	}

	return &BucketStorage{
		Name: bucket,
		S3:   s3.NewFromConfig(cfg),
	}, nil
}

func (s *BucketStorage) SetRequestListParameters(requestListParameters *RequestListParameters) {
	s.requestListParameters = requestListParameters
}

type RequestListParameters struct {
	ListType          int       `url:"list-type"`
	ContinuationToken string    `url:"continuation-token"`
	Delimiter         string    `url:"delimiter"`
	EncodingType      string    `url:"encoding-type"`
	FetchOwner        bool      `url:"fetch-owner"`
	MaxKeys           int       `url:"max-keys"`
	Prefix            string    `url:"prefix"`
	StartAfter        time.Time `url:"start-after"`
}
