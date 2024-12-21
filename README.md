
![mandelbrot_example](https://github.com/user-attachments/assets/a5db090a-e2a5-47bc-b6b3-3d84b619391f)

![julia_example]()

# The Mandelbrot Set

The **Mandelbrot Set** is a famous fractal structure in the complex plane, defined by all points  $z_0 = a + bi$ for which the sequence $z_n$ does not diverge to infinity. 

The sequence $z_n$ is given by the following recursive relation:

$z_{n+1} = z_n^2 + z_0$

The **Julia Sets** are similarly defined, taking in an extra constant $c$ and following the relation for each point $z_0$ on the plane:

$z_{n+1} = z_n^2 +c$

## Description

A Mandlebrot / Julia set generator written in go. Taking in resolution, iterations, boundary of region, toggle for julia set and julia set iteration value from the user, and returning an image of the generated fractal with each pixel colored by the max iteration attained.

The algorithm runs concurrently with the implementation of goroutines and worker pools for an optimised runtime.

## Installation

`git clone` into your desired directory, then run `go build .`

## Usage

The program accepts user input in the form of standard input. An empty or incorrect response will assign the default value.

Example Use:

```bash

$ Mandlebrot

> Enter picture width resolution in pixels: (Default 2000px)
$ 1920

> Enter desired number of iterations: (Default 500)
$ 10000 

> Enter the bottom left bound for the image 2 number components seperated by a space: (Default: -2 -2)
$ -1.5 -1

> Enter the top right bound for the image as 2 number components seperated by a space: (Default: 2 2)
$ 1.5 1

> Enter 'true' or 'false' to render Julia set inplace of Mandlebrot: (Default is 'false')
$ true 

> Enter the C value for the Julia set as 2 number components seperated by a space: (Default: 0.35 0.35)
$ -0.8 0.156
```

## Prerequisites

Golang: This project requires the Go programming language. You can download it from https://golang.org/dl/.
