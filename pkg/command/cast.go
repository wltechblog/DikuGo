package command

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// DoCast handles the cast command
func DoCast(character *types.Character, arguments string, world *world.World) error {
	// Check if character is a player
	if character.IsNPC {
		return nil
	}

	// Check if character is a mage or cleric
	if character.Class != types.CLASS_MAGIC_USER && character.Class != types.CLASS_CLERIC && character.Level < 21 {
		if character.Class == types.CLASS_WARRIOR {
			return fmt.Errorf("think you had better stick to fighting...")
		} else if character.Class == types.CLASS_THIEF {
			return fmt.Errorf("think you should stick to robbing and killing...")
		}
		return fmt.Errorf("you are not trained in the magical arts")
	}

	// Check if arguments are empty
	if arguments == "" {
		return fmt.Errorf("cast which what where?")
	}

	// Check if spell is enclosed in quotes
	if !strings.HasPrefix(arguments, "'") {
		return fmt.Errorf("magic must always be enclosed by the holy magic symbols: '")
	}

	// Find the closing quote
	endQuote := strings.Index(arguments[1:], "'")
	if endQuote == -1 {
		return fmt.Errorf("magic must always be enclosed by the holy magic symbols: '")
	}
	endQuote++ // Adjust for the offset from the first character

	// Extract the spell name
	spellName := strings.ToLower(arguments[1:endQuote])

	// Find the spell
	spellID := types.GetSpellByName(spellName)
	if spellID == types.SPELL_UNDEFINED {
		return fmt.Errorf("your lips move, but no magic appears")
	}

	// Check character's position
	minPosition := types.GetSpellPosition(spellID)
	if character.Position < minPosition {
		switch character.Position {
		case types.POS_SLEEPING:
			return fmt.Errorf("you dream about great magical powers")
		case types.POS_RESTING:
			return fmt.Errorf("you can't concentrate enough while resting")
		case types.POS_SITTING:
			return fmt.Errorf("you can't do this sitting!")
		case types.POS_FIGHTING:
			return fmt.Errorf("impossible! You can't concentrate enough!")
		default:
			return fmt.Errorf("it seems like you're in a pretty bad shape!")
		}
	}

	// Check if character knows the spell
	if character.Level < 21 {
		minLevel := types.GetSpellMinLevel(spellID, character.Class)
		if minLevel > character.Level {
			return fmt.Errorf("you are not powerful enough to cast that spell")
		}
	}

	// Check if character has enough mana
	manaCost := types.GetSpellMana(spellID)
	if character.ManaPoints < manaCost {
		return fmt.Errorf("you can't summon enough energy to cast the spell")
	}

	// Extract the target argument
	targetArg := ""
	if len(arguments) > endQuote+1 {
		targetArg = strings.TrimSpace(arguments[endQuote+1:])
	}

	// Find the target
	victim, obj, err := world.GetSpellTarget(character, targetArg, spellID)
	if err != nil {
		return err
	}

	// Say the spell words
	world.SaySpell(character, spellID)

	// Add delay
	delay := types.GetSpellDelay(spellID)
	world.AddDelay(character, delay)

	// Check for spell failure
	if character.Level < 21 {
		skillLevel := 0
		if level, ok := character.Spells[spellID]; ok {
			skillLevel = level
		}

		// Base 50% chance + 2% per level
		baseChance := 50 + (character.Level * 2)

		// Adjust based on skill level
		chance := baseChance + skillLevel

		// Cap at 95%
		if chance > 95 {
			chance = 95
		}

		// Roll for success
		if rand.Intn(100) >= chance {
			character.SendMessage("You lost your concentration!\r\n")
			character.ManaPoints -= (manaCost / 2)
			return nil
		}
	}

	// Cast the spell
	err = castSpell(world, spellID, character.Level, character, targetArg, types.SPELL_TYPE_SPELL, victim, obj)
	if err != nil {
		return err
	}

	// Deduct mana
	character.ManaPoints -= manaCost

	// Improve spell skill
	if character.Level < 21 && rand.Intn(100) > character.Spells[spellID] {
		character.Spells[spellID]++
		character.SendMessage(fmt.Sprintf("You feel more confident in your %s spell.\r\n", types.GetSpellName(spellID)))
	}

	return nil
}

// castSpell dispatches to the appropriate spell function
func castSpell(w *world.World, spellID, level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	switch spellID {
	case types.SPELL_ARMOR:
		return w.CastArmor(level, ch, arg, spellType, victim, obj)
	case types.SPELL_TELEPORT:
		return w.CastWordOfRecall(level, ch, arg, spellType, victim, obj) // Using word of recall for now
	case types.SPELL_BLESS:
		return w.CastBless(level, ch, arg, spellType, victim, obj)
	case types.SPELL_BLINDNESS:
		return w.CastBlindness(level, ch, arg, spellType, victim, obj)
	case types.SPELL_BURNING_HANDS:
		return w.CastBurningHands(level, ch, arg, spellType, victim, obj)
	case types.SPELL_CALL_LIGHTNING:
		return w.CastCallLightning(level, ch, arg, spellType, victim, obj)
	case types.SPELL_CHARM_PERSON:
		return w.CastCharmPerson(level, ch, arg, spellType, victim, obj)
	case types.SPELL_CHILL_TOUCH:
		return w.CastChillTouch(level, ch, arg, spellType, victim, obj)
	case types.SPELL_COLOR_SPRAY:
		return w.CastColorSpray(level, ch, arg, spellType, victim, obj)
	case types.SPELL_CONTROL_WEATHER:
		return w.CastControlWeather(level, ch, arg, spellType, victim, obj)
	case types.SPELL_CREATE_FOOD:
		return w.CastCreateFood(level, ch, arg, spellType, victim, obj)
	case types.SPELL_CREATE_WATER:
		return w.CastCreateWater(level, ch, arg, spellType, victim, obj)
	case types.SPELL_CURE_BLIND:
		return w.CastCureBlindness(level, ch, arg, spellType, victim, obj)
	case types.SPELL_CURE_CRITIC:
		return w.CastCureCritic(level, ch, arg, spellType, victim, obj)
	case types.SPELL_CURE_LIGHT:
		return w.CastCureLight(level, ch, arg, spellType, victim, obj)
	case types.SPELL_CURSE:
		return w.CastCurse(level, ch, arg, spellType, victim, obj)
	case types.SPELL_DETECT_EVIL:
		return w.CastDetectEvil(level, ch, arg, spellType, victim, obj)
	case types.SPELL_DETECT_INVISIBLE:
		return w.CastDetectInvisible(level, ch, arg, spellType, victim, obj)
	case types.SPELL_DETECT_MAGIC:
		return w.CastDetectMagic(level, ch, arg, spellType, victim, obj)
	case types.SPELL_DETECT_POISON:
		return w.CastDetectPoison(level, ch, arg, spellType, victim, obj)
	case types.SPELL_DISPEL_EVIL:
		return w.CastDispelEvil(level, ch, arg, spellType, victim, obj)
	case types.SPELL_EARTHQUAKE:
		return w.CastEarthquake(level, ch, arg, spellType, victim, obj)
	case types.SPELL_ENCHANT_WEAPON:
		return w.CastEnchantWeapon(level, ch, arg, spellType, victim, obj)
	case types.SPELL_ENERGY_DRAIN:
		return w.CastEnergyDrain(level, ch, arg, spellType, victim, obj)
	case types.SPELL_FIREBALL:
		return w.CastFireball(level, ch, arg, spellType, victim, obj)
	case types.SPELL_HARM:
		return w.CastHarm(level, ch, arg, spellType, victim, obj)
	case types.SPELL_HEAL:
		return w.CastHeal(level, ch, arg, spellType, victim, obj)
	case types.SPELL_INVISIBLE:
		return w.CastInvisible(level, ch, arg, spellType, victim, obj)
	case types.SPELL_LIGHTNING_BOLT:
		return w.CastLightningBolt(level, ch, arg, spellType, victim, obj)
	case types.SPELL_LOCATE_OBJECT:
		return w.CastLocateObject(level, ch, arg, spellType, victim, obj)
	case types.SPELL_MAGIC_MISSILE:
		return w.CastMagicMissile(level, ch, arg, spellType, victim, obj)
	case types.SPELL_POISON:
		return w.CastPoison(level, ch, arg, spellType, victim, obj)
	case types.SPELL_PROTECTION_FROM_EVIL:
		return w.CastProtectionFromEvil(level, ch, arg, spellType, victim, obj)
	case types.SPELL_REMOVE_CURSE:
		return w.CastRemoveCurse(level, ch, arg, spellType, victim, obj)
	case types.SPELL_SANCTUARY:
		return w.CastSanctuary(level, ch, arg, spellType, victim, obj)
	case types.SPELL_SHOCKING_GRASP:
		return w.CastShockingGrasp(level, ch, arg, spellType, victim, obj)
	case types.SPELL_SLEEP:
		return w.CastSleep(level, ch, arg, spellType, victim, obj)
	case types.SPELL_STRENGTH:
		return w.CastStrength(level, ch, arg, spellType, victim, obj)
	case types.SPELL_SUMMON:
		return w.CastSummon(level, ch, arg, spellType, victim, obj)
	case types.SPELL_WORD_OF_RECALL:
		return w.CastWordOfRecall(level, ch, arg, spellType, victim, obj)
	case types.SPELL_REMOVE_POISON:
		return w.CastRemovePoison(level, ch, arg, spellType, victim, obj)
	case types.SPELL_SENSE_LIFE:
		return w.CastSenseLife(level, ch, arg, spellType, victim, obj)
	default:
		return fmt.Errorf("sorry, this magic has not yet been implemented")
	}
}
