package service

import "github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"

const (
	Playing = -1
	Draw    = 0
	WinX    = 1
	WinO    = 2
)

func CheckGameStatus(board model.GameBoard) int {
	winner := checkWin(board)
	switch winner {
	case model.PlayerX:
		return WinX
	case model.PlayerO:
		return WinO
	default:
		if checkOverflow(board) {
			return Draw
		}
	}
	return Playing
}

func checkWin(board model.GameBoard) int {
	winCombinations := [8][3]int{
		{board[0][0], board[0][1], board[0][2]},
		{board[1][0], board[1][1], board[1][2]},
		{board[2][0], board[2][1], board[2][2]},

		{board[0][0], board[1][0], board[2][0]},
		{board[0][1], board[1][1], board[2][1]},
		{board[0][2], board[1][2], board[2][2]},

		{board[0][0], board[1][1], board[2][2]},
		{board[0][2], board[1][1], board[2][0]},
	}

	for _, comb := range winCombinations {
		if comb[0] != model.Empty && comb[0] == comb[1] && comb[0] == comb[2] {
			return comb[0]
		}
	}
	return 0
}

func checkOverflow(board model.GameBoard) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == model.Empty {
				return false
			}
		}
	}
	return true
}
