package main

import (
	"fmt"
	"unblock/core"
	"unblock/game"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	raylib.InitWindow(game.SCREEN_WIDTH, game.SCREEN_HEIGHT, "Test game")
	raylib.SetTargetFPS(60)

	game.LoadBlocks()

	var loading_img = raylib.LoadTexture("resources/assets/loading.png")

	reloadBtn := core.NewButton(
		raylib.NewVector2(game.SCREEN_WIDTH*0.5-225, game.SCREEN_HEIGHT*0.75),
		raylib.NewVector2(200, 50),
		"Reload", 24,
		func() {
			game.Blocks = nil
			if game.GameState == "WON" {
				game.GameState = "LOADING"
				game.PuzzleString = game.PuzzleGenerator.Generate(3000).Hash()
			}
			game.ClearMoveStack()
			game.LoadBlocks()

		})
	hintBtn := core.NewButton(
		raylib.NewVector2(game.SCREEN_WIDTH*0.5+25, game.SCREEN_HEIGHT*0.75+75),
		raylib.NewVector2(200, 50),
		"Hint | 2", 24,
		func() { fmt.Println("Hinting..") })
	undoBtn := core.NewButton(
		raylib.NewVector2(game.SCREEN_WIDTH*0.5-225, game.SCREEN_HEIGHT*0.75+75),
		raylib.NewVector2(200, 50),
		"Undo", 24,
		game.PopMove)
	quitBtn := core.NewButton(
		raylib.NewVector2(game.SCREEN_WIDTH*0.5+25, game.SCREEN_HEIGHT*0.75),
		raylib.NewVector2(200, 50),
		"Quit", 24,
		raylib.CloseWindow)

	for !raylib.WindowShouldClose() {
		reloadBtn.Update()
		hintBtn.Update()
		undoBtn.Update()
		quitBtn.Update()
		if game.GameState == "WON" {
			reloadBtn.Text = "NEW"
		} else {
			reloadBtn.Text = "Reload"
		}
		if game.GameState == "LOADING" {
			raylib.DrawTexture(loading_img, game.SCREEN_WIDTH/2-32, game.SCREEN_HEIGHT/2-32, raylib.White)
		} else {
			for _, blk := range game.Blocks {
				blk.Update()
			}

			raylib.BeginDrawing()
			raylib.DrawFPS(0, 0)
			raylib.ClearBackground(raylib.Color{40, 40, 40, 255})

			game.DrawBoard()
			for _, blk := range game.Blocks {
				blk.Draw()
			}

			reloadBtn.Draw()
			quitBtn.Draw()
			if game.GameState != "WON" {
				hintBtn.Draw()
				undoBtn.Draw()
			}

			game.DrawUI()
			raylib.EndDrawing()
		}
	}

	raylib.CloseWindow()
}
