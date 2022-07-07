package utils

import (
	"bufio"
	"os"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// FileReader is a structure for reading files in batches of lines
// reading files in this manner alleviates the stress on memory and garbage collector
type FileReader[T any] struct {
	file   *os.File
	reader *bufio.Reader
}

// NewFileReader returns a new file reader object
func NewFileReader[T any](filename string) *FileReader[T] {
	file, err := os.Open(filename)
	if err != nil {
		GracefulExit("Utils-File-Reader", err)
	}
	return &FileReader[T]{
		file:   file,
		reader: bufio.NewReader(file),
	}
}

// ReadLines reads the specified number of lines from the file (if possible)
// and binds each line to a specific JSON object specified by the generic parameter
func (r *FileReader[T]) ReadLines(count uint64) ([]*T, error) {
	jsonData := make([]*T, 0)
	for ; count > 0; count-- {
		line, _, err := r.reader.ReadLine()
		if err != nil {
			return jsonData, err
		}
		d := new(T)
		if err := json.Unmarshal(line, d); err != nil {
			return jsonData, err
		}
		jsonData = append(jsonData, d)
	}
	return jsonData, nil
}

// Close closes the file reader object
func (r *FileReader[T]) Close() error {
	return r.file.Close()
}
