# Active Context

## Current Focus
The current focus is on fixing and improving the combat system, particularly the damage calculation for NPCs.

## Recent Changes
1. **Combat System Fix**: Fixed the NPC damage calculation in the combat system. Previously, it was incorrectly using the hit dice values (Dice[0], Dice[1], Dice[2]) for damage calculation instead of the damage dice values (DamageType, AttackType).

2. **Mobstat Command Enhancement**: Updated the mobstat command to display both hit dice and damage dice values, making it clearer how the mob's HP and damage are calculated.

## Key Issues Addressed
- **Incorrect Damage Calculation**: NPCs were doing too much damage because the system was using hit dice values (which can be very high for HP calculation) instead of damage dice values.
- **Missing Information in Mobstat**: The mobstat command was not showing the damage dice values, making it difficult to understand how NPC damage was calculated.

## Implementation Details
1. In `pkg/combat/diku_combat.go`, modified the NPC damage calculation to use `DamageType` and `AttackType` instead of `Dice[0]` and `Dice[1]`.
2. Added a fallback to hit dice values if damage dice values are not set.
3. In `pkg/command/mobstat.go`, updated the output to show both hit dice and damage dice values.

## Current State
- The combat system now correctly uses damage dice values for NPC damage calculation.
- The mobstat command now shows both hit dice and damage dice values.
- NPCs like the rat now do damage consistent with their level and the original DikuMUD implementation.

## Next Steps
Potential areas for further improvement:
1. Review other aspects of the combat system for accuracy
2. Enhance the mobstat command with more detailed information
3. Add more comprehensive testing for the combat system
4. Consider adding configuration options for combat balance