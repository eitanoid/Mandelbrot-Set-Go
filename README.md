
![mandlebrot_cropped](https://github.com/user-attachments/assets/a5db090a-e2a5-47bc-b6b3-3d84b619391f)

# The Mandelbrot Set

The **Mandelbrot Set** is a famous fractal structure in the complex plane, defined by all points  $c = a + bi$ for which the sequence $Z_n$ does not diverge to infinity. 

The sequence $Z_n$ is given by the following recursive relation:

$Z_0 = 0$

$Z_{n+1} = Z_n^2 + c$

## Decription

A mandlebrot set generator using Go taking in a resolution and a max iteration count and outputting a PNG file of the mandlebrot set at that depth with a brightness adjusted based on the escape iteration of each point.

The algorithm runs concurrently with the implementation of goroutines and worker pools.

## Usage

You can pass the following flags to the executable:

```bash

usage: Mandlebrot [<flags>]

Flags:
-r    The resolution of the plot (default 2000)
-i    The number of iterations (default 500)

```
