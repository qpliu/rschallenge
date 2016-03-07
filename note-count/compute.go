package main

import (
	"errors"
)

func ComputeTotal(lowerBound, upperBound int, percents []int) (int, error) {
loop:
	for n := lowerBound; n < upperBound; n++ {
		for _, p := range percents {
			hits := (n*p + 5000) / 10000
			if hits*10000/n != p {
				continue loop
			}
		}
		return n, nil
	}
	return 0, errors.New("None found")
}
