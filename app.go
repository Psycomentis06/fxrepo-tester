package main

import (
	"context"
	"fxrepo_tester/src"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var cachedImagesCounter int
var cachedImagesMutex sync.Mutex

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

func (a *App) LoadFile(nrows int) (img []src.Image, err error) {
	path, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{})
	imgs, lErr := src.ParseImageCsvFile(path, nrows)
	a.cacheImages(&imgs)
	return imgs, lErr
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
		cachedImagesMutex.Lock()
		cachedImagesCounter++
		cachedImagesMutex.Unlock()
		eventData := make(map[string]interface{}, 2)
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
		}
		ch <- v
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
