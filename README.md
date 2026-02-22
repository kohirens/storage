# Storage

`go get kohirens/storage`

# SSS

Amazon Simple Storage Service (S3) for HTTP Session handling.

It uses RAM to store/retrieve data, only when you call Load or Save does
it send data across the network. This should provide good performance for the
average HTTP session use case.

However, if your storing large amounts of data then this may not be performant
for your use case.


AWS S3 Bucket Example:
```go
package main

import (
	"context"
	"fmt"
	"github.com/kohirens/www/session"
	"github.com/kohirens/storage"
	"os"
	"time"
)


bucket, ok := os.LookupEnv("S3_BUCKET_NAME")
if !ok {
    mainErr = fmt.Errorf("unset environment variable S3_BUCKET_NAME")
    return
}
//  to store session data.
store := storage.NewBucketStorage(bucket, context.Background())
// set where to store the session in the bucket.
sessionHandler.Prefix = "session"
// HTTP Session handler using RAM and then saving to Amazon S3 for longer-term.
sm := session.NewManager(store, time.Minute * 20)

sm.Set("test", []bytes("1234"))
fmt.Printf("returned session key info: %v", sm.Get("test"))
```