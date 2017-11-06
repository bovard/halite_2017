package strat

type ChlMessage int

const (
	NONE  ChlMessage = iota
	CANCELLED_PLANET_ASSIGNMENT_MIN
	CANCELLED_PLANET_ASSIGNMENT_TOO_CLOSE
	CANCELLED_PLANET_ASSIGNMENT_PLANET_TAKEN
	MOVING_TOWARD_PLANET
	MOVING_TOWARD_ENEMY
	COMBAT_WE_OUTNUMBER
	COMBAT_TIED
	COMBAT_OUTNUMBERED
	COMBAT_KILL_PRODUCTION
	COMBAT_OUTNUMBERED_AND_FAR_FROM_HOME
	COMBAT_SUICIDE_DUE_TO_LOWER_HEALTH
	COMBAT_TIED_SUICIDE_TO_GAIN_VALUE
)