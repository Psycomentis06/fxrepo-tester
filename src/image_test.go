package src

import (
	"fmt"
	"testing"
)

func TestParseImageCsvFile(t *testing.T) {
	a, e := ParseImageCsvFile("/home/psycomentis06/Documents/fxrepo tester data/unsplash-transformer/data/export/photos_with_tags.tsv", 10)
	if e != nil {
		fmt.Println("Error: ", e.Error())
	} 
  t.Log(len(a) > 0)
}

func TestSaveImage(t *testing.T) {
	a, _ := ParseImageCsvFile("/home/psycomentis06/Documents/fxrepo tester data/unsplash-transformer/data/export/photos_with_tags.tsv", 10)
  a[0].Save("/home/psycomentis06/.local/share/fxrepo_tester")
  a[1].Save("/home/psycomentis06/.local/share/fxrepo_tester")
  a[2].Save("/home/psycomentis06/.local/share/fxrepo_tester")
  a[3].Save("/home/psycomentis06/.local/share/fxrepo_tester")
}

// func TestPostImage(t *testing.T) {
//   img := Image{
//     Id: "id",
//     Title: "title",
//     Description: "desc",
//     ImageUrl: "/home/psycomentis06/.local/share/fxrepo_tester/images/_gycpm2K900",
//   }
//   img.PostToApi("http://localhost:9057/api/v1/file/image/new")
// }
