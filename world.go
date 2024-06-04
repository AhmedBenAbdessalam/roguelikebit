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
var userMessage *ecs.Component

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
	userMessage = manager.NewComponent()

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	orcImg, _, err := ebitenutil.NewImageFromFile("assets/orc.png")
	if err != nil {
		log.Fatal(err)
	}
	// add player entity
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
			Name:          "Battle Axe",
			MinimumDamage: 10,
			MaximumDamage: 20,
			ToHitBonus:    3,
		}).
		AddComponent(armor, &Armor{
			Name:       "Plat Armor",
			Defense:    15,
			ArmorClass: 18,
		}).
		AddComponent(name, &Name{
			Label: "Player",
		}).
		AddComponent(userMessage, &UserMessage{
			AttackMessage:    "",
			DeadMessage:      "",
			GameStateMessage: "",
		})
	for i, room := range startingLevel.Rooms {
		if room.X1 != startingRoom.X1 {
			mX, mY := room.Center()
			mobSpawn := GetDiceRoll(2)
			// add orc
			if mobSpawn == 1 {
				manager.NewEntity().
					AddComponent(monster, Monster{}).
					AddComponent(renderable, &Renderable{
						Image: orcImg,
					}).
					AddComponent(position, &Position{
						X: mX,
						Y: mY,
					}).
					AddComponent(name, &Name{
						Label: fmt.Sprintf("Orc #%d", i),
					}).
					AddComponent(health, &Health{
						MaxHealth:     30,
						CurrentHealth: 30,
					}).
					AddComponent(meleeWeapon, &MeleeWeapon{
						Name:          "Machete",
						MinimumDamage: 4,
						MaximumDamage: 8,
						ToHitBonus:    1,
					}).
					AddComponent(armor, &Armor{
						Name:       "Leather",
						Defense:    5,
						ArmorClass: 6,
					}).
					AddComponent(userMessage, &UserMessage{
						AttackMessage:    "",
						DeadMessage:      "",
						GameStateMessage: "",
					})

			} else {
				// add skeleton
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
					}).
					AddComponent(userMessage, &UserMessage{
						AttackMessage:    "",
						DeadMessage:      "",
						GameStateMessage: "",
					})
			}
		}
	}

	players := ecs.BuildTag(player, position, health, meleeWeapon, armor, name, userMessage)
	tags["players"] = players
	renderables := ecs.BuildTag(renderable, position)
	tags["renderables"] = renderables
	monsters := ecs.BuildTag(monster, position, health, meleeWeapon, armor, name, userMessage)
	tags["monsters"] = monsters
	messengers := ecs.BuildTag(userMessage, position)
	tags["messengers"] = messengers
	return manager, tags
}
