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
	winWidth  = 600
	winHeight = 500

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
	win, err := glfw.CreateWindow(winWidth, winHeight, "Nuklear Demo", nil, nil)
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
	sansFont := nk.NkFontAtlasAddFromBytes(atlas, MustAsset("assets/FreeSans.ttf"), 16, nil)
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
	state.gameState = MakeDefaultGameState()
	gs := state.gameState
	fmt.Printf("Begin I : %d C %d IP : %d CP : %d IM  : %d CM : %d IS : %d CS : %d\n",
		gs.inventory.iron_ore, gs.inventory.copper_ore, gs.inventory.iron_plates, gs.inventory.copper_plates,
		gs.production.iron_mines.count, gs.production.copper_mines.count, gs.production.iron_smelters.count, gs.production.copper_smelters.count)

	fpsTicker := time.NewTicker(time.Second / 100)
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

func createInventory(ctx *nk.Context, inventory Inventory) Inventory {
	bounds := nk.NkRect(20, 50, 200, 150)
	update := nk.NkBegin(ctx, "Inventory", bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)

	if update == 0 {
		return inventory
	}

	addLine(ctx, "Iron Ore", inventory.iron_ore, func() {
		inventory.iron_ore++
	})

	addLine(ctx, "Copper Ore", inventory.copper_ore, func() {
		inventory.copper_ore++
	})

	addLine(ctx, "Iron Plates", inventory.iron_plates, func() {
		inventory.iron_plates++
	})

	addLine(ctx, "Copper Ore", inventory.copper_plates, func() {
		inventory.copper_plates++
	})

	nk.NkEnd(ctx)

	return inventory
}

func createProduction(ctx *nk.Context, production Production) Production {
	bounds := nk.NkRect(240, 50, 200, 150)
	update := nk.NkBegin(ctx, "Production", bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)

	if update == 0 {
		return production
	}

	addLine(ctx, "Iron Mines", production.iron_mines.count, func() {
		production.iron_mines.count++
	})

	addLine(ctx, "Copper Mines", production.copper_mines.count, func() {
		production.copper_mines.count++
	})

	addLine(ctx, "Iron Smelter", production.iron_smelters.count, func() {
		production.iron_smelters.count++
	})

	addLine(ctx, "Copper Smelter", production.copper_smelters.count, func() {
		production.copper_smelters.count++
	})

	nk.NkEnd(ctx)
	return production
}

func addLine(ctx *nk.Context, name string, count int, f func()) {
	nk.NkLayoutRowBegin(ctx, nk.LayoutDynamic, 20, 3)
	{
		nk.NkLayoutRowPush(ctx, 80)
		nk.NkLabel(ctx, name+" : ", nk.TextLeft)

		nk.NkLayoutRowPush(ctx, 40)
		nk.NkLabel(ctx, strconv.Itoa(count), nk.TextRight)

		nk.NkLayoutRowPush(ctx, 20)
		if nk.NkButtonLabel(ctx, "+") > 0 {
			f()
		}
	}

	nk.NkLayoutRowEnd(ctx)
}

func gfxMain(win *glfw.Window, ctx *nk.Context, state *State) {
	nk.NkPlatformNewFrame()
	state.gameState = UpdateInventory(state.gameState)
	gs := state.gameState
	fmt.Printf("Begin I : %d C %d IP : %d CP : %d IM  : %d CM : %d IS : %d CS : %d\n",
		gs.inventory.iron_ore, gs.inventory.copper_ore, gs.inventory.iron_plates, gs.inventory.copper_plates,
		gs.production.iron_mines.count, gs.production.copper_mines.count, gs.production.iron_smelters.count, gs.production.copper_smelters.count)

	state.gameState.production = createProduction(ctx, state.gameState.production)
	state.gameState.inventory = createInventory(ctx, state.gameState.inventory)

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
}

func onError(code int32, msg string) {
	log.Printf("[glfw ERR]: error %d: %s", code, msg)
}
