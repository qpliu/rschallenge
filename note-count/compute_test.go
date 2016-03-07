package main

import (
	"testing"
)

func test(t *testing.T, lowerBound int, percents []int, expected int) {
	total, err := ComputeTotal(lowerBound, 5000, percents)
	if err != nil {
		t.Errorf("Error %v", err)
	} else if total != expected {
		t.Errorf("Expected %d, got %d", expected, total)
	}
}

func TestComputeTotal(t *testing.T) {
	test(t, 290, []int{9972, 9945, 9918, 9890, 9836, 9808, 9781, 9754, 9699, 9508, 9234, 9207}, 366)
	test(t, 335, []int{9980, 9885, 9733, 9180}, 525)
	test(t, 257, []int{9961, 9923, 9884, 9827, 9712, 9654}, 521)
	test(t, 610, []int{9983, 9950, 9934, 9786, 9754, 9622, 9065}, 610)
	test(t, 132, []int{9384, 9175, 9151, 9138, 8757}, 813)
}
