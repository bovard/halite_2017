package strat

import ( "../hlt")

type GameController struct {
	GameMap                 *hlt.Map
	ShipControllers         map[int]*ShipController
}

func (self *GameController) Update(gameMap *hlt.Map) {
	self.GameMap = gameMap
	myPlayer := gameMap.Players[gameMap.MyId]
	myShips  := myPlayer.Ships

	for i := 0; i < len(myShips); i++ {
		ship := myShips[i]
		_, contains := self.ShipControllers[ship.Entity.Id]
		if !contains {
			sc := ShipController {
				Ship:   &ship,
				Past:   nil,
				Id:     ship.Entity.Id,
				Planet: -1,
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
	var free [] hlt.Planet
	assignments := make(map[int]int)

	for _, p := range self.GameMap.Planets {
    	assignments[p.Entity.Id] = 0
    }

    for _, sc := range self.ShipControllers {
    	if sc.Planet != -1 {
    		assignments[sc.Planet] += 1
    	}
    }

	for _, p := range self.GameMap.Planets {
		assigned := assignments[p.Entity.Id]
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
				assigned := assignments[p.Entity.Id]
				if dist < closestDist && assigned < p.NumDockingSpots {
					closestDist = dist
					closest = p.Entity.Id
				}
			}
			if closest != -1 {
    			assignments[closest] += 1
				sc.Planet = closest
			}

		}
	}
}