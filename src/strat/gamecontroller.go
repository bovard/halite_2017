package strat

import ( "../hlt")

type GameController struct {
	GameMap                 hlt.Map
	ShipControllers         map[int]ShipController
	ShipToPlanetAssignments map[int][]int
}

func (self GameController) UpdatePlanets(planets []hlt.Planet) {
	for _, p := range planets {
		self.ShipToPlanetAssignments[p.Entity.Id] = make([]int, 0)
	}
}

func (self GameController) Update(gameMap hlt.Map) {
	self.GameMap = gameMap
	myPlayer := gameMap.Players[gameMap.MyId]
	myShips  := myPlayer.Ships

	for key, _ := range self.ShipControllers {
		sc := self.ShipControllers[key]
		sc.Alive = false
	}

	for i := 0; i < len(myShips); i++ {
		ship := myShips[i]
		_, contains := self.ShipControllers[ship.Entity.Id]
		if !contains {
			self.ShipControllers[ship.Entity.Id] = ShipController {
				Ship:   ship,
				Past:   nil,
				Id:     ship.Entity.Id,
				Planet: -1,
				Alive:  true,
			}
		} else {
			sc := self.ShipControllers[ship.Entity.Id]
			sc.Update(ship)
		}
	}

	for key, _ := range self.ShipControllers {
		sc := self.ShipControllers[key]
		if !sc.Alive {
			if sc.Planet != -1 {
				assigned := self.ShipToPlanetAssignments[sc.Planet]
				self.ShipToPlanetAssignments[sc.Planet] = remove(assigned, sc.Id)
			}
			delete(self.ShipControllers, key)
		}
	}
}

func remove(s []int, i int) []int {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func (self GameController) AssignToPlanets() {
	var free [] hlt.Planet
	for _, p := range self.GameMap.Planets {
		assigned := len(self.ShipToPlanetAssignments[p.Entity.Id])
		if (p.Owned == 0 || p.Owner == self.GameMap.MyId) && assigned < p.NumDockingSpots {
			free = append(free, p)
		}
	}

	for key, _ := range self.ShipControllers {
		sc := self.ShipControllers[key]
		if sc.Planet == -1 {
			closest := -1
			closestDist := 10000.0
			for _, p := range free {
				dist := sc.Ship.Entity.CalculateDistanceTo(p.Entity)
				assigned := len(self.ShipToPlanetAssignments[p.Entity.Id])
				if dist < closestDist && assigned < p.NumDockingSpots {
					closestDist = dist
					closest = p.Entity.Id
				}
			}
			if closest != -1 {
				assigned := self.ShipToPlanetAssignments[closest]
				assigned = append(assigned, sc.Id)
				self.ShipToPlanetAssignments[closest] = assigned
				sc.Planet = closest
			}

		}
	}
}