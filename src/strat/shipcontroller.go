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
	MISSION_RUSH_AND_DISTRACT
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
		if !s.IsAliveNextTurn() {
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

	canKillSuicideOnProduction := self.Info.ClosestDockedEnemyShipDistance < hlt.SHIP_MAX_SPEED && self.Info.ClosestDockedEnemyShip.IsAliveNextTurn() && self.HeadingIsClear(int(self.Info.ClosestDockedEnemyShipDistance+.5), self.Info.ClosestDockedEnemyShipDir, gameMap, self.Info.ClosestDockedEnemyShip.Id)
	canKillSuicideOnNearestEnemy := self.Info.ClosestEnemyShip.IsAliveNextTurn() && self.Info.ClosestEnemyShipDistance < hlt.SHIP_MAX_SPEED && self.HeadingIsClear(int(self.Info.ClosestEnemyShipDistance+.5), self.Info.ClosestEnemyShipDir, gameMap, self.Info.ClosestEnemyShip.Id)

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
	} else if self.Ship.FireNextTurn {
		message = ALREADY_FIRED
		dir := self.Info.ClosestEnemyShip.AngleTo(&self.Ship.Point)
		targetPos := self.Ship.Point.AddThrust(hlt.SHIP_MAX_SPEED + 1, dir)
		heading = self.MoveToPoint(&targetPos, gameMap)
	} else {
		log.Println("TOTAL enemies/allies", self.Info.TotalEnemies, self.Info.TotalAllies)
		message = MOVING_TOWARD_ENEMY
		enemyShipVel := self.Info.ClosestEnemyShip.Vel
		log.Println("TOTAL CHASING",turnComm.Chasing[self.Info.ClosestEnemyShip.Id])
		
		log.Println("SHOULD NOT CHASE?", turnComm.Chasing[self.Info.ClosestEnemyShip.Id], "?>=4 AND", self.Info.ClosestDockedEnemyShipDistance, "?< 200")
		if turnComm.Chasing[self.Info.ClosestEnemyShip.Id] >= 4 && self.Info.ClosestDockedEnemyShipDistance < 200 {
			log.Println("Too many chasing, going for docked ship")
			message = TOO_MANY_CHASING_MOVING_TOWARD_DOCKED_ENEMY
			heading = self.MoveToShip(self.Info.ClosestDockedEnemyShip, gameMap)
		} else if enemyShipVel.Magnitude() > 0 && self.Info.ClosestEnemyShip.IsAliveNextTurn() && self.Info.TotalAllies + 2 <= self.Info.TotalEnemies {
			message = MOVING_TO_CLOSEST_ENEMY_MAX_RANGE 
			newV := enemyShipVel.RescaleToMag(int(enemyShipVel.Magnitude() + .5) + int(hlt.SHIP_MAX_ATTACK_RANGE) + 1)
			targetP := self.Info.ClosestEnemyShip.AddVector(&newV)
			heading = self.MoveToPoint(&targetP, gameMap)
			turnComm.Chasing[self.Info.ClosestEnemyShip.Id] = turnComm.Chasing[self.Info.ClosestEnemyShip.Id] + 1
		} else {
			message = MOVING_TOWARD_ENEMY
			heading = self.MoveToShip(self.Info.ClosestEnemyShip, gameMap)
			turnComm.Chasing[self.Info.ClosestEnemyShip.Id] = turnComm.Chasing[self.Info.ClosestEnemyShip.Id] + 1
		}
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

	dir := self.Info.ClosestEnemyShip.AngleTo(&self.Ship.Point)
	targetPos := self.Ship.Point.AddThrust(hlt.SHIP_MAX_SPEED, dir)

	heading = self.MoveToPoint(&targetPos, gameMap)

	return message, heading
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


func (self *ShipController) SetRushPlanet(gameMap *hlt.GameMap) {
	mins := []float64{10000,10000,10000,10000,10000}
	minID := []int{-1,-1,-1,-1,-1}
	maxs := []float64{0,0,0,0,0}
	maxID := []int{-1,-1,-1,-1,-1}
	log.Println("mins/maxs", mins, maxs)
	for _, pid := range(gameMap.Planets) {
		p := gameMap.PlanetLookup[pid]
		if p.Owner == gameMap.MyId {
			for _, tid := range(gameMap.Planets) {
				t := gameMap.PlanetLookup[tid]
				if t.Owner == gameMap.MyId || t.Owner == 0 {
					continue
				}
				d := p.SqDistanceTo(&t.Point)
				if d < mins[t.Owner] {
					mins[t.Owner] = d
					minID[t.Owner] = t.Id
				} 
				if d > maxs[t.Owner] {
					maxs[t.Owner] = d
					maxID[t.Owner] = t.Id
				}
			}

		}
	}
	log.Println("mins/maxs", mins, maxs)
	log.Println("minID/maxID",minID,maxID)
	minIdx := -1
	minVal := 10000.0
	for i := 0; i < 4; i++ {
		if mins[i] < minVal {
			minVal = mins[i]
			minIdx = i
		}
	}

	self.TargetPlanet = maxID[minIdx]
}


func (self *ShipController) rushAndDistract(gameMap *hlt.GameMap) (ChlMessage, hlt.Heading) {
	heading := hlt.Heading{
		Magnitude: 0,
		Angle:     0,
	}
	message := RUSH
	p := gameMap.PlanetLookup[self.TargetPlanet]
	dToTarget := self.Ship.DistanceToCollision(&p.Entity)
	vToTarget := self.Ship.VectorTo(&p.Point)
	vToTarget = vToTarget.RescaleToMag(int(hlt.SHIP_MAX_SPEED))
	log.Println("Ship", self.Ship.Id, "rushing to", self.TargetPlanet)

	if self.Info.TotalEnemies > 0 {
		message = RUSH_AVOIDING_ENEMY
		vAway := self.Info.ClosestNonDockedEnemyShip.VectorTo(&self.Ship.Point)
		vAway = vAway.RescaleToMag(int(hlt.SHIP_MAX_SPEED))
		toGo := vToTarget.Add(&vAway)
		if self.Info.AlliedClosestPlanetDist < 1000 && self.Info.AlliedClosestPlanetDist < 1.5 * self.Info.EnemyClosestPlanetDist {
			awayFromOurP := self.Info.AlliedClosestPlanet.VectorTo(&self.Ship.Point)
			awayFromOurP = awayFromOurP.RescaleToMag((int(hlt.SHIP_MAX_SPEED)))
			toGo = awayFromOurP.Add(&vAway)
		}
		toGo = toGo.Add(&vAway)
		toGo = toGo.RescaleToMag(int(hlt.SHIP_MAX_SPEED) + 1)
		targetPos := self.Ship.AddVector(&toGo)
		heading = self.MoveToPoint(&targetPos, gameMap)
	} else if dToTarget < 20 && self.Info.ClosestDockedEnemyShipDistance < 100 {
		message = RUSH_KILLING_DOCKED
		heading = self.MoveToShip(self.Info.ClosestDockedEnemyShip, gameMap)
	} else {
		toP := self.Ship.AngleTo(&p.Point)
		newAng := toP - math.Pi/4
		if self.Id % 2 == 0 {
			newAng = toP + math.Pi/4

		}
		if self.HeadingIsClear(int(hlt.SHIP_MAX_SPEED), newAng, gameMap, -1) {
			heading = hlt.Heading{
				Magnitude: int(hlt.SHIP_MAX_SPEED),
				Angle: int(hlt.RadToDeg(newAng)),
			}
		} else {
			heading = self.MoveToPlanet(p, gameMap)
		}
	}

	return message, heading
}


func (self *ShipController) stupidRunAwayMeta(gameMap *hlt.GameMap) (ChlMessage, hlt.Heading) {
	/*
		heading := hlt.Heading{
			Magnitude: 0,
			Angle:     0,
		}
		message := HIDE_WE_ARE_LOSING
	*/
	// TODO: head to corner

	return self.runAway(gameMap)
}

func (self *ShipController) IsTargetPlanetStillValid(gameMap *hlt.GameMap) (bool, ChlMessage) {
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

func (self *ShipController) SetTarget(gameMap *hlt.GameMap) {
	if self.Mission == MISSION_FOUND_PLANET {
		planet := gameMap.PlanetLookup[self.TargetPlanet]
		self.Target = &planet.Point
	} else if self.TargetPlanet == -1 {
		self.Target = &self.Info.ClosestEnemyShip.Point
	} else if self.TargetPlanet != -1 {
		valid, _ := self.IsTargetPlanetStillValid(gameMap)
		if valid {
			planet := gameMap.PlanetLookup[self.TargetPlanet]
			self.Target = &planet.Point
		} else {
			self.Target = &self.Info.ClosestEnemyShip.Point
		}
	}
}

func (self *ShipController) Act(gameMap *hlt.GameMap, turnComm *TurnComm) string {

	log.Println("Ship ", self.Id, " Act. Planet is ", self.TargetPlanet)
	log.Println("ClosestEnemy is ", self.Info.ClosestEnemyShipDistance)

	heading := hlt.Heading{
		Magnitude: 0,
		Angle:     0,
	}
	message := NONE

	log.Println("Looking for", self.Info.ClosestEnemyShip.Id, "in", turnComm.Chasing)
	if _, ok := turnComm.Chasing[self.Info.ClosestEnemyShip.Id]; !ok {
		log.Println("didn't find it, setting to zero")
		turnComm.Chasing[self.Info.ClosestEnemyShip.Id] = 0
	}

	if self.Mission == STUPID_RUN_AWAY_META {
		message, heading = self.stupidRunAwayMeta(gameMap)
	} else if self.Mission == MISSION_RUSH_AND_DISTRACT {
		message, heading = self.rushAndDistract(gameMap)	
	} else if self.Info.TotalEnemies > 0 {
		message, heading = self.combat(gameMap, turnComm)
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
		valid, pmess := self.IsTargetPlanetStillValid(gameMap)
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
	} else if self.Info.ClosestDockedEnemyShipDistance < 200 && turnComm.Chasing[self.Info.ClosestEnemyShip.Id] > 4 {
		message = TOO_MANY_CHASING_FAR_AWAY_MOVING_TO_DOCKED_ENEMY
		heading = self.MoveToShip(self.Info.ClosestDockedEnemyShip, gameMap)
	} else {
		message = MOVING_TOWARD_ENEMY
		turnComm.Chasing[self.Info.ClosestEnemyShip.Id] = turnComm.Chasing[self.Info.ClosestEnemyShip.Id] + 1
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
	log.Println("turn complete heading:", heading, "message:", message)
	return heading.ToMoveCmd(self.Ship, int(message))
}
