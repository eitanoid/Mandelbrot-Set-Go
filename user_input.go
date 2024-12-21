package main

import (
	"bufio"
	"fmt"
	"os"
)

// return resolution iterations gif_iterations gif_delay lx ly ux uy julia jx jy
func user_input(res_default int, iter_default int, gifdelay_default int, lx_default float64, ly_default float64, ux_default float64, uy_default float64, julia_default bool, jx_default float64, jy_default float64) (int, int, int, int, float64, float64, float64, float64, bool, float64, float64) {
	var res_input, iter_input, frameiter_input, gifdelay_input int
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

	fmt.Println("If non-zero will create a gif, enter the number of iterations per frame: (Default 0)")
	scanner.Scan()
	arg = scanner.Text()
	_, err = fmt.Sscanf(arg, "%d", &frameiter_input)
	if err != nil { // set default and remove the rest of a bad string.
		frameiter_input = 0
	}

	if frameiter_input != 0 {
		fmt.Println("Enter GIF frame delay measured in 100th of a second: (Default 20)")
		scanner.Scan()
		arg = scanner.Text()
		_, err = fmt.Sscanf(arg, "%d", &gifdelay_input)
		if err != nil { // set default and remove the rest of a bad string.
			gifdelay_input = gifdelay_default
		}
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
	return res_input, iter_input, frameiter_input, gifdelay_input, lx_input, ly_input, ux_input, uy_input, julia_input, jx_input, jy_input
}
