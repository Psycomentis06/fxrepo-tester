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
	"strconv"
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
			Message:     "HTTP Error: " + err.Error(),
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

func CreateImagePost(c *http.Client, endpoint string, bodyData ImagePostCreateModel) (ImagePost, error) {
	imagePostJson, jsonErr := json.Marshal(bodyData)
	if jsonErr != nil {
		return ImagePost{}, jsonErr
	}
	body := bytes.NewBuffer(imagePostJson)
	request, requestError := http.NewRequest("POST", endpoint, body)
	if requestError != nil {
		return ImagePost{}, requestError
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	res, resErr := c.Do(request)
	if resErr != nil {
		return ImagePost{}, resErr
	}
	defer res.Body.Close()
	status := res.StatusCode
	if status == http.StatusInternalServerError {
		var springErr SpringBootResponseError
		err := json.NewDecoder(res.Body).Decode(&springErr)
		if err != nil {
			return ImagePost{}, err
		}
		return ImagePost{}, errors.New(springErr.Error)
	}
	if status != http.StatusCreated {
		var resErr HttpResponseError
		err := json.NewDecoder(res.Body).Decode(&resErr)
		if err != nil {
			return ImagePost{}, err
		}
		return ImagePost{}, errors.New(resErr.Message)
	}
	var resData HttpResponse[ImagePost]
	err := json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		return ImagePost{}, err
	}
	return resData.Data, nil
}

func (i *Image) CreateImageFile(c *http.Client, endpoint string) (ImageFile, error) {
	file, openErr := os.ReadFile(i.ImageUrl)
	if openErr != nil {
		return ImageFile{}, openErr
	}
	bodyBuff := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyBuff)
	partFile, partFileErr := writer.CreateFormFile("file", i.Id)
	if partFileErr != nil {
		return ImageFile{}, partFileErr
	}
	_, partErr := partFile.Write(file)
	if partErr != nil {
		return ImageFile{}, partErr
	}
	writerCloseErr := writer.Close()
	if writerCloseErr != nil {
		return ImageFile{}, writerCloseErr
	}

	request, reqError := http.NewRequest("POST", endpoint, bodyBuff)
	if reqError != nil {
		return ImageFile{}, reqError
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	res, resErr := c.Do(request)
	if resErr != nil {
		return ImageFile{}, resErr
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusCreated:
		var resData HttpResponse[ImageFile]
		err := json.NewDecoder(res.Body).Decode(&resData)
		if err != nil {
			return ImageFile{}, err
		}
		return resData.Data, nil
	default:
		b, _ := io.ReadAll(res.Body)
		log.Println(string(b))
		return ImageFile{}, errors.New("unhandled HTTP Status Code:  " + strconv.Itoa(res.StatusCode))
	}
}

func (i *Image) SaveToService(c *http.Client, endpoints *Endpoints) (ImagePost, error) {
	// Check categories and create them if they don't exist
	for _, cat := range i.Category {
		categoryEndpoint := endpoints.GetCategoryEndpoint + cat
		_, err := GetCategory(c, categoryEndpoint)
		if err != nil {
			log.Println("Category " + cat + " does not exist, creating...")
			categoryObj := Category{Name: cat, Description: ""}
			_, err := CreateCategory(c, endpoints.CreateCategoryEndpoint, categoryObj)
			if err != nil {
				log.Println("Error creating category: " + err.Error())
			}
		}
	}
	// Create ImageFile
	imageFile, err := i.CreateImageFile(c, endpoints.CreateImageFileEndpoint)
	if err != nil {
		log.Println("Error creating image file: " + err.Error())
		return ImagePost{}, err
	}

	if len(i.Title) == 0 {
		if len(i.Description) == 0 {
			i.Title = "Untitled"
		} else {
			i.Title = i.Description
		}
	}
	if len(i.Description) == 0 {
		if len(i.Title) == 0 {
			i.Description = "No description"
		} else {
			i.Description = i.Title
		}
	}
	var categoryName string
	if len(i.Category) == 0 {
		categoryName = "default"
	} else {
		categoryName = i.Category[0]
	}
	imagePost := ImagePostCreateModel{
		Title:    i.Title,
		Content:  i.Description,
		Public:   true,
		Nsfw:     false,
		Tags:     i.Tags,
		Image:    imageFile.Id,
		Category: categoryName,
	}
	post, err := CreateImagePost(c, endpoints.CreateImagePostEndpoint, imagePost)
	if err != nil {
		log.Println("Error creating image post: " + err.Error())
		return ImagePost{}, err
	}
	return post, nil
}
