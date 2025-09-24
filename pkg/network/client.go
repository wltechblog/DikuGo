package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wltechblog/DikuGo/pkg/command"
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/ui"
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
	StateMainMenu
	StatePlaying
	StateReadMOTD
	StateChangePassword
	StateConfirmNewPasswordChange
	StateReadStory
	StateDeleteCharacter
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
	shutdownCh      chan struct{} // Channel to signal shutdown
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
		shutdownCh:      make(chan struct{}),
	}
}

// Handle handles the client connection
func (c *Client) Handle() {
	defer c.Close()

	// Send welcome banner
	c.Write(ui.Banner)
	c.Write("By what name do you wish to be known? ")

	// Main loop
	for !c.Closed {
		// Use a select statement to handle both input and shutdown signals
		select {
		case <-c.shutdownCh:
			// Shutdown signal received
			log.Printf("Client %s received shutdown signal", c.ID)
			return
		default:
			// Set a read timeout to avoid blocking indefinitely
			c.Conn.SetReadDeadline(time.Now().Add(1 * time.Second))

			// Read input
			input, err := c.Read()
			if err != nil {
				// Check if it's a timeout error
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Timeout - continue the loop to check for shutdown
					continue
				}
				// Other error - client disconnected
				log.Printf("Error reading from client %s: %v", c.ID, err)
				break
			}

			// Clear the read deadline
			c.Conn.SetReadDeadline(time.Time{})

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
			case StateGetClass:
				c.HandleGetClass(input)
			case StateMainMenu:
				c.HandleMainMenu(input)
			case StateReadMOTD:
				c.HandleReadMOTD(input)
			case StateChangePassword:
				c.HandleChangePassword(input)
			case StateConfirmNewPasswordChange:
				c.HandleConfirmNewPasswordChange(input)
			case StateReadStory:
				c.HandleReadStory(input)
			case StateDeleteCharacter:
				c.HandleDeleteCharacter(input)
			case StatePlaying:
				c.HandleCommand(input)
			default:
				c.Write("Invalid state. Please try again: ")
				c.State = StateGetName
			}
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

	// Mark as closed first to prevent recursive calls
	c.Closed = true

	// Signal shutdown to the client goroutine
	select {
	case <-c.shutdownCh:
		// Channel already closed
	default:
		close(c.shutdownCh)
	}

	// Unregister client from the client registry
	if c.Character != nil {
		log.Printf("Unregistering client for character %s", c.Character.Name)
		UnregisterClient(c.Character)

		// Save character before disconnecting
		if c.World != nil {
			log.Printf("Saving character %s before disconnecting", c.Character.Name)
			err := c.World.SaveCharacter(c.Character)
			if err != nil {
				log.Printf("Error saving character %s during close: %v", c.Character.Name, err)
			}
		}
	}

	// Close the connection
	log.Printf("Closing connection for client %s", c.ID)
	err := c.Conn.Close()
	if err != nil {
		log.Printf("Error closing connection: %v", err)
	}

	log.Printf("Client %s closed", c.ID)
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

	// Password is correct, show menu
	c.Character = character
	c.Write(fmt.Sprintf("\r\nWelcome back, %s!\r\n", character.Name))
	c.Write(ui.Menu)
	c.State = StateMainMenu
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

	// Store name and password in InputBuf
	c.InputBuf = fmt.Sprintf("%s:%s", name, password)

	// Prompt for class selection
	c.Write("\r\nSelect a class:\r\n")
	c.Write("[1] Magic User - Masters of arcane magic\r\n")
	c.Write("[2] Cleric - Healers with divine magic\r\n")
	c.Write("[3] Thief - Stealthy rogues with special skills\r\n")
	c.Write("[4] Warrior - Strong fighters skilled in combat\r\n")
	c.Write("\r\nEnter your choice (1-4): ")
	c.State = StateGetClass
}

// HandleMainMenu handles the main menu state
func (c *Client) HandleMainMenu(input string) {
	switch input {
	case "0": // Exit
		c.Write("Goodbye!\r\n")
		c.Closed = true
		return
	case "1": // Enter the game
		// Register client with the client registry
		RegisterClient(c.Character, c)

		// Add character to world
		log.Printf("HandleMainMenu: Adding character %s to world", c.Character.Name)
		c.World.AddCharacter(c.Character)
		log.Printf("HandleMainMenu: Character %s added to world, in room: %v", c.Character.Name, c.Character.InRoom)

		// Set the World field in the character
		c.Character.World = c.World

		// Enter game
		c.State = StatePlaying

		// Show room description
		if c.Character.InRoom != nil {
			c.Write(fmt.Sprintf("\r\n%s\r\n%s\r\n", c.Character.InRoom.Name, c.Character.InRoom.Description))
		}

		c.Write("Enter your command: ")
	case "2": // Enter description
		c.Write("Enter a description for your character. Terminate with a '@'.\r\n")
		// TODO: Implement description editing
		c.Write(ui.Menu)
	case "3": // Read the background story
		c.Write(ui.Story)
		c.State = StateReadStory
	case "4": // Change password
		c.Write("Enter a new password: ")
		c.State = StateChangePassword
	case "5": // Delete character
		c.Write("\r\nWARNING: This will permanently delete your character!\r\n")
		c.Write("Type your character's name to confirm deletion: ")
		c.State = StateDeleteCharacter
	default:
		c.Write("Invalid choice.\r\n")
		c.Write(ui.Menu)
	}
}

// HandleReadMOTD handles the read MOTD state
func (c *Client) HandleReadMOTD(input string) {
	c.Write(ui.Menu)
	c.State = StateMainMenu
}

// HandleReadStory handles the read story state
func (c *Client) HandleReadStory(input string) {
	c.Write(ui.Menu)
	c.State = StateMainMenu
}

// HandleChangePassword handles the change password state
func (c *Client) HandleChangePassword(password string) {
	if len(password) < 5 {
		c.Write("Password too short. Please choose a password of at least 5 characters: ")
		return
	}

	c.InputBuf = password
	c.Write("Please confirm your new password: ")
	c.State = StateConfirmNewPasswordChange
}

// HandleGetClass handles the class selection state
func (c *Client) HandleGetClass(input string) {
	// Parse the class choice
	classChoice, err := strconv.Atoi(input)
	if err != nil || classChoice < 1 || classChoice > 4 {
		c.Write("Invalid choice. Please enter a number between 1 and 4: ")
		return
	}

	// Get name and password from InputBuf
	parts := strings.Split(c.InputBuf, ":")
	name := parts[0]
	password := parts[1]

	// Create character with the selected class
	character := &types.Character{
		Name:          name,
		Password:      password,
		Level:         1,
		Position:      types.POS_STANDING,
		ShortDesc:     name,
		Class:         classChoice,
		RoomVNUM:      -1, // No saved room, will be placed in default starting room
		Skills:        make(map[int]int),
		Spells:        make(map[int]int),
		Equipment:     make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory:     make([]*types.ObjectInstance, 0),
		LastSkillTime: make(map[int]time.Time),
		LastLogin:     time.Now(),
		Title:         " the newbie",
		Prompt:        "%h/%H hp %m/%M mana %v/%V mv> ",
	}

	// Initialize character abilities and stats based on class
	character.World = c.World

	// Initialize the character with appropriate stats
	c.World.InitializeNewCharacter(character)

	// Save character
	err = c.World.SaveCharacter(character)
	if err != nil {
		c.Write("Error creating character. Please try again: ")
		c.State = StateGetName
		return
	}

	// Set character and show menu
	c.Character = character

	// Display character stats
	c.Write(fmt.Sprintf("\r\nWelcome, %s the %s!\r\n",
		character.Name, getClassName(classChoice)))
	c.Write("\r\nYour starting attributes are:\r\n")
	c.Write(fmt.Sprintf("Strength:     %d\r\n", character.Abilities[0]))
	c.Write(fmt.Sprintf("Intelligence: %d\r\n", character.Abilities[1]))
	c.Write(fmt.Sprintf("Wisdom:       %d\r\n", character.Abilities[2]))
	c.Write(fmt.Sprintf("Dexterity:    %d\r\n", character.Abilities[3]))
	c.Write(fmt.Sprintf("Constitution: %d\r\n", character.Abilities[4]))
	c.Write(fmt.Sprintf("Charisma:     %d\r\n", character.Abilities[5]))
	c.Write("\r\nPress ENTER to continue...\r\n")

	c.Write(ui.Menu)
	c.State = StateMainMenu
}

// getClassName returns the name of a class
func getClassName(class int) string {
	return types.GetClassName(class)
}

// HandleDeleteCharacter handles the delete character state
func (c *Client) HandleDeleteCharacter(input string) {
	// Check if the input matches the character name
	if input != c.Character.Name {
		c.Write("Character name does not match. Deletion cancelled.\r\n")
		c.Write(ui.Menu)
		c.State = StateMainMenu
		return
	}

	// Delete the character
	err := c.World.DeleteCharacter(c.Character.Name)
	if err != nil {
		c.Write(fmt.Sprintf("Error deleting character: %s\r\n", err.Error()))
		c.Write(ui.Menu)
		c.State = StateMainMenu
		return
	}

	// Character deleted successfully
	c.Write(fmt.Sprintf("Character %s has been permanently deleted.\r\n", c.Character.Name))
	c.Character = nil
	c.Write("By what name do you wish to be known? ")
	c.State = StateGetName
}

// HandleConfirmNewPasswordChange handles the confirm new password change state
func (c *Client) HandleConfirmNewPasswordChange(confirm string) {
	if confirm != c.InputBuf {
		c.Write("Passwords don't match. Please try again.\r\n")
		c.Write("Enter a new password: ")
		c.State = StateChangePassword
		return
	}

	// Update password
	c.Character.Password = confirm

	// Save character
	err := c.World.SaveCharacter(c.Character)
	if err != nil {
		c.Write("Error saving password. Please try again.\r\n")
		c.Write(ui.Menu)
		c.State = StateMainMenu
		return
	}

	c.Write("Password changed successfully.\r\n")
	c.Write(ui.Menu)
	c.State = StateMainMenu
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
			// Save character
			c.World.SaveCharacter(c.Character)

			// Remove character from world
			c.World.RemoveCharacter(c.Character)

			// Unregister client
			UnregisterClient(c.Character)

			// Show menu
			c.Write(ui.Menu)
			c.State = StateMainMenu
			return
		}

		// Send error message to client
		c.Write(err.Error() + "\r\n")
	}

	// Check if the character has died and needs to return to the menu
	if c.Character != nil && c.Character.HasMessage("RETURN_TO_MENU") {
		// Clear the message
		c.Character.ClearMessage("RETURN_TO_MENU")

		// Save character
		c.World.SaveCharacter(c.Character)

		// Remove character from world
		c.World.RemoveCharacter(c.Character)

		// Unregister client
		UnregisterClient(c.Character)

		// Show menu
		c.Write(ui.Menu)
		c.State = StateMainMenu
		return
	}

	// Get the formatted prompt
	if c.CommandRegistry != nil && c.Character != nil {
		// Use the FormatPrompt function directly
		prompt := c.CommandRegistry.FormatPrompt(c.Character)
		c.Write(prompt)
	} else {
		c.Write("Enter your command: ")
	}
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
