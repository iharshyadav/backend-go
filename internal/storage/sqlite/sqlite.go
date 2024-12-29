package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/iharshyadav/backend/internal/config"
	"github.com/iharshyadav/backend/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite , error) {
	db , err := sql.Open("sqlite3",cfg.StoragePath)

	if err != nil {
		return nil , err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
	)`)

	if err != nil {
		return nil, err
	}
	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateUserInterface (name string , email string , age int) (int64 , error) {

	stmt , err := s.Db.Prepare("INSERT INTO user (name , email , age) VALUES (? , ? , ?)")

	if err != nil {
		return 0,err
	}

	defer stmt.Close()

	result , err := stmt.Exec(name , email , age)

	if err != nil {
		return 0 ,err
	}

	lastId , err := result.LastInsertId()

	if err != nil {
		return 0,err
	}

	return lastId , nil
}

func (s *Sqlite) GetUserById (id int64) (types.CreateUser, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM user WHERE id = ? LIMIT 1")
	if err != nil {
		return types.CreateUser{}, err
	}
	defer stmt.Close()

	var user types.CreateUser

	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name, &user.Email, &user.Age)

	if err != nil {

		if err == sql.ErrNoRows {
			return types.CreateUser{}, fmt.Errorf("no user found with id %s", fmt.Sprint(id))
		}
		return types.CreateUser{}, fmt.Errorf("query error: %w", err)
	}

	return user, nil
}

func (s *Sqlite) GetUsers() ([]types.CreateUser, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM user")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.CreateUser
	
	for rows.Next() {
		var user types.CreateUser
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Age)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}