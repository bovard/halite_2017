package hlt

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

type DockingStatus int

const (
	UNDOCKED DockingStatus = iota
	DOCKING
	DOCKED
	UNDOCKING
)

type Ship struct {
	Entity
	Born            Point
	Vel             Vector
	NextVel         Vector
	PlanetId        int
	Planet          Planet
	DockingStatus   DockingStatus
	DockingProgress float64
	WeaponCooldown  float64
}

func ParseShip(playerId int, tokens []string) (Ship, []string) {

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

	shipPoint := Point{
		X: shipX,
		Y: shipY,
	}

	shipEntity := Entity{
		Point:  shipPoint,
		Radius: .5,
		Health: shipHealth,
		Owner:  playerId,
		Id:     shipId,
	}

	shipVel := Vector{
		X: shipVelX,
		Y: shipVelY,
	}

	nextVel := Vector{
		X: 0,
		Y: 0,
	}

	ship := Ship{
		PlanetId:        shipPlanetId,
		DockingStatus:   IntToDockingStatus(shipDockingStatus),
		DockingProgress: shipDockingProgress,
		WeaponCooldown:  shipWeaponCooldown,
		Vel:             shipVel,
		NextVel:         nextVel,
		Entity:          shipEntity,
	}

	return ship, tokens[10:]
}

func IntToDockingStatus(i int) DockingStatus {
	statuses := [4]DockingStatus{UNDOCKED, DOCKING, DOCKED, UNDOCKING}
	return statuses[i]
}

func (ship *Ship) Thrust(magnitude float64, angle float64) string {
	log.Println("Thurst with mag ", magnitude, " and angle ", angle)
	angle = RadToDeg(angle)
	if angle < 0 {
		angle += 360
	} else if angle > 359 {
		angle -= 360
	}
	return fmt.Sprintf("t %s %s %s", strconv.Itoa(ship.Id), strconv.Itoa(int(magnitude)), strconv.Itoa(int(angle)))
}

func (ship *Ship) Dock(planet *Planet) string {
	return fmt.Sprintf("d %s %s", strconv.Itoa(ship.Id), strconv.Itoa(planet.Id))
}

func (ship *Ship) Undock() string {
	return fmt.Sprintf("u %s %s", strconv.Itoa(ship.Id))
}

func (ship *Ship) NavigateBasic(target *Entity) string {

	maxMove := ship.Point.DistanceTo(&target.Point) - (ship.Entity.Radius + target.Radius + .1)

	angle := ship.Point.AngleTo(&target.Point)
	speed := math.Min(maxMove, SHIP_MAX_SPEED)
	return ship.Thrust(speed, angle)
}

func (ship *Ship) CanDock(planet *Planet) bool {
	dist := ship.Point.DistanceTo(&planet.Point)

	return dist <= (planet.Radius + SHIP_DOCKING_RADIUS + .01)
}

func (ship *Ship) WillPathCollideWithPlanet(heading *Heading, planet *Planet) bool {
	if heading.Magnitude == 0 {
		return false
	}
	if ship.DistanceToCollision(&planet.Entity) > float64(heading.Magnitude) {
		return false
	}


	return true
}

func (ship *Ship) Navigate(target *Entity, gameMap GameMap) string {

	ob := gameMap.ObstaclesBetween(&ship.Entity, target)

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
				intermediateTarget := Point{
					X: x1,
					Y: y1,
				}
				intermediateEntity := Entity{
					Point: intermediateTarget,
				}
				ob1 := gameMap.ObstaclesBetween(&ship.Entity, &intermediateEntity)
				if !ob1 {
					ob2 := gameMap.ObstaclesBetween(&intermediateEntity, target)
					if !ob2 {
						totdist := math.Sqrt(math.Pow(x1-x0, 2)+math.Pow(y1-y0, 2)) + math.Sqrt(math.Pow(x1-x2, 2)+math.Pow(y1-y2, 2))
						if totdist < bestdist {
							bestdist = totdist
							bestTarget = &intermediateEntity

						}
					}
				}
			}
		}
		return ship.NavigateBasic(bestTarget)
	}

}

func (ship *Ship) NavigateSnail(speed int, angle float64, gameMap *GameMap) Heading {
	log.Println("NavigateSnail with speed ", speed, " and angle ", angle)

	// add points along the path
	for mag := 1; mag <= speed; mag++ {
		intermediatePos := ship.Entity.AddThrust(float64(mag), angle)
		intermediateEntity := Entity{
			Point:  intermediatePos,
			Radius: ship.Entity.Radius,
			Owner:  -1,
		}
		gameMap.Entities = append(gameMap.Entities, intermediateEntity)
	}

	return CreateHeading(speed, angle)
}

func (ship *Ship) BetterNavigate(target *Entity, gameMap *GameMap) Heading {
	log.Println("betternavigation from ", ship.Point, " to ", target.Point, " with id ", target.Id)

	maxTurn := (3 * math.Pi) / 2
	dTurn := math.Pi / 8

	startSpeed := int(math.Min(SHIP_MAX_SPEED, ship.Point.DistanceTo(&target.Point)-target.Radius-ship.Radius-.05))
	log.Println("setting start speed to ", startSpeed)
	baseAngle := ship.Point.AngleTo(&target.Point)

	if !gameMap.ObstaclesInPath(&ship.Entity, float64(startSpeed), baseAngle) {
		log.Println("Way is clear to planet!")
		return ship.NavigateSnail(startSpeed, baseAngle, gameMap)
	}

	for speed := startSpeed; speed >= 1; speed -- {
		log.Println("Trying speed, ", speed)
		for turn := dTurn; turn <= maxTurn; turn += dTurn {
			log.Println("Trying turn, ", turn)
			intermediateTargetLeft := ship.AddThrust(float64(speed), baseAngle+turn)
			obLeft := gameMap.ObstaclesInPath(&ship.Entity, float64(speed), baseAngle+turn)
			intermediateTargetRight := ship.AddThrust(float64(speed), baseAngle-turn)
			obRight := gameMap.ObstaclesInPath(&ship.Entity, float64(speed), baseAngle-turn)
			if !obLeft && !obRight {
				if intermediateTargetLeft.SqDistanceTo(&target.Point) < intermediateTargetRight.SqDistanceTo(&target.Point) {
					return ship.NavigateSnail(speed, baseAngle+turn, gameMap)
				} else {
					return ship.NavigateSnail(speed, baseAngle-turn, gameMap)
				}
			} else if !obLeft {
				return ship.NavigateSnail(speed, baseAngle+turn, gameMap)
			} else if !obRight {
				return ship.NavigateSnail(speed, baseAngle-turn, gameMap)
			}
		}
	}
	return Heading{
		Magnitude: 0,
		Angle:     0,
	}
}
