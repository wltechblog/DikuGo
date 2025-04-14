package command

import "github.com/wltechblog/DikuGo/pkg/types"

// directionName returns the name of a direction
func directionName(dir int) string {
	switch dir {
	case types.DIR_NORTH:
		return "north"
	case types.DIR_EAST:
		return "east"
	case types.DIR_SOUTH:
		return "south"
	case types.DIR_WEST:
		return "west"
	case types.DIR_UP:
		return "up"
	case types.DIR_DOWN:
		return "down"
	default:
		return "unknown"
	}
}
