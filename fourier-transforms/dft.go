package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"scientificgo.org/fft"
)

func toComplex(a []int) []complex128 {
	result := make([]complex128, len(a))
	for i, v := range a {
		result[i] = complex(float64(v), 0.0)
	}
	return result
}

func omegaPowers(N int, inverse bool) []complex128 {
	// Since omega = e^(-2πi/N) has only N values (roots), calculate it once
	pow := -1i * math.Pi * 2 / complex(float64(N), 0.0)

	if inverse {
		pow *= complex(-1.0, 0.0)
	}

	omega := cmplx.Exp(pow)

	omegaPowers := make([]complex128, N)

	for i := range omegaPowers {
		omegaPowers[i] = cmplx.Pow(omega, complex(float64(i), 0.0))
	}
	return omegaPowers
}

func dft(arr []complex128, inverse bool) []complex128 {
	// mth term of dft = sum of( kth term of signal x omega^(k*m))
	// where omega = e^(-2πi/N)
	// N = size of samples in original input signal
	// for m from 0 to N-1, k from 0 to N-1
	size := len(arr)

	result := make([]complex128, size)
	omPow := omegaPowers(size, inverse)
	fmt.Println("OmegaPowers when", inverse, omPow)
	for i, _ := range arr {
		for j, v := range arr {
			result[i] += v * omPow[j*i%size]
		}
		if inverse {
			result[i] *= complex(1.0/float64(size), 0.0)
		}
	}
	return result
}

func main() {
	arr := toComplex([]int{2, 3, 4})
	fmt.Println("arr ", arr)

	dft_arr := dft(arr, false)
	fmt.Println("dft ", dft_arr)
	idft_arr := dft(dft_arr, true)
	fmt.Println("idft", idft_arr)

	// IDFT of result from library FFT
	fmt.Println("idft2", dft([]complex128{complex(9, 0), complex(-1.5, 0.8660254037844384), complex(-1.5, -0.8660254037844384)}, true))

	fft_arr := fft.Fft(arr, false)
	fmt.Println("fft ", fft_arr)
	fmt.Println("ifft", fft.Fft(fft_arr, true))
}
