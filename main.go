package main

import (
	"embed"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/gonutz/di8"
	"github.com/gonutz/prototype/draw"
	"github.com/gonutz/w32"
)

//go:embed rsc/*
var rsc embed.FS

const (
	windowH            = 500
	kiwiW, kiwiH       = 343, 300
	ballW, ballH       = 60, 60
	leftKiwiPath       = "rsc/blue.png"
	leftKiwiShootPath  = "rsc/blue_shoot.png"
	rightKiwiPath      = "rsc/white.png"
	rightKiwiShootPath = "rsc/white_shoot.png"
	ballPath           = "rsc/ball.png"
	kiwiSpeed          = 15
	shootFrames        = 6
	shootCooldown      = shootFrames + 8
	minBallShootSpeed  = 30
	maxBallShootSpeed  = 50
	ballFriction       = 3
	leftGoalSoundPath  = "rsc/white_goal.wav"
	rightGoalSoundPath = "rsc/blue_goal.wav"
	leftWinSoundPath   = "rsc/blue_win.wav"
	rightWinSoundPath  = "rsc/white_win.wav"
	winScore           = 10
	winSoundCooldown   = 40
	backMusicPath      = "rsc/fuss_song.wav"
	blinkCooldown      = 30
)

var (
	leftKiwiShootX      = [2]int{86, 143}
	rightKiwiShootX     = [2]int{160, 218}
	ballHitBoxX         = [2]int{7, 52}
	leftShootSoundPaths = []string{
		"rsc/blue_shoot1.wav",
		"rsc/blue_shoot2.wav",
		"rsc/blue_shoot3.wav",
	}
	rightShootSoundPaths = []string{
		"rsc/white_shoot1.wav",
		"rsc/white_shoot2.wav",
		"rsc/white_shoot3.wav",
	}
	ballShootSoundPaths = []string{
		"rsc/ball_shoot1.wav",
		"rsc/ball_shoot2.wav",
	}
)

func main() {
	draw.OpenFile = func(path string) (io.ReadCloser, error) {
		return rsc.Open(path)
	}

	rand.Seed(time.Now().UnixNano())
	dinputInited := false
	r := w32.GetWindowRect(w32.GetDesktopWindow())
	windowW := int(r.Right - r.Left - 30)

	var left, right player
	ballX := (windowW - ballW) / 2
	ballRotation := 0
	ballVx := 0
	right.x = windowW - kiwiW
	scoringTimer := 0
	var leftWon, rightWon bool
	winSoundTimer := 0
	musicTimer := 0
	winShowRestartTimer := 0
	winRestartBlinkTimer := 0
	restartBlinking := false

	restart := func() {
		left.score = 0
		left.shootCooldown = 0
		left.shootFrames = 0
		left.x = 0
		right.score = 0
		right.shootCooldown = 0
		right.shootFrames = 0
		right.x = windowW - kiwiW
		scoringTimer = 0
		leftWon, rightWon = false, false
		winSoundTimer = 0
		winShowRestartTimer = 0
		winRestartBlinkTimer = 0
		restartBlinking = false
		ballRotation = 0
	}

	check(draw.RunWindow("Jolinas Kiwi Fußball", windowW, windowH, func(window draw.Window) {
		if !dinputInited {
			setWindowIcon()
			initDInput()
			dinputInited = true
		}

		if window.WasKeyPressed(draw.KeyEscape) {
			window.Close()
		}

		if musicTimer == 0 {
			window.PlaySoundFile(backMusicPath)
			musicTimer = 241
		}
		musicTimer--

		// query the controls from keyboard and gamepad
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

		if scoringTimer > 0 {
			scoringTimer--
			if scoringTimer == 0 {
				left.shootFrames = 0
				winSoundTimer = 0
				left.shootCooldown = 0
				left.x = 0
				right.shootFrames = 0
				right.shootCooldown = 0
				right.x = windowW - kiwiW
				ballX = (windowW - ballW) / 2
				ballVx = 0
				ballRotation = 0
				if left.score >= winScore {
					leftWon = true
				}
				if right.score >= winScore {
					rightWon = true
				}
				if leftWon || rightWon {
					winSoundTimer = winSoundCooldown
					winShowRestartTimer = winSoundCooldown + 90
				}
			}
		} else if leftWon || rightWon {
			if winShowRestartTimer == 0 {
				// if the restart instruction is showing and either player kicks
				// restart the game
				if window.WasKeyPressed(draw.KeyEnter) ||
					window.WasKeyPressed(draw.KeySpace) ||
					leftShootDown || rightShootDown {
					restart()
				}
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
			if leftShootDown && left.shootCooldown == 0 {
				// start shooting
				left.shootFrames = shootFrames
				left.shootCooldown = shootCooldown
				window.PlaySoundFile(leftShootSoundPaths[rand.Intn(len(leftShootSoundPaths))])
				// check ball collision
				shootLeft := left.x + leftKiwiShootX[0]
				shootRight := left.x + leftKiwiShootX[1]
				d := abs((ballLeft+ballRight)/2 - (shootLeft+shootRight)/2)
				if d < (ballRight-ballLeft)/2+(shootRight-shootLeft)/2 {
					ballVx += minBallShootSpeed + rand.Intn(maxBallShootSpeed-minBallShootSpeed)
					window.PlaySoundFile(ballShootSoundPaths[rand.Intn(len(ballShootSoundPaths))])
				}
			}
			if rightShootDown && right.shootCooldown == 0 {
				// start shooting
				right.shootFrames = shootFrames
				right.shootCooldown = shootCooldown
				window.PlaySoundFile(rightShootSoundPaths[rand.Intn(len(rightShootSoundPaths))])
				// check ball collision
				shootLeft := right.x + rightKiwiShootX[0]
				shootRight := right.x + rightKiwiShootX[1]
				d := abs((ballLeft+ballRight)/2 - (shootLeft+shootRight)/2)
				if d < (ballRight-ballLeft)/2+(shootRight-shootLeft)/2 {
					ballVx -= minBallShootSpeed + rand.Intn(maxBallShootSpeed-minBallShootSpeed)
					window.PlaySoundFile(ballShootSoundPaths[rand.Intn(len(ballShootSoundPaths))])
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
				if left.x < -kiwiW/2 {
					left.x = -kiwiW / 2
				}
				if left.x > windowW-kiwiW/4 {
					left.x = windowW - kiwiW/4
				}
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
				if right.x < -3*kiwiW/4 {
					right.x = -3 * kiwiW / 4
				}
				if right.x > windowW-kiwiW/2 {
					right.x = windowW - kiwiW/2
				}
			}
			// move ball
			ballX += ballVx
			ballRotation += ballVx
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
			leftGoal := ballX+ballHitBoxX[0] >= windowW
			rightGoal := ballX+ballHitBoxX[1] < 0
			if leftGoal || rightGoal {
				if leftGoal {
					left.score++
					window.PlaySoundFile(leftGoalSoundPath)
				} else {
					right.score++
					window.PlaySoundFile(rightGoalSoundPath)
				}
				scoringTimer = 60
			}
		}

		// draw everything
		window.FillRect(0, 0, windowW, windowH, draw.LightGreen)
		const scoreScale = 3
		score := fmt.Sprintf("%d : %d", left.score, right.score)
		scoreTextW, scoreTextH := window.GetScaledTextSize(score, scoreScale)
		window.DrawScaledText(score, (windowW-scoreTextW)/2, 10, scoreScale, draw.Black)
		if leftWon || rightWon {
			winSoundTimer--
			if winSoundTimer < 0 {
				winSoundTimer = 0
			}
			if winSoundTimer == 1 {
				if leftWon {
					window.PlaySoundFile(leftWinSoundPath)
				}
				if rightWon {
					window.PlaySoundFile(rightWinSoundPath)
				}
			}
			if leftWon || rightWon {
				winner := leftKiwiPath
				if rightWon {
					winner = rightKiwiPath
				}
				window.DrawImageFile(winner, (windowW-kiwiW)/2, scoreTextH+(windowH-scoreTextH-kiwiH)/2)
			}
			winShowRestartTimer--
			if winShowRestartTimer < 0 {
				winShowRestartTimer = 0
			}
			if winShowRestartTimer == 0 {
				winRestartBlinkTimer--
				if winRestartBlinkTimer < 0 {
					winRestartBlinkTimer = blinkCooldown
					restartBlinking = !restartBlinking
				}
				if restartBlinking {
					window.DrawScaledText(
						"Zum Neustart kicken/Enter/Leertaste",
						10,
						windowH-scoreTextH,
						2,
						draw.Black,
					)
				}
			}
		} else {
			// draw left kiwi
			leftPath := leftKiwiPath
			if left.shootFrames > 0 {
				leftPath = leftKiwiShootPath
			}
			window.DrawImageFile(leftPath, left.x, windowH-kiwiH-20)
			// draw ball
			window.DrawImageFileRotated(ballPath, ballX, windowH-ballH-10, ballRotation)
			// draw right kiwi
			rightPath := rightKiwiPath
			if right.shootFrames > 0 {
				rightPath = rightKiwiShootPath
			}
			window.DrawImageFile(rightPath, right.x, windowH-kiwiH)
		}
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

func setWindowIcon() {
	// the icon is contained in the .exe file as a resource, load it and set it
	// as the window icon so it appears in the top-left corner of the window and
	// when you alt+tab between windows
	const iconResourceID = 10
	iconHandle := w32.LoadImage(
		w32.GetModuleHandle(""),
		w32.MakeIntResource(iconResourceID),
		w32.IMAGE_ICON,
		0,
		0,
		w32.LR_DEFAULTSIZE|w32.LR_SHARED,
	)
	if iconHandle != 0 {
		window := w32.GetActiveWindow()
		w32.SendMessage(window, w32.WM_SETICON, w32.ICON_SMALL, uintptr(iconHandle))
		w32.SendMessage(window, w32.WM_SETICON, w32.ICON_SMALL2, uintptr(iconHandle))
		w32.SendMessage(window, w32.WM_SETICON, w32.ICON_BIG, uintptr(iconHandle))
	}
}
