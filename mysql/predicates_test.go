package sq

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestCustomPredicate_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		p           Predicate
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS()
	tests := []TT{
		{
			"basic",
			Predicatef("? = ?", u.USER_ID, "22"),
			nil,
			"users.user_id = ?",
			[]interface{}{"22"},
		},
		{
			"not",
			Predicatef("? = ?", u.USER_ID, "22").Not(),
			nil,
			"NOT users.user_id = ?",
			[]interface{}{"22"},
		},
		{
			"excludedTableQualifiers",
			Predicatef("? = ?", u.USER_ID, "22"),
			[]string{u.GetName()},
			"user_id = ?",
			[]interface{}{"22"},
		},
		{
			"Eq",
			Eq(u.USER_ID, 22),
			[]string{u.GetName()},
			"user_id = ?",
			[]interface{}{22},
		},
		{
			"Ne",
			Ne(u.USER_ID, 22),
			[]string{u.GetName()},
			"user_id <> ?",
			[]interface{}{22},
		},
		{
			"Gt",
			Gt(u.USER_ID, 22),
			[]string{u.GetName()},
			"user_id > ?",
			[]interface{}{22},
		},
		{
			"Ge",
			Ge(u.USER_ID, 22),
			[]string{u.GetName()},
			"user_id >= ?",
			[]interface{}{22},
		},
		{
			"Lt",
			Lt(u.USER_ID, 22),
			[]string{u.GetName()},
			"user_id < ?",
			[]interface{}{22},
		},
		{
			"Le",
			Le(u.USER_ID, 22),
			[]string{u.GetName()},
			"user_id <= ?",
			[]interface{}{22},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.p.AppendSQLExclude(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestExists(t *testing.T) {
	type TT struct {
		description string
		p           Predicate
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS()
	tests := []TT{
		{
			"basic",
			Exists(SelectOne().From(u).Where(u.EMAIL.LikeString("%@gmail.com"))),
			nil,
			"EXISTS(SELECT 1 FROM devlab.users WHERE users.email LIKE ?)",
			[]interface{}{"%@gmail.com"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.p.AppendSQLExclude(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestVariadicPredicate_AppendSQLExclude(t *testing.T) {
	type TT struct {
		description string
		p           Predicate
		exclude     []string
		wantQuery   string
		wantArgs    []interface{}
	}
	u := USERS()
	tests := []TT{
		{
			"no predicates",
			VariadicPredicate{},
			nil,
			"",
			nil,
		},
		{
			"one predicate",
			VariadicPredicate{
				Predicates: []Predicate{
					u.USER_ID.EqInt(1),
				},
			},
			nil,
			"users.user_id = ?",
			[]interface{}{1},
		},
		{
			"multiple predicates",
			VariadicPredicate{
				Predicates: []Predicate{
					u.USER_ID.EqInt(1),
					u.EMAIL.EqString("lorem ipsum"),
					u.DISPLAYNAME.EqString("lorem ipsum"),
				},
			},
			nil,
			"(users.user_id = ? AND users.email = ? AND users.displayname = ?)",
			[]interface{}{1, "lorem ipsum", "lorem ipsum"},
		},
		{
			"multiple predicates (explicit AND)",
			VariadicPredicate{
				Operator: PredicateAnd,
				Predicates: []Predicate{
					u.USER_ID.EqInt(1),
					u.EMAIL.EqString("lorem ipsum"),
					u.DISPLAYNAME.EqString("lorem ipsum"),
				},
			},
			nil,
			"(users.user_id = ? AND users.email = ? AND users.displayname = ?)",
			[]interface{}{1, "lorem ipsum", "lorem ipsum"},
		},
		{
			"multiple predicates (explicit OR)",
			VariadicPredicate{
				Operator: PredicateOr,
				Predicates: []Predicate{
					u.USER_ID.EqInt(1),
					u.EMAIL.EqString("lorem ipsum"),
					u.DISPLAYNAME.EqString("lorem ipsum"),
				},
			},
			nil,
			"(users.user_id = ? OR users.email = ? OR users.displayname = ?)",
			[]interface{}{1, "lorem ipsum", "lorem ipsum"},
		},
		{
			"multiple predicates (nested)",
			VariadicPredicate{
				Predicates: []Predicate{
					Or(
						u.USER_ID.EqInt(1),
						u.EMAIL.EqString("lorem ipsum"),
						u.DISPLAYNAME.EqString("lorem ipsum"),
					),
					Or(
						u.USER_ID.EqInt(2),
						u.EMAIL.EqString("dolor sit amet"),
						u.DISPLAYNAME.EqString("dolor sit amet"),
					),
				},
			},
			nil,
			"(" +
				"(users.user_id = ? OR users.email = ? OR users.displayname = ?)" +
				" AND (users.user_id = ? OR users.email = ? OR users.displayname = ?)" +
				")",
			[]interface{}{1, "lorem ipsum", "lorem ipsum", 2, "dolor sit amet", "dolor sit amet"},
		},
		{
			"one nil predicate",
			VariadicPredicate{
				Predicates: []Predicate{
					nil,
				},
			},
			nil,
			"NULL",
			nil,
		},
		{
			"one variadic predicate",
			VariadicPredicate{
				Predicates: []Predicate{
					VariadicPredicate{
						Predicates: []Predicate{
							u.USER_ID.EqInt(1),
							u.EMAIL.EqString("lorem ipsum"),
							u.DISPLAYNAME.EqString("lorem ipsum"),
						},
					},
				},
			},
			nil,
			"(users.user_id = ? AND users.email = ? AND users.displayname = ?)",
			[]interface{}{1, "lorem ipsum", "lorem ipsum"},
		},
		{
			"multiple predicates with toplevel enabled",
			VariadicPredicate{
				toplevel: true,
				Predicates: []Predicate{
					u.USER_ID.EqInt(1),
					u.EMAIL.EqString("lorem ipsum"),
					u.DISPLAYNAME.EqString("lorem ipsum"),
				},
			},
			nil,
			"users.user_id = ? AND users.email = ? AND users.displayname = ?",
			[]interface{}{1, "lorem ipsum", "lorem ipsum"},
		},
		{
			"not (one predicate)",
			Not(VariadicPredicate{
				Predicates: []Predicate{
					u.USER_ID.EqInt(1),
				},
			}),
			nil,
			"NOT users.user_id = ?",
			[]interface{}{1},
		},
		{
			"not (multiple predicates)",
			Not(VariadicPredicate{
				Predicates: []Predicate{
					u.USER_ID.EqInt(1),
					u.EMAIL.EqString("lorem ipsum"),
					u.DISPLAYNAME.EqString("lorem ipsum"),
				},
			}),
			nil,
			"NOT (users.user_id = ? AND users.email = ? AND users.displayname = ?)",
			[]interface{}{1, "lorem ipsum", "lorem ipsum"},
		},
		{
			"excludedTableQualifiers",
			VariadicPredicate{
				toplevel: true,
				Predicates: []Predicate{
					u.USER_ID.EqInt(1),
					u.EMAIL.EqString("lorem ipsum"),
					u.DISPLAYNAME.EqString("lorem ipsum"),
				},
			},
			[]string{u.GetName()},
			"user_id = ? AND email = ? AND displayname = ?",
			[]interface{}{1, "lorem ipsum", "lorem ipsum"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			buf := &strings.Builder{}
			var args []interface{}
			tt.p.AppendSQLExclude(buf, &args, nil, tt.exclude)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}
