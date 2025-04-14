package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// HelpCommand represents the help command
type HelpCommand struct {
	// Registry is the command registry
	Registry *Registry
}

// Execute executes the help command
func (c *HelpCommand) Execute(character *types.Character, args string) error {
	// If no arguments, show general help
	if args == "" {
		return c.showGeneralHelp(character)
	}

	// Show help for a specific command
	return c.showCommandHelp(character, args)
}

// showGeneralHelp shows general help
func (c *HelpCommand) showGeneralHelp(character *types.Character) error {
	// Build the help text
	var sb strings.Builder

	sb.WriteString("\r\nAvailable commands:\r\n")
	sb.WriteString("------------------\r\n")

	// Get all commands
	commands := make(map[string]bool)
	for name, cmd := range c.Registry.commands {
		// Skip aliases
		if name == cmd.Name() {
			commands[name] = true
		}
	}

	// Add each command to the list
	var commandList []string
	for name := range commands {
		commandList = append(commandList, name)
	}

	// Sort the commands
	// TODO: Sort the commands

	// Add the commands to the help text
	sb.WriteString(strings.Join(commandList, ", "))
	sb.WriteString("\r\n\r\n")
	sb.WriteString("Type 'help <command>' for help on a specific command.\r\n")

	// Send the help text to the character
	// TODO: Implement a way to send messages to characters
	// For now, just return the help text as an error
	return fmt.Errorf("%s", sb.String())
}

// showCommandHelp shows help for a specific command
func (c *HelpCommand) showCommandHelp(character *types.Character, commandName string) error {
	// Find the command
	cmd := c.Registry.Find(commandName)
	if cmd == nil {
		return fmt.Errorf("no help available for '%s'", commandName)
	}

	// Build the help text
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\r\nHelp for '%s':\r\n", cmd.Name()))
	sb.WriteString("----------------\r\n")

	// Add the command description
	// TODO: Add command descriptions
	sb.WriteString("No description available.\r\n")

	// Add the command aliases
	if len(cmd.Aliases()) > 0 {
		sb.WriteString(fmt.Sprintf("\r\nAliases: %s\r\n", strings.Join(cmd.Aliases(), ", ")))
	}

	// Send the help text to the character
	// TODO: Implement a way to send messages to characters
	// For now, just return the help text as an error
	return fmt.Errorf("%s", sb.String())
}

// Name returns the name of the command
func (c *HelpCommand) Name() string {
	return "help"
}

// Aliases returns the aliases of the command
func (c *HelpCommand) Aliases() []string {
	return []string{"?"}
}

// MinPosition returns the minimum position required to execute the command
func (c *HelpCommand) MinPosition() int {
	return types.POS_SLEEPING
}

// Level returns the minimum level required to execute the command
func (c *HelpCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *HelpCommand) LogCommand() bool {
	return false
}
