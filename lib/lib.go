package lib

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
)

const (
	JsonFormat = "json"
)

func ReadToFile(fileName string, format string, container interface{}) (err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}

	filePath := filepath.Join(currentDir, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case JsonFormat:
		dataType := reflect.TypeOf(container)
		if dataType.Kind() != reflect.Ptr {
			return errors.New("container should be a pointer")
		}

		result := reflect.New(dataType.Elem()).Interface()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(result)
		if err != nil {
			return err
		}

		reflect.ValueOf(container).Elem().Set(reflect.ValueOf(result).Elem())
		return nil
	default:
		return errors.New("format not defined")
	}
}
