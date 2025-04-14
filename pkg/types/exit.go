package types

// Exit flag constants
const (
	EX_ISDOOR    = (1 << 0)
	EX_CLOSED    = (1 << 1)
	EX_LOCKED    = (1 << 2)
	EX_PICKPROOF = (1 << 3)
	EX_HIDDEN    = (1 << 4)
)

// IsClosed returns true if the exit is closed
func (e *Exit) IsClosed() bool {
	return (e.Flags & EX_CLOSED) != 0
}

// IsLocked returns true if the exit is locked
func (e *Exit) IsLocked() bool {
	return (e.Flags & EX_LOCKED) != 0
}

// IsDoor returns true if the exit is a door
func (e *Exit) IsDoor() bool {
	return (e.Flags & EX_ISDOOR) != 0
}

// IsPickproof returns true if the exit is pickproof
func (e *Exit) IsPickproof() bool {
	return (e.Flags & EX_PICKPROOF) != 0
}

// IsHidden returns true if the exit is hidden
func (e *Exit) IsHidden() bool {
	return (e.Flags & EX_HIDDEN) != 0
}
