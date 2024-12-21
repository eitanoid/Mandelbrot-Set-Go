package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Complex struct {
	X, Y float64
}

type Mandelbrot_Point struct {
	Z         Complex
	C         Complex
	Iteration float64
}

type Mandelbrot_Plane struct {
	Plane    [][]Mandelbrot_Point // the whole grid
	Iterable []*Mandelbrot_Point  // slice of pointers to array
}

func (p *Mandelbrot_Plane) Init_plane(min_Z Complex, x_steps int, y_steps int, step_size float64, julia bool, julia_x float64, julia_y float64) {

	p.Plane = make([][]Mandelbrot_Point, y_steps)

	p.Iterable = make([]*Mandelbrot_Point, x_steps*y_steps) // initialize slice

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ { //populate slices using goroutines

		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			for y := w; y < y_steps; y += workers { // rows
				p.Plane[y] = make([]Mandelbrot_Point, x_steps)
				for x := 0; x < x_steps; x++ { // cols
					var point Mandelbrot_Point
					if julia {
						point = Mandelbrot_Point{
							Z:         Complex{X: min_Z.X + float64(x)*step_size, Y: min_Z.Y + float64(y)*step_size},
							C:         Complex{julia_x, julia_y},
							Iteration: 0}
					} else { // mandelbrot
						point = Mandelbrot_Point{
							Z:         Complex{0, 0},
							C:         Complex{X: min_Z.X + float64(x)*step_size, Y: min_Z.Y + float64(y)*step_size},
							Iteration: 0}
					}
					p.Plane[y][x] = point
					p.Iterable[y*x_steps+x] = &p.Plane[y][x]
				}
			}
		}(w, &wg)
	}
	wg.Wait()

}

func (p *Mandelbrot_Plane) Iterations(max_iterations int) {
	chunk_size := 10000
	num_points := len(p.Iterable)
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
					point.Iteration += float64(j + 1)
					iterable[i] = nil
					break
				}
			}
			if !diverged { // set the final position
				point.Z = Complex{Z[0], Z[1]}
				point.Iteration += float64(max_iterations)

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
		work_queue <- p.Iterable[start:end]

		percent := int(progress.Load()*10) / (num_points / chunk_size)

		if percent != current_progress {
			fmt.Printf("Iterations are %d %s complete  \n", percent*10, "%")
			current_progress = percent
		}

	}

	close(work_queue)
	wg.Wait()

	non_divergent := []*Mandelbrot_Point{} // clean up divergent points
	for _, point := range p.Iterable {
		if point != nil {
			non_divergent = append(non_divergent, point)
		}
	}
	p.Iterable = non_divergent
}
