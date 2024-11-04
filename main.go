package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"sync"
	"time"
)

// set plane of complex mandlebrot points
// array of pointers to non divergent points
// split the array based on the number of workers
// go routines to increment

//TODO:
// fix not working for rectangles

type complex struct {
	X, Y float64
}

type mandlebrot_point struct {
	Z         complex
	C         complex
	iteration float64
}

type mandlebrot_plane struct {
	plane    []mandlebrot_point  // the whole grid
	iterable []*mandlebrot_point // slice of pointers to array
}

var (
	rinput = flag.Int("r", 2000, "Set the resolution")
	iinput = flag.Int("i", 500, "Set the number of iterations")
)
var (
	min_Z          complex = complex{X: -2, Y: -2}
	step_size      float64
	x_steps        int //int(x_len / step_size)
	y_steps        int //int(y_len / step_size)
	max_iterations int
	workers        int
)

func (p *mandlebrot_plane) init_plane(min_Z complex, x_steps int, y_steps int) {

	p.plane = make([]mandlebrot_point, x_steps*y_steps)

	p.iterable = make([]*mandlebrot_point, x_steps*y_steps) // initialize slice

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ { //populate slices using goroutines
		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			for y := w; y < y_steps; y += workers { // rows
				for x := 0; x < x_steps; x++ { // cols
					point := mandlebrot_point{
						Z:         complex{0, 0},
						C:         complex{X: min_Z.X + float64(x)*step_size, Y: min_Z.Y + float64(y)*step_size},
						iteration: 0}
					p.plane[y*y_steps+x] = point
					p.iterable[y*y_steps+x] = &p.plane[y*y_steps+x]
				}
			}
		}(w, &wg)
	}
	wg.Wait()

}

func (p *mandlebrot_plane) iterations(max_iterations int) {
	chunk_size := 10000
	num_points := len(p.iterable)
	work_queue := make(chan []*mandlebrot_point)
	var wg sync.WaitGroup

	worker_iteration := func(iterable []*mandlebrot_point) {
		Z := [2]float64{}
		C := [2]float64{}
		var next_x, next_y float64

		for i, point := range iterable { //next generation
			diverged := false
			Z = [2]float64{point.Z.X, point.Z.Y}
			C = [2]float64{point.C.X, point.C.Y}

			for j := 0; j < max_iterations && !diverged; j++ {

				next_x = Z[0]*Z[0] - Z[1]*Z[1] + C[0]
				next_y = 2*Z[0]*Z[1] + C[1]
				Z[0], Z[1] = next_x, next_y

				if Z[0]*Z[0]+Z[1]*Z[1] >= 4 { // remove the pointers of divergent points
					diverged = true
					point.iteration += float64(j + 1)
					iterable[i] = nil
					break
				}
			}
			if !diverged { // set the final position
				point.Z = complex{Z[0], Z[1]}

			}
		}
	}

	for i := 0; i < workers; i++ { // worker pool
		wg.Add(1)
		go func() {
			defer wg.Done()
			for points, ok := <-work_queue; ok; points, ok = <-work_queue {
				worker_iteration(points)
			}
		}()
	}

	for i, start, end := 0, 0, 0; end < num_points; i++ { // split work
		start = i * chunk_size
		end = min(start+chunk_size, num_points)
		work_queue <- p.iterable[start:end]

	}

	close(work_queue)
	wg.Wait()

	non_divergent := []*mandlebrot_point{} // clean up divergent points
	for _, point := range p.iterable {
		if point != nil {
			non_divergent = append(non_divergent, point)
		}
	}
	p.iterable = non_divergent
}

func (p *mandlebrot_plane) plot_to_png() {
	width, height := x_steps, y_steps
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	for pos, point := range p.plane {
		row := pos % y_steps
		col := pos / x_steps
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
	flag.Parse()
	x_steps = *rinput
	y_steps = x_steps
	max_iterations = *iinput
	step_size = float64(4) / float64(x_steps)

	runtime.GOMAXPROCS(runtime.NumCPU())
	workers = runtime.NumCPU()
	mandlebrot_set := mandlebrot_plane{}

	init_time := time.Now()
	mandlebrot_set.init_plane(min_Z, x_steps, y_steps)
	fmt.Printf("Initialized %d points taking %dms \n", len(mandlebrot_set.iterable), time.Since(init_time).Milliseconds())

	now := time.Now()
	mandlebrot_set.iterations(max_iterations)

	end := time.Since(now).Milliseconds()

	fmt.Printf("%d workers completed %d iterations on %d points in %d ms \n", workers, max_iterations, x_steps*y_steps, end)

	plot_time := time.Now()

	mandlebrot_set.plot_to_png()
	fmt.Printf("Finished plotting, it took %dms \n	", time.Since(plot_time).Milliseconds())
}
