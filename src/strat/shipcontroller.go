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
	ShipNum int
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

	for _, p := range(self.Info.PossiblePlanetCollisions) {
		log.Println("Comparing with planet ", p.Id, " at loc ", p.Point)
		if self.Ship.WillCollideWith(&p.Entity, &v) {
			return false
		}
	}
	for _, s := range(self.Info.PossibleEnemyShipCollisions) {
		log.Println("Comparing with enemyShip ", s.Id, " at loc ", s.Point)
		if s.Id == target {
			continue
		}
		if self.Ship.WillCollideWith(&s.Entity, &v) {
			return false
		}
	}
	var nv hlt.Vector
	for _, s := range(self.Info.PossibleAlliedShipCollisions) {
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

func (self *ShipController) BetterHeadingIsClear(mag int, angle float64, gameMap *hlt.GameMap) bool {
	v := hlt.CreateVector(mag, angle)

	targetPos := self.Ship.Point.AddVector(&v)	
	if !gameMap.IsOnMap(&targetPos) {
		return false
	}

	for _, p := range(self.Info.PossiblePlanetCollisions) {
		log.Println("Comparing with planet ", p.Id, " at loc ", p.Point)
		if self.Ship.WillCollideWith(&p.Entity, &v) {
			return false
		}
	}
	for _, s := range(self.Info.PossibleEnemyShipCollisions) {
		log.Println("Comparing with enemyShip ", s.Id, " at loc ", s.Point)
		if self.Ship.WillCollideWith(&s.Entity, &v) {
			return false
		}
	}
	var nv hlt.Vector
	for _, s := range(self.Info.PossibleAlliedShipCollisions) {
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


func (self *ShipController) moveToLoop(pointTest func(*hlt.Point, float64, hlt.Point) bool, point *hlt.Point, radius float64, gameMap *hlt.GameMap, maxTurn float64, minSpeed int) hlt.Heading {
	startSpeed := int(math.Min(hlt.SHIP_MAX_SPEED, self.Ship.Point.DistanceTo(point)-radius-self.Ship.Radius-.05))
	log.Println("setting start speed to ", startSpeed)
	baseAngle := self.Ship.Point.AngleTo(point)
	dTurn := math.Pi / 16
	for speed := startSpeed; speed >= minSpeed; speed-- {
		log.Println("Trying speed, ", speed)
		for turn := dTurn; turn <= maxTurn; turn += dTurn {
			log.Println("Trying turn, ", turn)
			intermediateTargetLeft := self.Ship.AddThrust(float64(speed), baseAngle+turn)
			canGoLeft := pointTest(point, radius, intermediateTargetLeft) && self.BetterHeadingIsClear(speed, baseAngle+turn, gameMap)
			intermediateTargetRight := self.Ship.AddThrust(float64(speed), baseAngle-turn)
			canGoRight := pointTest(point, radius, intermediateTargetRight) && self.BetterHeadingIsClear(speed, baseAngle-turn, gameMap)
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


func (self *ShipController) moveTo(pointTest func(*hlt.Point, float64, hlt.Point) bool, point *hlt.Point, radius float64, gameMap *hlt.GameMap) hlt.Heading {
	log.Println("moveTo from ", self.Ship.Point, " to ", point, " with radius ", radius)

	firstTurn := math.Pi / 2 
	maxTurn := (3 * math.Pi) / 2

	startSpeed := int(math.Min(hlt.SHIP_MAX_SPEED, self.Ship.Point.DistanceTo(point)-radius-self.Ship.Radius-.05))
	log.Println("setting start speed to ", startSpeed)
	baseAngle := self.Ship.Point.AngleTo(point)

	if pointTest(point, radius, self.Ship.AddThrust(float64(startSpeed), baseAngle)) && self.BetterHeadingIsClear(startSpeed, baseAngle, gameMap) {
		log.Println("Way is clear to target!")
		return hlt.CreateHeading(startSpeed, baseAngle)
	}

	heading := self.moveToLoop(pointTest, point, radius, gameMap, firstTurn, int(math.Max(1, float64(startSpeed) - 1)))

	if (heading.Magnitude == 0) {
		heading = self.moveToLoop(pointTest, point, radius, gameMap, maxTurn, 1)
	}
	return heading 
}

func (self *ShipController) combat(gameMap *hlt.GameMap) (ChlMessage, hlt.Heading) {
	var message ChlMessage
	var heading hlt.Heading

	canKillSuicideOnProduction := self.Info.ClosestDockedEnemyShipDistance < hlt.SHIP_MAX_SPEED && self.HeadingIsClear(int(self.Info.ClosestDockedEnemyShipDistance + .5), self.Info.ClosestDockedEnemyShipDir, gameMap, self.Info.ClosestDockedEnemyShip.Id)

	if canKillSuicideOnProduction && self.Ship.Health <= 2.0 * hlt.SHIP_DAMAGE * (float64(self.Info.EnemiesInCombatRange) + float64(self.Info.EnemiesInThreatRange)) && self.Info.ClosestDockedEnemyShip.Health > hlt.SHIP_DAMAGE * float64(self.Info.AlliesInCombatRange + 1)  {
		message = COMBAT_SUICIDE_ON_PRODUCTION_DUE_TO_LOWER_HEALTH
		heading = self.UnsafeMoveToPoint(&self.Info.ClosestDockedEnemyShip.Point, gameMap, true)
	} else if self.Info.EnemiesInCombatRange != 0 && (self.Info.EnemiesInThreatRange + self.Info.EnemiesInActiveThreatRange) > 0 && self.Info.ClosestAlliedShipDistance < 100 {
		p := self.Info.ClosestAlliedShip.AddVector(&self.Info.ClosestAlliedShip.NextVel)
		message = MOVING_TO_ALLY
		heading = self.MoveToPoint(&p, gameMap)
	} else {
		message = MOVING_TOWARD_ENEMY
		heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
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
	if self.Info.EnemiesInCombatRange > 0 || self.Info.EnemiesInThreatRange > 0 || self.Info.EnemiesInActiveThreatRange > 0  {
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
