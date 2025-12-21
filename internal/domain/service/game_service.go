package service

import (
	"github.com/google/uuid"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
)

type GameService interface {
	ComputeNextMove(board model.GameBoard, player int) (int, int, error)
	ValidateGameBoard(current model.GameBoard, gameID uuid.UUID) error
	CheckGameCompletion(board model.GameBoard) bool
}
