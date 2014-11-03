package bpm

import (
	"azul3d.org/v1/audio"
	"fmt"
)

func comb(b audio.Buffer, delay int) (avg, n float64) {
	length := b.Len()
	var (
		i    int
		done int
	)
	for i = length; (i - delay) > 0; i-- {
		delayed := i - delay
		if delayed < 0 {
			break
		}
		done++
		avg += (float64(b.At(delayed)) + 1.0) / 2.0
	}
	if done == 0 {
		//fmt.Println(avg, done, avg / float64(done))
		return 0, 0
	}
	return avg, float64(done)
}

func Chunk(b audio.Buffer, combSize, combDelay int) (bpm float64) {
	if combDelay >= combSize {
		panic("combDelay >= combSize")
	}

	length := b.Len()
	low := 1000.0
	high := -1000.0
	avg := 0.0
	for i := 0; i < length; i++ {
		end := i + combSize
		if end > length {
			end = length
		}
		combed, _ := comb(b.Slice(i, end), combDelay)
		if combed < low {
			low = combed
		}
		if combed > high {
			high = combed
		}
		avg += combed
	}
	avg = avg / float64(length)
	diff := high - avg
	fmt.Println("low", low, "high", high, "diff", diff, "avg", avg)

	for i := 0; i < length; i++ {
		end := i + combSize
		if end > length {
			end = length
		}
		combed, _ := comb(b.Slice(i, end), combDelay)
		if combed > (high - diff*0.007) {
			bpm++
		}
	}
	fmt.Println("bpm", bpm)
	return
}
