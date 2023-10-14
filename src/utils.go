package src

import (
	"errors"
	"os"
)

const (
	IMAGES_DIR_NAME = "images"
	CACHE_DIR_NAME  = "cache"
)

func GetPage(array []int, chunk int, page int) ([]int, error) {
	if chunk <= 0 {
		return nil, errors.New("Chunk size must be greater than 0")
	}

	if page <= 0 {
		return nil, errors.New("Page number must be greater than 0")
	}

	chunksArraySize := (len(array) + chunk - 1) / chunk // Rounds up to account for incomplete chunks

	if page > chunksArraySize {
		return nil, errors.New("Page number bigger than current pages")
	}

	start := (page - 1) * chunk
	end := start + chunk
	if end > len(array) {
		end = len(array) // Don't go past the end of the array
	}

	return array[start:end], nil
}

func GetAppDirPath() string {
	home := os.Getenv("HOME")
	return home + "/.local/share/fxrepo_tester"
}

func GetAppSubDir(name string) string {
	dir := GetAppDirPath() + "/" + name
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		mkdirErr := os.MkdirAll(dir, os.ModePerm)
		if mkdirErr != nil {
			panic(mkdirErr)
		}
	}
	return dir
}

func GetImageCacheDir() string {
	return GetAppSubDir(IMAGES_DIR_NAME)
}

func GetCacheDir() string {
	return GetAppSubDir(CACHE_DIR_NAME)
}
