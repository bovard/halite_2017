package hlt

import (
	"math"
	"strconv"
	"fmt"

)

type DockingStatus int

const (
	UNDOCKED  DockingStatus = iota
	DOCKING
	DOCKED
	UNDOCKING
)

type Entity struct {
	X      float64
	Y      float64
	Radius float64
	Health float64
	Owner  int
	Id        int
	Distance float64
}

type Position struct {
	X, Y float64
}

type Planet struct {
	Entity
	NumDockingSpots    int 
	NumDockedShips     int
	CurrentProduction  float64
	RemainingResources float64
	DockedShipIds      [] int
	DockedShips        [] Ship
	Owned              float64
	Distance           float64
}

type Ship struct {
	Entity
	VelX float64
	VelY float64

	PlanetId        int
	Planet          Planet
	DockingStatus   DockingStatus
	DockingProgress float64
	WeaponCooldown  float64
}

func (self Entity) RotateAround(target Entity, angle float64) Entity {
	x1 := self.X - target.X
	y1 := self.Y - target.Y
    x2 := x1 * math.Cos(angle) - y1 * math.Sin(angle)
	y2 := x1 * math.Sin(angle) - y1 * math.Cos(angle)
	return Entity{
		X: x2 + target.X,
		Y: y2 + target.Y,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
}

func (self Entity) AddThrust(magnitude float64, angle float64) Entity {
	return Entity {
		X: self.X + magnitude * math.Cos(angle),
		Y: self.Y + magnitude * math.Sin(angle),
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
}

func (self Entity) GetMidPoint(target Entity) Entity {
	return Entity {
		X: (self.X + target.X)/2,
		Y: (self.Y + target.Y)/2,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
}

func (self Entity) CalculateDistanceTo(target Entity) float64 {
	// returns euclidean distance to target
	dx := target.X - self.X
	dy := target.Y - self.Y

	return math.Sqrt(dx*dx + dy*dy)
}

func (self Entity) CalculateAngleTo(target Entity) float64 {
	// returns angle in radians from self to target
	dx := target.X - self.X
	dy := target.Y - self.Y

	return math.Atan2(dy, dx)
}

func (self Entity) ClosestPointTo(target Entity, minDistance float64) Entity {
	// returns closest point to self that is at least minDistance from target
	dist := self.CalculateDistanceTo(target) - target.Radius - minDistance
	angle := target.CalculateAngleTo(self)
	return target.AddThrust(dist, angle)
}

func ParseShip(playerId int, tokens []string) (Ship, [] string) {

	shipId, _ := strconv.Atoi(tokens[0])
	shipX, _ := strconv.ParseFloat(tokens[1], 64)
	shipY, _ := strconv.ParseFloat(tokens[2], 64)
	shipHealth, _ := strconv.ParseFloat(tokens[3], 64)
	shipVelX, _ := strconv.ParseFloat(tokens[4], 64)
	shipVelY, _ := strconv.ParseFloat(tokens[5], 64)
	shipDockingStatus, _ := strconv.Atoi(tokens[6])
	shipPlanetId, _ := strconv.Atoi(tokens[7])
	shipDockingProgress, _ := strconv.ParseFloat(tokens[8], 64)
	shipWeaponCooldown, _ := strconv.ParseFloat(tokens[9], 64)

	shipEntity := Entity{
		X:      shipX,
		Y:      shipY,
		Radius: .5,
		Health: shipHealth,
		Owner:  playerId,
		Id:     shipId,
	}

	ship := Ship{
		PlanetId:        shipPlanetId,
		DockingStatus:   IntToDockingStatus(shipDockingStatus),
		DockingProgress: shipDockingProgress,
		WeaponCooldown:  shipWeaponCooldown,
		VelX:            shipVelX,
		VelY:            shipVelY,
		Entity:          shipEntity,
	}

	return ship, tokens[10:]
}

func ParsePlanet(tokens []string) (Planet, [] string) {

	planetId, _ := strconv.Atoi(tokens[0])
	planetX, _ := strconv.ParseFloat(tokens[1], 64)
	planetY, _ := strconv.ParseFloat(tokens[2], 64)
	planetHealth, _ := strconv.ParseFloat(tokens[3], 64)
	planetRadius, _ := strconv.ParseFloat(tokens[4], 64)
	planetNumDockingSpots, _ := strconv.ParseInt(tokens[5], 10, 32)
	planetCurrentProduction, _ := strconv.ParseFloat(tokens[6], 64)
	planetRemainingResources, _ := strconv.ParseFloat(tokens[7], 64)
	planetOwned, _ := strconv.ParseFloat(tokens[8], 64)
	planetOwner, _ := strconv.Atoi(tokens[9])
	planetNumDockedShips, _ := strconv.ParseInt(tokens[10], 10, 32)

	planetEntity := Entity{
		X:      planetX,
		Y:      planetY,
		Radius: planetRadius,
		Health: planetHealth,
		Owner:  planetOwner,
		Id:     planetId,
	}

	planet := Planet{
		NumDockingSpots:    int(planetNumDockingSpots),
		NumDockedShips:     int(planetNumDockedShips),
		CurrentProduction:  planetCurrentProduction,
		RemainingResources: planetRemainingResources,
		DockedShipIds:      nil,
		DockedShips:        nil,
		Owned:              planetOwned,
		Entity:             planetEntity,
	}

	for i := 0; i < int(planetNumDockedShips); i++ {
		dockedShipId, _ := strconv.Atoi(tokens[11+i])
		planet.DockedShipIds = append(planet.DockedShipIds, dockedShipId)
	}
	return planet, tokens[11+int(planetNumDockedShips):]
}

func IntToDockingStatus(i int) DockingStatus {
	statuses := [4]DockingStatus{UNDOCKED, DOCKING, DOCKED, UNDOCKING}
	return statuses[i]
}

func (ship Ship) Thrust(magnitude float64, angle float64) string {
	angle = RadToDeg(angle)
	if angle < 0 {
		angle += 360
	} else if angle > 359 {
		angle -= 360
	}
	return fmt.Sprintf("t %s %s %s", strconv.Itoa(ship.Id), strconv.Itoa(int(magnitude)), strconv.Itoa(int(angle)))
}

func (ship Ship) Dock(planet Planet) string {
	return fmt.Sprintf("d %s %s", strconv.Itoa(ship.Id), strconv.Itoa(planet.Id))
}

func (ship Ship) Undock() string {
	return fmt.Sprintf("u %s %s", strconv.Itoa(ship.Id))
}

func (ship Ship) NavigateBasic(target Entity) string {

	distance := ship.CalculateDistanceTo(target)

	angle := ship.CalculateAngleTo(target)
	speed := SHIP_MAX_SPEED
	if distance < SHIP_MAX_SPEED {
		speed = distance - (SHIP_RADIUS + .1)
	} 
	speed = math.Min(speed, distance)
	return ship.Thrust(speed, angle)
}

func (ship Ship) CanDock(planet Planet) bool {
	dist := ship.CalculateDistanceTo(planet.Entity)

	return dist <= (planet.Radius + SHIP_DOCKING_RADIUS + ship.Radius)
}

func (ship Ship) Navigate(target Entity, gameMap Map) string {


	ob := gameMap.ObstaclesBetween(ship.Entity, target)

	if !ob {
		return ship.NavigateBasic(target)
	} else {

		x0 := math.Min(ship.X, target.X)
		x2 := math.Max(ship.X, target.X)
		y0 := math.Min(ship.Y, target.Y)
		y2 := math.Max(ship.Y, target.Y)

		dx := (x2 - x0) / 5
		dy := (y2 - y0) / 5
		bestdist := 1000.0
		bestTarget := target

		for x1 := x0; x1 <= x2; x1 += dx {
			for y1 := y0; y1 <= y2; y1 += dy {
				intermediateTarget := Entity{
					X:      x1,
					Y:      y1,
					Radius: 0,
					Health: 0,
					Owner:  0,
					Id:     -1,
				}
				ob1 := gameMap.ObstaclesBetween(ship.Entity, intermediateTarget)
				if !ob1 {
					ob2 := gameMap.ObstaclesBetween(intermediateTarget, target)
					if !ob2 {
						totdist := math.Sqrt(math.Pow(x1-x0, 2)+math.Pow(y1-y0, 2)) + math.Sqrt(math.Pow(x1-x2, 2)+math.Pow(y1-y2, 2))
						if totdist < bestdist {
							bestdist = totdist
							bestTarget = intermediateTarget


						}
					}
				}
			}
		}


		return ship.NavigateBasic(bestTarget)

	}

}


func (ship Ship) BetterNavigate(target Entity, gameMap Map) string {


	ob := gameMap.ObstaclesBetween(ship.Entity, target)

	if !ob {
		return ship.NavigateBasic(target)
	} else {
		maxTurn := (3 * math.Pi) / 2
		dTurn := math.Pi / 8

		startSpeed := math.Min(SHIP_MAX_SPEED, ship.Entity.CalculateDistanceTo(target) - target.Radius - ship.Entity.Radius - .05)

		for speed := startSpeed; speed > .25; speed /= 2 {
			for turn := dTurn; turn <= maxTurn; turn += dTurn {
				intermediateTargetLeft := ship.Entity.AddThrust(speed, turn)
				obLeft := gameMap.ObstaclesInPath(ship.Entity, speed, turn)
				intermediateTargetRight := ship.Entity.AddThrust(speed, -turn)
				obRight := gameMap.ObstaclesInPath(ship.Entity, speed, -turn)
				if !obLeft && !obRight {
					if intermediateTargetLeft.CalculateDistanceTo(target) < intermediateTargetRight.CalculateDistanceTo(target) {
						return ship.NavigateBasic(intermediateTargetLeft)
					} else {
						return ship.NavigateBasic(intermediateTargetRight)
					}
				} else if !obLeft {
					return ship.NavigateBasic(intermediateTargetLeft)
				} else if !obRight {
					return ship.NavigateBasic(intermediateTargetRight)
				}
			}
		}
		return ship.NavigateBasic(target)
	}
}
