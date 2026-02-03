package service

import (
	"errors"

	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
)

func FindBestMove(game *model.CurrentGame, player model.GamePlayer) (int, int, error) {
	bestScore, bestRow, bestCol := -11, -1, -1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if game.Board[i][j] == model.BoardEmpty {
				game.Board[i][j] = player.Flag
				score := minimax(game, player, 0, false)
				game.Board[i][j] = model.BoardEmpty

				if score > bestScore {
					bestScore, bestRow, bestCol = score, i, j
				}
			}
		}
	}
	if bestScore == -11 || bestRow == -1 || bestCol == -1 {
		return -1, -1, errors.New("no valid moves available")
	}
	return bestRow, bestCol, nil
}

func minimax(game *model.CurrentGame, player model.GamePlayer, depth int, isMaximizing bool) int {
	status := CheckGameStatus(game)
	switch status.Type {
	case model.StatusWin:
		if *status.Player == player {
			return 10 - depth
		} else {
			return depth - 10
		}
	case model.StatusDraw:
		return 0
	}

	if isMaximizing {
		bestScore := -11
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if game.Board[i][j] == model.BoardEmpty {
					game.Board[i][j] = player.Flag
					score := minimax(game, player, depth+1, false)
					game.Board[i][j] = model.BoardEmpty
					bestScore = max(bestScore, score)
					if bestScore == 10 {
						return bestScore
					}
				}
			}
		}
		return bestScore
	} else {
		leastScore := 11
		opp := game.Players.GetOpponent(player).Flag

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if game.Board[i][j] == model.BoardEmpty {
					game.Board[i][j] = opp
					score := minimax(game, player, depth+1, true)
					game.Board[i][j] = model.BoardEmpty
					leastScore = min(leastScore, score)
					if leastScore == -10 {
						return leastScore
					}
				}
			}
		}
		return leastScore
	}
}
