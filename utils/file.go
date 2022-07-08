package utils

import (
	"bufio"
	"os"

	jsoniter "github.com/json-iterator/go"
)

// jsoniter has better peformance than standard encoding/json
// https://github.com/json-iterator/go
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
	var (
		line []byte
		err  error
		tmp  *T
	)
	for ; count > 0; count-- {
		line, err = r.readEntireLine()
		if err != nil {
			return jsonData, err
		}
		tmp = new(T)
		if err = json.Unmarshal(line, tmp); err != nil {
			return jsonData, err
		}
		jsonData = append(jsonData, tmp)
	}
	return jsonData, nil
}

// readEntireLine reads the entire line till it reaches a `\n` character or EOF
// this is required because *bufio.Reader.ReadLine() only returns a maximum of 65536 characters in each line at a time
// in case the number of characters exceed that, the returned []byte is broken in which case JSON unmarshalling will fail due to invalid syntax
// the function below ensures the entire line is read and JSON is parsed safely even if a break occurs due to bufio internals
// Reference -> https://devmarkpro.com/working-big-files-golang
func (r *FileReader[T]) readEntireLine() (ln []byte, err error) {
	var (
		isPrefix = true
		parts    []byte
	)
	ln = make([]byte, 0)
	for isPrefix && err == nil {
		parts, isPrefix, err = r.reader.ReadLine()
		ln = append(ln, parts...)
	}
	return ln, err
}

// Close closes the file reader object
func (r *FileReader[T]) Close() error {
	return r.file.Close()
}
