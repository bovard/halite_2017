package strat

import (
	"../hlt"
	"log"
	"math"
)

type ShipController struct {
	Ship   *hlt.Ship
	Past   []*hlt.Ship
	Id     int
	TargetPlanet int
}

func (self *ShipController) Update(ship *hlt.Ship) {
	self.Past = append(self.Past, self.Ship)
	self.Ship = ship
}

func (self *ShipController) MoveToPlanet(planet *hlt.Planet, gameMap *hlt.GameMap) hlt.Heading {
	return self.moveTo(&planet.Point, planet.Radius, gameMap)
}

func (self *ShipController) MoveToPoint(point *hlt.Point, gameMap *hlt.GameMap) hlt.Heading {
	return self.moveTo(point, 0.0, gameMap)
}

func (self *ShipController) MoveToShip(ship *hlt.Ship, gameMap *hlt.GameMap) hlt.Heading {
	return self.moveTo(&ship.Point, ship.Radius, gameMap)
}

func (self *ShipController) HeadingIsClear(mag int, angle float64, gameMap *hlt.GameMap) bool {
	v := hlt.CreateVector(mag, angle)
	for _, p := range(gameMap.Planets) {
		log.Println("Comparing with planet ", p.Id, " at loc ", p.Point)
		if self.Ship.WillCollideWith(&p.Entity, &v) {
			return false
		}
	}
	for _, s := range(gameMap.EnemyShips) {
		log.Println("Comparing with enemyShip ", s.Id, " at loc ", s.Point)
		if self.Ship.WillCollideWith(&s.Entity, &v) {
			return false
		}
	}
	var nv hlt.Vector
	for _, s := range(gameMap.MyShips) {
		log.Println("Comparing with friendly ship ", s.Id, " at loc ", s.Point, " with Vel ", s.NextVel)
		if self.Ship.Id == s.Id {
			continue
		}
		nv = v.Subtract(&s.NextVel)
		if self.Ship.WillCollideWith(&s.Entity, &nv) {
			return false
		}
	}
	return true
}

func (self *ShipController) UnsafeMoveToPoint(point *hlt.Point, gameMap *hlt.GameMap) hlt.Heading {
	log.Println("UnsafeMoveToPoint from ", self.Ship.Point, " to ", point)

	startSpeed := int(hlt.SHIP_MAX_SPEED)
	log.Println("setting start speed to ", startSpeed)
	baseAngle := self.Ship.Point.AngleTo(point)

	log.Println("Way is clear to target!")
	return hlt.CreateHeading(startSpeed, baseAngle)
}

func (self *ShipController) moveTo(point *hlt.Point, radius float64, gameMap *hlt.GameMap) hlt.Heading {
	log.Println("moveTo from ", self.Ship.Point, " to ", point, " with radius ", radius)

	maxTurn := (3 * math.Pi) / 2
	dTurn := math.Pi / 16

	startSpeed := int(math.Min(hlt.SHIP_MAX_SPEED, self.Ship.Point.DistanceTo(point)-radius-self.Ship.Radius-.05))
	log.Println("setting start speed to ", startSpeed)
	baseAngle := self.Ship.Point.AngleTo(point)

	if self.HeadingIsClear(startSpeed, baseAngle, gameMap) {
		log.Println("Way is clear to target!")
		return hlt.CreateHeading(startSpeed, baseAngle)
	}

	for speed := startSpeed; speed >= 1; speed -- {
		log.Println("Trying speed, ", speed)
		for turn := dTurn; turn <= maxTurn; turn += dTurn {
			log.Println("Trying turn, ", turn)
			intermediateTargetLeft := self.Ship.AddThrust(float64(speed), baseAngle+turn)
			obLeft := !self.HeadingIsClear(speed, baseAngle+turn, gameMap)
			intermediateTargetRight := self.Ship.AddThrust(float64(speed), baseAngle-turn)
			obRight := !self.HeadingIsClear(speed, baseAngle-turn, gameMap)
			if !obLeft && !obRight {
				if intermediateTargetLeft.SqDistanceTo(point) < intermediateTargetRight.SqDistanceTo(point) {
					return hlt.CreateHeading(speed, baseAngle+turn)
				} else {
					return hlt.CreateHeading(speed, baseAngle-turn)
				}
			} else if !obLeft {
				return hlt.CreateHeading(speed, baseAngle+turn)
			} else if !obRight {
				return hlt.CreateHeading(speed, baseAngle-turn)
			}
		}
	}
	return hlt.Heading {
		Magnitude: 0,
		Angle: 0,
	}
}


func (self *ShipController) combat(gameMap *hlt.GameMap, enemies []hlt.Entity) (ChlMessage, hlt.Heading) {
	enemiesInCombatRange := 0
	enemiesDockedInCombatRange := 0
	enemiesInThreatRange := 0
	alliesInCombatRange := 0
	alliesDockedInCombatRange := 0
	alliesInThreatRange := 0
	closestEnemyShip := gameMap.ShipLookup[enemies[0].Id]
	var message ChlMessage
	var heading hlt.Heading

	for _, s := range(append(gameMap.MyShips, gameMap.EnemyShips...)) {
		dist := self.Ship.DistanceToCollision(&s.Entity)
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
		}
	} 
	if (self.Ship.Health <= hlt.SHIP_MAX_HEALTH - hlt.SHIP_DAMAGE && closestEnemyShip.DockingStatus == hlt.DOCKED) {
		message = COMBAT_KILL_PRODUCTION
		heading = self.UnsafeMoveToPoint(&closestEnemyShip.Point, gameMap)
	} else if (alliesInCombatRange >= enemiesInCombatRange) {
		message = COMBAT_WE_OUTNUMBER
		//t := self.Ship.AddVector(&closestEnemyShip.Vel)
		//heading = self.MoveToPoint(&t, gameMap)
		heading = self.MoveToShip(closestEnemyShip, gameMap)
	} else if (alliesInCombatRange + 1 == enemiesInCombatRange ) {
		message = COMBAT_TIED
		//t := self.Ship.AddVector(&closestEnemyShip.Vel)
		//heading = self.MoveToPoint(&t, gameMap)
		heading = self.MoveToShip(closestEnemyShip, gameMap)
	} else {
		message = COMBAT_OUTNUMBERED
		//n := closestEnemyShip.Entity.Point.VectorTo(&self.Ship.Entity.Point)
		//t := self.Ship.AddVector(&n)
		//heading = self.MoveToPoint(&t, gameMap)
		heading = self.MoveToShip(closestEnemyShip, gameMap)
	}

	return message, heading

}

func (self *ShipController) Act(gameMap *hlt.GameMap) string {
	log.Println("Ship ", self.Id, " Act. Planet is ", self.TargetPlanet)
	enemies := gameMap.NearestEnemiesByDistance(*self.Ship)
	closestEnemy := enemies[0].Distance
	closestEnemyShip := gameMap.ShipLookup[enemies[0].Id]
	log.Println("ClosestEnemy is ", closestEnemy)
	heading := hlt.Heading {
		Magnitude: 0,
		Angle: 0,
	}
	message := NONE
	if self.TargetPlanet != -1 {
		planet := gameMap.PlanetsLookup[self.TargetPlanet]
		planetDist := self.Ship.Entity.DistanceToCollision(&planet.Entity)

		if closestEnemy < hlt.SHIP_MAX_ATTACK_RANGE - 1.0 {
			message, heading = self.combat(gameMap, enemies)
		} else if closestEnemy < 2 * hlt.SHIP_MAX_SPEED {
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet, enemy in min threshold")
			message = CANCELLED_PLANET_ASSIGNMENT_MIN
			heading = self.MoveToShip(closestEnemyShip, gameMap)
		} else if closestEnemy/2 < planetDist  {
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet, enemy too close")
			message = CANCELLED_PLANET_ASSIGNMENT_TOO_CLOSE
			heading = self.MoveToShip(closestEnemyShip, gameMap)
		} else if (planet.Owner > 0 && planet.Owner != gameMap.MyId){
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet, planet taken")
			message = CANCELLED_PLANET_ASSIGNMENT_PLANET_TAKEN
			heading = self.MoveToShip(closestEnemyShip, gameMap)
		} else {
			log.Println("Continuing with assigned planet")
			if self.Ship.CanDock(&planet) {
				log.Println("We can dock!")
				return self.Ship.Dock(&planet)
			} else {
				log.Println("moving toward planet", planet.Id)
				message = MOVING_TOWARD_PLANET
				heading = self.MoveToPlanet(&planet, gameMap)
			}
		}
	} else {
		message = MOVING_TOWARD_ENEMY
		heading = self.MoveToShip(closestEnemyShip, gameMap)
	}
	log.Println(heading)
	if heading.Magnitude > 0 {
		// TODO: figure out why these aren't the same thing!! :(
		self.Ship.NextVel = heading.ToVelocity()
		s := gameMap.ShipLookup[self.Ship.Id]
		s.NextVel = heading.ToVelocity()
	}
	return heading.ToMoveCmd(self.Ship, int(message))
}
