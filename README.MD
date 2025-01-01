# Filtering  and Sorting
This application will use go and graphql to create an api that query a mysql databsse. The graphql api endpoints will support filtering functionality. Jet library is used for dynamic sql query gneration.


graohql schema:
```gql
scalar Date
type Post {
    id: ID!
    title: String! 
    commentCount : Int @deprecated(reason: "No longer supported")
    characters: Int!
    text: String
    score: Float 
    completed: Boolean!
    datePublished: Date
    author: Author!
}

type Author {
    id: ID!
    name: String! 
    posts: [Post!]
    friends: [Author]
}

input StringFilter {
  equals: String
  contains: String
}

input IntFilter {
  equals: Int
  gt: Int
  gte: Int
  lt: Int
  lte: Int
}

input PostFilter {
  title: StringFilter
  characters: IntFilter
  isComplete: Boolean

  and: [PostFilter!]
  or: [PostFilter!]
  not: PostFilter
}

type PostAggregateResult {
  posts: [Post!]
  count: String
  avgScore: Float
}

enum SortableField {
  title
  characters
  datePublished
}

enum SortOrder {
  ASC
  DESC
}


input PostOrder {
  field: SortableField
  order: SortOrder = ASC
}

type Query {
  getPost(postID: ID!): Post
  getPosts(filter: PostFilter, order: PostOrder): [Post!]
  aggregatePost(filter: PostFilter): PostAggregateResult
}



```


Rferences:
https://dgraph.io/docs/graphql/queries/queries-overview/

https://www.youtube.com/watch?v=dDxUu-K2qdE

## initalize project
run 
```sh
mkdir go-graphql-filtering
cd go-graphql-filtering
go mod init github.com/dmrhimali/go-graphql-filtering
```

Next, create a `tools.go` file and add gqlgen as a tool dependency for your module.

**tools.go:**
```
//go:build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
)
```
To automatically add the dependency to your go.mod run

`go mod tidy`

## Building the server
Create the project skeleton

`go run github.com/99designs/gqlgen init`

This will create our suggested package layout. You 
can modify these paths in gqlgen.yml if you need to.

```sh
├── go.mod
├── go.sum
├── gqlgen.yml               - The gqlgen config file, knobs for controlling the generated code.
├── graph
│   ├── generated            - A package that only contains the generated runtime
│   │   └── generated.go
│   ├── model                - A package for all your graph models, generated or otherwise
│   │   └── models_gen.go
│   ├── resolver.go          - The root graph resolver type. This file wont get regenerated
│   ├── schema.graphqls      - Some schema. You can split the schema into as many graphql files as you like
│   └── schema.resolvers.go  - the resolver implementation for schema.graphql
└── server.go                - The entry point to your app. Customize it however you see fit
```

### sortorder and filter

we will do sort and filtering here
[filering video](https://www.youtube.com/watch?v=dDxUu-K2qdE)

[ordering](https://daily.dev/blog/graphql-order-by-a-beginners-guide)


### define custom types
We will need `date` type in our shema.graphqls which is not nativelysupported by go. so date needs to be mapped to go supported `time.Time` data type.

Have a look at [custom scalars](https://www.apollographql.com/docs/graphql-tools/scalars.html) and [this](https://github.com/i-love-flamingo/graphql/blob/master/marshaller.go#L15)

for this new data type, define data type, marshalling and unmarshalling functions:

**graph/scalars.go**
```go
package graph

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type Date string

var ErrUnexpectedValue = errors.New("unexpected value")

// MarshalDate marshals/returns time.Time to YYYY-MM-DD string
func MarshalDate(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(t.Format("2006-01-02")))
	})
}

// UnmarshalDate unmarshalls/saves YYYY-MM-DD as time.Time 
func UnmarshalDate(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(string); ok {
		if len(tmpStr) == 0 {
			return time.Time{}, nil
		}

		parse, err := time.Parse(time.DateOnly, tmpStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("date must be in format %q: %w", time.DateOnly, err)
		}

		return parse, nil
	}

	return time.Time{}, fmt.Errorf("%w: date should be a string", ErrUnexpectedValue)
}
```

define scalar inshema and use it in `datePublished` field.

***graph/schema.graphqls:**
```go
scalar Date

type Post {
    id: ID!
    title: String! 
    characters: Int!
    text: String
    score: Float 
    completed: Boolean!
    datePublished: Date
    author: Author!
}

type Author {
    id: ID!
    name: String! 
    posts: [Post!]
    friends: [Author]
}

input StringFilter {
  equals: String
  contains: String
}

input IntFilter {
  equals: Int
  gt: Int
  gte: Int
  lt: Int
  lte: Int
}

input PostFilter {
  title: StringFilter
  characters: IntFilter
  isComplete: Boolean

  and: [PostFilter!]
  or: [PostFilter!]
  not: PostFilter
}

type PostAggregateResult {
  posts: [Post!]
  count: String
  avgScore: Float
}

enum SortableField {
  title
  characters
  datePublished
}

enum SortOrder {
  ASC
  DESC
}


input PostOrder {
  field: SortableField
  order: SortOrder = ASC
}

type Query {
  getPost(postID: ID!): Post
  getPosts(filter: PostFilter, order: PostOrder): [Post!]
  aggregatePost(filter: PostFilter): PostAggregateResult
}
```

## Database setup
You have already started mysql instance in the previous step. Now we will need to create our hackernews database in that instance. To create the database run these commands.

`docker exec -it mysql bash`

It will open the bash terminal inside mysql instance.

In the next step we will open mysql repl as the root user:

`mysql -u root -p`

It will ask you for root password, enter `dbpass` and enter.

Now we are inside mysql repl. To create the database, run this command:

`CREATE DATABASE posts;`

output:
```sh
$ docker exec -it mysql bash
bash-5.1# mysql -u root -p
mysql> CREATE DATABASE posts;

mysql> SHOW DATABASES;
+--------------------+
| Database           |
+--------------------+
| posts              |
| information_schema |
| mysql              |
| performance_schema |
| sys   
                     |
+--------------------+
5 rows in set (0.02 sec)

mysql>
```

### Models and migrations
We need to create migrations for our app so every time our app runs it creates tables it needs to work properly, we are going to use golang-migrate package. Create a folder structure for our database files in the project root directory:

`internal/db/migrations`

we will create `author`, `post` and `author_post` tables. post has m:n relationship with author.

create folowing migration files. order of execution is reflected in prefixed number in file and thus table creation should be from independent to dependent tables order.

```sh
go get -u github.com/go-sql-driver/mysql
go get github.com/golang-migrate/migrate/v4/cmd/migrate
go build -tags 'mysql' -ldflags="-X main.Version=1.0.0" -o ~/go/bin/migrate github.com/golang-migrate/migrate/v4/cmd/migrate/
```

Following will create empty files in given order. file name prefix will reflect the execution order in `go migrate up`

```sh
cd internal/pkg/db/migrations/
$ ~/go/bin/migrate create -ext sql  -seq create_author_table
$ ~/go/bin/migrate create -ext sql  -seq create_post_table
$ ~/go/bin/migrate create -ext sql  -seq create_author_post_table
```

Each table has up(create) and down(drop/reverse) scripts. Fill in scripts as follows:

**internal/db/migrations/000001_create_author_table.up.sql:**
```sql
CREATE TABLE IF NOT EXISTS Author(
    ID INT NOT NULL UNIQUE AUTO_INCREMENT,
    name VARCHAR (127) NOT NULL,
    PRIMARY KEY (ID)
)
```

**internal/db/migrations/000001_create_post_table.down.sql:**
```sql
DROP TABLE IF EXISTS Author;
```

**internal/db/migrations/000002_create_post_table.up.sql:**
```sql
CREATE TABLE IF NOT EXISTS Post(
    ID INT NOT NULL UNIQUE AUTO_INCREMENT,
    title VARCHAR (127) NOT NULL,
    characters INT NOT NULL,
    text VARCHAR (127),
    score  FLOAT,
    completed BOOLEAN NOT NULL DEFAULT 0,
    datePublished DATE,
    PRIMARY KEY (ID)
)
```

**internal/db/migrations/000002_create_post_table.up.sql:**
```sql
DROP TABLE IF EXISTS Post;
```
 
**internal/db/migrations/000004_create_author_post_table.up.sql:**
```sql
CREATE TABLE IF NOT EXISTS Author_Post(
    author_id INT NOT NULL,
    post_id INT NOT NULL,
    PRIMARY KEY (author_id, post_id),
    FOREIGN KEY (author_id) REFERENCES Author(ID),
    FOREIGN KEY (post_id) REFERENCES Post(ID)
)
```

**internal/db/migrations/000004_create_author_post_table.down.sql:**
```sql
DROP TABLE IF EXISTS Author_Post;
```

Install **go mysql driver** and **golang-migrate** packages then create migrations:


### migrate table creations
run `~/go/bin/migrate -database mysql://root:dbpass@/posts -path internal/db/migrations up` from root director.

when you get `sql errors `or error:` Dirty database version `number. simply fix your sql error, delete content of schema_migratins table and rerun migration up command.

```sh
$ pwd
/Users/Rasanjalee/code/go-graphql-filtering
### Error migration
$ ~/go/bin/migrate -database mysql://root:dbpass@/posts -path internal/db/migrations up 
error: Dirty database version 1. Fix and force version.
```



in db:

```sql
mysql> use posts;

mysql> show tables;
+-------------------+
| Tables_in_posts   |
+-------------------+
| schema_migrations |
+-------------------+

mysql> select * from schema_migrations;
+---------+-------+
| version | dirty |
+---------+-------+
|       1 |     1 |
+---------+-------+
1 row in set (0.00 sec)

--delete everything from schema_migrations
mysql> delete from schema_migrations;
Query OK, 1 row affected (0.01 sec)

mysql> select * from schema_migrations;
Empty set (0.01 sec)
```

delete above row in schema_migrations and run `~/go/bin/migrate -database mysql://root:dbpass@/posts -path internal/db/migrations up`  until successful.

```sql
mysql> show tables;
+----------------------+
| Tables_in_posts      |
+----------------------+
| Author               |
| Author_Post          |
| Post                 |
| schema_migrations    |
+----------------------+
5 rows in set (0.00 sec)
```

It is cumbersome to write a long cmmand to migrate so lets create a Makefile in root directory to easily execute commands

**.env:**
```sh
MYSQL_USER=root
MYSQL_PASSWORD=dbpass
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=posts
```

**Makefile:**
```sh
include .env

#call make create_migration insert_table_post
create_migration:
	~/go/bin/migrate	create	 -dir internal/db/migrations	-ext sql	-seq $(ARG) 

migrate_up:
	~/go/bin/migrate	-path=internal/db/migrations	-database	"mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DB}"	-verbose	up

migrate_down:
	~/go/bin/migrate -path=internal/db/migrationss -database "mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DB}?"	-verbose	down

.PHONY: create_migration migrate_up migrate_down
```

in the last line are shortcuts for commands.

now you can run migrations like this:

migrate up: `make migrate_up`
migrate down: `make migrate down`
create empty migration files: `make create_migration`


wirte code to intializedb , migrate tables and close db.

**internal/db/migrations/mysql.go:**
```go
package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

var Db *sql.DB

func InitDB() {
	// Use root:dbpass@tcp(172.17.0.2)/hackernews, if you're using Windows.
	db, err := sql.Open("mysql", "root:dbpass@tcp(localhost)/posts")
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	Db = db
}

func CloseDB() error {
	return Db.Close()
}

func Migrate() {
	if err := Db.Ping(); err != nil {
		log.Fatal(err)
	}
	driver, _ := mysql.WithInstance(Db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations/mysql",
		"mysql",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

}

```

call them at server startup:

**server.go:**
```go
func main() {
    ...
	database.InitDB()
	defer database.CloseDB()
	database.Migrate()
    ...
}
```

## Database queries

lets add some data to tables:
```sh
$ cd ~/code/go-graphql-filtering

$ make create_migration  ARG="insert_post_table"
~/go/bin/migrate        create   -dir internal/db/migrations    -ext sql        -seq insert_post_table
/Users/Rasanjalee/code/go-graphql-filtering/internal/db/migrations/000005_insert_post_table.up.sql
/Users/Rasanjalee/code/go-graphql-filtering/internal/db/migrations/000005_insert_post_table.down.sql

$ make create_migration  ARG="insert_author_table"
~/go/bin/migrate        create   -dir internal/db/migrations    -ext sql        -seq insert_author_table
/Users/Rasanjalee/code/go-graphql-filtering/internal/db/migrations/000006_insert_author_table.up.sql
/Users/Rasanjalee/code/go-graphql-filtering/internal/db/migrations/000006_insert_author_table.down.sql

$ make create_migration  ARG="insert_author_post_table"
~/go/bin/migrate        create   -dir internal/db/migrations    -ext sql        -seq insert_author_post_table
/Users/Rasanjalee/code/go-graphql-filtering/internal/db/migrations/000007_insert_author_post_table.up.sql
/Users/Rasanjalee/code/go-graphql-filtering/internal/db/migrations/000007_insert_author_post_table.down.sql

```

add following:

**internal/db/migrations/000004_insert_author_table.up.sql:**
```sql
INSERT INTO Author (ID, name) VALUES
(1, 'Author 1'),
(2, 'Author 2'),
(3, 'Author 3'),
(4, 'Author 4'),
(5, 'Author 5'),
(6, 'Author 6'),
(7, 'Author 7'),
(8, 'Author 8'),
(9, 'Author 9'),
(10, 'Author 10');
```

**internal/db/migrations/000004_insert_author_table.down.sql:**
```sql
DELETE FROM Author;
```

**internal/db/migrations/000005_insert_post_table.up.sql**
```sql
INSERT INTO Post (ID, title, characters, text, score, completed, datePublished) VALUES 
(1, 'Post 1', 5, 'abcde', 4.5, true, '2023-12-31'),
(2, 'Post 2', 15, 'defgh', 3.5, false, '2021-12-31'),
(3, 'Post 3', 100, 'lmnop', 4.2, true, '2020-12-31'),
(4, 'Post 4', 500, 'ccccc', 2.0, false, '2019-12-31'),
(5, 'Post 5', 55, 'ddddd', 3.2, true, '2010-12-31'),
(6, 'Post 6', 65, 'ddddd', 5.0, true, '2020-12-31'),
(7, 'Post 7', 305, 'ddddd', 1.25, true, '2000-12-31'),
(8, 'Post 8', 450, 'ddddd', 4.0, true, '2011-12-31'),
(9, 'Post 9', 1000, 'ddddd', 3.0, true, '2012-12-31'),
(10, 'Post 10', 2, 'ddddd', 2.0, true, '2013-12-31');
```


**internal/db/migrations/000005_insert_post_table.down.sql**
```sql
DELETE FROM Post;
```

**internal/db/migrations/000008_insert_author_friend_author_table.up.sql:**
```sql
INSERT INTO Author_Post
(author_id, post_id)
VALUES
(1, 5),
(1, 8),
(2, 3),
(3, 3),
(4, 5),
(4, 1),
(6, 8),
(7, 8),
(8, 10),
(8, 2),
(9, 1),
(9, 5),
(10, 10);
```

**internal/db/migrations/000005_insert_author_post_table.down.sql**
```sql
DELETE FROM Author_Post;
```


```sql
mysql> select * from Author;
+----+-----------+
| ID | name      |
+----+-----------+
|  1 | Author 1  |
|  2 | Author 2  |
|  3 | Author 3  |
|  4 | Author 4  |
|  5 | Author 5  |
|  6 | Author 6  |
|  7 | Author 7  |
|  8 | Author 8  |
|  9 | Author 9  |
| 10 | Author 10 |
+----+-----------+
10 rows in set (0.00 sec)

mysql> select * from Post;
+----+---------+------------+-------+-------+-----------+---------------+
| ID | title   | characters | text  | score | completed | datePublished |
+----+---------+------------+-------+-------+-----------+---------------+
|  1 | Post 1  |          5 | abcde |   4.5 |         1 | 2023-12-31    |
|  2 | Post 2  |         15 | defgh |   3.5 |         0 | 2021-12-31    |
|  3 | Post 3  |        100 | lmnop |   4.2 |         1 | 2020-12-31    |
|  4 | Post 4  |        500 | ccccc |     2 |         0 | 2019-12-31    |
|  5 | Post 5  |         55 | ddddd |   3.2 |         1 | 2010-12-31    |
|  6 | Post 6  |         65 | ddddd |     5 |         1 | 2020-12-31    |
|  7 | Post 7  |        305 | ddddd |  1.25 |         1 | 2000-12-31    |
|  8 | Post 8  |        450 | ddddd |     4 |         1 | 2011-12-31    |
|  9 | Post 9  |       1000 | ddddd |     3 |         1 | 2012-12-31    |
| 10 | Post 10 |          2 | ddddd |     2 |         1 | 2013-12-31    |
+----+---------+------------+-------+-------+-----------+---------------+
10 rows in set (0.00 sec)

mysql> select * from Author_Post;
+-----------+---------+
| author_id | post_id |
+-----------+---------+
|         4 |       1 |
|         9 |       1 |
|         8 |       2 |
|         2 |       3 |
|         3 |       3 |
|         1 |       5 |
|         4 |       5 |
|         9 |       5 |
|         1 |       8 |
|         6 |       8 |
|         7 |       8 |
|         8 |      10 |
|        10 |      10 |
+-----------+---------+
13 rows in set (0.00 sec)
```

## Select query generation method

### write sql queries: ORM, QueryBuilder or Raw SQL?

[a good article](https://lyz-code.github.io/blue-book/architecture/orm_builder_query_or_raw_sql/)

#### Raw SQL
also called native SQL, is the most basic, most low-level form of database interaction. You tell the database what to do in the language of the database. Most developers should know basics of SQL

**Excels:**

1. **Flexibility**: As you are writing raw SQL code, you are not constrained by higher level abstractions.

1. **Performance**: You can use engine specific tricks to increase the performance and your queries will probably be simpler than the higher abstraction ones.

1. **Magic free**: It's easier to understand what your code does, as you scale up in the abstraction level, magic starts to appear which is nice if everything goes well, but it backfires when you encounter problems.

1. **No logic coupling**: As your models are not linked to the way you interact with the storage solution, it's easier to define a clean software architecture that follows the SOLID principles, which also allows to switch between different storage approaches.

**Cons:**

1. **SQL Injections**: As you are manually writing the queries, it's easier to fall into these vulnerabilities.
   
2. **Change management**: Databases change over time. With raw SQL, you typically don't get any support for that. You have to migrate the schema and all queries yourself.
   
3. **Query Extension**: If you have an analytical query, it's nice if you can apply slight modifications to it. It’s possible to extend a query when you have raw SQL, but it’s cumbersome. You need to touch the original query and add placeholders.
   
4. **Editor support**: As it's interpreted as a string in the editor, your editor is not able to detect typos, syntax highlight or auto complete the SQL code.
SQL knowledge: You need to know SQL to interact with the database.
Database Locking: You might use features which are specific to that database, which makes a future database switch harder.

#### query builder

Query builders are `libraries which are written in the programming language you use and use native classes and functions to build SQL queries`. Query builders typically have a fluent interface, so the queries are built by an object-oriented interface which uses method chaining.

```python
query = Query.from_(books) \
             .select("*") \
             .where(books.author_id == aid)
```

[Pypika](https://github.com/kayak/pypika) is an example for a Query Builder in Python. Note that the resulting query is still the same as in the raw code, built in another way, so abstraction level over using raw SQL is small.

**Excels:**

1. **Performance:** Same performance as using raw SQL.
1. **Magic free:** Same comprehension as using raw SQL.
1. **No logic coupling:** Same coupling as using raw SQL.
1. **Query Extension:** Given the fluent interface, it's easier to build, extend and reuse queries.


**Mitigates:**

1. **Flexibility:** You depend on the builder implementation of the language you are trying to use, but if the functionality you are trying to use is not there, you can always fall back to raw SQL.
1. **SQL Injections:** Query builders have mechanism to insert parameters into the queries in a safe way.
1. **Editor support:** The query builder prevents typos in the offered parts — .select, .from_ , .where, and as it's object oriented you have better syntax highlight and auto completion.
1. **Database Locking:** Query builders support different databases make database switch easier.
   
**Cons:**

1. **Change management:** Databases change over time. With raw SQL, you typically don't get any support for that. You have to migrate the schema and all queries yourself.
1. **SQL knowledge:** You need to know SQL to interact with the database.
1. **Query builder knowledge:** You need to know the library to interact with the database.

#### ORM

ORMs `create an object for each database table` and `allows you to interact between related objects`, in a way that you can use your object oriented programming to interact with the database even without knowing SQL.

[SQLAlchemy](https://lyz-code.github.io/blue-book/coding/python/sqlalchemy/) is an example for an ORM in Python.

This way, there is a language-native representation and thus the languages ecosystem features such as autocomplete and syntax-highlighting work.

**Excels:**

1. **Change management:** ORM come with helper programs like Alembic which can automatically detect when your models changed compared to the last known state of the database, thus it's `able to create schema migration files for you`.
1. **Query Extension:** They have a fluent interface used and developed by a lot of people, so it may have better support than query builders.
1. **SQL Injections:** As the ORM builds the queries by itself and it maintained by a large community, you're `less prone to suffer from this` vulnerabilities.
1. **Editor support:** As you are interacting with Python objects, you have full editor support for highlighting and auto-formatting, which reduces the maintenance by making the queries easier to read
1. **Database Locking:** ORM fully support different databases, so it's easy to switch between different database solutions.

**Mitigates:**

1. **SQL knowledge:** In theory you don't need to know SQL, in reality, you need to have some basic knowledge to build the tables and relationships, as well as while debugging.
   
**Cons:**

1. **Flexibility:** Being the `highest level of abstraction`, you are constrained by what the ORM solution offers, allowing you to write raw SQL and try to give enough features, so you don't notice it unless you're writing complex queries.
1. **Performance:** When you run queries with ORMs, you tend to get more than you need. This is translated in` fetching more information and executing more queries than the other solutions`. You can try to tweak it but it can be tricky, making it easy to create queries which are wrong in a subtle way.
They also` encounter the N+1 problem`, where you potentially run more queries than you need to fetch the same result.
1. **It's all magic:** ORMs are complex high level abstractions, so when you encounter errors or want to change the default behaviour, you're going to have a bad time.
1. **Big coupling:** ORM models already contain all the data you need, so you will be `tempted to use it outside of database related code`, which introduces a tight coupling between your business model and the storage solution, which decreases flexibility when changing storage drivers, makes testing harder, leads to software architectures that induce the big ball of mud by getting further from the SOLID principles.
1. **Learn the ORM:** ORMs are `complex to learn`, they have lots of features and different ways to achieve the same result, so it's hard to learn how to use them well, and usually there is no way to fulfill all your needs.
1. **Configure the ORM:** I've had a `hard time understanding the correct way to configure database connection` inside a packaged python program, both for the normal use and to define the test environment. I've first learned using the declarative way, and then I had to learn all over again for the classical mapping required by the use of the repository pattern.

#### Conclusion⚑
**Query Builders live in the sweet spot in the abstraction plane**. They give enough abstraction to ease the interaction with the database and mitigating security vulnerabilities while retaining the flexibility, performance and architecture cleanness of using raw SQL.

Although they require you to learn SQL and the query builder library, it will pay off as you develop your programs. In the end, if you don't expect to use this knowledge in the future, you may better use pandas to your small project than a SQL solution.

**Raw SQL should be used when:**
You don't mind spending some time learning SQL.
You plan to develop and maintain complex or different projects that use SQL to store data.

**Query builders should be used when:**
You don't want to learn SQL and need to create a small script that needs to perform a specific task.

**ORMs should be used when:**
Small projects where the developers are already familiar with the ORM.
Maintaining existing ORM code, although migrating to query builders should be evaluated.

we will use qury builders. 

#### some go query builders
[jet](https://github.com/go-jet/jet/blob/master/README.md#quick-start)

[sqlboiler](https://github.com/volatiletech/sqlboiler)

we will use [jet](https://github.com/go-jet/jet/blob/master/README.md#quick-start)
[wiki](https://github.com/go-jet/jet/wiki)

[squirrel](https://github.com/Masterminds/squirrel)


### go-jet Installation
Use the command bellow to add jet as a dependency into go.mod project:

`$ go get -u github.com/go-jet/jet/v2`

Jet generator can be installed in one of the following ways:

(Go1.16+) Install jet generator using go install:
`go install github.com/go-jet/jet/v2/cmd/jet@latest`

Jet generator is installed to the directory named by the `GOBIN` environment variable, which defaults to `$GOPATH/bin` or `$HOME/go/bin` if the GOPATH environment variable is not set.

```sh
ls $HOME/go/bin 
gopls       jet         migrate     staticcheck
```


Install jet generator to specific folder:
```sh
cd /Users/Rasanjalee/code
mkdir jet
git clone https://github.com/go-jet/jet.git
cd jet && go build -o /Users/Rasanjalee/jet ./cmd/jet
```

Make sure dir_path `/Users/Rasanjalee/jet` folder is added to the PATH environment variable.
```sh
vi ~/.zshrc
# add following to file and save
export PATH=$HOME/jet/:$PATH

#save
!wq

#source
source ~/.zshrc
```

### Generate SQL Builder and Model types

```sh
$ cd ~/code/go-graphql-filtering/graph

$ $HOME/jet -source=mysql -dsn="root:dbpass@tcp(localhost:3306)/posts" -path=./internal/db/.gen

Connecting to MySQL database...
Retrieving database information...
        FOUND 5 table(s), 0 view(s), 0 enum(s)
Destination directory: .gen/posts
Cleaning up destination directory...
Generating table model files...
Generating table sql builder files

```

check the .gen folder :

```sh
(base)  ✘ Rasanjalee@jayampathis-mbp  ~/code/go-graphql-filtering   main  tree ./internal/db/.gen

./internal/db/.gen
└── posts
    ├── model
    │   ├── author.go
    │   ├── author_post.go
    │   ├── post.go
    │   └── schema_migrations.go
    └── table
        ├── author.go
        ├── author_post.go
        ├── post.go
        ├── schema_migrations.go
        └── table_use_schema.go

4 directories, 9 files
```

Types from `table`, `view` and `enum` a**re used to write type safe SQL in Go**, and `model` types **are combined to store results of the SQL queries**.



### Code for queries

**graph/resolver.go:**
revert resolver to original state to remove mock data nad new resolver:

```go
package graph

type Resolver struct{}
```

**server.go:**
change Resolvers to default reolver `graph.Resolver`. Comment migrate as we manually did it.
```go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dmrhimali/go-graphql-filtering/graph"
	database "github.com/dmrhimali/go-graphql-filtering/internal/db"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.InitDB()
	defer database.CloseDB()
	database.Migrate()

	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.NewResolver()))
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

where

**internal/db/mysql.go:**
```go
package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

var Db *sql.DB

func InitDB() {
	// Use root:dbpass@tcp(172.17.0.2)/hackernews, if you're using Windows.
	db, err := sql.Open("mysql", "root:dbpass@tcp(localhost)/posts")
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	Db = db
}

func CloseDB() error {
	return Db.Close()
}

func Migrate() {
	if err := Db.Ping(); err != nil {
		log.Fatal(err)
	}
	driver, _ := mysql.WithInstance(Db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"mysql",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

}
```

Now let's write a database helper to execute and return query results to schema.resolver.

**graph/schema.resolvers.go:**
```go
package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"

	"github.com/dmrhimali/go-graphql-filtering/graph/model"
	database "github.com/dmrhimali/go-graphql-filtering/internal/db"
)

// GetPost is the resolver for the getPost field.
func (r *queryResolver) GetPost(ctx context.Context, postID string) (*model.Post, error) {
	var resultPost *model.Post
	resultPost, err := database.GetPost(postID)
	if err != nil {
		return nil, err
	}
	return resultPost, nil

}

// GetPosts is the resolver for the getPosts field.
func (r *queryResolver) GetPosts(ctx context.Context, filter *model.PostFilter, order *model.PostOrder) ([]*model.Post, error) {
	resultPosts, err := database.GetPosts(filter, order)
	if err != nil {
		return nil, err
	}
	return resultPosts, nil
}

// AggregatePost is the resolver for the aggregatePost field.
func (r *queryResolver) AggregatePost(ctx context.Context, filter *model.PostFilter) (*model.PostAggregateResult, error) {
	aggregatePosts, err := database.GetAggregatePosts(filter)
	if err != nil {
		return nil, err
	}
	return aggregatePosts, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

```


**internal/db/db_query_executor.go**
Note that we use `dynamic query building` to build queries and `printStatementInfo()` method allows you to print generated sql.

```go
package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dmrhimali/go-graphql-filtering/graph/model"
	graphModel "github.com/dmrhimali/go-graphql-filtering/graph/model"
	dbModel "github.com/dmrhimali/go-graphql-filtering/internal/db/.gen/posts/model"
	dbTable "github.com/dmrhimali/go-graphql-filtering/internal/db/.gen/posts/table"
	jetMysql "github.com/go-jet/jet/v2/mysql"
)

type GetPostDest struct {
	dbModel.Post

	Authors []struct {
		dbModel.Author
	}
}

type GetPostsDest []GetPostDest

type GetPostAggregateDest struct {
	Posts []struct {
		*GetPostDest
	}

	Count int32

	AvgScore float32
}

func GetPost(postID string) (*graphModel.Post, error) {
	//get postId as int64
	postIdVal, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	//query
	var projectionList jetMysql.ProjectionList
	projectionList = append(projectionList, dbTable.Post.ID)
	projectionList = append(projectionList, dbTable.Post.Title)
	projectionList = append(projectionList, dbTable.Post.Characters)
	projectionList = append(projectionList, dbTable.Post.Completed)
	projectionList = append(projectionList, dbTable.Post.DatePublished)
	projectionList = append(projectionList, dbTable.Post.Score)
	projectionList = append(projectionList, dbTable.Post.Text)
	projectionList = append(projectionList, dbTable.Author.AllColumns)
	query :=
		jetMysql.
			SELECT(projectionList).
			FROM(
				dbTable.Post.
					INNER_JOIN(dbTable.AuthorPost, dbTable.Post.ID.EQ(dbTable.AuthorPost.PostID)).
					INNER_JOIN(dbTable.Author, dbTable.Author.ID.EQ(dbTable.AuthorPost.AuthorID)),
			).
			WHERE(
				dbTable.Post.ID.EQ(jetMysql.Int(postIdVal)),
			)

	printStatementInfo(query)
    // Debug sql: 
    // ==============================
    // SELECT `Post`.`ID` AS "Post.ID",
    //     `Post`.title AS "Post.title",
    //     `Post`.characters AS "Post.characters",
    //     `Post`.completed AS "Post.completed",
    //     `Post`.`datePublished` AS "Post.datePublished",
    //     `Post`.score AS "Post.score",
    //     `Post`.text AS "Post.text",
    //     `Author`.`ID` AS "Author.ID",
    //     `Author`.name AS "Author.name"
    // FROM posts.`Post`
    //     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
    //     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
    // WHERE `Post`.`ID` = 1;


	//Store result into desired destination:
	var dest GetPostDest

	err = query.Query(Db, &dest)
	if err != nil {
		return nil, err
	}

	jsonSave("./internal/db/out/resultPost.json", dest)

	datePub := dest.DatePublished.Format("2006-01-02")
	id := strconv.FormatInt(int64(dest.ID), 10)

	var firstAuthor *graphModel.Author
	if len(dest.Authors) > 0 {
		firstAuthor = &graphModel.Author{
			ID:   strconv.FormatInt(int64(dest.Authors[0].ID), 10),
			Name: dest.Authors[0].Name,
		}
	}

	resultPost := &graphModel.Post{
		ID:            id,
		Title:         dest.Title,
		Characters:    int(dest.Characters),
		Text:          dest.Text,
		Score:         dest.Score,
		Completed:     dest.Completed,
		DatePublished: &datePub,
		Author:        firstAuthor,
	}
	return resultPost, nil
}

func GetPosts(filter *graphModel.PostFilter, order *graphModel.PostOrder) ([]*model.Post, error) {
	//get fields to return
	var projectionList jetMysql.ProjectionList
	projectionList = append(projectionList, dbTable.Post.ID)
	projectionList = append(projectionList, dbTable.Post.Title)
	projectionList = append(projectionList, dbTable.Post.Characters)
	projectionList = append(projectionList, dbTable.Post.Completed)
	projectionList = append(projectionList, dbTable.Post.DatePublished)
	projectionList = append(projectionList, dbTable.Post.Score)
	projectionList = append(projectionList, dbTable.Post.Text)
	projectionList = append(projectionList, dbTable.Author.AllColumns)
	//see https://github.com/go-jet/jet/wiki/SELECT#table-aliasing if you have relation in the same table like author-firend-author, manager-employee to write sql

	var orderByClause jetMysql.OrderByClause
	var whereClauseSqlExpression jetMysql.BoolExpression
	var err error

	query :=
		jetMysql.SELECT(
			projectionList,
		).FROM(
			dbTable.Post.
				INNER_JOIN(dbTable.AuthorPost, dbTable.Post.ID.EQ(dbTable.AuthorPost.PostID)).
				INNER_JOIN(dbTable.Author, dbTable.Author.ID.EQ(dbTable.AuthorPost.AuthorID)),
		)

	//get filtered query expression
	if filter != nil {
		whereClauseSqlExpression, err = GetFilterWhereExpression(filter)
		if err != nil {
			return nil, err
		}
		if whereClauseSqlExpression == nil {
			return nil, fmt.Errorf("error creating where clause")
		}

		query = query.WHERE(whereClauseSqlExpression)
	}

	//get ordered query expression
	if order != nil {
		orderByClause, err = GetOrderByClause(*order)
		if err != nil {
			return nil, err
		}
		if orderByClause == nil {
			return nil, fmt.Errorf("error creating order by clause")
		}
		query = query.ORDER_BY(orderByClause)
	}

	printStatementInfo(query)

	//execute and store results in dest:
	var dest []GetPostDest

	err = query.Query(Db, &dest)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	jsonSave("./internal/db/out/resultPosts.json", dest)

	var resultPosts []*model.Post
	for _, post := range dest {
		datePub := post.DatePublished.Format("2006-01-02")
		id := strconv.FormatInt(int64(post.ID), 10)

		var firstAuthor *graphModel.Author
		if len(post.Authors) > 0 {
			firstAuthor = &graphModel.Author{
				ID:   strconv.FormatInt(int64(post.Authors[0].ID), 10),
				Name: post.Authors[0].Name,
			}
		}
		resultPosts = append(resultPosts, &model.Post{
			ID:            id,
			Title:         post.Title,
			Characters:    int(post.Characters),
			Text:          post.Text,
			Score:         post.Score,
			Completed:     post.Completed,
			DatePublished: &datePub,
			Author:        firstAuthor,
		})
	}

	return resultPosts, nil
}

func GetOrderByClause(postOrder graphModel.PostOrder) (jetMysql.OrderByClause, error) {

	if postOrder.Field == nil {
		return nil, fmt.Errorf("sort field cannot be null")
	}

	var orderByClause jetMysql.OrderByClause
	var sortableField = *postOrder.Field
	if sortableField == graphModel.SortableFieldTitle {
		field := dbTable.Post.Title
		if postOrder.Order == nil || *postOrder.Order == graphModel.SortOrderAsc {
			orderByClause = field.ASC()
		} else {
			orderByClause = field.DESC()
		}

		return orderByClause, nil
	} else if sortableField == graphModel.SortableFieldCharacters {
		field := dbTable.Post.Characters
		if postOrder.Order == nil || *postOrder.Order == graphModel.SortOrderAsc {
			orderByClause = field.ASC()
		} else {
			orderByClause = field.DESC()
		}
		return orderByClause, nil
	} else if sortableField == graphModel.SortableFieldDatePublished {
		field := dbTable.Post.DatePublished
		if postOrder.Order == nil || *postOrder.Order == graphModel.SortOrderAsc {
			orderByClause = field.ASC()
		} else {
			orderByClause = field.DESC()
		}
		return orderByClause, nil
	} else {
		return nil, fmt.Errorf("Unsupported sort field")
	}
}

func GetFilterWhereExpression(filter *graphModel.PostFilter) (jetMysql.BoolExpression, error) {
	//get where clause
	var whereBoolExpression jetMysql.BoolExpression = jetMysql.Bool(true)

	if filter != nil {
		foundFilter := false
		//title filter
		if filter.Title != nil {
			titleBoolExpression, err := GetStringFilterBooleanExpression(dbTable.Post.Title, filter.Title)
			if err != nil || titleBoolExpression == nil {
				return nil, err
			}
			whereBoolExpression = whereBoolExpression.AND(titleBoolExpression)
			foundFilter = true
		}

		//character filter
		if filter.Characters != nil {
			var charactersBoolExpression jetMysql.BoolExpression
			var err error
			charactersBoolExpression, err = GetIntFilterBooleanExpression(dbTable.Post.Characters, filter.Characters)
			if err != nil || charactersBoolExpression == nil {
				return nil, err
			}

			whereBoolExpression = whereBoolExpression.AND(charactersBoolExpression)
			foundFilter = true
		}

		//isComplete filter
		if filter.IsComplete != nil {
			isCompleteBoolExpression, err := GetBooleanFilterBooleanExpression(dbTable.Post.Completed, filter.IsComplete)
			if err != nil || isCompleteBoolExpression == nil {
				return nil, err
			}
			whereBoolExpression = whereBoolExpression.AND(isCompleteBoolExpression)
			foundFilter = true
		}

		//and filter
		if filter.And != nil {
			nFilters := len(filter.And)
			if nFilters > 0 {
				andBoolExpression, err := GetAndFilterBooleanExpression(filter.And)
				if err != nil || andBoolExpression == nil {
					return nil, err
				}
				whereBoolExpression = whereBoolExpression.AND(andBoolExpression)
				foundFilter = true

			} else {
				return nil, fmt.Errorf("and filter is empty")
			}
		}

		//or filter
		if filter.Or != nil {
			nFilters := len(filter.Or)
			if nFilters > 0 {
				orBoolExpression, err := GetOrFilterBooleanExpression(filter.Or)
				if err != nil || orBoolExpression == nil {
					return nil, err
				}
				whereBoolExpression = whereBoolExpression.AND(orBoolExpression)
				foundFilter = true

			} else {
				return nil, fmt.Errorf("or filter is empty")
			}
		}

		//not filter
		if filter.Not != nil {
			notBoolExpression, err := GetNotFilterBooleanExpression(filter.Not)
			if err != nil || notBoolExpression == nil {
				return nil, err
			}
			whereBoolExpression = whereBoolExpression.AND(notBoolExpression)
			foundFilter = true
		}

		if foundFilter {
			return whereBoolExpression, nil
		}

		return nil, fmt.Errorf("at least one filter field must be specified")

	}
	return nil, fmt.Errorf("filter is empty")
}

func GetStringFilterBooleanExpression(dbTablefield jetMysql.ColumnString, filter *graphModel.StringFilter) (jetMysql.BoolExpression, error) {
	var stringExpression jetMysql.BoolExpression
	if filter != nil {
		if filter.Contains != nil { //CONTAINS
			val := *filter.Contains

			strVal := fmt.Sprintf("%%%s%%", val)

			stringExpression = dbTablefield.LIKE(jetMysql.String(strVal))

		} else if filter.Equals != nil { //EQULS
			strVal := *filter.Equals
			stringExpression = dbTablefield.EQ(jetMysql.String(strVal))
		}
		return stringExpression, nil
	}
	return nil, fmt.Errorf("string filter is empty")
}

func GetIntFilterBooleanExpression(dbTablefield jetMysql.ColumnInteger, filter *graphModel.IntFilter) (jetMysql.BoolExpression, error) {
	var intBoolExpression jetMysql.BoolExpression

	if filter != nil {
		if filter.Equals != nil {
			val := *filter.Equals
			intBoolExpression = dbTablefield.EQ(jetMysql.Int(int64(val)))
		} else if filter.Gt != nil {
			val := *filter.Gt
			intBoolExpression = dbTablefield.GT(jetMysql.Int(int64(val)))
		} else if filter.Gte != nil {
			val := *filter.Gte
			intBoolExpression = dbTablefield.GT_EQ(jetMysql.Int(int64(val)))
		} else if filter.Lt != nil {
			val := *filter.Lt
			intBoolExpression = dbTablefield.LT(jetMysql.Int(int64(val)))
		} else if filter.Lte != nil {
			val := *filter.Lte
			intBoolExpression = dbTablefield.LT_EQ(jetMysql.Int(int64(val)))
		}
		return intBoolExpression, nil
	}

	return nil, fmt.Errorf("int filter is empty")
}

func GetBooleanFilterBooleanExpression(dbTablefield jetMysql.ColumnBool, val *bool) (jetMysql.BoolExpression, error) {
	var booleanBoolExpression jetMysql.BoolExpression
	if val != nil {
		if *val {
			booleanBoolExpression = dbTablefield.IS_TRUE()
		} else {
			booleanBoolExpression = dbTablefield.IS_FALSE()
		}

		return booleanBoolExpression, nil
	}

	return nil, fmt.Errorf("filter is empty")
}

func GetAndFilterBooleanExpression(filters []*graphModel.PostFilter) (jetMysql.BoolExpression, error) {
	var andBoolExpression jetMysql.BoolExpression = jetMysql.Bool(true)
	nFilters := len(filters)
	if nFilters > 0 {
		for _, filter := range filters {
			if filter != nil {
				queryExpression, err := GetFilterWhereExpression(filter)

				if err != nil {
					return nil, err
				}
				andBoolExpression = andBoolExpression.AND(queryExpression)
			}
		}
		return andBoolExpression, nil
	}
	return nil, fmt.Errorf("and filter is empty")

}

func GetOrFilterBooleanExpression(filters []*graphModel.PostFilter) (jetMysql.BoolExpression, error) {
	var orBoolExpression jetMysql.BoolExpression = jetMysql.Bool(false)
	nFilters := len(filters)
	if nFilters > 0 {
		for _, filter := range filters {
			if filter != nil {

				queryExpression, err := GetFilterWhereExpression(filter)

				if err != nil {
					return nil, err
				}
				orBoolExpression = orBoolExpression.OR(queryExpression)
			}
		}
		return orBoolExpression, nil
	}
	return nil, fmt.Errorf("or filter is empty")
}

func GetNotFilterBooleanExpression(filter *graphModel.PostFilter) (jetMysql.BoolExpression, error) {
	var notBoolExpression jetMysql.BoolExpression

	if filter != nil {
		queryExpression, err := GetFilterWhereExpression(filter)
		if err != nil {
			return nil, err
		}
		notBoolExpression = queryExpression.IS_FALSE()
		return notBoolExpression, nil

	}

	return nil, fmt.Errorf("not filter is empty")
}

func GetAggregatePosts(filter *graphModel.PostFilter) (*graphModel.PostAggregateResult, error) {
	resultPosts, err := GetPosts(filter, nil)
	if err != nil {
		return nil, err
	}

	if resultPosts == nil {
		return nil, fmt.Errorf("no data found ")
	}

	count := len(resultPosts)
	countStr := strconv.FormatInt(int64(count), 10)

	avgScore := 0.0
	totalScore := 0.0
	for _, post := range resultPosts {
		totalScore += *post.Score
	}
	avgScore = totalScore / float64(count)

	aggregatePosts := &graphModel.PostAggregateResult{
		Posts:    resultPosts,
		Count:    &countStr,
		AvgScore: &avgScore,
	}
	return aggregatePosts, nil

}

func printStatementInfo(stmt jetMysql.SelectStatement) {
	query, args := stmt.Sql()

	fmt.Println("Parameterized query: ")
	fmt.Println("==============================")
	fmt.Println(query)
	fmt.Println("Arguments: ")
	fmt.Println(args)

	debugSQL := stmt.DebugSql()

	fmt.Println("\n\nDebug sql: ")
	fmt.Println("==============================")
	fmt.Println(debugSQL)
}

func jsonSave(path string, v interface{}) {
	jsonText, _ := json.MarshalIndent(v, "", "\t")

	err := os.WriteFile(path, jsonText, 0600)

	panicOnError(err)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}

}

func printJson(title string, v interface{}) {
	val, _ := json.MarshalIndent(v, "", "\t")
	fmt.Println(title, ":", string(val))
}
```

### Runserver

`go run server.go`

### issue queries 

in insomnia, post following queries to endpoint `http://localhost:8080/query`

I have exported all queries in `gql-queries-insomnia.json`. Simply import it to insomnia to try (settings> preferences> data tab> import)

### get post by id

query:

```gql
query {
  getPost ( postID: "1") {
		id
		title
		datePublished
		characters
		completed
		score
		text
		author {
			id
			name
			friends{
				id
				name
			}
		}
		
	}
}
```

response:
```json
{
	"data": {
		"getPost": {
			"id": "1",
			"title": "Post 1",
			"datePublished": "2023-12-31",
			"characters": 5,
			"completed": true,
			"score": 4.5,
			"text": "abcde",
			"author": {
				"id": "4",
				"name": "Author 4",
				"friends": null
			}
		}
	}
}
```

Debug sql: 

(from printStatementInfo() call in `internal/db/db_query_executor.go` )
```sql
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE `Post`.`ID` = 1;
```

### get all posts

query:
```gql
query {
  getPosts (filter: null, order: { field: characters, order: ASC}) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}
```

response:
```json
{
	"data": {
		"getPosts": [
			{
				"id": "10",
				"title": "Post 10",
				"characters": 2,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			},
			{
				"id": "1",
				"title": "Post 1",
				"characters": 5,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": null
				}
			},
			{
				"id": "2",
				"title": "Post 2",
				"characters": 15,
				"score": 3.5,
				"completed": false,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			},
			{
				"id": "5",
				"title": "Post 5",
				"characters": 55,
				"score": 3.2,
				"completed": true,
				"author": {
					"name": "Author 1",
					"friends": null
				}
			},
			{
				"id": "3",
				"title": "Post 3",
				"characters": 100,
				"score": 4.2,
				"completed": true,
				"author": {
					"name": "Author 2",
					"friends": null
				}
			},
			{
				"id": "8",
				"title": "Post 8",
				"characters": 450,
				"score": 4,
				"completed": true,
				"author": {
					"name": "Author 1",
					"friends": null
				}
			}
		]
	}
}
```

Debug sql: 
```sql
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
ORDER BY `Post`.characters ASC;
```

### get posts filtered by number of characters

query:
```gql
query GetPostsCharactersEquals{
  getPosts (
		filter: { 
			characters : {equals: 55}
		}
	) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}


query GetPostsCharactersGTE{
  getPosts (
		filter: { 
			characters : {lt: 55}
		}
	) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}

query GetPostsCharactersGT{
  getPosts (
		filter: { 
			characters : {gt: 55}
		}
	) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}

query GetPostsCharactersLTE{
  getPosts (
		filter: { 
			characters : {lte: 55}
		}
	) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}

query GetPostsCharactersLT{
  getPosts (
		filter: { 
			characters : {lt: 55}
		}
	) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}
```

response:
```json
//GetPostsCharactersEquals
{
	"data": {
		"getPosts": [
			{
				"id": "5",
				"title": "Post 5",
				"characters": 55,
				"score": 3.2,
				"completed": true,
				"author": {
					"name": "Author 1",
					"friends": null
				}
			}
		]
	}
}

//GetPostsCharactersGTE
{
	"data": {
		"getPosts": [
			{
				"id": "1",
				"title": "Post 1",
				"characters": 5,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": null
				}
			},
			{
				"id": "2",
				"title": "Post 2",
				"characters": 15,
				"score": 3.5,
				"completed": false,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			},
			{
				"id": "10",
				"title": "Post 10",
				"characters": 2,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			}
		]
	}
}

//GetPostsCharactersGT
{
	"data": {
		"getPosts": [
			{
				"id": "3",
				"title": "Post 3",
				"characters": 100,
				"score": 4.2,
				"completed": true,
				"author": {
					"name": "Author 2",
					"friends": null
				}
			},
			{
				"id": "8",
				"title": "Post 8",
				"characters": 450,
				"score": 4,
				"completed": true,
				"author": {
					"name": "Author 1",
					"friends": null
				}
			}
		]
	}
}

//GetPostsCharactersGTE
{
	"data": {
		"getPosts": [
			{
				"id": "1",
				"title": "Post 1",
				"characters": 5,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": null
				}
			},
			{
				"id": "2",
				"title": "Post 2",
				"characters": 15,
				"score": 3.5,
				"completed": false,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			},
			{
				"id": "10",
				"title": "Post 10",
				"characters": 2,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			}
		]
	}
}

//GetPostsCharactersLTE
{
	"data": {
		"getPosts": [
			{
				"id": "1",
				"title": "Post 1",
				"characters": 5,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": null
				}
			},
			{
				"id": "2",
				"title": "Post 2",
				"characters": 15,
				"score": 3.5,
				"completed": false,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			},
			{
				"id": "5",
				"title": "Post 5",
				"characters": 55,
				"score": 3.2,
				"completed": true,
				"author": {
					"name": "Author 1",
					"friends": null
				}
			},
			{
				"id": "10",
				"title": "Post 10",
				"characters": 2,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			}
		]
	}
}

//GetPostsCharactersLT
{
	"data": {
		"getPosts": [
			{
				"id": "1",
				"title": "Post 1",
				"characters": 5,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": null
				}
			},
			{
				"id": "2",
				"title": "Post 2",
				"characters": 15,
				"score": 3.5,
				"completed": false,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			},
			{
				"id": "10",
				"title": "Post 10",
				"characters": 2,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			}
		]
	}
}
```

Debug sql: 
```sql
-- LTE
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND (`Post`.characters <= 55);


-- LT
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND (`Post`.characters < 55);

-- GTE
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND (`Post`.characters >= 55);

-- GT
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND (`Post`.characters > 55);

-- EQ
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND (`Post`.characters = 55);
```

### getposts filtered by title

query:
```gql
query GetPostsTitleContains{
  getPosts (filter: { title : {contains: "1"}}) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}

query GetPostsTitleEquals{
  getPosts (filter: { title : {equals: "Post 1"}}) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}
```

response:
```json
//contains
{
	"data": {
		"getPosts": [
			{
				"id": "1",
				"title": "Post 1",
				"characters": 5,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": null
				}
			},
			{
				"id": "10",
				"title": "Post 10",
				"characters": 2,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			}
		]
	}
}

//equal
{
	"data": {
		"getPosts": [
			{
				"id": "1",
				"title": "Post 1",
				"characters": 5,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": null
				}
			}
		]
	}
}

```

Debug sql: 
```sql
---contains
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND (`Post`.title LIKE '%1%');

--equal
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND (`Post`.title = 'Post 1');
```

### get posts filtered by iscomplete

query:
```gql
query  GetPostsByIsCompleteFalse{
  getPosts (filter: { isComplete : false}) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}

query  GetPostsByIsCompleteTrue{
  getPosts (filter: { isComplete : true}) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}
```

response:
```json
//iscomplete = false
{
	"data": {
		"getPosts": [
			{
				"id": "2",
				"title": "Post 2",
				"characters": 15,
				"score": 3.5,
				"completed": false,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			}
		]
	}
}

//iscomplete=true
{
	"data": {
		"getPosts": [
			{
				"id": "5",
				"title": "Post 5",
				"characters": 55,
				"score": 3.2,
				"completed": true,
				"author": {
					"name": "Author 1",
					"friends": null
				}
			},
			{
				"id": "8",
				"title": "Post 8",
				"characters": 450,
				"score": 4,
				"completed": true,
				"author": {
					"name": "Author 1",
					"friends": null
				}
			},
			{
				"id": "3",
				"title": "Post 3",
				"characters": 100,
				"score": 4.2,
				"completed": true,
				"author": {
					"name": "Author 2",
					"friends": null
				}
			},
			{
				"id": "1",
				"title": "Post 1",
				"characters": 5,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": null
				}
			},
			{
				"id": "10",
				"title": "Post 10",
				"characters": 2,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			}
		]
	}
}
```

Debug sql: 
```sql
--iscomplete=false
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND `Post`.completed IS FALSE;

--iscomplete=true
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND `Post`.completed IS TRUE;

```

### get posts filtered by not

query:
```gql
query GetPostsNot{
  getPosts (filter: { not : {characters: {gte :50}}}) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}
```

response:
```json
{
	"data": {
		"getPosts": [
			{
				"id": "1",
				"title": "Post 1",
				"characters": 5,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": null
				}
			},
			{
				"id": "2",
				"title": "Post 2",
				"characters": 15,
				"score": 3.5,
				"completed": false,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			},
			{
				"id": "10",
				"title": "Post 10",
				"characters": 2,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			}
		]
	}
}
```

Debug sql: 
```sql
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND (TRUE AND (`Post`.characters >= 50)) IS FALSE;
```

### get posts filtered by and

query:
```gql
query {
  getPosts (
		filter: {
			and: [
			 	{ title: {contains: "1"}},
				{ characters : {lt: 55}}
				
			]
		},
		order: { field: title, order: ASC}
	
	) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}
```

response:
```json
{
	"data": {
		"getPosts": [
			{
				"id": "1",
				"title": "Post 1",
				"characters": 5,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": null
				}
			},
			{
				"id": "10",
				"title": "Post 10",
				"characters": 2,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			}
		]
	}
}
```

Debug sql: 
```sql
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND ((TRUE AND (TRUE AND (`Post`.title LIKE '%1%'))) AND (TRUE AND (`Post`.characters < 55)))
ORDER BY `Post`.title ASC;
```

### get posts filtered by or

query:
```gql
query {
  getPosts (
		filter: {
			or: [
			 	{ title: {contains: "11"}},
				{ characters : {lt: 55}}
				
			]
		},
		order: { field: title, order: ASC}
	
	) {
		id
		title
		characters
		score
		completed
		author {
			name
			friends {
				name
			}
		}
	}
}
```

response:
```json
{
	"data": {
		"getPosts": [
			{
				"id": "1",
				"title": "Post 1",
				"characters": 5,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": null
				}
			},
			{
				"id": "10",
				"title": "Post 10",
				"characters": 2,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			},
			{
				"id": "2",
				"title": "Post 2",
				"characters": 15,
				"score": 3.5,
				"completed": false,
				"author": {
					"name": "Author 8",
					"friends": null
				}
			}
		]
	}
}
```

Debug sql: 
```sql
SELECT `Post`.`ID` AS "Post.ID",
     `Post`.title AS "Post.title",
     `Post`.characters AS "Post.characters",
     `Post`.completed AS "Post.completed",
     `Post`.`datePublished` AS "Post.datePublished",
     `Post`.score AS "Post.score",
     `Post`.text AS "Post.text",
     `Author`.`ID` AS "Author.ID",
     `Author`.name AS "Author.name"
FROM posts.`Post`
     INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
     INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
WHERE TRUE AND ((FALSE OR (TRUE AND (`Post`.title LIKE '%11%'))) OR (TRUE AND (`Post`.characters < 55)))
ORDER BY `Post`.title ASC;
```

#### go-jet resources

[jet wiki](https://github.com/go-jet/jet/wiki)

[jet select](https://github.com/go-jet/jet/wiki/SELECT)

[jet expressions](https://github.com/go-jet/jet/wiki/Expressions)

[dynamic query generation](https://github.com/go-jet/jet/wiki/FAQ#how-to-construct-dynamic-condition)


## Graphql Directives
https://developer.ibm.com/articles/awb-graphql-directives-overview/

### Built-in graphql directives
The following list has all the built-in directives that are defined by the GraphQL specification:

- `@skip`: This directive can be used to exclude fields from a query operation conditionally.
- `@include`: Does the opposite of @skip and can be used to include fields in a query operation conditionally.
- `@deprecated`: This directive can mark a field or an enum value as deprecated and can provide a reason for deprecation to the client.
- `@specifiedBy`: This directive can be used to provide a URL for the specification of a custom scalar.

As GraphQL evolves, new execution capabilities might be introduced and exposed through directives. For example, the directives@defer and @stream have been announced by the GraphQL Working Group but aren't listed in the latest draft of the GraphQL specification.

GraphQL services and tools can also provide custom directives beyond those already mentioned.- z

### Custom GraphQL directives
Custom directives can be used to extend the functionality of GraphQL and are the preferred way to add custom behavior to a GraphQL API. Different GraphQL server and client implementations are already using custom directives to add additional functionality to GraphQL.

For example, with GraphQL APIs built with API Connect Essentials, custom directives are defined to connect with your data sources, such as the @rest, @dbquery, and @graphql directives. When using API Connect Essentials to develop your GraphQL API, you can use these directives to connect to REST APIs, databases, and other GraphQL APIs. Additionally, we have defined the @materializer and @sequence directives to mix and match data from multiple data sources in a single type.

### Directives that apply to the type system
Directives that apply to the type system are used to annotate a schema, object type, or field definition written in GraphQL SDL (schema definition language) when building a GraphQL server. Both built-in and custom directives can be used in type system directive locations, and GraphQL server implementations can then use these annotations to take additional actions. Therefore, type system directive locations are also called "schema directives" as they only exist on the GraphQL schema itself.

The following locations in a GraphQL schema are valid type system directive locations:

- SCHEMA
- SCALAR
- OBJECT
- FIELD_DEFINITION
- ARGUMENT_DEFINITION
- INTERFACE
- UNION
- ENUM & ENUM_VALUE
- INPUT_OBJECT & INPUT_FIELD_DEFINITION

`Examples` of directives that apply to the type system include `@deprecated`, a built-in directive that can `mark a field as deprecated`. Let's see what it looks like in a schema:

```gql
type User {
  id: ID!
  name: String! @deprecated(reason: "Use the firstName and lastName field instead")
  firstName: String!
  lastName: String!
  email: String!
}
```

In the example above, the @deprecated directive marks the name field as deprecated. The reason argument provides a reason for deprecation available to services that introspect the schema. The client could then, for example, warn the user that the field name is deprecated and shouldn't be used anymore.

Another example of a directive that can be applied to the type system is` @rest`, a custom directive only available in API Connect Essentials implementations. The `data for the User type in the example above is fetched from a REST API using the @rest directive`, which is defined in the following way on a query field in a GraphQL schema:

```gql
type User {
  id: ID!
  name: String!
  email: String!
}

type Query {
  user(id: ID!): User
    @rest(url: "https://jsonplaceholder.typicode.com/users/$id")
}
```


When you run an operation that includes the user field, the API Connect Essentials GraphQL API fetches the data from the REST API and returns it to the client. The @rest directive is applied to a type system location because it annotates the user field in the schema. It lets you declaratively define how the data for the user field should be fetched, rather than having to write a resolver function, as you might expect from other GraphQL server implementations.


### Directives that apply to execution
`You can use directives that apply to execution to modify the behavior of an operation, field, or fragment in runtime execution`. For example, directives that apply to execution can `include or exclude fields` or `perform additional data processing` before the response is returned.

Executable directive locations in GraphQL are:

- QUERY
- MUTATION
- SUBSCRIPTION
- FIELD
- FRAGMENT_DEFINITION & FRAGMENT_SPREAD
- INLINE_FRAGMENT
- VARIABLE_DEFINITION

Similar to directives that apply to the type system directives, both built-in and custom directives can be applied to executable locations. `Most built-in directives are executable, such as @skip and @include`, which you can use to include or exclude fields in an operation conditionally.

Let's see what the @include directive looks like in a query operation:

```gql
query me($showName: Boolean!) {
  me {
    id
    firstName @include(if: $showName)
    lastName @include(if: $showName)
    email
  }
}
```

`The @include directive conditionally includes the firstName and lastName fields in the response`. The if argument specifies a boolean value determining whether to include the field in the response. In the example above, the if argument is set to a variable $showName, which is defined in the operation variables. The variable's value can be set to true or false to include or exclude the fields in the response.


Another example of an executable directive is`@sort`, a custom directive only available in IBM API Connect Essentials implementations.

With the @sort directive, you can `sort the data returned by a field in a GraphQL operation`. You can use the @sort directive on a field that returns either a list of leaf fields or a list of objects.

For example, `let's say you have a products field that returns a list of tags`. You can use the @sort directive to sort the tags by alphabetical order, even when the data source doesn't suppport this sorting order:

```gql
query {
  products {
    tags # ['c', 'b', 'a', null]
  }
}
```

The list of tags will be transformed to this when the @sort directive is used:
```gql
query {
  products {
    tags @sort # [null, 'a', 'b', 'c']
  }
}
```

### Declaring directives
https://boopathi.blog/graphql-directives-an-introduction


`If you’re using GraphQL SDL` (Schema Definition Language) to define schema, a directive declaration would look like this -

```gql
directive @foo($arg: String!) on QUERY | MUTATION
```


`If you’re not using the schema language`, you can declare a directive using the `GraphQLDirective constructor` -

```js
import { GraphQLDirective, DirectiveLocation } from "graphql";

const fooDirective = new GraphQLDirective({
  name: "foo",
  args: {
    arg: { type: GraphQLString },
  },
  locations: [DirectiveLocation.QUERY, DirectiveLocation.MUTATION],
});
```


The declaration contains` three parts` that we can control.

- The **name** of the directive. In the above example, it is @foo
- The **arguments** of the directive and their types. In the above example, arg.
- The **places where the directive can be used**. In the above example, QUERY | MUTATION

### Directive Locations
The **possible values of where a directive can be defined is available in the GraphQL specification** — `DirectiveLocations`. As you can see in the specification, a `directive can be defined for one of the two categories of locations — Executable and TypeSystem`.

The location names in the `Executable` form are the query directives that the client can use. For example,
```gql
# In the server schema definitions,
directive @auth(token: String!) on QUERY | MUTATION

# and in the client,
query ($token: String!) @auth(token: $token) {
  ...queryFields
}
```

The location names in the TypeSystem form are the schema directives that the schema can use. But, what use does a schema directive have? To answer this question let’s start with a problem statement for our GraphQL servers.

### Executable directives locations
These are` directives that are used in the query`. For example, the `in-built directives` (at the time of this writing) such as `@skip`, `@include`, `@stream`, `@defer` are all executable directives.

The **executable directives are available for the locations listed below**. Consider @foo to be directive defined for the location mentioned in the 1st column.

```gql
directive @foo on LOCATION_IN_FIRST_COLUMN
Directive Location	Example
QUERY	query name @foo {}
MUTATION	mutation name @foo {}
SUBSCRIPTION	subscription name @foo {}
FIELD	query { product @foo {} }
FRAGMENT_DEFINITION	fragment x on Query @foo { }
FRAGMENT_SPREAD	query { ...x @foo }
INLINE_FRAGMENT	query { ... @foo { } }
VARIABLE_DEFINITION	query ($id: ID @foo) { }
```

### TypeSystem directives locations
These are **locations where the directives will be used in the schema**. The existing` in-built TypeSystem directives` in GraphQL are `@deprecated`, and `@specifiedBy`.

The type system directives are available for the below listed locations. Consider @foo to be the directive declared for the possible locations.
```gql
directive @foo on LOCATION_IN_FIRST_COLUMN
Directive Location	Example
SCHEMA	schema @foo { query: Query }
SCALAR	scalar x @foo
OBJECT	type Product @foo { }
FIELD_DEFINITION	type X { field: String @foo }
ARGUMENT_DEFINITION	type X { field(arg: Int @foo): String }
INTERFACE	interface X @foo {}
UNION	union X @foo = A | B
ENUM	enum X @foo { A B }
ENUM_VALUE	enum X { A @foo B }
INPUT_OBJECT	input X @foo { }
INPUT_FIELD_DEFINITION	input X { field: String @foo }
```

### Metadata
Let’s consider the following case.` You use persisted queries in your GraphQL server.` `Many teams build and persist these queries`.` Many versions of apps use different queries`. And also, the `queries could also be used in different platforms` — like web, iOS, or android.` You want to be able to associate failures to page the right team who built the query`. You want to draw metrics specific to an app version or platform. **All this meta can go into query directives while the author is writing persisting the query**.

As an example,

```gql
mutation addToCart($id: ID!)
@team(name: "cart")
@platform(name: IOS)
@appVersion(version: "5.15.2")
@sli(name: "add-to-cart") {
  addToCart(id: $id)
}
```


### Validation rules
Schema directives on input definitions are useful to `specify custom constraints on the input`. As examples,

#### @maxLength
```gql
type Query {
  hasTooManyCharacters(str: String! @maxLength(value: 64)): Boolean
}
```

#### @conformsRegex
```gql
type Foo {
  hasSpecialCharacter(str: String! @conformsRegex(regex: "[a-zA-Z]+")): Boolean
}
```

#### @validate, @minLength
```gql
type Mutation {
  createAccount(
    email: String! @validate(format: EMAIL)
    password: String! @validate(format: PASSWORD)
    name: String! @minLength(value: 1)
  ): CreateAccountPayload
}
```

There can be many more use-cases of custom directives for input validation depending on your business domain. The advantage of using directives here to specify these validation rules is that these are declarative. One does not need a separate document to check what the validation of each field is. The specification of these validation rules goes 1-1 with how we express it with these directives, which is a very nice property to have in our models.

### Auth
**@auth**
Though auth must be done in the request layer and not GraphQL layer, `it could be a useful tool for a GraphQL server that allows only one query document per request.`
```gql
query foo($token: String!) @auth(token: $token) {
  user {
    name
  }
}
```

**@isAuthenticated**
When we have a unified schema that covers the entire website or web app or app, certain parts of the schema might be public data and certain other parts might be user’s data. To express which parts are private data and require customer authentication, directives provide a great way to express this in declaration.
```gql
type Query {
  user: User @isAuthenticated
}
```
If you have `multiple levels of auth`, you can simply specify that using the directive’s parameter.
```gql
type Query {
  user: User @isAuthenticated(level: USER)

  allUsers(first: Int, after: String): [User!] @isAuthenticated(level: ADMIN)
}
```

### Marking sensitive data
Marking certain data as sensitive so that it doesn’t end up in our monitoring platforms is an important responsibility of a server. Directives are a great way to declare such fields.

I’ve written a dedicated blog post for this topic. It covers the @sensitive directive. You can read it here — [How to avoid logging sensitive data in GraphQL?](https://boopathi.blog/how-to-avoid-logging-sensitive-data-in-graphql)

https://ayoakinola.medium.com/implementing-the-skipauth-directive-in-graphql-with-golang-f292ed10d9b9
https://medium.com/@louisaldorio/adding-custom-tags-to-go-struct-with-custom-directives-and-gqlgen-plugins-c93322324c76

[gqlgen directives](https://gqlgen.com/reference/directives/)
