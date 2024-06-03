package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/norendren/go-fov/fov"
)

type Level struct {
	Tiles         []MapTile
	Rooms         []Rect
	PlayerVisible *fov.View
}

func NewLevel() Level {
	l := Level{}
	rooms := make([]Rect, 0)
	l.Rooms = rooms
	l.GenerateLevelTiles()
	l.PlayerVisible = fov.New()
	return l
}

type MapTile struct {
	PixelX     int
	PixelY     int
	Blocked    bool
	Image      *ebiten.Image
	IsRevealed bool
}

func (level *Level) GetIndexFromXY(x, y int) int {
	gd := NewGameData()
	return (y * gd.ScreenWidth) + x
}

func (level *Level) CreateTiles() []MapTile {
	gd := NewGameData()
	tiles := make([]MapTile, gd.ScreenWidth*gd.ScreenHeight)

	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < gd.ScreenHeight; y++ {
			index := level.GetIndexFromXY(x, y)
			wall, _, err := ebitenutil.NewImageFromFile("assets/wall.png")
			if err != nil {
				log.Fatal(err)
			}
			tile := MapTile{
				PixelX:     x * gd.TileWidth,
				PixelY:     y * gd.TileHeight,
				Blocked:    true,
				Image:      wall,
				IsRevealed: false,
			}
			tiles[index] = tile
		}
	}
	return tiles
}

func (level *Level) DrawLevel(screen *ebiten.Image) {
	gd := NewGameData()
	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < gd.ScreenHeight; y++ {
			isVisible := level.PlayerVisible.IsVisible(x, y)
			tile := level.Tiles[level.GetIndexFromXY(x, y)]
			if isVisible {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
				screen.DrawImage(tile.Image, op)
				level.Tiles[level.GetIndexFromXY(x, y)].IsRevealed = true
			} else if tile.IsRevealed {
				op := &colorm.DrawImageOptions{}
				op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
				var cm colorm.ColorM
				cm.Scale(1, 1, 1, .5)
				colorm.DrawImage(screen, tile.Image, cm, op)
			}
		}
	}
}

func (level *Level) createRoom(room Rect) {
	floor, _, err := ebitenutil.NewImageFromFile("assets/floor.png")
	if err != nil {
		log.Fatal(err)
	}
	for y := room.Y1; y < room.Y2+1; y++ {
		for x := room.X1; x < room.X2+1; x++ {
			index := level.GetIndexFromXY(x, y)
			level.Tiles[index].Image = floor
			level.Tiles[index].Blocked = false
		}
	}
}

// GenerateLevelTiles creates a new Dungeon Level Map
func (level *Level) GenerateLevelTiles() {
	MIN_SIZE := 6
	MAX_SIZE := 10
	MAX_ROOMS := 30

	gd := NewGameData()
	tiles := level.CreateTiles()
	level.Tiles = tiles

	for idx := 0; idx < MAX_ROOMS; idx++ {
		w := GetRandomBetween(MIN_SIZE, MAX_SIZE)
		h := GetRandomBetween(MIN_SIZE, MAX_SIZE)
		x := GetDiceRoll(gd.ScreenWidth - w - 2)
		y := GetDiceRoll(gd.ScreenHeight - h - 2)

		new_room := NewRect(x, y, w, h)
		okToAdd := true
		for _, otherRoom := range level.Rooms {
			if new_room.Intersect(otherRoom) {
				okToAdd = false
				break
			}
		}
		if okToAdd {
			level.createRoom(new_room)
			// create tunnels
			if len(level.Rooms) > 0 {
				previousRoom := level.Rooms[len(level.Rooms)-1]
				x1, y1 := previousRoom.Center()
				x2, y2 := new_room.Center()
				coinFlip := GetDiceRoll(2)
				if coinFlip == 1 {
					level.createHorizontalTunnel(x1, x2, y1)
					level.createVerticalTunnel(y1, y2, x2)
				} else {
					level.createVerticalTunnel(y1, y2, x1)
					level.createHorizontalTunnel(x1, x2, y2)
				}
			}
			level.Rooms = append(level.Rooms, new_room)
		}
	}
}

func (level *Level) createHorizontalTunnel(x1, x2, y int) {
	gd := NewGameData()
	floor, _, err := ebitenutil.NewImageFromFile("assets/floor.png")
	if err != nil {
		log.Fatal(err)
	}
	for x := min(x1, x2); x < max(x1, x2)+1; x++ {
		index := level.GetIndexFromXY(x, y)
		if index > 0 && index < gd.ScreenWidth*gd.ScreenHeight {
			level.Tiles[index].Blocked = false
			level.Tiles[index].Image = floor
		}
	}
}

func (level *Level) createVerticalTunnel(y1, y2, x int) {
	gd := NewGameData()
	floor, _, err := ebitenutil.NewImageFromFile("assets/floor.png")
	if err != nil {
		log.Fatal(err)
	}
	for y := min(y1, y2); y < max(y1, y2)+1; y++ {
		index := level.GetIndexFromXY(x, y)
		if index > 0 && index < gd.ScreenWidth*gd.ScreenHeight {
			level.Tiles[index].Blocked = false
			level.Tiles[index].Image = floor
		}
	}
}

func (level Level) InBounds(x, y int) bool {
	gd := NewGameData()
	if x < 0 || x > gd.ScreenWidth || y < 0 || y > gd.ScreenHeight {
		return false
	}
	return true
}

// TODO: Change this to check for WALL, not blocked
func (level Level) IsOpaque(x, y int) bool {
	idx := level.GetIndexFromXY(x, y)
	return level.Tiles[idx].Blocked
}
