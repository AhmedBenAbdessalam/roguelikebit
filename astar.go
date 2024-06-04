package main

import (
	"errors"
	"reflect"
)

// Node represents a given point on a map
// g is the total distance of the node from the start
// h is the estimated distance of the node from the ending
// f is the total value from the node (g+h)
type node struct {
	Parent   *node
	Position *Position
	g        int
	h        int
	f        int
}

func newNode(parent *node, position *Position) *node {
	n := node{}
	n.Parent = parent
	n.Position = position
	n.g = 0
	n.h = 0
	n.f = 0
	return &n
}

func (n *node) isEqual(other *node) bool {
	return n.Position.IsEqual(other.Position)
}

func reverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic(errors.New("data must be a slice type"))
	}
	valueLen := value.Len()
	for i := 0; i <= int((valueLen-1)/2); i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}

func isInSlice(s []*node, target *node) bool {
	for _, n := range s {
		if n.isEqual(target) {
			return true
		}
	}
	return false
}

// AStar implements the AStar Algorithm.
type AStar struct{}

// GetPath takes a level, the starting position and an ending position (the goal) and returns
// a list of Positions which is the path between the points.
func (as AStar) GetPath(level Level, start *Position, end *Position) []Position {
	gd := NewGameData()
	openList := make([]*node, 0)
	closedList := make([]*node, 0)

	//Create our starting point
	startNode := newNode(nil, start)
	endNodePlaceholder := newNode(nil, end)

	openList = append(openList, startNode)

	for {
		if len(openList) == 0 {
			break
		}
		//get current node
		currentNode := openList[0]
		currentIndex := 0

		//get the node with the smallest f value
		for index, item := range openList {
			if item.f < currentNode.f {
				currentNode = item
				currentIndex = index
			}
		}

		//Move from open to closed list
		openList = append(openList[:currentIndex], openList[currentIndex+1:]...)
		closedList = append(closedList, currentNode)

		//check to see if we reached our end
		//if so, we are done here
		if currentNode.isEqual(endNodePlaceholder) {
			path := make([]Position, 0)
			current := currentNode
			for {
				if current == nil {
					break
				}
				path = append(path, *current.Position)
				current = current.Parent
			}

			//reverse the path and return it
			reverseSlice(path)
			return path
		}

		edges := make([]*node, 0)
		//up edge
		if currentNode.Position.Y > 0 {
			tile := level.Tiles[level.GetIndexFromXY(currentNode.Position.X, currentNode.Position.Y-1)]
			if tile.TileType != WALL {
				upNodePosition := Position{
					X: currentNode.Position.X,
					Y: currentNode.Position.Y - 1,
				}
				newNode := newNode(currentNode, &upNodePosition)
				edges = append(edges, newNode)
			}
		}
		//down edge
		if currentNode.Position.Y < gd.ScreenHeight {
			tile := level.Tiles[level.GetIndexFromXY(currentNode.Position.X, currentNode.Position.Y+1)]
			if tile.TileType != WALL {
				downNodePosition := Position{
					X: currentNode.Position.X,
					Y: currentNode.Position.Y + 1,
				}
				newNode := newNode(currentNode, &downNodePosition)
				edges = append(edges, newNode)
			}
		}
		//left edge
		if currentNode.Position.X > 0 {
			tile := level.Tiles[level.GetIndexFromXY(currentNode.Position.X-1, currentNode.Position.Y)]
			if tile.TileType != WALL {
				leftNodePosition := Position{
					X: currentNode.Position.X - 1,
					Y: currentNode.Position.Y,
				}
				newNode := newNode(currentNode, &leftNodePosition)
				edges = append(edges, newNode)
			}
		}
		//right edge
		if currentNode.Position.X < gd.ScreenWidth {
			tile := level.Tiles[level.GetIndexFromXY(currentNode.Position.X+1, currentNode.Position.Y)]
			if tile.TileType != WALL {
				rightNodePosition := Position{
					X: currentNode.Position.X + 1,
					Y: currentNode.Position.Y,
				}
				newNode := newNode(currentNode, &rightNodePosition)
				edges = append(edges, newNode)
			}
		}
		//Now iterate through the edges and put them in the open list
		for _, edge := range edges {
			if isInSlice(closedList, edge) {
				continue
			}
			edge.g = currentNode.g + 1
			edge.h = edge.Position.GetManhattanDistance(endNodePlaceholder.Position)
			edge.f = edge.g + edge.h

			if isInSlice(openList, edge) {
				isFurther := false
				for _, n := range openList {
					if edge.g > n.g {
						isFurther = true
						break
					}
				}
				if isFurther {
					continue
				}

			}
			openList = append(openList, edge)
		}
	}
	return nil
}
