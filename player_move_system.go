package main

import "github.com/hajimehoshi/ebiten/v2"

func TryMovePlayer(g *Game) {
	turnTaken := false
	//get direction delta
	x, y := 0, 0
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		turnTaken = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		y = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		y = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		x = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		x = 1
	}
	level := g.Map.CurrentLevel
	for _, result := range g.World.Query(g.WorldTags["players"]) {
		pos := result.Components[position].(*Position)

		index := level.GetIndexFromXY(pos.X+x, pos.Y+y)
		//check if tile is blocked
		if !level.Tiles[index].Blocked {
			level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
			pos.X += x
			pos.Y += y
			level.Tiles[index].Blocked = true
			level.PlayerVisible.Compute(level, pos.X, pos.Y, 8)
		}
		if x != 0 || y != 0 || turnTaken {
			g.Turn = GetNextState(g.Turn)
			g.TurnCounter = 0

		}
	}
}
