package main

import "github.com/hajimehoshi/ebiten/v2"

func processPlayerMovement(g *Game, level Level) {
	//get direction delta
	delta := Position{0, 0}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		delta.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		delta.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		delta.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		delta.X += 1
	}
	for _, result := range g.World.Query(g.WorldTags["players"]) {
		pos := result.Components[position].(*Position)

		newPos := Position{X: pos.X + delta.X, Y: pos.Y + delta.Y}
		index := level.GetIndexFromXY(newPos.X, newPos.Y)
		//check if out of bound
		if index < 0 || index > len(level.Tiles)-1 {
			return
		}
		//check if tile is blocked
		if !level.Tiles[index].Blocked {
			*pos = newPos
		}
	}
}
