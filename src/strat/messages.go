package strat

type ChlMessage int

const (
	NONE  ChlMessage = iota
	CANCELLED_PLANET_ASSIGNMENT
	MOVING_TOWARD_PLANET
	MOVING_TOWARD_ENEMY
)