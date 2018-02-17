//+build !test

package main

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang-ui/nuklear/nk"
	"github.com/xlab/closer"
)

// https://github.com/golang-ui/nuklear
const (
	winWidth  = 360
	winHeight = 640

	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		closer.Fatalln(err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	win, err := glfw.CreateWindow(winWidth, winHeight, "FactoriGo Prototype", nil, nil)
	if err != nil {
		closer.Fatalln(err)
	}
	win.MakeContextCurrent()

	width, height := win.GetSize()
	log.Printf("glfw: created window %dx%d", width, height)

	if err := gl.Init(); err != nil {
		closer.Fatalln("opengl: init failed:", err)
	}
	gl.Viewport(0, 0, int32(width), int32(height))

	ctx := nk.NkPlatformInit(win, nk.PlatformInstallCallbacks)

	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	sansFont := nk.NkFontAtlasAddFromBytes(atlas, MustAsset("assets/FreeSans.ttf"), 14, nil)
	nk.NkFontStashEnd()
	if sansFont != nil {
		nk.NkStyleSetFont(ctx, sansFont.Handle())
	}

	exitC := make(chan struct{}, 1)
	doneC := make(chan struct{}, 1)
	closer.Bind(func() {
		close(exitC)
		<-doneC
	})

	state := &State{
		bgColor: nk.NkRgba(28, 48, 62, 255),
	}
	state.gameState = ReadDefaultGameState()

	fpsTicker := time.NewTicker(time.Second / 100)

	state.windows = map[string]Window{}

	for index, _ := range state.gameState.CurrentProduction.ProductionUnits {
		key, w := createProductionBox(ctx, state, index)
		state.windows[key] = w
	}

	for {
		select {
		case <-exitC:
			nk.NkPlatformShutdown()
			glfw.Terminate()
			fpsTicker.Stop()
			close(doneC)
			return
		case <-fpsTicker.C:
			if win.ShouldClose() {
				close(exitC)
				continue
			}
			glfw.PollEvents()
			gfxMain(win, ctx, state)
		}
	}
}

func createProductionBox(ctx *nk.Context, state *State, index int) (string, Window) {
	productionUnit := state.gameState.CurrentProduction.ProductionUnits[index]
	t := func() Production {
		bounds := nk.NkRect(winWidth-150, float32(80*index), 150, 80)
		update := nk.NkBegin(ctx, fmt.Sprintf("%s : %d", productionUnit.Name, productionUnit.Count()), bounds,
			nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowClosable|nk.WindowTitle)

		if update > 0 {
			/*
				nk.NkLayoutRowDynamic(ctx, 20, 1)
				{
					nk.NkLabel(ctx, fmt.Sprintf("Count : %d", productionUnit.Count()), nk.TextLeft)
				}

					nk.NkLayoutRowDynamic(ctx, 20, 1)
					{
						nk.NkLabel(ctx, "Generates : 10 x IO", nk.TextLeft)
					}

					nk.NkLayoutRowDynamic(ctx, 20, 1)
					{
						nk.NkLabel(ctx, fmt.Sprintf("Time : %d", productionUnit.TicksPerCycle), nk.TextLeft)
					}
			*/

			nk.NkLayoutRowDynamic(ctx, 20, 1)
			{
				if nk.NkButtonLabel(ctx, "Build") > 0 {
					state.gameState.CurrentProduction.ProductionUnits[index], state.gameState.CurrentInventory = BuilNewProductionUnit(productionUnit, state.gameState.CurrentInventory)
				}
			}
			/*
				nk.NkLayoutRowDynamic(ctx, 20, 1)
				{
					if nk.NkButtonLabel(ctx, "Close") > 0 {

						oldWindow := Window{state.windows[productionUnit.Name].Update, false}
						state.windows[productionUnit.Name] = oldWindow
					}
				}
			*/
		}
		nk.NkEnd(ctx)

		return state.gameState.CurrentProduction
	}
	return productionUnit.Name, Window{t, true}
}

func createInventory(ctx *nk.Context, inventory Inventory) Inventory {
	bounds := nk.NkRect(0, 0, 240, 350)
	update := nk.NkBegin(ctx, "Inventory", bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)

	if update == 0 {
		return inventory
	}

	keys := []string{"stone", "wood", "coal", "iron_ore", "copper_ore", "iron_plates", "copper_plates", "copper_wire", "circuit_board"}

	for _, key := range keys {

		addLine(ctx, inventory.Items[key], func() {

		},
			func() {
				_, inventory = ApplyInventoryItemChange(inventory, NewInventoryChange(key, 1))
			}, inventory.Items[key].CanBeMadeManually)
	}

	nk.NkEnd(ctx)

	return inventory
}

func addLine(ctx *nk.Context, printable Printable, info func(), createNew func(), actionable bool) {
	nk.NkLayoutRowBegin(ctx, nk.LayoutDynamic, 20, 4)
	{
		nk.NkLayoutRowPush(ctx, 80)
		nk.NkLabel(ctx, fmt.Sprintf("%s :", printable), nk.TextLeft)

		nk.NkLayoutRowPush(ctx, 30)
		nk.NkLabel(ctx, strconv.Itoa(printable.Count()), nk.TextRight)

		/*
			nk.NkLayoutRowPush(ctx, 30)
			if nk.NkButtonLabel(ctx, "info") > 0 {
				info()
			}
		*/

		nk.NkLayoutRowPush(ctx, 30)
		if actionable && nk.NkButtonLabel(ctx, "+") > 0 {
			createNew()
		}
	}

	nk.NkLayoutRowEnd(ctx)
}

func gfxMain(win *glfw.Window, ctx *nk.Context, state *State) {
	nk.NkPlatformNewFrame()
	state.gameState = UpdateInventory(state.gameState)
	state.gameState.CurrentInventory = createInventory(ctx, state.gameState.CurrentInventory)

	for _, window := range state.windows {
		if window.IsVisible {
			state.gameState.CurrentProduction = window.Update()
		}
	}
	// Render
	bg := make([]float32, 4)
	nk.NkColorFv(bg, state.bgColor)
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
	win.SwapBuffers()
}

type State struct {
	bgColor   nk.Color
	prop      int32
	gameState GameState

	windows map[string]Window `json:"items"`
}

type Window struct {
	Update    func() Production
	IsVisible bool
}

func onError(code int32, msg string) {
	log.Printf("[glfw ERR]: error %d: %s", code, msg)
}
