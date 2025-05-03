package game

import (
	"math"
	"unblock/core"

	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/sameer-manek/rush"
)

var BoardOrigin = raylib.NewVector2(SCREEN_WIDTH*0.5-BOARD_SIZE*0.5, SCREEN_HEIGHT*0.45-BOARD_SIZE*0.5)
var BoardBounds = core.GetBounds(BoardOrigin, raylib.NewVector2(BOARD_SIZE, BOARD_SIZE))
var PuzzleSize = 6
var CellSize = BOARD_SIZE / PuzzleSize
var BlockPadding = float64(CellSize) * 0.05

var PuzzleGenerator = &rush.Generator{
	Width:       PuzzleSize,
	Height:      PuzzleSize,
	PrimaryRow:  2,
	PrimarySize: 2,
}
var PuzzleString = PuzzleGenerator.Generate(3000).Hash()

var PuzzlesWon []string

func DrawBoard() {
	for i := range PuzzleSize {
		for j := range PuzzleSize {
			if i == PuzzleGenerator.PrimaryRow && j > (PuzzleSize-1-PuzzleGenerator.PrimarySize) {
				color := raylib.Color{84, 255, 84, 50}
				if GameState == "WON" {
					color.A = 100
				}
				raylib.DrawRectangle(
					int32(BoardOrigin.X+float32(j*CellSize)),
					int32(BoardOrigin.Y+float32(i*CellSize)),
					int32(CellSize), int32(CellSize), color)
			}
			raylib.DrawRectangleLines(
				int32(BoardOrigin.X+float32(j*CellSize)),
				int32(BoardOrigin.Y+float32(i*CellSize)),
				int32(CellSize), int32(CellSize), raylib.White)
		}
	}
}

func ReloadBlocks() {
	Blocks = nil
	LoadBlocks()
}

func LoadBlocks() {
	runes := []rune(PuzzleString)
	for i, c := range runes {
		if c == 'o' || c == '.' {
			continue
		}
		row := i / PuzzleSize
		col := int(math.Mod(float64(i), float64(PuzzleSize)))

		if c == 'x' {
			Blocks = append(Blocks, NewBlock(
				col, row, 1, 1,
				'N', raylib.White, false))
		} else {
			var color raylib.Color
			var axis rune
			var sizeX = 1
			var sizeY = 1

			if c == 'A' {
				color = raylib.Red
			} else {
				color = raylib.Yellow
			}

			// check horizontal
			p2 := PuzzleSize * PuzzleSize
			for x := i + 1; x < p2; x++ {
				if runes[x] != c {
					break
				}
				sizeX++
				runes[x] = 'o'
			}

			// check vertical
			for x := i + PuzzleSize; x < p2; x += PuzzleSize {

				if runes[x] != c {
					break
				}
				sizeY++
				runes[x] = 'o'
			}

			if sizeX > 1 {
				axis = 'X'
			} else if sizeY > 1 {
				axis = 'Y'
			}

			Blocks = append(Blocks, NewBlock(
				col, row, sizeX, sizeY,
				axis, color, c == 'A'))
			// check vertical
		}

	}
}
