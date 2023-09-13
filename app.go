package main

import (
	"context"
	"fxrepo_tester/src"
	"log"
	"math"
	"os"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
  APP_DIR_PATH = "/.local/share/fxrepo_tester"
)

// App struct
type App struct {
	ctx context.Context
}

type LoadFileReturn struct {
  Images []src.Image
  Err error
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}


func initAppDir() {
  e := os.Getenv("HOME")
  if len(e) > 0 {
    dirPath := e + APP_DIR_PATH + "/" + src.IMAGES_DIR_NAME
    _, err := os.Stat(dirPath)
    if err != nil && os.IsNotExist(err) {
      mkdirErr := os.MkdirAll(dirPath, os.ModePerm)
      if mkdirErr != nil {
        panic(mkdirErr)
      }
    }
  }
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
  initAppDir()
}

func (a *App) LoadFile(nrows int) (img []src.Image,err error) {
	path, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{})
  imgs, lErr := src.ParseImageCsvFile(path, nrows)
  a.cacheImages(&imgs)
  return imgs, lErr
}

func (a *App) cacheImages(imgs *[]src.Image) {
  runtime.EventsEmit(a.ctx, "cache-start")
  homeDir, homeDirErr := os.UserHomeDir()
  if homeDirErr != nil {
    panic("Home dir Error: " + homeDirErr.Error())
  }
  const chunkSize = 10
  ch := make(chan src.Image, len(*imgs))
  var mainWg, workerWg sync.WaitGroup
  chunksArraySize := int(math.Round(float64(len(*imgs)) / float64(chunkSize)))
  page := 1
  for i := 0; i < chunksArraySize; i++ {
    var imgsSlice []src.Image 
    if (i <= 0) {
      imgsSlice = (*imgs)[:chunkSize]
    } else {
      if chunksArraySize - 1 == i && len(*imgs) % chunkSize != 0 {
        imgsSlice = (*imgs)[page*chunkSize - chunkSize:]
      } else {
        imgsSlice = (*imgs)[page*chunkSize - chunkSize: chunkSize * page]
      }
    }
    workerWg.Add(1)
    go func(){
      defer workerWg.Done()
      a.cacheImage(&imgsSlice, ch, homeDir, len(*imgs))
    }()
    page++
  }
  go func() {
    workerWg.Wait()
    close(ch)
  }()

  mainWg.Add(1)
  go func() {
    defer mainWg.Done()
    for range ch {}
  }()

  mainWg.Wait()
  runtime.EventsEmit(a.ctx, "cache-end")
}

func (a *App) cacheImage(imgs *[]src.Image, ch chan src.Image, homeDir string, totalImages int) {
  for index, v := range *imgs {
    eventData := make(map[string]interface{}, 2)
    eventData["image"] = v
    eventData["totalImages"] = totalImages
    runtime.EventsEmit(a.ctx, "cache-event", eventData)
    _ , httpEr, statusCode := v.Save(homeDir + APP_DIR_PATH)
    if httpEr != nil {
      if (statusCode == 429) {
        log.Println("Rate limit exceeded")
        return
      }
    }
    (*imgs)[index].ImageUrl = homeDir + "/" + APP_DIR_PATH + "/images/" + v.ImageUrl
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
