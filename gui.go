//+build !test

package main

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/olevegard/nuklear/nk"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
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

	var err error
	sdl.Init(sdl.INIT_EVERYTHING)

	win, err := sdl.CreateWindow("FactoriGo Prototype", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		closer.Fatalln(err)
	}
	defer win.Destroy()

	context, err := sdl.GLCreateContext(win)
	if err != nil {
		closer.Fatalln(err)
	}

	width, height := win.GetSize()
	log.Printf("SDL2: created window %dx%d", width, height)

	if err := gl.Init(); err != nil {
		closer.Fatalln("opengl: init failed:", err)
	}
	gl.Viewport(0, 0, int32(width), int32(height))

	ctx := nk.NkPlatformInit(win, context, nk.PlatformInstallCallbacks)

	state := &State{
		bgColor: nk.NkRgba(28, 48, 62, 255),
	}

	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	state.smallFont = nk.NkFontAtlasAddFromBytes(atlas, MustAsset("assets/FreeSans.ttf"), 12, nil)
	state.bigFont = nk.NkFontAtlasAddFromBytes(atlas, MustAsset("assets/FreeSans.ttf"), 16, nil)
	nk.NkFontStashEnd()
	if state.smallFont != nil {
		nk.NkStyleSetFont(ctx, state.smallFont.Handle())
	}

	state.gameState = ReadDefaultGameState()
	renderer, _ := win.GetRenderer()

	for key, _ := range state.gameState.CurrentInventory.Items {
		newItem := state.gameState.CurrentInventory.Items[key]
		newItem.Icon = LoadIcon(renderer, fmt.Sprintf("assets/%s.png", key))
		state.gameState.CurrentInventory.Items[key] = newItem
	}

	exitC := make(chan struct{}, 1)
	doneC := make(chan struct{}, 1)
	closer.Bind(func() {
		close(exitC)
		<-doneC
	})

	fpsTicker := time.NewTicker(time.Second / 30)
	for {
		select {
		case <-exitC:
			nk.NkPlatformShutdown()
			sdl.Quit()
			fpsTicker.Stop()
			close(doneC)
			return
		case <-fpsTicker.C:
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					close(exitC)
				}
			}
			gfxMain(win, ctx, state)
		}
	}
}

func addProductionUnitLine(ctx *nk.Context, state *State, productionUnit ProductionUnit) ProductionUnit {
	nk.NkStyleSetFont(ctx, state.bigFont.Handle())
	nk.NkLayoutRowDynamic(ctx, 20, 1)
	{
		nk.NkLabel(ctx, productionUnit.Name, nk.TextCentered)
	}
	nk.NkStyleSetFont(ctx, state.smallFont.Handle())

	nk.NkLayoutRowDynamic(ctx, 20, 1)
	{
		nk.NkLabel(ctx, fmt.Sprintf("Count : %d", productionUnit.Count()), nk.TextCentered)
	}

	nk.NkLayoutRowDynamic(ctx, 20, 1)
	{
		nk.NkLabel(ctx, fmt.Sprintf("Time Left : %d", productionUnit.TicksRemaining), nk.TextLeft)
	}

	// TODO: The section for Generates and Requirements should be created from a function instead of duplicating
	// TODO: Find a way to display how much it generates each cycle ( change amount * ProductionUnit count)
	nk.NkLayoutRowDynamic(ctx, 20, 1)
	{
		nk.NkLabel(ctx, "Generates : ", nk.TextLeft)
	}

	nk.NkLayoutRowDynamic(ctx, 20, 5)
	{
		for _, change := range productionUnit.RecipeChangeSet {
			item := state.gameState.CurrentInventory.Items[change.InventoryItemId]
			if icon, ok := item.Icon.(nk.Image); ok {
				nk.NkImage(ctx, icon)
			} else {
				fmt.Printf("ERROR: Item.Icon is not NkImage \n\tActual type is %v\n", item.Icon)
			}
		}
	}

	nk.NkLayoutRowDynamic(ctx, 12, 5)
	{
		for _, change := range productionUnit.RecipeChangeSet {
			nk.NkLabel(ctx, strconv.Itoa(change.ChangeAmount), nk.TextCentered)
		}
	}

	nk.NkLayoutRowDynamic(ctx, 20, 1)
	{
		nk.NkLabel(ctx, "Requirements : ", nk.TextLeft)
	}
	nk.NkLayoutRowDynamic(ctx, 20, 5)
	{
		for _, change := range productionUnit.ChangeSetForBuildingNew {
			item := state.gameState.CurrentInventory.Items[change.InventoryItemId]
			if icon, ok := item.Icon.(nk.Image); ok {
				nk.NkImage(ctx, icon)
			} else {
				fmt.Printf("ERROR: Item.Icon is not NkImage \n\tActual type is %v\n", item.Icon)
			}
		}
	}

	nk.NkLayoutRowDynamic(ctx, 12, 5)
	{
		for _, change := range productionUnit.ChangeSetForBuildingNew {
			nk.NkLabel(ctx, strconv.Itoa(change.ChangeAmount*-1), nk.TextCentered)
		}
	}

	nk.NkLayoutRowDynamic(ctx, 20, 1)
	{
		if CanBuilNewProductionUnit(productionUnit, state.gameState.CurrentInventory) {
			if nk.NkButtonLabel(ctx, "Build") > 0 {
				productionUnit, state.gameState.CurrentInventory = BuilNewProductionUnit(productionUnit, state.gameState.CurrentInventory)
				return productionUnit
			}
		} else {
			if nk.NkButtonLabel(ctx, "Can't build") > 0 {
			}
		}
	}

	return productionUnit
}

func createProductionBox(ctx *nk.Context, state *State) Production {
	gameState := state.gameState
	bounds := nk.NkRect(winWidth-200, 0, 200, winHeight)
	update := nk.NkBegin(ctx, "Production", bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowClosable|nk.WindowTitle)

	if update > 0 {
		for index, unit := range gameState.CurrentProduction.ProductionUnits {
			gameState.CurrentProduction.ProductionUnits[index] = addProductionUnitLine(ctx, state, unit)
		}
	}

	nk.NkEnd(ctx)

	return gameState.CurrentProduction
}

func createInventory(ctx *nk.Context, inventory Inventory) Inventory {
	bounds := nk.NkRect(0, 0, 145, 350)
	update := nk.NkBegin(ctx, "Inventory", bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)

	if update == 0 {
		return inventory
	}

	keys := []string{"stone", "wood", "coal", "iron_ore", "copper_ore", "iron_plate", "copper_plate", "copper_wire", "circuit_board", "iron_gear", "simple_robotic_warrior"}

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
		nk.NkLayoutRowPush(ctx, 60)
		nk.NkLabel(ctx, fmt.Sprintf("%s :", printable), nk.TextLeft)

		nk.NkLayoutRowPush(ctx, 20)
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

func gfxMain(win *sdl.Window, ctx *nk.Context, state *State) {
	nk.NkPlatformNewFrame()
	if (state.ticksSinceStart % 20) == 0 {
		state.gameState = UpdateInventory(state.gameState)
	}
	state.gameState.CurrentInventory = createInventory(ctx, state.gameState.CurrentInventory)

	for _, window := range state.windows {
		if window.IsVisible {
			state.gameState.CurrentProduction = window.Update()
		}
	}
	createProductionBox(ctx, state)
	// Render
	bg := make([]float32, 4)
	nk.NkColorFv(bg, state.bgColor)
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
	sdl.GLSwapWindow(win)
}

type State struct {
	bgColor         nk.Color
	prop            int32
	gameState       GameState
	smallFont       *nk.Font
	bigFont         *nk.Font
	ticksSinceStart uint64

	wood          nk.Image
	coal          nk.Image
	stone         nk.Image
	circuit_board nk.Image
	iron_ore      nk.Image
	copper_ore    nk.Image
	iron_plate    nk.Image
	copper_plate  nk.Image
	copper_wire   nk.Image

	windows map[string]Window `json:"items"`
}

type Window struct {
	Update    func() Production
	IsVisible bool
}

func onError(code int32, msg string) {
	log.Printf("[glfw ERR]: error %d: %s", code, msg)
}

func LoadIcon(renderer *sdl.Renderer, filename string) nk.Image {
	surf, err := img.Load(filename)

	if err != nil {
		fmt.Printf("Failed loading image \n\tFilename : ''%s', \n\tReason : %s\n", filename, err)
	}
	rect := sdl.Rect{}
	surf.GetClipRect(&rect)
	var textureId uint32

	gl.GenTextures(1, &textureId)
	gl.BindTexture(gl.TEXTURE_2D, textureId)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, rect.W, rect.H, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&surf.Pixels()[0]))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return nk.NkImageId(int32(textureId))
}
