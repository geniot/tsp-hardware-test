package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ButtonCode = uint8
type ResourceKey = int

const (
	APP_NAME          = "TSP Hardware Test"
	APP_VERSION       = "0.1"
	TSP_SCREEN_WIDTH  = 1280
	TSP_SCREEN_HEIGHT = 720
)

var (
	COLOR_RED    = sdl.Color{R: 192, G: 64, B: 64, A: 255}
	COLOR_GREEN  = sdl.Color{R: 64, G: 192, B: 64, A: 255}
	COLOR_GRAY   = sdl.Color{R: 192, G: 192, B: 192, A: 255}
	COLOR_WHITE  = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	COLOR_PURPLE = sdl.Color{R: 255, G: 0, B: 255, A: 255}
	COLOR_YELLOW = sdl.Color{R: 255, G: 255, B: 0, A: 255}
	COLOR_BLUE   = sdl.Color{R: 0, G: 255, B: 255, A: 255}
	COLOR_BLACK  = sdl.Color{R: 0, G: 0, B: 0, A: 255}

	BACKGROUND_COLOR = COLOR_BLACK
)

const (
	RESOURCE_BGR_KEY           = ResourceKey(iota)
	RESOURCE_CIRCLE_YELLOW_KEY = ResourceKey(iota)
	RESOURCE_CROSS_YELLOW_KEY  = ResourceKey(iota)
)

const (
	BUTTON_CODE_MENU           = ButtonCode(5)
	BUTTON_CODE_START          = ButtonCode(6)
	BUTTON_CODE_LEFT_JOYSTICK  = ButtonCode(14)
	BUTTON_CODE_RIGHT_JOYSTICK = ButtonCode(15)
)

const (
	SCREEN_LEFT_UP_X     = int32(361)
	SCREEN_LEFT_UP_Y     = int32(192)
	SCREEN_RIGHT_DOWN_X  = int32(919)
	SCREEN_RIGHT_DOWN_Y  = int32(508)
	SCREEN_WIDTH         = SCREEN_RIGHT_DOWN_X - SCREEN_LEFT_UP_X
	SCREEN_HEIGHT        = SCREEN_RIGHT_DOWN_Y - SCREEN_LEFT_UP_Y
	JOYSTICK_INITIAL_MAX = 0.9
)

var (
	Reactors = map[ButtonCode]*ImageDescriptor{
		BUTTON_CODE_LEFT_JOYSTICK: {
			OffsetX:     283,
			OffsetY:     391,
			ResourceKey: RESOURCE_CIRCLE_YELLOW_KEY,
		},
		BUTTON_CODE_RIGHT_JOYSTICK: {
			OffsetX:     955,
			OffsetY:     391,
			ResourceKey: RESOURCE_CIRCLE_YELLOW_KEY,
		},
	}
)
