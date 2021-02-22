package users

import (
	"database/sql"
	"regexp"
	"testing"
	"reflect"
	"fmt"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func createRows(user *User) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "email", "passhash", "username", "firstname", "lastname", "photourl"})
	rows.AddRow(user.ID, user.Email, user.PassHash, user.UserName,
				user.FirstName, user.LastName, user.PhotoURL)
	return rows
}


func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}
	defer db.Close()

	expectedUser := &User{
		ID: 1,
		Email: "test@uw.edu",
		UserName: "hansol7",
		FirstName: "Hansol",
		LastName:  "Kim",
	}
	
	row := createRows(expectedUser)
	mock.ExpectQuery(regexp.QuoteMeta(GetID)).WithArgs(expectedUser.ID).WillReturnRows(row)
	store := NewMySQLStore(db)

	user, err := store.GetByID(expectedUser.ID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if err == nil && !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("user returned does not match expected user")
	}

	mock.ExpectQuery(regexp.QuoteMeta(GetID)).WithArgs(2).WillReturnError(sql.ErrNoRows)
	_, err = store.GetByID(2)

	if err == nil {
		t.Errorf("expected error: %v", sql.ErrNoRows)
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("error with sql mock expectation : %v", err)
	}
}

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}
	defer db.Close()

	expectedUser := &User{
		ID: 1,
		Email: "test@uw.edu",
		UserName: "hansol7",
		FirstName: "Hansol",
		LastName:  "Kim",
	}

	row := createRows(expectedUser)
	mock.ExpectQuery(regexp.QuoteMeta(GetEmail)).WithArgs(expectedUser.Email).WillReturnRows(row)
	store := NewMySQLStore(db)
	user, err := store.GetByEmail(expectedUser.Email)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if err == nil && !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("user returned does not match expected user")
	}

	mock.ExpectQuery(regexp.QuoteMeta(GetEmail)).WithArgs("idk@uw.edu").WillReturnError(sql.ErrNoRows)
	_, err = store.GetByEmail("idk@uw.edu")

	if err == nil {
		t.Errorf("expected error: %v", sql.ErrNoRows)
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("error with sql mock expectation : %v", err)
	}
}

func TestGetByUserName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}
	defer db.Close()

	expectedUser := &User{
		ID: 1,
		Email: "test@uw.edu",
		UserName: "hansol7",
		FirstName: "Hansol",
		LastName:  "Kim",
	}

	row := createRows(expectedUser)
	mock.ExpectQuery(regexp.QuoteMeta(GetUserName)).WithArgs(expectedUser.UserName).WillReturnRows(row)
	store := NewMySQLStore(db)
	user, err := store.GetByUserName(expectedUser.UserName)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if err == nil && !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("user returned does not match expected user")
	}
	mock.ExpectQuery(regexp.QuoteMeta(GetUserName)).WithArgs("ughugh").WillReturnError(sql.ErrNoRows)
	_, err = store.GetByUserName("ughugh")

	if err == nil {
		t.Errorf("expected error: %v", sql.ErrNoRows)
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("error with sql mock expectation : %v", err)
	}
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}
	defer db.Close()

	nu := &NewUser{
		Email:        "test@uw.edu",
		UserName:     "hansol7",
		Password:     "test123",
		PasswordConf: "test123",
		FirstName:    "Hansol",
		LastName:     "Kim",
	}
	user, err := nu.ToUser()
	if err != nil {
		t.Errorf("error adding user: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(SQLInsert)).WithArgs(user.Email, user.PassHash,
		user.UserName, user.FirstName, user.LastName, user.PhotoURL).WillReturnResult(sqlmock.NewResult(1,1))
	store := NewMySQLStore(db)
	returned, err := store.Insert(user)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if err == nil && !reflect.DeepEqual(returned, user) {
		t.Errorf("Returned user not equal to expected user")
	}

	insertErr := fmt.Errorf("Error with insert")
	invalidUser := &User{
		ID:        5,
		Email:     "testlala@uw.edu",
		PassHash:  nil,
		UserName:  "nope",
		FirstName: "bye",
		LastName:  "lala",
		PhotoURL:  "what",
	}
	mock.ExpectExec(regexp.QuoteMeta(SQLInsert)).WithArgs(invalidUser.Email, invalidUser.PassHash,
		invalidUser.UserName, invalidUser.FirstName, invalidUser.LastName,
		invalidUser.PhotoURL).WillReturnError(insertErr)
	_, err = store.Insert(invalidUser)
	if err == nil {
		t.Errorf("Expected error: %v", insertErr)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}
	defer db.Close()
	expectedUser := &User{
		ID: 1,
		FirstName: "Hansol",
		LastName: "Kim",
	}
	store := NewMySQLStore(db)
	update := &Updates {
		"Caleb",
		"Trapp",
	}
	mock.ExpectExec(regexp.QuoteMeta(SQLUpdate)).WithArgs(update.FirstName, update.LastName, expectedUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	row := createRows(expectedUser)
	mock.ExpectQuery(regexp.QuoteMeta(GetID)).WithArgs(expectedUser.ID).WillReturnRows(row)

	updated, err := store.Update(expectedUser.ID, update)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if err == nil && !reflect.DeepEqual(expectedUser, updated) {
		t.Errorf("User returned does not match expected user")
	}
	updateErr := fmt.Errorf("error updating: %v", err)
	mock.ExpectExec(regexp.QuoteMeta(SQLUpdate)).WithArgs(update.FirstName, update.LastName, 2).WillReturnError(updateErr)

	_, err = store.Update(2, update)
	if err == nil {
		t.Errorf("Expected error: %v", updateErr)
	}

	invalidUpdate := &Updates{}
	mock.ExpectExec(regexp.QuoteMeta(SQLUpdate)).WithArgs(invalidUpdate.FirstName, invalidUpdate.LastName, 1).WillReturnError(updateErr)
	_, err = store.Update(expectedUser.ID, invalidUpdate)
	if err == nil {
		t.Errorf("Expected error: %v", updateErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}
	defer db.Close()

	expectedUser := &User{
		ID: 1,
		Email: "test@uw.edu",
		UserName: "hansol7",
		FirstName: "Hansol",
		LastName:  "Kim",
	}

	store := NewMySQLStore(db)
	mock.ExpectExec(regexp.QuoteMeta(SQLDelete)).
	WithArgs(expectedUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	err = store.Delete(expectedUser.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	deleteErr := fmt.Errorf("error deleting")
	err = store.Delete(5)
	if err == nil {
		t.Errorf("Expected error: %v", deleteErr)
	}
}