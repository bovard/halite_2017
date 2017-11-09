package strat

import (
	"../hlt"
)

type GameTurnInfo struct {
	ShipCountDeltaToLeader int
	MyShipCount int
	MaxOpponentShipCount int
}


func CreateGameTurnInfo(gameMap *hlt.GameMap) GameTurnInfo {
	myId := gameMap.MyId
	myShipCount := len(gameMap.Players[myId].Ships)

	maxOpponentCount := 0
	for idx, _ := range(gameMap.Players) {
		if idx == myId {
			continue
		} else if len(gameMap.Players[idx].Ships) > maxOpponentCount {
			maxOpponentCount = len(gameMap.Players[idx].Ships)
		}

	}

	return GameTurnInfo {
		ShipCountDeltaToLeader: myShipCount - maxOpponentCount,
		MyShipCount: myShipCount,
		MaxOpponentShipCount: maxOpponentCount,
	}
}
