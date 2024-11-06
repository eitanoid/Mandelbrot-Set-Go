
![mandlebrot_cropped](https://github.com/user-attachments/assets/a5db090a-e2a5-47bc-b6b3-3d84b619391f)

# The Mandelbrot Set

The **Mandelbrot Set** is a famous fractal structure in the complex plane, defined by all points  $c = a + bi$ for which the sequence $Z_n$ does not diverge to infinity. 

The sequence $Z_n$ is given by the following recursive relation:

$Z_0 = 0$

$Z_{n+1} = Z_n^2 + c$

## Description

A mandlebrot set generator using Go taking in a user inputs for a region in the complex plane, resolution and a number of iterations and outputting a PNG file of the mandlebrot set at that depth with a brightness adjusted based on the escape iteration of each point.

The algorithm runs concurrently with the implementation of goroutines and worker pools for an optimised runtime.

## Usage

You can pass the following flags to the executable:

```bash

$Mandlebrot -h

usage of mandlebrot:

-r    The width resolution of the plot (default 2000)
-i    The number of iterations (default 500)
-lx   The left x bound of the image (default -2)
-ly   The bottom y bound of the image (default -2)
-ux   The right x bound of the image (default 2)
-uy   The top y bound of the image (default 2)

```

## Prerequisites

Golang: This project requires the Go programming language. You can download it from https://golang.org/dl/.
