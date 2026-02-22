package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/stdlib/test"
	"github.com/kohirens/www/storage"
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

func ExampleNewLocalStorage() {
	var (
		log     = &logger.Standard{}
		mainErr error
	)

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
