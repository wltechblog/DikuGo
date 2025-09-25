# Enhanced Equipment System

## Overview

The DikuGo equipment system has been enhanced to provide better handling of equipment commands, following the original DikuMUD mechanics more closely while adding quality-of-life improvements.

## New Commands

### 🤲 **Hold Command**

The `hold` command allows players to hold items in their hands, with special handling for light sources.

**Usage:**
```
hold <item>
grab <item>  (alias)
```

**Features:**
- **Light Sources**: Automatically detects `ITEM_LIGHT` objects and places them in the `WEAR_LIGHT` position
- **Holdable Items**: Places items with `ITEM_WEAR_HOLD` flag in the `WEAR_HOLD` position
- **Class Restrictions**: Respects class-based item restrictions
- **Alignment Checks**: Handles anti-alignment items (zapping effect)

**Examples:**
```
> hold torch
Ok.
(Torch goes to WEAR_LIGHT position)

> hold wand
Ok.
(Wand goes to WEAR_HOLD position)

> hold sword
You can't hold a long sword.
(Swords can't be held, only wielded)
```

### ⚔️ **Wield Command**

The `wield` command is specifically for wielding weapons, separate from the general wear command.

**Usage:**
```
wield <weapon>
```

**Features:**
- **Weapon-Specific**: Only works with items that have `ITEM_WEAR_WIELD` flag
- **Class Restrictions**: Respects weapon restrictions (e.g., clerics can't use edged weapons)
- **Strength Checks**: Future enhancement for weapon weight requirements
- **Alignment Checks**: Handles anti-alignment weapons

**Examples:**
```
> wield sword
Ok.

> wield dagger
Ok.

> wield shield
You can't wield a wooden shield.
(Shields are worn, not wielded)
```

## Enhanced Wear Command

### 🔄 **Multiple Position Support**

The wear command now intelligently handles items that can be worn in multiple positions, following the original DikuMUD behavior.

**Finger Items:**
- First ring goes to left finger (`WEAR_FINGER_L`)
- Second ring goes to right finger (`WEAR_FINGER_R`)
- Provides specific messages: "You put the ring on your left finger."

**Wrist Items:**
- First bracelet goes to left wrist (`WEAR_WRIST_L`)
- Second bracelet goes to right wrist (`WEAR_WRIST_R`)
- Provides specific messages: "You wear the bracelet around your left wrist."

**Neck Items:**
- First necklace goes to first neck position (`WEAR_NECK_1`)
- Second necklace goes to second neck position (`WEAR_NECK_2`)
- Uses simple "Ok." message like original DikuMUD

### 📦 **Enhanced "Wear All"**

The `wear all` command has been improved to iterate through inventory more intelligently:

**Features:**
- **Smart Position Selection**: Uses `findBestWearPosition()` to find optimal positions
- **Multiple Position Handling**: Properly handles rings, bracelets, and necklaces
- **Class Restrictions**: Skips items the character can't use
- **Alignment Checks**: Skips items that would zap the character
- **Detailed Feedback**: Shows exactly what was worn and where

**Example:**
```
> wear all
You put the gold ring on your left finger.
You put the silver ring on your right finger.
You wear the leather bracelet around your left wrist.
You wear the chain mail on your body.
You wear the iron helmet on your head.
```

## Position-Specific Messages

### 🎯 **Authentic DikuMUD Messages**

The system now provides position-specific messages that match the original DikuMUD:

| Position | Message Format |
|----------|----------------|
| Left Finger | "You put [item] on your left finger." |
| Right Finger | "You put [item] on your right finger." |
| Left Wrist | "You wear [item] around your left wrist." |
| Right Wrist | "You wear [item] around your right wrist." |
| Neck | "Ok." |
| Shield | "You start using [item]." |
| Other | "You wear [item] on your [position]." |

## Technical Implementation

### 🔧 **Core Functions**

**`findBestWearPosition(character, wearFlags)`**
- Intelligently finds the best available position for an item
- Handles multiple positions (fingers, wrists, neck)
- Checks if positions are already occupied
- Returns -1 if no suitable position is available

**`getWearMessage(obj, position)`**
- Returns position-specific wear messages
- Matches original DikuMUD message formats
- Handles special cases like shields and jewelry

**Legacy Compatibility**
- `findWearPosition()` function maintained for backward compatibility
- Existing code continues to work without modification

### 🧪 **Comprehensive Testing**

**Test Coverage:**
- Light source holding (torches, lanterns)
- Holdable item handling (wands, staves)
- Weapon wielding with class restrictions
- Multiple finger ring wearing
- Multiple wrist bracelet wearing
- Position-specific message validation
- Error condition handling

## Equipment Positions

### 📍 **Available Positions**

| Position | Constant | Description | Multiple |
|----------|----------|-------------|----------|
| Light | `WEAR_LIGHT` | Light sources | No |
| Left Finger | `WEAR_FINGER_L` | Rings, jewelry | Yes (2 total) |
| Right Finger | `WEAR_FINGER_R` | Rings, jewelry | Yes (2 total) |
| Neck 1 | `WEAR_NECK_1` | Necklaces, amulets | Yes (2 total) |
| Neck 2 | `WEAR_NECK_2` | Necklaces, amulets | Yes (2 total) |
| Body | `WEAR_BODY` | Armor, clothing | No |
| Head | `WEAR_HEAD` | Helmets, hats | No |
| Legs | `WEAR_LEGS` | Pants, leggings | No |
| Feet | `WEAR_FEET` | Boots, shoes | No |
| Hands | `WEAR_HANDS` | Gloves, gauntlets | No |
| Arms | `WEAR_ARMS` | Sleeves, bracers | No |
| Shield | `WEAR_SHIELD` | Shields | No |
| About | `WEAR_ABOUT` | Cloaks, capes | No |
| Waist | `WEAR_WAIST` | Belts, sashes | No |
| Left Wrist | `WEAR_WRIST_L` | Bracelets, bands | Yes (2 total) |
| Right Wrist | `WEAR_WRIST_R` | Bracelets, bands | Yes (2 total) |
| Wield | `WEAR_WIELD` | Weapons | No |
| Hold | `WEAR_HOLD` | Held items | No |

## Usage Examples

### 🎮 **Complete Equipment Session**

```
> inventory
You are carrying:
  a torch
  a magic wand
  a long sword
  a wooden shield
  a gold ring
  a silver ring
  a leather bracelet
  a chain bracelet

> hold torch
Ok.

> wield sword
Ok.

> wear shield
You start using a wooden shield.

> wear all
You put the gold ring on your left finger.
You put the silver ring on your right finger.
You wear the leather bracelet around your left wrist.
You wear the chain bracelet around your right wrist.

> equipment
You are using:
  light source        a torch
  right finger        a silver ring
  left finger         a gold ring
  shield              a wooden shield
  left wrist          a leather bracelet
  right wrist         a chain bracelet
  wielded             a long sword
  held                a magic wand
```

## Benefits

### ✨ **Player Experience**

1. **Intuitive Commands**: Separate `hold` and `wield` commands for different item types
2. **Smart Positioning**: Automatic selection of best available positions
3. **Clear Feedback**: Position-specific messages show exactly where items are worn
4. **Efficient Workflow**: `wear all` handles multiple items intelligently

### 🔧 **Technical Benefits**

1. **Original Compatibility**: Matches DikuMUD behavior exactly
2. **Extensible Design**: Easy to add new positions or item types
3. **Robust Testing**: Comprehensive test coverage ensures reliability
4. **Clean Code**: Well-structured functions with clear responsibilities

The enhanced equipment system provides a much more polished and authentic DikuMUD experience while maintaining full backward compatibility with existing code and data files.
