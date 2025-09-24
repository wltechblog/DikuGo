package command

import (
	"errors"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// Command represents a command that can be executed by a character
type Command interface {
	// Execute executes the command
	Execute(character *types.Character, args string) error

	// Name returns the name of the command
	Name() string

	// Aliases returns the aliases of the command
	Aliases() []string

	// MinPosition returns the minimum position required to execute the command
	MinPosition() int

	// Level returns the minimum level required to execute the command
	Level() int

	// LogCommand returns whether the command should be logged
	LogCommand() bool
}

// Registry is a registry of commands
type Registry struct {
	commands map[string]Command
}

// NewRegistry creates a new command registry
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

// Register registers a command with the registry
func (r *Registry) Register(cmd Command) {
	// Register the command by its name
	r.commands[cmd.Name()] = cmd

	// Register the command by its aliases
	for _, alias := range cmd.Aliases() {
		r.commands[alias] = cmd
	}
}

// Find finds a command by name or alias
func (r *Registry) Find(name string) Command {
	return r.commands[name]
}

// Execute executes a command
func (r *Registry) Execute(character *types.Character, input string) error {
	// Parse the command and arguments
	cmdName, args := parseCommand(input)
	if cmdName == "" {
		return nil
	}

	// Find the command
	cmd := r.Find(cmdName)
	if cmd == nil {
		return ErrCommandNotFound
	}

	// Check if the character meets the minimum position requirement
	if character.Position < cmd.MinPosition() {
		// Return a more specific error message based on the character's position
		switch character.Position {
		case types.POS_DEAD:
			return errors.New("Lie still; you are DEAD!!!")
		case types.POS_MORTALLY, types.POS_INCAP:
			return errors.New("You are in a pretty bad shape, unable to do anything!")
		case types.POS_STUNNED:
			return errors.New("You are too stunned to do that!")
		case types.POS_SLEEPING:
			return errors.New("You can't do that while sleeping!")
		case types.POS_RESTING:
			return errors.New("You can't do that while resting!")
		case types.POS_SITTING:
			return errors.New("You can't do that while sitting!")
		case types.POS_FIGHTING:
			return errors.New("You can't do that while fighting!")
		default:
			return errors.New("You are in the wrong position for that!")
		}
	}

	// Check if the character meets the minimum level requirement
	if character.Level < cmd.Level() {
		return ErrInsufficientLevel
	}

	// Execute the command
	return cmd.Execute(character, args)
}

// parseCommand parses a command string into a command name and arguments
func parseCommand(input string) (string, string) {
	// Find the first space
	for i, c := range input {
		if c == ' ' {
			return input[:i], input[i+1:]
		}
	}

	// No space found, the whole input is the command name
	return input, ""
}

// FormatPrompt formats a character's prompt
func (r *Registry) FormatPrompt(ch *types.Character) string {
	return FormatPrompt(ch)
}
