package main

import (
	"bufio"
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

// set plane of complex mandelbrot points
// array of pointers to non divergent points
// split the array based on the number of workers
// go routines to iterate the points by chunks

//TODO:
// user can input color function
// dynamic visualisation with a visual library like raylib or turn to gif

type Complex struct {
	X, Y float64
}

var (
	julia_c = Complex{X: 0.35, Y: 0.35} // an interesting julia set value
)

type Mandelbrot_Point struct {
	Z         Complex
	C         Complex
	iteration float64
}

type Mandelbrot_Plane struct {
	plane    [][]Mandelbrot_Point // the whole grid
	iterable []*Mandelbrot_Point  // slice of pointers to array
}

var ( //user info
	res_in  int
	iter_in int
	lx_in   float64
	ly_in   float64
	ux_in   float64
	uy_in   float64
	julia   bool    = false
	j_x     float64 // c value used for the julia set c = Complex{j_x,j_y}
	j_y     float64
)

func user_input(res_default int, iter_default int, lx_default float64, ly_default float64, ux_default float64, uy_default float64, julia_default bool, jx_default float64, jy_default float64) (int, int, float64, float64, float64, float64, bool, float64, float64) {
	var res_input, iter_input int
	var lx_input, ly_input, ux_input, uy_input float64
	var julia_input bool
	var jx_input, jy_input float64

	var arg string // current arg
	var err error

	scanner := bufio.NewScanner(os.Stdin)

	// NOTE: these input functions dont respond to empty strings. Maybe scan into a string and process it?

	fmt.Println("Enter picture width resolution in pixels: (Default 2000px)")
	scanner.Scan()
	arg = scanner.Text()
	_, err = fmt.Sscanf(arg, "%d", &res_input)
	if err != nil { // set default and remove the rest of a bad string.
		res_input = res_default
	}

	fmt.Println("Enter desired number of iterations: (Default 500)")
	scanner.Scan()
	arg = scanner.Text()
	_, err = fmt.Sscanf(arg, "%d", &iter_input)
	if err != nil { // set default and remove the rest of a bad string.
		iter_input = iter_default
	}

	fmt.Println("Enter the bottom left bound for the image 2 number components seperated by a space: (Default: -2 -2)")
	scanner.Scan()
	arg = scanner.Text()
	_, err = fmt.Sscanf(arg, "%f %f", &lx_input, &ly_input)
	if err != nil { // set default and remove the rest of a bad string.
		lx_input = lx_default
		ly_input = ly_default
	}

	fmt.Println("Enter the top right bound for the image as 2 number components seperated by a space: (Default: 2 2)")
	scanner.Scan()
	arg = scanner.Text()
	_, err = fmt.Sscanf(arg, "%f %f", &ux_input, &uy_input)
	if err != nil { // set default and remove the rest of a bad string.
		ux_input = ux_default
		uy_input = uy_default
	}

	fmt.Println("Enter 'true' or 'false' to render Julia set inplace of Mandelbrot: (Default is 'false')")
	scanner.Scan()
	arg = scanner.Text()
	_, err = fmt.Sscanf(arg, "%t", &julia_input)
	if err != nil { // set default and remove the rest of a bad string.
		julia_input = julia_default
	}

	if julia_input {
		fmt.Println("Enter the C value for the Julia set as 2 number components seperated by a space: (Default: 0.35 0.35)")
		scanner.Scan()
		arg = scanner.Text()
		_, err = fmt.Sscanf(arg, "%f %f", &jx_input, &jy_input)
		if err != nil { // set default and remove the rest of a bad string.
			jx_input = jx_default
			jy_input = jy_default
		}
	}
	// fmt.Println("Enter picture width resolution in pixels: (Default 2000px)")
	// 	_, err = fmt.Scanf("%d", &res_input)
	// 	if err != nil { // set default and remove the rest of a bad string.
	// 		fmt.Scan(&void)
	// 		res_input = res_default
	// 	}

	// 	fmt.Println("Enter desired number of iterations: (Default 500)")
	// 		_, err = fmt.Scanf("%d", &iter_input)
	// 		if err != nil { // set default and remove the itert of a bad string.
	// 			fmt.Scan(&void)
	// 			iter_input = iter_default
	// 		}

	// fmt.Println("Enter the bottom left bound for the image as a complex number seperated by a space: (Default: -2 -2)")
	// _, err = fmt.Scanf("%f %f", &lx_input, &ly_input)
	// if err != nil { // set default and remove the rest of a bad string.
	// 	fmt.Scan(&void)
	// 	lx_input = lx_default
	// 	ly_input = ly_default
	// }

	// fmt.Println("Enter the top right bound for the image as a complex number seperated by a space: (Default: 2 2)")
	// _, err = fmt.Scanf("%f %f", &ux_input, &uy_input)
	// if err != nil { // set default and remove the rest of a bad string.
	// 	fmt.Scan(&void)
	// 	ux_input = ux_default
	// 	uy_input = uy_default
	// }

	return res_input, iter_input, lx_input, ly_input, ux_input, uy_input, julia_input, jx_input, jy_input
}

var (
	workers int
)

func (p *Mandelbrot_Plane) Init_plane(min_Z Complex, x_steps int, y_steps int, step_size float64, julia bool, julia_x float64, julia_y float64) {

	p.plane = make([][]Mandelbrot_Point, y_steps)

	p.iterable = make([]*Mandelbrot_Point, x_steps*y_steps) // initialize slice

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ { //populate slices using goroutines

		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			for y := w; y < y_steps; y += workers { // rows
				p.plane[y] = make([]Mandelbrot_Point, x_steps)
				for x := 0; x < x_steps; x++ { // cols
					var point Mandelbrot_Point
					if julia {
						point = Mandelbrot_Point{
							Z:         Complex{X: min_Z.X + float64(x)*step_size, Y: min_Z.Y + float64(y)*step_size},
							C:         Complex{julia_x, julia_y},
							iteration: 0}
					} else { // mandelbrot
						point = Mandelbrot_Point{
							Z:         Complex{0, 0},
							C:         Complex{X: min_Z.X + float64(x)*step_size, Y: min_Z.Y + float64(y)*step_size},
							iteration: 0}
					}
					p.plane[y][x] = point
					p.iterable[y*x_steps+x] = &p.plane[y][x]
				}
			}
		}(w, &wg)
	}
	wg.Wait()

}

func (p *Mandelbrot_Plane) Iterations(max_iterations int) {
	chunk_size := 10000
	num_points := len(p.iterable)
	work_queue := make(chan []*Mandelbrot_Point)
	var progress atomic.Int32
	var wg sync.WaitGroup

	worker_iteration := func(iterable []*Mandelbrot_Point) {
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
				point.Z = Complex{Z[0], Z[1]}

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

	non_divergent := []*Mandelbrot_Point{} // clean up divergent points
	for _, point := range p.iterable {
		if point != nil {
			non_divergent = append(non_divergent, point)
		}
	}
	p.iterable = non_divergent
}

func (p *Mandelbrot_Plane) Plot_to_PNG(x_steps int, y_steps int) {
	width, height := x_steps, y_steps
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	for y, row := range p.plane {
		// var speed uint8 = 1 // arbitrary number
		for x, point := range row {
			// iterations := point.iteration
			var (
				red   uint8 = 1 // (speed * uint8(iterations)) % 255 // multily the colir intensity by color values out of phase (255/3 = 85).
				blue  uint8 = 1 // (speed*uint8(iterations) + 85) % 255
				green uint8 = 1 // (speed*uint8(iterations) + 85*2) % 255
			)
			color_val := uint8(255 - (255 / (1 + 0.05*point.iteration)))
			img.Set(x, y_steps-y, color.RGBA{color_val * red, color_val * blue, color_val * green, 255})
		}
	}

	file, err := os.Create("mandelbrot.png")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

}

func main() {

	res_in, iter_in, lx_in, ly_in, ux_in, uy_in, julia_in, julia_x, julia_y := user_input(2000, 500, -2, -2, 2, 2, false, 0.35, 0.35)
	fmt.Println(res_in, iter_in, lx_in, ly_in, ux_in, uy_in)
	min_Z := Complex{X: lx_in, Y: ly_in}
	x_len := ux_in - lx_in
	y_len := uy_in - ly_in
	x_steps := res_in
	max_iterations := iter_in
	step_size := float64(x_len) / float64(x_steps)
	y_steps := int(y_len / step_size)

	if lx_in >= ux_in || ly_in >= uy_in { // check user input
		fmt.Println("upper bound must be larger than lower bound.")
		return
	}

	if max_iterations <= 0 || res_in <= 0 {
		fmt.Println("integers must be positive")
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	workers = runtime.NumCPU()
	mandelbrot_set := Mandelbrot_Plane{}

	init_time := time.Now()
	mandelbrot_set.Init_plane(min_Z, x_steps, y_steps, step_size, julia_in, julia_x, julia_y)
	fmt.Printf("Initialized %d points taking %dms \n", len(mandelbrot_set.iterable), time.Since(init_time).Milliseconds())

	now := time.Now()
	mandelbrot_set.Iterations(max_iterations)
	end := time.Since(now).Milliseconds()
	fmt.Printf("%d workers completed %d iterations on %d points in %d ms \n", workers, max_iterations, x_steps*y_steps, end)

	plot_time := time.Now()
	mandelbrot_set.Plot_to_PNG(x_steps, y_steps)
	fmt.Printf("Finished plotting, it took %dms \n	", time.Since(plot_time).Milliseconds())
}
