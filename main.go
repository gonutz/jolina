package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gonutz/di8"
	"github.com/gonutz/prototype/draw"
	"github.com/gonutz/w32"
)

const (
	windowW, windowH    = 1800, 500
	kiwiW, kiwiH        = 343, 300
	ballW, ballH        = 60, 60
	leftKiwiPath        = "rsc/blue.png"
	leftKiwiShootPath   = "rsc/blue_shoot.png"
	rightKiwiPath       = "rsc/white.png"
	rightKiwiShootPath  = "rsc/white_shoot.png"
	ballPath            = "rsc/ball.png"
	kiwiSpeed           = 15
	shootFrames         = 6
	shootCooldown       = shootFrames + 8
	minBallShootSpeed   = 30
	maxBallShootSpeed   = 50
	ballFriction        = 3
	leftShootSoundPath  = "rsc/blue_shoot.wav"
	rightShootSoundPath = "rsc/white_shoot.wav"
	ballShootSoundPath  = "rsc/ball_shoot.wav"
)

var (
	leftKiwiShootX  = [2]int{90, 140}
	rightKiwiShootX = [2]int{160, 220}
	ballHitBoxX     = [2]int{5, 50}
)

func main() {
	rand.Seed(time.Now().UnixNano())
	dinputInited := false
	var left, right player
	ballX := (windowW - ballW) / 2
	ballVx := 0
	right.x = windowW - kiwiW
	scoringTimer := 0
	check(draw.RunWindow("Jolina Kiwi Fußball", windowW, windowH, func(window draw.Window) {
		if !dinputInited {
			initDInput()
			dinputInited = true
		}

		if window.WasKeyPressed(draw.KeyEscape) {
			window.Close()
		}

		if scoringTimer > 0 {
			scoringTimer--
			if scoringTimer == 0 {
				left.shootFrames = 0
				left.shootCooldown = 0
				left.x = 0
				right.shootFrames = 0
				right.shootCooldown = 0
				right.x = windowW - kiwiW
				ballX = (windowW - ballW) / 2
				ballVx = 0
			}
		} else {
			// shoot
			left.shootFrames--
			if left.shootFrames < 0 {
				left.shootFrames = 0
			}
			left.shootCooldown--
			if left.shootCooldown < 0 {
				left.shootCooldown = 0
			}
			right.shootFrames--
			if right.shootFrames < 0 {
				right.shootFrames = 0
			}
			right.shootCooldown--
			if right.shootCooldown < 0 {
				right.shootCooldown = 0
			}
			ballLeft := ballX + ballHitBoxX[0]
			ballRight := ballX + ballHitBoxX[1]
			checkDevice := func(dev *di8.Device) (shoot, left, right bool) {
				var state di8.JOYSTATE
				if err := dev.GetDeviceState(&state); err == nil {
					pos := axisPos(uint32(state.X))
					if pos < -0.9 {
						left = true
					}
					if pos > 0.9 {
						right = true
					}
					n, err := dev.GetDeviceData(devBuf[:], 0)
					if err == nil {
						for i := range devBuf[:n] {
							if di8.JOFS_BUTTON0 <= devBuf[i].Ofs && devBuf[i].Ofs <= di8.JOFS_BUTTON31 &&
								devBuf[i].Data&0xFF != 0 {
								shoot = true
							}
						}
					} else if err != nil &&
						(err.Code() == di8.ERR_INPUTLOST || err.Code() == di8.ERR_NOTACQUIRED) {
						dev.Acquire()
					}
				} else if err != nil &&
					(err.Code() == di8.ERR_INPUTLOST || err.Code() == di8.ERR_NOTACQUIRED) {
					dev.Acquire()
				}
				return
			}
			leftShootDown := window.WasKeyPressed(draw.KeyW)
			leftLeftDown := window.IsKeyDown(draw.KeyA)
			leftRightDown := window.IsKeyDown(draw.KeyD)
			// see the controller input
			if len(devices) > 0 {
				shoot, left, right := checkDevice(devices[0])
				leftShootDown = leftShootDown || shoot
				leftLeftDown = leftLeftDown || left
				leftRightDown = leftRightDown || right
			}
			rightShootDown := window.WasKeyPressed(draw.KeyUp)
			rightLeftDown := window.IsKeyDown(draw.KeyLeft)
			rightRightDown := window.IsKeyDown(draw.KeyRight)
			if len(devices) > 1 {
				shoot, left, right := checkDevice(devices[1])
				rightShootDown = rightShootDown || shoot
				rightLeftDown = rightLeftDown || left
				rightRightDown = rightRightDown || right
			}
			if leftShootDown && left.shootCooldown == 0 {
				// start shooting
				left.shootFrames = shootFrames
				left.shootCooldown = shootCooldown
				window.PlaySoundFile(leftShootSoundPath)
				// check ball collision
				shootLeft := left.x + leftKiwiShootX[0]
				shootRight := left.x + leftKiwiShootX[1]
				d := abs((ballLeft+ballRight)/2 - (shootLeft+shootRight)/2)
				if d < (ballRight-ballLeft)/2+(shootRight-shootLeft)/2 {
					ballVx += minBallShootSpeed + rand.Intn(maxBallShootSpeed-minBallShootSpeed)
					window.PlaySoundFile(ballShootSoundPath)
				}
			}
			if rightShootDown && right.shootCooldown == 0 {
				// start shooting
				right.shootFrames = shootFrames
				right.shootCooldown = shootCooldown
				window.PlaySoundFile(rightShootSoundPath)
				// check ball collision
				shootLeft := right.x + rightKiwiShootX[0]
				shootRight := right.x + rightKiwiShootX[1]
				d := abs((ballLeft+ballRight)/2 - (shootLeft+shootRight)/2)
				if d < (ballRight-ballLeft)/2+(shootRight-shootLeft)/2 {
					ballVx -= minBallShootSpeed + rand.Intn(maxBallShootSpeed-minBallShootSpeed)
					window.PlaySoundFile(ballShootSoundPath)
				}
			}
			// move left player
			if left.shootCooldown == 0 {
				vx := 0
				if leftLeftDown {
					vx -= kiwiSpeed
				}
				if leftRightDown {
					vx += kiwiSpeed
				}
				left.x += vx
			}
			// move right player
			if right.shootCooldown == 0 {
				vx := 0
				if rightLeftDown {
					vx -= kiwiSpeed
				}
				if rightRightDown {
					vx += kiwiSpeed
				}
				right.x += vx
			}
			// move ball
			ballX += ballVx
			if ballVx > 0 {
				ballVx -= ballFriction
				if ballVx < 0 {
					ballVx = 0
				}
			} else if ballVx < 0 {
				ballVx += ballFriction
				if ballVx > 0 {
					ballVx = 0
				}
			}
			leftGoal := ballX+ballHitBoxX[0] > windowW
			rightGoal := ballX+ballHitBoxX[1] < 0
			if leftGoal || rightGoal {
				if leftGoal {
					left.score++
				} else {
					right.score++
				}
				scoringTimer = 60
			}
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
		window.DrawImageFile(ballPath, ballX, windowH-ballH)
		// draw right kiwi
		rightPath := rightKiwiPath
		if right.shootFrames > 0 {
			rightPath = rightKiwiShootPath
		}
		window.DrawImageFile(rightPath, right.x, windowH-kiwiH)
		const scoreScale = 3
		score := fmt.Sprintf("%d : %d", left.score, right.score)
		w, _ := window.GetScaledTextSize(score, scoreScale)
		window.DrawScaledText(score, (windowW-w)/2, 10, scoreScale, draw.Black)
	}))
	closeDInput()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type player struct {
	x             int
	shootFrames   int
	shootCooldown int
	score         int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var (
	dinput  *di8.DirectInput
	devices []*di8.Device
	devBuf  [32]di8.DEVICEOBJECTDATA
)

func initDInput() {
	var err error
	dinput, err = di8.Create(di8.HINSTANCE(uintptr(w32.GetModuleHandle(""))))
	if err != nil {
		return
	}
	var insts []di8.DEVICEINSTANCE
	dinput.EnumDevices(
		di8.DEVCLASS_GAMECTRL,
		func(inst *di8.DEVICEINSTANCE, ref uintptr) uintptr {
			insts = append(insts, *inst)
			return 1
		},
		0,
		di8.EDFL_ATTACHEDONLY,
	)
	window := di8.HWND(w32.GetActiveWindow())
	for i := range insts {
		if len(devices) == 2 {
			break
		}
		dev, err := dinput.CreateDevice(insts[i].GuidInstance)
		if err == nil {
			err = dev.SetCooperativeLevel(
				window,
				di8.SCL_EXCLUSIVE|di8.SCL_FOREGROUND,
			)
			if err != nil {
				continue
			}
			err = dev.SetDataFormat(&di8.Joystick)
			if err != nil {
				continue
			}
			err = dev.SetProperty(
				di8.PROP_BUFFERSIZE,
				di8.NewPropDWord(0, di8.PH_DEVICE, 32),
			)
			if err != nil {
				continue
			}
			err = dev.Acquire()
			if err != nil {
				continue
			}
			devices = append(devices, dev)
		}
	}
}

func closeDInput() {
	for i := range devices {
		devices[i].Release()
	}
	if dinput != nil {
		dinput.Release()
	}
}

func axisPos(data uint32) float64 {
	n := int(data) - 0xFFFF/2
	if n < 0 {
		return float64(n) / 32767
	} else {
		return float64(n) / 32768
	}
}
