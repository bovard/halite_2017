package ships

import (
	"../../hlt"
	"log"
)


func (self *ShipController) NormalSetTarget(gameMap *hlt.GameMap) {

	if self.TargetPlanet == -1 {
		self.Target = &self.Info.ClosestEnemyShip.Point
	} else if self.TargetPlanet != -1 {
		valid, _ := self.IsTargetPlanetStillValid(gameMap)
		if valid {
			planet := gameMap.PlanetLookup[self.TargetPlanet]
			self.Target = &planet.Point
		} else {
			log.Println("invalid planet, moving to enemy")
			self.Target = &self.Info.ClosestEnemyShip.Point
			self.TargetPlanet = -1
		}
	}

}