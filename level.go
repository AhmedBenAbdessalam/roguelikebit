package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/norendren/go-fov/fov"
)

var levelHeight int = 0

type TileType int

const (
	WALL TileType = iota
	FLOOR
)

var floor *ebiten.Image
var wall *ebiten.Image

func loadTileImages() {
	if floor != nil && wall != nil {
		return
	}
	var err error
	floor, _, err = ebitenutil.NewImageFromFile("assets/floor.png")
	if err != nil {
		log.Fatal(err)
	}

	wall, _, err = ebitenutil.NewImageFromFile("assets/wall.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Level struct {
	Tiles         []*MapTile
	Rooms         []Rect
	PlayerVisible *fov.View
}

func NewLevel() Level {
	l := Level{}
	loadTileImages()
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
	TileType   TileType
}

func (level *Level) GetIndexFromXY(x, y int) int {
	gd := NewGameData()
	return (y * gd.ScreenWidth) + x
}

func (level *Level) CreateTiles() []*MapTile {
	gd := NewGameData()
	levelHeight = gd.ScreenHeight - gd.UIHeight
	tiles := make([]*MapTile, gd.ScreenWidth*levelHeight)

	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < levelHeight; y++ {
			index := level.GetIndexFromXY(x, y)
			tile := MapTile{
				PixelX:     x * gd.TileWidth,
				PixelY:     y * gd.TileHeight,
				Blocked:    true,
				Image:      wall,
				IsRevealed: false,
				TileType:   WALL,
			}
			tiles[index] = &tile
		}
	}
	return tiles
}

func (level *Level) DrawLevel(screen *ebiten.Image) {
	gd := NewGameData()
	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < levelHeight; y++ {
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
	for y := room.Y1; y < room.Y2+1; y++ {
		for x := room.X1; x < room.X2+1; x++ {
			index := level.GetIndexFromXY(x, y)
			level.Tiles[index].Image = floor
			level.Tiles[index].Blocked = false
			level.Tiles[index].TileType = FLOOR
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
		y := GetDiceRoll(levelHeight - h - 2)

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
	for x := min(x1, x2); x < max(x1, x2)+1; x++ {
		index := level.GetIndexFromXY(x, y)
		if index > 0 && index < gd.ScreenWidth*levelHeight {
			level.Tiles[index].Blocked = false
			level.Tiles[index].Image = floor
			level.Tiles[index].TileType = FLOOR
		}
	}
}

func (level *Level) createVerticalTunnel(y1, y2, x int) {
	gd := NewGameData()
	for y := min(y1, y2); y < max(y1, y2)+1; y++ {
		index := level.GetIndexFromXY(x, y)
		if index > 0 && index < gd.ScreenWidth*levelHeight {
			level.Tiles[index].Blocked = false
			level.Tiles[index].Image = floor
			level.Tiles[index].TileType = FLOOR
		}
	}
}

func (level Level) InBounds(x, y int) bool {
	gd := NewGameData()
	if x < 0 || x > gd.ScreenWidth || y < 0 || y > levelHeight {
		return false
	}
	return true
}

func (level Level) IsOpaque(x, y int) bool {
	idx := level.GetIndexFromXY(x, y)
	return level.Tiles[idx].TileType == WALL
}
