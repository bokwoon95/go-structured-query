package sq

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/matryer/is"
)

func TestMain(t *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_PORT := os.Getenv("POSTGRES_PORT")
	POSTGRES_NAME := os.Getenv("POSTGRES_NAME")
	txdb.Register("txdb", "postgres", fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_PORT, POSTGRES_NAME))

	os.Exit(t.Run())
}

type Option struct {
	Value   string `json:"Value"`
	Display string `json:"Display"`
}

type Subquestion struct {
	Name string `json:"Name"`
	Text string `json:"Text"`
}

type Question struct {
	Type         string        `json:"Type"`
	Text         string        `json:"Text"`
	Name         string        `json:"Name"`
	Options      []Option      `json:"Options"`
	Subquestions []Subquestion `json:"Subquestions"`
}

func (question Question) Value() (driver.Value, error) {
	b, err := json.Marshal(question)
	return driver.Value(string(b)), err
}

func (question *Question) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case nil:
		return nil
	case string:
		err = json.Unmarshal([]byte(v), &question)
		return err
	case []byte:
		err = json.Unmarshal(v, &question)
		return err
	default:
		return fmt.Errorf("value %#v from database is neither a string nor NULL", value)
	}
}

type Questions []Question

func (questions Questions) Value() (driver.Value, error) {
	b, err := json.Marshal(questions)
	return driver.Value(string(b)), err
}

func (questions *Questions) Scan(value interface{}) error {
	var err error
	switch v := value.(type) {
	case nil:
		return nil
	case string:
		err = json.Unmarshal([]byte(v), questions)
	case []byte:
		err = json.Unmarshal(v, questions)
	default:
		return fmt.Errorf("value %#v from database is neither a string nor NULL", value)
	}
	return err
}

func TestRow_ScanInto(t *testing.T) {
	if testing.Short() {
		return
	}
	is := is.New(t)
	db, err := sql.Open("txdb", "Row_ScanInto")
	is.NoErr(err)
	a := APPLICATIONS()
	type Data struct {
		wantBool        bool
		wantNullBool    sql.NullBool
		wantFloat64     float64
		wantNullFloat64 sql.NullFloat64
		wantInt32       int32
		wantNullInt32   sql.NullInt32
		wantInt64       int64
		wantNullInt64   sql.NullInt64
		wantInt         int
		wantString      string
		wantNullString  sql.NullString

		gotBool        bool
		gotNullBool    sql.NullBool
		gotFloat64     float64
		gotNullFloat64 sql.NullFloat64
		gotInt32       int32
		gotNullInt32   sql.NullInt32
		gotInt64       int64
		gotInt         int
		gotNullInt64   sql.NullInt64
		gotString      string
		gotNullString  sql.NullString
		gotTime        time.Time
		gotNullTime    sql.NullTime
	}
	data := Data{
		wantBool:        true,                                          // applications.submitted
		wantNullBool:    sql.NullBool{Valid: true, Bool: true},         // applications.submitted
		wantFloat64:     1,                                             // applications.application_id
		wantNullFloat64: sql.NullFloat64{Valid: true, Float64: 1},      // applications.application_id
		wantInt32:       1,                                             // applications.application_id
		wantNullInt32:   sql.NullInt32{Valid: true, Int32: 1},          // applications.application_id
		wantInt64:       1,                                             // applications.application_id
		wantInt:         1,                                             // applications.application_id
		wantNullInt64:   sql.NullInt64{Valid: true, Int64: 1},          // applications.application_id
		wantString:      "gemini",                                      // applications.project_level
		wantNullString:  sql.NullString{Valid: true, String: "gemini"}, // applications.project_level
	}
	err = WithDefaultLog(Lverbose).
		From(a).
		Where(a.APPLICATION_ID.EqInt(1)).
		SelectRowx(func(row *Row) {
			row.ScanInto(&data.gotBool, a.SUBMITTED)
			row.ScanInto(&data.gotNullBool, a.SUBMITTED)
			row.ScanInto(&data.gotFloat64, a.APPLICATION_ID)
			row.ScanInto(&data.gotNullFloat64, a.APPLICATION_ID)
			row.ScanInto(&data.gotInt32, a.APPLICATION_ID)
			row.ScanInto(&data.gotNullInt32, a.APPLICATION_ID)
			row.ScanInto(&data.gotInt64, a.APPLICATION_ID)
			row.ScanInto(&data.gotInt, a.APPLICATION_ID)
			row.ScanInto(&data.gotNullInt64, a.APPLICATION_ID)
			row.ScanInto(&data.gotString, a.PROJECT_LEVEL)
			row.ScanInto(&data.gotNullString, a.PROJECT_LEVEL)
			row.ScanInto(&data.gotTime, a.CREATED_AT)
			row.ScanInto(&data.gotNullTime, a.CREATED_AT)
		}).
		Fetch(db)
	is.NoErr(err)
	is.Equal(data.wantBool, data.gotBool)
	is.Equal(data.wantNullBool, data.gotNullBool)
	is.Equal(data.wantFloat64, data.gotFloat64)
	is.Equal(data.wantNullFloat64, data.gotNullFloat64)
	is.Equal(data.wantInt32, data.gotInt32)
	is.Equal(data.wantNullInt32, data.gotNullInt32)
	is.Equal(data.wantInt64, data.gotInt64)
	is.Equal(data.wantInt, data.gotInt)
	is.Equal(data.wantNullInt64, data.gotNullInt64)
	is.Equal(data.wantString, data.gotString)
	is.Equal(data.wantNullString, data.gotNullString)
	is.True(data.gotNullTime.Valid)

	var gotQuestions Questions
	var wantQuestions = Questions{
		Question{
			Type:         "long text",
			Name:         "readme",
			Text:         "<b>Project Readme</b>",
			Options:      []Option{},
			Subquestions: []Subquestion{},
		},
		Question{
			Type:         "long text",
			Name:         "log",
			Text:         "<b>Project Log</b>",
			Options:      []Option{},
			Subquestions: []Subquestion{},
		},
		Question{
			Type:         "short text",
			Name:         "poster",
			Text:         "<b>Poster Link</b>",
			Options:      nil,
			Subquestions: nil,
		},
		Question{
			Type:         "short text",
			Name:         "video",
			Text:         "<b>Video Link</b>",
			Options:      nil,
			Subquestions: nil,
		},
	}

	// Custom driver.Valuer struct
	f := FORMS()
	err = WithDefaultLog(Lverbose).
		From(f).
		Where(f.FORM_ID.EqInt(3)).
		SelectRowx(func(row *Row) {
			row.ScanInto(&gotQuestions, f.QUESTIONS)
		}).
		Fetch(db)
	is.NoErr(err)
	is.Equal(wantQuestions, gotQuestions)

	// Intentionally scan into wrong struct
	err = WithDefaultLog(0).
		From(f).
		Where(f.FORM_ID.EqInt(3)).
		SelectRowx(func(row *Row) {
			row.ScanInto(&struct{}{}, f.QUESTIONS)
		}).
		Fetch(db)
	is.True(err != nil)
}

func TestRow_ScanArray(t *testing.T) {
	if testing.Short() {
		return
	}
	is := is.New(t)
	db, err := sql.Open("txdb", "Row_ScanArray")
	is.NoErr(err)
	ur := USER_ROLES()
	var userID int
	var userRoles []string
	var userRoleIDs []int64
	err = WithDefaultLog(Lverbose).
		From(ur).
		Where(ur.USER_ID.EqInt(16)).
		GroupBy(ur.USER_ID).
		SelectRowx(func(row *Row) {
			userID = row.Int(ur.USER_ID)
			row.ScanArray(&userRoleIDs, Fieldf("array_agg(? ORDER BY ?)", ur.USER_ROLE_ID, ur.USER_ROLE_ID.Asc()))
			row.ScanArray(&userRoles, Fieldf("array_agg(? ORDER BY ?)", ur.ROLE, ur.ROLE.Asc()))
		}).
		Fetch(db)
	is.NoErr(err)
	is.Equal(userID, 16)
	is.Equal(userRoleIDs, []int64{16, 18})
	is.Equal(userRoles, []string{"applicant", "student"})

	// Intentionally scan into unsupported slice type
	var roles []int
	err = WithDefaultLog(Lverbose).
		From(ur).
		Where(ur.USER_ID.EqInt(16)).
		GroupBy(ur.USER_ID).
		SelectRowx(func(row *Row) {
			row.ScanArray(&roles, Fieldf("array_agg(? ORDER BY ?)", ur.USER_ROLE_ID, ur.USER_ROLE_ID.Asc()))
		}).
		Fetch(db)
	is.True(err != nil)
}

func TestRow_ScanJSON(t *testing.T) {
	if testing.Short() {
		return
	}
	is := is.New(t)
	db, err := sql.Open("txdb", "Row_ScanJSON")
	is.NoErr(err)

	ur := USER_ROLES()

	type User struct {
		UserID int
		Cohort string
	}

	var users []User

	field := Fieldf("jsonb_agg(jsonb_build_object(?))", []interface{}{
		Literal(`'UserID'`), ur.USER_ID,
		Literal(`'Cohort'`), ur.COHORT,
	})

	err = WithDefaultLog(Lverbose).
		SelectRowx(func(r *Row) {
			r.ScanJSON(&users, field)
		}).
		From(ur).
		Where(ur.USER_ID.EqInt(16)).
		GroupBy(ur.USER_ID).
		Fetch(db)

	is.NoErr(err)
	is.True(len(users) > 0)

	var user User

	err = WithDefaultLog(Lverbose).
		SelectRowx(func(r *Row) {
			r.ScanJSON(&user, field)
		}).
		From(ur).
		Where(ur.USER_ID.EqInt(16)).
		GroupBy(ur.USER_ID).
		Fetch(db)

	is.True(err != nil)
}

func TestRow_Assorted(t *testing.T) {
	if testing.Short() {
		return
	}
	is := is.New(t)
	db, err := sql.Open("txdb", "Row_Assorted")
	is.NoErr(err)
	a := APPLICATIONS()
	type Data struct {
		wantBool         bool
		wantBoolValid    bool
		wantFloat64      float64
		wantFloat64Valid bool
		wantInt64        int64
		wantInt64Valid   bool
		wantInt          int
		wantIntValid     bool
		wantString       string
		wantStringValid  bool

		gotBool         bool
		gotBoolValid    bool
		gotFloat64      float64
		gotFloat64Valid bool
		gotInt64        int64
		gotInt64Valid   bool
		gotInt          int
		gotIntValid     bool
		gotString       string
		gotStringValid  bool
		gotTime         time.Time
		gotTimeValid    bool
	}
	data := Data{
		wantBool:         true,
		wantBoolValid:    true,
		wantFloat64:      1,
		wantFloat64Valid: true,
		wantInt64:        1,
		wantInt64Valid:   true,
		wantInt:          1,
		wantIntValid:     true,
		wantString:       "gemini",
		wantStringValid:  true,
	}
	err = WithDefaultLog(Lverbose).
		From(a).
		Where(a.APPLICATION_ID.EqInt(1)).
		SelectRowx(func(row *Row) {
			data.gotBool = row.Bool(a.SUBMITTED)
			data.gotBoolValid = row.BoolValid(a.SUBMITTED)
			data.gotFloat64 = row.Float64(a.APPLICATION_ID)
			data.gotFloat64Valid = row.Float64Valid(a.APPLICATION_ID)
			data.gotInt64 = row.Int64(a.APPLICATION_ID)
			data.gotInt64Valid = row.Int64Valid(a.APPLICATION_ID)
			data.gotInt = row.Int(a.APPLICATION_ID)
			data.gotIntValid = row.IntValid(a.APPLICATION_ID)
			data.gotString = row.String(a.PROJECT_LEVEL)
			data.gotStringValid = row.StringValid(a.PROJECT_LEVEL)
			data.gotTime = row.Time(a.CREATED_AT)
			data.gotTimeValid = row.TimeValid(a.CREATED_AT)
		}).
		Fetch(db)
	is.NoErr(err)
	is.Equal(data.wantBool, data.gotBool)
	is.Equal(data.wantBoolValid, data.gotBoolValid)
	is.Equal(data.wantFloat64, data.gotFloat64)
	is.Equal(data.wantFloat64Valid, data.gotFloat64Valid)
	is.Equal(data.wantInt64, data.gotInt64)
	is.Equal(data.wantInt64Valid, data.gotInt64Valid)
	is.Equal(data.wantInt, data.gotInt)
	is.Equal(data.wantIntValid, data.gotIntValid)
	is.Equal(data.wantString, data.gotString)
	is.Equal(data.wantStringValid, data.gotStringValid)
	is.True(data.gotTimeValid)
}

func TestRowFunctions(t *testing.T) {
	if testing.Short() {
		return
	}
	is := is.New(t)
	db, err := sql.Open("txdb", "RowFunctions")
	is.NoErr(err)
	a := APPLICATIONS()
	type Data struct {
		wantBool         bool
		wantBoolValid    bool
		wantNullBool     sql.NullBool
		wantFloat64      float64
		wantFloat64Valid bool
		wantNullFloat64  sql.NullFloat64
		wantInt64        int64
		wantInt64Valid   bool
		wantNullInt64    sql.NullInt64
		wantInt          int
		wantIntValid     bool
		wantString       string
		wantStringValid  bool
		wantNullString   sql.NullString

		gotBool         bool
		gotBoolValid    bool
		gotNullBool     sql.NullBool
		gotFloat64      float64
		gotFloat64Valid bool
		gotNullFloat64  sql.NullFloat64
		gotInt64        int64
		gotInt64Valid   bool
		gotNullInt64    sql.NullInt64
		gotInt          int
		gotIntValid     bool
		gotString       string
		gotStringValid  bool
		gotNullString   sql.NullString
		gotTime         time.Time
		gotTimeValid    bool
		gotNullTime     sql.NullTime
	}
	data := Data{
		wantBool:         true,
		wantBoolValid:    true,
		wantNullBool:     sql.NullBool{Bool: true, Valid: true},
		wantFloat64:      1,
		wantFloat64Valid: true,
		wantNullFloat64:  sql.NullFloat64{Float64: 1, Valid: true},
		wantInt64:        1,
		wantInt64Valid:   true,
		wantNullInt64:    sql.NullInt64{Int64: 1, Valid: true},
		wantInt:          1,
		wantIntValid:     true,
		wantString:       "gemini",
		wantStringValid:  true,
		wantNullString:   sql.NullString{String: "gemini", Valid: true},
	}
	err = WithDefaultLog(Lverbose).
		From(a).
		Where(a.APPLICATION_ID.EqInt(1)).
		SelectRowx(func(row *Row) {
			data.gotBool = row.Bool(a.SUBMITTED)
			data.gotBoolValid = row.BoolValid(a.SUBMITTED)
			data.gotNullBool = row.NullBool(a.SUBMITTED)
			data.gotFloat64 = rowFloat64(row, a.APPLICATION_ID)
			data.gotFloat64Valid = rowFloat64Valid(row, a.APPLICATION_ID)
			data.gotNullFloat64 = row.NullFloat64(a.APPLICATION_ID)
			data.gotInt64 = rowInt64(row, a.APPLICATION_ID)
			data.gotInt64Valid = rowInt64Valid(row, a.APPLICATION_ID)
			data.gotNullInt64 = row.NullInt64(a.APPLICATION_ID)
			data.gotInt = rowInt(row, a.APPLICATION_ID)
			data.gotIntValid = rowIntValid(row, a.APPLICATION_ID)
			data.gotString = rowString(row, a.PROJECT_LEVEL)
			data.gotStringValid = rowStringValid(row, a.PROJECT_LEVEL)
			data.gotNullString = row.NullString(a.PROJECT_LEVEL)
			data.gotTime = rowTime(row, a.CREATED_AT)
			data.gotTimeValid = rowTimeValid(row, a.CREATED_AT)
			data.gotNullTime = row.NullTime(a.CREATED_AT)
		}).
		Fetch(db)
	is.NoErr(err)
	is.Equal(data.wantBool, data.gotBool)
	is.Equal(data.wantBoolValid, data.gotBoolValid)
	is.Equal(data.wantFloat64, data.gotFloat64)
	is.Equal(data.wantFloat64Valid, data.gotFloat64Valid)
	is.Equal(data.wantInt64, data.gotInt64)
	is.Equal(data.wantInt64Valid, data.gotInt64Valid)
	is.Equal(data.wantInt, data.gotInt)
	is.Equal(data.wantIntValid, data.gotIntValid)
	is.Equal(data.wantString, data.gotString)
	is.Equal(data.wantStringValid, data.gotStringValid)
	is.True(data.gotTimeValid)
}

func TestRowUUIDs(t *testing.T) {
	if testing.Short() {
		return
	}

	is := is.New(t)
	db, err := sql.Open("txdb", "RowUUIDs")
	is.NoErr(err)
	m := MEDIA()

	// hardcoded in init.sql
	wantUUID, err := uuid.Parse("dbc59321-22ee-4613-bb27-8b2ba5be3109")

	is.NoErr(err)

	type Data struct {
		wantUUID [16]byte
		gotUUID  [16]byte
	}

	data := Data{
		wantUUID: [16]byte(wantUUID),
	}

	err = WithDefaultLog(Lverbose).
		From(m).
		Limit(1).
		SelectRowx(func(r *Row) {
			data.gotUUID = r.UUID(m.UUID)
		}).
		Fetch(db)

	is.NoErr(err)
	is.Equal(data.gotUUID, data.wantUUID)
}
