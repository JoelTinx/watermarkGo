package watermarkGo

import (
	"fmt"
	"image"
	"image/png"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"runtime"
)

func init() {
  //os.Mkdir("." + string(filepath.Separator) + "output",0777)
  os.Mkdir("output",0777) // Output files proccess
}

func whenError(err error) {
	if err != nil {
		panic(err)
		return
	}
}

func setWaterMark(pathImgBG string, pathImgFG string, despX int, despY int) {
	fileImgBack, err := os.Open(pathImgBG)
	whenError(err)
	defer fileImgBack.Close()

	fileImageFront, err := os.Open(pathImgFG)
	whenError(err)
	defer fileImageFront.Close()

	imgBack, err := jpeg.Decode(fileImgBack)
	whenError(err)
	sbx, sby := imgBack.Bounds().Dx(), imgBack.Bounds().Dy()

  imgWaterMark, err := png.Decode(fileImageFront)
  whenError(err)
  swx, swy := imgWaterMark.Bounds().Dx(), imgWaterMark.Bounds().Dy()

  iX, iY := sbx - swx - despX, sby - swy - despY

  outputFile := image.NewRGBA(image.Rect(0, 0, sbx, sby))
  if sbx > swx && sby > swy {
    for x := 0; x < sbx; x++ {
      for y := 0; y < sby; y++ {
				if valida(x, iX, y, iY, swx, swy) {
          _, _, _, a := imgWaterMark.At(x - iX, y - iY).RGBA()
          if a != 0 {
            outputFile.Set(x, y, imgWaterMark.At(x - iX, y - iY))
          } else {
						outputFile.Set(x, y, imgBack.At(x, y))
					}
        } else {
					outputFile.Set(x, y, imgBack.At(x, y))
				}
      }
    }
  } else {
    fmt.Println("Imagen de sello de agua es más grande que el fondo")
  }

  outputFileName := strings.Replace(pathImgBG, "\\", "_", -1)

	newFile, err := os.Create("output/" + outputFileName)
  whenError(err)
	defer newFile.Close()

	err = png.Encode(newFile, outputFile)
  whenError(err)

  fmt.Println("Imagen procesada ", pathImgBG)
}

func valida(x, iX, y, iY, swx, swy int) bool {
	return (x > iX && y > iY) && (x < (iX + swx) && y < (iY + swy))
}

/*******************************************************************************
Public Functions
*******************************************************************************/
func SetImageWaterMark(pathImgBG, pathImgFG string) {
	setWaterMark(pathImgBG, pathImgFG, 20, 20)
}

func SetDirWaterMark(pathDirBG, pathImgFG string) {
	files := []string{}

  filepath.Walk(pathDirBG, func(path string, info os.FileInfo, err error) error {
    if !info.IsDir() {
      files = append(files, path)
    }
    return err
  })

  filesNameChannel := make(chan string)
  numCores := runtime.NumCPU()
  runtime.GOMAXPROCS(numCores)

  var wg sync.WaitGroup

  for i := 0; i < numCores; i++ {
    go processImageInDir(filesNameChannel, &wg, pathImgFG)
  }

  wg.Add(len(files))
  for _, fileName := range files {
    filesNameChannel <- fileName
  }

  wg.Wait()
  fmt.Println("\nPROCESS DONE!")
}

func processImageInDir(fileChannel chan string, wg *sync.WaitGroup, pathImgFG string)  {
	for {
    fileName := <- fileChannel
    fmt.Println("Imagen Recibida: ", fileName)
    setWaterMark(fileName, pathImgFG, 25, 25)
    wg.Done()
  }
}
