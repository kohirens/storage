package storage

import (
	"github.com/kohirens/stdlib/test"
	"os"
	"reflect"
	"testing"
)

const tmpDir = "tmp"

func TestMain(m *testing.M) {
	test.ResetDir(tmpDir, os.ModeDir|os.ModePerm)
	os.Exit(m.Run())
}

func TestLocalStorage_Load(runner *testing.T) {
	cases := []struct {
		name        string
		fName       string
		WorkDir     string
		want        []byte
		wantNewErr  bool
		wantSaveErr bool
		wantLoadErr bool
	}{
		{
			"save-then-load-success",
			"test-01",
			"tmp",
			[]byte("1234"),
			false,
			false,
			false,
		},
	}
	for _, c := range cases {
		runner.Run(c.name, func(t *testing.T) {
			s, e1 := NewLocalStorage(c.WorkDir)

			if (e1 != nil) != c.wantNewErr {
				t.Errorf("NewLocalStorage() error = %v, wantNewErr %v", e1, c.wantNewErr)
				return
			}

			if err := s.Save(c.fName, c.want); (err != nil) != c.wantSaveErr {
				t.Errorf("Render() error = %v, wantSaveErr %v", err, c.wantSaveErr)
				return
			}

			got, err := s.Load(c.fName)
			if (err != nil) != c.wantLoadErr {
				t.Errorf("Render() error = %v, wantLoadErr %v", err, c.wantLoadErr)
				return
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Get() got %s, want %s", got, c.want)
				return
			}
		})
	}
}
