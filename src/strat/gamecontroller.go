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

	gameMap.UpdateShipsFromHistory(self.GameMap)
	gameMap.LookaheadCalculations()

	self.GameMap = gameMap

	self.Info = CreateGameTurnInfo(gameMap)

	for _, id := range gameMap.MyShips {
		ship := gameMap.ShipLookup[id]
		_, contains := self.ShipControllers[ship.Entity.Id]
		if !contains {
			sc := ShipController{
				Ship:         ship,
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
			sc.Update(ship)
		}
	}

	for key, sc := range self.ShipControllers {
		contains := false
		for _, id := range gameMap.MyShips {
			if sc.Id == id {
				contains = true
			}

		}
		if !contains {
			delete(self.ShipControllers, key)
		}
	}
}

func (self *GameController) AssignToPlanets() {
	var free []*hlt.Planet
	assignments := make(map[int]int)

	for _, id := range self.GameMap.Planets {
		assignments[id] = 0
	}

	for _, sc := range self.ShipControllers {
		if sc.TargetPlanet != -1 {
			assignments[sc.TargetPlanet] += 1
		}
	}

	for _, id := range self.GameMap.Planets {
		assigned := assignments[id]
		p := self.GameMap.PlanetLookup[id]
		if (p.Owned == 0 || p.Owner == self.GameMap.MyId) && assigned < p.NumDockingSpots {
			free = append(free, p)
		}
	}

	for _, sc := range self.ShipControllers {
		if sc.ShipNum == 5 || sc.ShipNum % 17 == 0 {
			if sc.TargetPlanet == -1 {
				sc.SetRushPlanet(self.GameMap)
			}
			sc.Mission = MISSION_RUSH_AND_DISTRACT
			continue
		}

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
			log.Println("Assigning", sc.Ship.Id, "to", closest)
			sc.TargetPlanet = closest
			if sc.ShipNum%15 == 0 && self.Info.ShipCountDeltaToLeader != 0 {
				sc.Mission = MISSION_FOUND_PLANET
			}
		}
	}
}

func (self *GameController) UpdateShipTurnInfos() {
	for _, sc := range self.ShipControllers {
		sc.UpdateInfo(self.GameMap)
		sc.SetTarget(self.GameMap)
	}
}

func (self *GameController) Act(turn int) []string {
	self.UpdateShipTurnInfos()

	if self.Info.ActivateStupidRunAwayMeta {
		log.Println("Activating StupidRunAwayMetaa")
		self.StupidRunAwayMeta()
	}

	if turn == 1 {
		self.GameStart()
	} else {
		self.AssignToPlanets()
	}
	return self.ExecuteShipTurn(turn)
}

func (self *GameController) StupidRunAwayMeta() {
	for _, sc := range self.ShipControllers {
		sc.Mission = STUPID_RUN_AWAY_META
	}
}

func (self *GameController) GameStart() {
	bestTargetDist := 1000000.0
	targetPlanet := -1
	for _, id := range self.GameMap.MyShips {
		ship := self.GameMap.ShipLookup[id]
		nearestPlanets := self.GameMap.NearestPlanetsByDistance(ship)
		for _, p := range nearestPlanets {
			if int(nearestPlanets[0].Distance/7.0) > int(p.Distance/7.0)+4 {
				continue
			}
			if targetPlanet == -1 && p.NumDockingSpots >= 3 {
				if p.Distance < bestTargetDist {
					bestTargetDist = p.Distance
					targetPlanet = p.Id
				}
			}
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
}

func (self *GameController) GetSCsInOrder() []*ShipController {
	scs := []*ShipController{}
	for _, sc := range self.ShipControllers {
		sc.Distance = sc.Ship.SqDistanceTo(sc.Target)
		scs = append(scs, sc)
	}

	sort.Sort(byDistSc(scs))

	return scs
}

func (self *GameController) ExecuteShipTurn(turn int) []string {
	commandQueue := []string{}

	turnComm := GetTurnComm()
	scs := self.GetSCsInOrder()

	log.Println("Chasing is", turnComm.Chasing)

	for _, sc := range scs {
		ship := sc.Ship
		log.Println("Ship", sc.Id, "turn", turn, "with ship num",sc.ShipNum)
		log.Println(sc.Id, "is assigned to planet ", sc.TargetPlanet)
		log.Println(ship)
		log.Println("Ship is located at ", ship.Point)
		log.Println("With Vel ", ship.Vel, " and mag ", ship.Vel.Magnitude())
		if ship.DockingStatus == hlt.UNDOCKED {
			cmd := sc.Act(self.GameMap, &turnComm)
			log.Println(cmd)
			commandQueue = append(commandQueue, cmd)
		} else if sc.Mission == STUPID_RUN_AWAY_META && ship.DockingStatus == hlt.DOCKED {
			cmd := ship.Undock()
			log.Println(cmd)
			commandQueue = append(commandQueue, cmd)
		}
	}

	log.Println("Chasing is now", turnComm.Chasing)

	return commandQueue
}
