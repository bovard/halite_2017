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
}

func CreateGameTurnInfo(gameMap *hlt.GameMap) GameTurnInfo {
	myId := gameMap.MyId
	myShipCount := len(gameMap.Players[myId].Ships)

	maxOpponentCount := 0
	minOpponentCount := 1000000
	for idx := range gameMap.Players {
		if idx == myId {
			continue
		}  
		numShips := len(gameMap.Players[idx].Ships)
		if numShips > maxOpponentCount {
			maxOpponentCount = numShips
		} 
		if numShips < minOpponentCount {
			minOpponentCount = numShips
		}

	}

	activateStupidRunAwayMeta := len(gameMap.Players) > 2 && gameMap.Turn > 100 && myShipCount * 3 < maxOpponentCount

	return GameTurnInfo{
		ShipCountDeltaToLeader:    myShipCount - maxOpponentCount,
		MyShipCount:               myShipCount,
		MaxOpponentShipCount:      maxOpponentCount,
		MinOpponentShipCount:      minOpponentCount,
		ActivateStupidRunAwayMeta: activateStupidRunAwayMeta,
	}
}
