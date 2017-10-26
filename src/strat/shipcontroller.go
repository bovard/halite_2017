package strat

import ( "../hlt")

type ShipController struct {
	Ship     hlt.Ship
	Past     [] hlt.Ship
	Id       int
	Planet   int
	Alive    bool
}

func (self ShipController) Update(ship hlt.Ship) {
	self.Past = append(self.Past, self.Ship)
	self.Ship = ship
}