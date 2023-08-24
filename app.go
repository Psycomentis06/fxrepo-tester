package main

import (
	"context"
	"fxrepo_tester/src"
	"log"
	"os"
	"time"

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
  runtime.EventsEmit(a.ctx, "cache-start")
  go func(){ a.cacheImages(imgs) }()
  return imgs, lErr
}

func (a *App) cacheImages(imgs []src.Image) {
  homeDir, homeDirErr := os.UserHomeDir()
  if homeDirErr != nil {
    panic("Home dir Error: " + homeDirErr.Error())
  }
  for _, v := range imgs {
    eventData := make(map[string]interface{})
    eventData["image"] = v
    eventData["totalImages"] = len(imgs)
    runtime.EventsEmit(a.ctx, "cache-event", eventData)
    isSaved, httpEr, statusCode := v.Save(homeDir + "/" + APP_DIR_PATH)
    if httpEr != nil {
      if (statusCode == 429) {
        log.Println("Rate limit exceeded")
        break
      }
    }
    if !isSaved {
      time.Sleep(time.Duration(500) * time.Millisecond)
    }
  }
  runtime.EventsEmit(a.ctx, "cache-end")
}

func (a *App) LoadFileAsync(nrows int) {
  runtime.EventsEmit(a.ctx, "load-file-start")
  go func() {
    obj, e := a.LoadFile(nrows)
    runtime.EventsEmit(a.ctx, "load-file-end", LoadFileReturn{Images: obj, Err: e})
  }()
}
