# Kiwi Kick

![Screenshot](https://raw.githubusercontent.com/gonutz/jolina/master/screenshot.png)

Download the [installer here](https://github.com/gonutz/jolina/releases/download/v1.0.0/Jolina.Setup.exe).

This is a game I did in cooperation with my friend Jolina who is seven years old at the time of this game. She helped me out with the graphics and sound for this little soccer game.

The game runs on Windows only.

# Controls

Blue kiwi:

`A` `D`: move left/right

`W`: kick the ball

White kiwi:

`LEFT` `RIGHT`: move left/right

`UP`: kick the ball

There is also controller support, if you plug in one or two game controllers before running the game, they will be recognized automatically. Use the standard left/right axis for movement and press any key on the pad to kick the ball.

# Build

To build the project you need to have [the Go programming language](https://golang.org/dl/) installed. You also need [Git](https://git-scm.com/downloads). To build and run the program, type this in the command line:

```
go get github.com/gonutz/jolina
cd %GOPATH%\src\github.com\gonutz\jolina
build.bat
jolina.exe
```

There is also an [Inno Setup](http://www.jrsoftware.org/isinfo.php) script to build the Windows installer after you have built the game.