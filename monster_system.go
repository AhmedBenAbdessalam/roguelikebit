package main

import (
	"log"

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
		m := mon.Components[monster].(Monster)
		pos := mon.Components[position].(*Position)
		monsterVision := fov.New()
		monsterVision.Compute(l, pos.X, pos.Y, 8)
		if monsterVision.IsVisible(playerPosition.X, playerPosition.Y) {
			log.Printf("%s is scared to the bone\n", m.Name)
		}
	}
	game.Turn = PlayerTurn

}
