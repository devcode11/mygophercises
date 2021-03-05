package main

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		img := columnChart(10, 33, 64, 89, 12)
		png.Encode(w, img)
	})
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func columnChart(data ...int) image.Image {
	h, w := 100, len(data)*60+10
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	//border
	for x := 0; x < w; x++ {
		img.SetRGBA(x, 0, color.RGBA{0, 0, 0, 255})
		img.SetRGBA(x, h-1, color.RGBA{0, 0, 0, 255})
	}

	for y := 0; y < w; y++ {
		img.SetRGBA(0, y, color.RGBA{0, 0, 0, 255})
		img.SetRGBA(w-1, y, color.RGBA{0, 0, 0, 255})
	}

	columnsSrc := image.NewUniform(color.RGBA{100, 100, 255, 255})

	for i, dp := range data {
		if dp > 100 {
			panic(errors.New("data point must be less than or equal to 100"))
		}

		x0, y0 := i*60+10, 100-dp
		x1, y1 := (i+1)*60, 100

		draw.Draw(img, image.Rect(x0, y0, x1, y1), columnsSrc, image.Point{}, draw.Src)

		// for x := i*60 + 10; x < (i+1)*60; x++ {
		// 	for y := 100; y > (100 - dp); y-- {
		// 		img.Set(x, y, color.RGBA{180, 180, 255, 255})
		// 	}
		// }
	}

	return img
}

