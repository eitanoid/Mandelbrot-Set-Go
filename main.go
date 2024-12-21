package main

import (
	"fmt"
	"runtime"
	"time"
)

// set plane of complex mandelbrot points
// array of pointers to non divergent points
// split the array based on the number of workers
// go routines to iterate the points by chunks

//TODO:
// user can input color function
// dynamic visualisation with a visual library like raylib or turn to gif

var (
	workers int // can limit max goroutines
)

func main() {
	// handle user input
	res_in, iter_in, gif_frame_iterations_in, gifdelay_in, lx_in, ly_in, ux_in, uy_in, julia_in, julia_x, julia_y := user_input(2000, 500, 20, -2, -2, 2, 2, false, 0.35, 0.35)
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

	runtime.GOMAXPROCS(runtime.NumCPU()) // use all cpu cores
	workers = runtime.NumCPU()
	mandelbrot_set := Mandelbrot_Plane{}

	init_time := time.Now()
	mandelbrot_set.Init_plane(min_Z, x_steps, y_steps, step_size, julia_in, julia_x, julia_y)
	fmt.Printf("Initialized %d points taking %dms \n", len(mandelbrot_set.Iterable), time.Since(init_time).Milliseconds())

	if gif_frame_iterations_in != 0 { // gif if non 0 image otherwise
		mandelbrot_set.Plot_to_GIF(gif_frame_iterations_in, max_iterations, gifdelay_in)
	} else {
		now := time.Now()
		mandelbrot_set.Iterations(max_iterations)
		end := time.Since(now).Milliseconds()
		fmt.Printf("%d workers completed %d iterations on %d points in %d ms \n", workers, max_iterations, x_steps*y_steps, end)

		plot_time := time.Now()
		Save_Image(mandelbrot_set.Plot_to_Image(max_iterations), "output.png")
		fmt.Printf("Finished plotting, it took %dms \n	", time.Since(plot_time).Milliseconds())

	}
}
