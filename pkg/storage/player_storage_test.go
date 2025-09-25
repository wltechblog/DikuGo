package storage

import (
	"os"
	"testing"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestSavePlayerWithCircularReferences(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dikugo_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file player storage
	storage, err := NewFilePlayerStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create a test character with equipment and inventory
	player := &types.Character{
		Name:          "TestPlayer",
		Level:         5,
		Gold:          100,
		HP:            50,
		MaxHitPoints:  50,
		Equipment:     make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory:     make([]*types.ObjectInstance, 0),
		Conditions:    [3]int{10, 15, 0},
		Skills:        make(map[int]int),
		Spells:        make(map[int]int),
		LastSkillTime: make(map[int]time.Time),
	}

	// Create a sword with circular references
	sword := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1001,
			Name:      "sword",
			ShortDesc: "a steel sword",
			Type:      types.ITEM_WEAPON,
		},
		CarriedBy: player, // Circular reference
		WornBy:    player, // Circular reference
		WornOn:    types.WEAR_WIELD,
		Value:     [4]int{0, 0, 0, 0},
	}

	// Create a container with items
	bag := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1002,
			Name:      "bag",
			ShortDesc: "a leather bag",
			Type:      types.ITEM_CONTAINER,
		},
		CarriedBy: player, // Circular reference
		Contains:  make([]*types.ObjectInstance, 0),
		Value:     [4]int{100, 0, 0, 0},
	}

	// Create an item inside the bag
	potion := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1003,
			Name:      "potion",
			ShortDesc: "a healing potion",
			Type:      types.ITEM_POTION,
		},
		InObj: bag, // Circular reference to container
		Value: [4]int{50, 1, 0, 0},
	}

	// Set up the circular references
	bag.Contains = append(bag.Contains, potion)
	player.Equipment[types.WEAR_WIELD] = sword
	player.Inventory = append(player.Inventory, bag)

	// Test saving the player (this should not cause a circular reference error)
	err = storage.SavePlayer(player)
	if err != nil {
		t.Fatalf("Failed to save player with circular references: %v", err)
	}

	// Test loading the player back
	loadedPlayer, err := storage.LoadPlayer("TestPlayer")
	if err != nil {
		t.Fatalf("Failed to load player: %v", err)
	}

	// Verify the player data was preserved
	if loadedPlayer.Name != "TestPlayer" {
		t.Errorf("Expected name 'TestPlayer', got '%s'", loadedPlayer.Name)
	}

	if loadedPlayer.Level != 5 {
		t.Errorf("Expected level 5, got %d", loadedPlayer.Level)
	}

	if loadedPlayer.Gold != 100 {
		t.Errorf("Expected gold 100, got %d", loadedPlayer.Gold)
	}

	// Verify equipment was restored
	if loadedPlayer.Equipment[types.WEAR_WIELD] == nil {
		t.Error("Expected sword in WEAR_WIELD slot, got nil")
	} else {
		loadedSword := loadedPlayer.Equipment[types.WEAR_WIELD]
		if loadedSword.Prototype.Name != "sword" {
			t.Errorf("Expected sword name 'sword', got '%s'", loadedSword.Prototype.Name)
		}
		// Verify circular references were restored
		if loadedSword.CarriedBy != loadedPlayer {
			t.Error("Expected sword.CarriedBy to point to player")
		}
		if loadedSword.WornBy != loadedPlayer {
			t.Error("Expected sword.WornBy to point to player")
		}
		if loadedSword.WornOn != types.WEAR_WIELD {
			t.Errorf("Expected sword.WornOn to be %d, got %d", types.WEAR_WIELD, loadedSword.WornOn)
		}
	}

	// Verify inventory was restored
	if len(loadedPlayer.Inventory) != 1 {
		t.Errorf("Expected 1 item in inventory, got %d", len(loadedPlayer.Inventory))
	} else {
		loadedBag := loadedPlayer.Inventory[0]
		if loadedBag.Prototype.Name != "bag" {
			t.Errorf("Expected bag name 'bag', got '%s'", loadedBag.Prototype.Name)
		}
		// Verify circular references were restored
		if loadedBag.CarriedBy != loadedPlayer {
			t.Error("Expected bag.CarriedBy to point to player")
		}

		// Verify container contents were restored
		if len(loadedBag.Contains) != 1 {
			t.Errorf("Expected 1 item in bag, got %d", len(loadedBag.Contains))
		} else {
			loadedPotion := loadedBag.Contains[0]
			if loadedPotion.Prototype.Name != "potion" {
				t.Errorf("Expected potion name 'potion', got '%s'", loadedPotion.Prototype.Name)
			}
			// Verify circular references were restored
			if loadedPotion.InObj != loadedBag {
				t.Error("Expected potion.InObj to point to bag")
			}
			if loadedPotion.CarriedBy != nil {
				t.Error("Expected potion.CarriedBy to be nil (items in containers are not carried directly)")
			}
		}
	}
}

func TestSavePlayerWithoutCircularReferences(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dikugo_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file player storage
	storage, err := NewFilePlayerStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create a simple test character without circular references
	player := &types.Character{
		Name:          "SimplePlayer",
		Level:         1,
		Gold:          50,
		HP:            20,
		MaxHitPoints:  20,
		Equipment:     make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory:     make([]*types.ObjectInstance, 0),
		Conditions:    [3]int{0, 0, 0},
		Skills:        make(map[int]int),
		Spells:        make(map[int]int),
		LastSkillTime: make(map[int]time.Time),
	}

	// Test saving the simple player
	err = storage.SavePlayer(player)
	if err != nil {
		t.Fatalf("Failed to save simple player: %v", err)
	}

	// Test loading the player back
	loadedPlayer, err := storage.LoadPlayer("SimplePlayer")
	if err != nil {
		t.Fatalf("Failed to load simple player: %v", err)
	}

	// Verify the player data was preserved
	if loadedPlayer.Name != "SimplePlayer" {
		t.Errorf("Expected name 'SimplePlayer', got '%s'", loadedPlayer.Name)
	}

	if loadedPlayer.Level != 1 {
		t.Errorf("Expected level 1, got %d", loadedPlayer.Level)
	}

	if loadedPlayer.Gold != 50 {
		t.Errorf("Expected gold 50, got %d", loadedPlayer.Gold)
	}
}
