package storage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ParseShops parses the shop file and returns a slice of shops
func ParseShops(filename string) ([]*types.Shop, error) {
	// Debug: Print the filename
	log.Printf("Parsing shop file: %s", filename)
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
	var producingItems []int
	var buyTypes []int
	var messages []string
	var shopNum int
	var roomVnum int
	var mobileVnum int
	var profitBuy float64
	var profitSell float64
	var openHour int
	var closeHour int
	var itemCount int
	var messageCount int

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

			// Parse shop number
			vnumStr := strings.TrimPrefix(line, "#")
			vnumStr = strings.TrimSuffix(vnumStr, "~")
			log.Printf("Parsing shop with VNUM string: '%s'", vnumStr)
			var err error
			shopNum, err = strconv.Atoi(vnumStr)
			if err != nil {
				return nil, fmt.Errorf("invalid shop number on line %d: %w", lineNum, err)
			}

			// Reset variables for the new shop
			producingItems = make([]int, 0)
			buyTypes = make([]int, 0)
			messages = make([]string, 0)
			roomVnum = 0
			mobileVnum = 0
			profitBuy = 0
			profitSell = 0
			openHour = 0
			closeHour = 0
			itemCount = 0
			messageCount = 0

			// Set the state to read producing items next
			state = "producing"
			continue
		}

		// Process the line based on the current state
		switch state {
		case "producing":
			// Parse producing items (5 items)
			item, err := strconv.Atoi(line)
			if err != nil {
				log.Printf("Warning: Invalid producing item on line %d: %s", lineNum, line)
			} else {
				if item != -1 {
					producingItems = append(producingItems, item)
				}
			}

			itemCount++
			if itemCount >= 5 {
				state = "profitbuy"
				itemCount = 0
			}

		case "profitbuy":
			// Parse profit buy
			profitBuy, err = strconv.ParseFloat(line, 64)
			if err != nil {
				log.Printf("Warning: Invalid profit buy on line %d: %s", lineNum, line)
				profitBuy = 1.0 // Default value
			}
			state = "profitsell"

		case "profitsell":
			// Parse profit sell
			profitSell, err = strconv.ParseFloat(line, 64)
			if err != nil {
				log.Printf("Warning: Invalid profit sell on line %d: %s", lineNum, line)
				profitSell = 1.0 // Default value
			}
			state = "buytypes"

		case "buytypes":
			// Parse buy types (5 types)
			buyType, err := strconv.Atoi(line)
			if err != nil {
				log.Printf("Warning: Invalid buy type on line %d: %s", lineNum, line)
			} else {
				if buyType != -1 {
					buyTypes = append(buyTypes, buyType)
				}
			}

			itemCount++
			if itemCount >= 5 {
				state = "messages"
				itemCount = 0
			}

		case "messages":
			// Parse messages (7 messages)
			messages = append(messages, line)
			messageCount++
			if messageCount >= 7 {
				state = "temper1"
			}

		case "temper1":
			// Skip temper1
			state = "temper2"

		case "temper2":
			// Skip temper2
			state = "keeper"

		case "keeper":
			// Parse keeper (mobile VNUM)
			mobileVnum, err = strconv.Atoi(line)
			if err != nil {
				log.Printf("Warning: Invalid keeper VNUM on line %d: %s", lineNum, line)
				mobileVnum = 0 // Default value
			}
			state = "withwho"

		case "withwho":
			// Skip withwho
			state = "roomvnum"

		case "roomvnum":
			// Parse room VNUM
			roomVnum, err = strconv.Atoi(line)
			if err != nil {
				log.Printf("Warning: Invalid room VNUM on line %d: %s", lineNum, line)
				roomVnum = 0 // Default value
			}
			state = "openhour1"

		case "openhour1":
			// Parse open hour 1
			openHour, err = strconv.Atoi(line)
			if err != nil {
				log.Printf("Warning: Invalid open hour on line %d: %s", lineNum, line)
				openHour = 0 // Default value
			}
			state = "closehour1"

		case "closehour1":
			// Parse close hour 1
			closeHour, err = strconv.Atoi(line)
			if err != nil {
				log.Printf("Warning: Invalid close hour on line %d: %s", lineNum, line)
				closeHour = 0 // Default value
			}
			state = "openhour2"

		case "openhour2":
			// Skip open hour 2
			state = "closehour2"

		case "closehour2":
			// Skip close hour 2
			// Create the shop now that we have all the data
			currentShop = &types.Shop{
				VNUM:       shopNum,
				RoomVNUM:   roomVnum,
				MobileVNUM: mobileVnum,
				ProfitBuy:  profitBuy,
				ProfitSell: profitSell,
				OpenHour:   openHour,
				CloseHour:  closeHour,
				Producing:  make([]int, len(producingItems)),
				BuyTypes:   make([]int, len(buyTypes)),
				Messages:   make([]string, len(messages)),
			}

			// Copy the slices to avoid reference issues
			copy(currentShop.Producing, producingItems)
			copy(currentShop.BuyTypes, buyTypes)
			copy(currentShop.Messages, messages)

			log.Printf("Created shop #%d: Room VNUM = %d, Keeper VNUM = %d, Items: %v",
				currentShop.VNUM, currentShop.RoomVNUM, currentShop.MobileVNUM, currentShop.Producing)

			// Wait for the next shop
			state = "waitfornext"

		case "waitfornext":
			// Do nothing, wait for the next shop
		}
	}

	// Add the last shop if it exists
	if currentShop != nil && state == "waitfornext" {
		shops = append(shops, currentShop)
	}

	// Debug: Print all parsed shops
	log.Printf("Parsed %d shops from file", len(shops))
	for _, shop := range shops {
		log.Printf("Parsed shop #%d: Room VNUM = %d, Keeper VNUM = %d, Items: %v",
			shop.VNUM, shop.RoomVNUM, shop.MobileVNUM, shop.Producing)
	}

	return shops, nil
}
