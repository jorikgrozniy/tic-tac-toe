package service

import "tic-tac-toe/internal/domain/model"

func CheckGameStatus(game *model.CurrentGame) model.GameStatus {
	if !game.Players.IsBoth() {
		return model.GameStatus{
			Type: model.StatusWaiting,
		}
	} else {
		switch checkWin(game.Board) {
		case model.BoardX:
			return model.GameStatus{
				Type:   model.StatusWin,
				Player: game.Players.GetX(),
			}
		case model.BoardO:
			return model.GameStatus{
				Type:   model.StatusWin,
				Player: game.Players.GetO(),
			}
		}

		if isBoardFilled(game.Board) {
			return model.GameStatus{
				Type: model.StatusDraw,
			}
		}

		x, o := countFlags(game.Board)
		switch x {
		case o:
			return model.GameStatus{
				Type:   model.StatusTurn,
				Player: game.Players.GetX(),
			}
		case o + 1:
			return model.GameStatus{
				Type:   model.StatusTurn,
				Player: game.Players.GetO(),
			}
		}
	}

	return model.GameStatus{}
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
		if comb[0] != model.BoardEmpty && comb[0] == comb[1] && comb[0] == comb[2] {
			return comb[0]
		}
	}
	return 0
}

func isBoardFilled(board model.GameBoard) bool {
	x, o := countFlags(board)
	return x+o == 9
}

func countFlags(board model.GameBoard) (int, int) {
	x, o := 0, 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch board[i][j] {
			case model.BoardX:
				x++
			case model.BoardO:
				o++
			}
		}
	}
	return x, o
}
