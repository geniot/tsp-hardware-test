package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/tevino/abool/v2"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/debug"
)

type Application struct {
	settings           *Settings
	resources          map[ResourceKey]*SurfTexture
	sdlWindow          *sdl.Window
	sdlRenderer        *sdl.Renderer
	sdlGameController  *sdl.GameController
	font               *ttf.Font
	joysticks          [16]*sdl.Joystick
	pressedKeysCodes   mapset.Set[sdl.Keycode]
	pressedButtonCodes mapset.Set[ButtonCode]
	axisValues         [4]float64
	maxAxisValues      [4]float64
	minAxisValues      [4]float64
	isRunning          *abool.AtomicBool
}

func NewApplication() *Application {
	return &Application{
		pressedKeysCodes:   mapset.NewSet[sdl.Keycode](),
		pressedButtonCodes: mapset.NewSet[ButtonCode](),
		isRunning:          abool.New(),
		resources:          make(map[int]*SurfTexture),
		settings:           NewSettings(),
		maxAxisValues:      [4]float64{JOYSTICK_INITIAL_MAX, JOYSTICK_INITIAL_MAX, JOYSTICK_INITIAL_MAX, JOYSTICK_INITIAL_MAX},
		minAxisValues:      [4]float64{-JOYSTICK_INITIAL_MAX, -JOYSTICK_INITIAL_MAX, -JOYSTICK_INITIAL_MAX, -JOYSTICK_INITIAL_MAX},
	}
}

func (app *Application) Start(args []string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Unhandled error: %v\n", r)
			log.Println("Stack trace:")
			debug.PrintStack()
			os.Exit(-1)
		}
	}()

	orPanic(sdl.Init(sdl.INIT_VIDEO | sdl.INIT_JOYSTICK | sdl.INIT_GAMECONTROLLER))

	orPanic(ttf.Init())
	app.font = orPanicRes(ttf.OpenFontRW(LoadMediaFile("pixelberry.ttf"), 1, 20))

	sdl.JoystickEventState(sdl.ENABLE)
	for i := 0; i < sdl.NumJoysticks(); i++ {
		if sdl.IsGameController(i) {
			app.sdlGameController = sdl.GameControllerOpen(i)
		}
	}
	println(runtime.GOARCH)
	if runtime.GOARCH == "arm64" { //most likely it's a TSP device
		orPanic(app.sdlGameController != nil)
	}

	app.sdlWindow = orPanicRes(sdl.CreateWindow(
		APP_NAME+" "+APP_VERSION,
		int32(app.settings.WindowPosX), int32(app.settings.WindowPosY),
		int32(app.settings.WindowWidth), int32(app.settings.WindowHeight),
		uint32(app.settings.WindowState)))

	app.sdlRenderer = orPanicRes(sdl.CreateRenderer(app.sdlWindow, -1, sdl.RENDERER_PRESENTVSYNC|sdl.RENDERER_ACCELERATED))
	app.initResources() //should be called after the creation of sdlRenderer

	app.isRunning.Set()
	for app.isRunning.IsSet() {
		app.UpdateEvents()
		app.UpdatePhysics()
		app.UpdateView()
	}
	app.releaseResources()
}

func (app *Application) Stop() {
	app.isRunning.UnSet()
}

func (app *Application) releaseResources() {
	app.settings.Save(app.sdlWindow)
	app.sdlGameController.Close()
	app.font.Close()
	ttf.Quit()
	sdl.Quit()
}

func (app *Application) UpdateEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {

		case *sdl.JoyAxisEvent:
			// Convert the value to a -1.0 - 1.0 range
			value := float64(t.Value) / 32768.0
			ind := If(t.Axis == 3 || t.Axis == 4, t.Axis-1, t.Axis)
			app.axisValues[ind] = value
			app.maxAxisValues[ind] = math.Max(value, app.maxAxisValues[ind])
			app.minAxisValues[ind] = math.Min(value, app.minAxisValues[ind])
			break

		case *sdl.ControllerButtonEvent:
			if t.State == sdl.PRESSED {
				println(t.Button)
				app.pressedButtonCodes.Add(t.Button)
			} else {
				app.pressedButtonCodes.Remove(t.Button)
			}
			break

		case *sdl.JoyDeviceAddedEvent:
			// Open joystick for use
			app.joysticks[int(t.Which)] = sdl.JoystickOpen(int(t.Which))
			if app.joysticks[int(t.Which)] != nil {
				fmt.Println("Joystick", t.Which, "connected")
			}
			break
		case *sdl.JoyDeviceRemovedEvent:
			if joystick := app.joysticks[int(t.Which)]; joystick != nil {
				joystick.Close()
			}
			fmt.Println("Joystick", t.Which, "disconnected")
			break

		case *sdl.KeyboardEvent:
			if t.Repeat > 0 {
				break
			}
			if t.State == sdl.PRESSED {
				app.pressedKeysCodes.Add(t.Keysym.Sym)
			} else { // if t.State == sdl.RELEASED {
				app.pressedKeysCodes.Remove(t.Keysym.Sym)
			}
			break

		case *sdl.WindowEvent:
			if t.Event == sdl.WINDOWEVENT_CLOSE {
				app.settings.SaveWindowState(app.sdlWindow)
			}
			break

		case *sdl.QuitEvent:
			app.Stop()
			break
		}
	}
}

func (app *Application) UpdatePhysics() {
	if app.pressedKeysCodes.Contains(sdl.K_q) || (app.pressedButtonCodes.Contains(BUTTON_CODE_MENU) && app.pressedButtonCodes.Contains(BUTTON_CODE_START)) {
		app.Stop()
	}
}

func (app *Application) UpdateView() {
	orPanic(app.sdlRenderer.SetDrawColorArray(BACKGROUND_COLOR.R, BACKGROUND_COLOR.G, BACKGROUND_COLOR.B, BACKGROUND_COLOR.A))
	orPanic(app.sdlRenderer.Clear())

	orWarn(app.sdlRenderer.Copy(app.resources[RESOURCE_BGR_KEY].T, nil, &sdl.Rect{X: 0, Y: 0, W: app.resources[RESOURCE_BGR_KEY].W, H: app.resources[RESOURCE_BGR_KEY].H}))

	app.renderJoystick(Reactors[BUTTON_CODE_LEFT_JOYSTICK].OffsetX, Reactors[BUTTON_CODE_LEFT_JOYSTICK].OffsetY, 0, 1, sdl.K_l)
	app.renderJoystick(Reactors[BUTTON_CODE_RIGHT_JOYSTICK].OffsetX, Reactors[BUTTON_CODE_RIGHT_JOYSTICK].OffsetY, 2, 3, sdl.K_r)

	//used to debug joystick values
	//for i := 0; i < len(app.axisValues); i++ {
	//	var val = float64(app.axisValues[i])
	//	var valStr = strconv.FormatFloat(val, 'f', 6, 64)
	//	app.drawText(valStr, int32(If(val < 0, 5, 18)), int32(i*30))
	//}
	//var procent = (app.axisValues[0] * 100 / If(app.axisValues[0] > 0, app.maxAxisValues[0], -app.minAxisValues[0])) / 100
	//var valStr = strconv.FormatFloat(procent, 'f', 6, 64)
	//app.drawText(valStr, 5, 120)

	for val := range app.pressedButtonCodes.Iter() {
		if Reactors[val] != nil {
			width := If(Reactors[val].Width == 0, app.resources[Reactors[val].ResourceKey].W, Reactors[val].Width)
			height := If(Reactors[val].Height == 0, app.resources[Reactors[val].ResourceKey].H, Reactors[val].Height)
			orWarn(app.sdlRenderer.Copy(app.resources[Reactors[val].ResourceKey].T, nil, &sdl.Rect{X: Reactors[val].OffsetX, Y: Reactors[val].OffsetY, W: width, H: height}))
		}
	}
	app.sdlRenderer.Present()
}

func (app *Application) renderJoystick(posX, posY int32, axisIndexX, axisIndexY int, debugKeyCode sdl.Keycode) {
	var (
		axisX                 = app.axisValues[axisIndexX]
		axisY                 = app.axisValues[axisIndexY]
		roundedX              = toFixed(axisX, 3)
		roundedY              = toFixed(axisY, 3)
		correctedX            = (app.axisValues[axisIndexX] * 100 / If(app.axisValues[axisIndexX] > 0, app.maxAxisValues[axisIndexX], -app.minAxisValues[axisIndexX])) / 100
		correctedY            = (app.axisValues[axisIndexY] * 100 / If(app.axisValues[axisIndexY] > 0, app.maxAxisValues[axisIndexY], -app.minAxisValues[axisIndexY])) / 100
		correctedScreenWidth  = SCREEN_WIDTH - app.resources[RESOURCE_CROSS_YELLOW_KEY].W
		correctedScreenHeight = SCREEN_HEIGHT - app.resources[RESOURCE_CROSS_YELLOW_KEY].H
	)
	//drawing yellow joystick circles
	if app.pressedKeysCodes.Contains(debugKeyCode) || (roundedX != 0 || roundedY != 0) {
		orWarn(app.sdlRenderer.Copy(app.resources[RESOURCE_CIRCLE_YELLOW_KEY].T, nil,
			&sdl.Rect{X: posX, Y: posY, W: app.resources[RESOURCE_CIRCLE_YELLOW_KEY].W, H: app.resources[RESOURCE_CIRCLE_YELLOW_KEY].H}))
	}
	//cross-hairs
	if roundedX != 0 || roundedY != 0 {
		orWarn(app.sdlRenderer.Copy(app.resources[RESOURCE_CROSS_YELLOW_KEY].T, nil,
			&sdl.Rect{
				X: SCREEN_LEFT_UP_X + correctedScreenWidth/2 + int32(float64(correctedScreenWidth/2)*correctedX),
				Y: SCREEN_LEFT_UP_Y + correctedScreenHeight/2 + int32(float64(correctedScreenHeight/2)*correctedY),
				W: app.resources[RESOURCE_CROSS_YELLOW_KEY].W,
				H: app.resources[RESOURCE_CROSS_YELLOW_KEY].H}))
	}
}

func (app *Application) initResources() {
	app.resources[RESOURCE_BGR_KEY] = LoadSurfTexture("bgr.png", app.sdlRenderer)
	app.resources[RESOURCE_CIRCLE_YELLOW_KEY] = LoadSurfTexture("circle_yellow.png", app.sdlRenderer)
	app.resources[RESOURCE_CROSS_YELLOW_KEY] = LoadSurfTexture("cross_yellow.png", app.sdlRenderer)
}

func (app *Application) drawText(val string, x, y int32) {
	textSurface, _ := app.font.RenderUTF8Blended(val, COLOR_BLACK)
	defer textSurface.Free()
	textTexture, _ := app.sdlRenderer.CreateTextureFromSurface(textSurface)
	orWarn(app.sdlRenderer.Copy(textTexture, nil, &sdl.Rect{X: x, Y: y, W: textSurface.W, H: textSurface.H}))
	defer orWarn(textTexture.Destroy())
}
