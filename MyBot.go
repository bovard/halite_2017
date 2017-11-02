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
	botName := "current"

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
		GameMap:                 &gameMap,
		ShipControllers:         make(map[int]*strat.ShipController),
	}
	for true {
		gameMap = conn.UpdateMap()
		gc.Update(&gameMap)
		gc.AssignToPlanets()
		commandQueue := []string{}

		myPlayer := gameMap.Players[gameMap.MyId]
		myShips := myPlayer.Ships

		for i := 0; i < len(myShips); i++ {
			ship := myShips[i]
			sc := *gc.ShipControllers[ship.Entity.Id]
			log.Println(sc.Id, "is assigned to planet ", sc.Planet)
			log.Println(sc.Ship.Entity.X, sc.Ship.Entity.Y)
			log.Println(ship.Entity.X, ship.Entity.Y)
			log.Println(ship.Point)
			if ship.DockingStatus == hlt.UNDOCKED {
				cmd := sc.Act(&gameMap)
				log.Println(cmd)
				commandQueue = append(commandQueue, cmd)
			}
		}
		conn.SubmitCommands(commandQueue)
		gameturn++
	}
}
