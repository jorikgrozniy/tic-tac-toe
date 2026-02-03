package model

import (
	"github.com/google/uuid"
)

type CurrentGame struct {
	ID      uuid.UUID
	Board   GameBoard
	Status  GameStatus
	Players GamePlayers
	Type    GameType
}

type GamePlayer struct {
	ID   *uuid.UUID
	IsAI bool
	Flag int
}

func (p *GamePlayer) IsBot() bool {
	return p.ID == nil && p.IsAI
}

func (p *GamePlayer) IsExist() bool {
	return p != nil
}

type GamePlayers struct {
	P1 GamePlayer
	P2 GamePlayer
}

func (p *GamePlayers) IsBoth() bool {
	return (p.P1.ID != nil || p.P1.IsAI) && (p.P2.ID != nil || p.P2.IsAI)
}

func (p *GamePlayers) AddPlayer(player GamePlayer) {
	if p.P1.ID == nil && !p.P1.IsAI {
		p.P1 = player
	} else if p.P2.ID == nil && !p.P2.IsAI {
		p.P2 = player
	}
}

func (p *GamePlayers) GetByID(id uuid.UUID) *GamePlayer {
	if *p.P1.ID == id {
		return &p.P1
	} else if *p.P2.ID == id {
		return &p.P2
	}
	return nil
}

func (p *GamePlayers) GetX() *GamePlayer {
	if p.P1.Flag == BoardX {
		return &p.P1
	} else if p.P2.Flag == BoardX {
		return &p.P2
	}
	return nil
}

func (p *GamePlayers) GetO() *GamePlayer {
	if p.P1.Flag == BoardO {
		return &p.P1
	} else if p.P2.Flag == BoardO {
		return &p.P2
	}
	return nil
}

func (p *GamePlayers) GetOpponent(player GamePlayer) *GamePlayer {
	switch player {
	case p.P1:
		return &p.P2
	case p.P2:
		return &p.P1
	}
	return nil
}

type GameStatus struct {
	Type   GameStatusType
	Player *GamePlayer
}

type GameStatusType string

const (
	StatusWaiting GameStatusType = "waiting"
	StatusTurn    GameStatusType = "turn"
	StatusWin     GameStatusType = "win"
	StatusDraw    GameStatusType = "draw"
)

type GameType string

const (
	TypePVP GameType = "pvp"
	TypePVE GameType = "pve"
)
