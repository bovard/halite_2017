package main

import (
	"./src/hlt"
	"./src/strat"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
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
	gameMap := conn.UpdateMap(0)
	gameturn := 1
	gc := strat.GameController{
		GameMap:         &gameMap,
		ShipControllers: make(map[int]*strat.ShipController),
		ShipNumIdx:      1,
	}

	for true {
		log.Println("Game Turn: ", gameturn)
		start := time.Now()
		newMap := conn.UpdateMap(gameturn)
		gc.Update(&newMap)
		gameMap = newMap

		commandQueue := gc.Act(gameturn)

		conn.SubmitCommands(commandQueue)
		gameturn++
		t := time.Now()
		log.Println("this turn took", t.Sub(start))
	}
}
