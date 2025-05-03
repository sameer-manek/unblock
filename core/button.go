package core

import raylib "github.com/gen2brain/raylib-go/raylib"

type Button struct {
	Position raylib.Vector2
	Size     raylib.Vector2
	Text     string
	FontSize uint8
	Scale    float32
	Action   func()
}

func NewButton(pos, size raylib.Vector2, txt string, font_size uint8, action func()) *Button {
	return &Button{
		Position: pos,
		Size:     size,
		Text:     txt,
		FontSize: font_size,
		Action:   action,
		Scale:    1.0,
	}
}

func (btn *Button) Update() {
	if raylib.IsMouseButtonDown(raylib.MouseButtonLeft) && btn.IsPointColliding(raylib.GetMousePosition()) {
		btn.Scale = 0.9
	}
	if raylib.IsMouseButtonUp(raylib.MouseButtonLeft) && btn.IsPointColliding(raylib.GetMousePosition()) {
		btn.Scale = 1.0
	}
	if raylib.IsMouseButtonPressed(raylib.MouseButtonLeft) && btn.IsPointColliding(raylib.GetMousePosition()) {
		btn.Action()
	}
}

func (btn *Button) Draw() {
	// draw button
	button_rect := raylib.NewRectangle(btn.Position.X, btn.Position.Y, btn.Size.X*btn.Scale, btn.Size.Y*btn.Scale)
	raylib.DrawRectangleRounded(button_rect, 0.25, 0, raylib.DarkBlue)
	x := int32(btn.Position.X) + int32(btn.Size.X*0.5) - raylib.MeasureText(btn.Text, int32(btn.FontSize))/2
	raylib.DrawText(
		btn.Text,
		x,
		int32(btn.Position.Y)+(int32(btn.Size.Y*0.5)-int32(btn.FontSize/2)),
		int32(float32(btn.FontSize)*btn.Scale),
		raylib.White)
}

func (btn *Button) GetBounds() *Bounds {
	return &Bounds{
		Start: btn.Position,
		End:   raylib.Vector2Add(btn.Position, btn.Size),
	}
}

func (btn *Button) IsPointColliding(point raylib.Vector2) bool {
	bounds := btn.GetBounds()

	return point.X > bounds.Start.X &&
		point.X < bounds.End.X &&
		point.Y > bounds.Start.Y &&
		point.Y < bounds.End.Y
}
