package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
	"github.com/matryer/is"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MYSQL_USER := os.Getenv("MYSQL_USER")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	MYSQL_PORT := os.Getenv("MYSQL_PORT")
	MYSQL_NAME := os.Getenv("MYSQL_NAME")
	txdb.Register("txdb", "mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s?parseTime=true", MYSQL_USER, MYSQL_PASSWORD, MYSQL_PORT, MYSQL_NAME))
}

func TestBuildTables(t *testing.T) {
	if testing.Short() {
		return
	}

	db, err := sql.Open("txdb", "BuildTables")

	is := is.New(t)
	is.NoErr(err)

	config := Config{
		DB: db,
		Package: "tables",
		Schemas: []string{"devlab"},
		Exclude: nil,
		Logger: &mockLogger{},
	}

	var writer strings.Builder
	err = BuildTables(config, &writer)
	
	is.NoErr(err)

	out := writer.String()

	is.Equal(out, expectedTables)
}

const expectedTables = `// Code generated by 'sqgen-mysql tables'; DO NOT EDIT.
package tables

import (
	sq "github.com/bokwoon95/go-structured-query/mysql"
)

// TABLE_APPLICATIONS references the devlab.applications table.
type TABLE_APPLICATIONS struct {
	*sq.TableInfo
	APPLICATION_DATA     sq.JSONField
	APPLICATION_FORM_ID  sq.NumberField
	APPLICATION_ID       sq.NumberField
	COHORT               sq.StringField
	CREATED_AT           sq.TimeField
	CREATOR_USER_ROLE_ID sq.NumberField
	DELETED_AT           sq.TimeField
	MAGICSTRING          sq.StringField
	PROJECT_IDEA         sq.StringField
	PROJECT_LEVEL        sq.StringField
	STATUS               sq.StringField
	SUBMITTED            sq.BooleanField
	TEAM_ID              sq.NumberField
	TEAM_NAME            sq.StringField
	UPDATED_AT           sq.TimeField
}

// APPLICATIONS creates an instance of the devlab.applications table.
func APPLICATIONS() TABLE_APPLICATIONS {
	tbl := TABLE_APPLICATIONS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "applications",
	}}
	tbl.APPLICATION_DATA = sq.NewJSONField("application_data", tbl.TableInfo)
	tbl.APPLICATION_FORM_ID = sq.NewNumberField("application_form_id", tbl.TableInfo)
	tbl.APPLICATION_ID = sq.NewNumberField("application_id", tbl.TableInfo)
	tbl.COHORT = sq.NewStringField("cohort", tbl.TableInfo)
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.CREATOR_USER_ROLE_ID = sq.NewNumberField("creator_user_role_id", tbl.TableInfo)
	tbl.DELETED_AT = sq.NewTimeField("deleted_at", tbl.TableInfo)
	tbl.MAGICSTRING = sq.NewStringField("magicstring", tbl.TableInfo)
	tbl.PROJECT_IDEA = sq.NewStringField("project_idea", tbl.TableInfo)
	tbl.PROJECT_LEVEL = sq.NewStringField("project_level", tbl.TableInfo)
	tbl.STATUS = sq.NewStringField("status", tbl.TableInfo)
	tbl.SUBMITTED = sq.NewBooleanField("submitted", tbl.TableInfo)
	tbl.TEAM_ID = sq.NewNumberField("team_id", tbl.TableInfo)
	tbl.TEAM_NAME = sq.NewStringField("team_name", tbl.TableInfo)
	tbl.UPDATED_AT = sq.NewTimeField("updated_at", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_APPLICATIONS) As(alias string) TABLE_APPLICATIONS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_APPLICATIONS_STATUS_ENUM references the devlab.applications_status_enum table.
type TABLE_APPLICATIONS_STATUS_ENUM struct {
	*sq.TableInfo
	STATUS sq.StringField
}

// APPLICATIONS_STATUS_ENUM creates an instance of the devlab.applications_status_enum table.
func APPLICATIONS_STATUS_ENUM() TABLE_APPLICATIONS_STATUS_ENUM {
	tbl := TABLE_APPLICATIONS_STATUS_ENUM{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "applications_status_enum",
	}}
	tbl.STATUS = sq.NewStringField("status", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_APPLICATIONS_STATUS_ENUM) As(alias string) TABLE_APPLICATIONS_STATUS_ENUM {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_COHORT_ENUM references the devlab.cohort_enum table.
type TABLE_COHORT_ENUM struct {
	*sq.TableInfo
	COHORT sq.StringField
}

// COHORT_ENUM creates an instance of the devlab.cohort_enum table.
func COHORT_ENUM() TABLE_COHORT_ENUM {
	tbl := TABLE_COHORT_ENUM{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "cohort_enum",
	}}
	tbl.COHORT = sq.NewStringField("cohort", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_COHORT_ENUM) As(alias string) TABLE_COHORT_ENUM {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_FEEDBACK_ON_TEAMS references the devlab.feedback_on_teams table.
type TABLE_FEEDBACK_ON_TEAMS struct {
	*sq.TableInfo
	CREATED_AT          sq.TimeField
	DELETED_AT          sq.TimeField
	EVALUATEE_TEAM_ID   sq.NumberField
	EVALUATOR_TEAM_ID   sq.NumberField
	FEEDBACK_DATA       sq.JSONField
	FEEDBACK_FORM_ID    sq.NumberField
	FEEDBACK_ID_ON_TEAM sq.NumberField
	OVERRIDE_OPEN       sq.BooleanField
	SUBMITTED           sq.BooleanField
	UPDATED_AT          sq.TimeField
}

// FEEDBACK_ON_TEAMS creates an instance of the devlab.feedback_on_teams table.
func FEEDBACK_ON_TEAMS() TABLE_FEEDBACK_ON_TEAMS {
	tbl := TABLE_FEEDBACK_ON_TEAMS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "feedback_on_teams",
	}}
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.DELETED_AT = sq.NewTimeField("deleted_at", tbl.TableInfo)
	tbl.EVALUATEE_TEAM_ID = sq.NewNumberField("evaluatee_team_id", tbl.TableInfo)
	tbl.EVALUATOR_TEAM_ID = sq.NewNumberField("evaluator_team_id", tbl.TableInfo)
	tbl.FEEDBACK_DATA = sq.NewJSONField("feedback_data", tbl.TableInfo)
	tbl.FEEDBACK_FORM_ID = sq.NewNumberField("feedback_form_id", tbl.TableInfo)
	tbl.FEEDBACK_ID_ON_TEAM = sq.NewNumberField("feedback_id_on_team", tbl.TableInfo)
	tbl.OVERRIDE_OPEN = sq.NewBooleanField("override_open", tbl.TableInfo)
	tbl.SUBMITTED = sq.NewBooleanField("submitted", tbl.TableInfo)
	tbl.UPDATED_AT = sq.NewTimeField("updated_at", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_FEEDBACK_ON_TEAMS) As(alias string) TABLE_FEEDBACK_ON_TEAMS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_FEEDBACK_ON_USERS references the devlab.feedback_on_users table.
type TABLE_FEEDBACK_ON_USERS struct {
	*sq.TableInfo
	CREATED_AT             sq.TimeField
	DELETED_AT             sq.TimeField
	EVALUATEE_USER_ROLE_ID sq.NumberField
	EVALUATOR_TEAM_ID      sq.NumberField
	FEEDBACK_DATA          sq.JSONField
	FEEDBACK_FORM_ID       sq.NumberField
	FEEDBACK_ID_ON_USER    sq.NumberField
	OVERRIDE_OPEN          sq.BooleanField
	SUBMITTED              sq.BooleanField
	UPDATED_AT             sq.TimeField
}

// FEEDBACK_ON_USERS creates an instance of the devlab.feedback_on_users table.
func FEEDBACK_ON_USERS() TABLE_FEEDBACK_ON_USERS {
	tbl := TABLE_FEEDBACK_ON_USERS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "feedback_on_users",
	}}
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.DELETED_AT = sq.NewTimeField("deleted_at", tbl.TableInfo)
	tbl.EVALUATEE_USER_ROLE_ID = sq.NewNumberField("evaluatee_user_role_id", tbl.TableInfo)
	tbl.EVALUATOR_TEAM_ID = sq.NewNumberField("evaluator_team_id", tbl.TableInfo)
	tbl.FEEDBACK_DATA = sq.NewJSONField("feedback_data", tbl.TableInfo)
	tbl.FEEDBACK_FORM_ID = sq.NewNumberField("feedback_form_id", tbl.TableInfo)
	tbl.FEEDBACK_ID_ON_USER = sq.NewNumberField("feedback_id_on_user", tbl.TableInfo)
	tbl.OVERRIDE_OPEN = sq.NewBooleanField("override_open", tbl.TableInfo)
	tbl.SUBMITTED = sq.NewBooleanField("submitted", tbl.TableInfo)
	tbl.UPDATED_AT = sq.NewTimeField("updated_at", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_FEEDBACK_ON_USERS) As(alias string) TABLE_FEEDBACK_ON_USERS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_FORMS references the devlab.forms table.
type TABLE_FORMS struct {
	*sq.TableInfo
	CREATED_AT sq.TimeField
	DELETED_AT sq.TimeField
	FORM_ID    sq.NumberField
	NAME       sq.StringField
	PERIOD_ID  sq.NumberField
	QUESTIONS  sq.JSONField
	SUBSECTION sq.StringField
	UPDATED_AT sq.TimeField
}

// FORMS creates an instance of the devlab.forms table.
func FORMS() TABLE_FORMS {
	tbl := TABLE_FORMS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "forms",
	}}
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.DELETED_AT = sq.NewTimeField("deleted_at", tbl.TableInfo)
	tbl.FORM_ID = sq.NewNumberField("form_id", tbl.TableInfo)
	tbl.NAME = sq.NewStringField("name", tbl.TableInfo)
	tbl.PERIOD_ID = sq.NewNumberField("period_id", tbl.TableInfo)
	tbl.QUESTIONS = sq.NewJSONField("questions", tbl.TableInfo)
	tbl.SUBSECTION = sq.NewStringField("subsection", tbl.TableInfo)
	tbl.UPDATED_AT = sq.NewTimeField("updated_at", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_FORMS) As(alias string) TABLE_FORMS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_FORMS_AUTHORIZED_ROLES references the devlab.forms_authorized_roles table.
type TABLE_FORMS_AUTHORIZED_ROLES struct {
	*sq.TableInfo
	FORM_ID sq.NumberField
	ROLE    sq.StringField
}

// FORMS_AUTHORIZED_ROLES creates an instance of the devlab.forms_authorized_roles table.
func FORMS_AUTHORIZED_ROLES() TABLE_FORMS_AUTHORIZED_ROLES {
	tbl := TABLE_FORMS_AUTHORIZED_ROLES{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "forms_authorized_roles",
	}}
	tbl.FORM_ID = sq.NewNumberField("form_id", tbl.TableInfo)
	tbl.ROLE = sq.NewStringField("role", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_FORMS_AUTHORIZED_ROLES) As(alias string) TABLE_FORMS_AUTHORIZED_ROLES {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_MEDIA references the devlab.media table.
type TABLE_MEDIA struct {
	*sq.TableInfo
	CREATED_AT  sq.TimeField
	DATA        sq.BinaryField
	DELETED_AT  sq.TimeField
	DESCRIPTION sq.StringField
	NAME        sq.StringField
	TYPE        sq.StringField
	UPDATED_AT  sq.TimeField
	UUID        sq.BinaryField
}

// MEDIA creates an instance of the devlab.media table.
func MEDIA() TABLE_MEDIA {
	tbl := TABLE_MEDIA{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "media",
	}}
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.DATA = sq.NewBinaryField("data", tbl.TableInfo)
	tbl.DELETED_AT = sq.NewTimeField("deleted_at", tbl.TableInfo)
	tbl.DESCRIPTION = sq.NewStringField("description", tbl.TableInfo)
	tbl.NAME = sq.NewStringField("name", tbl.TableInfo)
	tbl.TYPE = sq.NewStringField("type", tbl.TableInfo)
	tbl.UPDATED_AT = sq.NewTimeField("updated_at", tbl.TableInfo)
	tbl.UUID = sq.NewBinaryField("uuid", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_MEDIA) As(alias string) TABLE_MEDIA {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_MILESTONE_ENUM references the devlab.milestone_enum table.
type TABLE_MILESTONE_ENUM struct {
	*sq.TableInfo
	MILESTONE sq.StringField
}

// MILESTONE_ENUM creates an instance of the devlab.milestone_enum table.
func MILESTONE_ENUM() TABLE_MILESTONE_ENUM {
	tbl := TABLE_MILESTONE_ENUM{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "milestone_enum",
	}}
	tbl.MILESTONE = sq.NewStringField("milestone", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_MILESTONE_ENUM) As(alias string) TABLE_MILESTONE_ENUM {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_MIME_TYPE_ENUM references the devlab.mime_type_enum table.
type TABLE_MIME_TYPE_ENUM struct {
	*sq.TableInfo
	TYPE sq.StringField
}

// MIME_TYPE_ENUM creates an instance of the devlab.mime_type_enum table.
func MIME_TYPE_ENUM() TABLE_MIME_TYPE_ENUM {
	tbl := TABLE_MIME_TYPE_ENUM{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "mime_type_enum",
	}}
	tbl.TYPE = sq.NewStringField("type", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_MIME_TYPE_ENUM) As(alias string) TABLE_MIME_TYPE_ENUM {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_PERIODS references the devlab.periods table.
type TABLE_PERIODS struct {
	*sq.TableInfo
	COHORT     sq.StringField
	CREATED_AT sq.TimeField
	DELETED_AT sq.TimeField
	END_AT     sq.TimeField
	MILESTONE  sq.StringField
	PERIOD_ID  sq.NumberField
	STAGE      sq.StringField
	START_AT   sq.TimeField
	UPDATED_AT sq.TimeField
}

// PERIODS creates an instance of the devlab.periods table.
func PERIODS() TABLE_PERIODS {
	tbl := TABLE_PERIODS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "periods",
	}}
	tbl.COHORT = sq.NewStringField("cohort", tbl.TableInfo)
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.DELETED_AT = sq.NewTimeField("deleted_at", tbl.TableInfo)
	tbl.END_AT = sq.NewTimeField("end_at", tbl.TableInfo)
	tbl.MILESTONE = sq.NewStringField("milestone", tbl.TableInfo)
	tbl.PERIOD_ID = sq.NewNumberField("period_id", tbl.TableInfo)
	tbl.STAGE = sq.NewStringField("stage", tbl.TableInfo)
	tbl.START_AT = sq.NewTimeField("start_at", tbl.TableInfo)
	tbl.UPDATED_AT = sq.NewTimeField("updated_at", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_PERIODS) As(alias string) TABLE_PERIODS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_PROJECT_CATEGORY_ENUM references the devlab.project_category_enum table.
type TABLE_PROJECT_CATEGORY_ENUM struct {
	*sq.TableInfo
	PROJECT_CATEGORY sq.StringField
}

// PROJECT_CATEGORY_ENUM creates an instance of the devlab.project_category_enum table.
func PROJECT_CATEGORY_ENUM() TABLE_PROJECT_CATEGORY_ENUM {
	tbl := TABLE_PROJECT_CATEGORY_ENUM{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "project_category_enum",
	}}
	tbl.PROJECT_CATEGORY = sq.NewStringField("project_category", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_PROJECT_CATEGORY_ENUM) As(alias string) TABLE_PROJECT_CATEGORY_ENUM {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_PROJECT_LEVEL_ENUM references the devlab.project_level_enum table.
type TABLE_PROJECT_LEVEL_ENUM struct {
	*sq.TableInfo
	PROJECT_LEVEL sq.StringField
}

// PROJECT_LEVEL_ENUM creates an instance of the devlab.project_level_enum table.
func PROJECT_LEVEL_ENUM() TABLE_PROJECT_LEVEL_ENUM {
	tbl := TABLE_PROJECT_LEVEL_ENUM{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "project_level_enum",
	}}
	tbl.PROJECT_LEVEL = sq.NewStringField("project_level", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_PROJECT_LEVEL_ENUM) As(alias string) TABLE_PROJECT_LEVEL_ENUM {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_ROLE_ENUM references the devlab.role_enum table.
type TABLE_ROLE_ENUM struct {
	*sq.TableInfo
	ROLE sq.StringField
}

// ROLE_ENUM creates an instance of the devlab.role_enum table.
func ROLE_ENUM() TABLE_ROLE_ENUM {
	tbl := TABLE_ROLE_ENUM{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "role_enum",
	}}
	tbl.ROLE = sq.NewStringField("role", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_ROLE_ENUM) As(alias string) TABLE_ROLE_ENUM {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_SESSIONS references the devlab.sessions table.
type TABLE_SESSIONS struct {
	*sq.TableInfo
	CREATED_AT sq.TimeField
	HASH       sq.StringField
	USER_ID    sq.NumberField
}

// SESSIONS creates an instance of the devlab.sessions table.
func SESSIONS() TABLE_SESSIONS {
	tbl := TABLE_SESSIONS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "sessions",
	}}
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.HASH = sq.NewStringField("hash", tbl.TableInfo)
	tbl.USER_ID = sq.NewNumberField("user_id", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_SESSIONS) As(alias string) TABLE_SESSIONS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_STAGE_ENUM references the devlab.stage_enum table.
type TABLE_STAGE_ENUM struct {
	*sq.TableInfo
	STAGE sq.StringField
}

// STAGE_ENUM creates an instance of the devlab.stage_enum table.
func STAGE_ENUM() TABLE_STAGE_ENUM {
	tbl := TABLE_STAGE_ENUM{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "stage_enum",
	}}
	tbl.STAGE = sq.NewStringField("stage", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_STAGE_ENUM) As(alias string) TABLE_STAGE_ENUM {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_SUBMISSIONS references the devlab.submissions table.
type TABLE_SUBMISSIONS struct {
	*sq.TableInfo
	CREATED_AT         sq.TimeField
	DELETED_AT         sq.TimeField
	OVERRIDE_OPEN      sq.BooleanField
	POSTER             sq.StringField
	README             sq.StringField
	SUBMISSION_DATA    sq.JSONField
	SUBMISSION_FORM_ID sq.NumberField
	SUBMISSION_ID      sq.NumberField
	SUBMITTED          sq.BooleanField
	TEAM_ID            sq.NumberField
	UPDATED_AT         sq.TimeField
	VIDEO              sq.StringField
}

// SUBMISSIONS creates an instance of the devlab.submissions table.
func SUBMISSIONS() TABLE_SUBMISSIONS {
	tbl := TABLE_SUBMISSIONS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "submissions",
	}}
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.DELETED_AT = sq.NewTimeField("deleted_at", tbl.TableInfo)
	tbl.OVERRIDE_OPEN = sq.NewBooleanField("override_open", tbl.TableInfo)
	tbl.POSTER = sq.NewStringField("poster", tbl.TableInfo)
	tbl.README = sq.NewStringField("readme", tbl.TableInfo)
	tbl.SUBMISSION_DATA = sq.NewJSONField("submission_data", tbl.TableInfo)
	tbl.SUBMISSION_FORM_ID = sq.NewNumberField("submission_form_id", tbl.TableInfo)
	tbl.SUBMISSION_ID = sq.NewNumberField("submission_id", tbl.TableInfo)
	tbl.SUBMITTED = sq.NewBooleanField("submitted", tbl.TableInfo)
	tbl.TEAM_ID = sq.NewNumberField("team_id", tbl.TableInfo)
	tbl.UPDATED_AT = sq.NewTimeField("updated_at", tbl.TableInfo)
	tbl.VIDEO = sq.NewStringField("video", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_SUBMISSIONS) As(alias string) TABLE_SUBMISSIONS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_SUBMISSIONS_CATEGORIES references the devlab.submissions_categories table.
type TABLE_SUBMISSIONS_CATEGORIES struct {
	*sq.TableInfo
	CATEGORY      sq.StringField
	SUBMISSION_ID sq.NumberField
}

// SUBMISSIONS_CATEGORIES creates an instance of the devlab.submissions_categories table.
func SUBMISSIONS_CATEGORIES() TABLE_SUBMISSIONS_CATEGORIES {
	tbl := TABLE_SUBMISSIONS_CATEGORIES{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "submissions_categories",
	}}
	tbl.CATEGORY = sq.NewStringField("category", tbl.TableInfo)
	tbl.SUBMISSION_ID = sq.NewNumberField("submission_id", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_SUBMISSIONS_CATEGORIES) As(alias string) TABLE_SUBMISSIONS_CATEGORIES {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_TEAM_EVALUATION_PAIRS references the devlab.team_evaluation_pairs table.
type TABLE_TEAM_EVALUATION_PAIRS struct {
	*sq.TableInfo
	EVALUATEE_TEAM_ID sq.NumberField
	EVALUATOR_TEAM_ID sq.NumberField
}

// TEAM_EVALUATION_PAIRS creates an instance of the devlab.team_evaluation_pairs table.
func TEAM_EVALUATION_PAIRS() TABLE_TEAM_EVALUATION_PAIRS {
	tbl := TABLE_TEAM_EVALUATION_PAIRS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "team_evaluation_pairs",
	}}
	tbl.EVALUATEE_TEAM_ID = sq.NewNumberField("evaluatee_team_id", tbl.TableInfo)
	tbl.EVALUATOR_TEAM_ID = sq.NewNumberField("evaluator_team_id", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_TEAM_EVALUATION_PAIRS) As(alias string) TABLE_TEAM_EVALUATION_PAIRS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_TEAM_EVALUATIONS references the devlab.team_evaluations table.
type TABLE_TEAM_EVALUATIONS struct {
	*sq.TableInfo
	CREATED_AT              sq.TimeField
	DELETED_AT              sq.TimeField
	EVALUATEE_SUBMISSION_ID sq.NumberField
	EVALUATION_DATA         sq.JSONField
	EVALUATION_FORM_ID      sq.NumberField
	EVALUATOR_TEAM_ID       sq.NumberField
	OVERRIDE_OPEN           sq.BooleanField
	SUBMITTED               sq.BooleanField
	TEAM_EVALUATION_ID      sq.NumberField
	UPDATED_AT              sq.TimeField
}

// TEAM_EVALUATIONS creates an instance of the devlab.team_evaluations table.
func TEAM_EVALUATIONS() TABLE_TEAM_EVALUATIONS {
	tbl := TABLE_TEAM_EVALUATIONS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "team_evaluations",
	}}
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.DELETED_AT = sq.NewTimeField("deleted_at", tbl.TableInfo)
	tbl.EVALUATEE_SUBMISSION_ID = sq.NewNumberField("evaluatee_submission_id", tbl.TableInfo)
	tbl.EVALUATION_DATA = sq.NewJSONField("evaluation_data", tbl.TableInfo)
	tbl.EVALUATION_FORM_ID = sq.NewNumberField("evaluation_form_id", tbl.TableInfo)
	tbl.EVALUATOR_TEAM_ID = sq.NewNumberField("evaluator_team_id", tbl.TableInfo)
	tbl.OVERRIDE_OPEN = sq.NewBooleanField("override_open", tbl.TableInfo)
	tbl.SUBMITTED = sq.NewBooleanField("submitted", tbl.TableInfo)
	tbl.TEAM_EVALUATION_ID = sq.NewNumberField("team_evaluation_id", tbl.TableInfo)
	tbl.UPDATED_AT = sq.NewTimeField("updated_at", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_TEAM_EVALUATIONS) As(alias string) TABLE_TEAM_EVALUATIONS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_TEAMS references the devlab.teams table.
type TABLE_TEAMS struct {
	*sq.TableInfo
	ADVISER_USER_ROLE_ID sq.NumberField
	COHORT               sq.StringField
	CREATED_AT           sq.TimeField
	DELETED_AT           sq.TimeField
	MENTOR_USER_ROLE_ID  sq.NumberField
	PROJECT_IDEA         sq.StringField
	PROJECT_LEVEL        sq.StringField
	STATUS               sq.StringField
	TEAM_DATA            sq.JSONField
	TEAM_ID              sq.NumberField
	TEAM_NAME            sq.StringField
	UPDATED_AT           sq.TimeField
}

// TEAMS creates an instance of the devlab.teams table.
func TEAMS() TABLE_TEAMS {
	tbl := TABLE_TEAMS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "teams",
	}}
	tbl.ADVISER_USER_ROLE_ID = sq.NewNumberField("adviser_user_role_id", tbl.TableInfo)
	tbl.COHORT = sq.NewStringField("cohort", tbl.TableInfo)
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.DELETED_AT = sq.NewTimeField("deleted_at", tbl.TableInfo)
	tbl.MENTOR_USER_ROLE_ID = sq.NewNumberField("mentor_user_role_id", tbl.TableInfo)
	tbl.PROJECT_IDEA = sq.NewStringField("project_idea", tbl.TableInfo)
	tbl.PROJECT_LEVEL = sq.NewStringField("project_level", tbl.TableInfo)
	tbl.STATUS = sq.NewStringField("status", tbl.TableInfo)
	tbl.TEAM_DATA = sq.NewJSONField("team_data", tbl.TableInfo)
	tbl.TEAM_ID = sq.NewNumberField("team_id", tbl.TableInfo)
	tbl.TEAM_NAME = sq.NewStringField("team_name", tbl.TableInfo)
	tbl.UPDATED_AT = sq.NewTimeField("updated_at", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_TEAMS) As(alias string) TABLE_TEAMS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_TEAMS_STATUS_ENUM references the devlab.teams_status_enum table.
type TABLE_TEAMS_STATUS_ENUM struct {
	*sq.TableInfo
	STATUS sq.StringField
}

// TEAMS_STATUS_ENUM creates an instance of the devlab.teams_status_enum table.
func TEAMS_STATUS_ENUM() TABLE_TEAMS_STATUS_ENUM {
	tbl := TABLE_TEAMS_STATUS_ENUM{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "teams_status_enum",
	}}
	tbl.STATUS = sq.NewStringField("status", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_TEAMS_STATUS_ENUM) As(alias string) TABLE_TEAMS_STATUS_ENUM {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_USER_EVALUATIONS references the devlab.user_evaluations table.
type TABLE_USER_EVALUATIONS struct {
	*sq.TableInfo
	CREATED_AT              sq.TimeField
	DELETED_AT              sq.TimeField
	EVALUATEE_SUBMISSION_ID sq.NumberField
	EVALUATION_DATA         sq.JSONField
	EVALUATION_FORM_ID      sq.NumberField
	EVALUATOR_USER_ROLE_ID  sq.NumberField
	OVERRIDE_OPEN           sq.BooleanField
	SUBMITTED               sq.BooleanField
	UPDATED_AT              sq.TimeField
	USER_EVALUATION_ID      sq.NumberField
}

// USER_EVALUATIONS creates an instance of the devlab.user_evaluations table.
func USER_EVALUATIONS() TABLE_USER_EVALUATIONS {
	tbl := TABLE_USER_EVALUATIONS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "user_evaluations",
	}}
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.DELETED_AT = sq.NewTimeField("deleted_at", tbl.TableInfo)
	tbl.EVALUATEE_SUBMISSION_ID = sq.NewNumberField("evaluatee_submission_id", tbl.TableInfo)
	tbl.EVALUATION_DATA = sq.NewJSONField("evaluation_data", tbl.TableInfo)
	tbl.EVALUATION_FORM_ID = sq.NewNumberField("evaluation_form_id", tbl.TableInfo)
	tbl.EVALUATOR_USER_ROLE_ID = sq.NewNumberField("evaluator_user_role_id", tbl.TableInfo)
	tbl.OVERRIDE_OPEN = sq.NewBooleanField("override_open", tbl.TableInfo)
	tbl.SUBMITTED = sq.NewBooleanField("submitted", tbl.TableInfo)
	tbl.UPDATED_AT = sq.NewTimeField("updated_at", tbl.TableInfo)
	tbl.USER_EVALUATION_ID = sq.NewNumberField("user_evaluation_id", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_USER_EVALUATIONS) As(alias string) TABLE_USER_EVALUATIONS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_USER_ROLES references the devlab.user_roles table.
type TABLE_USER_ROLES struct {
	*sq.TableInfo
	COHORT       sq.StringField
	CREATED_AT   sq.TimeField
	DELETED_AT   sq.TimeField
	ROLE         sq.StringField
	UPDATED_AT   sq.TimeField
	USER_ID      sq.NumberField
	USER_ROLE_ID sq.NumberField
}

// USER_ROLES creates an instance of the devlab.user_roles table.
func USER_ROLES() TABLE_USER_ROLES {
	tbl := TABLE_USER_ROLES{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "user_roles",
	}}
	tbl.COHORT = sq.NewStringField("cohort", tbl.TableInfo)
	tbl.CREATED_AT = sq.NewTimeField("created_at", tbl.TableInfo)
	tbl.DELETED_AT = sq.NewTimeField("deleted_at", tbl.TableInfo)
	tbl.ROLE = sq.NewStringField("role", tbl.TableInfo)
	tbl.UPDATED_AT = sq.NewTimeField("updated_at", tbl.TableInfo)
	tbl.USER_ID = sq.NewNumberField("user_id", tbl.TableInfo)
	tbl.USER_ROLE_ID = sq.NewNumberField("user_role_id", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_USER_ROLES) As(alias string) TABLE_USER_ROLES {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_USER_ROLES_APPLICANTS references the devlab.user_roles_applicants table.
type TABLE_USER_ROLES_APPLICANTS struct {
	*sq.TableInfo
	APPLICANT_DATA    sq.JSONField
	APPLICANT_FORM_ID sq.NumberField
	APPLICATION_ID    sq.NumberField
	USER_ROLE_ID      sq.NumberField
}

// USER_ROLES_APPLICANTS creates an instance of the devlab.user_roles_applicants table.
func USER_ROLES_APPLICANTS() TABLE_USER_ROLES_APPLICANTS {
	tbl := TABLE_USER_ROLES_APPLICANTS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "user_roles_applicants",
	}}
	tbl.APPLICANT_DATA = sq.NewJSONField("applicant_data", tbl.TableInfo)
	tbl.APPLICANT_FORM_ID = sq.NewNumberField("applicant_form_id", tbl.TableInfo)
	tbl.APPLICATION_ID = sq.NewNumberField("application_id", tbl.TableInfo)
	tbl.USER_ROLE_ID = sq.NewNumberField("user_role_id", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_USER_ROLES_APPLICANTS) As(alias string) TABLE_USER_ROLES_APPLICANTS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_USER_ROLES_STUDENTS references the devlab.user_roles_students table.
type TABLE_USER_ROLES_STUDENTS struct {
	*sq.TableInfo
	STUDENT_DATA sq.JSONField
	TEAM_ID      sq.NumberField
	USER_ROLE_ID sq.NumberField
}

// USER_ROLES_STUDENTS creates an instance of the devlab.user_roles_students table.
func USER_ROLES_STUDENTS() TABLE_USER_ROLES_STUDENTS {
	tbl := TABLE_USER_ROLES_STUDENTS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "user_roles_students",
	}}
	tbl.STUDENT_DATA = sq.NewJSONField("student_data", tbl.TableInfo)
	tbl.TEAM_ID = sq.NewNumberField("team_id", tbl.TableInfo)
	tbl.USER_ROLE_ID = sq.NewNumberField("user_role_id", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_USER_ROLES_STUDENTS) As(alias string) TABLE_USER_ROLES_STUDENTS {
	tbl.TableInfo.Alias = alias
	return tbl
}

// TABLE_USERS references the devlab.users table.
type TABLE_USERS struct {
	*sq.TableInfo
	DISPLAYNAME sq.StringField
	EMAIL       sq.StringField
	PASSWORD    sq.StringField
	USER_ID     sq.NumberField
}

// USERS creates an instance of the devlab.users table.
func USERS() TABLE_USERS {
	tbl := TABLE_USERS{TableInfo: &sq.TableInfo{
		Schema: "devlab",
		Name:   "users",
	}}
	tbl.DISPLAYNAME = sq.NewStringField("displayname", tbl.TableInfo)
	tbl.EMAIL = sq.NewStringField("email", tbl.TableInfo)
	tbl.PASSWORD = sq.NewStringField("password", tbl.TableInfo)
	tbl.USER_ID = sq.NewNumberField("user_id", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_USERS) As(alias string) TABLE_USERS {
	tbl.TableInfo.Alias = alias
	return tbl
}
`