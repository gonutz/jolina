//+build ignore

package main

import (
	"image"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	bw, err := readImage("ball.png")
	check(err)
	gray := toGray(bw)
	hist := grayHistogram(gray)
	threshold(gray, thresholdOtsu(&hist))
	writeImage(gray, "threshold.png")
	mask := makeCircleMask(3)
	opened := openBinary(gray, mask)
	writeImage(opened, "opened.png")
	closed := closeBinary(gray, mask)
	writeImage(closed, "closed.png")

	writeImage(closeBinary(openBinary(closeBinary(gray, mask), mask), mask), "final.png")
}

func makeCircleMask(diameter int) mask {
	img := image.NewGray(image.Rect(0, 0, diameter, diameter))
	outline := ellipseOutline(0, 0, diameter, diameter)
	for _, p := range outline {
		img.Pix[img.PixOffset(p.X, p.Y)] = 255
	}
	center := image.Pt(diameter/2, diameter/2)
	floodFillWhite(img, center)
	return mask{img: img, center: center}
}

func writeImage(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func readImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func threshold(img *image.Gray, at uint8) {
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			i := img.PixOffset(x, y)
			if img.Pix[i] < at {
				img.Pix[i] = 0
			} else {
				img.Pix[i] = 255
			}
		}
	}
}

func thresholdOtsu(hist *histogram) uint8 {
	var sum float64
	for i := 0; i < 256; i++ {
		sum += float64(i * hist.counts[i])
	}

	var sumB float64
	var wb, wf int

	var varMax float64
	threshold := 0

	for t := 0; t < 256; t++ {
		wb += hist.counts[t]
		if wb == 0 {
			continue
		}
		wf = hist.total - wb
		if wf == 0 {
			break
		}

		sumB += float64(t * hist.counts[t])

		mb := sumB / float64(wb)
		mf := (sum - sumB) / float64(wf)

		varBetween := float64(wb) * float64(wf) * (mb - mf) * (mb - mf)

		if varBetween > varMax {
			varMax = varBetween
			threshold = t
		}
	}
	return uint8(threshold)
}

func grayHistogram(img *image.Gray) histogram {
	var hist histogram
	b := img.Bounds()
	hist.total = b.Dx() * b.Dy()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			hist.counts[img.Pix[img.PixOffset(x, y)]]++
		}
	}
	return hist
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func toGray(img image.Image) *image.Gray {
	gray := image.NewGray(img.Bounds())
	if nrgba, ok := img.(*image.NRGBA); ok {
		done := make(chan bool)
		do := func(b image.Rectangle) {
			for x := b.Min.X; x < b.Max.X; x++ {
				for y := b.Min.Y; y < b.Max.Y; y++ {
					i := nrgba.PixOffset(x, y)
					r := uint32(nrgba.Pix[i])
					r |= r << 8
					g := uint32(nrgba.Pix[i+1])
					g |= g << 8
					b := uint32(nrgba.Pix[i+2])
					b |= b << 8
					gray.Pix[gray.PixOffset(x, y)] = uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
				}
			}
			done <- true
		}
		b := gray.Bounds()
		go do(image.Rect(b.Min.X, b.Min.Y, b.Min.X+b.Dx()/2, b.Min.Y+b.Dy()/2))
		go do(image.Rect(b.Min.X+b.Dx()/2, b.Min.Y, b.Max.X, b.Min.Y+b.Dy()/2))
		go do(image.Rect(b.Min.X, b.Min.Y+b.Dy()/2, b.Min.X+b.Dx()/2, b.Max.Y))
		go do(image.Rect(b.Min.X+b.Dx()/2, b.Min.Y+b.Dy()/2, b.Max.X, b.Max.Y))
		<-done
		<-done
		<-done
		<-done
		return gray
	}
	draw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, draw.Src)
	return gray
}

type histogram struct {
	counts [256]int
	total  int
}

func openBinary(img *image.Gray, mask mask) *image.Gray {
	return dilateBinary(erodeBinary(img, mask), mask)
}

func openBinaryInverse(img *image.Gray, mask mask) *image.Gray {
	return inverted(dilateBinary(erodeBinary(inverted(img), mask), mask))
}

func closeBinary(img *image.Gray, mask mask) *image.Gray {
	return erodeBinary(dilateBinary(img, mask), mask)
}

func closeBinaryInverse(img *image.Gray, mask mask) *image.Gray {
	return inverted(erodeBinary(dilateBinary(inverted(img), mask), mask))
}

type mask struct {
	img    *image.Gray
	center image.Point
}

func dilateBinary(img *image.Gray, mask mask) *image.Gray {
	imgBounds := img.Bounds()
	maskBounds := mask.img.Bounds()
	addLeft := mask.center.X - maskBounds.Min.X
	addTop := mask.center.Y - maskBounds.Min.Y
	addRight := maskBounds.Max.X - mask.center.X - 1
	addBottom := maskBounds.Max.Y - mask.center.Y - 1

	outBounds := imgBounds
	outBounds.Min.X -= addLeft
	outBounds.Min.Y -= addTop
	outBounds.Max.X += addRight
	outBounds.Max.Y += addBottom
	out := image.NewGray(outBounds)

	done := make(chan bool)
	do := func(imgBounds image.Rectangle) {
		for x := imgBounds.Min.X; x < imgBounds.Max.X; x++ {
			for y := imgBounds.Min.Y; y < imgBounds.Max.Y; y++ {
				if img.Pix[img.PixOffset(x, y)] != 0 {
					// pixel is set -> apply mask here
					for mx := maskBounds.Min.X; mx < maskBounds.Max.X; mx++ {
						for my := maskBounds.Min.Y; my < maskBounds.Max.Y; my++ {
							if mask.img.Pix[mask.img.PixOffset(mx, my)] != 0 {
								// mask pixel is set as well -> set pixel in result
								ox := x + mx - mask.center.X
								oy := y + my - mask.center.Y
								out.Pix[out.PixOffset(ox, oy)] = 255
							}
						}
					}
				}
			}
		}
		done <- true
	}
	b := imgBounds
	go do(image.Rect(b.Min.X, b.Min.Y, b.Min.X+b.Dx()/2, b.Min.Y+b.Dy()/2))
	go do(image.Rect(b.Min.X+b.Dx()/2, b.Min.Y, b.Max.X, b.Min.Y+b.Dy()/2))
	go do(image.Rect(b.Min.X, b.Min.Y+b.Dy()/2, b.Min.X+b.Dx()/2, b.Max.Y))
	go do(image.Rect(b.Min.X+b.Dx()/2, b.Min.Y+b.Dy()/2, b.Max.X, b.Max.Y))
	<-done
	<-done
	<-done
	<-done

	return out.SubImage(img.Bounds()).(*image.Gray)
}

func erodeBinary(img *image.Gray, mask mask) *image.Gray {
	imgBounds := img.Bounds()
	maskBounds := mask.img.Bounds()
	addLeft := mask.center.X - maskBounds.Min.X
	addTop := mask.center.Y - maskBounds.Min.Y
	addRight := maskBounds.Max.X - mask.center.X - 1
	addBottom := maskBounds.Max.Y - mask.center.Y - 1

	outBounds := imgBounds
	outBounds.Min.X += addLeft
	outBounds.Min.Y += addTop
	outBounds.Max.X -= addRight
	outBounds.Max.Y -= addBottom
	out := image.NewGray(outBounds)

	done := make(chan bool)
	do := func(outBounds image.Rectangle) {
		for x := outBounds.Min.X; x < outBounds.Max.X; x++ {
			for y := outBounds.Min.Y; y < outBounds.Max.Y; y++ {
				for mx := maskBounds.Min.X; mx < maskBounds.Max.X; mx++ {
					for my := maskBounds.Min.Y; my < maskBounds.Max.Y; my++ {
						if mask.img.Pix[mask.img.PixOffset(mx, my)] != 0 &&
							img.Pix[img.PixOffset(x+mx-mask.center.X, y+my-mask.center.Y)] == 0 {
							goto skip
						}
					}
				}
				out.Pix[out.PixOffset(x, y)] = 255
			skip:
			}
		}
		done <- true
	}

	b := outBounds
	go do(image.Rect(b.Min.X, b.Min.Y, b.Min.X+b.Dx()/2, b.Min.Y+b.Dy()/2))
	go do(image.Rect(b.Min.X+b.Dx()/2, b.Min.Y, b.Max.X, b.Min.Y+b.Dy()/2))
	go do(image.Rect(b.Min.X, b.Min.Y+b.Dy()/2, b.Min.X+b.Dx()/2, b.Max.Y))
	go do(image.Rect(b.Min.X+b.Dx()/2, b.Min.Y+b.Dy()/2, b.Max.X, b.Max.Y))
	<-done
	<-done
	<-done
	<-done

	return out
}

func inverted(img *image.Gray) *image.Gray {
	i := copyGray(img)
	invert(i)
	return i
}

func invert(img *image.Gray) {
	for i := range img.Pix {
		img.Pix[i] = 255 - img.Pix[i]
	}
}

func copyGray(img *image.Gray) *image.Gray {
	cpy := image.NewGray(img.Bounds())
	draw.Draw(cpy, cpy.Bounds(), img, img.Bounds().Min, draw.Src)
	return cpy
}

// ellipseOutline returns a list of pixel positions that mark the outline of the
// requested ellipse.
//
//       jih
//     lk   gf
//    m       e
//     no   cd
//       pab
//
func ellipseOutline(x, y, w, h int) []image.Point {
	quarter := quaterEllipsePoints(w, h)
	xPivot, yPivot := 0, 0
	if w%2 == 0 {
		xPivot = 1
	}
	if h%2 == 0 {
		yPivot = 1
	}
	dx, dy := x+w/2, y+h/2
	p := make([]image.Point, 0, len(quarter)*4)
	for i := range quarter {
		p = append(p, image.Pt(
			quarter[i].X+dx,
			quarter[i].Y+dy,
		))
	}
	for i := len(quarter) - 1 - (1 - yPivot); i >= 0; i-- {
		p = append(p, image.Pt(
			quarter[i].X+dx,
			-quarter[i].Y-yPivot+dy,
		))
	}
	for i := 1 - xPivot; i < len(quarter); i++ {
		p = append(p, image.Pt(
			-quarter[i].X-xPivot+dx,
			-quarter[i].Y-yPivot+dy,
		))
	}
	for i := len(quarter) - 1 - (1 - yPivot); i >= 1-xPivot; i-- {
		p = append(p, image.Pt(
			-quarter[i].X-xPivot+dx,
			quarter[i].Y+dy,
		))
	}
	return p
}

func quaterEllipsePoints(w, h int) (p []image.Point) {
	if w <= 0 || h <= 0 {
		return nil
	}

	a, b := (w-1)/2, (h-1)/2
	x, y := 0, b
	a2, b2 := a*a, b*b

	crit1 := -(a2/4 + a%2 + b2)
	crit2 := -(b2/4 + b%2 + a2)
	crit3 := -(b2/4 + b%2)
	t := -a2 * y
	dxt := 2 * b2 * x
	dyt := -2 * a2 * y
	d2xt := 2 * b2
	d2yt := 2 * a2

	for y >= 0 && x <= a {
		p = append(p, image.Pt(x, y))
		if t+b2*x <= crit1 || t+a2*y <= crit3 {
			x++
			dxt += d2xt
			t += dxt
		} else if t-a2*y > crit2 {
			y--
			dyt += d2yt
			t += dyt
		} else {
			x++
			dxt += d2xt
			t += dxt
			y--
			dyt += d2yt
			t += dyt
		}
	}
	return
}

func floodFillWhite(img *image.Gray, start image.Point) {
	b := img.Bounds()
	img.Pix[img.PixOffset(start.X, start.Y)] = 255
	q := []image.Point{start}
	var neighbors [4]image.Point
	for len(q) > 0 {
		pt := q[0]
		q = q[1:]
		// add black neighbors to the queue
		neighbors[0].X = pt.X + 1
		neighbors[0].Y = pt.Y
		neighbors[1].X = pt.X - 1
		neighbors[1].Y = pt.Y
		neighbors[2].X = pt.X
		neighbors[2].Y = pt.Y + 1
		neighbors[3].X = pt.X
		neighbors[3].Y = pt.Y - 1
		for _, n := range neighbors {
			i := img.PixOffset(n.X, n.Y)
			if n.In(b) && img.Pix[i] == 0 {
				img.Pix[i] = 255
				q = append(q, n)
			}
		}
	}
}
