package main

import "github.com/gonutz/prototype/draw"

const (
	windowW, windowH   = 1800, 500
	kiwiW, kiwiH       = 343, 300
	ballW, ballH       = 60, 60
	leftKiwiPath       = "rsc/blue.png"
	leftKiwiShootPath  = "rsc/blue_shoot.png"
	rightKiwiPath      = "rsc/white.png"
	rightKiwiShootPath = "rsc/white_shoot.png"
	ballPath           = "rsc/ball.png"
	kiwiSpeed          = 15
	shootFrames        = 6
)

func main() {
	var left, right player
	right.x = windowW - kiwiW
	check(draw.RunWindow("Jolina Kiwi Fußball", windowW, windowH, func(window draw.Window) {
		if window.WasKeyPressed(draw.KeyEscape) {
			window.Close()
		}

		if window.WasKeyPressed(draw.KeyW) && left.shootFrames == 0 {
			left.shootFrames = shootFrames
		}
		if window.WasKeyPressed(draw.KeyUp) && right.shootFrames == 0 {
			right.shootFrames = shootFrames
		}
		// move left player
		{
			vx := 0
			if window.IsKeyDown(draw.KeyA) {
				vx -= kiwiSpeed
			}
			if window.IsKeyDown(draw.KeyD) {
				vx += kiwiSpeed
			}
			left.x += vx
		}
		// move right player
		{
			vx := 0
			if window.IsKeyDown(draw.KeyLeft) {
				vx -= kiwiSpeed
			}
			if window.IsKeyDown(draw.KeyRight) {
				vx += kiwiSpeed
			}
			right.x += vx
		}

		left.shootFrames--
		if left.shootFrames < 0 {
			left.shootFrames = 0
		}
		right.shootFrames--
		if right.shootFrames < 0 {
			right.shootFrames = 0
		}

		// draw everything
		window.FillRect(0, 0, windowW, windowH, draw.LightGreen)
		// draw left kiwi
		leftPath := leftKiwiPath
		if left.shootFrames > 0 {
			leftPath = leftKiwiShootPath
		}
		window.DrawImageFile(leftPath, left.x, windowH-kiwiH)
		// draw ball
		window.DrawImageFile(ballPath, (windowW-ballW)/2, windowH-ballH)
		// draw right kiwi
		rightPath := rightKiwiPath
		if right.shootFrames > 0 {
			rightPath = rightKiwiShootPath
		}
		window.DrawImageFile(rightPath, right.x, windowH-kiwiH)
	}))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type player struct {
	x           int
	shootFrames int
}
