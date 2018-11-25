/*
Package paths is a simple library written in Go made to handle 2D pathfinding for games. All you need to do is generate a Grid,
specify which cells aren't walkable, optionally change the cost on specific cells, and finally get a path from one cell to
another.
*/
package paths

import (
	"fmt"
	"sort"
)

// A Cell represents a point on a Grid map. It has an X and Y value for the position, a Cost, which influences which Cells are
// ideal for paths, Walkable, which indicates if the tile can be walked on or should be avoided, and a Character, which indicates
// which rune character the Cell is represented by.
type Cell struct {
	X, Y      int
	Cost      int
	Walkable  bool
	Character rune
}

func (cell Cell) String() string {
	return fmt.Sprintf("X:%d Y:%d Cost:%d Walkable:%t Character:%s", cell.X, cell.Y, cell.Cost, cell.Walkable, cell.Character)
}

// Grid represents a "map" composed of individual Cells at each point in the map.
type Grid struct {
	Data [][]*Cell
}

// NewGrid returns a new Grid of width x height size.
func NewGrid(width, height int) *Grid {
	m := &Grid{}
	for y := 0; y < height; y++ {
		m.Data = append(m.Data, []*Cell{})
		for x := 0; x < width; x++ {
			m.Data[y] = append(m.Data[y], &Cell{x, y, 1, true, ' '})
		}
	}
	return m
}

// DataToString returns a string, used to easily identify the Grid map.
func (m *Grid) DataToString() string {
	s := ""
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			s += string(m.Data[y][x].Character) + " "
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

// GetCellsByCharacter returns a slice of pointers to Cells that all have the character provided.
func (m *Grid) GetCellsByCharacter(char rune) []*Cell {

	cells := make([]*Cell, 0)

	for y := 0; y < m.Height(); y++ {

		for x := 0; x < m.Width(); x++ {

			c := m.Get(x, y)
			if c.Character == char {
				cells = append(cells, c)
			}

		}

	}

	return cells

}

// GetCellsByCost returns a slice of pointers to Cells that all have the Cost value provided.
func (m *Grid) GetCellsByCost(cost int) []*Cell {

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

// GetPath returns a Path, from the starting Cell to the destination Cell. diagonals controls whether diagonal Cells are
// considered when creating the Path.
func (m *Grid) GetPath(start, dest *Cell, diagonals bool) Path {

	type Node struct {
		Cell   *Cell
		Parent *Node
		Cost   int
	}

	openNodes := []*Node{&Node{Cell: dest, Cost: dest.Cost}}

	checkedNodes := make([]*Cell, 0)

	hasBeenAdded := func(cell *Cell) bool {

		for _, c := range checkedNodes {
			if cell == c {
				return true
			}
		}
		return false

	}

	path := Path{}

	if !dest.Walkable {
		return path
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
			for true {
				path = append(path, t.Cell)
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
				checkedNodes = append(checkedNodes, n.Cell)
			}
		}
		if node.Cell.X < m.Width()-1 {
			c := m.Get(node.Cell.X+1, node.Cell.Y)
			n := &Node{c, node, c.Cost + node.Cost}
			if n.Cell.Walkable && !hasBeenAdded(n.Cell) {
				openNodes = append(openNodes, n)
				checkedNodes = append(checkedNodes, n.Cell)
			}
		}

		if node.Cell.Y > 0 {
			c := m.Get(node.Cell.X, node.Cell.Y-1)
			n := &Node{c, node, c.Cost + node.Cost}
			if n.Cell.Walkable && !hasBeenAdded(n.Cell) {
				openNodes = append(openNodes, n)
				checkedNodes = append(checkedNodes, n.Cell)
			}
		}
		if node.Cell.Y < m.Height()-1 {
			c := m.Get(node.Cell.X, node.Cell.Y+1)
			n := &Node{c, node, c.Cost + node.Cost}
			if n.Cell.Walkable && !hasBeenAdded(n.Cell) {
				openNodes = append(openNodes, n)
				checkedNodes = append(checkedNodes, n.Cell)
			}
		}

		// Do the same thing for diagonals.
		if diagonals {

			if node.Cell.X > 0 && node.Cell.Y > 0 {
				c := m.Get(node.Cell.X-1, node.Cell.Y-1)
				n := &Node{c, node, c.Cost + node.Cost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) {
					openNodes = append(openNodes, n)
					checkedNodes = append(checkedNodes, n.Cell)
				}
			}
			if node.Cell.X < m.Width()-1 && node.Cell.Y > 0 {
				c := m.Get(node.Cell.X+1, node.Cell.Y-1)
				n := &Node{c, node, c.Cost + node.Cost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) {
					openNodes = append(openNodes, n)
					checkedNodes = append(checkedNodes, n.Cell)
				}
			}

			if node.Cell.X > 0 && node.Cell.Y < m.Height()-1 {
				c := m.Get(node.Cell.X-1, node.Cell.Y+1)
				n := &Node{c, node, c.Cost + node.Cost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) {
					openNodes = append(openNodes, n)
					checkedNodes = append(checkedNodes, n.Cell)
				}
			}
			if node.Cell.X < m.Width()-1 && node.Cell.Y < m.Height()-1 {
				c := m.Get(node.Cell.X+1, node.Cell.Y+1)
				n := &Node{c, node, c.Cost + node.Cost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) {
					openNodes = append(openNodes, n)
					checkedNodes = append(checkedNodes, n.Cell)
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

// DataAsRuneArrays returns a 2D array of runes for each Cell in the Grid. The first axis is the Y axis.
func (m *Grid) DataAsRuneArrays() [][]rune {

	runes := [][]rune{}

	for y := 0; y < m.Height(); y++ {
		runes = append(runes, []rune{})
		for x := 0; x < m.Width(); x++ {
			runes[y] = append(runes[y], m.Data[y][x].Character)
		}
	}

	return runes

}

// NewGridFromStringArrays creates a Grid map from a 1D array of strings. Each string becomes a row of Cells, each
// with one rune as its character.
func NewGridFromStringArrays(arrays []string) *Grid {

	m := &Grid{}

	for y := 0; y < len(arrays); y++ {
		m.Data = append(m.Data, []*Cell{})
		stringLine := []rune(arrays[y])
		for x := 0; x < len(arrays[y]); x++ {
			m.Data[y] = append(m.Data[y], &Cell{X: x, Y: y, Cost: 1, Walkable: true, Character: stringLine[x]})
		}
	}

	return m

}

// NewGridFromRuneArrays creates a Grid map from a 2D array of runes. Each individual Rune becomes a Cell in the resulting
// Grid.
func NewGridFromRuneArrays(arrays [][]rune) *Grid {

	m := &Grid{}

	for y := 0; y < len(arrays); y++ {
		m.Data = append(m.Data, []*Cell{})
		for x := 0; x < len(arrays[y]); x++ {
			m.Data[y] = append(m.Data[y], &Cell{X: x, Y: y, Cost: 1, Walkable: true, Character: arrays[y][x]})
		}
	}

	return m

}

// A Path is a slice of pointers to Cells.
type Path []*Cell

// TotalCost returns the total cost of the Path.
func (p Path) TotalCost() int {

	cost := 0
	for _, cell := range p {
		cost += cell.Cost
	}
	return cost

}

// Reverse reverses the Path.
func (p *Path) Reverse() {

	np := Path{}

	for c := len(*p) - 1; c >= 0; c-- {
		np = append(np, (*p)[c])
	}

	*p = np

}
