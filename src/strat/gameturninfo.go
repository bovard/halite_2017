package strat

import (
	"../hlt"
	"math"
)

type GameTurnInfo struct {
	ShipCountDeltaToLeader    int
	MyShipCount               int
	MaxOpponentShipCount      int
	MinOpponentShipCount      int
	ActivateStupidRunAwayMeta bool
	NumEnemyPlanets           int
	NumEnemies                int
	PrimaryOpponentDied       bool
	MinEnemyDist              float64
}

func CreateGameTurnInfo(gameMap *hlt.GameMap, oldGameMap *hlt.GameMap) GameTurnInfo {
	myId := gameMap.MyId

	myShipCount := len(gameMap.MyShips)
	shipDiff := len(oldGameMap.EnemyShips) - len(gameMap.EnemyShips) 
	shipPer := float64(len(gameMap.EnemyShips)) / float64(len(oldGameMap.EnemyShips))

	primaryOpponentDied := shipDiff > 20 && shipPer < .3

	maxOpponentCount := 0
	minOpponentCount := 1000000
	numOpponents := 0
	for idx := range gameMap.Players {
		playerId := gameMap.Players[idx].Id
		if playerId == myId {
		}
		numShips := len(gameMap.Players[idx].Ships)
		if numShips > 0 {
			numOpponents += 1
		}
		if numShips > maxOpponentCount {
			maxOpponentCount = numShips
		}
		if numShips < minOpponentCount {
			minOpponentCount = numShips
		}

	}

	numEnemyPlanets := 0
	minEnemyDist := 10000.0
	for _, p := range gameMap.PlanetLookup {
		if p.Owned == 1 {
			if p.Owner == gameMap.MyId {
				for _, s := range gameMap.ShipLookup {
					d := p.SqDistanceTo(&s.Point)
					if d < minEnemyDist {
						minEnemyDist = d
					}
				}
			} else {
				numEnemyPlanets++
			}
		}
	}
	minEnemyDist = math.Sqrt(minEnemyDist)

	activateStupidRunAwayMeta := numOpponents > 1 && ((gameMap.Turn > 100 && myShipCount*3 < maxOpponentCount) || (gameMap.Turn > 50 && (myShipCount < 10 && maxOpponentCount > 30)))

	return GameTurnInfo{
		ShipCountDeltaToLeader:    myShipCount - maxOpponentCount,
		MyShipCount:               myShipCount,
		MaxOpponentShipCount:      maxOpponentCount,
		MinOpponentShipCount:      minOpponentCount,
		ActivateStupidRunAwayMeta: activateStupidRunAwayMeta,
		NumEnemyPlanets:           numEnemyPlanets,
		NumEnemies:                numOpponents,
		PrimaryOpponentDied:       primaryOpponentDied,
		MinEnemyDist: 			   minEnemyDist,
	}
}
