// Stefan Nilsson 2013-02-27
// modified by Vilhelm Prytz 2023-03-14

// original runtime
// $ time go run julia.go
// go run julia.go  13,56s user 0,17s system 101% cpu 13,504 total

// modified runtime
// $ time go run julia.go
// go run julia.go  11,99s user 0,20s system 143% cpu 8,466 total

// This program creates pictures of Julia sets (en.wikipedia.org/wiki/Julia_set).
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"
	"strconv"
	"sync"
)

type ComplexFunc func(complex128) complex128

var Funcs []ComplexFunc = []ComplexFunc{
	func(z complex128) complex128 { return z*z - 0.61803398875 },
	func(z complex128) complex128 { return z*z + complex(0, 1) },
	func(z complex128) complex128 { return z*z + complex(-0.835, -0.2321) },
	func(z complex128) complex128 { return z*z + complex(0.45, 0.1428) },
	func(z complex128) complex128 { return z*z*z + 0.400 },
	func(z complex128) complex128 { return cmplx.Exp(z*z*z) - 0.621 },
	func(z complex128) complex128 { return (z*z+z)/cmplx.Log(z) + complex(0.268, 0.060) },
	func(z complex128) complex128 { return cmplx.Sqrt(cmplx.Sinh(z*z)) + complex(0.065, 0.122) },
}

func main() {
	// setup waitgroup
	wg := new(sync.WaitGroup)
	wg.Add(len(Funcs))
	for n, fn := range Funcs {
		go CreatePng("picture-"+strconv.Itoa(n)+".png", fn, 1024, wg)
	}

	// wait for all goroutines to finish
	wg.Wait()
}

// CreatePng creates a PNG picture file with a Julia image of size n x n.
func CreatePng(filename string, f ComplexFunc, n int, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	err = png.Encode(file, Julia(f, n))
	return
}

// Julia returns an image of size n x n of the Julia set for f.
func Julia(f ComplexFunc, n int) image.Image {
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	img := image.NewRGBA(bounds)
	s := float64(n / 4)

	var wg sync.WaitGroup
	for outerI := bounds.Min.X; outerI < bounds.Max.X; outerI++ {
		for outerJ := bounds.Min.Y; outerJ < bounds.Max.Y; outerJ++ {
			wg.Add(1)
			go func(i, j int) {
				n := Iterate(f, complex(float64(i)/s, float64(j)/s), 256)
				r := uint8(0)
				g := uint8(0)
				b := uint8(n % 32 * 8)
				img.Set(i, j, color.RGBA{r, g, b, 255})
				wg.Done()
			}(outerI, outerJ)
		}
	}
	wg.Wait()
	return img
}

// Iterate sets z_0 = z, and repeatedly computes z_n = f(z_{n-1}), n â‰¥ 1,
// until |z_n| > 2  or n = max and returns this n.
func Iterate(f ComplexFunc, z complex128, max int) (n int) {
	for ; n < max; n++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			break
		}
		z = f(z)
	}
	return
}
