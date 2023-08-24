package src

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
  "log"
)

const (
	COLUMN_NAME_IMAGE_ID          = "id"
	COLUMN_NAME_IMAGE_TITLE       = "title"
	COLUMN_NAME_IMAGE_DESCRIPTION = "description"
	COLUMN_NAME_IMAGE_CATEGORY    = "category"
	COLUMN_NAME_IMAGE_TAGS        = "tags"
	COLUMN_NAME_IMAGE_URL         = "image_url"

  IMAGES_DIR_NAME = "images"
)

type Image struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	ImageUrl    string   `json:"image_url"`
	Description string   `json:"description"`
	Category    []string `json:"category"`
	Tags        []string `json:"tags"`
}

func ParseImageCsvFile(path string, nrows int) (image []Image, err error) {
	file, fileOpenErr := os.Open(path)
	if fileOpenErr != nil {
		return []Image{}, fileOpenErr
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var imageArray []Image
	var imageIdIndex, imageTitle, imageDescription, imageUrl, imageTags, imageCategory int
	var indexError error
	for i := 0; scanner.Scan(); i++ {
		st := scanner.Text()
		colValues := strings.Split(st, "\t")
		if i == 0 {
			imageIdIndex, indexError = getColIndex(COLUMN_NAME_IMAGE_ID, colValues)
			imageTitle, indexError = getColIndex(COLUMN_NAME_IMAGE_TITLE, colValues)
			imageDescription, indexError = getColIndex(COLUMN_NAME_IMAGE_DESCRIPTION, colValues)
			imageUrl, indexError = getColIndex(COLUMN_NAME_IMAGE_URL, colValues)
			imageTags, indexError = getColIndex(COLUMN_NAME_IMAGE_TAGS, colValues)
			imageCategory, indexError = getColIndex(COLUMN_NAME_IMAGE_CATEGORY, colValues)
			if indexError != nil {
				return []Image{}, indexError
			}
		} else {
			var imgId, imgTitle, imgDesc, imgUrl string
			var imgTags, imgCateg []string
			if len(colValues) > imageIdIndex {
				imgId = colValues[imageIdIndex]
			} else {
				imgId = ""
			}
			if len(colValues) > imageTitle {
				imgTitle = colValues[imageTitle]
			} else {
				imgTitle = ""
			}
			if len(colValues) > imageDescription {
				imgDesc = colValues[imageDescription]
			} else {
				imgDesc = ""
			}
			if len(colValues) > imageUrl {
				imgUrl = colValues[imageUrl]
			} else {
				imgUrl = ""
			}
			if len(colValues) > imageTags {
				imgTags = strings.Split(colValues[imageTags], ",")
			} else {
				imgTags = []string{}
			}
			if len(colValues) > imageCategory {
				imgCateg = strings.Split(colValues[imageCategory], ",")
			} else {
				imgCateg = []string{}
			}
			img := Image{
				Id:          imgId,
				Title:       imgTitle,
				Description: imgDesc,
				Tags:        imgTags,
				Category:    imgCateg,
				ImageUrl:    imgUrl,
			}
			imageArray = append(imageArray, img)
      if nrows > 0 && i == nrows {
        break
      }
		}
	}
	return imageArray, nil
}

func getColIndex(name string, colNames []string) (index int, err error) {
	for i, v := range colNames {
		if v == name {
			return i, nil
		}
	}
	return -1, errors.New("Column " + name + " not found in given array")
}


func (i* Image) Save(dir string) (bool, error, int) {
  imagePath := dir + "/" + IMAGES_DIR_NAME + "/" + i.Id
  _, openErr := os.Open(imagePath)
  if openErr != nil {
    log.Printf("Image %s not found locally. Trying to retrieve it\n", i.Id)
    if (strings.HasPrefix(i.ImageUrl, "http")) {
      resp, httpErr := http.Get(i.ImageUrl)
      if httpErr != nil {
        log.Println("HTTP Error : " + httpErr.Error())
        return false, httpErr, resp.StatusCode
      }
      defer resp.Body.Close()
      log.Printf("Downloading file from %s success\n", i.ImageUrl)
      file, fileErr := os.Create(imagePath)
      if fileErr != nil {
        log.Println("Create file Error: " + fileErr.Error())
      }
      defer file.Close()
      log.Printf("Create file %s \n", i.Id)
      _, copyErr := io.Copy(file, resp.Body)
      if copyErr != nil {
        fmt.Println("Copy file Error: " + copyErr.Error())
      } else {
        log.Printf("Copying image to %s \n", imagePath)
        return true, nil, resp.StatusCode
      }
    }
  } 
  log.Println("Image is already saved skip creation")
  return false, nil, -1
}
