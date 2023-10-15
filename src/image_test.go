package src

import (
	"fmt"
	"net/http"
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
	a[0].Save()
	a[1].Save()
	a[2].Save()
	a[3].Save()
}

func TestCreateImageFile(t *testing.T) {
	img := Image{
		Id:          "id",
		Title:       "title",
		Description: "desc",
		ImageUrl:    "/home/psycomentis06/.local/share/fxrepo_tester/images/_gycpm2K900",
	}
	client := &http.Client{}
	file, err := img.CreateImageFile(client, "http://localhost:9057/api/v1/file/image/new")
	if err != nil {
		t.Errorf("Error while creating ImageFile: %s", err.Error())
	} else {
		t.Logf("ImageFile created successfully: ID => %s", file.Id)
	}
}

func TestImage_SaveToService(t *testing.T) {
	img := Image{
		Id:          "_gycpm2K900",
		Title:       "Sea Landscape",
		Description: "Open sea landscape",
		ImageUrl:    "/home/psycomentis06/.local/share/fxrepo_tester/images/_gycpm2K900",
		Tags:        []string{"landscape", "sea", "forest", "gray-blue sky"},
		Category:    []string{"landscape", "sea"},
	}
	client := &http.Client{}
	endpoints := &Endpoints{
		CreateImageFileEndpoint: "http://localhost:9057/v1/api/file/image/new",
		CreateImagePostEndpoint: "http://localhost:9057/v1/api/post/image/new",
		CreateCategoryEndpoint:  "http://localhost:9057/v1/api/category/new",
		GetCategoryEndpoint:     "http://localhost:9057/v1/api/category/",
	}
	imagePost, err := img.SaveToService(client, endpoints)
	if err != nil {
		t.Errorf("Save image to service failed: %s", err.Error())
	} else {
		var imageUrl string
		for _, variant := range imagePost.Image.Variants {
			if variant.Original {
				imageUrl = variant.Url
			}
		}
		t.Logf("Image saved successfully: ID => %s, Image Url => %s", imagePost.Id, imageUrl)
	}
}
