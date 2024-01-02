package todolist

import (
	"animedown/core/search"
	"fmt"
)

type TodoList struct {
	list [][]string
}

func New() *TodoList {
	return &TodoList{
		list: nil,
	}
}

func (t *TodoList) Add(row []string) error {
	if t.IsExist(row) {
		return fmt.Errorf("row already exist")
	}
	t.list = append(t.list, row)
	return nil
}

func (t *TodoList) Delete(index int) {
	t.list = append(t.list[:index], t.list[index+1:]...)
}

func (t *TodoList) Get(index int) []string {
	return t.list[index]
}

func (t *TodoList) GetList() [][]string {
	return t.list
}

func (t *TodoList) Len() int {
	return len(t.list)
}

func (t *TodoList) Clear() {
	t.list = nil
}

func (t *TodoList) IsEmpty() bool {
	return len(t.list) == 0
}

func (t *TodoList) IsExist(row []string) bool {
	magnetIndex := search.GetAllFilterIndex(search.FilterMagnet)
	for _, r := range t.list {
		if r[magnetIndex] == row[magnetIndex] {
			return true
		}
	}
	return false
}

func (t *TodoList) GetMagnet(index int) string {
	magnetIndex := search.GetAllFilterIndex(search.FilterMagnet)
	return t.list[index][magnetIndex]
}
