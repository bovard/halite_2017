package ships

import (
	"../../hlt"
	"log"
)


func (self *ShipController) SettlerSetTarget(gameMap *hlt.GameMap) {
	if self.TargetPlanet == -1 {
		log.Println("We don't have a valid target, switching to normal")
		self.Mission = MISSION_NORMAL
		self.NormalSetTarget(gameMap)
	} else if valid, _ := self.IsTargetPlanetStillValid(gameMap); valid {
		log.Println("our target planet is no longer a valid target, switching to normal")
		self.Mission = MISSION_NORMAL
		self.NormalSetTarget(gameMap)
	} else {
		planet := gameMap.PlanetLookup[self.TargetPlanet]
		self.Target = &planet.Point
	}
}