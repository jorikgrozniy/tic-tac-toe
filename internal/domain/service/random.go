package service

import (
	"math/rand"

	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
)

func SetRandomPlayerFlags(p *model.GamePlayers) {
	p.P1.Flag = getRandomFlag()
	if p.P1.Flag == model.BoardX {
		p.P2.Flag = model.BoardO
	} else {
		p.P2.Flag = model.BoardX
	}
}

func getRandomFlag() int {
	return 1 + rand.Intn(2)
}
