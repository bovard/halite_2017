package strat

import (
	"../hlt"
	"log"
)

type ShipController struct {
	Ship   *hlt.Ship
	Past   []*hlt.Ship
	Id     int
	Planet int
}

func (self *ShipController) Update(ship *hlt.Ship) {
	self.Past = append(self.Past, self.Ship)
	self.Ship = ship
}

func (self *ShipController) MoveToPlanet(planet *hlt.Planet, gameMap *hlt.GameMap) {

}

func (self *ShipController) Act(gameMap *hlt.GameMap) string {
	log.Println("Ship ", self.Id, " Act. Planet is ", self.Planet)
	enemies := gameMap.NearestEnemiesByDistance(*self.Ship)
	closetEnemy := enemies[0].Distance
	log.Println("ClosestEnemy is ", closetEnemy)
	heading := hlt.Heading {
		Magnitude: 0,
		Angle: 0,
	}
	message := NONE
	if self.Planet != -1 {
		planet := gameMap.PlanetsLookup[self.Planet]
		planetDist := self.Ship.Entity.DistanceToCollision(&planet.Entity)
		log.Println("Conditions for docking")
		log.Println(closetEnemy/2, " ?< ", planetDist)
		log.Println(closetEnemy, " ?< ", 2 * hlt.SHIP_MAX_SPEED)
		log.Println(planet.Owner)
		if closetEnemy/2 < planetDist || closetEnemy < 2 * hlt.SHIP_MAX_SPEED || (planet.Owner > 0 && planet.Owner != gameMap.MyId) {
			self.Planet = -1
			log.Println("Cancelling assigned planet, moving to enemy")
			message = CANCELLED_PLANET_ASSIGNMENT
			heading = self.Ship.BetterNavigate(&enemies[0], gameMap)
		} else {
			log.Println("Continuing with assigned planet")
			if self.Ship.CanDock(&planet) {
				log.Println("We can dock!")
				return self.Ship.Dock(&planet)
			} else {
				log.Println("moving toward planet", planet.Id)
				message = MOVING_TOWARD_PLANET
				heading = self.Ship.BetterNavigate(&planet.Entity, gameMap)
			}
		}
	} else {
		message = MOVING_TOWARD_ENEMY
		heading = self.Ship.BetterNavigate(&enemies[0], gameMap)
	}
	log.Println(heading)
	return heading.ToMoveCmd(self.Ship, int(message))
}
