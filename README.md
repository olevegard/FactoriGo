# FactoriGo
An attempt at making an automation game similar to Factorio in Golang. But unlike Factorio, it won't be a 3d game where the player can walk around. Instead it'll be a simple menu based game where the player needs to balance buildings. The player will then have to balance the different type of buildings while also making sure they are powered. 


## Implemented so far:
- Simple GUI
- Can click to make buildings and increment inventory
- Buildings generate new items every x ticks, each building type has a different value for x
- "Recipe" system that allows for having a set of changes applied to inventoru ( ie "add 2 copper plates", "subtract 1 copper ore" }
- Recipe system used for creating new building and for when buildings generate new items\
- Inventory system with a map of inventory items with count and title

## Todo:
- Make and implement a better GUI
- Small windows for each production unit. 
- - These should show if they are getting what they need
- Miners should be harder to make, maybe requring the user to fist get stone -> make stone furnace -> get iron ore -> get coal -> make iron plates form iron ore -> build miners like in factorio
- Upgrades
- New materials
- Produce robotic warriors
- Robotic warriors fight enemies and brings back loot needed for research
- After a while ( or after some trigger like the player making warriros ) enemies start attacking
- Peaceful mode, where enemies never attack
- All aspects of warriors can be upgraded
- Can make new types of warriors ( possibly with different strenghts and weakneses)
- Different robots requires different materials
- Make game run on Android
- Make game run on iOS

## Investage ideas
- Making an "advanced mode" where buildings need to be conenected ie. iron miners need to be connected to smelters.
- - connections should be possbile to upgrade
- - Find a way to split this up. Maybe a button to display a group of production units and the connections between them.
- End game idea : colonize a new planet with more materials but also more powerful enemies. Player keep inventory / buildings?

## Other :
- System for generating coverage reports. Can use ` go test -coverprofile cover.out -tags test && go tool cover -html=cover.out -o cover.html
`, but should be in a script

## Development guidelines
- Functions should have no side effects ( take in an argument, return a copy, don't change the supplied arguments when pointer types are involved ( ie maps )
- Update logic should be entirely separate from game logic
- Code should follow go idioms as far as possible
- Update logic should have 100% test coverage if possible
- Update logic should be developed using a test driven approach
- Should have system tests for all platforms covering as much as possible.
