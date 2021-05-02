
# paths

![paths](https://user-images.githubusercontent.com/4733521/48970683-21882880-efc4-11e8-9b60-670f46c6fd77.gif)

[GoDocs](https://pkg.go.dev/github.com/SolarLune/paths?tab=doc)

## What is paths?

paths is a pathfinding library written in Golang created mainly for video games. Its main feature is simple best-first and shortest-cost path finding.

## Why is it called that?

Because I finally learned how to spell.

## Why did you create paths?

Because I needed to do pathfinding for a game and couldn't really find any pre-made Golang libraries that seemed to do this, so... Yet again, here we are.

## How do I install it?

Just go get it:

`go get github.com/SolarLune/paths`

## How do I use it?

paths is based around defining a Grid, which consists of a rectangular series of Cells. Each Cell occupies a single X and Y position in space, and has a couple of properties that influence pathfinding, which are Cost and Walkability. If a Cell isn't walkable, then it is considered an obstacle that paths generated must circumvent. All Cells default to a Cost of 1; pathfinding will prioritize lower-cost Cells.

```go
import "github.com/SolarLune/paths"

func Init() {

    // This line creates a new Grid, comprised (internally) of Cells. The size is 10x10. Each Cell's 
    // "world" size is 16x16. By default, all Cells are walkable, have a cost of 1, and a blank rune 
    // of ' ' associated with them.
    firstMap = paths.NewGrid(10, 10, 16, 16)
    
    // You can also create the Grid from an array of strings (which are interpreted as arrays of runes), 
    // if you already have it:
    layout := []string{
        "xxxxxxxxxx",
        "x        x",
        "x xxxxxx x",
        "x xg   x x",
        "x xgxx x x",
        "x gggx x x",
        "x xxxx   x",
        "x  xgg x x",
        "xg ggx x x",
        "xxxxxxxxxx",
    }

    secondMap := paths.NewGridFromStringArrays(layout, 16, 16)

    // After creating the Grid, you can edit it using the Grid's functions. Note that here, we're using 'x' 
    // to get Cells that have the rune for the lowercase x character 'x', not the string "x".
    firstMap.SetWalkable('x', false)

    // You can also loop through them by using the `GetCells` functions thusly...
    for _, goop := range secondMap.GetCellsByRune('g') {
        goop.Cost = 5
    }

    // This gets a new Path from the Cell occupied by a starting position [24, 21], to another [99, 78]. The last boolean argument,
    // false, indicates whether moving diagonally is fine when creating the Path.
    firstPath := GameMap.GetPath(24, 21, 99, 78, false)

    // You can also get a path using references to the Cells directly.
    secondPath := GameMap.GetPathFromCell(GameMap.Get(1, 1), GameMap.Get(6, 3), false)

    // After that, you can use Path.Current() and Path.Next() to get the current and next Cells on the Path. When you determine that 
    // the pathfinding agent has reached that Cell, you can kick the Path forward with path.Advance().

    // And that's it!

}

```
---

And that's about it! If you want to see more info or examples, feel free to examine the main.go and world.go tests to see how the test is set up.

[You can check out the GoDoc link here, as well.](https://pkg.go.dev/github.com/SolarLune/paths?tab=doc)

You can also run the example by installing SDL with the instructions [here](https://github.com/veandco/go-sdl2#requirements)
and then:

```
$ cd ./example
$ go run ./
```

## Dependencies?

For the actual package, there are no external dependencies.

For the tests, paths requires veandco's sdl2 port to create the window, handle input, and draw the shapes and text.

## Shout-out Time!

Props to whoever made arcadepi.ttf! It's a nice font.

Thanks a lot to the SDL2 team for development.

Thanks to veandco for maintaining the Golang SDL2 port, as well!
