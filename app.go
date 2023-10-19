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
	//a.cacheImages(&imgs)
	a.cacheImagesV2(&imgs)
	images = &imgs
	return imgs, lErr
}

var saveImagesWg sync.WaitGroup

func (a *App) SubmitImagesV2(mainServerHost string) {
	if images == nil {
		log.Println("Can't submit images, images pointer is nil")
		return
	}
	queue := make(chan *src.Image)
	runtime.EventsEmit(a.ctx, "save-images-start")
	for i := 0; i < 10; i++ {
		go a.saveImageWorker(queue, mainServerHost, i, len(*images))
	}

	for i := 0; i < len(*images); i++ {
		saveImagesWg.Add(1)
		a.saveImageEnqueue(&(*images)[i], queue)
	}
	saveImagesWg.Wait()
	log.Printf("Finished saving %d images", len(*images))
	runtime.EventsEmit(a.ctx, "save-images-end")
}

func (a *App) saveImageWorker(ch chan *src.Image, mainServerHost string, goroutineIndex int, totalImages int) {
	for {
		select {
		case image := <-ch:
			client := &http.Client{}
			endpoints := &src.Endpoints{
				CreateImageFileEndpoint: mainServerHost + "/api/v1/file/image/new",
				CreateImagePostEndpoint: mainServerHost + "/api/v1/post/image/new",
				CreateCategoryEndpoint:  mainServerHost + "/api/v1/category/new",
				GetCategoryEndpoint:     mainServerHost + "/api/v1/category/",
			}
			imagePost, saveError := image.SaveToService(client, endpoints)
			if saveError != nil {
				log.Println("Save image to server failed: ", saveError)
			} else {
				savedImagesMutex.Lock()
				savedImagesCounter++
				savedImagesMutex.Unlock()
				imageSave := make(map[string]interface{}, 3)
				imageSave["image"] = imagePost
				imageSave["savedImages"] = savedImagesCounter
				imageSave["totalImages"] = totalImages
				saveImagesWg.Done()
				runtime.EventsEmit(a.ctx, "image-saved", imageSave)
			}
		}
	}
}

func (a *App) saveImageEnqueue(img *src.Image, queue chan *src.Image) {
	go func() {
		queue <- img
	}()
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
				imageSave["savedImages"] = savedImagesCounter
				imageSave["totalImages"] = numImages
				runtime.EventsEmit(a.ctx, "image-saved", imageSave)
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

var cacheImageWg sync.WaitGroup

func (a *App) cacheImagesV2(imgs *[]src.Image) {
	queue := make(chan *src.Image)
	runtime.EventsEmit(a.ctx, "cache-start")
	for i := 0; i < 10; i++ {
		go a.cacheWorker(queue, i, len(*imgs))
	}

	for i := 0; i < len(*imgs); i++ {
		cacheImageWg.Add(1)
		cacheEnqueue(&(*imgs)[i], queue)
	}
	cacheImageWg.Wait()
	log.Printf("Finished caching %d images\n", len(*imgs))
	runtime.EventsEmit(a.ctx, "cache-end")
}

func cacheEnqueue(img *src.Image, queue chan *src.Image) {
	go func() {
		queue <- img
	}()
}

func (a *App) cacheWorker(ch chan *src.Image, goroutineIndex int, totalImages int) {
	for {
		select {
		case image := <-ch:
			log.Printf("===================== Start goroutine N: %d ===================== \n", goroutineIndex)
			imgSaveErr := image.Save()
			if imgSaveErr != nil {
				if imgSaveErr.StatusCode == 429 {
					log.Println("Rate limit exceeded")
					return
				}
				if imgSaveErr.OriginalErr != nil {
					filePath := src.GetAppDirPath() + "/" + image.Id
					rmFileErr := os.Remove(filePath)
					if rmFileErr == nil {
						log.Println("Remove corrupted image: ", image.Id)
					}
					log.Println("Failed to save image", imgSaveErr)
					return
				}
			} else {
				log.Printf("Image %s saved successfully\n", image.Id)
				cachedImagesMutex.Lock()
				cachedImagesCounter++
				cachedImagesMutex.Unlock()
				eventData := make(map[string]interface{}, 3)
				eventData["cachedImages"] = cachedImagesCounter
				eventData["image"] = image
				eventData["totalImages"] = totalImages
				image.ImageUrl = src.GetImageCacheDir() + "/" + image.Id
				runtime.EventsEmit(a.ctx, "cache-event", eventData)
			}
			log.Printf("===================== End goroutine N: %d ===================== \n\r", goroutineIndex)
			cacheImageWg.Done()
		}
	}
}

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
