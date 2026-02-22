package storage

import (
	"context"
	"os"
	"testing"
)

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
