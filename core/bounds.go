package core

import rl "github.com/gen2brain/raylib-go/raylib"

type Bounds struct {
	Start rl.Vector2
	End   rl.Vector2
}

func GetBounds(pos rl.Vector2, size rl.Vector2) *Bounds {
	end := rl.Vector2Add(pos, size)
	return &Bounds{
		Start: pos,
		End:   end,
	}
}

func (b *Bounds) checkPointInBounds(pos *rl.Vector2) bool {
	return b.Start.X <= pos.X &&
		b.End.X >= pos.X &&
		b.Start.Y <= pos.Y &&
		b.End.Y >= pos.Y
}

func (b *Bounds) checkBoxCollision(boxBounds *Bounds) bool {
	return (b.Start.X <= boxBounds.End.X || b.End.X >= boxBounds.Start.X) &&
		(b.Start.Y <= boxBounds.End.Y || b.End.Y >= boxBounds.Start.Y)
}

func (b *Bounds) checkBoxInBounds(boxBounds *Bounds) bool {
	return (b.Start.X <= boxBounds.Start.X || b.End.X >= boxBounds.End.X) &&
		(b.Start.Y <= boxBounds.Start.Y || b.End.Y >= boxBounds.End.Y)
}
