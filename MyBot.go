package main

import (
	"./src/hlt"
	"./src/strat"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	logging := true
	botName := "bovard"

	conn := hlt.NewConnection(botName)

	// set up logging
	if logging {
		fname := strconv.Itoa(conn.PlayerTag) + "_gamelog.log"
		f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}
	gameMap := conn.UpdateMap()
	gameturn := 1
	gc := strat.GameController {
		GameMap:                 gameMap,
		ShipControllers:         make(map[int]strat.ShipController),
		ShipToPlanetAssignments: make(map[int][]int),
	}
	gc.UpdatePlanets(gameMap.Planets)
	for true {
		gameMap = conn.UpdateMap()
		gc.Update(gameMap)
		commandQueue := []string{}

		myPlayer := gameMap.Players[gameMap.MyId]
		myShips := myPlayer.Ships

		for i := 0; i < len(myShips); i++ {
			log.Printf("Ship #%v\n", i)
			ship := myShips[i]
			if ship.DockingStatus == hlt.UNDOCKED {
				log.Printf("Starting turn\n")
				cmd := hlt.SmarterBasicBot(ship, gameMap)
				log.Printf("CMD %v\n", cmd)
				commandQueue = append(commandQueue, cmd)
			}
		}
		log.Printf("Turn %v\n", gameturn)
		conn.SubmitCommands(commandQueue)
		gameturn++
	}
}
