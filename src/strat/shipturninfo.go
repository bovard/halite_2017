package strat

import (
	"../hlt"
)

type ShipTurnInfo struct {
	PossibleEnemyShipCollisions, PossibleAlliedShipCollisions []*hlt.Ship
	PossiblePlanetCollisions []hlt.Planet
	EnemiesInCombatRange, EnemiesDockedInCombatRange, EnemiesInThreatRange, AlliesInCombatRange, AlliesDockedInCombatRange, AlliesInThreatRange int
	ClosestEnemyShipDistance, ClosestEnemyShipDir, ClosestAlliedShipDistance, ClosestAlliedShipDir float64
	ClosestEnemyShip, ClosestAlliedShip hlt.Ship
}
