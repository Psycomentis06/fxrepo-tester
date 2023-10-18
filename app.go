package main

import (
	"context"
	"fxrepo_tester/src"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var cachedImagesCounter, savedImagesCounter int
var cachedImagesMutex, savedImagesMutex sync.Mutex

// App struct
type App struct {
	ctx context.Context
}

type LoadFileReturn struct {
	Images []src.Image
	Err    error
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func initAppDir() {
	src.GetAppDirPath()
	src.GetCacheDir()
	src.GetImageCacheDir()
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	initAppDir()
}

var images *[]src.Image

func (a *App) LoadFile(nrows int) (img []src.Image, err error) {
	path, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{})
	imgs, lErr := src.ParseImageCsvFile(path, nrows)
	a.cacheImages(&imgs)
	images = &imgs
	return imgs, lErr
}

func (a *App) SubmitImages(mainServerHost string) {
	if images == nil {
		log.Println("Can't submit images, images pointer is nil")
		return
	}
	runtime.EventsEmit(a.ctx, "save-images-start")
	imgs := *images
	var wgMain, wgWorker sync.WaitGroup
	ch := make(chan src.ImagePost, len(imgs))
	const chunkSize = 10
	numImages := len(imgs)
	chunksArraySize := (numImages / chunkSize) + 1
	for i := 0; i < chunksArraySize; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > numImages {
			end = numImages
		}
		imgsSlice := (imgs)[start:end]
		wgWorker.Add(1)
		go func() {
			defer wgWorker.Done()
			client := &http.Client{}
			endpoints := &src.Endpoints{
				CreateImageFileEndpoint: mainServerHost + "/api/v1/file/image/new",
				CreateImagePostEndpoint: mainServerHost + "/api/v1/post/image/new",
				CreateCategoryEndpoint:  mainServerHost + "/api/v1/category/new",
				GetCategoryEndpoint:     mainServerHost + "/api/v1/category/",
			}
			for _, img := range imgsSlice {
				savedImagePost, saveErr := img.SaveToService(client, endpoints)
				if saveErr != nil {
					log.Println(saveErr)
					return
				}
				savedImagesMutex.Lock()
				savedImagesCounter++
				savedImagesMutex.Unlock()
				imageSave := make(map[string]interface{}, 3)
				imageSave["image"] = savedImagePost
				imageSave["savedCounter"] = savedImagesCounter
				imageSave["totalImages"] = numImages
				runtime.EventsEmit(a.ctx, "image-saved", savedImagePost)
				ch <- savedImagePost
			}
		}()
	}

	go func() {
		wgWorker.Wait()
		close(ch)
	}()

	wgMain.Add(1)
	go func() {
		defer wgMain.Done()
		for range ch {
		}
	}()
	wgMain.Wait()
	runtime.EventsEmit(a.ctx, "save-images-end")
}

// func (a *App) cacheFileMetadata(path) {
//
// }

func (a *App) cacheImages(imgs *[]src.Image) {
	runtime.EventsEmit(a.ctx, "cache-start")
	const chunkSize = 10
	ch := make(chan src.Image, len(*imgs))
	var mainWg, workerWg sync.WaitGroup
	numImages := len(*imgs)
	chunksArraySize := (numImages + chunkSize - 1) / chunkSize

	for i := 0; i < chunksArraySize; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > numImages {
			end = numImages
		}

		imgsSlice := (*imgs)[start:end]

		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			a.cacheImage(&imgsSlice, ch, numImages)
		}()
	}

	go func() {
		workerWg.Wait()
		close(ch)
	}()

	mainWg.Add(1)
	go func() {
		defer mainWg.Done()
		for range ch {
		}
	}()

	mainWg.Wait()
	runtime.EventsEmit(a.ctx, "cache-end")
}

func (a *App) cacheImage(imgs *[]src.Image, ch chan src.Image, totalImages int) {
	for index, v := range *imgs {

		eventData := make(map[string]interface{}, 3)
		eventData["cachedImages"] = cachedImagesCounter
		eventData["image"] = v
		eventData["totalImages"] = totalImages
		runtime.EventsEmit(a.ctx, "cache-event", eventData)
		imgSaveErr := v.Save()
		if imgSaveErr != nil {
			if imgSaveErr.StatusCode == 429 {
				log.Println("Rate limit exceeded")
				return
			}
			if imgSaveErr.OriginalErr != nil {
				filePath := src.GetAppDirPath() + "/" + v.Id
				rmFileErr := os.Remove(filePath)
				if rmFileErr == nil {
					log.Println("Remove corrupted image: ", v.Id)
				}
				return
			}
		} else {
			(*imgs)[index].ImageUrl = src.GetImageCacheDir() + "/" + v.Id
			log.Println("Cached image: ", v.Id)
			cachedImagesMutex.Lock()
			cachedImagesCounter++
			cachedImagesMutex.Unlock()
			ch <- v
		}
	}
}

func (a *App) LoadFileAsync(nrows int) {
	runtime.EventsEmit(a.ctx, "load-file-start")
	go func() {
		obj, e := a.LoadFile(nrows)
		runtime.EventsEmit(a.ctx, "load-file-end", LoadFileReturn{Images: obj, Err: e})
	}()
}

func (a *App) OpenImage(path string) {
	// Supports only gwenview
	cmd := exec.Command("gwenview", "-f", path)
	_, cmdErr := cmd.Output()
	if cmdErr != nil {
		log.Println("Failed to open image.")
	}
}

func (a *App) PingServer(host string) bool {
	_, err := http.Get(host)
	if err != nil {
		return false
	}
	return true
}
