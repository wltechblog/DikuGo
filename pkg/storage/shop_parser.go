package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ParseShops parses the shop file and returns a slice of shops
func ParseShops(filename string) ([]*types.Shop, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	var shops []*types.Shop
	var currentShop *types.Shop
	var lineNum int
	var state string
	var skipUntilNextShop bool

	// Read the file line by line
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for end of file marker
		if line == "$~" {
			break
		}

		// Check for new shop
		if strings.HasPrefix(line, "#") {
			// Save the current shop if it exists
			if currentShop != nil {
				shops = append(shops, currentShop)
			}

			// Reset the skip flag
			skipUntilNextShop = false

			// Parse shop number
			vnumStr := strings.TrimPrefix(line, "#")
			vnum, err := strconv.Atoi(strings.TrimSuffix(vnumStr, "~"))
			if err != nil {
				return nil, fmt.Errorf("invalid shop number on line %d: %w", lineNum, err)
			}

			// Create a new shop
			currentShop = &types.Shop{
				VNUM:      vnum,
				BuyTypes:  []int{},
				Producing: []int{},
				Messages:  []string{},
			}

			// Set the state to read producing items next
			state = "producing"
			continue
		}

		// Skip lines until we find a new shop if we're in skip mode
		if skipUntilNextShop {
			continue
		}

		// Process the line based on the current state
		switch state {
		case "producing":
			// Parse producing items
			producingItem, err := strconv.Atoi(line)
			if err != nil {
				// Skip this shop
				skipUntilNextShop = true
				continue
			}

			if producingItem != -1 {
				currentShop.Producing = append(currentShop.Producing, producingItem)
			}

			// Check if we've read all 5 producing items
			if len(currentShop.Producing) >= 5 {
				state = "profitbuy"
			}
		case "buytypes":
			// Parse buy types
			buyType, err := strconv.Atoi(line)
			if err != nil {
				// Skip this shop
				skipUntilNextShop = true
				continue
			}

			if buyType != -1 {
				currentShop.BuyTypes = append(currentShop.BuyTypes, buyType)
			}

			// Check if we've read all 5 buy types
			if len(currentShop.BuyTypes) >= 5 {
				state = "profitbuy"
			}
		case "profitbuy":
			// Parse profit buy
			profitBuy, err := strconv.ParseFloat(line, 64)
			if err != nil {
				// Skip this shop
				skipUntilNextShop = true
				continue
			}

			currentShop.ProfitBuy = profitBuy
			state = "profitsell"
		case "profitsell":
			// Parse profit sell
			profitSell, err := strconv.ParseFloat(line, 64)
			if err != nil {
				// Skip this shop
				skipUntilNextShop = true
				continue
			}

			currentShop.ProfitSell = profitSell
			state = "openhour"
		case "openhour":
			// Parse open hour
			openHour, err := strconv.Atoi(line)
			if err != nil {
				// Skip this shop
				skipUntilNextShop = true
				continue
			}

			currentShop.OpenHour = openHour
			state = "closehour"
		case "closehour":
			// Parse close hour
			closeHour, err := strconv.Atoi(line)
			if err != nil {
				// Skip this shop
				skipUntilNextShop = true
				continue
			}

			currentShop.CloseHour = closeHour
			state = "keeper"
		case "keeper":
			// Parse keeper (mobile VNUM)
			keeper, err := strconv.Atoi(line)
			if err != nil {
				// Skip this shop
				skipUntilNextShop = true
				continue
			}

			currentShop.MobileVNUM = keeper
			state = "roomvnum"
		case "roomvnum":
			// Parse room VNUM
			roomVnum, err := strconv.Atoi(line)
			if err != nil {
				// Skip this shop
				skipUntilNextShop = true
				continue
			}

			currentShop.RoomVNUM = roomVnum

			// Read the shop messages
			for i := 0; i < 7; i++ {
				if !scanner.Scan() {
					// End of file
					break
				}
				lineNum++
				message := scanner.Text()
				currentShop.Messages = append(currentShop.Messages, message)
			}

			// We're done with this shop, skip until we find a new one
			skipUntilNextShop = true
		}
	}

	// Add the last shop if it exists and we're not skipping it
	if currentShop != nil && !skipUntilNextShop {
		shops = append(shops, currentShop)
	}

	return shops, nil
}
