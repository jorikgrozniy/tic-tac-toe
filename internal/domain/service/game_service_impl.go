package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/repository"
)

var (
	ErrPrevBoardNotFound = errors.New("previous board not found")
	ErrBoardAltered      = errors.New("the board was altered")
	ErrOpponentMoved     = errors.New("the opponent's move was made")
	ErrTooManyMoves      = errors.New("more than one move was made")
	ErrMoveAlreadyMade   = errors.New("this move was already made")
	ErrMoveWhileWaiting  = errors.New("cannot move while waiting")
	ErrNotYourTurn       = errors.New("not your turn")
	ErrAccessDenied      = errors.New("access denied to this game")
	ErrGameOver          = errors.New("the game is over")
	ErrNoSuchGameType    = errors.New("no such game type")
	ErrInternal          = errors.New("internal error")
	ErrGameInProgress    = errors.New("game in progress")
	ErrInitPlayers       = errors.New("error initializing players")
	ErrCannotJoinGame    = errors.New("cannot join this game")
	ErrNotPVPGame        = errors.New("not a pvp game")
)

type gameServiceImpl struct {
	repo repository.GameRepository
}

func NewGameService(repo repository.GameRepository) GameService {
	return &gameServiceImpl{
		repo: repo,
	}
}

func (s *gameServiceImpl) ComputeNextMove(game *model.CurrentGame, player model.GamePlayer) (int, int, error) {
	return FindBestMove(game, player)
}

func (s *gameServiceImpl) ValidateGameBoard(game *model.CurrentGame, player model.GamePlayer) error {
	prevGame, err := s.repo.Retrieve(game.ID)
	if err != nil {
		return ErrPrevBoardNotFound
	}

	identical := true
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if prevGame.Board[i][j] != model.BoardEmpty && prevGame.Board[i][j] != game.Board[i][j] {
				return ErrBoardAltered
			}
			if prevGame.Board[i][j] == model.BoardEmpty && game.Board[i][j] == game.Players.GetOpponent(player).Flag {
				return ErrOpponentMoved
			}
			if prevGame.Board[i][j] != game.Board[i][j] {
				if !identical {
					return ErrTooManyMoves
				}
				identical = false
			}
		}
	}
	if identical {
		return ErrMoveAlreadyMade
	}
	return nil
}

func (s *gameServiceImpl) CheckGameCompletion(game *model.CurrentGame) model.GameStatus {
	return CheckGameStatus(game)
}

func (s *gameServiceImpl) MakeMove(reqGame *model.CurrentGame, userID uuid.UUID) error {
	prevGame, err := s.repo.Retrieve(reqGame.ID)
	if err != nil {
		return err
	}

	player := prevGame.Players.GetByID(userID)
	if player == nil {
		return ErrAccessDenied
	}

	opponent := reqGame.Players.GetOpponent(*player)
	reqGame.Players = prevGame.Players
	reqGame.Type = prevGame.Type

	if err := s.ValidateGameBoard(reqGame, *player); err != nil {
		return err
	}

	reqGame.Status = s.CheckGameCompletion(reqGame)
	switch prevGame.Status.Type {
	case model.StatusWaiting:
		return ErrMoveWhileWaiting
	case model.StatusTurn:
		if *prevGame.Status.Player == *player {
			if opponent.IsExist() && opponent.IsBot() {
				row, col, err := s.ComputeNextMove(reqGame, *opponent)
				if err != nil {
					return err
				}

				s.placeFlag(&reqGame.Board, row, col, opponent.Flag)
				reqGame.Status = s.CheckGameCompletion(reqGame)
			}

			if err := s.repo.Save(reqGame); err != nil {
				return err
			}
		} else {
			return ErrNotYourTurn
		}
	default:
		return ErrGameOver
	}

	return nil
}

func (s *gameServiceImpl) CreateGame(userID uuid.UUID, t string) (*model.CurrentGame, error) {
	gameType := model.GameType(t)
	if gameType != model.TypePVE && gameType != model.TypePVP {
		return nil, ErrNoSuchGameType
	}

	game, err := s.newGame(userID, gameType)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(game); err != nil {
		return nil, err
	}

	return game, nil
}

func (s *gameServiceImpl) GetGame(gameID string, userID uuid.UUID) (*model.CurrentGame, error) {
	parsedID, err := uuid.Parse(gameID)
	if err != nil {
		return nil, err
	}

	game, err := s.repo.Retrieve(parsedID)
	if err != nil {
		return nil, err
	}

	player := game.Players.GetByID(userID)
	if player != nil || (game.Status.Type == model.StatusWaiting && game.Type == model.TypePVP) {
		return game, nil
	}

	return nil, ErrAccessDenied
}

func (s *gameServiceImpl) GetAvailableGames() ([]model.GameInfo, int, error) {
	games, err := s.repo.FindAvailableGames()
	return games, len(games), err
}

func (s *gameServiceImpl) JoinGame(gameID string, userID uuid.UUID) error {
	parsedGameID, err := uuid.Parse(gameID)
	if err != nil {
		return err
	}

	game, err := s.repo.Retrieve(parsedGameID)
	if err != nil {
		return err
	}

	if game.Type != model.TypePVP {
		return ErrNotPVPGame
	}

	switch game.Status.Type {
	case model.StatusWaiting:
		if !game.Players.IsBoth() {
			player := model.GamePlayer{
				ID:   &userID,
				IsAI: false,
				Flag: model.BoardEmpty,
			}

			game.Players.AddPlayer(player)
			if err := s.initPlayers(game); err != nil {
				return ErrCannotJoinGame
			}
		} else {
			return ErrInternal
		}
	default:
		return ErrGameInProgress
	}

	if err := s.repo.Save(game); err != nil {
		return err
	}

	return nil
}

func (s *gameServiceImpl) initPlayers(game *model.CurrentGame) error {
	switch game.Type {
	case model.TypePVE:
		if game.Status.Type == model.StatusWaiting {
			game.Players.P1.Flag = model.BoardX
			game.Players.P2.Flag = model.BoardO
			game.Players.P2.IsAI = true
			game.Status.Type = model.StatusTurn
			game.Status.Player = &game.Players.P1
		} else {
			return ErrInitPlayers
		}
	case model.TypePVP:
		if game.Players.IsBoth() && !game.Players.P1.IsBot() && !game.Players.P2.IsBot() && game.Status.Type == model.StatusWaiting {
			SetRandomPlayerFlags(&game.Players)
			game.Status.Type = model.StatusTurn
			game.Status.Player = game.Players.GetX()
		} else {
			return ErrInitPlayers
		}
	default:
		return ErrNoSuchGameType
	}

	return nil
}

func (s *gameServiceImpl) newGame(userID uuid.UUID, t model.GameType) (*model.CurrentGame, error) {
	game := &model.CurrentGame{
		ID: uuid.New(),
		Board: model.GameBoard{
			{model.BoardEmpty, model.BoardEmpty, model.BoardEmpty},
			{model.BoardEmpty, model.BoardEmpty, model.BoardEmpty},
			{model.BoardEmpty, model.BoardEmpty, model.BoardEmpty},
		},
		Status: model.GameStatus{
			Type: model.StatusWaiting,
		},
		Players: model.GamePlayers{
			P1: model.GamePlayer{
				ID:   &userID,
				IsAI: false,
				Flag: model.BoardEmpty,
			},
		},
		Type: t,
	}

	var err error
	if t == model.TypePVE {
		err = s.initPlayers(game)
	}

	return game, err
}

func (s *gameServiceImpl) placeFlag(board *model.GameBoard, row, col, player int) {
	board[row][col] = player
}
