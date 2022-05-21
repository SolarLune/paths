/*
Package paths is a simple library written in Go made to handle 2D pathfinding for games. All you need to do is generate a Grid,
specify which cells aren't walkable, optionally change the cost on specific cells, and finally get a path from one cell to
another.
*/
package paths

import (
	"fmt"
	"math"
	"sort"
)

// A Cell represents a point on a Grid map. It has an X and Y value for the position, a Cost, which influences which Cells are
// ideal for paths, Walkable, which indicates if the tile can be walked on or should be avoided, and a Rune, which indicates
// which rune character the Cell is represented by.
type Cell struct {
	X, Y     int
	Cost     float64
	Walkable bool
	Rune     rune
}

func (cell Cell) String() string {
	return fmt.Sprintf("X:%d Y:%d Cost:%f Walkable:%t Rune:%s(%d)", cell.X, cell.Y, cell.Cost, cell.Walkable, string(cell.Rune), int(cell.Rune))
}

// Grid represents a "map" composed of individual Cells at each point in the map.
// Data is a 2D array of Cells.
// CellWidth and CellHeight indicate the size of Cells for Cell Position <-> World Position translation.
type Grid struct {
	Data                  [][]*Cell
	CellWidth, CellHeight int
}

// NewGrid returns a new Grid of (gridWidth x gridHeight) size. cellWidth and cellHeight changes the size of each Cell in the Grid.
// This is used to translate world position to Cell positions (i.e. the Cell position [2, 5] with a CellWidth and CellHeight of
// [16, 16] would be the world positon [32, 80]).
func NewGrid(gridWidth, gridHeight, cellWidth, cellHeight int) *Grid {

	m := &Grid{CellWidth: cellWidth, CellHeight: cellHeight}

	for y := 0; y < gridHeight; y++ {
		m.Data = append(m.Data, []*Cell{})
		for x := 0; x < gridWidth; x++ {
			m.Data[y] = append(m.Data[y], &Cell{x, y, 1, true, ' '})
		}
	}
	return m
}

// NewGridFromStringArrays creates a Grid map from a 1D array of strings. Each string becomes a row of Cells, each
// with one rune as its character. cellWidth and cellHeight changes the size of each Cell in the Grid. This is used to
// translate world position to Cell positions (i.e. the Cell position [2, 5] with a CellWidth and CellHeight of
// [16, 16] would be the world positon [32, 80]).
func NewGridFromStringArrays(arrays []string, cellWidth, cellHeight int) *Grid {

	m := &Grid{CellWidth: cellWidth, CellHeight: cellHeight}

	for y := 0; y < len(arrays); y++ {
		m.Data = append(m.Data, []*Cell{})
		stringLine := []rune(arrays[y])
		for x := 0; x < len(arrays[y]); x++ {
			m.Data[y] = append(m.Data[y], &Cell{X: x, Y: y, Cost: 1, Walkable: true, Rune: stringLine[x]})
		}
	}

	return m

}

// NewGridFromRuneArrays creates a Grid map from a 2D array of runes. Each individual Rune becomes a Cell in the resulting
// Grid. cellWidth and cellHeight changes the size of each Cell in the Grid. This is used to translate world position to Cell
// positions (i.e. the Cell position [2, 5] with a CellWidth and CellHeight of [16, 16] would be the world positon [32, 80]).
func NewGridFromRuneArrays(arrays [][]rune, cellWidth, cellHeight int) *Grid {

	m := &Grid{CellWidth: cellWidth, CellHeight: cellHeight}

	for y := 0; y < len(arrays); y++ {
		m.Data = append(m.Data, []*Cell{})
		for x := 0; x < len(arrays[y]); x++ {
			m.Data[y] = append(m.Data[y], &Cell{X: x, Y: y, Cost: 1, Walkable: true, Rune: arrays[y][x]})
		}
	}

	return m

}

// DataToString returns a string, used to easily identify the Grid map.
func (m *Grid) DataToString() string {
	s := ""
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			s += string(m.Data[y][x].Rune) + " "
		}
		s += "\n"
	}
	return s
}

// Get returns a pointer to the Cell in the x and y position provided.
func (m *Grid) Get(x, y int) *Cell {
	if x < 0 || y < 0 || x >= m.Width() || y >= m.Height() {
		return nil
	}
	return m.Data[y][x]
}

// Height returns the height of the Grid map.
func (m *Grid) Height() int {
	return len(m.Data)
}

// Width returns the width of the Grid map.
func (m *Grid) Width() int {
	return len(m.Data[0])
}

// CellsByRune returns a slice of pointers to Cells that all have the character provided.
func (m *Grid) CellsByRune(char rune) []*Cell {

	cells := make([]*Cell, 0)

	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			c := m.Get(x, y)
			if c.Rune == char {
				cells = append(cells, c)
			}
		}
	}

	return cells

}

// AllCells returns a single slice of pointers to all Cells contained in the Grid's 2D Data array.
func (m *Grid) AllCells() []*Cell {

	cells := make([]*Cell, 0)

	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			cells = append(cells, m.Get(x, y))
		}
	}

	return cells

}

// CellsByCost returns a slice of pointers to Cells that all have the Cost value provided.
func (m *Grid) CellsByCost(cost float64) []*Cell {

	cells := make([]*Cell, 0)

	for y := 0; y < m.Height(); y++ {

		for x := 0; x < m.Width(); x++ {

			c := m.Get(x, y)
			if c.Cost == cost {
				cells = append(cells, c)
			}

		}

	}

	return cells

}

// CellsByWalkable returns a slice of pointers to Cells that all have the Cost value provided.
func (m *Grid) CellsByWalkable(walkable bool) []*Cell {

	cells := make([]*Cell, 0)

	for y := 0; y < m.Height(); y++ {

		for x := 0; x < m.Width(); x++ {

			c := m.Get(x, y)
			if c.Walkable == walkable {
				cells = append(cells, c)
			}

		}

	}

	return cells

}

// SetWalkable sets walkability across all cells in the Grid with the specified rune.
func (m *Grid) SetWalkable(char rune, walkable bool) {

	for y := 0; y < m.Height(); y++ {

		for x := 0; x < m.Width(); x++ {
			cell := m.Get(x, y)
			if cell.Rune == char {
				cell.Walkable = walkable
			}
		}

	}

}

// SetCost sets the movement cost across all cells in the Grid with the specified rune.
func (m *Grid) SetCost(char rune, cost float64) {

	for y := 0; y < m.Height(); y++ {

		for x := 0; x < m.Width(); x++ {
			cell := m.Get(x, y)
			if cell.Rune == char {
				cell.Cost = cost
			}
		}

	}

}

// GridToWorld converts from a grid position to world position, multiplying the value by the CellWidth and CellHeight of the Grid.
func (m *Grid) GridToWorld(x, y int) (float64, float64) {
	rx := float64(x * m.CellWidth)
	ry := float64(y * m.CellHeight)
	return rx, ry
}

// WorldToGrid converts from a grid position to world position, multiplying the value by the CellWidth and CellHeight of the Grid.
func (m *Grid) WorldToGrid(x, y float64) (int, int) {
	tx := int(math.Floor(x / float64(m.CellWidth)))
	ty := int(math.Floor(y / float64(m.CellHeight)))
	return tx, ty
}

// GetPathFromCells returns a Path, from the starting Cell to the destination Cell. diagonals controls whether moving diagonally
// is acceptable when creating the Path. wallsBlockDiagonals indicates whether to allow diagonal movement "through" walls that are
// positioned diagonally.
func (m *Grid) GetPathFromCells(start, dest *Cell, diagonals, wallsBlockDiagonals bool) *Path {

	type Node struct {
		Cell   *Cell
		Parent *Node
		Cost   float64
	}

	openNodes := []*Node{&Node{Cell: dest, Cost: dest.Cost}}

	// checkedNodes := make([]*Cell, 0)
	checkedNodes := make(map[*Cell]struct{})

	hasBeenAdded := func(cell *Cell) bool {
		_, ok := checkedNodes[cell]
		return ok
		// for _, c := range checkedNodes {
		// 	if cell == c {
		// 		return true
		// 	}
		// }
		// return false

	}

	path := &Path{}

	if !start.Walkable || !dest.Walkable {
		return nil
	}

	for {

		// If the list of openNodes (nodes to check) is at 0, then we've checked all Nodes, and so the function can quit.
		if len(openNodes) == 0 {
			break
		}

		node := openNodes[0]
		openNodes = openNodes[1:]

		// If we've reached the start, then we've constructed our Path going from the destination to the start; we just have
		// to loop through each Node and go up, adding it and its parents recursively to the path.
		if node.Cell == start {

			var t = node
			for {
				path.Cells = append(path.Cells, t.Cell)
				t = t.Parent
				if t == nil {
					break
				}
			}

			break
		}

		// Otherwise, we add the current node's neighbors to the list of cells to check, and list of cells that have already been
		// checked (so we don't get nodes being checked multiple times).
		if node.Cell.X > 0 {
			c := m.Get(node.Cell.X-1, node.Cell.Y)
			n := &Node{c, node, c.Cost + node.Cost}
			if n.Cell.Walkable && !hasBeenAdded(n.Cell) {
				openNodes = append(openNodes, n)
				checkedNodes[n.Cell] = struct{}{}
				// checkedNodes = append(checkedNodes, n.Cell)
			}
		}
		if node.Cell.X < m.Width()-1 {
			c := m.Get(node.Cell.X+1, node.Cell.Y)
			n := &Node{c, node, c.Cost + node.Cost}
			if n.Cell.Walkable && !hasBeenAdded(n.Cell) {
				openNodes = append(openNodes, n)
				checkedNodes[n.Cell] = struct{}{}
				// checkedNodes = append(checkedNodes, n.Cell)
			}
		}

		if node.Cell.Y > 0 {
			c := m.Get(node.Cell.X, node.Cell.Y-1)
			n := &Node{c, node, c.Cost + node.Cost}
			if n.Cell.Walkable && !hasBeenAdded(n.Cell) {
				openNodes = append(openNodes, n)
				checkedNodes[n.Cell] = struct{}{}
				// checkedNodes = append(checkedNodes, n.Cell)
			}
		}
		if node.Cell.Y < m.Height()-1 {
			c := m.Get(node.Cell.X, node.Cell.Y+1)
			n := &Node{c, node, c.Cost + node.Cost}
			if n.Cell.Walkable && !hasBeenAdded(n.Cell) {
				openNodes = append(openNodes, n)
				checkedNodes[n.Cell] = struct{}{}
				// checkedNodes = append(checkedNodes, n.Cell)
			}
		}

		// Do the same thing for diagonals.
		if diagonals {

			diagonalCost := .414 // Diagonal movement is slightly slower, so we should prioritize straightaways if possible

			up := m.Get(node.Cell.X, node.Cell.Y-1).Walkable
			down := m.Get(node.Cell.X, node.Cell.Y+1).Walkable
			left := m.Get(node.Cell.X-1, node.Cell.Y).Walkable
			right := m.Get(node.Cell.X+1, node.Cell.Y).Walkable

			if node.Cell.X > 0 && node.Cell.Y > 0 {
				c := m.Get(node.Cell.X-1, node.Cell.Y-1)
				n := &Node{c, node, c.Cost + node.Cost + diagonalCost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (!wallsBlockDiagonals || (left && up)) {
					openNodes = append(openNodes, n)
					checkedNodes[n.Cell] = struct{}{}
					// checkedNodes = append(checkedNodes, n.Cell)
				}
			}

			if node.Cell.X < m.Width()-1 && node.Cell.Y > 0 {
				c := m.Get(node.Cell.X+1, node.Cell.Y-1)
				n := &Node{c, node, c.Cost + node.Cost + diagonalCost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (!wallsBlockDiagonals || (right && up)) {
					openNodes = append(openNodes, n)
					checkedNodes[n.Cell] = struct{}{}
					// checkedNodes = append(checkedNodes, n.Cell)
				}
			}

			if node.Cell.X > 0 && node.Cell.Y < m.Height()-1 {
				c := m.Get(node.Cell.X-1, node.Cell.Y+1)
				n := &Node{c, node, c.Cost + node.Cost + diagonalCost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (!wallsBlockDiagonals || (left && down)) {
					openNodes = append(openNodes, n)
					checkedNodes[n.Cell] = struct{}{}
					// checkedNodes = append(checkedNodes, n.Cell)
				}
			}

			if node.Cell.X < m.Width()-1 && node.Cell.Y < m.Height()-1 {
				c := m.Get(node.Cell.X+1, node.Cell.Y+1)
				n := &Node{c, node, c.Cost + node.Cost + diagonalCost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (!wallsBlockDiagonals || (right && down)) {
					openNodes = append(openNodes, n)
					checkedNodes[n.Cell] = struct{}{}
					// checkedNodes = append(checkedNodes, n.Cell)
				}
			}

		}

		// We sort the list of nodes by the cost to make the ones with lower cost checked first. That means that the function
		// automatically favors paths that are shorter (and so the "top" Cell has the shortest Cost), or Paths that cross over
		// the lowest-cost Cells (and so the constructed Path might be longer, but have a lower overall Cost).
		sort.Slice(openNodes, func(i, j int) bool {
			return openNodes[i].Cost < openNodes[j].Cost
		})

	}

	return path

}

// GetPath returns a Path, from the starting world X and Y position to the ending X and Y position. diagonals controls whether
// moving diagonally is acceptable when creating the Path. wallsBlockDiagonals indicates whether to allow diagonal movement "through" walls
// that are positioned diagonally. This is essentially just a smoother way to get a Path from GetPathFromCells().
func (m *Grid) GetPath(startX, startY, endX, endY float64, diagonals bool, wallsBlockDiagonals bool) *Path {

	sx, sy := m.WorldToGrid(startX, startY)
	sc := m.Get(sx, sy)
	ex, ey := m.WorldToGrid(endX, endY)
	ec := m.Get(ex, ey)

	if sc != nil && ec != nil {
		return m.GetPathFromCells(sc, ec, diagonals, wallsBlockDiagonals)
	}
	return nil
}

// DataAsStringArray returns a 2D array of runes for each Cell in the Grid. The first axis is the Y axis.
func (m *Grid) DataAsStringArray() []string {

	data := []string{}

	for y := 0; y < m.Height(); y++ {
		data = append(data, "")
		for x := 0; x < m.Width(); x++ {
			data[y] += string(m.Data[y][x].Rune)
		}
	}

	return data

}

// DataAsRuneArrays returns a 2D array of runes for each Cell in the Grid. The first axis is the Y axis.
func (m *Grid) DataAsRuneArrays() [][]rune {

	runes := [][]rune{}

	for y := 0; y < m.Height(); y++ {
		runes = append(runes, []rune{})
		for x := 0; x < m.Width(); x++ {
			runes[y] = append(runes[y], m.Data[y][x].Rune)
		}
	}

	return runes

}

// A Path is a struct that represents a path, or sequence of Cells from point A to point B. The Cells list is the list of Cells contained in the Path,
// and the CurrentIndex value represents the current step on the Path. Using Path.Next() and Path.Prev() advances and walks back the Path by one step.
type Path struct {
	Cells        []*Cell
	CurrentIndex int
}

// TotalCost returns the total cost of the Path (i.e. is the sum of all of the Cells in the Path).
func (p *Path) TotalCost() float64 {

	cost := 0.0
	for _, cell := range p.Cells {
		cost += cell.Cost
	}
	return cost

}

// Reverse reverses the Cells in the Path.
func (p *Path) Reverse() {

	np := []*Cell{}

	for c := len(p.Cells) - 1; c >= 0; c-- {
		np = append(np, p.Cells[c])
	}

	p.Cells = np

}

// Restart restarts the Path, so that calling path.Current() will now return the first Cell in the Path.
func (p *Path) Restart() {
	p.CurrentIndex = 0
}

// Current returns the current Cell in the Path.
func (p *Path) Current() *Cell {
	return p.Cells[p.CurrentIndex]

}

// Next returns the next cell in the path. If the Path is at the end, Next() returns nil.
func (p *Path) Next() *Cell {

	if p.CurrentIndex < len(p.Cells)-1 {
		return p.Cells[p.CurrentIndex+1]
	}
	return nil

}

// Advance advances the path by one cell.
func (p *Path) Advance() {

	p.CurrentIndex++
	if p.CurrentIndex >= len(p.Cells) {
		p.CurrentIndex = len(p.Cells) - 1
	}

}

// Prev returns the previous cell in the path. If the Path is at the start, Prev() returns nil.
func (p *Path) Prev() *Cell {

	if p.CurrentIndex > 0 {
		return p.Cells[p.CurrentIndex-1]
	}
	return nil

}

// Same returns if the Path shares the exact same cells as the other specified Path.
func (p *Path) Same(otherPath *Path) bool {

	if p == nil || otherPath == nil || len(p.Cells) != len(otherPath.Cells) {
		return false
	}

	for i := range p.Cells {
		if len(otherPath.Cells) <= i || p.Cells[i] != otherPath.Cells[i] {
			return false
		}
	}

	return true

}

// Length returns the length of the Path (how many Cells are in the Path).
func (p *Path) Length() int {
	return len(p.Cells)
}

// Get returns the Cell of the specified index in the Path. If the index is outside of the
// length of the Path, it returns -1.
func (p *Path) Get(index int) *Cell {
	if index < len(p.Cells) {
		return p.Cells[index]
	}
	return nil
}

// Index returns the index of the specified Cell in the Path. If the Cell isn't contained
// in the Path, it returns -1.
func (p *Path) Index(cell *Cell) int {
	for i, c := range p.Cells {
		if c == cell {
			return i
		}
	}
	return -1
}

// SetIndex sets the index of the Path, allowing you to safely manually manipulate the Path
// as necessary. If the index exceeds the bounds of the Path, it will be clamped.
func (p *Path) SetIndex(index int) {

	if index >= len(p.Cells) {
		p.CurrentIndex = len(p.Cells) - 1
	} else if index < 0 {
		p.CurrentIndex = 0
	} else {
		p.CurrentIndex = index
	}

}

// AtStart returns if the Path's current index is 0, the first Cell in the Path.
func (p *Path) AtStart() bool {
	return p.CurrentIndex == 0
}

// AtEnd returns if the Path's current index is the last Cell in the Path.
func (p *Path) AtEnd() bool {
	return p.CurrentIndex >= len(p.Cells)-1
}
