package hlt

import (
	"sort"
	"strconv"
	"strings"
)

type GameMap struct {
	MyId, Width, Height int
	Planets             []Planet /// preallocating for speed, assuming we cant have > 100 planets
	Players             [4]Player
	Entities            []Entity
	EnemyShips          []*Ship
	MyShips             []*Ship
	ShipLookup          map[int]*Ship
	PlanetsLookup       map[int]Planet
}

type Player struct {
	Id    int
	Ships []Ship /// preallocating for speed, assuming we cant have > 10k ships.
}

func ParsePlayer(tokens []string) (Player, []string) {
	playerId, _ := strconv.Atoi(tokens[0])
	playerNumShips, _ := strconv.ParseFloat(tokens[1], 64)

	player := Player{
		Id:    playerId,
		Ships: []Ship{},
	}

	tokens = tokens[2:]
	for i := 0; float64(i) < playerNumShips; i++ {
		ship, tokensnew := ParseShip(playerId, tokens)
		tokens = tokensnew
		player.Ships = append(player.Ships, ship)
	}

	return player, tokens
}

func ParseGameString(gameString string, self GameMap) GameMap {
	tokens := strings.Split(gameString, " ")
	numPlayers, _ := strconv.Atoi(tokens[0])
	tokens = tokens[1:]

	for i := 0; i < numPlayers; i++ {
		player, tokensnew := ParsePlayer(tokens)
		tokens = tokensnew
		self.Players[player.Id] = player
		for j := 0; j < len(player.Ships); j++ {
			self.ShipLookup[player.Ships[j].Id] = &player.Ships[j]
			self.Entities = append(self.Entities, player.Ships[j].Entity)
			if i == self.MyId {
				self.MyShips = append(self.MyShips, &player.Ships[j])
			} else {
				self.EnemyShips = append(self.EnemyShips, &player.Ships[j])
			}
		}
	}

	numPlanets, _ := strconv.Atoi(tokens[0])
	tokens = tokens[1:]

	for i := 0; i < numPlanets; i++ {
		planet, tokensnew := ParsePlanet(tokens)
		tokens = tokensnew
		self.Planets = append(self.Planets, planet)
		self.PlanetsLookup[planet.Entity.Id] = planet
		self.Entities = append(self.Entities, planet.Entity)
	}

	return self
}

func (gameMap *GameMap) UpdateShipsFromHistory(lastFrame *GameMap) {
	for _, ship := range append(gameMap.MyShips, gameMap.EnemyShips...) {
		if oldShip, ok := lastFrame.ShipLookup[ship.Id]; ok {
			ship.Born = oldShip.Born
			ship.Vel = oldShip.Point.VectorTo(&ship.Point)
		} else {
			ship.Born = ship.Point
			ship.Vel = Vector{ 
				X: 0,
				Y: 0,
			}
		}
	}

}


func (gameMap *GameMap) NearestPlanetsByDistance(ship *Ship) []Planet {
	planets := gameMap.Planets

	for i := 0; i < len(planets); i++ {
		planets[i].Distance = ship.Entity.DistanceToCollision(&planets[i].Entity)
	}

	sort.Sort(byDist(planets))

	return planets
}

func (self *GameMap) IsOnMap(p *Point) bool {
	if p.X <= .5 || p.Y <= .5 {
		return false
	} else if p.X >= float64(self.Width) - .5 || p.Y >= float64(self.Height) - .5 {
		return false
	}
	return true
}

func (gameMap GameMap) NearestShipsByDistance(ship *Ship, ships []*Ship) []Ship {
	var enemies []Ship
	for _, e := range ships {
		enemies = append(enemies, *e)
	}

	for i := 0; i < len(enemies); i++ {
		enemies[i].Distance = ship.Entity.DistanceToCollision(&enemies[i].Entity)
	}

	sort.Sort(byDistShip(enemies))

	return enemies
}

func (gameMap GameMap) NearestEnemiesByDistance(ship Ship) []Entity {
	entities := gameMap.Entities
	var enemies []Entity
	for _, e := range entities {
		if e.Owner != gameMap.MyId && e.Owner != -1 && e.Radius < 1 {
			enemies = append(enemies, e)
		}
	}

	for i := 0; i < len(enemies); i++ {
		enemies[i].Distance = ship.Entity.DistanceToCollision(&enemies[i])
	}

	sort.Sort(byDistEntity(enemies))

	return enemies
}

type byDistEntity []Entity

func (a byDistEntity) Len() int           { return len(a) }
func (a byDistEntity) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDistEntity) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

type byDistShip []Ship

func (a byDistShip) Len() int           { return len(a) }
func (a byDistShip) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDistShip) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

type byDist []Planet

func (a byDist) Len() int           { return len(a) }
func (a byDist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDist) Less(i, j int) bool { return a[i].Distance < a[j].Distance }
