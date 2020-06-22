package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestCTEs_AppendSQL(t *testing.T) {
	type TT struct {
		description string
		ctes        CTEs
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "recursive"
			// real query from https://dba.stackexchange.com/a/257091
			folder_ids := NewRecursiveCTE("folder_ids", Queryf(
				"SELECT f.id FROM folder AS f WHERE f.id = @folderId"+
					" UNION ALL"+
					" SELECT f.id FROM folder_ids JOIN folder AS f ON f.parent_id = folder_ids.id",
			))
			user_ids := NewRecursiveCTE("user_ids", Queryf(
				"SELECT SUBSTRING_INDEX(@user_ids, ',' , 1)"+
					", CONCAT(SUBSTRING(@user_ids FROM1 + LOCATE(',', @user_ids)), ',')"+
					" UNION ALL"+
					" SELECT SUBSTRING_INDEX(slack, ',', 1)"+
					", SUBSTRING(slack FROM 2 + LENGTH(SUBSTRING_INDEX(slack, ',', 1)))"+
					" FROM user_ids WHERE slack",
			), "user_id", "slack")
			wantQuery := "WITH RECURSIVE folder_ids AS (" +
				"SELECT f.id FROM folder AS f WHERE f.id = @folderId" +
				" UNION ALL" +
				" SELECT f.id FROM folder_ids JOIN folder AS f ON f.parent_id = folder_ids.id" +
				")" +
				", user_ids (user_id, slack) AS (" +
				"SELECT SUBSTRING_INDEX(@user_ids, ',' , 1)" +
				", CONCAT(SUBSTRING(@user_ids FROM1 + LOCATE(',', @user_ids)), ',')" +
				" UNION ALL" +
				" SELECT SUBSTRING_INDEX(slack, ',', 1)" +
				", SUBSTRING(slack FROM 2 + LENGTH(SUBSTRING_INDEX(slack, ',', 1)))" +
				" FROM user_ids WHERE slack" +
				")"
			return TT{desc, CTEs{folder_ids, user_ids}, wantQuery, nil}
		}(),
		{
			"empty cte",
			CTEs{NewCTE("empty_query", nil)},
			"WITH empty_query AS (NULL)",
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.ctes.AppendSQL(buf, &args)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestCTE_AppendSQL(t *testing.T) {
	type TT struct {
		description string
		cte         CTE
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{
		func() TT {
			desc := "basic"
			folder_ids := NewRecursiveCTE("folder_ids", Queryf(
				"SELECT f.id FROM folder AS f WHERE f.id = @folderId"+
					" UNION ALL"+
					" SELECT f.id FROM folder_ids JOIN folder AS f ON f.parent_id = folder_ids.id",
			))
			wantQuery := "folder_ids"
			return TT{desc, folder_ids, wantQuery, nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.cte.AppendSQL(buf, &args)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestCTE_BasicTesting(t *testing.T) {
	is := is.New(t)
	var buf = &strings.Builder{}
	var args []interface{}

	// AppendSQL
	cte := NewRecursiveCTE("folder_ids", Queryf(
		"SELECT f.id FROM folder AS f WHERE f.id = @folderId"+
			" UNION ALL"+
			" SELECT f.id FROM folder_ids JOIN folder AS f ON f.parent_id = folder_ids.id",
	))
	cte.AppendSQL(buf, &args)
	is.Equal("folder_ids", buf.String())

	// GetAlias
	is.Equal("", cte.GetAlias())

	// GetName
	is.Equal("folder_ids", cte.GetName())

	// Get
	f := cte.Get("supercalifragilisticexpialidocious")
	buf.Reset()
	f.AppendSQLExclude(buf, &args, nil)
	is.Equal("folder_ids.supercalifragilisticexpialidocious", buf.String())

	// AliasedCTE
	aliased := cte.As("hurr_durr")

	// AliasedCTE AppendSQL
	buf.Reset()
	aliased.AppendSQL(buf, &args)
	is.Equal("folder_ids AS hurr_durr", buf.String()+" AS "+aliased.GetAlias())

	// AliasedCTE GetAlias
	is.Equal("hurr_durr", aliased.GetAlias())

	// AliasedCTE GetName
	is.Equal("folder_ids", aliased.GetName())

	// AliasedCTE Get
	f = aliased.Get("supercalifragilisticexpialidocious")
	buf.Reset()
	f.AppendSQLExclude(buf, &args, nil)
	is.Equal("hurr_durr.supercalifragilisticexpialidocious", buf.String())
}
