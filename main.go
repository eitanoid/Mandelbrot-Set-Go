package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
	"sync"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// set plane of complex mandlebrot points
// array of pointers to non divergent points
// split the array based on the number of workers
// go routines to increment

type complex struct {
	X, Y float64
}

type mandlebrot_point struct {
	Z         complex
	C         complex
	iteration float64
}

type mandlebrot_plane struct {
	plane    map[mandlebrot_point]struct{} // the whole grid
	iterable [][]*float64                  // points which did not yet diverge
}

var (
	min_Z     complex = complex{X: -2, Y: -1.5}
	max_Z     complex = complex{X: 1, Y: 1.5}
	increment int     = 3000
)

const (
	max_iterations int = 6000
	workers        int = 20
)

func (p *mandlebrot_plane) init_plane(min_Z complex, max_Z complex, increments int) {

	x_step := (max_Z.X - min_Z.X) / float64(increments)
	y_step := (max_Z.Y - min_Z.Y) / float64(increments)
	p.plane = make(map[mandlebrot_point]struct{})

	for i := 0; i < increments; i++ { // rows
		for j := 0; j < increments; j++ { // cols
			point := mandlebrot_point{
				Z:         complex{0, 0},
				C:         complex{X: min_Z.X + float64(j)*x_step, Y: min_Z.Y + float64(i)*y_step},
				iteration: 0}
			p.plane[point] = struct{}{}

			point_pointers := []*float64{&point.Z.X, &point.Z.Y, &point.C.X, &point.C.Y, &point.iteration}
			p.iterable = append(p.iterable, point_pointers)
		}
	}

}

func (p *mandlebrot_plane) iterations(iterations int) {
	var wg sync.WaitGroup
	var divide_iterables [workers][][]*float64

	num_points := len(p.iterable)
	slice_size := num_points / workers

	worker_iteration := func(iterable [][]*float64) {
		defer wg.Done()

		Z := [2]float64{}
		C := [2]float64{}
		var next_x, next_y float64
		for i, point := range iterable { //next generation
			diverged := false
			Z = [2]float64{*point[0], *point[1]}
			C = [2]float64{*point[2], *point[3]}

			for j := 0; j < iterations; j++ {

				next_x = Z[0]*Z[0] - Z[1]*Z[1] + C[0]
				next_y = 2*Z[0]*Z[1] + C[1]
				Z[0], Z[1] = next_x, next_y

				if next_x*next_x+next_y*next_y >= 4 { // remove the pointers of divergent points
					diverged = true
					*point[4] = *point[4] + float64(j)
					break
				}
			}
			if !diverged {
				*point[0] = Z[0]
				*point[1] = Z[1]
			} else {
				iterable[i] = nil
			}

		}
	}

	for i := 0; i < workers; i++ {
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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	mandlebrot_set := mandlebrot_plane{}
	mandlebrot_set.init_plane(min_Z, max_Z, increment)

	now := time.Now()
	mandlebrot_set.iterations(max_iterations)
	//fmt.Printf("Iterated, %d points remain \n", len(mandlebrot_set.iterable))

	end := time.Since(now).Milliseconds()
	fmt.Printf("%d workers iterating %d times on %d points took %d ms \n", workers, max_iterations, increment*increment, end)

	plot := func(mandlebrot_set mandlebrot_plane) {
		p := plot.New() // plotting

		p.Title.Text = "Mandelbrot Set"
		p.X.Label.Text = "Re"
		p.Y.Label.Text = "Im"

		points := make(plotter.XYs, 0)

		var real, imag float64
		for _, pointers := range mandlebrot_set.iterable { // add data
			real = *pointers[2]
			imag = *pointers[3]
			points = append(points, plotter.XY{X: real, Y: imag})
		}

		// Create a scatter plot.
		s, err := plotter.NewScatter(points)
		if err != nil {
			log.Fatalf("plotter.NewScatter() failed: %v", err)
		}
		s.GlyphStyle.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255} // white
		s.GlyphStyle.Radius = 0.15
		p.Add(s)

		p.X.Min, p.Y.Min = -2, -1
		p.X.Max, p.Y.Max = 1, 2

		// Save the plot to a PNG file.
		if err := p.Save(64*vg.Inch, 64*vg.Inch, "mandelbrot.png"); err != nil {
			log.Fatalf("p.Save() failed: %v", err)
		}

	}
	plot(mandlebrot_set)
}
