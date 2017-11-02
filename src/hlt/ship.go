package hlt

import (
	"math"
	"strconv"
	"fmt"
	"log"
)

type DockingStatus int

const (
	UNDOCKED  DockingStatus = iota
	DOCKING
	DOCKED
	UNDOCKING
)

type Ship struct {
	Entity
	Vel Vector
	NextVel Vector
	PlanetId        int
	Planet          Planet
	DockingStatus   DockingStatus
	DockingProgress float64
	WeaponCooldown  float64
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

	shipPoint := Point {
		X: shipX,
		Y: shipY,
	}

	shipEntity := Entity{
		Point: shipPoint,
		Radius: .5,
		Health: shipHealth,
		Owner:  playerId,
		Id:     shipId,
	}

	shipVel := Vector {
		X: shipVelX,
		Y: shipVelY,
	}

	nextVel := Vector { 
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

func (ship *Ship) WillPathCollideWithPlanet(thrust float64, angle float64, planet *Planet) bool {
	if (thrust == 0) {
		return false
	}

	return true
}

func (ship *Ship) Navigate(target *Entity, gameMap Map) string {


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
					X:      x1,
					Y:      y1,
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


func (ship *Ship) NavigateSnail(target *Point, gameMap *Map) string {

	maxMove := ship.Point.DistanceTo(target) - (ship.Entity.Radius + .1)

	angle := ship.AngleTo(target)
	speed := math.Min(maxMove, SHIP_MAX_SPEED)

	// add points along the path
	for mag := 1.0; mag <= speed; mag++ {
		intermediatePos := ship.Entity.AddThrust(mag, angle)
		intermediateEntity := Entity {
			Point:  intermediatePos,
			Radius: ship.Entity.Radius,
			Owner: -1,
		}
		gameMap.Entities = append(gameMap.Entities, intermediateEntity)
	}
	// add the terminal location 
	intermediatePos := ship.Entity.AddThrust(speed, angle)
	intermediateEntity := Entity {
		Point:  intermediatePos,
		Radius: ship.Entity.Radius,
		Owner: -1,
	}
	gameMap.Entities = append(gameMap.Entities, intermediateEntity)

	return ship.Thrust(speed, angle)
}

func (ship *Ship) BetterNavigate(target *Entity, gameMap *Map) string {
	log.Println("betternavigation from ", ship.Point, " to ", target.Point, " with id ", target.Id)
	
	maxTurn := (3 * math.Pi) / 2
	dTurn := math.Pi / 8

	startSpeed := math.Min(SHIP_MAX_SPEED, ship.Point.DistanceTo(&target.Point) - target.Radius - ship.Radius - .05)
	log.Println("setting start speed to ", startSpeed)
	baseAngle := ship.Point.AngleTo(&target.Point)

	intermediateTarget := ship.Entity.AddThrust(startSpeed, baseAngle)
	if !gameMap.ObstaclesInPath(&ship.Entity, startSpeed, baseAngle) {
		log.Println("Way is clear to planet!")
		return ship.NavigateSnail(&intermediateTarget, gameMap)
	}

	for speed := startSpeed; speed > .25; speed /= 2 {
		log.Println("Trying speed, ", speed)
		for turn := dTurn; turn <= maxTurn; turn += dTurn {
			log.Println("Trying turn, ", turn)
			intermediateTargetLeft := ship.AddThrust(speed, baseAngle + turn)
			obLeft := gameMap.ObstaclesInPath(&ship.Entity, speed, baseAngle + turn)
			intermediateTargetRight := ship.AddThrust(speed, baseAngle - turn)
			obRight := gameMap.ObstaclesInPath(&ship.Entity, speed, baseAngle - turn)
			if !obLeft && !obRight {
				if intermediateTargetLeft.SqDistanceTo(&target.Point) < intermediateTargetRight.SqDistanceTo(&target.Point) {
					return ship.NavigateSnail(&intermediateTargetLeft, gameMap)
				} else {
					return ship.NavigateSnail(&intermediateTargetRight, gameMap)
				}
			} else if !obLeft {
				return ship.NavigateSnail(&intermediateTargetLeft, gameMap)
			} else if !obRight {
				return ship.NavigateSnail(&intermediateTargetRight, gameMap)
			}
		}
	}
	return ""
}