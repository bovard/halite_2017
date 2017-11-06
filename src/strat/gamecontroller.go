package strat

import (
	"../hlt"
	"log"
)

type GameController struct {
	GameMap         *hlt.GameMap
	ShipControllers map[int]*ShipController
}

func (self *GameController) Update(gameMap *hlt.GameMap) {
	self.GameMap = gameMap
	myPlayer := gameMap.Players[gameMap.MyId]
	myShips := myPlayer.Ships

	for i := 0; i < len(myShips); i++ {
		ship := myShips[i]
		_, contains := self.ShipControllers[ship.Entity.Id]
		if !contains {
			sc := ShipController{
				Ship:   &ship,
				Past:   nil,
				Id:     ship.Entity.Id,
				TargetPlanet: -1,
			}
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
		}
	}
	log.Println("REprinting planet assignments")
	for key, sc := range self.ShipControllers {
		log.Println(key, " is assigned to ", sc.TargetPlanet, " with status ", sc.Ship.DockingStatus)
	}
	log.Println("End reprinting docking assignments")
}
