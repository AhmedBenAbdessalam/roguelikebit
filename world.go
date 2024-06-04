package main

import (
	"fmt"
	"log"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var position *ecs.Component
var renderable *ecs.Component
var monster *ecs.Component
var health *ecs.Component
var meleeWeapon *ecs.Component
var armor *ecs.Component
var name *ecs.Component

func InitializeWorld(startingLevel Level) (*ecs.Manager, map[string]ecs.Tag) {
	//Get First Room
	startingRoom := startingLevel.Rooms[0]
	x, y := startingRoom.Center()
	tags := make(map[string]ecs.Tag)
	manager := ecs.NewManager()

	player := manager.NewComponent()
	monster = manager.NewComponent()
	position = manager.NewComponent()
	renderable = manager.NewComponent()
	movable := manager.NewComponent()
	health = manager.NewComponent()
	meleeWeapon = manager.NewComponent()
	armor = manager.NewComponent()
	name = manager.NewComponent()

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	manager.NewEntity().
		AddComponent(player, Player{}).
		AddComponent(renderable, &Renderable{
			Image: playerImg,
		}).
		AddComponent(movable, Movable{}).
		AddComponent(position, &Position{
			X: x,
			Y: y,
		}).
		AddComponent(health, &Health{
			MaxHealth:     30,
			CurrentHealth: 30,
		}).
		AddComponent(meleeWeapon, &MeleeWeapon{
			Name:          "Fake Excalibur",
			MinimumDamage: 10,
			MaximumDamage: 30,
			ToHitBonus:    2,
		}).
		AddComponent(armor, &Armor{
			Name:       "old hood",
			Defense:    1,
			ArmorClass: 1,
		}).
		AddComponent(name, &Name{
			Label: "Player",
		})
	for i, room := range startingLevel.Rooms {
		if room.X1 != startingRoom.X1 {
			mX, mY := room.Center()
			manager.NewEntity().
				AddComponent(monster, Monster{}).
				AddComponent(renderable, &Renderable{
					Image: skeletonImg,
				}).
				AddComponent(position, &Position{
					X: mX,
					Y: mY,
				}).
				AddComponent(name, &Name{
					Label: fmt.Sprintf("Skeleton #%d", i),
				}).
				AddComponent(health, &Health{
					MaxHealth:     10,
					CurrentHealth: 10,
				}).
				AddComponent(meleeWeapon, &MeleeWeapon{
					Name:          "Short Sword",
					MinimumDamage: 2,
					MaximumDamage: 6,
					ToHitBonus:    0,
				}).
				AddComponent(armor, &Armor{
					Name:       "Bone",
					Defense:    3,
					ArmorClass: 4,
				})
		}
	}

	players := ecs.BuildTag(player, position, health, meleeWeapon, armor, name)
	tags["players"] = players
	renderables := ecs.BuildTag(renderable, position)
	tags["renderables"] = renderables
	monsters := ecs.BuildTag(monster, position, health, meleeWeapon, armor, name)
	tags["monsters"] = monsters
	return manager, tags
}
