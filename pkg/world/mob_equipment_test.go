package world

import (
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestMobEquipment(t *testing.T) {
	// Create a mock storage
	storage := NewMockStorage()

	// Create a test room
	room := &types.Room{
		VNUM:        3001,
		Name:        "Test Room",
		Description: "This is a test room.",
		Characters:  make([]*types.Character, 0),
	}
	storage.rooms = append(storage.rooms, room)

	// Create a test world
	world, err := NewWorld(nil, storage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Create a test object prototype
	swordProto := &types.Object{
		VNUM:        2001,
		Name:        "sword",
		ShortDesc:   "a test sword",
		Description: "A test sword is lying here.",
		Type:        types.ITEM_WEAPON,
		WearFlags:   types.ITEM_WEAR_TAKE | types.ITEM_WEAR_WIELD,
		Value:       [4]int{0, 1, 6, 0}, // Type, dice, sides, flags
	}
	world.objects[2001] = swordProto

	// Create a test armor prototype
	armorProto := &types.Object{
		VNUM:        2002,
		Name:        "armor",
		ShortDesc:   "a test armor",
		Description: "A test armor is lying here.",
		Type:        types.ITEM_ARMOR,
		WearFlags:   types.ITEM_WEAR_TAKE | types.ITEM_WEAR_BODY,
		Value:       [4]int{5, 0, 0, 0}, // AC bonus
	}
	world.objects[2002] = armorProto

	// Create a test mobile prototype with equipment
	mobProto := &types.Mobile{
		VNUM:        1001,
		Name:        "testmob",
		ShortDesc:   "a test mob",
		LongDesc:    "A test mob is standing here.",
		Description: "This is a test mob for unit testing.",
		ActFlags:    types.ACT_ISNPC,
		Level:       5,
		Position:    types.POS_STANDING,
		Gold:        100,
		Experience:  500,
		AC:          [3]int{5, 5, 5},
		HitRoll:     2,
		DamRoll:     2,
		Dice:        [3]int{2, 8, 5}, // 2d8+5 hit points
		Equipment: []types.MobEquipment{
			{
				ObjectVNUM: 2001,
				Position:   types.WEAR_WIELD,
				Chance:     100,
			},
			{
				ObjectVNUM: 2002,
				Position:   types.WEAR_BODY,
				Chance:     100,
			},
		},
	}
	world.mobiles[1001] = mobProto

	// Create a mob from the prototype
	mob := world.CreateMobFromPrototype(1001, room)
	if mob == nil {
		t.Fatalf("Failed to create mob from prototype")
	}

	// Check if the mob has the sword equipped
	if mob.Equipment[types.WEAR_WIELD] == nil {
		t.Errorf("Expected mob to have a sword equipped in WEAR_WIELD position")
	} else {
		sword := mob.Equipment[types.WEAR_WIELD]
		if sword.Prototype.VNUM != 2001 {
			t.Errorf("Expected mob to have sword with VNUM 2001, got %d", sword.Prototype.VNUM)
		}
	}

	// Check if the mob has the armor equipped
	if mob.Equipment[types.WEAR_BODY] == nil {
		t.Errorf("Expected mob to have armor equipped in WEAR_BODY position")
	} else {
		armor := mob.Equipment[types.WEAR_BODY]
		if armor.Prototype.VNUM != 2002 {
			t.Errorf("Expected mob to have armor with VNUM 2002, got %d", armor.Prototype.VNUM)
		}
	}

	// Test partial equipment chance
	// Create a new mob prototype with partial equipment chance
	mobProto2 := &types.Mobile{
		VNUM:        1002,
		Name:        "testmob2",
		ShortDesc:   "a test mob 2",
		LongDesc:    "A test mob 2 is standing here.",
		Description: "This is a test mob 2 for unit testing.",
		ActFlags:    types.ACT_ISNPC,
		Level:       5,
		Position:    types.POS_STANDING,
		Gold:        100,
		Experience:  500,
		AC:          [3]int{5, 5, 5},
		HitRoll:     2,
		DamRoll:     2,
		Dice:        [3]int{2, 8, 5}, // 2d8+5 hit points
		Equipment: []types.MobEquipment{
			{
				ObjectVNUM: 2001,
				Position:   types.WEAR_WIELD,
				Chance:     0, // 0% chance to equip
			},
		},
	}
	world.mobiles[1002] = mobProto2

	// Create a mob from the second prototype
	mob2 := world.CreateMobFromPrototype(1002, room)
	if mob2 == nil {
		t.Fatalf("Failed to create mob2 from prototype")
	}

	// Check that the mob doesn't have the sword equipped (0% chance)
	if mob2.Equipment[types.WEAR_WIELD] != nil {
		t.Errorf("Expected mob2 to not have a sword equipped (0%% chance)")
	}
}
