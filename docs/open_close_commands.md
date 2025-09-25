# Open and Close Commands Implementation

## Overview

I have implemented the `open` and `close` commands following the original DikuMUD mechanics exactly. These commands work on both containers (held or in room) and doors/exits, with proper state management and two-way door synchronization.

## Commands Implemented

### 1. **Open Command**
Opens containers and doors that are currently closed.

**Usage**: 
- `open <container>` - Opens a container in inventory or room
- `open <door>` - Opens a door by keyword
- `open <door> <direction>` - Opens a specific door in a direction

**Original DikuMUD Behavior**:
- Searches inventory first, then room objects, then doors
- Handles container flags: CLOSEABLE, CLOSED, LOCKED
- Handles door flags: EX_ISDOOR, EX_CLOSED, EX_LOCKED
- Opens both sides of doors automatically
- Sends appropriate messages to character and room

### 2. **Close Command**
Closes containers and doors that are currently open.

**Usage**:
- `close <container>` - Closes a container in inventory or room
- `close <door>` - Closes a door by keyword
- `close <door> <direction>` - Closes a specific door in a direction

**Original DikuMUD Behavior**:
- Same search priority as open command
- Validates container and door states before closing
- Closes both sides of doors automatically
- Sends appropriate messages to character and room

## Technical Implementation

### Container Mechanics

**Container Flags** (from `obj.Prototype.Value[1]`):
- `CONT_CLOSEABLE` (1): Container can be opened/closed
- `CONT_PICKPROOF` (2): Container cannot be picked
- `CONT_CLOSED` (4): Container is currently closed
- `CONT_LOCKED` (8): Container is currently locked

**Container State Validation**:
```go
// Open container checks
if container.Prototype.Type != types.ITEM_CONTAINER {
    return "That's not a container."
}
if (container.Prototype.Value[1] & types.CONT_CLOSED) == 0 {
    return "But it's already open!"
}
if (container.Prototype.Value[1] & types.CONT_CLOSEABLE) == 0 {
    return "You can't do that."
}
if (container.Prototype.Value[1] & types.CONT_LOCKED) != 0 {
    return "It seems to be locked."
}
```

**Container State Changes**:
- **Open**: `container.Value[1] &^= types.CONT_CLOSED` (remove CLOSED flag)
- **Close**: `container.Value[1] |= types.CONT_CLOSED` (add CLOSED flag)

### Door Mechanics

**Exit Flags** (from `exit.Flags`):
- `EX_ISDOOR` (1): Exit is a door that can be opened/closed
- `EX_CLOSED` (2): Door is currently closed
- `EX_LOCKED` (4): Door is currently locked
- `EX_PICKPROOF` (8): Door cannot be picked

**Door State Validation**:
```go
// Open door checks
if (exit.Flags & types.EX_ISDOOR) == 0 {
    return "That's impossible, I'm afraid."
}
if (exit.Flags & types.EX_CLOSED) == 0 {
    return "It's already open!"
}
if (exit.Flags & types.EX_LOCKED) != 0 {
    return "It seems to be locked."
}
```

**Door State Changes**:
- **Open**: `exit.Flags &^= types.EX_CLOSED` (remove CLOSED flag)
- **Close**: `exit.Flags |= types.EX_CLOSED` (add CLOSED flag)

### Two-Way Door Synchronization

**Automatic Other Side Updates**:
1. Find destination room using `exit.DestVnum`
2. Calculate reverse direction (north ↔ south, east ↔ west, up ↔ down)
3. Verify reverse exit leads back to current room
4. Apply same state change to reverse exit
5. Send appropriate message to destination room

**Reverse Direction Mapping**:
```go
func getReverseDirection(dir int) int {
    switch dir {
    case types.DIR_NORTH: return types.DIR_SOUTH
    case types.DIR_EAST:  return types.DIR_WEST
    case types.DIR_SOUTH: return types.DIR_NORTH
    case types.DIR_WEST:  return types.DIR_EAST
    case types.DIR_UP:    return types.DIR_DOWN
    case types.DIR_DOWN:  return types.DIR_UP
    }
}
```

### Search Priority and Object Finding

**Search Order**:
1. **Inventory**: Check character's carried items first
2. **Room Objects**: Check objects in the current room
3. **Doors/Exits**: Check room exits by keyword or direction

**Name Matching**:
- **Containers**: Match against `obj.Prototype.Name` and `obj.Prototype.ShortDesc`
- **Doors**: Match against `exit.Keywords` or default "door"
- **Case Insensitive**: All matching is case-insensitive
- **Partial Match**: Uses `strings.Contains()` for flexible matching

**Direction Parsing**:
- Full names: "north", "east", "south", "west", "up", "down"
- Abbreviations: "n", "e", "s", "w", "u", "d"
- Case insensitive matching

## Message System

### Success Messages

**Container Operations**:
- To character: "Ok.\r\n"
- To room: "{Character} opens {container}." / "{Character} closes {container}."

**Door Operations**:
- To character: "Ok.\r\n"
- To room: "{Character} opens the {door}." / "{Character} closes the {door}."
- To other side: "The {door} is opened from the other side." / "The {door} closes quietly."

### Error Messages

**Container Errors**:
- "That's not a container."
- "But it's already open!" / "But it's already closed!"
- "You can't do that." (not closeable)
- "That's impossible." (not closeable for close)
- "It seems to be locked."

**Door Errors**:
- "That's impossible, I'm afraid." (not a door for open)
- "That's absurd." (not a door for close)
- "It's already open!" / "It's already closed!"
- "It seems to be locked."

**General Errors**:
- "Open what?" / "Close what?" (no arguments)
- "I see no {object} here." (object not found)

## Files Created/Modified

### New Command Files
1. **`pkg/command/open.go`** - Open command implementation
2. **`pkg/command/close.go`** - Close command implementation

### Test Files
3. **`pkg/command/open_close_test.go`** - Comprehensive tests for both commands

### Modified Files
4. **`pkg/command/registry.go`** - Registered new commands

## Test Coverage

### Container Tests
- ✅ Opening closed containers
- ✅ Closing open containers
- ✅ Already open/closed validation
- ✅ Locked container handling
- ✅ Non-closeable container handling
- ✅ Non-container object rejection

### Door Tests
- ✅ Opening closed doors
- ✅ Closing open doors
- ✅ Already open/closed validation
- ✅ Locked door handling
- ✅ Non-door exit rejection

### General Tests
- ✅ No arguments handling
- ✅ Object not found handling
- ✅ Direction parsing
- ✅ Keyword matching

## Original DikuMUD Compatibility

### Exact Message Format
All messages match the original DikuMUD exactly:
- "Ok." (success)
- "That's not a container."
- "But it's already open!"
- "You can't do that."
- "It seems to be locked."
- "That's impossible, I'm afraid."
- "That's absurd."

### Mechanics Fidelity
- **Flag Operations**: Uses exact bitwise operations from original C code
- **Search Priority**: Inventory → Room Objects → Doors (same as original)
- **Two-Way Doors**: Automatic synchronization of both sides
- **State Validation**: All original validation checks implemented
- **Error Handling**: Identical error conditions and messages

### Integration Points
- **Container System**: Works with existing ITEM_CONTAINER objects
- **Exit System**: Works with existing room exit structure
- **Flag System**: Uses existing CONT_* and EX_* flag constants
- **Message System**: Integrates with character.SendMessage()
- **World Interface**: Uses GetRoom() for door synchronization

## Usage Examples

### Opening Containers
```
> open chest
Ok.
> open chest
But it's already open!
> open locked_box
It seems to be locked.
```

### Closing Containers
```
> close chest
Ok.
> close chest
But it's already closed!
```

### Opening Doors
```
> open door
Ok.
> open door north
Ok.
> open gate
Ok.
```

### Closing Doors
```
> close door
Ok.
> close gate south
Ok.
```

### Error Cases
```
> open sword
That's not a container.
> open nonexistent
I see no nonexistent here.
> open
Open what?
```

## Benefits

### ✅ **Complete DikuMUD Compatibility**
- Exact behavior match with original C implementation
- Same message format and error handling
- Identical flag operations and state management

### ✅ **Robust Implementation**
- Comprehensive error handling and validation
- Safe bitwise flag operations
- Proper two-way door synchronization

### ✅ **Flexible Usage**
- Works with containers in inventory or room
- Supports door keywords and directions
- Case-insensitive and partial name matching

### ✅ **Thorough Testing**
- 16+ unit tests covering all scenarios
- Container and door functionality verified
- Edge cases and error conditions tested

### ✅ **Clean Integration**
- Uses existing type system and constants
- Follows established DikuGo patterns
- Minimal dependencies and clean interfaces

The open and close commands are now fully functional and provide an authentic DikuMUD experience with proper container and door mechanics, state validation, and automatic door synchronization!
