package strat

import ( "../hlt")

type ShipController struct {
	Ship     *hlt.Ship
	Past     [] *hlt.Ship
	Id       int
	Planet   int
	Alive    bool
}

func (self *ShipController) Update(ship *hlt.Ship) {
	self.Past = append(self.Past, self.Ship)
	self.Ship = ship
}


func (self *ShipController) Act(gameMap *hlt.Map) string {
	enemies := gameMap.NearestEnemiesByDistance(*self.Ship)
	closetEnemy := enemies[0].Distance
	if self.Planet != -1 {
		planet := gameMap.PlanetsLookup[self.Planet]
		planetDist := self.Ship.Entity.CalculateDistanceTo(planet.Entity)
		if closetEnemy / 2 < planetDist {
			self.Planet = -1
			return self.Ship.BetterNavigate(enemies[0], *gameMap)
		} else {
			if self.Ship.CanDock(planet) {
				return self.Ship.Dock(planet)
			} else {
				//return self.Ship.NavigateBasic(planet.Entity)
				return self.Ship.BetterNavigate(planet.Entity, *gameMap)
			}
		}
	} else {
		return self.Ship.BetterNavigate(enemies[0], *gameMap)
	}
	return ""
}