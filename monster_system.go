package main

import (
	"github.com/norendren/go-fov/fov"
)

func UpdateMonster(game *Game) {
	l := game.Map.CurrentLevel
	playerPosition := Position{}

	for _, p := range game.World.Query(game.WorldTags["players"]) {
		pos := p.Components[position].(*Position)
		playerPosition = *pos
	}

	for _, mon := range game.World.Query(game.WorldTags["monsters"]) {
		pos := mon.Components[position].(*Position)
		monsterVision := fov.New()
		monsterVision.Compute(l, pos.X, pos.Y, 8)
		if monsterVision.IsVisible(playerPosition.X, playerPosition.Y) {
			if pos.GetManhattanDistance(&playerPosition) == 1 {
				//The monster is right next to the player. Just smack him down
				AttackSystem(game, pos, &playerPosition)
				if mon.Components[health].(*Health).CurrentHealth <= 0 {
					//this monster is dead
					//clear the tile
					t := l.Tiles[l.GetIndexFromXY(pos.X, pos.Y)]
					t.Blocked = false
				}
			} else {
				astar := AStar{}
				path := astar.GetPath(l, pos, &playerPosition)
				if len(path) > 1 {
					nextTile := l.Tiles[l.GetIndexFromXY(path[1].X, path[1].Y)]
					if !nextTile.Blocked {
						l.Tiles[l.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
						pos.X = path[1].X
						pos.Y = path[1].Y
						nextTile.Blocked = true
					}
				}
			}
		}
	}
	game.Turn = PlayerTurn

}
