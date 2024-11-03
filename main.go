package main

import (
	"fmt"
	"image/color"
	"runtime"
	"sync"
	"time"

	"image"
	"image/png"
	"os"
)

// set plane of complex mandlebrot points
// array of pointers to non divergent points
// split the array based on the number of workers
// go routines to increment

//TODO:
// restrict y and x iterations to be the same
// reduce memory use

type complex struct {
	X, Y float64
}

type mandlebrot_point struct {
	Z         complex
	C         complex
	iteration float64
	id        float64
}

type mandlebrot_plane struct {
	plane    []mandlebrot_point // the whole grid
	iterable [][]*float64       // points which did not yet diverge, each point is a 5 slice of pointers
}

var (
	min_Z     complex = complex{X: -2, Y: -2}
	max_Z     complex = complex{X: 2, Y: 2}
	increment int     = 1000
)

const (
	max_iterations int = 5000
	workers        int = 20
)

func (p *mandlebrot_plane) init_plane(min_Z complex, max_Z complex, increments int) {

	x_step := (max_Z.X - min_Z.X) / float64(increments)
	y_step := (max_Z.Y - min_Z.Y) / float64(increments)
	p.plane = make([]mandlebrot_point, increments*increments)

	p.iterable = make([][]*float64, increments*increments) // initialize slice

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ { //populate slices using goroutines
		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			for i := w; i < increments; i += workers { // rows
				for j := 0; j < increments; j++ { // cols
					point := mandlebrot_point{
						Z:         complex{0, 0},
						C:         complex{X: min_Z.X + float64(j)*x_step, Y: min_Z.Y + float64(i)*y_step},
						iteration: 0,
						id:        float64(i*increments + j)} // stopgap solution because the pointer for point[4] isnt working
					p.plane[i*increments+j] = point

					point_pointers := []*float64{&point.Z.X, &point.Z.Y, &point.C.X, &point.C.Y, &point.id}
					p.iterable[i*increments+j] = point_pointers
				}
			}
		}(w, &wg)
	}
	wg.Wait()

}

func (p *mandlebrot_plane) iterations(max_iterations int) {
	var wg sync.WaitGroup
	var divide_iterables [workers][][]*float64

	num_points := len(p.iterable)
	slice_size := num_points / workers

	worker_iteration := func(iterable [][]*float64) { //TODO: pointer issue, potentially assign pointers in this block
		defer wg.Done()
		Z := [2]float64{}
		C := [2]float64{}
		var next_x, next_y float64
		for i, point := range iterable { //next generation
			diverged := false
			Z = [2]float64{*point[0], *point[1]}
			C = [2]float64{*point[2], *point[3]}

			for j := 0; j < max_iterations && !diverged; j++ {

				next_x = Z[0]*Z[0] - Z[1]*Z[1] + C[0]
				next_y = 2*Z[0]*Z[1] + C[1]
				Z[0], Z[1] = next_x, next_y

				if Z[0]*Z[0]+Z[1]*Z[1] >= 4 { // remove the pointers of divergent points
					diverged = true
					//*point[4] = float64(j + 1) // this bit doesnt work for whatever reason so added id attribute to points
					p.plane[int(*point[4])].iteration = float64(j + 1)
					iterable[i] = nil
					break
				}
			}
			if !diverged { // set the final position
				*point[0] = Z[0]
				*point[1] = Z[1]
			}

		}
	}

	for i := 0; i < workers; i++ { // split work
		start := i * slice_size
		end := min(start+slice_size, num_points)
		divide_iterables[i] = p.iterable[start:end]
		wg.Add(1)
		go worker_iteration(divide_iterables[i])

	}
	wg.Wait()

	non_divergent := [][]*float64{} // clean up divergent points
	for _, pointers := range p.iterable {
		if pointers != nil {
			non_divergent = append(non_divergent, pointers)
		}
	}
	p.iterable = non_divergent
}

func (p *mandlebrot_plane) plot_to_png() {
	width, height := increment, increment
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	for pos, point := range p.plane {
		row := pos % increment
		col := pos / increment
		color_val := uint8(255 - (255 / (1 + 0.05*point.iteration)))

		img.Set(row, col, color.RGBA{color_val, color_val, color_val, 255})

	}

	file, err := os.Create("mandlebrot.png")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	mandlebrot_set := mandlebrot_plane{}

	init_time := time.Now()
	mandlebrot_set.init_plane(min_Z, max_Z, increment)
	fmt.Printf("Finished initialising the plane with %d points, it took %dms \n", len(mandlebrot_set.iterable), time.Since(init_time).Milliseconds())

	now := time.Now()
	mandlebrot_set.iterations(max_iterations)

	end := time.Since(now).Milliseconds()

	fmt.Printf("%d workers iterating %d times on %d points took %d ms \n", workers, max_iterations, increment*increment, end)

	plot_time := time.Now()

	mandlebrot_set.plot_to_png()
	fmt.Printf("Finished plotting, it took %dms \n	", time.Since(plot_time).Milliseconds())
}
