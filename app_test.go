package main

import (
	"fmt"
	"math"
	"testing"
)

func TestCacheImage(t *testing.T)  {
  imgsl := make([]int, 55)
  for i := 0; i < len(imgsl); i++ {
    imgsl[i] = i * 10
  }
  imgs := &imgsl
  const chunkSize = 10
  chunksArraySize := int(math.Ceil(float64(len(*imgs)) / float64(chunkSize)))
  page := 1
  for i := 0; i < chunksArraySize; i++ {
    var imgsSlice []int
    if (i <= 0) {
      imgsSlice = (*imgs)[:chunkSize]
    } else {
      if chunksArraySize - 1 == i && len(*imgs) % chunkSize != 0 {
        imgsSlice = (*imgs)[page*chunkSize - chunkSize:]
      } else {
        imgsSlice = (*imgs)[page*chunkSize - chunkSize: chunkSize * page]
      }
    }
    page++
    fmt.Println(imgsSlice)
    fmt.Println("*************************")
  }
}
