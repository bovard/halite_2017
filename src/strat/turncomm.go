package strat

type TurnComm struct {
	Chasing map[int]int
}

func GetTurnComm() TurnComm {
	return TurnComm{
		Chasing: make(map[int]int),
	}
}
