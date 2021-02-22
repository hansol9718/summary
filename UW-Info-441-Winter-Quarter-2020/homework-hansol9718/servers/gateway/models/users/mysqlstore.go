package users

import (
	"database/sql"
	"fmt"
)

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore(db *sql.DB) *MySQLStore {
	if db == nil {
		panic("nil database pointer passed to NewSqlStore")
	}
	return &MySQLStore{
		db: db,
	}
}

const GetID = "select * from users where id=?"
const GetEmail = "select * from users where email=?"
const GetUserName = "select * from users where username=?"
const SQLInsert = "insert into users(email, passHash, username, firstname, lastname, photoUrl) values (?,?,?,?,?,?)"
const SQLUpdate = "update users set firstname=?, lastname=? where id=?"
const SQLDelete = "delete from users where id=?"

func (s *MySQLStore) GetByID(id int64) (*User, error) {
	u := &User{}
	err := s.db.QueryRow(GetID, id).Scan(&u.ID, &u.Email, &u.PassHash, &u.UserName,
		&u.FirstName, &u.LastName, &u.PhotoURL)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *MySQLStore) GetByEmail(email string) (*User, error) {
	u := &User{}
	err := s.db.QueryRow(GetEmail, email).Scan(&u.ID, &u.Email, &u.PassHash, &u.UserName,
		&u.FirstName, &u.LastName, &u.PhotoURL)
	if err != nil {
		return nil, err
	}
	return u, nil
}


func (s *MySQLStore) GetByUserName(username string) (*User, error) {
	u := &User{}
	err := s.db.QueryRow(GetUserName, username).Scan(&u.ID, &u.Email, &u.PassHash, &u.UserName,
		&u.FirstName, &u.LastName, &u.PhotoURL)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *MySQLStore) Insert(u *User) (*User, error) {
	result, err := s.db.Exec(SQLInsert, u.Email, u.PassHash, u.UserName, u.FirstName, u.LastName, u.PhotoURL)
	if err != nil {
		return nil, err
	}
    // gets the new ID from the result that the database returns, 
    // usually returns a whole record but depends on the specific DB 
	newID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = newID
	return u, nil
}

func (s *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	_, err := s.db.Exec(SQLUpdate, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}
func (s *MySQLStore) Delete(id int64) error {
	_, err := s.db.Exec(SQLDelete, id)
	if err != nil {
		return fmt.Errorf("error deleteing from database %v", err)
	}
	return nil
}


