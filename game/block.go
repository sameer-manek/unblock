package game

import (
	"math"
	"unblock/core"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Block struct {
	Position raylib.Vector2
	Size     raylib.Vector2
	Axis     rune
	Color    raylib.Color
	IsMain   bool
}

var Blocks []*Block

var forward_block *raylib.Vector2
var backward_block *raylib.Vector2

type Move struct {
	Block   *Block
	LastPos raylib.Vector2
}

var MoveStack []Move

func PushMove(mv Move) {
	MoveStack = append(MoveStack, mv)
}

func PopMove() {
	l := len(MoveStack)
	if l > 0 {
		mv := MoveStack[l-1]
		mv.Block.Position = mv.LastPos
		MoveStack = MoveStack[:l-1]
	}
}

func ClearMoveStack() {
	Moves = 0
	MoveStack = nil
	GameState = "IDLE"
}

var MovingBlock *Block
var TouchOffset *raylib.Vector2
var MovingBlockPos raylib.Vector2

func NewBlock(row, col, width, height int, axis rune, color raylib.Color, isMain bool) *Block {
	return &Block{
		Position: raylib.NewVector2(
			BoardOrigin.X+float32(row)*float32(CellSize)+float32(BlockPadding),
			BoardOrigin.Y+float32(col)*float32(CellSize)+float32(BlockPadding)),
		Size: raylib.NewVector2(
			float32(width)*float32(CellSize)-float32(BlockPadding)*2,
			float32(height)*float32(CellSize)-float32(BlockPadding)*2),
		Axis: axis, Color: color, IsMain: isMain,
	}
}

func (block *Block) Update() {
	if raylib.IsMouseButtonDown(raylib.MouseLeftButton) && GameState != "WON" {
		if block.IsPointColliding(raylib.GetMousePosition()) && MovingBlock == nil {
			GameState = "PLAYING"
			MovingBlock = block
			offset := raylib.Vector2Subtract(raylib.GetMousePosition(), block.Position)
			TouchOffset = &offset
			MovingBlockPos = MovingBlock.Position
		}

		if MovingBlock != nil {
			var newpos raylib.Vector2
			if MovingBlock.Axis == 'X' {
				newpos = raylib.NewVector2(raylib.GetMousePosition().X-TouchOffset.X, MovingBlock.Position.Y)
			}
			if MovingBlock.Axis == 'Y' {
				newpos = raylib.NewVector2(MovingBlock.Position.X, raylib.GetMousePosition().Y-TouchOffset.Y)
			}
			MovingBlock.MoveToPos(newpos)
		}
	}

	if raylib.IsMouseButtonUp(raylib.MouseButtonLeft) && MovingBlock != nil {
		MovingBlock.Settle()
		if MovingBlock.Position.X != MovingBlockPos.X || MovingBlock.Position.Y != MovingBlockPos.Y {
			Moves += 1
			mv := Move{
				Block:   MovingBlock,
				LastPos: MovingBlockPos,
			}
			PushMove(mv)
		}
		MovingBlock = nil
		TouchOffset = nil
		forward_block = nil
		backward_block = nil
	}
}

func (block *Block) Draw() {
	rect := raylib.NewRectangle(block.Position.X, block.Position.Y, block.Size.X, block.Size.Y)
	raylib.DrawRectangleRounded(rect, 0.25, 0, block.Color)
}

func (block *Block) GetBounds() *core.Bounds {
	return &core.Bounds{
		Start: block.Position,
		End:   raylib.Vector2Add(block.Position, block.Size),
	}
}

func (block *Block) MoveToPos(newpos raylib.Vector2) {
	if !block.IsBlocked(newpos) && block.IsPosOnBoard(&newpos) {
		block.Position = newpos
		if block.IsMain == true && GameState != "WON" &&
			block.Position.X+block.Size.X-BoardOrigin.X > float32(BOARD_SIZE)-(float32(CellSize)*0.5) {
			GameState = "WON"
		}
	} else {
		if block.Axis == 'X' {
			if block.Position.X < newpos.X && forward_block == nil {
				forward_block = &newpos
			} else if block.Position.X > newpos.X && backward_block == nil {
				backward_block = &newpos
			}
		} else if block.Axis == 'Y' {
			if block.Position.Y < newpos.Y && forward_block == nil {
				forward_block = &newpos
			} else if block.Position.Y > newpos.Y && backward_block == nil {
				backward_block = &newpos
			}
		}
	}
}

func (block *Block) IsBlocked(pos raylib.Vector2) bool {
	is_blocked := false
	if block.Axis == 'X' {
		if (block.Position.X < pos.X && forward_block != nil && pos.X >= forward_block.X) ||
			(block.Position.X < pos.X && backward_block != nil && pos.X <= backward_block.X) {
			is_blocked = true
		}
	} else if block.Axis == 'Y' {
		if (block.Position.Y < pos.Y && forward_block != nil && pos.Y >= forward_block.Y) ||
			(block.Position.Y < pos.Y && backward_block != nil && pos.Y <= backward_block.Y) {
			is_blocked = true
		}
	}
	newpos_offset := raylib.Vector2Subtract(pos, block.Position)
	for _, blk := range Blocks {
		if is_blocked {
			break
		}
		if block != blk {
			is_blocked = block.WillBoxCollide(blk.GetBounds(), newpos_offset)
		}
	}
	return is_blocked
}

func (block *Block) WillBoxCollide(box *core.Bounds, offset raylib.Vector2) bool {
	bounds := block.GetBounds()
	bounds.Start = raylib.Vector2Add(bounds.Start, offset)
	bounds.End = raylib.Vector2Add(bounds.End, offset)
	// Check if the two boxes overlap on the X-axis
	if bounds.Start.X < box.End.X && bounds.End.X > box.Start.X {
		// Check if the two boxes overlap on the Y-axis
		if bounds.Start.Y < box.End.Y && bounds.End.Y > box.Start.Y {
			return true
		}
	}
	return false
}

func (block *Block) IsPointColliding(point raylib.Vector2) bool {
	bounds := block.GetBounds()

	return point.X > bounds.Start.X &&
		point.X < bounds.End.X &&
		point.Y > bounds.Start.Y &&
		point.Y < bounds.End.Y
}

func (block *Block) IsPosOnBoard(pos *raylib.Vector2) bool {
	return pos.X > BoardBounds.Start.X &&
		pos.X+block.Size.X < BoardBounds.End.X &&
		pos.Y > BoardBounds.Start.Y &&
		pos.Y+block.Size.Y < BoardBounds.End.Y
}

func (block *Block) Settle() {
	if block.Axis == 'X' {
		var row uint8
		if float32(math.Mod(float64(block.Position.X-BoardOrigin.X), float64(float32(CellSize)))) < float32(CellSize)*0.5 {
			row = uint8(math.Floor(float64(block.Position.X-BoardOrigin.X) / float64(float32(CellSize))))
		} else {
			row = uint8(math.Ceil(float64(block.Position.X-BoardOrigin.X) / float64(float32(CellSize))))
		}
		block.Position.X = BoardOrigin.X + float32(row)*float32(CellSize) + float32(BlockPadding)
	}

	if block.Axis == 'Y' {
		var col uint8
		if float32(math.Mod(float64(block.Position.Y-BoardOrigin.Y), float64(float32(CellSize)))) < float32(CellSize)*0.5 {
			col = uint8(math.Floor(float64(block.Position.Y-BoardOrigin.Y) / float64(float32(CellSize))))
		} else {
			col = uint8(math.Ceil(float64(block.Position.Y-BoardOrigin.Y) / float64(float32(CellSize))))
		}
		block.Position.Y = BoardOrigin.Y + float32(col)*float32(CellSize) + float32(BlockPadding)
	}
}
