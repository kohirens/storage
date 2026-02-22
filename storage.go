package storage

import "github.com/kohirens/stdlib/logger"

// Storage Save data for long term.
type Storage interface {
	// Exist Verification the file is in storage.
	Exist(name string) bool
	// New method Storage.List will list files in a directory.

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
	List(location string) ([]string, error)
	// Load Retrieve data from storage.
	Load(filename string) ([]byte, error)
	// Location Get the location in storage. This does not check for existence.
	Location(filename string) string
	// Save Write data to storage.
	Save(filename string, data []byte) error
	// Remove data from storage.
	Remove(filename string) error
}

var Log = &logger.Standard{}
