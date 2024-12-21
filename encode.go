package main

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"os"
)

func Load_Image(path string) image.Image {
	file, _ := os.Open(path)
	defer file.Close()
	img, _ := png.Decode(file)
	return img
}

func Save_Image(img image.Image, path string) {
	file, _ := os.Create(path)
	defer file.Close()
	png.Encode(file, img)
}

func (p *Mandelbrot_Plane) Plot_to_Image(max_iter int) image.Image {
	height, width := len(p.Plane), len(p.Plane[0])
	// width, height := x_steps, y_steps
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	for y, row := range p.Plane {
		// var speed uint8 = 1 // arbitrary number
		for x, point := range row {
			// iterations := point.iteration
			if point.Iteration != float64(max_iter) {
				var (
					red   uint8 = 1 // (speed * uint8(iterations)) % 255 // multily the colir intensity by color values out of phase (255/3 = 85).
					blue  uint8 = 1 // (speed*uint8(iterations) + 85) % 255
					green uint8 = 1 // (speed*uint8(iterations) + 85*2) % 255
				)
				color_val := uint8(255 - (255 / (1 + 0.05*point.Iteration)))
				img.Set(x, height-y, color.RGBA{color_val * red, color_val * blue, color_val * green, 255})
			} else {
				img.Set(x, height-y, color.RGBA{0, 0, 0, 255}) // black if didn't diverge.
			}
		}
	}
	return img
}

func (p *Mandelbrot_Plane) Plot_to_GIF(iter_per_frame int, max_iter int, delay int) {
	//TODO: output/ needs to be a temp directory not a permenant one!

	gif_images := []*image.Paletted{}
	delays := []int{}

	// create the gif
	for i := 1; i*iter_per_frame <= max_iter; i++ {
		p.Iterations(iter_per_frame)
		img := p.Plot_to_Image(iter_per_frame * i)
		bounds := img.Bounds()
		palettedImage := image.NewPaletted(bounds, palette.Plan9) // paletter?
		draw.FloydSteinberg.Draw(palettedImage, bounds, img, image.Point{})
		gif_images = append(gif_images, palettedImage)
		delays = append(delays, delay)
	}
	outGif := &gif.GIF{
		Image: gif_images,
		Delay: delays,
	}
	file, err := os.Create("out.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = gif.EncodeAll(file, outGif)
	if err != nil {
		log.Fatal(err)
	}
}
