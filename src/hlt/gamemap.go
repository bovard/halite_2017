package hlt

import (
	"log"
	"sort"
	"strconv"
	"strings"
)

type GameMap struct {
	MyId, Width, Height int
	Turn                int
	Players             [4]Player
	Planets             []int
	EnemyShips          []int
	MyShips             []int
	ShipLookup          map[int]*Ship
	PlanetLookup        map[int]*Planet
}

type Player struct {
	Id    int
	Ships []int /// preallocating for speed, assuming we cant have > 10k ships.
}

func (self *GameMap) GetNECorner() Point {
	return Point{
		X: float64(self.Width),
		Y: 0,
	}
}

func (self *GameMap) GetSECorner() Point {
	return Point{
		X: float64(self.Width),
		Y: float64(self.Height),
	}
}

func (self *GameMap) GetSWCorner() Point {
	return Point{
		X: 0,
		Y: float64(self.Height),
	}
}

func (self *GameMap) GetNWCorner() Point {
	return Point{
		X: 0,
		Y: 0,
	}
}

func (self *GameMap) ParsePlayer(tokens []string) (Player, []string) {
	playerId, _ := strconv.Atoi(tokens[0])
	playerNumShips, _ := strconv.ParseFloat(tokens[1], 64)

	player := Player{
		Id:    playerId,
		Ships: []int{},
	}

	tokens = tokens[2:]
	for i := 0; float64(i) < playerNumShips; i++ {
		ship, tokensnew := ParseShip(playerId, tokens)
		tokens = tokensnew
		player.Ships = append(player.Ships, ship.Id)
		self.ShipLookup[ship.Id] = &ship
		if ship.Owner == self.MyId {
			self.MyShips = append(self.MyShips, ship.Id)
		} else {
			self.EnemyShips = append(self.EnemyShips, ship.Id)
		}
	}

	return player, tokens
}

func (self *GameMap) ParseGameString(gameString string) {
	tokens := strings.Split(gameString, " ")
	numPlayers, _ := strconv.Atoi(tokens[0])
	tokens = tokens[1:]

	for i := 0; i < numPlayers; i++ {
		player, tokensnew := self.ParsePlayer(tokens)
		tokens = tokensnew
		self.Players[player.Id] = player
	}

	numPlanets, _ := strconv.Atoi(tokens[0])
	tokens = tokens[1:]

	for i := 0; i < numPlanets; i++ {
		planet, tokensnew := ParsePlanet(tokens)
		tokens = tokensnew
		self.Planets = append(self.Planets, planet.Id)
		self.PlanetLookup[planet.Id] = &planet
	}
}

func (self *GameMap) UpdateShipsFromHistory(lastFrame *GameMap) {
	log.Println("coming gamemap from turn", self.Turn, "to old turn", lastFrame.Turn)
	for _, id := range append(self.MyShips, self.EnemyShips...) {
		ship := self.ShipLookup[id]
		log.Println("Updating Ship", id)
		log.Println("we are currently at", ship.Point)
		if oldShip, ok := lastFrame.ShipLookup[id]; ok {
			log.Println("FOUND OLD MATCH at loc", oldShip.Point)
			ship.Born = oldShip.Born
			ship.Vel = oldShip.Point.VectorTo(&ship.Point)
			ship.LastPos = oldShip.Point
		}
	}
}

func (self *GameMap) LookaheadCalculations() {

	allShips := []*Ship{}
	for _, s := range self.ShipLookup {
		allShips = append(allShips, s)
	}

	for _, s := range self.ShipLookup {
		for _, t := range allShips {
			t.Distance = s.SqDistanceTo(&t.Point)
		}
		sort.Sort(byDistShip(allShips))
		numEnemies := 0
		for _, t := range allShips {
			if t.Distance > SHIP_SQ_MAX_ATTACK_DISTANCE {
				break
			} else if t.Owner != s.Owner {
				numEnemies++
			}
		}
		if numEnemies > 0 {
			s.FireNextTurn = true
			dmg := SHIP_DAMAGE / float64(numEnemies)
			for _, t := range allShips {
				if t.Distance > SHIP_SQ_MAX_ATTACK_DISTANCE {
					break
				} else if t.Owner != s.Owner {
					toApply := self.ShipLookup[t.Id]
					toApply.IncomingDamage += dmg
				}
			}
		}
	}

}

func (self *GameMap) NearestPlanetsByDistance(ship *Ship) []*Planet {
	planets := []*Planet{}

	for _, id := range self.Planets {
		planet := self.PlanetLookup[id]
		planet.Distance = ship.Entity.DistanceToCollision(&planet.Entity)
		planets = append(planets, planet)
	}

	sort.Sort(byDist(planets))

	return planets
}

func (self *GameMap) IsOnMap(p *Point) bool {
	if p.X <= .5 || p.Y <= .5 {
		return false
	} else if p.X >= float64(self.Width)-.5 || p.Y >= float64(self.Height)-.5 {
		return false
	}
	return true
}

func (self *GameMap) NearestShipsByDistance(ship *Ship, ships []int) []*Ship {
	var sortedShips []*Ship
	for _, id := range ships {
		s := self.ShipLookup[id]
		sortedShips = append(sortedShips, s)
	}

	for i := 0; i < len(sortedShips); i++ {
		sortedShips[i].Distance = ship.DistanceToCollision(&sortedShips[i].Entity)
	}

	sort.Sort(byDistShip(sortedShips))

	return sortedShips
}

type byDist []*Planet

func (a byDist) Len() int           { return len(a) }
func (a byDist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDist) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

type byDistShip []*Ship

func (a byDistShip) Len() int           { return len(a) }
func (a byDistShip) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDistShip) Less(i, j int) bool { return a[i].Distance < a[j].Distance }
