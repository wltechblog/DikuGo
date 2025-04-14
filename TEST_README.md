# DikuGo Test Scripts

These TinTin++ scripts are designed to test various aspects of the DikuGo MUD implementation.

## Prerequisites

- TinTin++ installed on your system
- DikuGo server running on localhost:4000

## Available Test Scripts

### 1. test_navigation_fixed.tt

Tests basic navigation through the game world, including:
- Moving in cardinal directions
- Using the goto command
- Looking at rooms and objects

### 2. test_equipment_fixed.tt

Tests the equipment system, including:
- Inventory and equipment commands
- Get/drop commands
- Wear/remove commands
- Testing different equipment positions

### 3. test_simple.tt

A simplified test script that tests basic navigation, extra descriptions, and equipment functionality.

### 4. test_tintin.tt

The simplest test script that uses only basic TinTin++ commands with fixed delays. This is the most reliable script to use if you're having issues with the other scripts.

## How to Run

1. Start the DikuGo server:
   ```
   go run cmd/dikugo/main.go
   ```

2. In a separate terminal, run one of the test scripts with TinTin++:
   ```
   tt++ test_tintin.tt
   ```
   or
   ```
   tt++ test_simple.tt
   ```

3. The script will automatically:
   - Connect to the server
   - Create or log in as a test user
   - Run through the test sequence
   - Disconnect when complete

## Important Note About TinTin++ Syntax

TinTin++ scripts should use semicolons to separate commands on a single line rather than using multi-line format. For example:

```
#ACTION {pattern} {#ECHO {message}; #SEND {command}}
```

Rather than:

```
#ACTION {pattern} {
    #ECHO {message}
    #SEND {command}
}
```

The test_tintin.tt script uses the simplest possible syntax and should work with any version of TinTin++.

## Manual Testing

You can also use these scripts as a basis for manual testing. To do this:

1. Connect to the server using TinTin++:
   ```
   tt++
   #SESSION {DikuGoTest} {localhost} {4000}
   ```

2. After logging in, you can test various commands manually:
   ```
   look
   north
   look
   goto 3001
   look fountain
   get all
   inventory
   wear all
   equipment
   remove all
   drop all
   testexits
   ```

## Customizing Tests

Feel free to modify these scripts to test additional features or specific scenarios.
