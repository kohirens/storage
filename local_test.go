package storage

import (
	"bytes"
	"os"
	"testing"
)

func TestLocalStorage(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		workDir  string
		filename string
		wantErr  bool
	}{
		{
			name:     "TestLocalStorage",
			content:  []byte("test1234"),
			workDir:  tmpDir,
			filename: "test1234.txt",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorage{
				WorkDir: tt.workDir,
			}

			err := s.Save(tt.filename, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalStorage.Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			content, err2 := s.Load(tt.filename)
			if (err2 != nil) != tt.wantErr {
				t.Errorf("LocalStorage.Load(%v) error = %v, wantErr %v", tt.filename, err2, tt.wantErr)
			}

			if !bytes.Equal(content, tt.content) {
				t.Errorf("LocalStorage.Load(%v) = %v, want %v", tt.filename, content, tt.content)
			}

			if err3 := s.Remove(tt.filename); (err3 != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err3, tt.wantErr)
			}

			_, err4 := os.Stat(tmpDir + "/" + tt.filename)
			if os.IsExist(err4) {
				t.Errorf("LocalStorage file %v should have been removed", tt.filename)
			}
		})
	}
}

func TestLocalStorage_Exist(t *testing.T) {
	_ = os.WriteFile(tmpDir+"/file-exists.txt", []byte("test1234"), 0777)
	tests := []struct {
		name     string
		workDir  string
		filename string
		want     bool
	}{
		{
			name:     "has_filename",
			workDir:  tmpDir,
			filename: "file-exists.txt",
			want:     true,
		},
		{
			name:     "filename_not_found",
			workDir:  tmpDir,
			filename: "does-not-exist.txt",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorage{
				WorkDir: tt.workDir,
			}
			if got := s.Exist(tt.filename); got != tt.want {
				t.Errorf("Exist() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestLocalStorage_List(t *testing.T) {
	_ = os.MkdirAll(tmpDir+"/list/dir-01", 0777)
	s := &LocalStorage{
		WorkDir: tmpDir,
	}
	if e := s.Save("/list/file-01.txt", []byte("01")); e != nil {
		t.Fatal(e)
	}
	if e := s.Save("/list/file-02.txt", []byte("02")); e != nil {
		t.Fatal(e)
	}

	tests := []struct {
		name    string
		workDir string
		want    []string
		wantErr bool
	}{
		{
			name:    "can_list_files",
			workDir: "list",
			want:    []string{"file-01.txt", "file-02.txt", "dir-01"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := s.List(tt.workDir)
			if tt.wantErr != (gotErr != nil) {
				t.Errorf("List() err %v, want %v", gotErr, tt.wantErr)
			}

			gotEmAll := 0
			wantEmAll := len(tt.want)

			for _, v := range got {
				for _, w := range tt.want {
					if w == v {
						gotEmAll++
					}
				}
			}

			if gotEmAll != wantEmAll {
				t.Errorf("List() got %v, want %v", got, tt.want)
			}
		})
	}
}
