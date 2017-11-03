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
	gc := strat.GameController{
		GameMap:         &gameMap,
		ShipControllers: make(map[int]*strat.ShipController),
	}

	var newGameMap hlt.GameMap
	for true {
		newGameMap = conn.UpdateMap()
		newGameMap.UpdateShipsFromHistory(&gameMap)
		gameMap = newGameMap

		gc.Update(&gameMap)
		gc.AssignToPlanets()
		commandQueue := []string{}

		myPlayer := gameMap.Players[gameMap.MyId]
		myShips := myPlayer.Ships

		for i := 0; i < len(myShips); i++ {
			ship := myShips[i]
			sc := *gc.ShipControllers[ship.Entity.Id]
			log.Println(sc.Id, "is assigned to planet ", sc.Planet)
			log.Println("Ship is located at ", ship.Point)
			log.Println("With Vel ", ship.Vel, " and mag ", ship.Vel.Magnitude())
			if sc.Planet != -1 {
				targetPlanet := gameMap.PlanetsLookup[sc.Planet]
				log.Println("planet location is ", targetPlanet.Point, ", d = ", ship.DistanceToCollision(&targetPlanet.Entity))
				rad := ship.Point.AngleTo(&targetPlanet.Point)
				log.Println(int(360+hlt.RadToDeg(rad)) % 360)
			}
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
