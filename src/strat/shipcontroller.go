package strat

import (
	"../hlt"
	"log"
	"math"
)


type Mission int

const (
	MISSION_NORMAL Mission = iota
	MISSION_FOUND_PLANET
)

type ShipController struct {
	Ship   *hlt.Ship
	Past   []*hlt.Ship
	Id     int
	TargetPlanet int
	Mission Mission
	Info ShipTurnInfo
}

func (self *ShipController) Update(ship *hlt.Ship) {
	self.Past = append(self.Past, self.Ship)
	self.Ship = ship
}

func (self *ShipController) MoveToPlanet(planet *hlt.Planet, gameMap *hlt.GameMap) hlt.Heading {
	return self.moveTo(genericPointTest, &planet.Point, planet.Radius, gameMap)
}

func (self *ShipController) MoveToPoint(point *hlt.Point, gameMap *hlt.GameMap) hlt.Heading {
	return self.moveTo(genericPointTest, point, 0.0, gameMap)
}

func (self *ShipController) MoveToShip(ship *hlt.Ship, gameMap *hlt.GameMap) hlt.Heading {
	return self.moveTo(genericPointTest, &ship.Point, ship.Radius, gameMap)
}

func (self * ShipController) MoveToDockingRange(planet *hlt.Planet, gameMap *hlt.GameMap) hlt.Heading {
	return self.moveTo(inDockingRangePointTest, &planet.Point, planet.Radius, gameMap)
}

func genericPointTest(point *hlt.Point, radius float64, newPoint hlt.Point) bool {
	return true;
}

func inDockingRangePointTest(planetPoint *hlt.Point, planetRadius float64, pointToEval hlt.Point) bool {
	return planetPoint.DistanceTo(&pointToEval) <= hlt.SHIP_DOCKING_RADIUS + hlt.SHIP_RADIUS + planetRadius
}


func (self *ShipController) HeadingIsClear(mag int, angle float64, gameMap *hlt.GameMap, target int) bool {
	v := hlt.CreateVector(mag, angle)

	targetPos := self.Ship.Point.AddVector(&v)	
	if !gameMap.IsOnMap(&targetPos) {
		return false
	}

	for _, p := range(gameMap.Planets) {
		log.Println("Comparing with planet ", p.Id, " at loc ", p.Point)
		if self.Ship.WillCollideWith(&p.Entity, &v) {
			return false
		}
	}
	for _, s := range(gameMap.EnemyShips) {
		log.Println("Comparing with enemyShip ", s.Id, " at loc ", s.Point)
		if s.Id == target {
			continue
		}
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

func (self *ShipController) BetterHeadingIsClear(mag int, angle float64, gameMap *hlt.GameMap, possiblePlanetCollisions []hlt.Planet, possibleEnemyShipCollisions []*hlt.Ship, possibleAlliedShipCollisions []*hlt.Ship) bool {
	v := hlt.CreateVector(mag, angle)

	targetPos := self.Ship.Point.AddVector(&v)	
	if !gameMap.IsOnMap(&targetPos) {
		return false
	}

	for _, p := range(possiblePlanetCollisions) {
		log.Println("Comparing with planet ", p.Id, " at loc ", p.Point)
		if self.Ship.WillCollideWith(&p.Entity, &v) {
			return false
		}
	}
	for _, s := range(possibleEnemyShipCollisions) {
		log.Println("Comparing with enemyShip ", s.Id, " at loc ", s.Point)
		if self.Ship.WillCollideWith(&s.Entity, &v) {
			return false
		}
	}
	var nv hlt.Vector
	for _, s := range(possibleAlliedShipCollisions) {
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

func (self *ShipController) UnsafeMoveToPoint(point *hlt.Point, gameMap *hlt.GameMap, maxSpeed bool) hlt.Heading {
	log.Println("UnsafeMoveToPoint from ", self.Ship.Point, " to ", point)

	startSpeed := int(math.Min(hlt.SHIP_MAX_SPEED, self.Ship.Point.DistanceTo(point)))
	if maxSpeed {
		startSpeed = int(hlt.SHIP_MAX_SPEED)
	}
	log.Println("setting start speed to ", startSpeed)
	baseAngle := self.Ship.Point.AngleTo(point)

	return hlt.CreateHeading(startSpeed, baseAngle)
}

func (self *ShipController) moveTo(pointTest func(*hlt.Point, float64, hlt.Point) bool, point *hlt.Point, radius float64, gameMap *hlt.GameMap) hlt.Heading {
	log.Println("moveTo from ", self.Ship.Point, " to ", point, " with radius ", radius)

	// TODO: why can't we do this with pointers :(
	// when using *hlt.Planet it always defaults to the last element
	// i guess since slices are already pointers?
	possiblePlanetCollisions := []hlt.Planet{}
	for _, p := range(gameMap.Planets) {
		log.Println(p.Id, self.Ship.DistanceToCollision(&p.Entity), " ?<= ", hlt.SHIP_MAX_SPEED)
		if self.Ship.DistanceToCollision(&p.Entity) <= hlt.SHIP_MAX_SPEED {
			possiblePlanetCollisions = append(possiblePlanetCollisions, p)
		}
	}
	log.Println(len(possiblePlanetCollisions))
	for _, p := range(possiblePlanetCollisions) {
		log.Println(p.Id)
	}

	possibleEnemyShipCollisions := []*hlt.Ship{}
	for _, s := range(gameMap.EnemyShips) {
		if self.Ship.DistanceToCollision(&s.Entity) <= 2 * hlt.SHIP_MAX_SPEED {
			possibleEnemyShipCollisions = append(possibleEnemyShipCollisions, s)
		}
	}

	possibleAlliedShipCollisions := []*hlt.Ship{}
	for _, s := range(gameMap.MyShips) {
		if self.Ship.DistanceToCollision(&s.Entity) <= 2 * hlt.SHIP_MAX_SPEED {
			possibleAlliedShipCollisions = append(possibleAlliedShipCollisions, s)
		}
	}


	maxTurn := (3 * math.Pi) / 2
	dTurn := math.Pi / 16

	startSpeed := int(math.Min(hlt.SHIP_MAX_SPEED, self.Ship.Point.DistanceTo(point)-radius-self.Ship.Radius-.05))
	log.Println("setting start speed to ", startSpeed)
	baseAngle := self.Ship.Point.AngleTo(point)

	if pointTest(point, radius, self.Ship.AddThrust(float64(startSpeed), baseAngle)) && self.BetterHeadingIsClear(startSpeed, baseAngle, gameMap, possiblePlanetCollisions, possibleEnemyShipCollisions, possibleAlliedShipCollisions) {
		log.Println("Way is clear to target!")
		return hlt.CreateHeading(startSpeed, baseAngle)
	}

	for speed := startSpeed; speed >= 1; speed -- {
		log.Println("Trying speed, ", speed)
		for turn := dTurn; turn <= maxTurn; turn += dTurn {
			log.Println("Trying turn, ", turn)
			intermediateTargetLeft := self.Ship.AddThrust(float64(speed), baseAngle+turn)
			canGoLeft := pointTest(point, radius, intermediateTargetLeft) && self.BetterHeadingIsClear(speed, baseAngle+turn, gameMap, possiblePlanetCollisions, possibleEnemyShipCollisions, possibleAlliedShipCollisions)
			intermediateTargetRight := self.Ship.AddThrust(float64(speed), baseAngle-turn)
			canGoRight := pointTest(point, radius, intermediateTargetRight) && self.BetterHeadingIsClear(speed, baseAngle-turn, gameMap, possiblePlanetCollisions, possibleEnemyShipCollisions, possibleAlliedShipCollisions)
			if canGoLeft && canGoRight {
				if intermediateTargetLeft.SqDistanceTo(point) < intermediateTargetRight.SqDistanceTo(point) {
					return hlt.CreateHeading(speed, baseAngle+turn)
				} else {
					return hlt.CreateHeading(speed, baseAngle-turn)
				}
			} else if canGoLeft {
				return hlt.CreateHeading(speed, baseAngle+turn)
			} else if canGoRight {
				return hlt.CreateHeading(speed, baseAngle-turn)
			}
		}
	}
	return hlt.Heading {
		Magnitude: 0,
		Angle: 0,
	}
}


func (self *ShipController) combat(gameMap *hlt.GameMap) (ChlMessage, hlt.Heading) {
	var message ChlMessage
	var heading hlt.Heading

	if (((self.Ship.Health < hlt.SHIP_MAX_HEALTH && self.Info.AlliesInCombatRange == 0)  || self.Info.EnemiesInCombatRange > self.Info.AlliesInCombatRange) && self.Info.ClosestDockedEnemyShipDistance < 5.0 && self.HeadingIsClear(int(self.Info.ClosestDockedEnemyShipDistance + .5), self.Info.ClosestDockedEnemyShipDir, gameMap, self.Info.ClosestDockedEnemyShip.Id) ) {
		message = COMBAT_KILL_PRODUCTION
		heading = self.UnsafeMoveToPoint(&self.Info.ClosestEnemyShip.Point, gameMap, true)
	} else if (self.Info.ClosestEnemyShipDistance <= 2 && int(self.Ship.Health/hlt.SHIP_MAX_HEALTH) < int(self.Info.ClosestEnemyShip.Health/hlt.SHIP_MAX_HEALTH) && self.Info.ClosestDockedEnemyShipDistance < 5.0 && self.HeadingIsClear(int(self.Info.ClosestDockedEnemyShipDistance + .5), self.Info.ClosestDockedEnemyShipDir, gameMap, self.Info.ClosestDockedEnemyShip.Id) ) {
		message = COMBAT_SUICIDE_ON_PRODUCTION_DUE_TO_LOWER_HEALTH
		heading = self.UnsafeMoveToPoint(&self.Info.ClosestDockedEnemyShip.Point, gameMap, true)
	} else if (self.Info.ClosestEnemyShipDistance <= 2 && int(self.Ship.Health/hlt.SHIP_MAX_HEALTH) < int(self.Info.ClosestEnemyShip.Health/hlt.SHIP_MAX_HEALTH) && self.HeadingIsClear(int(self.Info.ClosestEnemyShipDistance + .5), self.Info.ClosestEnemyShipDir, gameMap, self.Info.ClosestEnemyShip.Id) ) {
		message = COMBAT_SUICIDE_DUE_TO_LOWER_HEALTH
		heading = self.UnsafeMoveToPoint(&self.Info.ClosestEnemyShip.Point, gameMap, false)
	} else if (self.Info.AlliesInCombatRange >= self.Info.EnemiesInCombatRange) {
		message = COMBAT_WE_OUTNUMBER
		//t := self.Ship.AddVector(&self.Info.ClosestEnemyShip.Vel)
		//heading = self.MoveToPoint(&t, gameMap)
		heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
	} else if (self.Info.AlliesInCombatRange + 1 == self.Info.EnemiesInCombatRange ) {
		if (self.Info.ClosestEnemyShipDistance <= 2 && int(self.Ship.Health/hlt.SHIP_MAX_HEALTH) < int(self.Info.ClosestEnemyShip.Health/hlt.SHIP_MAX_HEALTH) && self.HeadingIsClear(int(self.Info.ClosestEnemyShipDistance + .5), self.Info.ClosestEnemyShipDir, gameMap, self.Info.ClosestEnemyShip.Id)) {
			message = COMBAT_TIED_SUICIDE_TO_GAIN_VALUE
			heading = self.UnsafeMoveToPoint(&self.Info.ClosestEnemyShip.Point, gameMap, false)
		} else if self.Info.ClosestDockedEnemyShipDistance < 2 * hlt.SHIP_MAX_SPEED && self.Info.EnemiesInThreatRange == 0 && self.Info.EnemiesInCombatRange == 1 {
			message = COMBAT_TIED_GOING_TO_HURT_PRODUCTION
			heading = self.MoveToShip(self.Info.ClosestDockedEnemyShip, gameMap)
		}else {
			message = COMBAT_TIED
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		}
	} else {
		if (self.Info.AlliedClosestPlanetDist > 2 * self.Info.EnemyClosestPlanetDist) {
			message = COMBAT_OUTNUMBERED_AND_FAR_FROM_HOME
			away := self.Info.ClosestEnemyShip.Entity.Point.VectorTo(&self.Ship.Entity.Point)
			away = away.RescaleToMag(int(hlt.SHIP_MAX_SPEED))
			t := self.Ship.AddVector(&away)
			heading = self.MoveToPoint(&t, gameMap)

		} else {
			message = COMBAT_OUTNUMBERED
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		}
	}

	return message, heading

}

func (self *ShipController) Act(gameMap *hlt.GameMap) string {
	self.Info = CreateShipTurnInfo(self.Ship, gameMap)

	log.Println("Ship ", self.Id, " Act. Planet is ", self.TargetPlanet)
	log.Println("ClosestEnemy is ", self.Info.ClosestEnemyShipDistance)

	heading := hlt.Heading {
		Magnitude: 0,
		Angle: 0,
	}
	message := NONE
	if self.Info.ClosestEnemyShipDistance <= hlt.SHIP_MAX_ATTACK_RANGE - 1.0 {
			message, heading = self.combat(gameMap)
	} else if self.Mission == MISSION_FOUND_PLANET {
		planet := gameMap.PlanetsLookup[self.TargetPlanet]
		log.Println("Continuing with assigned planet")
		if self.Ship.CanDock(&planet) {
			log.Println("We can dock!")
			return self.Ship.Dock(&planet)
		}  
		h := self.MoveToDockingRange(&planet, gameMap)
		if h.Magnitude > 0 {
			log.Println("can move to docking range of", planet.Id)
			message = MOVE_TO_DOCKING
			heading = h
		} else {
			log.Println("moving toward planet", planet.Id)
			message = MOVING_TOWARD_PLANET
			heading = self.MoveToPlanet(&planet, gameMap)
		}	
	} else if self.TargetPlanet != -1 {
		planet := gameMap.PlanetsLookup[self.TargetPlanet]
		planetDist := self.Ship.Entity.DistanceToCollision(&planet.Entity)

		if self.Info.ClosestEnemyShipDistance < 2 * hlt.SHIP_MAX_SPEED {
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet, enemy in min threshold")
			message = CANCELLED_PLANET_ASSIGNMENT_MIN
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		} else if self.Info.ClosestEnemyShipDistance / 2 < planetDist  {
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet, enemy too close")
			message = CANCELLED_PLANET_ASSIGNMENT_TOO_CLOSE
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		} else if (planet.Owner > 0 && planet.Owner != gameMap.MyId){
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet, planet taken")
			message = CANCELLED_PLANET_ASSIGNMENT_PLANET_TAKEN
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		} else if (self.Info.EnemyClosestPlanetDist < hlt.SHIP_MAX_SPEED) {
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet, enemy planet too close")
			message = CANCELLED_PLANET_ASSIGNMENT_TOO_CLOSE_TO_ENEMEY_PLANET
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		} else {
			log.Println("Continuing with assigned planet")
			if self.Ship.CanDock(&planet) {
				log.Println("We can dock!")
				return self.Ship.Dock(&planet)
			}  
			h := self.MoveToDockingRange(&planet, gameMap)
			if h.Magnitude > 0 {
				log.Println("can move to docking range of", planet.Id)
				message = MOVE_TO_DOCKING
				heading = h
			} else {
				log.Println("moving toward planet", planet.Id)
				message = MOVING_TOWARD_PLANET
				heading = self.MoveToPlanet(&planet, gameMap)
			}
		}
	} else {
		message = MOVING_TOWARD_ENEMY
		heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
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
