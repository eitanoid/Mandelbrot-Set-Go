
| Mandelbrot | Julia | Julia GIF |
| :-         | :-    | :-        |
| ![mandelbrot_example](https://github.com/eitanoid/Mandelbrot-Set-Go/blob/main/pictures/mandelbrot.png) | ![julia_example](https://github.com/eitanoid/Mandlebrot-Set-Go/blob/main/pictures/julia.png) | ![gif example](https://github.com/eitanoid/Mandelbrot-Set-Go/blob/main/pictures/julia_gif.gif) |

# The Mandelbrot Set

The **Mandelbrot Set** is a famous fractal structure in the complex plane, defined by all points  $z_0 = a + bi$ for which the sequence $z_n$ does not diverge to infinity. 

The sequence $z_n$ is given by the following recursive relation:

$z_{n+1} = z_n^2 + z_0$

The **Julia Sets** are similarly defined, taking in an extra constant $c$ and following the relation for each point $z_0$ on the plane:

$z_{n+1} = z_n^2 +c$

## Description

A Mandelbrot / Julia set generator written in go. Taking in resolution, iterations, boundary of region, toggle for julia set and julia set iteration value from the user, and returning an image of the generated fractal with each pixel colored by the max iteration attained.

The program can also output a GIF, taking in the maximum iterations desired, the number of iterations per frame of the GIF and the delay between GIF frames.

The algorithm runs concurrently with the implementation of goroutines and worker pools for an optimised runtime.

## Installation

`git clone` into your desired directory, then run `go build .`

## Usage

The program accepts user input in the form of standard input. An empty or incorrect response will assign the default value.

Example Use:

```bash

$ Mandelbrot

> Enter picture width resolution in pixels: (Default 2000)
$ 1920

> Enter desired number of iterations: (Default 500)
$ 10000 

> If non-zero will create a gif, enter the number of iterations per frame: (Default 0)
$ 0

#If above is non-zero the following is asked:
> Enter GIF frame delay measured in 100th of a second: (Default 20)
$ 30

> Enter the bottom left bound for the image 2 number components seperated by a space: (Default: -2 -2)
$ -1.5 -1

> Enter the top right bound for the image as 2 number components seperated by a space: (Default: 2 2)
$ 1.5 1

> Enter 'true' or 'false' to render Julia set inplace of Mandelbrot: (Default is 'false')
$ true 

# If answered true to previous questions asks this one:
> Enter the C value for the Julia set as 2 number components seperated by a space: (Default: 0.35 0.35)
$ -0.8 0.156
```

## Warning

This project will generate gifs at any size, however at present no compression is used, which can result in MASSIVE file sizes at large gif resolutions.


## Prerequisites

Golang: This project requires the Go programming language. You can download it from https://golang.org/dl/.
