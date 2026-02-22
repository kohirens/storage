package storage

import "fmt"

type ErrDirNoExist struct {
	data string
}

func (e *ErrDirNoExist) Error() string {
	return fmt.Sprintf(stderr.DirNoExist, e.data)
}

type ErrReadFile struct {
	data string
}

func (e *ErrReadFile) Error() string {
	return fmt.Sprintf(stderr.ReadFile, e.data)
}

type ErrDecodeJSON struct {
	data string
}

func (e *ErrDecodeJSON) Error() string {
	return fmt.Sprintf(stderr.DecodeJSON, e.data)
}

type ErrEncodeJSON struct {
	data string
}

func (e *ErrEncodeJSON) Error() string {
	return fmt.Sprintf(stderr.EncodeJSON, e.data)
}

type ErrWriteFile struct {
	data string
}

func (e *ErrWriteFile) Error() string {
	return fmt.Sprintf(stderr.WriteFile, e.data)
}
