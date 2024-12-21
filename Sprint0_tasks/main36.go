package main

import (
	"time"
)

type Task struct {
	summary     string
	description string
	deadline    time.Time
	priority    int
}
type Note struct {
	title string
	text  string
}
type ToDoList struct {
	name  string
	tasks []Task
	notes []Note
}

func (td ToDoList) TasksCount() int {
	return len(td.tasks)
}

func (td ToDoList) NotesCount() int {
	return len(td.notes)
}

func (ctt ToDoList) CountTopPrioritiesTasks() int {
	count := 0
	for _, task := range ctt.tasks {
		if task.IsTopPriority() { // Вызываем метод у объекта task
			count++
		}
	}
	return count
}

func (ctt ToDoList) CountOverdueTasks() int {
	count := 0
	for _, task := range ctt.tasks {
		if task.IsOverdue() { // Вызываем метод у объекта task
			count++
		}
	}
	return count
}

func (t Task) IsOverdue() bool {
	return time.Now().After(t.deadline)
}

func (t Task) IsTopPriority() bool {
	return t.priority >= 4
}
