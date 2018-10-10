package helper

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/crypto/bcrypt"
)

func DecodeDataFromJsonFile(f *os.File, data interface{}) error {
	jsonParser := jsoniter.NewDecoder(f)
	err := jsonParser.Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

func CheckStringElementInSlice(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}

func GetPosStringElementInSlice(list []string, str string) int {
	for i, item := range list {
		if item == str {
			return i
		}
	}
	return -1
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
