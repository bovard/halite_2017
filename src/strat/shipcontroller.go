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
	STUPID_RUN_AWAY_META
)

type ShipController struct {
	Ship         *hlt.Ship
	Past         []*hlt.Ship
	Id           int
	TargetPlanet int
	Mission      Mission
	Info         ShipTurnInfo
	ShipNum      int
	Distance     float64
	Target       *hlt.Point
}

type byDistSc []*ShipController

func (a byDistSc) Len() int           { return len(a) }
func (a byDistSc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDistSc) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

func (self *ShipController) Update(ship *hlt.Ship) {
	self.Past = append(self.Past, self.Ship)
	self.Ship = ship
}

func (self *ShipController) MoveToPlanet(planet *hlt.Planet, gameMap *hlt.GameMap) hlt.Heading {
	closestDockingPoint := self.Ship.ClosestPointTo(&planet.Entity, hlt.SHIP_DOCKING_RADIUS)
	return self.moveTo(genericPointTest, &closestDockingPoint, 0, gameMap)
}

func (self *ShipController) MoveToPoint(point *hlt.Point, gameMap *hlt.GameMap) hlt.Heading {
	return self.moveTo(genericPointTest, point, 0.0, gameMap)
}

func (self *ShipController) MoveToShip(ship *hlt.Ship, gameMap *hlt.GameMap) hlt.Heading {
	return self.moveTo(genericPointTest, &ship.Point, ship.Radius, gameMap)
}

func (self *ShipController) MoveToDockingRange(planet *hlt.Planet, gameMap *hlt.GameMap) hlt.Heading {
	return self.moveTo(inDockingRangePointTest, &planet.Point, planet.Radius, gameMap)
}

func genericPointTest(point *hlt.Point, radius float64, newPoint hlt.Point) bool {
	return true
}

func inDockingRangePointTest(planetPoint *hlt.Point, planetRadius float64, pointToEval hlt.Point) bool {
	return planetPoint.DistanceTo(&pointToEval) <= hlt.SHIP_DOCKING_RADIUS+hlt.SHIP_RADIUS+planetRadius
}

func (self *ShipController) HeadingIsClear(mag int, angle float64, gameMap *hlt.GameMap, target int) bool {
	v := hlt.CreateRoundedVector(mag, angle)

	targetPos := self.Ship.Point.AddVector(&v)
	if !gameMap.IsOnMap(&targetPos) {
		return false
	}

	for _, p := range self.Info.PossiblePlanetCollisions {
		log.Println("Comparing with planet ", p.Id, " at loc ", p.Point)
		if self.Ship.WillCollideWith(&p.Entity, &v) {
			return false
		}
	}
	var nv hlt.Vector
	for _, s := range self.Info.PossibleEnemyShipCollisions {
		log.Println("Comparing with enemyShip ", s.Id, " at loc ", s.Point)
		log.Println("Enemey ship LastPos", s.LastPos, "velocity", s.Vel)
		if s.Id == target {
			continue
		}
		if self.Ship.WillCollideWith(&s.Entity, &v) {
			return false
		}
		if s.DockingStatus == hlt.UNDOCKED {
			// check if the closest enemy ship moves toward us, will will collide?
			if s.Id == self.Info.ClosestEnemyShip.Id && self.Info.ClosestEnemyShipDistance <= 2*hlt.SHIP_MAX_SPEED {
				tm := hlt.CreateVector(int(self.Info.ClosestEnemyShipDistance), self.Info.ClosestEnemyShip.AngleTo(&self.Ship.Point))
				nv = v.Subtract(&tm)
				if self.Ship.WillCollideWith(&s.Entity, &nv) {
					return false
				}
			}
			if s.Vel.Magnitude() > 0 {
				log.Println("DOING IT")
				nv = v.Subtract(&s.Vel)
				if self.Ship.WillCollideWith(&s.Entity, &nv) {
					return false
				}
			}
		}

	}
	for _, s := range self.Info.PossibleAlliedShipCollisions {
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
			canGoLeft := pointTest(point, radius, intermediateTargetLeft) && self.HeadingIsClear(speed, baseAngle+turn, gameMap, -1)
			intermediateTargetRight := self.Ship.AddThrust(float64(speed), baseAngle-turn)
			canGoRight := pointTest(point, radius, intermediateTargetRight) && self.HeadingIsClear(speed, baseAngle-turn, gameMap, -1)
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

	return hlt.Heading{
		Magnitude: 0,
		Angle:     0,
	}
}

func (self *ShipController) moveTo(pointTest func(*hlt.Point, float64, hlt.Point) bool, point *hlt.Point, radius float64, gameMap *hlt.GameMap) hlt.Heading {
	log.Println("moveTo from ", self.Ship.Point, " to ", point, " with radius ", radius)

	firstTurn := math.Pi / 2
	maxTurn := (3 * math.Pi) / 2

	startSpeed := int(math.Min(hlt.SHIP_MAX_SPEED, self.Ship.Point.DistanceTo(point)-radius-self.Ship.Radius-.05))
	log.Println("setting start speed to ", startSpeed)
	baseAngle := self.Ship.Point.AngleTo(point)

	if pointTest(point, radius, self.Ship.AddThrust(float64(startSpeed), baseAngle)) && self.HeadingIsClear(startSpeed, baseAngle, gameMap, -1) {
		log.Println("Way is clear to target!")
		return hlt.CreateHeading(startSpeed, baseAngle)
	}

	heading := self.moveToLoop(pointTest, point, radius, gameMap, firstTurn, int(math.Max(1, float64(startSpeed)-1)))

	if heading.Magnitude == 0 {
		heading = self.moveToLoop(pointTest, point, radius, gameMap, maxTurn, 1)
	}
	return heading
}

func (self *ShipController) combat(gameMap *hlt.GameMap) (ChlMessage, hlt.Heading) {
	var message ChlMessage
	var heading hlt.Heading

	canKillSuicideOnProduction := self.Info.ClosestDockedEnemyShipDistance < hlt.SHIP_MAX_SPEED && self.HeadingIsClear(int(self.Info.ClosestDockedEnemyShipDistance+.5), self.Info.ClosestDockedEnemyShipDir, gameMap, self.Info.ClosestDockedEnemyShip.Id)
	canKillSuicideOnNearestEnemy := self.Info.ClosestEnemyShipDistance < hlt.SHIP_MAX_SPEED && self.HeadingIsClear(int(self.Info.ClosestEnemyShipDistance+.5), self.Info.ClosestEnemyShipDir, gameMap, self.Info.ClosestEnemyShip.Id)

	if canKillSuicideOnProduction && self.Ship.Health <= 2.0*hlt.SHIP_DAMAGE*(float64(self.Info.EnemiesInCombatRange)+float64(self.Info.EnemiesInThreatRange)) && self.Info.ClosestDockedEnemyShip.Health > hlt.SHIP_DAMAGE/2 {
		message = COMBAT_SUICIDE_ON_PRODUCTION_DUE_TO_LOWER_HEALTH
		heading = self.UnsafeMoveToPoint(&self.Info.ClosestDockedEnemyShip.Point, gameMap, true)
	} else if canKillSuicideOnNearestEnemy && self.Info.AlliesInCombatRange == 0 && self.Info.EnemiesInCombatRange > 0 && int(self.Ship.Health/hlt.SHIP_MAX_HEALTH) < int(self.Info.ClosestEnemyShip.Health/hlt.SHIP_MAX_HEALTH) {
		message = COMBAT_SUICIDE_DUE_TO_LOWER_HEALTH
		heading = self.UnsafeMoveToPoint(&self.Info.ClosestEnemyShip.Point, gameMap, true)
	} else if self.Info.ClosestEnemyShip.Vel.Magnitude() > 0 && !self.Info.ClosestEnemyShipClosingDistance && self.Info.ClosestDockedEnemyShipDistance < 2*hlt.SHIP_MAX_SPEED {
		message = CHASING_DOWN_ENEMY
		heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		if self.HeadingIsClear(int(hlt.SHIP_MAX_SPEED), heading.GetAngleInRads(), gameMap, self.Info.ClosestEnemyShip.Id) {
			heading.Magnitude = int(hlt.SHIP_MAX_SPEED)
		}
	} else if self.Info.EnemiesInCombatRange > 1 {
		message = MOVING_TO_BETTER_LOCAL
		// free to move to opimal spot
		if math.Abs(self.Info.ClosestDockedEnemyShipDir-self.Ship.AngleTo(&self.Info.EnemiesByDist[1].Point)) < 2*math.Pi/3 {
			v := self.Info.EnemiesByDist[1].VectorTo(&self.Info.ClosestEnemyShip.Point)
			v = v.RescaleToMagFloat(hlt.SHIP_MAX_ATTACK_RANGE + .95)
			p := self.Info.ClosestEnemyShip.AddVector(&v)
			heading = self.MoveToPoint(&p, gameMap)
		} else {
			fromCloset := self.Info.ClosestEnemyShip.VectorTo(&self.Ship.Point)
			fromOther := self.Info.EnemiesByDist[1].VectorTo(&self.Ship.Point)
			v := fromCloset.Add(&fromOther)
			v = v.RescaleToMagFloat(hlt.SHIP_MAX_SPEED + .1)
			p := self.Ship.AddVector(&v)
			heading = self.MoveToPoint(&p, gameMap)
		}
	} else if self.Info.EnemiesInCombatRange == 1 && self.Info.EnemiesInThreatRange == 1 {
		message = MOVING_TO_MAX_RANGE
		v := self.Info.ClosestEnemyShip.VectorTo(&self.Ship.Point)
		v = v.RescaleToMagFloat(hlt.SHIP_MAX_SPEED + .1)
		p := self.Info.ClosestEnemyShip.AddVector(&v)
		heading = self.MoveToPoint(&p, gameMap)
	} else {
		log.Println("TOTAL enemies/allies", self.Info.TotalEnemies, self.Info.TotalAllies)
		message = MOVING_TOWARD_ENEMY
		heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
	}

	return message, heading
}

func (self *ShipController) UpdateInfo(gameMap *hlt.GameMap) {
	self.Info = CreateShipTurnInfo(self.Ship, gameMap)
}


func (self *ShipController) runAway(gameMap *hlt.GameMap) (ChlMessage, hlt.Heading) {
	heading := hlt.Heading{
		Magnitude: 0,
		Angle:     0,
	}
	message := RUN_AWAY

}


func nextCorner(current hlt.Point, gameMap *hlt.GameMap) {
	if current.X == 0 && current.Y == 0 {
		return hlt.Point {
			X: gameMap.Width,
			Y: 0,
		} 
	} else if current.X == gameMap.Width && current.Y == 0 {
		return hlt.Point {
			X: gameMap.Width,
			Y: gameMap.Height,
		}
	} else if current.X == gameMap.Width && current.Y == gameMap.Height {
		return hlt.Point {
			X: 0,
			Y: gameMap.Height,
		}
	} else if current.X == 0 && current.Y == gameMap.Height {
		return hlt.Point {
			X: 0,
			Y: 0,
		}
	}
}

func (self *ShipController) stupidRunAwayMeta(gameMap *hlt.GameMap) (ChlMessage, hlt.Heading) {
	heading := hlt.Heading{
		Magnitude: 0,
		Angle:     0,
	}
	message := HIDE_WE_ARE_LOSING	



	return message, heading
}


func (self *ShipController) SetTarget(gameMap *hlt.GameMap) {
	if self.MISSION == MISSION_FOUND_PLANET {
		planet := gameMap.PlanetLookup[self.TargetPlanet]
		self.Target = &planet.Point
	} else if self.TargetPlanet != -1 {
		planet := gameMap.PlanetLookup[self.TargetPlanet]
		self.Target = &planet.Point
	}
}

func (self *ShipController) Act(gameMap *hlt.GameMap) string {

	log.Println("Ship ", self.Id, " Act. Planet is ", self.TargetPlanet)
	log.Println("ClosestEnemy is ", self.Info.ClosestEnemyShipDistance)

	heading := hlt.Heading{
		Magnitude: 0,
		Angle:     0,
	}
	message := NONE
	if self.Mission == STUPID_RUN_AWAY_META {
		message, heading = self.stupidRunAwayMeta(gameMap)
	} else if self.TotalEnemies > 0 {
		message, heading = self.combat(gameMap)
	} else if self.Mission == MISSION_FOUND_PLANET {
		planet := gameMap.PlanetLookup[self.TargetPlanet]
		log.Println("Continuing with assigned planet")
		if self.Ship.CanDock(planet) {
			log.Println("We can dock!")
			return self.Ship.Dock(planet)
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
	} else if self.TargetPlanet != -1 {
		planet := gameMap.PlanetLookup[self.TargetPlanet]
		planetDist := self.Ship.Entity.DistanceToCollision(&planet.Entity)

		if self.Info.ClosestEnemyShipDistance < 2*hlt.SHIP_MAX_SPEED {
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet, enemy in min threshold")
			message = CANCELLED_PLANET_ASSIGNMENT_MIN
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		} else if self.Info.ClosestEnemyShipDistance/2 < planetDist {
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet, enemy too close")
			message = CANCELLED_PLANET_ASSIGNMENT_TOO_CLOSE
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		} else if planet.Owner > 0 && planet.Owner != gameMap.MyId {
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet, planet taken")
			message = CANCELLED_PLANET_ASSIGNMENT_PLANET_TAKEN
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		} else if self.Info.EnemyClosestPlanetDist < hlt.SHIP_MAX_SPEED {
			self.TargetPlanet = -1
			log.Println("Cancelling assigned planet, enemy planet too close")
			message = CANCELLED_PLANET_ASSIGNMENT_TOO_CLOSE_TO_ENEMEY_PLANET
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
		} else {
			log.Println("Continuing with assigned planet")
			if self.Ship.CanDock(planet) {
				log.Println("We can dock!")
				return self.Ship.Dock(planet)
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
	} else {
		message = MOVING_TOWARD_ENEMY
		heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
	}
	log.Println(heading)
	if heading.Magnitude > 0 {
		s := gameMap.ShipLookup[self.Ship.Id]
		log.Println("Compare pointers")
		log.Println(&self.Ship, &s)
		// TODO: figure out why these aren't the same thing!! :(
		self.Ship.NextVel = heading.ToVelocity()
		s.NextVel = heading.ToVelocity()
	}
	return heading.ToMoveCmd(self.Ship, int(message))
}
