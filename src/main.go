package main

import (
	"embed"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

var (
	//go:embed media/*
	mediaList embed.FS
)

const (
	upCode     = 1
	rightCode  = 2
	downCode   = 3
	leftCode   = 4
	xCode      = 5
	aCode      = 6
	bCode      = 7
	yCode      = 8
	l1Code     = 9
	l2Code     = 10
	r1Code     = 11
	r2Code     = 12
	selectCode = 13
	menuCode   = 14
	startCode  = 15
)

const (
	screenLeftUpX      = int32(361)
	screenLeftUpY      = int32(192)
	screenRightDownX   = int32(919)
	screenRightDownY   = int32(508)
	screenWidth        = screenRightDownX - screenLeftUpX
	screenHeight       = screenRightDownY - screenLeftUpY
	joystickInitialMax = 0.9
)

func main() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(1280, 720, "TrimUI Hardware Test")
	rl.InitAudioDevice()
	rl.SetTargetFPS(60)

	var (
		gamePadId                                          int32   = 0
		shouldExit                                                 = false
		crossBytes                                                 = orPanicRes(mediaList.ReadFile("media/cross_yellow.png"))
		crossTexture                                               = rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", crossBytes, int32(len(crossBytes))))
		bgrBytes                                                   = orPanicRes(mediaList.ReadFile("media/bgr.png"))
		bgrTexture                                                 = rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", bgrBytes, int32(len(bgrBytes))))
		soundBytes                                                 = orPanicRes(mediaList.ReadFile("media/sound.wav"))
		sound                                                      = rl.LoadSoundFromWave(rl.LoadWaveFromMemory(".wav", soundBytes, int32(len(soundBytes))))
		keyColor                                                   = rl.Lime
		x1, y1, x2, y2                                     float64 = 0, 0, 0, 0
		roundedX1, roundedY1, roundedX2, roundedY2         float64 = 0, 0, 0, 0
		correctedX1, correctedY1, correctedX2, correctedY2 float64 = 0, 0, 0, 0
		maxX1, maxY1, maxX2, maxY2                                 = joystickInitialMax, joystickInitialMax, joystickInitialMax, joystickInitialMax
		minX1, minY1, minX2, minY2                                 = -joystickInitialMax, -joystickInitialMax, -joystickInitialMax, -joystickInitialMax
	)

	for !rl.WindowShouldClose() && !shouldExit {

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawTexture(bgrTexture, 0, 0, rl.White)

		rl.DrawText("MENU+SELECT : Play Sound", 5, 695, 20, rl.DarkBlue)
		rl.DrawText("MENU+START  : Exit", 1060, 695, 20, rl.DarkBlue)

		//functional
		if rl.IsGamepadButtonDown(gamePadId, menuCode) || rl.IsKeyDown(rl.KeyM) {
			rl.DrawCircle(331, 484, 11, keyColor)
		}
		if rl.IsGamepadButtonDown(gamePadId, selectCode) || rl.IsKeyDown(rl.KeyS) {
			rl.DrawCircle(947, 484, 11, keyColor)
		}
		if rl.IsGamepadButtonDown(gamePadId, startCode) || rl.IsKeyDown(rl.KeyT) {
			rl.DrawCircle(991, 484, 11, keyColor)
		}
		//letters
		if rl.IsGamepadButtonDown(gamePadId, aCode) || rl.IsKeyDown(rl.KeyA) {
			rl.DrawCircle(1023, 292, 13, keyColor)
		}
		if rl.IsGamepadButtonDown(gamePadId, bCode) || rl.IsKeyDown(rl.KeyB) {
			rl.DrawCircle(988, 327, 13, keyColor)
		}
		if rl.IsGamepadButtonDown(gamePadId, xCode) || rl.IsKeyDown(rl.KeyX) {
			rl.DrawCircle(988, 258, 13, keyColor)
		}
		if rl.IsGamepadButtonDown(gamePadId, yCode) || rl.IsKeyDown(rl.KeyY) {
			rl.DrawCircle(954, 292, 13, keyColor)
		}
		//arrows
		if rl.IsGamepadButtonDown(gamePadId, leftCode) || rl.IsKeyDown(rl.KeyLeft) {
			rl.DrawRectangle(249, 280, 25, 24, keyColor)
		}
		if rl.IsGamepadButtonDown(gamePadId, rightCode) || rl.IsKeyDown(rl.KeyRight) {
			rl.DrawRectangle(303, 280, 25, 24, keyColor)
		}
		if rl.IsGamepadButtonDown(gamePadId, upCode) || rl.IsKeyDown(rl.KeyUp) {
			rl.DrawRectangle(277, 253, 24, 24, keyColor)
		}
		if rl.IsGamepadButtonDown(gamePadId, downCode) || rl.IsKeyDown(rl.KeyDown) {
			rl.DrawRectangle(277, 308, 24, 24, keyColor)
		}
		//shoulders
		if rl.IsGamepadButtonDown(gamePadId, l1Code) || rl.IsKeyDown(rl.KeyOne) {
			rl.DrawRectangleRounded(rl.Rectangle{X: 218, Y: 82, Width: 110, Height: 20}, 0.9, 10, keyColor)
		}
		if rl.IsGamepadButtonDown(gamePadId, r1Code) || rl.IsKeyDown(rl.KeyTwo) {
			rl.DrawRectangleRounded(rl.Rectangle{X: 955, Y: 82, Width: 110, Height: 20}, 0.9, 10, keyColor)
		}
		if rl.IsGamepadButtonDown(gamePadId, l2Code) || rl.IsKeyDown(rl.KeyThree) {
			rl.DrawRectangleRounded(rl.Rectangle{X: 235, Y: 45, Width: 90, Height: 20}, 0.9, 10, keyColor)
		}
		if rl.IsGamepadButtonDown(gamePadId, r2Code) || rl.IsKeyDown(rl.KeyFour) {
			rl.DrawRectangleRounded(rl.Rectangle{X: 955, Y: 45, Width: 90, Height: 20}, 0.9, 10, keyColor)
		}
		//joysticks
		x1 = float64(rl.GetGamepadAxisMovement(gamePadId, rl.GamepadAxisLeftX))
		y1 = float64(rl.GetGamepadAxisMovement(gamePadId, rl.GamepadAxisLeftY))
		x2 = float64(rl.GetGamepadAxisMovement(gamePadId, rl.GamepadAxisRightX))
		y2 = float64(rl.GetGamepadAxisMovement(gamePadId, rl.GamepadAxisRightY))

		maxX1 = math.Max(x1, maxX1)
		maxY1 = math.Max(x1, maxY1)
		maxX2 = math.Max(x2, maxX2)
		maxY2 = math.Max(y2, maxY2)

		minX1 = math.Min(x1, minX1)
		minY1 = math.Min(x1, minY1)
		minX2 = math.Min(x2, minX2)
		minY2 = math.Min(y2, minY2)

		roundedX1 = toFixed(x1, 3)
		roundedY1 = toFixed(y1, 3)
		roundedX2 = toFixed(x2, 3)
		roundedY2 = toFixed(y2, 3)

		correctedX1 = x1 / If(x1 > 0, maxX1, -minX1)
		correctedY1 = y1 / If(y1 > 0, maxY1, -minY1)
		correctedX2 = x2 / If(x2 > 0, maxX2, -minX2)
		correctedY2 = y2 / If(y2 > 0, maxY2, -minY2)

		if (roundedX1 != 0 && roundedY1 != 0) || rl.IsKeyDown(rl.KeyL) {
			rl.DrawCircle(304, 412, 20, keyColor)
		}

		if (roundedX2 != 0 && roundedY2 != 0) || rl.IsKeyDown(rl.KeyR) {
			rl.DrawCircle(977, 412, 20, keyColor)
		}

		if roundedX1 != 0 && roundedY1 != 0 {
			rl.DrawTexture(
				crossTexture,
				screenLeftUpX+(screenWidth-crossTexture.Width)/2+int32(float64((screenWidth-crossTexture.Width)/2)*correctedX1),
				screenLeftUpY+(screenHeight-crossTexture.Height)/2+int32(float64((screenHeight-crossTexture.Height)/2)*correctedY1),
				rl.White)
		}
		if roundedX2 != 0 && roundedY2 != 0 {
			rl.DrawTexture(
				crossTexture,
				screenLeftUpX+(screenWidth-crossTexture.Width)/2+int32(float64((screenWidth-crossTexture.Width)/2)*correctedX2),
				screenLeftUpY+(screenHeight-crossTexture.Height)/2+int32(float64((screenHeight-crossTexture.Height)/2)*correctedY2),
				rl.White)
		}
		if (rl.IsGamepadButtonDown(gamePadId, menuCode) && rl.IsGamepadButtonDown(gamePadId, selectCode)) || rl.IsKeyDown(rl.KeyU) {
			rl.PlaySound(sound)
		}

		//exit
		if rl.IsGamepadButtonDown(gamePadId, menuCode) && rl.IsGamepadButtonDown(gamePadId, startCode) {
			shouldExit = true //see WindowShouldClose, it checks if KeyEscape pressed or Close icon pressed
		}
		rl.EndDrawing()
	}
	rl.UnloadTexture(bgrTexture)
	rl.UnloadTexture(crossTexture)
	rl.CloseWindow()
}

func orPanic(err interface{}) {
	switch v := err.(type) {
	case error:
		if v != nil {
			panic(err)
		}
	case bool:
		if !v {
			panic("condition failed: != true")
		}
	}
}

func orPanicRes[T any](res T, err interface{}) T {
	orPanic(err)
	return res
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func If[T any](cond bool, vTrue, vFalse T) T {
	if cond {
		return vTrue
	}
	return vFalse
}
