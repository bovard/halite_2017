package ships

import (
	"../../hlt"
	"log"
)

func (self *ShipController) NormalSetTarget(gameMap *hlt.GameMap) {
	if self.TargetPlanet == -1 {
		self.Target = &self.Info.ClosestEnemyShip.Point
	} else if self.TargetPlanet != -1 {
		valid, _ := self.IsTargetPlanetStillMoveTarget(gameMap)
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

func (self *ShipController) NormalAct(gameMap *hlt.GameMap, turnComm *TurnComm) (ChlMessage, hlt.Heading) {
	heading := hlt.Heading{
		Magnitude: 0,
		Angle:     0,
	}
	message := NONE

	if self.Info.TotalEnemies > 0 {
		log.Println("We are in combat!")
		return self.combat(gameMap, turnComm)
	}

	if self.TargetPlanet != -1 {
		valid, pmess := self.IsTargetPlanetStillMoveTarget(gameMap)
		if !valid {
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet,", pmess)
			message = pmess
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		} else {
			planet := gameMap.PlanetLookup[self.TargetPlanet]
			log.Println("Continuing with assigned planet")
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
		}
	} else if self.Info.ClosestDockedEnemyShipDistance < 200 && turnComm.Chasing[self.Info.ClosestEnemyShip.Id] > 4 {
		message = TOO_MANY_CHASING_FAR_AWAY_MOVING_TO_DOCKED_ENEMY
		heading = self.MoveToShip(self.Info.ClosestDockedEnemyShip, gameMap)
	} else {
		message = MOVING_TOWARD_ENEMY
		turnComm.Chasing[self.Info.ClosestEnemyShip.Id] = turnComm.Chasing[self.Info.ClosestEnemyShip.Id] + 1
		heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
	}

	return message, heading
}

func (self *ShipController) IsTargetPlanetStillMoveTarget(gameMap *hlt.GameMap) (bool, ChlMessage) {
	message := NONE
	valid := true

	if self.TargetPlanet == -1 {
		return false, message
	}

	if _, ok := gameMap.PlanetLookup[self.TargetPlanet]; !ok {
		return false, message
	}

	planet := gameMap.PlanetLookup[self.TargetPlanet]
	planetDist := self.Ship.DistanceToCollision(&planet.Entity)

	if self.Info.ClosestEnemyShipDistance < 2*hlt.SHIP_MAX_SPEED {
		valid = false
		log.Println("Cancelling assigned planet, enemy in min threshold")
		message = CANCELLED_PLANET_ASSIGNMENT_MIN
	} else if self.Info.ClosestEnemyShipDistance/2 < planetDist {
		valid = false
		log.Println("Cancelling assigned planet, enemy too close")
		message = CANCELLED_PLANET_ASSIGNMENT_TOO_CLOSE
	} else if planet.Owner > 0 && planet.Owner != gameMap.MyId {
		valid = false
		log.Println("Cancelling assigned planet, planet taken")
		message = CANCELLED_PLANET_ASSIGNMENT_PLANET_TAKEN
	} else if self.Info.EnemyClosestPlanetDist < hlt.SHIP_MAX_SPEED {
		valid = false
		log.Println("Cancelling assigned planet, enemy planet too close")
		message = CANCELLED_PLANET_ASSIGNMENT_TOO_CLOSE_TO_ENEMEY_PLANET
	}
	return valid, message
}
