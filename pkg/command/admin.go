package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// GotoCommand is a command that teleports a character to a specific room
type GotoCommand struct {
	vnum int
}

// Name returns the name of the command
func (c *GotoCommand) Name() string {
	return "goto"
}

// Aliases returns the aliases of the command
func (c *GotoCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *GotoCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *GotoCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *GotoCommand) LogCommand() bool {
	return true
}

// Execute executes the goto command
func (c *GotoCommand) Execute(character *types.Character, args string) error {
	// Parse the room VNUM from args
	vnumStr := strings.TrimSpace(args)
	if vnumStr == "" {
		return fmt.Errorf("usage: goto <room vnum>")
	}

	vnum, err := strconv.Atoi(vnumStr)
	if err != nil {
		return fmt.Errorf("invalid room vnum: %s", vnumStr)
	}

	// Get the world from the character
	world, ok := character.World.(interface{ GetRoom(int) *types.Room })
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Get the destination room
	destRoom := world.GetRoom(vnum)
	if destRoom == nil {
		return fmt.Errorf("room %d does not exist", vnum)
	}

	// Remove character from current room
	if character.InRoom != nil {
		for i, ch := range character.InRoom.Characters {
			if ch == character {
				character.InRoom.Characters = append(character.InRoom.Characters[:i], character.InRoom.Characters[i+1:]...)
				break
			}
		}
	}

	// Add character to destination room
	character.InRoom = destRoom
	destRoom.Characters = append(destRoom.Characters, character)

	// Show the room to the character
	lookCmd := &LookCommand{}
	return lookCmd.Execute(character, "")
}

// RstatCommand is a command that shows detailed information about a room
type RstatCommand struct {
	vnum int
}

// Name returns the name of the command
func (c *RstatCommand) Name() string {
	return "rstat"
}

// Aliases returns the aliases of the command
func (c *RstatCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *RstatCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *RstatCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *RstatCommand) LogCommand() bool {
	return true
}

// Execute executes the rstat command
func (c *RstatCommand) Execute(character *types.Character, args string) error {
	var room *types.Room

	// Parse the room VNUM from args if provided
	vnumStr := strings.TrimSpace(args)
	if vnumStr == "" {
		// Use character's current room
		room = character.InRoom
	} else {
		// Parse the VNUM
		vnum, err := strconv.Atoi(vnumStr)
		if err != nil {
			return fmt.Errorf("invalid room vnum: %s", vnumStr)
		}

		// Get the world from the character
		world, ok := character.World.(interface{ GetRoom(int) *types.Room })
		if !ok {
			return fmt.Errorf("world interface not available")
		}

		// Get the room by VNUM
		room = world.GetRoom(vnum)
	}

	if room == nil {
		return fmt.Errorf("room does not exist")
	}

	// Build the room stats
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\r\n--- Room Stats for Room %d ---\r\n", room.VNUM))
	sb.WriteString(fmt.Sprintf("Name: %s\r\n", room.Name))
	sb.WriteString(fmt.Sprintf("Description: %s\r\n", room.Description))
	sb.WriteString(fmt.Sprintf("Flags: %d\r\n", room.Flags))
	sb.WriteString(fmt.Sprintf("Sector Type: %d\r\n", room.SectorType))

	// Show exits
	sb.WriteString("\r\nExits:\r\n")
	for dir := 0; dir < 6; dir++ {
		exit := room.Exits[dir]
		if exit != nil {
			dirName := directionName(dir)
			sb.WriteString(fmt.Sprintf("  %s: To Room %d, Flags: %d, Key: %d\r\n",
				dirName, exit.DestVnum, exit.Flags, exit.Key))
			sb.WriteString(fmt.Sprintf("     Description: %s\r\n", exit.Description))
			sb.WriteString(fmt.Sprintf("     Keywords: %s\r\n", exit.Keywords))
		}
	}

	// Show characters in the room
	sb.WriteString("\r\nCharacters:\r\n")
	for _, ch := range room.Characters {
		if ch.IsNPC {
			// Show more details for NPCs
			sb.WriteString(fmt.Sprintf("  NPC: %s (VNUM: %d)\r\n",
				ch.Name, ch.Prototype.VNUM))
			sb.WriteString(fmt.Sprintf("     Short: %s\r\n", ch.ShortDesc))
			sb.WriteString(fmt.Sprintf("     Long: %s\r\n", ch.LongDesc))
			sb.WriteString(fmt.Sprintf("     Level: %d\r\n", ch.Level))
		} else {
			// Just show the name for players
			sb.WriteString(fmt.Sprintf("  Player: %s\r\n", ch.Name))
		}
	}

	// Show objects in the room
	sb.WriteString("\r\nObjects:\r\n")
	for _, obj := range room.Objects {
		if obj.Prototype != nil {
			sb.WriteString(fmt.Sprintf("  %s\r\n", obj.Prototype.Name))
		} else {
			sb.WriteString("  Unknown object\r\n")
		}
	}

	// Send the stats to the character
	character.SendMessage(sb.String())
	return nil
}
