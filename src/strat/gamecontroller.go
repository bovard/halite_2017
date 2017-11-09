package strat

import (
	"../hlt"
	"log"
	"sort"
)

type GameController struct {
	GameMap         *hlt.GameMap
	ShipControllers map[int]*ShipController
	ShipNumIdx      int
	Info            GameTurnInfo
}

func (self *GameController) Update(gameMap *hlt.GameMap) {
	self.GameMap = gameMap
	myPlayer := gameMap.Players[gameMap.MyId]
	myShips := myPlayer.Ships

	self.Info = CreateGameTurnInfo(gameMap)

	for i := 0; i < len(myShips); i++ {
		ship := myShips[i]
		_, contains := self.ShipControllers[ship.Entity.Id]
		if !contains {
			sc := ShipController{
				Ship:         &ship,
				Past:         nil,
				Id:           ship.Entity.Id,
				TargetPlanet: -1,
				Mission:      MISSION_NORMAL,
				ShipNum:      self.ShipNumIdx,
			}
			self.ShipNumIdx++
			self.ShipControllers[ship.Entity.Id] = &sc
		} else {
			sc := self.ShipControllers[ship.Entity.Id]
			sc.Update(&ship)
		}
	}

	for key, sc := range self.ShipControllers {
		contains := false
		for i := 0; i < len(myShips); i++ {
			if sc.Id == myShips[i].Entity.Id {
				contains = true
			}

		}
		if !contains {
			delete(self.ShipControllers, key)
		}
	}
}

func remove(s []int, i int) []int {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (self *GameController) AssignToPlanets() {
	var free []hlt.Planet
	assignments := make(map[int]int)

	for _, p := range self.GameMap.Planets {
		assignments[p.Entity.Id] = 0
	}

	for _, sc := range self.ShipControllers {
		if sc.TargetPlanet != -1 {
			assignments[sc.TargetPlanet] += 1
		}
	}

	for _, p := range self.GameMap.Planets {
		assigned := assignments[p.Entity.Id]
		if (p.Owned == 0 || p.Owner == self.GameMap.MyId) && assigned < p.NumDockingSpots {
			free = append(free, p)
		}
	}

	log.Println("Printing planet assignments")
	for key, sc := range self.ShipControllers {
		log.Println(key, " is assigned to ", sc.TargetPlanet, " with status ", sc.Ship.DockingStatus)
	}
	log.Println("End docking assignments")

	for _, sc := range self.ShipControllers {
		log.Println("Looking to make assignment for ship ", sc.Id)
		if sc.TargetPlanet != -1 {
			log.Println("already assigned to ", sc.TargetPlanet)
			continue
		}
		closest := -1
		closestDist := 10000.0
		for _, p := range free {
			dist := sc.Ship.DistanceToCollision(&p.Entity)
			assigned := assignments[p.Entity.Id]
			log.Println("Planet ", p.Id, " is ", dist, " away and has ", assigned, " of ", p.NumDockingSpots, " used")
			if dist < closestDist && assigned < p.NumDockingSpots {
				closestDist = dist
				closest = p.Id
			}
		}
		if closest != -1 {
			assignments[closest] += 1
			sc.TargetPlanet = closest
			if sc.ShipNum%15 == 0 && self.Info.ShipCountDeltaToLeader > 2 {
				sc.Mission = MISSION_FOUND_PLANET
			}
		}
	}
	log.Println("Reprinting planet assignments")
	for key, sc := range self.ShipControllers {
		log.Println(key, " is assigned to ", sc.TargetPlanet, " with status ", sc.Ship.DockingStatus)
	}
	log.Println("End reprinting docking assignments")
}

func (self *GameController) UpdateShipInfos() {
	for _, sc := range self.ShipControllers {
		sc.UpdateInfo(self.GameMap)
	}
}

func (self *GameController) Act(turn int) []string {
	self.UpdateShipInfos()
	if turn == 1 {
		return self.GameStart()
	} else {
		self.AssignToPlanets()
		return self.NormalTurn()
	}
}

func (self *GameController) GameStart() []string {
	centerShip := self.GameMap.Players[self.GameMap.MyId].Ships[1]
	nearestPlanets := self.GameMap.NearestPlanetsByDistance(&centerShip)
	nearestPlanetDist := nearestPlanets[0].Distance
	targetPlanet := -1
	for _, p := range nearestPlanets {
		if int(nearestPlanetDist/7.0) > int(p.Distance/7.0)+4 {
			continue
		}
		if targetPlanet == -1 && p.NumDockingSpots >= 3 {
			targetPlanet = p.Id
		}
	}

	if targetPlanet != -1 {
		for _, sc := range self.ShipControllers {
			sc.Mission = MISSION_FOUND_PLANET
			sc.TargetPlanet = targetPlanet
		}
	} else {
		self.AssignToPlanets()
	}

	return self.NormalTurn()
}

func (self *GameController) GetSCsInOrder() []*ShipController {
	scs := []*ShipController{}
	for _, sc := range self.ShipControllers {
		if sc.TargetPlanet != -1 {
			p := self.GameMap.PlanetsLookup[sc.TargetPlanet]
			sc.Distance = sc.Ship.DistanceToCollision(&p.Entity)
		} else {
			sc.Distance = sc.Info.ClosestEnemyShipDistance
		}
		scs = append(scs, sc)
	}

	sort.Sort(byDistSc(scs))

	return scs
}

func (self *GameController) NormalTurn() []string {
	commandQueue := []string{}

	scs := self.GetSCsInOrder()

	for _, sc := range scs {
		ship := sc.Ship
		log.Println(sc.Id, "is assigned to planet ", sc.TargetPlanet)
		log.Println("Ship is located at ", ship.Point)
		log.Println("With Vel ", ship.Vel, " and mag ", ship.Vel.Magnitude())
		if sc.TargetPlanet != -1 {
			targetPlanet := self.GameMap.PlanetsLookup[sc.TargetPlanet]
			log.Println("planet location is ", targetPlanet.Point, ", d = ", ship.DistanceToCollision(&targetPlanet.Entity))
			rad := ship.Point.AngleTo(&targetPlanet.Point)
			log.Println("angle to planet is ", int(360+hlt.RadToDeg(rad))%360)
		}
		if ship.DockingStatus == hlt.UNDOCKED {
			cmd := sc.Act(self.GameMap)
			log.Println(cmd)
			commandQueue = append(commandQueue, cmd)
		}
	}
	return commandQueue
}
