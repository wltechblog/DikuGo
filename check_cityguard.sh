#!/bin/bash
# Connect to the game and check the cityguard stats
(
  sleep 1
  echo "admin"
  sleep 1
  echo "Y"
  sleep 1
  echo "password"
  sleep 1
  echo "password"
  sleep 1
  echo "goto 3111"
  sleep 1
  echo "look"
  sleep 1
  echo "kill cityguard"
  sleep 1
  echo "quit"
  sleep 1
  echo "Y"
) | telnet localhost 4000
