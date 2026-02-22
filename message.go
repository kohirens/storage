package storage

var stderr = struct {
	AwsConfig,
	DecodeJSON,
	DeleteObject,
	DirNoExist,
	EncodeJSON,
	ReadFile,
	RequestListParameters,
	ListFiles,
	LoadKey,
	ReadObject,
	RemoveFile,
	PutObject,
	WriteFile string
}{
	AwsConfig:             "failed to load AWS config: %v",
	DecodeJSON:            "cannot decode JSON: %v",
	DeleteObject:          "cannot delete object: %v",
	DirNoExist:            "%v directory does not exist",
	EncodeJSON:            "cannot encode JSON: %v",
	ReadFile:              "cannot read file %v",
	RequestListParameters: "RequestListParameters has not been set, it is a requirement to call BucketStorage.SetRequestListParameters(*) before calling BucketStorage.List(*). This is a condition to keep the Storage API consistent across mediums while allowing S3 to have its differences.",
	ListFiles:             "cannot list files %v",
	LoadKey:               "cannot load object key %v in bucket %v: %v",
	PutObject:             "cannot put object: %v",
	ReadObject:            "cannot read object: %v",
	RemoveFile:            "cannot remove file %v",
	WriteFile:             "attempting to write, but cannot %v",
}
var stdout = struct {
	Load,
	LoadKey,
	SaveKey string
}{
	Load:    "loading %v",
	LoadKey: "loading object from key %v",
	SaveKey: "saving object to key %v",
}
