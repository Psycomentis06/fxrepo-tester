package src

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const (
	COLUMN_NAME_IMAGE_ID          = "id"
	COLUMN_NAME_IMAGE_TITLE       = "title"
	COLUMN_NAME_IMAGE_DESCRIPTION = "description"
	COLUMN_NAME_IMAGE_CATEGORY    = "category"
	COLUMN_NAME_IMAGE_TAGS        = "tags"
	COLUMN_NAME_IMAGE_URL         = "image_url"
)

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

func (i *Image) Save() *ImageSaveError {
	imagePath := GetImageCacheDir() + "/" + i.Id
	if _, openErr := os.Stat(imagePath); openErr == nil {
		return nil
	}
	if !strings.HasPrefix(i.ImageUrl, "http") {
		return &ImageSaveError{
			Message:     "Invalid URL format",
			OriginalErr: nil,
			StatusCode:  -1,
		}
	}

	resp, err := http.Get(i.ImageUrl)
	if err != nil {
		return &ImageSaveError{
			Message:     "HTTP Error",
			OriginalErr: err,
			StatusCode:  -1,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &ImageSaveError{
			Message:     "HTTP Request Failed",
			OriginalErr: nil,
			StatusCode:  resp.StatusCode,
		}
	}

	file, err := os.Create(imagePath)
	if err != nil {
		return &ImageSaveError{
			Message:     "Create file Error",
			OriginalErr: err,
			StatusCode:  -1,
		}
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return &ImageSaveError{
			Message:     "Copy file Error",
			OriginalErr: err,
			StatusCode:  -1,
		}
	}
	return nil
}

func (i *Image) PostToApi(endpoint string) {
	file, openErr := os.ReadFile(i.ImageUrl)
	if openErr != nil {
		log.Println(openErr.Error())
	}
	bodyBuff := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyBuff)
	partFile, partFileErr := writer.CreateFormFile("file", i.Id)
	if partFileErr != nil {
		log.Println("Part error")
	}
	_, partErr := partFile.Write(file)
	if partErr != nil {
		log.Println(partErr.Error())
	}
	writerCloseErr := writer.Close()
	if writerCloseErr != nil {
		log.Println(writerCloseErr)
	}

	request, reqError := http.NewRequest("POST", endpoint, bodyBuff)
	if reqError != nil {
		log.Println(reqError.Error())
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, resErr := client.Do(request)
	if resErr != nil {
		log.Println(resErr)
	}
	defer res.Body.Close()

	var resData map[string]interface{}
	json.NewDecoder(res.Body).Decode(&resData)
	log.Println(resData)

	//  imageFileValues := url.Values{
	//    "file": {i.ImageUrl},
	//  }
	//  res, httpErr := http.PostForm(endpoint, imageFileValues)
	//  if httpErr != nil {
	//    log.Println(httpErr.Error())
	//  }
	//  defer res.Body.Close()
	// var resData map[string]interface{}
	// json.NewDecoder(res.Body).Decode(&resData)
	// log.Println(resData)
}
