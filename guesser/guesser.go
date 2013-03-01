package guesser

import (
	"image"
	"image/color"
)

const (
	Purple = iota
	Orange
	Blue
	White
	Black
)

type ScoreBar byte

func (s ScoreBar) String() string {
	switch s {
	case Purple:
		return "Purple"
	case Orange:
		return "Orange"
	case Blue:
		return "Blue"
	}
	return "Unknown"
}

func Reduce(img image.Image) image.Image {
	grid := reducedGrid(img)
	w := img.Bounds().Max.X - img.Bounds().Min.X
	h := img.Bounds().Max.Y - img.Bounds().Min.Y
	rimg := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			switch grid[x*h+y] {
			case Purple:
				rimg.Set(x, y, color.NRGBA{255, 0, 255, 255})
			case Orange:
				rimg.Set(x, y, color.NRGBA{255, 240, 0, 255})
			case Blue:
				rimg.Set(x, y, color.NRGBA{0, 0, 255, 255})
			case White:
				rimg.Set(x, y, color.White)
			case Black:
				rimg.Set(x, y, color.Black)
			}
		}
	}
	scorebarX0, scorebarX1, scorebarYs := scorebar(w, h, grid)
	for x := scorebarX0; x < scorebarX1; x++ {
		for y := scorebarYs[x][0]; y < scorebarYs[x][1]; y++ {
			rimg.Set(x, y, color.NRGBA{0, 255, 0, 255})
		}
	}
	return rimg
}

func Guess(img image.Image) ScoreBar {
	grid := reducedGrid(img)
	w := img.Bounds().Max.X - img.Bounds().Min.X
	h := img.Bounds().Max.Y - img.Bounds().Min.Y
	scorebarX0, scorebarX1, scorebarYs := scorebar(w, h, grid)
	counts := [5]int{0, 0, 0, 0, 0}
	for x := scorebarX0; x < scorebarX1; x++ {
		for y := scorebarYs[x][0]; y < scorebarYs[x][1]; y++ {
			counts[grid[x*h+y]]++
		}
	}
	switch {
	case counts[Purple] > counts[Orange] && counts[Purple] > counts[Blue]:
		return Purple
	case counts[Orange] > counts[Purple] && counts[Orange] > counts[Blue]:
		return Orange
	case counts[Blue] > counts[Purple] && counts[Blue] > counts[Orange]:
		return Blue
	}
	return Black
}

func scorebar(w, h int, grid []ScoreBar) (int, int, [][2]int) {
	stripes := make([][2]int, w/2)
	stripes[0][0], stripes[0][1] = stripe(0, w, h, grid)
	stripeWidth := 0
	stripeHeight := 0
	maxStripeX := 0
	maxStripeHeight := 0
	maxStripeWidth := 0
	for x := 1; x < w/2; x++ {
		stripes[x][0], stripes[x][1] = stripe(x, w, h, grid)
		if stripes[x][0] < stripes[x-1][1] && stripes[x][1] > stripes[x-1][0] {
			stripeWidth++
			if stripeHeight < stripes[x][1]-stripes[x][0] {
				stripeHeight = stripes[x][1] - stripes[x][0]
			}
		} else if stripeWidth > 0 {
			if stripeWidth > maxStripeWidth && stripeHeight > maxStripeHeight {
				maxStripeX = x - stripeWidth
				maxStripeHeight = stripeHeight
				maxStripeWidth = stripeWidth
			}
			stripeWidth = 0
			stripeHeight = 0
		}
	}
	return maxStripeX, maxStripeX + maxStripeWidth, stripes
}

func stripe(x, w, h int, grid []ScoreBar) (int, int) {
	blacks := 0
	y0 := 0
	y1 := 0
	y0result := 0
	y1result := 0
	for y := 0; y < h; y++ {
		if grid[x*h+y] != Black {
			if blacks != 0 {
				if blacks > 5 {
					y0 = y
					y1 = y
				} else {
					y0 = 0
					y1 = 0
				}
			} else if y0 > 0 {
				y1 = y
			}
			blacks = 0
		} else {
			blacks++
			if y1 > y0 && blacks > 5 {
				if y1-y0 > y1result-y0result {
					y1result = y1
					y0result = y0
					y1 = 0
					y0 = 0
				}
			}
		}
	}
	return y0result, y1result
}

func reducedGrid(img image.Image) []ScoreBar {
	dark, bright := brightnessRange(img)
	x0 := img.Bounds().Min.X
	x1 := img.Bounds().Max.X
	y0 := img.Bounds().Min.Y
	y1 := img.Bounds().Max.Y
	w := x1 - x0
	h := y1 - y0
	grid := make([]ScoreBar, w*h)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			grid[x*h+y] = classify(img.At(x+x0, y+y0), dark, bright)
		}
	}

	// blur
	blurred := make([]ScoreBar, w*h)
	for x := 1; x < w-1; x++ {
		for y := 3; y < h-3; y++ {
			block := [5]int{0, 0, 0, 0, 0}
			for i := -1; i < 2; i++ {
				for j := -3; j < 4; j++ {
					block[grid[(x+i)*h+y+j]]++
				}
			}
			switch {
			case block[Purple] > 10 && block[Purple] > block[Orange] && block[Purple] > block[Blue]:
				blurred[x*h+y] = Purple
			case block[Orange] > 10 && block[Orange] > block[Purple] && block[Orange] > block[Blue]:
				blurred[x*h+y] = Orange
			case block[Blue] > 10 && block[Blue] > block[Purple] && block[Blue] > block[Orange]:
				blurred[x*h+y] = Blue
			case block[Black] < 9:
				blurred[x*h+y] = White
			default:
				blurred[x*h+y] = Black
			}
		}
	}

	return blurred
}

func brightnessRange(img image.Image) (uint32, uint32) {
	x0 := img.Bounds().Min.X
	x1 := img.Bounds().Max.X
	y0 := img.Bounds().Min.Y
	y1 := img.Bounds().Max.Y
	r, g, b, _ := img.At(x0, y0).RGBA()
	min := r + g + b
	max := r + g + b
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			r, g, b, _ = img.At(x, y).RGBA()
			if min > r+g+b {
				min = r + g + b
			}
			if max < r+g+b {
				max = r + g + b
			}
		}
	}
	min1 := min
	min2 := (min + max) / 2
	min3 := max
	max1 := min
	max2 := (min + max) / 2
	max3 := max
	for i := 0; i < 10; i++ {
		countMin := 0
		countMax := 0
		for x := x0; x < x1; x++ {
			for y := y0; y < y1; y++ {
				r, g, b, _ = img.At(x, y).RGBA()
				if r+g+b > min2 {
					countMin++
				}
				if r+g+b > max2 {
					countMax++
				}
			}
		}
		if countMin > (x1-x0)*(y1-y0)/2 {
			min3 = min2
			min2 = (min1 + min2) / 2
		} else {
			min1 = min2
			min2 = (min2 + min3) / 2
		}
		if countMax > (x1-x0)*(y1-y0)/6 {
			max1 = max2
			max2 = (max2 + max3) / 2
		} else {
			max3 = max2
			max2 = (max1 + max2) / 2
		}
	}
	return min2, max2
}

func classify(c color.Color, dark, bright uint32) ScoreBar {
	r, g, b, _ := c.RGBA()
	switch {
	case r+g+b < dark/3:
		return Black
	case r > g && b > g:
		return Purple
	case r > b && g > b:
		return Orange
	case b > r && b > g:
		return Blue
	}
	return White
}
