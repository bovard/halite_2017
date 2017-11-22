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
	} else if valid := self.CanWeDockOnTargetPlanet(gameMap); !valid {
		log.Println("our target planet is no longer a valid target, switching to normal")
		self.TargetPlanet = -1
		self.Mission = MISSION_NORMAL
		self.NormalSetTarget(gameMap)
	} else {
		planet := gameMap.PlanetLookup[self.TargetPlanet]
		self.Target = &planet.Point
	}
}


func (self *ShipController) SettlerAct(gameMap *hlt.GameMap, turnComm *TurnComm) (ChlMessage, hlt.Heading) {
	heading := hlt.Heading{
		Magnitude: 0,
		Angle:     0,
	}
	message := NONE

	if self.Info.TotalEnemies > 0 {
		log.Println("We are in combat!")
		return self.combat(gameMap, turnComm)
	}

	planet := gameMap.PlanetLookup[self.TargetPlanet]
	log.Println("Continuing with assigned planet", planet.Id)
	if self.Ship.CanDock(planet) {
		log.Println("We can dock!")
		return DOCK, heading
	}
	h := self.MoveToDockingRange(planet, gameMap)
	if h.Magnitude > 0 {
		log.Println("can move to docking range of", planet.Id)
		message = MOVE_TO_DOCKING
		heading = h
	} else {
		log.Println("moving toward planet", planet.Id)
		message = MOVING_TOWARD_PLANET
		heading = self.MoveToPlanet(planet, gameMap)
	}
	return message, heading
}

func (self *ShipController) CanWeDockOnTargetPlanet(gameMap *hlt.GameMap) bool {
	if self.TargetPlanet == -1 {
		log.Println("Can't dock, planet -1 invalid")
		return false
	}

	if _, ok := gameMap.PlanetLookup[self.TargetPlanet]; !ok {
		log.Println("Can't dock, planet doesn't exist")
		return false
	}

	planet := gameMap.PlanetLookup[self.TargetPlanet]
	if planet.Owned == 1 && planet.Owner != gameMap.MyId {
		log.Println("Can't dock, enemy owns planet")
		return false
	}

	if planet.Owned == 1 && planet.Owner == gameMap.MyId && planet.NumDockedShips == planet.NumDockingSpots {
		log.Println("Can't dock, not enough open spots")
		return false
	}

	log.Println("can dock, yay")
	return true

}
