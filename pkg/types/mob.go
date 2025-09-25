package types

// Mobile action flag constants
const (
	ACT_SPEC       = (1 << 0) // Special routine to be called if exist
	ACT_SENTINEL   = (1 << 1) // This mobile not to be moved
	ACT_SCAVENGER  = (1 << 2) // Pick up stuff lying around
	ACT_ISNPC      = (1 << 3) // This bit is set for use with IS_NPC()
	ACT_NICE_THIEF = (1 << 4) // Set if a thief should NOT be killed
	ACT_AGGRESSIVE = (1 << 5) // Set if automatic attack on NPCs
	ACT_STAY_ZONE  = (1 << 6) // MOB Must stay inside its own zone
	ACT_WIMPY      = (1 << 7) // MOB Will flee when injured, and if aggressive only attack sleeping players
	ACT_FOLLOWER   = (1 << 8) // MOB is a follower/pet
)

// Mobile-related constants

// IsScavenger returns true if the mobile is a scavenger
func (m *Mobile) IsScavenger() bool {
	return (m.ActFlags & ACT_SCAVENGER) != 0
}

// IsSentinel returns true if the mobile is a sentinel
func (m *Mobile) IsSentinel() bool {
	return (m.ActFlags & ACT_SENTINEL) != 0
}

// IsAggressive returns true if the mobile is aggressive
func (m *Mobile) IsAggressive() bool {
	return (m.ActFlags & ACT_AGGRESSIVE) != 0
}

// IsStayZone returns true if the mobile stays in its zone
func (m *Mobile) IsStayZone() bool {
	return (m.ActFlags & ACT_STAY_ZONE) != 0
}

// IsWimpy returns true if the mobile is wimpy
func (m *Mobile) IsWimpy() bool {
	return (m.ActFlags & ACT_WIMPY) != 0
}

// IsNiceThief returns true if the mobile is nice to thieves
func (m *Mobile) IsNiceThief() bool {
	return (m.ActFlags & ACT_NICE_THIEF) != 0
}

// HasSpecProc returns true if the mobile has a special procedure
func (m *Mobile) HasSpecProc() bool {
	return (m.ActFlags & ACT_SPEC) != 0
}
