[![GoDoc-postgres](https://img.shields.io/badge/pkg.go.dev-postgres-blue)](https://pkg.go.dev/github.com/bokwoon95/go-structured-query/postgres)
[![GoDoc-mysql](https://img.shields.io/badge/pkg.go.dev-mysql-blue)](https://pkg.go.dev/github.com/bokwoon95/go-structured-query/mysql)
![CI](https://github.com/bokwoon95/go-structured-query/workflows/CI/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bokwoon95/go-structured-query)](https://goreportcard.com/report/github.com/bokwoon95/go-structured-query)
[![Coverage Status](https://coveralls.io/repos/github/bokwoon95/go-structured-query/badge.svg?branch=master)](https://coveralls.io/github/bokwoon95/go-structured-query?branch=master)
<!-- [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://raw.githubusercontent.com/bokwoon95/go-structured-query/master/LICENSE) -->

<div align="center"><h1>sq (Structured Query)</h1></div>
<div align="center"><h5>🎯🏆 sq is a code-generated, type safe query builder and struct mapper for Go. 🏆🎯</h5></div>
<div align="center">
<!-- <a href="https://bokwoon95.github.io/sq/quickstart">Quickstart</a> -->
<!-- <span>&nbsp;&nbsp;•&nbsp;&nbsp;</span> -->
<a href="https://bokwoon95.github.io/sq/">Documentation</a>
<span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
<a href="https://bokwoon95.github.io/sq/basics/tables-fields-and-predicates.html#query-builder-reference">Reference</a>
<span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
<a href="http://bokwoon95.github.io/sq/#query-building">Examples</a>
</div>
<hr>

This package provides type safe querying on top of Go's database/sql. It is essentially a glorified string builder, but automates things in all the right places to make working with SQL queries pleasant and boilerplate-free.

- <b>Avoid magic strings.</b> SQL queries written in Go are full of [magic strings](https://deviq.com/magic-strings/): strings specified directly within application code that have an impact on the application's behavior. Specifically, you have to hardcode table or column names over and over into your queries (even ORMs are guilty of this!). Such magic strings are prone to typos and hard to change as your database schema changes. sq generates table structs from your database and ensures that whatever query you write is always reflective of what's actually in your database. [more info](https://bokwoon95.github.io/sq/basics/generating-table-types.html)

- <b>Better null handling</b>. Handling NULLs is a bit of a pain in the ass in Go. You have to either use pointers (cannot be used in HTML templates) or sql.NullXXX structs (extra layer of indirection). sq scans NULLs as zero values, while still offering you the ability to check if the column was NULL. [more info](https://bokwoon95.github.io/sq/basics/struct-mapping.html#nulls)

- <b>The mapper function *is* the SELECT clause</b>.
    - database/sql requires you to repeat the list of columns twice in the exact same order, once for SELECT-ing and once for scanning. If you mess the order up, that's an error.
    - Reflection-based mapping (struct tags) has you defining a set of possible column names to map, and then requires you repeat those columns names again in your query. If you mistype a column name in the struct tag, that's an error. If you SELECT a column that's not present in the struct, that's an error.
    - In sq whatever you SELECT is automatically mapped. **This means you just have to write your query, execute it and if there were no errors, the data is already in your Go variables.** No iterating rows, no specifying column scan order, no error checking three times. *Write your query, run it, you're done*.
    - [more info](https://bokwoon95.github.io/sq/basics/struct-mapping.html)

## Getting started
```bash
go get github.com/bokwoon95/go-structured-query
```
You will also need the dialect-specific code generator
```bash
# Postgres
go get github.com/bokwoon95/go-structured-query/cmd/sqgen-postgres

# MySQL
go get github.com/bokwoon95/go-structured-query/cmd/sqgen-mysql
```
Generate tables from your database
```bash
# for more options, check out --help

# Postgres
sqgen-postgres tables --database 'postgres://name:pass@localhost:5432/dbname?sslmode=disable' --overwrite

# MySQL
sqgen-postgres tables --database 'name:pass@tcp(127.0.0.1:3306)/dbname' --schema dbname --overwrite
```

For an example of what the generated file looks like, check out [postgres/devlab\_tables\_test.go](postgres/devlab_tables_test.go).

### Importing sq

Each SQL dialect has its own sq package. Import the sq package for the dialect you are using accordingly:
```go
// Postgres
import (
    sq "github.com/bokwoon95/go-structured-query/postgres"
)

// MySQL
import (
    sq "github.com/bokwoon95/go-structured-query/mysql"
)
```

## Examples
You just want to see code, right? Here's some.

#### SELECT
```sql
-- SQL
SELECT u.user_id, u.name, u.email, u.created_at
FROM public.users AS u
WHERE u.name = 'Bob';
```
```go
// Go
u := tables.USERS().As("u") // table is code generated
var user User
var users []User
err := sq.
    From(u).
    Where(u.NAME.EqString("Bob")).
    Selectx(func(row *sq.Row) {
        user.UserID = row.Int(u.USER_ID)
        user.Name = row.String(u.NAME)
        user.Email = row.String(u.EMAIL)
        user.CreatedAt = row.Time(u.CREATED_AT)
    }, func() {
        users = append(users, user)
    }).
    Fetch(db)
if err != nil {
    // handle error
}
```

#### INSERT
```sql
-- SQL
INSERT INTO public.users (name, email)
VALUES ('Bob', 'bob@email.com'), ('Alice', 'alice@email.com'), ('Eve', 'eve@email.com');
```
```go
// Go
u := tables.USERS().As("u") // table is code generated
users := []User{
    {Name: "Bob",   Email: "bob@email.com"},
    {Name: "Alice", Email: "alice@email.com"},
    {Name: "Eve  ", Email: "eve@email.com"},
}
rowsAffected, err := sq.
    InsertInto(u).
    Valuesx(func(col *sq.Column) {
        for _, user := range users {
            col.SetString(u.NAME, user.Name)
            col.SetString(u.EMAIL, user.Email)
        }
    }).
    Exec(db, sq.ErowsAffected)
if err != nil {
    // handle error
}
```

#### UPDATE
```sql
-- SQL
UPDATE public.users
SET name = 'Bob', password = 'qwertyuiop'
WHERE email = 'bob@email.com';
```
```go
// Go
u := tables.USERS().As("u") // table is code generated
user := User{
    Name:     "Bob",
    Email:    "bob@email.com",
    Password: "qwertyuiop",
}
rowsAffected, err := sq.
    Update(u).
    Setx(func(col *sq.Column) {
        col.SetString(u.NAME, user.Name)
        col.SetString(u.PASSWORD, user.Password)
    }).
    Where(u.EMAIL.EqString(user.Email)).
    Exec(db, sq.ErowsAffected)
if err != nil {
    // handle error
}
```

#### DELETE
```sql
-- SQL
DELETE FROM public.users AS u
USING public.user_roles AS ur
JOIN public.user_roles_students AS urs ON urs.user_role_id = ur.user_role_id
WHERE u.user_id = ur.user_id AND urs.team_id = 15;
```
```go
// Go
u   := tables.USERS().As("u")                 // tables are code generated
ur  := tables.USER_ROLES().As("ur")           // tables are code generated
urs := tables.USER_ROLES_STUDENTS().As("urs") // tables are code generated
rowsAffected, err := sq.
    DeleteFrom(u).
    Using(ur).
    Join(urs, urs.USER_ROLE_ID.Eq(ur.USER_ROLE_ID)).
    Where(
        u.USER_ID.Eq(ur.USER_ID),
        urs.TEAM_ID.EqInt(15),
    ).
    Exec(db, sq.ErowsAffected)
if err != nil {
    // handle error
}
```

For more information, check out the [Basics](http://bokwoon95.github.io/sq/#basics).

For a list of example queries, check out [Query Building](http://bokwoon95.github.io/sq/#query-building).

## Project Status
The external API is considered stable. Any changes will only be add to the API (like support for custom loggers and structured logging). If you have any feature requests or if you find bugs do open a [new issue](https://github.com/bokwoon95/go-structured-query/issues).

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md)

## Appendix

### Why this exists

I wrote this because I needed a more convenient way to scan database rows into nested structs, some of which exist twice in the same struct due to self joined tables.
That made sqlx's StructScan unsuitable ([e.g. cannot handle `type Child struct { Father Person; Mother Person; }`](https://jmoiron.github.io/sqlx/#advancedScanning)).
database/sql's way of scanning is really verbose especially since I had about ~25 fields to scan into, some of which could be null. That's a lot of sql Null structs needed! Because I had opted to -not- pollute my domain structs with `sql.NullInt64`/ `sql.NullString` etc, I had to create a ton of intermediate Null structs just to contain the possible null fields, then transfer their zero value back into the domain struct. There had to be a better way. I just wanted their zero values, since everything in Go accomodates the zero value.

sq is therefore a data mapper first, and query builder second. I try my best to make the query builder as faithful to SQL as possible, but the main reason for its existence was always the [struct mapping](http://bokwoon95.github.io/sq/basics/struct-mapping.html).

### The case for ALL\_CAPS
Here are the reasons why ALL\_CAPS is used for table and column names over the idiomatic MixedCaps:
1) [jOOQ](https://www.jooq.org/doc/latest/manual/getting-started/use-cases/jooq-as-a-sql-builder-with-code-generation/) does it.
2) It's SQL. It's fine if it doesn't follow Go convention, because it isn't Go.
    - Go requires exported fields by capitalized.
    - SQL, being case insensitive, generally uses underscores as word delimiters.
    - ALL\_CAPS is a blend that satisfies both Go's export rules and SQL's naming conventions.
    - In my opinion, it is also easier to read because table and column names visually stand out from application code.
3) Avoids clashing with interface methods. For a struct to fit the Table interface, it has to possess the methods `GetAlias()` and `GetName()`. This means that no columns can be called 'GetAlias' or 'GetName' because it would clash with the interface methods. This is sidestepped by following an entirely different naming scheme for columns i.e. ALL\_CAPS.

### On SQL Type Safety

sq makes no effort to check the semantics of your SQL queries at runtime. Any type checking is entirely enforced by what methods that you can call and argument types that you can pass to these methods. For example, You can call Asc()/Desc() and NullsFirst()/NullsLast() on any selected field and it would pass the type checker, because Asc()/Desc()/NullsFirst()/NullsLast() still return a Field interface:
```go
u := tables.USERS().As("u")
sq.Select(u.USER_ID, u.USERNAME.Asc().NullsLast()).From(u)
```
which would translate to
```sql
SELECT u.user_id, u.username ASC NULLS LAST FROM users AS u
-- obviously wrong, you can't use ASC NULLS LAST inside the SELECT clause
```
The above example passes the Go type checker, so sq will happily build the query string -- even if that SQL query is sematically wrong. In practice, as long as you aren't trying to actively do the wrong thing (like in the above example), the limited type safety will prevent you from making the most common types of errors.

It also means the query builder will never fail: there's no boilerplate error checking required. Any semantic errors will be deferred to the database to point it out to you.

### Dialect agnostic query builder?
sq is not dialect agnostic. This means I can add your favorite dialect specific SQL features without the headache of cross-dialect compatibility. It also makes contributions easier, as you just have to focus on your own SQL dialect and not care about the others.
