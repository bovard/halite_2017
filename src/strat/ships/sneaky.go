package ships

import (
	"../../hlt"
	"log"
	"math"
)

func (self *ShipController) SneakySetTarget(gameMap *hlt.GameMap) {
	if self.Info.EnemyClosestPlanetDist < 20 && self.Info.ClosestAlly.Distance < 20 {
		log.Println("Allies nearby, switch to normal mode")
		self.TargetPlanet = -1
		self.Mission = MISSION_NORMAL
		self.NormalSetTarget(gameMap)
	} else {

	}
}

func (self *ShipController) SneakyAct(gameMap *hlt.GameMap, turnComm *TurnComm) (ChlMessage, hlt.Heading) {
	return self.rushAndDistract(gameMap)
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
		vAway := self.Info.ClosestNonDockedEnemy.Ship.VectorTo(&self.Ship.Point)
		vAway = vAway.RescaleToMag(int(hlt.SHIP_MAX_SPEED))
		toGo := vToTarget.Add(&vAway)
		if self.Info.AlliedClosestPlanetDist < 1000 && self.Info.AlliedClosestPlanetDist < 1.5*self.Info.EnemyClosestPlanetDist {
			awayFromOurP := self.Info.AlliedClosestPlanet.VectorTo(&self.Ship.Point)
			awayFromOurP = awayFromOurP.RescaleToMag((int(hlt.SHIP_MAX_SPEED)))
			toGo = awayFromOurP.Add(&vAway)
		}
		toGo = toGo.Add(&vAway)
		toGo = toGo.RescaleToMag(int(hlt.SHIP_MAX_SPEED) + 1)
		targetPos := self.Ship.AddVector(&toGo)
		heading = self.MoveToPoint(&targetPos, gameMap)
	} else if dToTarget < 20 && self.Info.ClosestDockedEnemy.Distance < 100 {
		message = RUSH_KILLING_DOCKED
		heading = self.MoveToShip(self.Info.ClosestDockedEnemy.Ship, gameMap)
	} else {
		toP := self.Ship.AngleTo(&p.Point)
		newAng := toP - math.Pi/4
		if self.Id%2 == 0 {
			newAng = toP + math.Pi/4

		}
		if self.HeadingIsClear(int(hlt.SHIP_MAX_SPEED), newAng, gameMap, -1) {
			heading = hlt.Heading{
				Magnitude: int(hlt.SHIP_MAX_SPEED),
				Angle:     int(hlt.RadToDeg(newAng)),
			}
		} else {
			heading = self.MoveToPlanet(p, gameMap)
		}
	}

	return message, heading
}

func (self *ShipController) SetRushPlanet(gameMap *hlt.GameMap) {
	mins := []float64{10000, 10000, 10000, 10000}
	minID := []int{-1, -1, -1, -1}
	maxs := []float64{0, 0, 0, 0}
	maxID := []int{-1, -1, -1, -1}
	log.Println("mins/maxs", mins, maxs)
	for _, pid := range gameMap.Planets {
		log.Println("Looking for playet", pid)
		p := gameMap.PlanetLookup[pid]
		log.Println("Looking at planet", p.Id)
		if p.Owned == 1 && p.Owner == gameMap.MyId {
			log.Println("We own planet", p.Id)
			for _, tid := range gameMap.Planets {
				log.Println("  Looking for playet", tid)
				t := gameMap.PlanetLookup[tid]
				log.Println("  Looking at planet", t.Id)
				if (t.Owner == gameMap.MyId && t.Owned == 1) || t.Owned == 0 {
					continue
				}
				log.Println("  Enemy Owns Planet", t.Id)
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
	log.Println("minID/maxID", minID, maxID)
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
