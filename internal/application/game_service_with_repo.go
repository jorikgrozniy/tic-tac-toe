package application

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/service"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/repository"
)

type GameServiceWithRepo struct {
	Repo repository.GameRepository
}

func NewGameAppService(repo repository.GameRepository) *GameServiceWithRepo {
	return &GameServiceWithRepo{
		Repo: repo,
	}
}

func (s *GameServiceWithRepo) ComputeNextMove(board model.GameBoard, player int) (int, int, error) {
	return service.FindBestMove(board, player)
}

func (s *GameServiceWithRepo) ValidateGameBoard(current model.GameBoard, gameID uuid.UUID) error {
	prevGame, err := s.Repo.Retrieve(gameID)
	if err != nil {
		return errors.New("previous board not found")
	}

	identical := true
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if prevGame.Board[i][j] != model.Empty && prevGame.Board[i][j] != current[i][j] {
				return errors.New("the board was altered")
			}
			// hardcoded "X" for player
			if prevGame.Board[i][j] == model.Empty && current[i][j] == model.PlayerO {
				return errors.New("the opponent's move was made")
			}
			if prevGame.Board[i][j] != current[i][j] {
				if !identical {
					return errors.New("more than one move was made")
				}
				identical = false
			}
		}
	}
	if identical {
		return errors.New("this move was already made")
	}
	return nil
}

func (s *GameServiceWithRepo) CheckGameCompletion(board model.GameBoard) int {
	return service.CheckGameStatus(board)
}

func (s *GameServiceWithRepo) ProcessMakeMove(reqGame *model.CurrentGame) (*model.CurrentGame, int, error) {
	if err := s.ValidateGameBoard(reqGame.Board, reqGame.ID); err != nil {
		return nil, -2, err
	}

	status := s.CheckGameCompletion(reqGame.Board)
	if status == service.Playing {
		// hardcoded "O" for AI
		row, col, err := s.ComputeNextMove(reqGame.Board, model.PlayerO)
		if err != nil {
			return nil, -2, err
		}
		// again
		s.makeMove(&reqGame.Board, row, col, model.PlayerO)
		status = s.CheckGameCompletion(reqGame.Board)

		if err := s.Repo.Save(reqGame); err != nil {
			return nil, -2, err
		}
	} else if prevGame, _ := s.Repo.Retrieve(reqGame.ID); s.CheckGameCompletion(prevGame.Board) != service.Playing {
		return nil, -2, errors.New("the game is over")
	}
	return reqGame, status, nil
}

func (s *GameServiceWithRepo) ProcessCreateGame() (*model.CurrentGame, error) {
	game := s.newGame()
	if err := s.Repo.Save(game); err != nil {
		return nil, err
	}
	return game, nil
}

func (s *GameServiceWithRepo) ProcessGetGame(gameID string) (*model.CurrentGame, int, error) {
	UUID, err := uuid.Parse(gameID)
	if err != nil {
		return nil, -2, err
	}

	game, err := s.Repo.Retrieve(UUID)
	if err != nil {
		return nil, -2, err
	}

	status := s.CheckGameCompletion(game.Board)
	return game, status, nil
}

func (s *GameServiceWithRepo) makeMove(board *model.GameBoard, row, col, player int) {
	board[row][col] = player
}

func (s *GameServiceWithRepo) newGame() *model.CurrentGame {
	return &model.CurrentGame{
		ID: uuid.New(),
		Board: model.GameBoard{
			{model.Empty, model.Empty, model.Empty},
			{model.Empty, model.Empty, model.Empty},
			{model.Empty, model.Empty, model.Empty},
		},
	}
}
