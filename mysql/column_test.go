package sq

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestColumnInsert(t *testing.T) {
	is := is.New(t)
	type User struct {
		UserID      int
		DisplayName string
		Email       string
		Password    string
	}
	users := []User{
		{
			UserID:      1,
			DisplayName: "one",
			Email:       "one",
			Password:    "one",
		},
		{
			UserID:      2,
			DisplayName: "two",
			Email:       "two",
			Password:    "two",
		},
		{
			UserID:      3,
			DisplayName: "three",
			Email:       "three",
			Password:    "three",
		},
	}
	col := &Column{mode: colmodeInsert}
	u := USERS()
	for _, user := range users {
		col.Set(u.USER_ID, user.UserID)
		col.Set(u.DISPLAYNAME, user.DisplayName)
		col.Set(u.EMAIL, user.Email)
		col.Set(u.PASSWORD, user.Password)
	}
	is.Equal(Fields{u.USER_ID, u.DISPLAYNAME, u.EMAIL, u.PASSWORD}, col.insertColumns)
	is.Equal(
		RowValues{
			{users[0].UserID, users[0].DisplayName, users[0].Email, users[0].Password},
			{users[1].UserID, users[1].DisplayName, users[1].Email, users[1].Password},
			{users[2].UserID, users[2].DisplayName, users[2].Email, users[2].Password},
		},
		col.rowValues,
	)
}

func TestColumnUpdate(t *testing.T) {
	is := is.New(t)
	type User struct {
		UserID      int
		DisplayName string
		Email       string
		Password    string
	}
	col := &Column{mode: colmodeUpdate}
	u := USERS()
	user := User{
		UserID:      1,
		DisplayName: "one",
		Email:       "one",
		Password:    "one",
	}
	col.Set(u.USER_ID, user.UserID)
	col.Set(u.DISPLAYNAME, user.DisplayName)
	col.Set(u.EMAIL, user.Email)
	col.Set(u.PASSWORD, user.Password)
	is.Equal(
		Assignments{
			u.USER_ID.Set(user.UserID),
			u.DISPLAYNAME.Set(user.DisplayName),
			u.EMAIL.Set(user.Email),
			u.PASSWORD.Set(user.Password),
		},
		col.assignments,
	)
}

func TestColumn_Basic(t *testing.T) {
	is := is.New(t)
	now := time.Now()
	a := APPLICATIONS().As("a")
	col := &Column{mode: colmodeInsert}
	col.SetBool(a.SUBMITTED, true)
	col.SetFloat64(a.TEAM_ID, 3.0)
	col.SetInt(a.APPLICATION_ID, 2)
	col.SetInt64(a.APPLICATION_FORM_ID, 4)
	col.SetTime(a.CREATED_AT, now)
	is.Equal(
		Fields{a.SUBMITTED, a.TEAM_ID, a.APPLICATION_ID, a.APPLICATION_FORM_ID, a.CREATED_AT},
		col.insertColumns,
	)
	is.Equal(
		RowValues{{true, 3.0, 2, int64(4), now}},
		col.rowValues,
	)
}
