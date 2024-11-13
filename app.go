package main

import (
	"cmp"
	"context"
	"math/rand/v2"
	"slices"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Field struct {
	turn          int
	Mark          mark `json:"mark"`
	WillBeRemoved bool `json:"willBeRemoved"`
}

type mark string

const (
	empty    mark = ""
	user          = "O"
	computer      = "X"
)

type App struct {
	ctx    context.Context
	turn   int
	Fields [3][3]Field `json:"fields"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) getRemoveMark(m mark) (int, int) {
	turns := make([][3]int, 0, 3)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if a.Fields[i][j].Mark == m {
				turns = append(turns, [3]int{i, j, a.Fields[i][j].turn})
			}
		}
	}

	slices.SortFunc(turns, func(a, b [3]int) int {
		return cmp.Compare(a[2], b[2])
	})
	return turns[0][0], turns[0][1]
}

func (a *App) isFull(m mark) bool {
	count := 0

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if a.Fields[i][j].Mark == m {
				count++
			}
		}
	}
	return count >= 3
}

func (a *App) nextComputerCoord() (int, int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if a.Fields[i][j].Mark == empty {
				if a.checkWin(computer, i, j) {
					return i, j
				}
			}
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if a.Fields[i][j].Mark == empty {
				if a.checkWin(user, i, j) {
					return i, j
				}
			}
		}
	}

	emptyCoords := make([][2]int, 0, 9)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if a.Fields[i][j].Mark == empty {
				emptyCoords = append(emptyCoords, [2]int{i, j})
			}
		}
	}

	emptyCoord := emptyCoords[rand.IntN(len(emptyCoords))]
	return emptyCoord[0], emptyCoord[1]
}

func (a *App) checkWin(m mark, coord ...int) bool {
	if len(coord) == 2 {
		prev := a.Fields[coord[0]][coord[1]].Mark
		a.Fields[coord[0]][coord[1]].Mark = m
		defer func() {
			a.Fields[coord[0]][coord[1]].Mark = prev
		}()
	}

	for i := 0; i < 3; i++ {
		if a.Fields[i][0].Mark == m && a.Fields[i][1].Mark == m && a.Fields[i][2].Mark == m {
			return true
		}
	}

	for j := 0; j < 3; j++ {
		if a.Fields[0][j].Mark == m && a.Fields[1][j].Mark == m && a.Fields[2][j].Mark == m {
			return true
		}
	}

	if a.Fields[0][0].Mark == m && a.Fields[1][1].Mark == m && a.Fields[2][2].Mark == m {
		return true
	}

	if a.Fields[0][2].Mark == m && a.Fields[1][1].Mark == m && a.Fields[2][0].Mark == m {
		return true
	}
	return false
}

func (a *App) ClearFields() {
	a.turn = 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			a.Fields[i][j].Mark = empty
			a.Fields[i][j].WillBeRemoved = false
		}
	}
}

func (a *App) CheckUserWin() bool {
	win := a.checkWin(user)
	if win {
		result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.QuestionDialog,
			Title:         "You win!",
			Message:       "Do you want to reset the game?",
			DefaultButton: "Yes",
			Buttons:       []string{"Yes", "No"},
		})
		if err != nil {
			return win
		}

		if result == "Yes" {
			a.ClearFields()
		}
	}
	return win
}

func (a *App) CheckComputerWin() bool {
	win := a.checkWin(computer)
	if win {
		result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.QuestionDialog,
			Title:         "Computer win!",
			Message:       "Do you want to reset the game?",
			DefaultButton: "Yes",
			Buttons:       []string{"Yes", "No"},
		})
		if err != nil {
			return win
		}

		if result == "Yes" {
			a.ClearFields()
		}
	}
	return win
}

func (a *App) GetFields() [3][3]Field {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			a.Fields[i][j].WillBeRemoved = false
		}
	}

	if a.isFull(user) {
		i, j := a.getRemoveMark(user)
		a.Fields[i][j].WillBeRemoved = true
	}
	if a.isFull(computer) {
		i, j := a.getRemoveMark(computer)
		a.Fields[i][j].WillBeRemoved = true
	}
	return a.Fields
}

func (a *App) SetComputerMark() {
	row, col := a.nextComputerCoord()

	if a.isFull(computer) {
		i, j := a.getRemoveMark(computer)
		a.Fields[i][j].Mark = empty
	}

	a.turn++
	a.Fields[row][col].Mark = computer
	a.Fields[row][col].turn = a.turn
}

func (a *App) SetUserMark(n int) bool {
	row, col := n/3, n%3
	if a.Fields[row][col].Mark != empty {
		return false
	}

	if a.isFull(user) {
		i, j := a.getRemoveMark(user)
		a.Fields[i][j].Mark = empty
	}

	a.turn++
	a.Fields[row][col].Mark = user
	a.Fields[row][col].turn = a.turn
	return true
}
