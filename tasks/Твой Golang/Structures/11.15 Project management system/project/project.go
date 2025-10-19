package project

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	StatusActive bool = true
	StatusClosed bool = false
)

type Task struct {
	Id          uuid.UUID
	Title       string
	Description string
	Status      bool
}

type Project struct {
	Id    uuid.UUID
	Name  string
	Tasks []Task
}

func New(id uuid.UUID, name string) (*Project, error) {
	if len(name) == 0 {
		return nil, errors.New("создание структуры с нулевой длиной")
	}

	return &Project{
		Id:    id,
		Name:  name,
		Tasks: nil,
	}, nil
}

func NewTask(id uuid.UUID, title string, description string) (*Task, error) {
	if len(title) == 0 {
		return nil, errors.New("у задачи отсутствует заголовок")
	}

	if len(description) == 0 {
		return nil, errors.New("у задачи отсутствует описание")
	}

	return &Task{
		Id:          id,
		Title:       title,
		Description: description,
		Status:      StatusActive,
	}, nil
}

func (p *Project) AddTask(task Task) error {
	for _, val := range p.Tasks {
		if val.Id == task.Id {
			return errors.New("id добавляемой задачи уже существует в списке")
		}
	}

	p.Tasks = append(p.Tasks, task)

	return nil
}

func (p *Project) UpdateTask(task Task) error {
	for i := range p.Tasks {
		if p.Tasks[i].Id == task.Id {
			p.Tasks[i].Description = task.Description
			p.Tasks[i].Title = task.Title
			p.Tasks[i].Status = task.Status
			return nil
		}
	}

	return errors.New("заданной задачи не существует")
}

func (t *Task) Close() error {
	if t.Status {
		t.Status = StatusClosed
		return nil
	}

	return errors.New("задача уже была выполнена")
}

func (t *Task) UpdateDescription(description string) error {
	if len(description) == 0 {
		return errors.New("новое описание пустое")
	}

	t.Description = description
	return nil
}

func (p Project) FilterTasksByStatus(status bool) []Task {
	tasks := []Task{}

	for _, val := range p.Tasks {
		if val.Status == status {
			tasks = append(tasks, val)
		}
	}

	return tasks
}

func (p Project) PrintInfo() {
	fmt.Printf("Id проекта: %d\n", p.Id)
	fmt.Printf("Имя проекта: %s\n", p.Name)
	fmt.Printf("Задачи проекта:\n")

	for _, val := range p.Tasks {
		fmt.Printf("Id задачи: %d\n", val.Id)
		fmt.Printf("Заголовок задачи: %s\n", val.Title)
		fmt.Printf("Описание задачи: %s\n", val.Description)
		fmt.Printf("Статус готовности задачи: %t\n", val.Status)
	}
}
