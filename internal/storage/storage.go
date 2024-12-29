package storage

import "github.com/iharshyadav/backend/internal/types"

type Storage interface {
	CreateUserInterface(name string, email string, age int) (int64, error)
	GetUserById(id int64) (types.CreateUser,error)
	GetUsers() ([]types.CreateUser,error)
}