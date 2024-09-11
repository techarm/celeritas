package data

import (
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	db2 "github.com/upper/db/v4"
)

func TestNew(t *testing.T) {
	fakeDB, mock, _ := sqlmock.New()
	defer fakeDB.Close()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("test"))
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("test"))

	_ = os.Setenv("DATATYPE_TYPE", "postgres")
	m := New(fakeDB)
	if fmt.Sprintf("%T", m) != "data.Models" {
		t.Error("Wrong type", fmt.Sprintf("%T", m))
	}

	_ = os.Setenv("DATATYPE_TYPE", "mysql")
	m = New(fakeDB)
	if fmt.Sprintf("%T", m) != "data.Models" {
		t.Error("Wrong type", fmt.Sprintf("%T", m))
	}
}

func TestGetInsertID(t *testing.T) {
	var id db2.ID
	id = int64(1)

	returnedID := getInsertID(id)
	if fmt.Sprintf("%T", returnedID) != "int" {
		t.Error("Wrong type", fmt.Sprintf("%T", returnedID))
	}

	id = 1
	returnedID = getInsertID(id)
	if fmt.Sprintf("%T", returnedID) != "int" {
		t.Error("Wrong type", fmt.Sprintf("%T", returnedID))
	}
}
