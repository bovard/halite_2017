package ships

import (
	"../../hlt"
	"log"
	"math"
)

type Mission int

const (
	MISSION_NORMAL Mission = iota
	MISSION_SETTLER
	MISSION_RUN_AWAY
	MISSION_SNEAKY
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

type ByDistSc []*ShipController

func (a ByDistSc) Len() int           { return len(a) }
func (a ByDistSc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistSc) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

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
		if !s.IsAliveNextTurn() {
			continue
		}
		if self.Ship.WillCollideWith(&s.Entity, &v) {
			return false
		}
		if s.DockingStatus == hlt.UNDOCKED {
			// check if the closest enemy ship moves toward us, will will collide?
			if s.Id == self.Info.ClosestEnemy.Ship.Id && self.Info.ClosestEnemy.Distance <= 2*hlt.SHIP_MAX_SPEED {
				tm := hlt.CreateVector(int(self.Info.ClosestEnemy.Distance), self.Info.ClosestEnemy.Ship.AngleTo(&self.Ship.Point))
				nv = v.Subtract(&tm)
				if self.Ship.WillCollideWith(&s.Entity, &nv) {
					return false
				}
			}
			if s.Vel.Magnitude() > 0 {
				nv = v.Subtract(&s.Vel)
				if self.Ship.WillCollideWith(&s.Entity, &nv) {
					return false
				}
			}
		}

	}
	for _, s := range self.Info.PossibleAlliedShipCollisions {
		log.Println("Comparing with friendly ship ", s.Id, " at loc ", s.Point, " with Vel ", s.NextVel)
		if !s.IsAliveNextTurn() {
			log.Println("It's not alive next turn!")
			continue
		}
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

func (self *ShipController) combat(gameMap *hlt.GameMap, turnComm *TurnComm) (ChlMessage, hlt.Heading) {
	var message ChlMessage
	var heading hlt.Heading

	canKillSuicideOnProduction := self.Info.ClosestDockedEnemy.Distance < hlt.SHIP_MAX_SPEED && self.Info.ClosestDockedEnemy.Ship.IsAliveNextTurn() && self.HeadingIsClear(int(self.Info.ClosestDockedEnemy.Distance+.5), self.Info.ClosestDockedEnemy.Direction, gameMap, self.Info.ClosestDockedEnemy.Ship.Id)
	canKillSuicideOnNearestEnemy := self.Info.ClosestNonDockedEnemy.Distance < hlt.SHIP_MAX_SPEED && self.Info.ClosestNonDockedEnemy.Ship.IsAliveNextTurn() && self.HeadingIsClear(int(self.Info.ClosestNonDockedEnemy.Distance+.5), self.Info.ClosestNonDockedEnemy.Direction, gameMap, self.Info.ClosestNonDockedEnemy.Ship.Id)

	if canKillSuicideOnProduction && self.Ship.Health <= 2.0*hlt.SHIP_DAMAGE*(float64(self.Info.EnemiesInCombatRange)+float64(self.Info.EnemiesInThreatRange)) && self.Info.ClosestDockedEnemy.Ship.Health > hlt.SHIP_DAMAGE/2 {
		message = COMBAT_SUICIDE_ON_PRODUCTION_DUE_TO_LOWER_HEALTH
		heading = self.UnsafeMoveToPoint(&self.Info.ClosestDockedEnemy.Ship.Point, gameMap, true)
	} else if canKillSuicideOnNearestEnemy && self.Info.AlliesInCombatRange == 0 && self.Info.EnemiesInCombatRange > 0 && int(self.Ship.Health/hlt.SHIP_MAX_HEALTH) < int(self.Info.ClosestEnemy.Ship.Health/hlt.SHIP_MAX_HEALTH) {
		message = COMBAT_SUICIDE_DUE_TO_LOWER_HEALTH
		heading = self.UnsafeMoveToPoint(&self.Info.ClosestEnemy.Ship.Point, gameMap, true)
	} else if self.Info.ClosestEnemy.Ship.Vel.Magnitude() > 0 && !self.Info.ClosestEnemyShipClosingDistance && self.Info.ClosestDockedEnemy.Distance < 2*hlt.SHIP_MAX_SPEED {
		message = CHASING_DOWN_ENEMY
		heading = self.MoveToShip(self.Info.ClosestEnemy.Ship, gameMap)
		if self.HeadingIsClear(int(hlt.SHIP_MAX_SPEED), heading.GetAngleInRads(), gameMap, self.Info.ClosestEnemy.Ship.Id) {
			heading.Magnitude = int(hlt.SHIP_MAX_SPEED)
		}
	} else if self.Ship.FireNextTurn {
		message = ALREADY_FIRED
		dir := self.Info.ClosestEnemy.Ship.AngleTo(&self.Ship.Point)
		targetPos := self.Ship.Point.AddThrust(hlt.SHIP_MAX_SPEED+1, dir)
		heading = self.MoveToPoint(&targetPos, gameMap)
	} else {
		log.Println("TOTAL enemies/allies", self.Info.TotalEnemies, self.Info.TotalAllies)
		message = MOVING_TOWARD_ENEMY
		enemyShipVel := self.Info.ClosestEnemy.Ship.Vel
		log.Println("TOTAL CHASING", turnComm.Chasing[self.Info.ClosestEnemy.Ship.Id])

		log.Println("SHOULD NOT CHASE?", turnComm.Chasing[self.Info.ClosestEnemy.Ship.Id], "?>=4 AND", self.Info.ClosestDockedEnemy.Distance, "?< 200")
		if turnComm.Chasing[self.Info.ClosestEnemy.Ship.Id] >= 4 && self.Info.ClosestDockedEnemy.Distance < 200 {
			log.Println("Too many chasing, going for docked ship")
			message = TOO_MANY_CHASING_MOVING_TOWARD_DOCKED_ENEMY
			heading = self.MoveToShip(self.Info.ClosestDockedEnemy.Ship, gameMap)
		} else if enemyShipVel.Magnitude() > 0 && self.Info.ClosestEnemy.Ship.IsAliveNextTurn() && self.Info.TotalAllies+2 <= self.Info.TotalEnemies {
			message = MOVING_TO_CLOSEST_ENEMY_MAX_RANGE
			newV := enemyShipVel.RescaleToMag(int(enemyShipVel.Magnitude()+.5) + int(hlt.SHIP_MAX_ATTACK_RANGE) + 1)
			targetP := self.Info.ClosestEnemy.Ship.AddVector(&newV)
			heading = self.MoveToPoint(&targetP, gameMap)
			turnComm.Chasing[self.Info.ClosestEnemy.Ship.Id] = turnComm.Chasing[self.Info.ClosestEnemy.Ship.Id] + 1
		} else {
			message = MOVING_TOWARD_ENEMY
			heading = self.MoveToShip(self.Info.ClosestEnemy.Ship, gameMap)
			turnComm.Chasing[self.Info.ClosestEnemy.Ship.Id] = turnComm.Chasing[self.Info.ClosestEnemy.Ship.Id] + 1
		}
	}

	return message, heading
}

func (self *ShipController) UpdateInfo(gameMap *hlt.GameMap) {
	self.Info = CreateShipTurnInfo(self.Ship, gameMap)
}

func nextCorner(current hlt.Point, gameMap *hlt.GameMap) hlt.Point {
	ne := gameMap.GetNECorner()
	nw := gameMap.GetNWCorner()
	se := gameMap.GetSECorner()
	sw := gameMap.GetSWCorner()
	if current.Equals(&ne) {
		return nw
	} else if current.Equals(&nw) {
		return sw
	} else if current.Equals(&sw) {
		return se
	} else if current.Equals(&se) {
		return ne
	}
	return ne
}

func (self *ShipController) SetTarget(gameMap *hlt.GameMap) {
	if self.Ship.DockingStatus != hlt.UNDOCKED {
		log.Println("We are docking, docked, or undocking")
		self.Target = &self.Info.PlanetsByDist[0].Point
		self.TargetPlanet = self.Info.PlanetsByDist[0].Id
		log.Println("Target planet is",self.TargetPlanet)
	} else if self.Mission == MISSION_NORMAL {
		self.NormalSetTarget(gameMap)
	} else if self.Mission == MISSION_SETTLER {
		self.SettlerSetTarget(gameMap)
	} else if self.Mission == MISSION_RUN_AWAY {
		self.RunAwaySetTarget(gameMap)
	} else if self.Mission == MISSION_SNEAKY {
		self.SneakySetTarget(gameMap)
	}
}

func (self *ShipController) Act(gameMap *hlt.GameMap, turnComm *TurnComm) string {

	log.Println("Ship ", self.Id, " Act. Planet is ", self.TargetPlanet)
	log.Println("ClosestEnemy is ", self.Info.ClosestEnemy.Distance)

	heading := hlt.Heading{
		Magnitude: 0,
		Angle:     0,
	}
	message := NONE

	log.Println("Looking for", self.Info.ClosestEnemy.Ship.Id, "in", turnComm.Chasing)
	if _, ok := turnComm.Chasing[self.Info.ClosestEnemy.Ship.Id]; !ok {
		log.Println("didn't find it, setting to zero")
		turnComm.Chasing[self.Info.ClosestEnemy.Ship.Id] = 0
	}

	if self.Mission == MISSION_RUN_AWAY {
		message, heading = self.RunAwayAct(gameMap, turnComm)
	} else if self.Mission == MISSION_SNEAKY {
		message, heading = self.SneakyAct(gameMap, turnComm)
	} else if self.Mission == MISSION_NORMAL {
		message, heading = self.NormalAct(gameMap, turnComm)
	} else if self.Mission == MISSION_SETTLER {
		message, heading = self.SettlerAct(gameMap, turnComm)
	}

	log.Println("Should we dock?", message, DOCK)
	if message == DOCK {
		log.Println("You should dock")
		planet := gameMap.PlanetLookup[self.TargetPlanet]
		return self.Ship.Dock(planet)
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
	log.Println("turn complete heading:", heading, "message:", message)
	return heading.ToMoveCmd(self.Ship, int(message))
}
