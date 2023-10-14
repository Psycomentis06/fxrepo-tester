package src

import (
	"fmt"
	"testing"
)

func TestGetPage(t *testing.T) {
	imgsl := make([]int, 55)
	for i := 0; i < len(imgsl); i++ {
		imgsl[i] = i * 10
	}
	imgs := &imgsl
	s, _ := GetPage(*imgs, 10, 2)
	fmt.Println(s)
}

func TestGetPaths(t *testing.T) {
	t.Log(GetAppDirPath())
	t.Log(GetCacheDir())
	t.Log(GetImageCacheDir())
}
