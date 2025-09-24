package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// PromptCommand allows players to set their prompt
type PromptCommand struct{}

// Execute runs the prompt command
func (c *PromptCommand) Execute(ch *types.Character, args string) error {
	if args == "" {
		// Display current prompt
		ch.SendMessage(fmt.Sprintf("Your current prompt is: %s\r\n", ch.Prompt))
		ch.SendMessage("Available prompt codes:\r\n")
		ch.SendMessage("%h - Current hit points\r\n")
		ch.SendMessage("%H - Maximum hit points\r\n")
		ch.SendMessage("%m - Current mana points\r\n")
		ch.SendMessage("%M - Maximum mana points\r\n")
		ch.SendMessage("%v - Current move points\r\n")
		ch.SendMessage("%V - Maximum move points\r\n")
		ch.SendMessage("%c - Your condition (only shown in combat)\r\n")
		ch.SendMessage("%C - Opponent's condition (only shown in combat)\r\n")
		ch.SendMessage("%t - Game time\r\n")
		ch.SendMessage("%% - A percent sign\r\n")
		ch.SendMessage("Example: %h/%H hp %m/%M mana %v/%V mv>\r\n")
		ch.SendMessage("Would display: 250/250 hp 267/267 mana 112/147 mv>\r\n")
		return nil
	}

	// Set the prompt
	if args == "default" {
		ch.Prompt = ">"
		ch.SendMessage("Prompt set to default.\r\n")
	} else {
		ch.Prompt = args
		ch.SendMessage(fmt.Sprintf("Prompt set to: %s\r\n", args))
	}

	return nil
}

// Name returns the name of the command
func (c *PromptCommand) Name() string {
	return "prompt"
}

// Aliases returns the aliases for the command
func (c *PromptCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *PromptCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *PromptCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *PromptCommand) LogCommand() bool {
	return false
}

// FormatPrompt formats a character's prompt with the appropriate values
func FormatPrompt(ch *types.Character) string {
	if ch.Prompt == "" {
		ch.Prompt = ">" // Default prompt
	}

	// If the prompt is just the default, return it
	if ch.Prompt == ">" {
		return ch.Prompt + " "
	}

	var result strings.Builder

	// Process each character in the prompt
	for i := 0; i < len(ch.Prompt); i++ {
		c := ch.Prompt[i]

		// Handle % escape sequences
		if c == '%' && i+1 < len(ch.Prompt) {
			next := ch.Prompt[i+1]
			switch next {
			case 'h':
				// Current hit points
				result.WriteString(fmt.Sprintf("%d", ch.HP))
				i++ // Skip the next character
			case 'H':
				// Maximum hit points
				result.WriteString(fmt.Sprintf("%d", ch.MaxHitPoints))
				i++ // Skip the next character
			case 'm':
				// Current mana points
				result.WriteString(fmt.Sprintf("%d", ch.ManaPoints))
				i++ // Skip the next character
			case 'M':
				// Maximum mana points
				result.WriteString(fmt.Sprintf("%d", ch.MaxManaPoints))
				i++ // Skip the next character
			case 'v':
				// Current move points
				result.WriteString(fmt.Sprintf("%d", ch.MovePoints))
				i++ // Skip the next character
			case 'V':
				// Maximum move points
				result.WriteString(fmt.Sprintf("%d", ch.MaxMovePoints))
				i++ // Skip the next character
			case 'c':
				// Character's condition (only in combat)
				if ch.Fighting != nil {
					result.WriteString(getConditionString(ch.HP, ch.MaxHitPoints))
				}
				i++ // Skip the next character
			case 'C':
				// Opponent's condition (only in combat)
				if ch.Fighting != nil {
					result.WriteString(getConditionString(ch.Fighting.HP, ch.Fighting.MaxHitPoints))
				}
				i++ // Skip the next character
			case 't':
				// Game time
				if w, ok := ch.World.(interface {
					GetTimeString() string
				}); ok {
					result.WriteString(w.GetTimeString())
				} else {
					// Fallback if we can't get the time
					result.WriteString("??AM")
				}
				i++ // Skip the next character
			case '%':
				// Escaped percent sign
				result.WriteRune('%')
				i++ // Skip the next character
			default:
				// Unknown escape sequence, just add the % and continue
				result.WriteRune('%')
			}
		} else {
			// Regular character, add as-is
			result.WriteRune(rune(c))
		}
	}

	// Always end with a > and a space
	if result.Len() == 0 || result.String()[result.Len()-1] != '>' {
		result.WriteString(">")
	}

	if result.Len() == 0 || result.String()[result.Len()-1] != ' ' {
		result.WriteString(" ")
	}

	return result.String()
}

// getConditionString returns a string describing a character's condition based on HP percentage
func getConditionString(hp, maxHP int) string {
	percent := float64(hp) / float64(maxHP) * 100

	switch {
	case percent >= 100:
		return "perfect"
	case percent >= 90:
		return "excellent"
	case percent >= 75:
		return "good"
	case percent >= 50:
		return "fair"
	case percent >= 30:
		return "wounded"
	case percent >= 15:
		return "bad"
	case percent >= 0:
		return "awful"
	default:
		return "dying"
	}
}

// getDirectionName returns the name of a direction
func getDirectionName(dir int) string {
	directions := map[int]string{
		types.DIR_NORTH: "north",
		types.DIR_EAST:  "east",
		types.DIR_SOUTH: "south",
		types.DIR_WEST:  "west",
		types.DIR_UP:    "up",
		types.DIR_DOWN:  "down",
	}

	if name, ok := directions[dir]; ok {
		return name
	}
	return "unknown"
}
