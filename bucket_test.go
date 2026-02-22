package storage

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www/session"
	"github.com/kohirens/www/storage"
)

// TestBucketStorage_List Test all the features of Save and List functionality.
func TestBucketStorage_List(tr *testing.T) {
	s, e1 := NewBucketStorage(
		os.Getenv("S3_BUCKET_NAME"),
		context.Background(),
	)
	if e1 != nil {
		tr.Fatal(e1)
	}
	s.Prefix = "list"
	if e := s.Save("/file-01.txt", []byte("01")); e != nil {
		tr.Fatal(e)
	}
	if e := s.Save("/file-02.txt", []byte("02")); e != nil {
		tr.Fatal(e)
	}

	cases := []struct {
		name                  string
		Name                  string
		location              string
		requestListParameters *RequestListParameters
		want                  []string
		wantErr               bool
	}{
		{
			name:     "can_list_files",
			location: "",
			want:     []string{"file-01.txt", "file-02.txt"},
			requestListParameters: &RequestListParameters{
				Prefix: s.Prefix,
			},
			wantErr: false,
		},
	}
	for _, tc := range cases {
		tr.Run(tc.name, func(t *testing.T) {
			s.SetRequestListParameters(tc.requestListParameters)
			got, gotErr := s.List(tc.location)
			if (gotErr != nil) != tc.wantErr {
				t.Errorf("List() error %v, wantErr %v", gotErr, tc.wantErr)
			}

			gotEmAll := 0
			wantEmAll := len(tc.want)

			for _, v := range got {
				for _, w := range tc.want {
					if w == v {
						gotEmAll++
					}
				}
			}

			if gotEmAll != wantEmAll {
				t.Errorf("List() got %v, want %v", got, tc.want)
			}
		})
	}
}

// TestBucketStorage_Load Test all the features of Save and Load functionality.
func TestBucketStorage_Load(tr *testing.T) {
	s, e1 := NewBucketStorage(
		os.Getenv("S3_BUCKET_NAME"),
		context.Background(),
	)
	if e1 != nil {
		tr.Fatal(e1)
	}
	s.Prefix = "list"
	if e := s.Save("/file-03.txt", []byte("03")); e != nil {
		tr.Fatal(e)
	}

	cases := []struct {
		name     string
		Name     string
		location string
		want     string
		wantErr  bool
	}{
		{
			name:     "can_load_file",
			location: "/file-03.txt",
			want:     "03",
			wantErr:  false,
		},
		{
			name:     "can_load_files",
			location: "file-does-not-exist.txt",
			wantErr:  true,
		},
	}
	for _, tc := range cases {
		tr.Run(tc.name, func(t *testing.T) {
			got, gotErr := s.Load(tc.location)
			if (gotErr != nil) != tc.wantErr {
				t.Errorf("Load() error %v, wantErr %v", gotErr, tc.wantErr)
			}

			if string(got) != tc.want {
				t.Errorf("Load() got %s, want %v", got, tc.want)
			}
		})
	}
}

func ExampleNewBucketStorage() {
	var (
		log     = &logger.Standard{}
		mainErr error
	)

	defer func() {
		if mainErr != nil {
			log.Errf("main error: %s", mainErr)
		}
	}()

	bucket, ok := os.LookupEnv("S3_BUCKET_NAME")
	if !ok {
		mainErr = fmt.Errorf("unset environment variable S3_BUCKET_NAME")
		return
	}
	//  to store session data.
	store, e1 := storage.NewBucketStorage(bucket, context.Background())
	if e1 != nil {
		mainErr = fmt.Errorf("init bucket storage: %s", e1.Error())
		return
	}
	// set where to store the session in the bucket.
	// HTTP Session handler using RAM and then saving to Amazon S3 for longer-term.
	sm := session.NewManager(store, "session", time.Minute*20)

	sm.Set("test", []byte("1234"))
	fmt.Printf("returned session key info: %s", sm.Get("test"))

	// Output:
	// returned session key info: 1234
}
