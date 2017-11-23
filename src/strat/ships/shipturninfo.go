package ships

import (
	"../../hlt"
)

type ShipInfo struct {
	Exists              bool
	Distance, Direction float64
	Ship                *hlt.Ship
}

type ShipTurnInfo struct {
	PossibleEnemyShipCollisions, PossibleAlliedShipCollisions []*hlt.Ship
	PossiblePlanetCollisions                                  []*hlt.Planet
	TotalEnemies, TotalAllies                                 int
	EnemiesInCombatRange, EnemiesDockedInCombatRange          int
	EnemiesInThreatRange, EnemiesInActiveThreatRange          int
	AlliesInCombatRange, AlliesDockedInCombatRange            int
	AlliesInThreatRange, AlliesInActiveThreatRange            int
	ClosestNonDockedEnemy                                     *ShipInfo
	ClosestDockedEnemy                                        *ShipInfo
	ClosestEnemy                                              *ShipInfo
	ClosestAlly                                               *ShipInfo
	ClosestEnemyShipClosingDistance                           bool
	PlanetsByDist                                             []*hlt.Planet
	EnemiesByDist, AlliesByDist                               []*hlt.Ship
	AlliedClosestPlanet, EnemyClosestPlanet                   *hlt.Planet
	AlliedClosestPlanetDist, EnemyClosestPlanetDist           float64
}

func CreateShipTurnInfo(ship *hlt.Ship, gameMap *hlt.GameMap) ShipTurnInfo {

	possiblePlanetCollisions := []*hlt.Planet{}
	for _, p := range gameMap.PlanetLookup {
		if ship.DistanceToCollision(&p.Entity) <= hlt.SHIP_MAX_SPEED {
			possiblePlanetCollisions = append(possiblePlanetCollisions, p)
		}
	}

	possibleEnemyShipCollisions := []*hlt.Ship{}
	for _, id := range gameMap.EnemyShips {
		s := gameMap.ShipLookup[id]
		if ship.DistanceToCollision(&s.Entity) <= 2*hlt.SHIP_MAX_SPEED {
			possibleEnemyShipCollisions = append(possibleEnemyShipCollisions, s)
		}
	}

	possibleAlliedShipCollisions := []*hlt.Ship{}
	for _, id := range gameMap.MyShips {
		s := gameMap.ShipLookup[id]
		if ship.DistanceToCollision(&s.Entity) <= 2*hlt.SHIP_MAX_SPEED {
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
	closestAlly := ShipInfo{
		Distance: 10000.0,
		Exists:   false,
	}
	if len(allies) > 1 {
		ship := gameMap.ShipLookup[allies[1].Id]
		closestAlly = ShipInfo{
			Distance:  allies[1].Distance,
			Exists:    true,
			Ship:      ship,
			Direction: ship.AngleTo(&ship.Point),
		}
	}
	closestEnemyShip := gameMap.ShipLookup[enemies[0].Id]
	closestEnemy := ShipInfo{
		Exists:    true,
		Distance:  enemies[0].Distance,
		Ship:      closestEnemyShip,
		Direction: ship.AngleTo(&closestEnemyShip.Point),
	}
	v2e := ship.Vel.Add(&closestEnemyShip.Vel)
	closestEnemyShipClosingDistance := v2e.Magnitude() < ship.Vel.Magnitude()

	closestDockedEnemy := ShipInfo{
		Exists:   false,
		Distance: 10000,
	}
	closestNonDockedEnemy := ShipInfo{
		Exists:   false,
		Distance: 10000,
	}

	foundNon := false
	foundDocked := false
	for _, e := range enemies {
		if e.DockingStatus == hlt.UNDOCKED && !foundNon {
			closestDockedEnemy = ShipInfo{
				Exists:    true,
				Ship:      e,
				Distance:  e.Distance,
				Direction: ship.AngleTo(&e.Point),
			}
			foundNon = true

		}
		if e.DockingStatus != hlt.UNDOCKED && !foundDocked {
			closestNonDockedEnemy = ShipInfo{
				Exists:    true,
				Ship:      e,
				Distance:  e.Distance,
				Direction: ship.AngleTo(&e.Point),
			}
			foundDocked = true
		}
		if foundDocked && foundNon {
			break
		}
	}

	planets := gameMap.NearestPlanetsByDistance(ship)

	var alliedClosestPlanet *hlt.Planet
	alliedClosestPlanetDist := 100000000.0
	var enemyClosestPlanet *hlt.Planet
	enemyClosestPlanetDist := 100000000.0

	for _, p := range planets {
		if p.Owned == 1 && p.Owner == gameMap.MyId {
			if p.Distance < alliedClosestPlanetDist {
				alliedClosestPlanetDist = p.Distance
				alliedClosestPlanet = p
			}
		} else if p.Owned == 1 {
			if p.Distance < enemyClosestPlanetDist {
				enemyClosestPlanetDist = p.Distance
				enemyClosestPlanet = p
			}
		}
	}

	for _, id := range append(gameMap.MyShips, gameMap.EnemyShips...) {
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
		} else if dist <= hlt.SHIP_MAX_SPEED+hlt.SHIP_MAX_ATTACK_RANGE {
			if s.DockingStatus == hlt.UNDOCKED {
				if s.Owner == gameMap.MyId {
					alliesInThreatRange++
				} else {
					enemiesInThreatRange++
				}
			}
		} else if dist <= 2*hlt.SHIP_MAX_SPEED+hlt.SHIP_MAX_ATTACK_RANGE {
			if s.DockingStatus == hlt.UNDOCKED {
				if s.Owner == gameMap.MyId {
					alliesInActiveThreatRange++
				} else {
					enemiesInActiveThreatRange++
				}
			}
		}
	}

	return ShipTurnInfo{
		PossibleEnemyShipCollisions:     possibleEnemyShipCollisions,
		PossibleAlliedShipCollisions:    possibleAlliedShipCollisions,
		PossiblePlanetCollisions:        possiblePlanetCollisions,
		EnemiesInCombatRange:            enemiesInCombatRange,
		EnemiesDockedInCombatRange:      enemiesDockedInCombatRange,
		EnemiesInThreatRange:            enemiesInThreatRange,
		EnemiesInActiveThreatRange:      enemiesInActiveThreatRange,
		TotalEnemies:                    enemiesInCombatRange + enemiesInThreatRange + enemiesInActiveThreatRange,
		AlliesInCombatRange:             alliesInCombatRange,
		AlliesDockedInCombatRange:       alliesDockedInCombatRange,
		AlliesInThreatRange:             alliesInThreatRange,
		AlliesInActiveThreatRange:       alliesInActiveThreatRange,
		TotalAllies:                     alliesInCombatRange + alliesInThreatRange + alliesInActiveThreatRange,
		ClosestAlly:                     &closestAlly,
		ClosestEnemy:                    &closestEnemy,
		ClosestNonDockedEnemy:           &closestNonDockedEnemy,
		ClosestDockedEnemy:              &closestDockedEnemy,
		ClosestEnemyShipClosingDistance: closestEnemyShipClosingDistance,
		PlanetsByDist:                   planets,
		EnemiesByDist:                   enemies,
		AlliesByDist:                    allies,
		AlliedClosestPlanet:             alliedClosestPlanet,
		EnemyClosestPlanet:              enemyClosestPlanet,
		AlliedClosestPlanetDist:         alliedClosestPlanetDist,
		EnemyClosestPlanetDist:          enemyClosestPlanetDist,
	}

}
