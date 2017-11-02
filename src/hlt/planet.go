package hlt

import (
	"strconv"
)

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

	planetPoint := Point{
		X: planetX,
		Y: planetY,
	}

	planetEntity := Entity{
		Point:  planetPoint,
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