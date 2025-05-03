package game

import (
	"strconv"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var GameState = "IDLE"
var Moves = 0

func DrawUI() {
	var color raylib.Color
	if GameState == "WON" {
		color = raylib.Green
	} else {
		color = raylib.White
	}
	// Draw moves
	moves_message := "MOVES: " + strconv.Itoa(Moves)
	x := SCREEN_WIDTH/2 - raylib.MeasureText(moves_message, 32)/2
	raylib.DrawText(moves_message, x, SCREEN_HEIGHT*0.1, 32, color)

	// Draw state
	x = SCREEN_WIDTH/2 - raylib.MeasureText(GameState, 18)/2
	raylib.DrawText(GameState, x, SCREEN_HEIGHT*0.15, 18, color)

	// Draw button

}
