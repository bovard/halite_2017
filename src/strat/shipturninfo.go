package strat

import (
	"../hlt"
)

type ShipTurnInfo struct {
	PossibleEnemyShipCollisions, PossibleAlliedShipCollisions []*hlt.Ship
	PossiblePlanetCollisions []*hlt.Planet
	TotalEnemies, TotalAllies, EnemiesInCombatRange, EnemiesDockedInCombatRange, EnemiesInThreatRange, EnemiesInActiveThreatRange, AlliesInCombatRange, AlliesDockedInCombatRange, AlliesInThreatRange, AlliesInActiveThreatRange int
	ClosestDockedEnemyShipDistance, ClosestDockedEnemyShipDir, ClosestEnemyShipDistance, ClosestEnemyShipDir, ClosestAlliedShipDistance, ClosestAlliedShipDir float64
	ClosestDockedEnemyShip, ClosestEnemyShip, ClosestAlliedShip *hlt.Ship
	ClosestEnemyShipClosingDistance bool
	PlanetsByDist []*hlt.Planet
	EnemiesByDist, AlliesByDist []*hlt.Ship
	AlliedClosestPlanet, EnemyClosestPlanet *hlt.Planet
	AlliedClosestPlanetDist, EnemyClosestPlanetDist float64
}

func CreateShipTurnInfo(ship *hlt.Ship, gameMap *hlt.GameMap) ShipTurnInfo {

	possiblePlanetCollisions := []*hlt.Planet{}
	for _, p := range(gameMap.PlanetLookup) {
		if ship.DistanceToCollision(&p.Entity) <= hlt.SHIP_MAX_SPEED {
			possiblePlanetCollisions = append(possiblePlanetCollisions, p)
		}
	}

	possibleEnemyShipCollisions := []*hlt.Ship{}
	for _, id := range(gameMap.EnemyShips) {
		s := gameMap.ShipLookup[id]
		if ship.DistanceToCollision(&s.Entity) <= 2 * hlt.SHIP_MAX_SPEED {
			possibleEnemyShipCollisions = append(possibleEnemyShipCollisions, s)
		}
	}

	possibleAlliedShipCollisions := []*hlt.Ship{}
	for _, id := range(gameMap.MyShips) {
		s := gameMap.ShipLookup[id]
		if ship.DistanceToCollision(&s.Entity) <= 2 * hlt.SHIP_MAX_SPEED {
			possibleAlliedShipCollisions = append(possibleAlliedShipCollisions, s)
		}
	}

	enemies := gameMap.NearestShipsByDistance(ship, gameMap.EnemyShips)
	allies := gameMap.NearestShipsByDistance(ship, gameMap.MyShips)
	enemiesInCombatRange := 0
	enemiesDockedInCombatRange := 0
	enemiesInThreatRange := 0
	enemiesInActiveThreatRange := 0
	alliesInCombatRange := 0
	alliesDockedInCombatRange := 0
	alliesInThreatRange := 0
	alliesInActiveThreatRange := 0
	closestAlliedShipDistance := allies[0].Distance
	closestAlliedShip := gameMap.ShipLookup[allies[0].Id]
	closestAlliedShipDir := ship.AngleTo(&closestAlliedShip.Point)
	if len(allies) > 1 {
		closestAlliedShipDistance = allies[1].Distance
		closestAlliedShip = gameMap.ShipLookup[allies[1].Id]
		closestAlliedShipDir = ship.AngleTo(&closestAlliedShip.Point)
	}
	closestEnemyShipDistance := enemies[0].Distance
	closestEnemyShip := gameMap.ShipLookup[enemies[0].Id]
	closestEnemyShipDir := ship.AngleTo(&closestEnemyShip.Point)
	v2e := ship.Vel.Add(&closestEnemyShip.Vel)
	closestEnemyShipClosingDistance := v2e.Magnitude() < ship.Vel.Magnitude()
	var closestDockedEnemyShip *hlt.Ship 
	var closestDockedEnemyShipDir float64
	closestDockedEnemyShipDistance := 100000000.0
	for _, e := range(enemies) {
		s := gameMap.ShipLookup[enemies[0].Id]
		if s.DockingStatus != hlt.UNDOCKED {
			closestDockedEnemyShip = s
			closestDockedEnemyShipDistance = e.Distance
			closestDockedEnemyShipDir = ship.AngleTo(&closestDockedEnemyShip.Point)
			break;
		}
	}

	planets := gameMap.NearestPlanetsByDistance(ship)

	var alliedClosestPlanet *hlt.Planet
	alliedClosestPlanetDist := 100000000.0
	var enemyClosestPlanet *hlt.Planet
	enemyClosestPlanetDist :=  100000000.0

	for _, p := range(planets) {
		if p.Owner == gameMap.MyId {
			if p.Distance < alliedClosestPlanetDist {
				alliedClosestPlanetDist = p.Distance
				alliedClosestPlanet = p
			} 
		} else if p.Owner != 0 {
			if p.Distance < enemyClosestPlanetDist {
				enemyClosestPlanetDist = p.Distance
				enemyClosestPlanet = p
			} 
		}
	}

	for _, id := range(append(gameMap.MyShips, gameMap.EnemyShips...)) {
		if id == ship.Id {
			continue
		}
		s := gameMap.ShipLookup[id]
		dist := ship.DistanceToCollision(&s.Entity)
		if dist <= hlt.SHIP_MAX_ATTACK_RANGE {
			if s.DockingStatus == hlt.UNDOCKED {
				if s.Owner == gameMap.MyId {
					alliesInCombatRange++
				} else {
					enemiesInCombatRange++
				}
			} else {
				if s.Owner == gameMap.MyId {
					alliesDockedInCombatRange++
				} else {
					enemiesDockedInCombatRange++
				}
			}
		} else if dist <= hlt.SHIP_MAX_SPEED + hlt.SHIP_MAX_ATTACK_RANGE {
			if s.DockingStatus == hlt.UNDOCKED {
				if s.Owner == gameMap.MyId {
					alliesInThreatRange++
				} else {
					enemiesInThreatRange++
				}
			}
		} else if dist <= 2 *hlt.SHIP_MAX_SPEED + hlt.SHIP_MAX_ATTACK_RANGE {
			if s.DockingStatus == hlt.UNDOCKED {
				if s.Owner == gameMap.MyId {
					alliesInActiveThreatRange++
				} else {
					enemiesInActiveThreatRange++
				}
			}
		}
	} 

	return ShipTurnInfo {
		PossibleEnemyShipCollisions: possibleEnemyShipCollisions, 
		PossibleAlliedShipCollisions: possibleAlliedShipCollisions,
		PossiblePlanetCollisions: possiblePlanetCollisions,
		EnemiesInCombatRange: enemiesInCombatRange, 
		EnemiesDockedInCombatRange: enemiesDockedInCombatRange, 
		EnemiesInThreatRange: enemiesInThreatRange, 
		EnemiesInActiveThreatRange: enemiesInActiveThreatRange,
		TotalEnemies: enemiesInCombatRange + enemiesInThreatRange + enemiesInActiveThreatRange,
		AlliesInCombatRange: alliesInCombatRange, 
		AlliesDockedInCombatRange: alliesDockedInCombatRange, 
		AlliesInThreatRange: alliesInThreatRange,
		AlliesInActiveThreatRange: alliesInActiveThreatRange,
		TotalAllies: alliesInCombatRange + alliesInThreatRange + alliesInActiveThreatRange,
		ClosestEnemyShipDistance: closestEnemyShipDistance,
		ClosestDockedEnemyShipDistance: closestDockedEnemyShipDistance,
		ClosestEnemyShipDir: closestEnemyShipDir,
		ClosestEnemyShipClosingDistance: closestEnemyShipClosingDistance,
		ClosestDockedEnemyShipDir: closestDockedEnemyShipDir,
		ClosestAlliedShipDistance: closestAlliedShipDistance,
		ClosestAlliedShipDir: closestAlliedShipDir,
		ClosestEnemyShip: closestEnemyShip,
		ClosestDockedEnemyShip: closestDockedEnemyShip,
		ClosestAlliedShip: closestAlliedShip ,
		PlanetsByDist: planets,
		EnemiesByDist: enemies,
		AlliesByDist: allies,
		AlliedClosestPlanet: alliedClosestPlanet, 
		EnemyClosestPlanet: enemyClosestPlanet,
		AlliedClosestPlanetDist: alliedClosestPlanetDist, 
		EnemyClosestPlanetDist: enemyClosestPlanetDist ,
	}

}