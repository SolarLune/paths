package main

import (
	"fmt"
	"math"

	"github.com/SolarLune/paths"
	"github.com/veandco/go-sdl2/sdl"
)

type PathDrawer struct {
	World            *World1
	StartX, StartY   int
	TargetX, TargetY int
	MouseCX, MouseCY int
	mouseClicked     bool
	Path             *paths.Path
}

func NewPathDrawer(x, y int, world *World1) *PathDrawer {
	f := &PathDrawer{StartX: x, StartY: y, World: world}
	f.Path = &paths.Path{}
	return f
}

func (pd *PathDrawer) Update() {

	x, y, lb := sdl.GetMouseState()

	w, h := window.GetSize()
	dx := float32(w) / float32(screenWidth)
	dy := float32(h) / float32(screenHeight)

	pd.MouseCX = int(math.Floor(float64(float32(x)/dx) / 16))
	pd.MouseCY = int(math.Floor(float64(float32(y)/dy) / 16))

	pd.TargetX = pd.MouseCX * 16
	pd.TargetY = pd.MouseCY * 16

	if lb == sdl.Button(sdl.BUTTON_LEFT) {
		if !pd.mouseClicked && pd.World.GameMap.Get(pd.MouseCX, pd.MouseCY) != nil {
			pd.StartX = pd.TargetX
			pd.StartY = pd.TargetY
			pd.mouseClicked = true
		}
	} else if lb == sdl.Button(sdl.BUTTON_RIGHT) {
		if !pd.mouseClicked {
			c := pd.World.GameMap.Get(pd.MouseCX, pd.MouseCY)
			if c != nil {
				c.Walkable = !c.Walkable
				pd.mouseClicked = true
			}
		}
	} else {
		pd.mouseClicked = false
	}

	if keyboard.KeyPressed(sdl.K_s) {
		c := pd.World.GameMap.Get(pd.MouseCX, pd.MouseCY)
		if c != nil {
			c.Cost += 5
		}
	}
	if keyboard.KeyPressed(sdl.K_a) {
		c := pd.World.GameMap.Get(pd.MouseCX, pd.MouseCY)
		if c != nil {
			c.Cost -= 5
		}
		if c.Cost < 1 {
			c.Cost = 1
		}
	}

}

func (pd *PathDrawer) Draw() {

	sc := pd.World.GameMap.Get(pd.StartX/16, pd.StartY/16)
	tc := pd.World.GameMap.Get(pd.TargetX/16, pd.TargetY/16)

	if tc != nil {

		newPath := pd.World.GameMap.GetPathFromCells(sc, tc, false)
		if !newPath.Same(pd.Path) {
			pd.Path = newPath
		}

		if pd.Path != nil {
			for i, c := range pd.Path.Cells {
				renderer.SetDrawColor(255-uint8(i*8%100), 255-uint8(i*8%100), 0, 255)
				renderer.FillRect(&sdl.Rect{int32(c.X * 16), int32(c.Y * 16), 16, 16})
			}

			if pd.Path.AtEnd() {
				pd.Path.Restart()
			}
			c := pd.Path.Next()
			renderer.SetDrawColor(255, 0, 0, 255)
			renderer.FillRect(&sdl.Rect{int32(c.X * 16), int32(c.Y * 16), 16, 16})
		}

	}

	renderer.SetDrawColor(0, 0, 255, 255)
	renderer.DrawRect(&sdl.Rect{int32(pd.TargetX), int32(pd.TargetY), 16, 16})

}

type World1 struct {
	GameMap    *paths.Grid
	PathDrawer *PathDrawer
	Tileset    *sdl.Texture
}

func (world *World1) Create() {

	layout := []string{
		"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		"x   f x  x     x               x",
		"x xxxxxxxxxxxx x               x",
		"x            x x               x",
		"x x xxxxxxxx x x               x",
		"x x x x    x x x       xxxx    x",
		"x x x   xx x x x               x",
		"x x xxxxxx x x x        xx     x",
		"x xxx      x x x               x",
		"x              xxxxxxxxx  xxxxxx",
		"x xxx      xxx x                ",
		"x            x x                ",
		"x xxxxxxxxxx x x xxxxxxxxxxxx   ",
		"x              x                ",
		"xxxxxxxxxxxxxxxx                ",
	}

	world.GameMap = paths.NewGridFromStringArrays(layout, 16, 16)

	spawn := world.GameMap.CellsByRune('f')[0]
	for _, cell := range world.GameMap.CellsByRune('x') {
		cell.Walkable = false
	}
	world.PathDrawer = NewPathDrawer(spawn.X*16, spawn.Y*16, world)

}

func (world *World1) Update() {

	world.PathDrawer.Update()

}

func (world *World1) Draw() {

	for y := 0; y < world.GameMap.Height(); y++ {
		for x := 0; x < world.GameMap.Width(); x++ {
			if !world.GameMap.Get(x, y).Walkable {
				renderer.SetDrawColor(255, 255, 255, 255)
			} else {
				c := world.GameMap.Get(x, y).Cost
				renderer.SetDrawColor(0, 0, 40+uint8(c*8), 50)
			}
			renderer.FillRect(&sdl.Rect{int32(x * 16), int32(y * 16), 16, 16})
		}
	}

	world.PathDrawer.Draw()

	if drawHelpText {

		c := world.GameMap.Get(world.PathDrawer.MouseCX, world.PathDrawer.MouseCY)

		cellInfo := "Current cell: Not on map"

		if c != nil {
			cellInfo = fmt.Sprintf("Current cell: X:%d, Y:%d, Walkable:%t, Cost:%d", c.X, c.Y, c.Walkable, c.Cost)
		}

		pathInfo := "Path cost: No path"
		if world.PathDrawer.Path != nil {
			pathInfo = fmt.Sprintf("Path cost: %d", world.PathDrawer.Path.TotalCost())
		}

		DrawText(16, 240,
			"paths example",
			"Mouse position: End point of path",
			"Left click: Set starting point of path",
			"Right click: Set walkability of cell",
			"A key: Lower cost of cell",
			"S key: Raise cost of cell",
			cellInfo,
		)

		DrawText(320, 240,
			pathInfo,
		)

	}
}

func (world *World1) Destroy() {

}
