package world

import (
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestMagicalEquipment(t *testing.T) {
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

	// Create a test character
	character := NewCharacter("TestCharacter", false)
	character.World = world
	character.InRoom = room
	character.HitRoll = 0
	character.DamRoll = 0
	character.ArmorClass = [3]int{10, 10, 10}
	character.Abilities = [6]int{10, 10, 10, 10, 10, 10}

	// Create a magical sword with +2 hitroll and +3 damroll
	swordProto := &types.Object{
		VNUM:        2001,
		Name:        "sword",
		ShortDesc:   "a magical sword",
		Description: "A magical sword is lying here.",
		Type:        types.ITEM_WEAPON,
		WearFlags:   types.ITEM_WEAR_TAKE | types.ITEM_WEAR_WIELD,
		ExtraFlags:  types.ITEM_MAGIC,
		Value:       [4]int{0, 1, 6, 0}, // Type, dice, sides, flags
		Affects: [types.MAX_OBJ_AFFECT]struct {
			Location int
			Modifier int
		}{
			{Location: types.APPLY_HITROLL, Modifier: 2},
			{Location: types.APPLY_DAMROLL, Modifier: 3},
		},
	}
	world.objects[2001] = swordProto

	// Create a magical armor with -10 AC
	armorProto := &types.Object{
		VNUM:        2002,
		Name:        "armor",
		ShortDesc:   "a magical armor",
		Description: "A magical armor is lying here.",
		Type:        types.ITEM_ARMOR,
		WearFlags:   types.ITEM_WEAR_TAKE | types.ITEM_WEAR_BODY,
		ExtraFlags:  types.ITEM_MAGIC,
		Value:       [4]int{10, 0, 0, 0}, // AC bonus
		Affects: [types.MAX_OBJ_AFFECT]struct {
			Location int
			Modifier int
		}{
			{Location: types.APPLY_AC, Modifier: -10},
			{Location: types.APPLY_STR, Modifier: 1},
		},
	}
	world.objects[2002] = armorProto

	// Create a magical ring with +1 to all abilities
	ringProto := &types.Object{
		VNUM:        2003,
		Name:        "ring",
		ShortDesc:   "a magical ring",
		Description: "A magical ring is lying here.",
		Type:        types.ITEM_ARMOR,
		WearFlags:   types.ITEM_WEAR_TAKE | types.ITEM_WEAR_FINGER,
		ExtraFlags:  types.ITEM_MAGIC,
		Affects: [types.MAX_OBJ_AFFECT]struct {
			Location int
			Modifier int
		}{
			{Location: types.APPLY_STR, Modifier: 1},
			{Location: types.APPLY_DEX, Modifier: 1},
		},
	}
	world.objects[2003] = ringProto

	// Create object instances
	sword := world.CreateObjectFromPrototype(2001)
	armor := world.CreateObjectFromPrototype(2002)
	ring := world.CreateObjectFromPrototype(2003)

	// Test equipping the sword
	sword.WornBy = character
	sword.WornOn = types.WEAR_WIELD
	character.Equipment[types.WEAR_WIELD] = sword

	// Apply the sword's affects
	world.ApplyObjectAffects(character, sword, true)

	// Check if the character's hitroll and damroll were increased
	if character.HitRoll != 2 {
		t.Errorf("Expected hitroll to be 2, got %d", character.HitRoll)
	}
	if character.DamRoll != 3 {
		t.Errorf("Expected damroll to be 3, got %d", character.DamRoll)
	}

	// Test equipping the armor
	armor.WornBy = character
	armor.WornOn = types.WEAR_BODY
	character.Equipment[types.WEAR_BODY] = armor

	// Apply the armor's affects
	world.ApplyObjectAffects(character, armor, true)

	// Check if the character's AC was decreased and strength increased
	if character.ArmorClass[0] != 0 {
		t.Errorf("Expected AC to be 0, got %d", character.ArmorClass[0])
	}
	if character.Abilities[types.ABILITY_STR] != 11 {
		t.Errorf("Expected strength to be 11, got %d", character.Abilities[types.ABILITY_STR])
	}

	// Test equipping the ring
	ring.WornBy = character
	ring.WornOn = types.WEAR_FINGER_R
	character.Equipment[types.WEAR_FINGER_R] = ring

	// Apply the ring's affects
	world.ApplyObjectAffects(character, ring, true)

	// Check if the character's abilities were increased
	if character.Abilities[types.ABILITY_STR] != 12 {
		t.Errorf("Expected strength to be 12, got %d", character.Abilities[types.ABILITY_STR])
	}
	if character.Abilities[types.ABILITY_DEX] != 11 {
		t.Errorf("Expected dexterity to be 11, got %d", character.Abilities[types.ABILITY_DEX])
	}

	// Test removing the sword
	world.ApplyObjectAffects(character, sword, false)
	character.Equipment[types.WEAR_WIELD] = nil
	sword.WornBy = nil
	sword.WornOn = -1

	// Check if the character's hitroll and damroll were decreased
	if character.HitRoll != 0 {
		t.Errorf("Expected hitroll to be 0, got %d", character.HitRoll)
	}
	if character.DamRoll != 0 {
		t.Errorf("Expected damroll to be 0, got %d", character.DamRoll)
	}

	// Test removing the armor
	world.ApplyObjectAffects(character, armor, false)
	character.Equipment[types.WEAR_BODY] = nil
	armor.WornBy = nil
	armor.WornOn = -1

	// Check if the character's AC was increased and strength decreased
	if character.ArmorClass[0] != 10 {
		t.Errorf("Expected AC to be 10, got %d", character.ArmorClass[0])
	}
	if character.Abilities[types.ABILITY_STR] != 11 {
		t.Errorf("Expected strength to be 11, got %d", character.Abilities[types.ABILITY_STR])
	}

	// Test removing the ring
	world.ApplyObjectAffects(character, ring, false)
	character.Equipment[types.WEAR_FINGER_R] = nil
	ring.WornBy = nil
	ring.WornOn = -1

	// Check if the character's abilities were decreased
	if character.Abilities[types.ABILITY_STR] != 10 {
		t.Errorf("Expected strength to be 10, got %d", character.Abilities[types.ABILITY_STR])
	}
	if character.Abilities[types.ABILITY_DEX] != 10 {
		t.Errorf("Expected dexterity to be 10, got %d", character.Abilities[types.ABILITY_DEX])
	}
}
