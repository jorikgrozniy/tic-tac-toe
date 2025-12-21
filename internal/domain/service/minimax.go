package service

import (
	"errors"

	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
)

func FindBestMove(board model.GameBoard, player int) (int, int, error) {
	bestScore, bestRow, bestCol := -11, -1, -1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == model.Empty {
				board[i][j] = player
				score := minimax(board, player, 0, false)
				board[i][j] = model.Empty

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

func minimax(board model.GameBoard, player, depth int, isMaximizing bool) int {
	status := CheckGameStatus(board)
	if status == player {
		return 10 - depth
	} else if status != Draw && status != Playing {
		return depth - 10
	} else if status == Draw {
		return 0
	}

	if isMaximizing {
		bestScore := -11
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == model.Empty {
					board[i][j] = player
					score := minimax(board, player, depth+1, false)
					board[i][j] = model.Empty
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
		opp := model.PlayerX
		if player == model.PlayerX {
			opp = model.PlayerO
		}

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == model.Empty {
					board[i][j] = opp
					score := minimax(board, player, depth+1, true)
					board[i][j] = model.Empty
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
