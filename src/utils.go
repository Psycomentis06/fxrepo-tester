package src

import (
	"errors"
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
