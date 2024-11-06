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
	"sync/atomic"
	"time"
)

// set plane of complex mandlebrot points
// array of pointers to non divergent points
// split the array based on the number of workers
// go routines to iterate the points by chunks

//TODO:
//

type complex struct {
	X, Y float64
}

type mandlebrot_point struct {
	Z         complex
	C         complex
	iteration float64
}

type mandlebrot_plane struct {
	plane    [][]mandlebrot_point // the whole grid
	iterable []*mandlebrot_point  // slice of pointers to array
}

var ( // user input
	rinput     = flag.Int("r", 2000, "Set the picture resolution")
	iinput     = flag.Int("i", 500, "Set the number of iterations")
	minx_input = flag.Float64("lx", -2, "Set the lower x bound for the picture")
	miny_input = flag.Float64("ly", -2, "Set the lower y bound for the picture")
	maxx_input = flag.Float64("ux", 2, "Set the upper x bound for the picture")
	maxy_input = flag.Float64("uy", 2, "Set the upper y bound for the picture")
)
var (
	workers int
)

func (p *mandlebrot_plane) init_plane(min_Z complex, x_steps int, y_steps int, step_size float64) {

	p.plane = make([][]mandlebrot_point, y_steps)

	p.iterable = make([]*mandlebrot_point, x_steps*y_steps) // initialize slice

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ { //populate slices using goroutines

		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			for y := w; y < y_steps; y += workers { // rows
				p.plane[y] = make([]mandlebrot_point, x_steps)
				for x := 0; x < x_steps; x++ { // cols
					point := mandlebrot_point{
						Z:         complex{0, 0},
						C:         complex{X: min_Z.X + float64(x)*step_size, Y: min_Z.Y + float64(y)*step_size},
						iteration: 0}
					p.plane[y][x] = point
					p.iterable[y*x_steps+x] = &p.plane[y][x]
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
	var progress atomic.Int32
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
		progress.Add(1)
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

	var current_progress int = 0
	for i, start, end := 0, 0, 0; end < num_points; i++ { // split work
		start = i * chunk_size
		end = min(start+chunk_size, num_points)
		work_queue <- p.iterable[start:end]

		percent := int(progress.Load()*10) / (num_points / chunk_size)

		if percent != current_progress {
			fmt.Printf("Iterations are %d %s complete  \n", percent*10, "%")
			current_progress = percent
		}

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

func (p *mandlebrot_plane) plot_to_png(x_steps int, y_steps int) {
	width, height := x_steps, y_steps
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	for y, row := range p.plane {
		for x, point := range row {
			color_val := uint8(255 - (255 / (1 + 0.05*point.iteration)))
			img.Set(x, y_steps-y, color.RGBA{color_val, color_val, color_val, 255})
		}
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

	flag.Parse() //accept user input

	min_Z := complex{X: *minx_input, Y: *miny_input}
	x_len := *maxx_input - *minx_input
	y_len := *maxy_input - *miny_input
	x_steps := *rinput
	max_iterations := *iinput
	step_size := float64(x_len) / float64(x_steps)
	y_steps := int(y_len / step_size)

	if *minx_input >= *maxx_input || *miny_input >= *maxy_input { // check user input
		fmt.Println("upper bound must be larger than lower bound.")
		return
	}

	if max_iterations <= 0 || *rinput <= 0 {
		fmt.Println("integers must be positive")
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	workers = runtime.NumCPU()
	mandlebrot_set := mandlebrot_plane{}

	init_time := time.Now()
	mandlebrot_set.init_plane(min_Z, x_steps, y_steps, step_size)
	fmt.Printf("Initialized %d points taking %dms \n", len(mandlebrot_set.iterable), time.Since(init_time).Milliseconds())

	now := time.Now()
	mandlebrot_set.iterations(max_iterations)
	end := time.Since(now).Milliseconds()
	fmt.Printf("%d workers completed %d iterations on %d points in %d ms \n", workers, max_iterations, x_steps*y_steps, end)

	plot_time := time.Now()
	mandlebrot_set.plot_to_png(x_steps, y_steps)
	fmt.Printf("Finished plotting, it took %dms \n	", time.Since(plot_time).Milliseconds())
}
