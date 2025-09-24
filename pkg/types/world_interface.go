package types

// WorldInterface defines the methods that a World implementation must provide
type WorldInterface interface {
	// Character methods
	AddCharacter(character *Character)
	RemoveCharacter(character *Character)
	GetCharacter(name string) *Character
	SaveCharacter(character *Character) error
	
	// Room methods
	GetRoom(vnum int) *Room
	CharToRoom(character *Character, room *Room)
	CharFromRoom(character *Character)
	
	// Object methods
	GetObject(vnum int) *Object
	CreateObject(vnum int) *ObjectInstance
	ObjectToChar(obj *ObjectInstance, ch *Character)
	ObjectToRoom(obj *ObjectInstance, room *Room)
	ExtractObj(obj *ObjectInstance)
	
	// Affect methods
	AffectToChar(ch *Character, af *Affect)
	AffectFromChar(ch *Character, type_ int)
	AffectedBySpell(ch *Character, type_ int) bool
	ApplyObjectAffects(ch *Character, obj *ObjectInstance, add bool)
	
	// Message methods
	SendMessageToCharacter(character *Character, message string)
	Act(msg string, hide bool, ch *Character, obj *ObjectInstance, vict *Character, msgType int)
}
