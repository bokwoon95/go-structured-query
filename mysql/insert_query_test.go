package sq

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/matryer/is"
)

func TestInsertQuery_ToSQL(t *testing.T) {
	type TT struct {
		description string
		q           InsertQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS().As("u")
	tests := []TT{
		{"empty", InsertQuery{}, "INSERT INTO NULL", nil},
		{"Into", WithDefaultLog(Linterpolate).InsertInto(u), "INSERT INTO devlab.users", nil},
		{
			"Insert Values",
			WithDefaultLog(Lverbose).
				InsertInto(u).
				Columns(u.DISPLAYNAME, u.EMAIL).
				Values("aaa", "aaa@email.com").
				Values("bbb", "bbb@email.com").
				OnDuplicateKeyUpdate(
					u.DISPLAYNAME.Set(u.DISPLAYNAME),
					u.EMAIL.Set(u.EMAIL),
				),
			"INSERT INTO devlab.users (displayname, email)" +
				" VALUES (?, ?), (?, ?)" +
				" ON DUPLICATE KEY UPDATE" +
				" displayname = displayname, email = email",
			[]interface{}{"aaa", "aaa@email.com", "bbb", "bbb@email.com"},
		},
		{
			"Insert InsertRow",
			WithDefaultLog(Linterpolate).
				InsertInto(u).
				InsertRow(
					u.DISPLAYNAME.SetString("aaa"),
					u.EMAIL.SetString("aaa@email.com"),
				).
				InsertRow(
					u.EMAIL.SetString("bbb"),
					u.EMAIL.SetString("bbb@email.com"),
				).
				OnDuplicateKeyUpdate(
					u.DISPLAYNAME.Set(Values(u.DISPLAYNAME)),
					u.EMAIL.Set(Values(u.EMAIL)),
				),
			"INSERT INTO devlab.users (displayname, email)" +
				" VALUES (?, ?), (?, ?)" +
				" ON DUPLICATE KEY UPDATE" +
				" displayname = VALUES(displayname), email = VALUES(email)",
			[]interface{}{"aaa", "aaa@email.com", "bbb", "bbb@email.com"},
		},
		{
			"Insert Select",
			WithLog(customLogger, 0).
				InsertInto(u).
				Columns(u.DISPLAYNAME, u.EMAIL).
				Select(
					Select(u.DISPLAYNAME, u.EMAIL).
						From(u).
						Where(u.USER_ID.In([]int{1, 2, 3})),
				).
				OnDuplicateKeyUpdate(
					u.DISPLAYNAME.Set(Values(u.DISPLAYNAME)),
					u.EMAIL.Set(Values(u.EMAIL)),
				),
			"INSERT INTO devlab.users (displayname, email)" +
				" SELECT u.displayname, u.email FROM devlab.users AS u" +
				" WHERE u.user_id IN (?, ?, ?)" +
				" ON DUPLICATE KEY UPDATE" +
				" displayname = VALUES(displayname), email = VALUES(email)",
			[]interface{}{1, 2, 3},
		},
		func() TT {
			desc := "Insert Ignore"
			u := USERS()
			q := WithDefaultLog(0).
				InsertIgnoreInto(u).
				Columns(u.DISPLAYNAME, u.EMAIL).
				Values("aaa", "aaa@email.com").
				Values("bbb", "bbb@email.com")
			wantQuery := "INSERT IGNORE INTO devlab.users (displayname, email) VALUES (?, ?), (?, ?)"
			wantArgs := []interface{}{"aaa", "aaa@email.com", "bbb", "bbb@email.com"}
			return TT{desc, q, wantQuery, wantArgs}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			var _ Query = tt.q
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestInsertQuery_Exec(t *testing.T) {
	if testing.Short() {
		return
	}
	is := is.New(t)
	db, err := sql.Open("txdb", "DeleteQuery_Exec")
	is.NoErr(err)
	defer db.Close()
	u := USERS()

	// Missing DB
	_, _, err = InsertInto(u).
		Columns(u.DISPLAYNAME, u.EMAIL).
		Values("aaa", "aaa@email.com").
		Exec(nil, ElastInsertID|ErowsAffected)
	is.True(err != nil)

	// SQL syntax error
	// use a tempDB so we don't foul up the current db transaction with the error
	tempDB, err := sql.Open("txdb", randomString(8))
	is.NoErr(err)
	_, _, err = WithLog(customLogger, Linterpolate).
		WithDB(tempDB).
		InsertInto(u).
		Columns(u.DISPLAYNAME, u.EMAIL).
		Values(Fieldf("ERROR")).
		Exec(nil, ElastInsertID|ErowsAffected)
	is.True(err != nil)
	tempDB.Close()

	// simulate timeout
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	_, _, err = WithDB(db).
		InsertInto(u).
		Columns(u.DISPLAYNAME, u.EMAIL).
		Values("aaa", "aaa@email.com").
		ExecContext(ctx, nil, ElastInsertID|ErowsAffected)
	is.True(errors.Is(err, context.DeadlineExceeded))

	// rowsAffected
	lastInsertID, rowsAffected, err := WithDefaultLog(Lverbose).
		WithDB(db).
		InsertInto(u).
		Columns(u.DISPLAYNAME, u.EMAIL).
		Values("aaa", "aaa@email.com").
		Exec(nil, ElastInsertID|ErowsAffected)
	is.NoErr(err)
	is.Equal(int64(1), rowsAffected)
	var id int64
	err = From(u).
		Where(u.EMAIL.EqString("aaa@email.com")).
		SelectRowx(func(row *Row) { id = row.Int64(u.USER_ID) }).
		Fetch(db)
	is.NoErr(err)
	is.Equal(id, lastInsertID)

	// again
	_, rowsAffected, err = WithLog(customLogger, Lverbose).
		WithDB(db).
		InsertInto(u).
		Columns(u.DISPLAYNAME, u.EMAIL).
		Values("aaa", "aaa@email.com").
		OnDuplicateKeyUpdate(
			u.DISPLAYNAME.Set(u.DISPLAYNAME),
			u.EMAIL.Set(u.EMAIL),
		).
		Exec(nil, ElastInsertID|ErowsAffected)
	is.NoErr(err)
	is.Equal(int64(0), rowsAffected)
}

func TestInsertQuery_Basic(t *testing.T) {
	is := is.New(t)
	var q InsertQuery
	q = InsertIgnoreInto(nil)
	is.Equal(true, q.Ignore)
	is.Equal(nil, q.IntoTable)
}
