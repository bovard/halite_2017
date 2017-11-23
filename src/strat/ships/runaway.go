package ships

import (
	"../../hlt"
)

func (self *ShipController) RunAwaySetTarget(gameMap *hlt.GameMap) {
	self.Target = &self.Info.ClosestEnemy.Ship.Point
}

func (self *ShipController) RunAwayAct(gameMap *hlt.GameMap, turnComm *TurnComm) (ChlMessage, hlt.Heading) {
	return self.stupidRunAwayMeta(gameMap)
}

func (self *ShipController) runAway(gameMap *hlt.GameMap) (ChlMessage, hlt.Heading) {
	heading := hlt.Heading{
		Magnitude: 0,
		Angle:     0,
	}
	message := RUN_AWAY

	dir := self.Info.ClosestEnemy.Ship.AngleTo(&self.Ship.Point)
	targetPos := self.Ship.Point.AddThrust(hlt.SHIP_MAX_SPEED, dir)

	heading = self.MoveToPoint(&targetPos, gameMap)

	return message, heading
}

func (self *ShipController) stupidRunAwayMeta(gameMap *hlt.GameMap) (ChlMessage, hlt.Heading) {
	/*
		heading := hlt.Heading{
			Magnitude: 0,
			Angle:     0,
		}
		message := HIDE_WE_ARE_LOSING
	*/
	// TODO: head to corner

	return self.runAway(gameMap)
}
