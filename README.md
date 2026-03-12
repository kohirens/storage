# Storage

This library provides an interface that is designed to be blue-print to
accessing data in various kinds of storage. It provides concrete structures for
local and AWS S3 bucket storages.

It should be trivial to add support for any other kind of storage. It stores
data it retrieves in RAM. So Load and Save need to implement methods for
getting and putting data to and from storage respectively.

## Status

[![quality-control](https://github.com/kohirens/storage/actions/workflows/quality-control.yml/badge.svg)](https://github.com/kohirens/storage/actions/workflows/quality-control.yml)

## Use Local Storge

It reads and stores data on a local drive. So performance is based on the type
of drive.

**Usage Example**:

```go

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/storage"
)

var (
	log     = &logger.Standard{}
	mainErr error
)

func main () {
	defer func() {
		if mainErr != nil {
			log.Errf("main error: %s", mainErr)
		}
	}()

	wd, e1 := filepath.Abs("/tmp")
	if e1 != nil {
		panic(fmt.Errorf("invalid dir: %v", e1.Error()))
	}

	storageDir := wd + "/storage"
	// make some directories to prevent dir exist errors.
	_ = os.MkdirAll(storageDir, 0777)

	//  to store session data.
	store, e1 := storage.NewLocalStorage(storageDir)
	if e1 != nil {
		mainErr = fmt.Errorf("init local storage: %s", e1.Error())
		return
	}
	if e := store.Save("test-01.txt", []byte("1234")); e != nil {
		mainErr = e
		return
	}
	data, e2 := store.Load("test-01.txt")
	if e2 != nil {
		mainErr = fmt.Errorf("load test data: %s", e2.Error())
		return
	}
	fmt.Printf("returned local data: %s", data)

	// Output:
	// returned local data: 1234
}
```


## Use AWS S3 Storage

Amazon Simple Storage Service (S3) for HTTP Session handling.

It stores data it retrieves from S3 in RAM. Plus it only sends data across the
network when you call Load or Save. This should provide good performance for the
average HTTP session use case.

However, if your storing large amounts of data then this may not be performant
for your use case.

**Usage Example**:

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/storage"
	"github.com/kohirens/www/session"
)

var (
	log = &logger.Standard{}
)

func main() {
	var mainErr error

	defer func() {
		if mainErr != nil {
			log.Errf("error: %v", mainErr)
		}
	}()
	bucket, ok := os.LookupEnv("S3_BUCKET_NAME")
	if !ok {
		mainErr = fmt.Errorf("unset environment variable S3_BUCKET_NAME")
		return
	}
	
	// Initialize storage backed by AWS S3.
	store, e1 := storage.NewBucketStorage(bucket, context.Background())
	if e1 !=  nil {
		mainErr = e1
		return
    }
	
	// Set where to store the session in the bucket.
	// The Session handler using RAM and then saving to Amazon S3 for
	// longer-term.
	sm := session.NewManager(store, "session", time.Minute*20)

	sm.Set("test", []byte("1234"))
	fmt.Printf("returned session key info: %s", sm.Get("test"))
}
```