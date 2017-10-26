package hlt

func StrategyBasicBot(ship Ship, gameMap Map) string {
	planets := gameMap.NearestPlanetsByDistance(ship)

	for i := 0; i < len(planets); i++ {
		planet := planets[i]
		if (planet.Owned == 0 || planet.Owner == gameMap.MyId) && planet.NumDockedShips < planet.NumDockingSpots && planet.Id%2 == ship.Id%2 {
			if ship.CanDock(planet) {
				return ship.Dock(planet)
			} else {
				return ship.Navigate(ship.ClosestPointTo(planet.Entity, 3), gameMap)
			}
		}
	}

	return ""
}

func SmarterBasicBot(ship Ship, gameMap Map) string {
	if ship.DockingStatus == DOCKED {
		return ""
	}
	planets := gameMap.NearestPlanetsByDistance(ship)
	enemies := gameMap.NearestEnemiesByDistance(ship)
	if enemies[0].Distance < planets[0].Distance {
		return ship.BetterNavigate(ship.ClosestPointTo(enemies[0], SHIP_MAX_ATTACK_RANGE - .1), gameMap)
	} 

	for i := 0; i < len(planets); i++ {
		planet := planets[i]
		if (planet.Owned == 0 || planet.Owner == gameMap.MyId) && planet.NumDockedShips < planet.NumDockingSpots {
			if ship.CanDock(planet) {
				return ship.Dock(planet)
			} else {
				return ship.BetterNavigate(ship.ClosestPointTo(planet.Entity, 3), gameMap)
			}
		}
	}

	return ship.BetterNavigate(ship.ClosestPointTo(enemies[0], SHIP_MAX_ATTACK_RANGE - .1), gameMap)
}
