package strat

import ( 
  "log"
  "../hlt"
)

type ShipController struct {
	Ship     *hlt.Ship
	Past     [] *hlt.Ship
	Id       int
	Planet   int
}

func (self *ShipController) Update(ship *hlt.Ship) {
	self.Past = append(self.Past, self.Ship)
	self.Ship = ship
}


func (self *ShipController) Act(gameMap *hlt.Map) string {
	log.Println("Ship ", self.Id, " Act. Planet is ", self.Planet)
	enemies := gameMap.NearestEnemiesByDistance(*self.Ship)
	closetEnemy := enemies[0].Distance
	if self.Planet != -1 {
		planet := gameMap.PlanetsLookup[self.Planet]
		planetDist := self.Ship.Entity.DistanceToCollision(&planet.Entity)
		log.Println("distance to enemy is ", closetEnemy.Distance, )
		if closetEnemy / 2 < planetDist {
			self.Planet = -1
			return self.Ship.BetterNavigate(&enemies[0], gameMap)
		} else {
			if self.Ship.CanDock(&planet) {
				return self.Ship.Dock(&planet)
			} else {
				log.Println("moving toward planet", planet.Id)
				//return self.Ship.NavigateBasic(planet.Entity)
				return self.Ship.BetterNavigate(&planet.Entity, gameMap)
			}
		}
	} else {
		return self.Ship.BetterNavigate(&enemies[0], gameMap)
	}
	return ""
}