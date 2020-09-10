[![GoDoc-postgres](https://img.shields.io/badge/godoc-postgres-blue)](https://godoc.org/github.com/bokwoon95/go-structured-query/postgres)
[![GoDoc-mysql](https://img.shields.io/badge/godoc-mysql-blue)](https://godoc.org/github.com/bokwoon95/go-structured-query/mysql)
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

- <b>Avoid magic strings.</b> SQL queries written in Go are full of [magic strings](https://deviq.com/magic-strings/): strings specified directly within application code that have an impact on the application's behavior. Specifically, you have to hardcode table or column names over and over into your queries (even ORMs are guilty of this!). Such magic strings are prone to typos and hard to refactor correctly. sq generates table structs from your database and ensures that whatever query you write is always reflective of what's actually in your database.

- <b>Better null handling</b>. Handling NULLs is a bit of a pain in the ass in Go. You have to either use pointers (cannot be used in HTML templates) or sql.NullXXX structs (extra layer of indirection). sq scans NULLs as zero values, while still offering you the ability to check if the column was NULL. [more info](https://bokwoon95.github.io/sq/basics/struct-mapping.html#nulls)

- <b>The mapper doubles as the SELECT clause</b>.
    - database/sql requires you to specify the list of columns twice in the exact same order, once for SELECT-ing and once for scanning. If you mess the order up, that's an error.
    - Reflection-based mapping (struct tags) has you defining a set of possible column names to map, and then requires you to adhere to those column names. If you mistype a column name in the struct tag, that's an error. If you SELECT a column that's not present in the struct, that's an error.
    - In sq it's the other way around. You are free to choose any column to map, not just the ones defined in the struct. Columns are guaranteed to exist because they were generated from the database schema. Any columns that you do map are automatically added to the SELECT clause, so you don't have to specify the columns twice. And you can map columns to any Go variable, you aren't constrained to struct fields.
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
sqgen-postgres tables --database 'postgres://name:pass@localhost:5432/dbname?sslmode=disable'

# MySQL
sqgen-postgres tables --database 'name:pass@tcp(127.0.0.1:3306)/dbname' --schema dbname
```

For an example of what the generated file looks like, check out [postgres/devlab\_tables\_test.go](postgres/devlab_tables_test.go).

## Basics
In sq, there are three entities that you will be interacting with the most: a table, a field and a predicate.

- <b>A Table is:</b> anything you can SELECT FROM or JOIN.
    ```go
    users := tables.USERS()
    // FROM public.users
    From(users)
    //   └───┘
    //   Table

    // JOIN (SELECT users.name FROM users) AS subquery ON 1 = 1
    selectQuery := Select(users.NAME).From(users).As("subquery")
    Join(selectQuery, Int(1).EqInt(1))
    //   └─────────┘
    //      Table
    ```
    - Tables can be aliased.
        ```go
        u := tables.USERS().As("u")
        // FROM public.users AS u
        From(u)
        ```
    - There are two specialisations of Table: BaseTable and Query.
        - A <b>BaseTable</b> is a table that actually exists in the database, and not some subquery or common table expression. These are code generated, so you don't have to worry about creating them.
        - A <b>Query</b> is an instance of the SELECT, INSERT, UPDATE or DELETE query builder.
- <b>A Field is:</b> any SQL expression that can be present in the SELECT clause. This is often a table column, but it can also be a literal value or an expression made up of other expressions.
    ```go
    users := tables.USERS()
    // SELECT users.user_id, users.name
    Select(users.USER_ID, users.NAME)
    //     └───────────┘  └────────┘
    //         Field         Field

    // SELECT 'lorem ipsum', COUNT(*)
    Select(String("lorem ipsum"), Count())
    //     └───────────────────┘  └─────┘
    //             Field           Field

    // SELECT COALESCE(users.score, users.previous_score))
    Select(Coalesce(users.SCORE, users.PREVIOUS_SCORE))
    //     │        └─────────┘  └──────────────────┘│
    //     │           Field             Field       │
    //     └─────────────────────────────────────────┘
    //                        Field
    ```
    - Like Tables, Fields can also be aliased
        ```go
        u := tables.USERS().As("u")
        // SELECT u.user_id AS uid FROM public.users AS u
        Select(u.USER_ID.As("uid")).From(u)
        ```
- <b>A Predicate is:</b> something that evaluates to true or false (or NULL) in SQL. A Predicate is often made up of Fields, but a Predicate is also a Field itself.
    ```go
    u, s := tables.USERS(), tables.STUDENTS()
    // WHERE users.user_id = students.user_id
    Where( u.USER_ID.Eq(s.USER_ID) )
    //     └─────────────────────┘
    //            Predicate

    // WHERE users.user_id = students.user_id AND users.user_id <> 33
    Where( u.USER_ID.Eq(s.USER_ID), u.USER_ID.NeInt(33) )
    //     └─────────────────────┘  └─────────────────┘
    //            Predicate              Predicate

    // WHERE (users.user_id = 1 OR users.user_id > students.user_id) AND 1 = 1
    Where( Or( u.USER_ID.EqInt(1), u.USER_ID.Gt(s.USER_ID) ), Int(1).EqInt(1) )
    //     │   └────────────────┘  └─────────────────────┘ │  └─────────────┘
    //     │        Predicate             Predicate        │     Predicate
    //     └───────────────────────────────────────────────┘
    //                       Predicate
    ```
    Predicates passed to the WHERE clause are implictly AND-ed together. If you want to OR them together, wrap the predicates in `sq.Or()`.
    To invert a Predicate, use `sq.Not` e.g. `A.Eq(B) => sq.Not(A.Eq(B))`.

For more information, check out the [Basics](http://bokwoon95.github.io/sq/#basics).

For a list of example queries, check out [Query Building](http://bokwoon95.github.io/sq/#query-building).

## Project Status
This project is not v1 yet, I would like more people to try it out first and give me feedback.

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md)

## Appendix

### Why this exists

I wrote this because I needed a more convenient way to scan database rows into nested structs, some of which exist twice in the same struct due to self joined tables.
That made sqlx's StructScan unsuitable ([e.g. cannot handle `type Child struct { Father Person; Mother Person; }`](https://jmoiron.github.io/sqlx/#advancedScanning)).
database/sql's way of scanning is really verbose especially since I had about ~25 fields to scan into, some of which could be null. That's a lot of sql Null structs needed! Because I had opted to -not- pollute my domain structs with `sql.NullInt64`/ `sql.NullString` etc, I had to create a ton of intermediate Null structs just to contain the possible null fields, then transfer their zero value back into the domain struct. There had to be a better way. I just wanted their zero values, since everything in Go accomodates the zero value.

As a result I view sq as a data mapper first, and query builder second. I try my best to make the query builder as faithful to SQL as possible, but the main reason for its existence was always the [struct mapping](http://bokwoon95.github.io/sq/basics/struct-mapping.html).

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

It also means the query builder will never fail: you don't have to check for any errors from it.

### Dialect agnostic query builder?
sq is not dialect agnostic. This means I can add your favorite dialect specific features without the headache of cross-dialect compatibility. It also makes contributions easier, as you just have to focus on your own SQL dialect and not care about the others.
