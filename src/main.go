package main

import (
	"embed"
	mapset "github.com/deckarep/golang-set/v2"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	//go:embed media/*
	mediaList embed.FS
)

func main() {
	rl.InitWindow(1280, 720, "TrimUI Hardware Test")

	file, _ := mediaList.Open("media/bgr.png")
	stat, _ := file.Stat()
	size := stat.Size()
	buf := make([]byte, size)
	if _, err := file.Read(buf); err != nil {
		println(err.Error())
	}

	img := rl.LoadImageFromMemory(".png", buf, int32(size))
	texTspPad := rl.LoadTextureFromImage(img)
	pressedButtonCodes := mapset.NewSet[int32]()
	shouldExit := false
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() && !shouldExit {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		if rl.IsGamepadAvailable(0) {

			pressedButtonCodes.Add(rl.GetGamepadButtonPressed())
			for val := range pressedButtonCodes.Iter() {
				if !rl.IsGamepadButtonDown(0, val) {
					pressedButtonCodes.Remove(val)
				}
			}
			if pressedButtonCodes.Contains(14) && pressedButtonCodes.Contains(15) {
				shouldExit = true
			}
		}

		rl.DrawTexture(texTspPad, 0, 0, rl.White)
		rl.EndDrawing()
	}
	rl.UnloadTexture(texTspPad)
	rl.CloseWindow()
}
