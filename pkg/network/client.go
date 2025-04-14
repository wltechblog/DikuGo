package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wltechblog/DikuGo/pkg/command"
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// ConnectionState represents the state of a client connection
type ConnectionState int

const (
	StateGetName ConnectionState = iota
	StateGetPassword
	StateConfirmPassword
	StateGetNewPassword
	StateConfirmNewPassword
	StateGetEmail
	StateGetGender
	StateGetClass
	StateGetRace
	StateGetAlignment
	StateGetStats
	StateConfirmCharacter
	StatePlaying
	StateDisconnected
)

// Client represents a connected client
type Client struct {
	ID              string
	Conn            net.Conn
	Reader          *bufio.Reader
	Writer          *bufio.Writer
	State           ConnectionState
	Character       *types.Character
	World           *world.World
	InputBuf        string
	OutputBuf       []string
	LastInput       time.Time
	LoginTime       time.Time
	Closed          bool
	CommandRegistry *command.Registry
}

// NewClient creates a new client instance
func NewClient(conn net.Conn, w *world.World, cmdRegistry *command.Registry) *Client {
	return &Client{
		ID:              uuid.New().String(),
		Conn:            conn,
		Reader:          bufio.NewReader(conn),
		Writer:          bufio.NewWriter(conn),
		State:           StateGetName,
		World:           w,
		OutputBuf:       make([]string, 0),
		LastInput:       time.Now(),
		LoginTime:       time.Now(),
		Closed:          false,
		CommandRegistry: cmdRegistry,
	}
}

// Handle handles the client connection
func (c *Client) Handle() {
	defer c.Close()

	// Send welcome message
	c.Write("\r\nWelcome to DikuGo!\r\n")
	c.Write("By what name do you wish to be known? ")

	// Main loop
	for !c.Closed {
		// Read input
		input, err := c.Read()
		if err != nil {
			log.Printf("Error reading from client %s: %v", c.ID, err)
			break
		}

		// Update last activity
		c.LastInput = time.Now()

		// Handle input based on state
		switch c.State {
		case StateGetName:
			c.HandleGetName(input)
		case StateGetPassword:
			c.HandleGetPassword(input)
		case StateConfirmPassword:
			c.HandleConfirmPassword(input)
		case StateGetNewPassword:
			c.HandleGetNewPassword(input)
		case StateConfirmNewPassword:
			c.HandleConfirmNewPassword(input)
		case StatePlaying:
			c.HandleCommand(input)
		default:
			c.Write("Invalid state. Please try again: ")
			c.State = StateGetName
		}
	}

	// Clean up
	if c.Character != nil {
		// Remove character from world
		c.World.RemoveCharacter(c.Character)

		// Save character
		err := c.World.SaveCharacter(c.Character)
		if err != nil {
			log.Printf("Error saving character %s: %v", c.Character.Name, err)
		}
	}
}

// Read reads a line from the client
func (c *Client) Read() (string, error) {
	// Read a line
	line, err := c.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Trim whitespace
	line = strings.TrimSpace(line)

	return line, nil
}

// Write writes a message to the client
func (c *Client) Write(message string) {
	_, err := c.Writer.WriteString(message)
	if err != nil {
		log.Printf("Error writing to client: %v", err)
		return
	}
	err = c.Writer.Flush()
	if err != nil {
		log.Printf("Error flushing writer: %v", err)
	}
}

// Close closes the client connection
func (c *Client) Close() {
	if c.Closed {
		return
	}

	c.Closed = true
	err := c.Conn.Close()
	if err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}

// HandleGetName handles the get name state
func (c *Client) HandleGetName(name string) {
	// Check if name is valid
	if !isValidName(name) {
		c.Write("Invalid name. Please try again: ")
		return
	}

	// Check if character exists
	if c.World.CharacterExists(name) {
		// Character exists, ask for password
		c.Write("Password: ")
		c.InputBuf = name
		c.State = StateGetPassword
		return
	}

	// Character doesn't exist, create a new one
	c.Write(fmt.Sprintf("I don't know %s. Is this a new character? (Y/N) ", name))
	c.InputBuf = name
	c.State = StateConfirmPassword
}

// HandleGetPassword handles the get password state
func (c *Client) HandleGetPassword(password string) {
	// Get character
	character, err := c.World.GetCharacter(c.InputBuf)
	if err != nil {
		c.Write("Error loading character. Please try again: ")
		c.State = StateGetName
		return
	}

	// Check password
	if character.Password != password {
		c.Write("Wrong password. Please try again: ")
		c.State = StateGetName
		return
	}

	// Password is correct, enter game
	c.Character = character
	c.State = StatePlaying

	// World field will be set in AddCharacter

	// Add character to world
	log.Printf("HandleGetPassword: Adding character %s to world", character.Name)
	c.World.AddCharacter(character)
	log.Printf("HandleGetPassword: Character %s added to world, in room: %v", character.Name, character.InRoom)

	// Send welcome message
	c.Write(fmt.Sprintf("\r\nWelcome back, %s!\r\n", character.Name))

	// Show room description
	if character.InRoom != nil {
		c.Write(fmt.Sprintf("\r\n%s\r\n%s\r\n", character.InRoom.Name, character.InRoom.Description))
	}

	c.Write("Enter your command: ")
}

// HandleConfirmPassword handles the confirm password state
func (c *Client) HandleConfirmPassword(confirm string) {
	// Check if user wants to create a new character
	if strings.ToLower(confirm) != "y" {
		c.Write("By what name do you wish to be known? ")
		c.State = StateGetName
		return
	}

	// Ask for password
	c.Write("Please choose a password: ")
	c.State = StateGetNewPassword
}

// HandleGetNewPassword handles the get new password state
func (c *Client) HandleGetNewPassword(password string) {
	// Check if password is valid
	if len(password) < 5 {
		c.Write("Password too short. Please choose a password of at least 5 characters: ")
		return
	}

	// Store password
	c.InputBuf = fmt.Sprintf("%s:%s", c.InputBuf, password)
	c.Write("Please confirm your password: ")
	c.State = StateConfirmNewPassword
}

// HandleConfirmNewPassword handles the confirm new password state
func (c *Client) HandleConfirmNewPassword(confirm string) {
	// Check if passwords match
	parts := strings.Split(c.InputBuf, ":")
	name := parts[0]
	password := parts[1]

	if confirm != password {
		c.Write("Passwords don't match. Please choose a password: ")
		c.InputBuf = name
		c.State = StateGetNewPassword
		return
	}

	// Create character
	character := &types.Character{
		Name:      name,
		Password:  password,
		Level:     1,
		Position:  types.POS_STANDING,
		ShortDesc: name,
	}

	// Save character
	err := c.World.SaveCharacter(character)
	if err != nil {
		c.Write("Error creating character. Please try again: ")
		c.State = StateGetName
		return
	}

	// Enter game
	c.Character = character
	c.State = StatePlaying

	// Set the World field in the character
	character.World = c.World

	// Add character to world
	log.Printf("Client: Adding character %s to world", character.Name)
	c.World.AddCharacter(character)
	log.Printf("Client: Character %s added to world, in room: %v", character.Name, character.InRoom)

	// Send welcome message
	c.Write(fmt.Sprintf("\r\nWelcome, %s!\r\n", character.Name))

	// Show room description
	if character.InRoom != nil {
		c.Write(fmt.Sprintf("\r\n%s\r\n%s\r\n", character.InRoom.Name, character.InRoom.Description))
	}

	c.Write("Enter your command: ")
}

// HandleCommand handles a game command
func (c *Client) HandleCommand(input string) {
	if input == "" {
		c.Write("Enter your command: ")
		return
	}

	// Execute the command
	err := c.CommandRegistry.Execute(c.Character, input)
	if err != nil {
		// Check if the error is a quit command
		if strings.HasSuffix(err.Error(), "QUIT") {
			// Close the connection
			c.Closed = true
			return
		}

		// Send error message to client
		c.Write(err.Error() + "\r\n")
	}

	c.Write("Enter your command: ")
}

// isValidName checks if a name is valid
func isValidName(name string) bool {
	// Check if name is empty
	if name == "" {
		return false
	}

	// Check if name is too short
	if len(name) < 3 {
		return false
	}

	// Check if name is too long
	if len(name) > 12 {
		return false
	}

	// Check if name contains only letters
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			return false
		}
	}

	return true
}
