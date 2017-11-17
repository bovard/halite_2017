package strat

import (
	"../hlt"
)

type GameTurnInfo struct {
	ShipCountDeltaToLeader    int
	MyShipCount               int
	MaxOpponentShipCount      int
	MinOpponentShipCount      int
	ActivateStupidRunAwayMeta bool
	NumEnemyPlanets           int
	NumOpponents              int
}

func CreateGameTurnInfo(gameMap *hlt.GameMap) GameTurnInfo {
	myId := gameMap.MyId
	myShipCount := len(gameMap.Players[myId].Ships)

	maxOpponentCount := 0
	minOpponentCount := 1000000
	numOpponents := 0
	for idx := range gameMap.Players {
		if idx == myId {
			continue
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
	for _, p := range gameMap.PlanetLookup {
		if p.Owner != gameMap.MyId {
			numEnemyPlanets ++
		}

	}

	activateStupidRunAwayMeta := numOpponents > 1 && ((gameMap.Turn > 100 && myShipCount*3 < maxOpponentCount) || (gameMap.Turn > 50 && (myShipCount < 10 && maxOpponentCount > 30)))

	return GameTurnInfo{
		ShipCountDeltaToLeader:    myShipCount - maxOpponentCount,
		MyShipCount:               myShipCount,
		MaxOpponentShipCount:      maxOpponentCount,
		MinOpponentShipCount:      minOpponentCount,
		ActivateStupidRunAwayMeta: activateStupidRunAwayMeta,
		NumEnemyPlanets: 		   numEnemyPlanets,
		NumOpponents:              numOpponents,
	}
}
