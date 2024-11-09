package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"scientificgo.org/fft"
)

/*
for inverse, omega = e ^(i2π)/N // change in sign in power
*/
func omegaPowers(size int, inverse bool) []complex128 {
	result := make([]complex128, size)
	pow := -2i * math.Pi/ complex(float64(size), 0.0)
	if inverse {
		pow = -pow
	}
	omega := cmplx.Exp(pow)
	for i := range result {
		result[i] = cmplx.Pow(omega, complex(float64(i), 0.0))
	}

	return result
}

/*
kth term of DFT = Summation (m from 0 to N-1) mth term of input * omega^(km)
Where omega = e^(-2πi/N) // - sign for anticlockwise

for fast fourier transform, radix 2 method assumes N is a multiple of 2
kth term = Ek + omega^k * Ok
(k + N/2)th term = Ek - omega^k * Ok
where Ek = kth term of FFT of original even terms
and Ok = kth term of FFT of original odd terms
And k = 0 to (N/2) - 1
*/
func myfft(arr []complex128, inverse bool) []complex128 {
	if len(arr) == 1 {
		return arr
	}

	var even, odd []complex128

	for i:= range arr {
		if i %2 == 0 {
			even = append(even, arr[i])
		} else {
			odd = append(odd, arr[i])
		}
	}

	even = myfft(even,false)
	odd = myfft(odd,false)

	//fmt.Println("even",even,"odd", odd)

	N := len(arr)

	omegaPowers := omegaPowers(N, inverse)
	result := make([]complex128, N)
	//fmt.Println("Calculating result", (N/2) - 1))
	//fmt.Println("omega values", omegaPowers)
	for i:=0; i<=(N/2)-1; i++ {
		secondTerm := omegaPowers[i] * odd[i]
		result[i] = even[i] + secondTerm
		result[i+N/2] = even[i] - secondTerm
	}
	//fmt.Println("even", even)
	//fmt.Println("odd",odd)
	//fmt.Println("result", result)

	if inverse {
		for i:=range result {
			result[i] /= complex(float64(N), 0.0)
		}
	}
	return result
}

func main() {
	arr := toComplex([]int{5, 3, 2, 1})
	fmt.Println("arr    ", arr)

	fft_arr := myfft(arr, false)
	fmt.Println("my fft ", fft_arr)
	ifft_arr := myfft(fft_arr, true)
	fmt.Println("my ifft", ifft_arr)

	lib_fft := fft.Fft(arr, false)
	fmt.Println("fft    ", lib_fft)
	fmt.Println("ifft   ", fft.Fft(lib_fft, true))
}

func toComplex(arr []int) []complex128 {
	result := make([]complex128, len(arr))
	for i,v := range arr {
		result[i] = complex(float64(v), 0.0)
	}
	return result
}
