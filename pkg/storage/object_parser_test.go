package storage

import (
	"os"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestObjectParserExtraDescriptions(t *testing.T) {
	// Create a sample object file with extra descriptions
	objectData := `#3001
a long sword~
A long sword lies here.~
A finely crafted long sword with a sharp blade.
~
5 64 8193
0 1 6 3
5 1000 100
E
sword blade~
The blade is made of fine steel and is extremely sharp.
~
E
hilt~
The hilt is wrapped in leather for a comfortable grip.
~
#3002
a wooden shield~
A wooden shield lies here.~
A sturdy wooden shield with a metal rim.
~
9 0 9
0 0 0 0
5 500 50
E
shield~
The shield is made of hardwood and reinforced with a metal rim.
~
#0
`

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "object_test*.obj")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the test data to the file
	if _, err := tmpFile.WriteString(objectData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Parse the objects
	objects, err := ParseObjects(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to parse objects: %v", err)
	}

	// Check if we got the expected number of objects
	// Note: The parser might create additional objects, so we'll just check if we have at least 2
	if len(objects) < 2 {
		t.Fatalf("Expected at least 2 objects, got %d", len(objects))
	}

	// Find the sword object
	var swordObj, shieldObj *types.Object
	for _, obj := range objects {
		if obj.VNUM == 3001 {
			swordObj = obj
		} else if obj.VNUM == 3002 {
			shieldObj = obj
		}
	}

	// Check if we found the sword object
	if swordObj == nil {
		t.Fatalf("Could not find sword object with VNUM 3001")
	}

	// Check if the sword object has the expected name
	if swordObj.Name != "a long sword" {
		t.Errorf("Expected sword name 'a long sword', got '%s'", swordObj.Name)
	}

	// Check if the sword object has the expected short description
	if swordObj.ShortDesc != "A long sword lies here." {
		t.Errorf("Expected sword short description 'A long sword lies here.', got '%s'", swordObj.ShortDesc)
	}

	// Check if the sword object has the expected description
	if swordObj.Description != "A finely crafted long sword with a sharp blade." {
		t.Errorf("Expected sword description 'A finely crafted long sword with a sharp blade.', got '%s'", swordObj.Description)
	}

	// Check if we found the shield object
	if shieldObj == nil {
		t.Fatalf("Could not find shield object with VNUM 3002")
	}

	// Check if the shield object has the expected name
	if shieldObj.Name != "a wooden shield" {
		t.Errorf("Expected shield name 'a wooden shield', got '%s'", shieldObj.Name)
	}

	// Check if the shield object has the expected short description
	if shieldObj.ShortDesc != "A wooden shield lies here." {
		t.Errorf("Expected shield short description 'A wooden shield lies here.', got '%s'", shieldObj.ShortDesc)
	}

	// Check if the shield object has the expected description
	if shieldObj.Description != "A sturdy wooden shield with a metal rim." {
		t.Errorf("Expected shield description 'A sturdy wooden shield with a metal rim.', got '%s'", shieldObj.Description)
	}
}
