package hlt

import (
	"fmt"
	"log"
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
	Distance        float64
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


func (ship *Ship) CanDock(planet *Planet) bool {
	dist := ship.DistanceToCollision(&planet.Entity)

	return dist <= SHIP_DOCKING_RADIUS
}

