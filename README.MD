# Filtering 

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

## Define your schema

update schema:


```gql
type Post {
    id: ID!
    title: String! @search
    text: String
    score: Float @search
    completed: Boolean @search
    datePublished: DateTime @search(by: [year])
    author: Author!
}

type Author {
    id: ID!
    name: String! @search
    posts: [Post!]
    friends: [Author]
}
```

### sortorder and filter
https://www.youtube.com/watch?v=dDxUu-K2qdE

//https://daily.dev/blog/graphql-order-by-a-beginners-guide for orderibf



# directives
https://nelsondominguez.hashnode.dev/inheritance-in-graphql-when-and-how-to-use-it

Directives are used to a`dd metadata and functionality to GraphQL types and/or their fields`. This metadata can be read at runtime to modify the behaviour of a GraphQL schema.

There are **two types** of directives in GraphQL:

1. **Schema directives**
Directives whose sole purpose is to` attach metadata` to a type or field.
E.g: @`deprecated`

2. **Operation directives**
Directives that m`anipulate the output of a GraphQL operation` (query, mutation).
E.g: `@skip`, `@include`

**Use cases** for directives include but are not limited to:

- including or excluding fields
- manipulating query arguments
- controlling access to specific fields

A GraphQL directive is a pure Schema Definition Language (`SDL`) feature.

**Exmple:**

We declare directives using the following SDL syntax:
```gql
# A directive to enforce a min/max length on string-valued field.
directive @length(min: Int, max: Int) on FIELD_DEFINITION

type Tweet {
  # The content of a tweet cannot exceed a length of 160 characters.
  content: String! @length(max: 160)
}

```

# custom types
date

Have a look at custom scalars: https://www.apollographql.com/docs/graphql-tools/scalars.html

date needs to be mapped to go time.Time data type.
https://github.com/i-love-flamingo/graphql/blob/master/marshaller.go#L15



Query with all parameters optional can't be called with empty (). It's expected since GraphQL specification requires argument list to be non-empty:



## query

### get all posts
```gql
query {
  getPosts (filter: {}) {
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
				"id": "19",
				"title": "Post 18",
				"characters": 180,
				"score": 1.2,
				"completed": false,
				"author": {
					"name": "Author 9",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "21",
				"title": "Post 20",
				"characters": 200,
				"score": 1.4,
				"completed": true,
				"author": {
					"name": "Author 1",
					"friends": [
						{
							"name": "Author 8"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "49",
				"title": "Post 48",
				"characters": 480,
				"score": 1.1,
				"completed": false,
				"author": {
					"name": "Author 1",
					"friends": [
						{
							"name": "Author 8"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "51",
				"title": "Post 50",
				"characters": 500,
				"score": 2.3,
				"completed": false,
				"author": {
					"name": "Author 3",
					"friends": [
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "59",
				"title": "Post 58",
				"characters": 580,
				"score": 2.9,
				"completed": false,
				"author": {
					"name": "",
					"friends": null
				}
			},
			{
				"id": "74",
				"title": "Post 73",
				"characters": 730,
				"score": 3.1,
				"completed": true,
				"author": {
					"name": "",
					"friends": null
				}
			},
			{
				"id": "80",
				"title": "Post 79",
				"characters": 790,
				"score": 2.7,
				"completed": true,
				"author": {
					"name": "Author 17",
					"friends": [
						{
							"name": "Author 18"
						},
						{
							"name": "Author 4"
						},
						{
							"name": "Author 7"
						}
					]
				}
			},
			{
				"id": "5",
				"title": "Post 4",
				"characters": 40,
				"score": 4.4,
				"completed": true,
				"author": {
					"name": "Author 10",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 7"
						}
					]
				}
			},
			{
				"id": "91",
				"title": "Post 90",
				"characters": 900,
				"score": 2.2,
				"completed": false,
				"author": {
					"name": "Author 15",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 6"
						},
						{
							"name": "Author 8"
						}
					]
				}
			},
			{
				"id": "37",
				"title": "Post 36",
				"characters": 360,
				"score": 1.9,
				"completed": true,
				"author": {
					"name": "Author 12",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "41",
				"title": "Post 40",
				"characters": 400,
				"score": 2.5,
				"completed": false,
				"author": {
					"name": "",
					"friends": null
				}
			},
			{
				"id": "70",
				"title": "Post 69",
				"characters": 690,
				"score": 4.6,
				"completed": true,
				"author": {
					"name": "Author 9",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "87",
				"title": "Post 86",
				"characters": 860,
				"score": 3,
				"completed": false,
				"author": {
					"name": "Author 18",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 4"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "17",
				"title": "Post 16",
				"characters": 160,
				"score": 4.9,
				"completed": true,
				"author": {
					"name": "",
					"friends": null
				}
			},
			{
				"id": "22",
				"title": "Post 21",
				"characters": 210,
				"score": 4.3,
				"completed": false,
				"author": {
					"name": "Author 5",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						}
					]
				}
			},
			{
				"id": "29",
				"title": "Post 28",
				"characters": 280,
				"score": 3.9,
				"completed": true,
				"author": {
					"name": "Author 6",
					"friends": [
						{
							"name": "Author 4"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 1"
						}
					]
				}
			},
			{
				"id": "31",
				"title": "Post 30",
				"characters": 300,
				"score": 1.2,
				"completed": false,
				"author": {
					"name": "Author 13",
					"friends": [
						{
							"name": "Author 15"
						},
						{
							"name": "Author 2"
						},
						{
							"name": "Author 4"
						}
					]
				}
			},
			{
				"id": "36",
				"title": "Post 35",
				"characters": 350,
				"score": 3.1,
				"completed": false,
				"author": {
					"name": "Author 4",
					"friends": [
						{
							"name": ""
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 17"
						}
					]
				}
			},
			{
				"id": "65",
				"title": "Post 64",
				"characters": 640,
				"score": 1.2,
				"completed": false,
				"author": {
					"name": "",
					"friends": null
				}
			},
			{
				"id": "76",
				"title": "Post 75",
				"characters": 750,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 1",
					"friends": [
						{
							"name": "Author 8"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "81",
				"title": "Post 80",
				"characters": 800,
				"score": 2.3,
				"completed": false,
				"author": {
					"name": "Author 9",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "6",
				"title": "Post 5",
				"characters": 50,
				"score": 2.9,
				"completed": true,
				"author": {
					"name": "Author 12",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "99",
				"title": "Post 98",
				"characters": 980,
				"score": 2.2,
				"completed": true,
				"author": {
					"name": "Author 5",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						}
					]
				}
			},
			{
				"id": "98",
				"title": "Post 97",
				"characters": 970,
				"score": 3,
				"completed": true,
				"author": {
					"name": "Author 13",
					"friends": [
						{
							"name": "Author 15"
						},
						{
							"name": "Author 2"
						},
						{
							"name": "Author 4"
						}
					]
				}
			},
			{
				"id": "14",
				"title": "Post 13",
				"characters": 130,
				"score": 1.7,
				"completed": false,
				"author": {
					"name": "Author 12",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "72",
				"title": "Post 71",
				"characters": 710,
				"score": 4.9,
				"completed": true,
				"author": {
					"name": "Author 16",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 6"
						}
					]
				}
			},
			{
				"id": "88",
				"title": "Post 87",
				"characters": 870,
				"score": 2.6,
				"completed": false,
				"author": {
					"name": "Author 14",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 9"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "92",
				"title": "Post 91",
				"characters": 910,
				"score": 3.3,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": [
						{
							"name": "Author 8"
						},
						{
							"name": "Author 17"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "94",
				"title": "Post 93",
				"characters": 930,
				"score": 3.3,
				"completed": false,
				"author": {
					"name": "Author 13",
					"friends": [
						{
							"name": "Author 15"
						},
						{
							"name": "Author 2"
						},
						{
							"name": "Author 4"
						}
					]
				}
			},
			{
				"id": "3",
				"title": "Post 2",
				"characters": 20,
				"score": 4.3,
				"completed": true,
				"author": {
					"name": "Author 10",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 7"
						}
					]
				}
			},
			{
				"id": "7",
				"title": "Post 6",
				"characters": 60,
				"score": 4.7,
				"completed": false,
				"author": {
					"name": "Author 9",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "28",
				"title": "Post 27",
				"characters": 270,
				"score": 2.1,
				"completed": false,
				"author": {
					"name": "Author 2",
					"friends": [
						{
							"name": "Author 12"
						},
						{
							"name": ""
						},
						{
							"name": "Author 6"
						}
					]
				}
			},
			{
				"id": "47",
				"title": "Post 46",
				"characters": 460,
				"score": 4.8,
				"completed": true,
				"author": {
					"name": "Author 3",
					"friends": [
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "66",
				"title": "Post 65",
				"characters": 650,
				"score": 2.6,
				"completed": true,
				"author": {
					"name": "Author 11",
					"friends": [
						{
							"name": "Author 6"
						},
						{
							"name": "Author 5"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "75",
				"title": "Post 74",
				"characters": 740,
				"score": 2.3,
				"completed": true,
				"author": {
					"name": "Author 18",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 4"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "79",
				"title": "Post 78",
				"characters": 780,
				"score": 3,
				"completed": true,
				"author": {
					"name": "Author 18",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 4"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "100",
				"title": "Post 99",
				"characters": 990,
				"score": 1.7,
				"completed": false,
				"author": {
					"name": "Author 10",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 7"
						}
					]
				}
			},
			{
				"id": "1",
				"title": "Post 0",
				"characters": 0,
				"score": 1.2,
				"completed": false,
				"author": {
					"name": "Author 7",
					"friends": [
						{
							"name": "Author 0"
						},
						{
							"name": "Author 1"
						},
						{
							"name": "Author 8"
						}
					]
				}
			},
			{
				"id": "54",
				"title": "Post 53",
				"characters": 530,
				"score": 4,
				"completed": true,
				"author": {
					"name": "Author 19",
					"friends": [
						{
							"name": "Author 2"
						},
						{
							"name": "Author 17"
						},
						{
							"name": ""
						}
					]
				}
			},
			{
				"id": "60",
				"title": "Post 59",
				"characters": 590,
				"score": 1.7,
				"completed": true,
				"author": {
					"name": "Author 9",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "62",
				"title": "Post 61",
				"characters": 610,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 6",
					"friends": [
						{
							"name": "Author 4"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 1"
						}
					]
				}
			},
			{
				"id": "83",
				"title": "Post 82",
				"characters": 820,
				"score": 2.3,
				"completed": true,
				"author": {
					"name": "Author 12",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "24",
				"title": "Post 23",
				"characters": 230,
				"score": 4.6,
				"completed": false,
				"author": {
					"name": "Author 15",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 6"
						},
						{
							"name": "Author 8"
						}
					]
				}
			},
			{
				"id": "61",
				"title": "Post 60",
				"characters": 600,
				"score": 3.1,
				"completed": true,
				"author": {
					"name": "Author 5",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						}
					]
				}
			},
			{
				"id": "67",
				"title": "Post 66",
				"characters": 660,
				"score": 1.4,
				"completed": false,
				"author": {
					"name": "Author 19",
					"friends": [
						{
							"name": "Author 2"
						},
						{
							"name": "Author 17"
						},
						{
							"name": ""
						}
					]
				}
			},
			{
				"id": "68",
				"title": "Post 67",
				"characters": 670,
				"score": 4.3,
				"completed": true,
				"author": {
					"name": "Author 4",
					"friends": [
						{
							"name": ""
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 17"
						}
					]
				}
			},
			{
				"id": "86",
				"title": "Post 85",
				"characters": 850,
				"score": 3.3,
				"completed": false,
				"author": {
					"name": "Author 5",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						}
					]
				}
			},
			{
				"id": "57",
				"title": "Post 56",
				"characters": 560,
				"score": 2.6,
				"completed": false,
				"author": {
					"name": "Author 17",
					"friends": [
						{
							"name": "Author 18"
						},
						{
							"name": "Author 4"
						},
						{
							"name": "Author 7"
						}
					]
				}
			},
			{
				"id": "16",
				"title": "Post 15",
				"characters": 150,
				"score": 3.1,
				"completed": false,
				"author": {
					"name": "Author 0",
					"friends": [
						{
							"name": "Author 13"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 2"
						}
					]
				}
			},
			{
				"id": "23",
				"title": "Post 22",
				"characters": 220,
				"score": 2.7,
				"completed": true,
				"author": {
					"name": "Author 7",
					"friends": [
						{
							"name": "Author 0"
						},
						{
							"name": "Author 1"
						},
						{
							"name": "Author 8"
						}
					]
				}
			},
			{
				"id": "32",
				"title": "Post 31",
				"characters": 310,
				"score": 2.5,
				"completed": false,
				"author": {
					"name": "Author 10",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 7"
						}
					]
				}
			},
			{
				"id": "44",
				"title": "Post 43",
				"characters": 430,
				"score": 1.6,
				"completed": false,
				"author": {
					"name": "Author 12",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "58",
				"title": "Post 57",
				"characters": 570,
				"score": 4.5,
				"completed": false,
				"author": {
					"name": "Author 8",
					"friends": [
						{
							"name": "Author 8"
						},
						{
							"name": "Author 17"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "63",
				"title": "Post 62",
				"characters": 620,
				"score": 4.9,
				"completed": false,
				"author": {
					"name": "Author 3",
					"friends": [
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "71",
				"title": "Post 70",
				"characters": 700,
				"score": 1.9,
				"completed": false,
				"author": {
					"name": "Author 13",
					"friends": [
						{
							"name": "Author 15"
						},
						{
							"name": "Author 2"
						},
						{
							"name": "Author 4"
						}
					]
				}
			},
			{
				"id": "9",
				"title": "Post 8",
				"characters": 80,
				"score": 2.3,
				"completed": true,
				"author": {
					"name": "Author 5",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						}
					]
				}
			},
			{
				"id": "26",
				"title": "Post 25",
				"characters": 250,
				"score": 4.7,
				"completed": false,
				"author": {
					"name": "Author 5",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						}
					]
				}
			},
			{
				"id": "82",
				"title": "Post 81",
				"characters": 810,
				"score": 3.1,
				"completed": true,
				"author": {
					"name": "Author 16",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 6"
						}
					]
				}
			},
			{
				"id": "85",
				"title": "Post 84",
				"characters": 840,
				"score": 2.6,
				"completed": false,
				"author": {
					"name": "Author 12",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "10",
				"title": "Post 9",
				"characters": 90,
				"score": 1.3,
				"completed": true,
				"author": {
					"name": "Author 7",
					"friends": [
						{
							"name": "Author 0"
						},
						{
							"name": "Author 1"
						},
						{
							"name": "Author 8"
						}
					]
				}
			},
			{
				"id": "12",
				"title": "Post 11",
				"characters": 110,
				"score": 1.5,
				"completed": false,
				"author": {
					"name": "Author 4",
					"friends": [
						{
							"name": ""
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 17"
						}
					]
				}
			},
			{
				"id": "20",
				"title": "Post 19",
				"characters": 190,
				"score": 4,
				"completed": true,
				"author": {
					"name": "Author 1",
					"friends": [
						{
							"name": "Author 8"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "34",
				"title": "Post 33",
				"characters": 330,
				"score": 2.8,
				"completed": false,
				"author": {
					"name": "Author 17",
					"friends": [
						{
							"name": "Author 18"
						},
						{
							"name": "Author 4"
						},
						{
							"name": "Author 7"
						}
					]
				}
			},
			{
				"id": "45",
				"title": "Post 44",
				"characters": 440,
				"score": 4.5,
				"completed": true,
				"author": {
					"name": "",
					"friends": null
				}
			},
			{
				"id": "48",
				"title": "Post 47",
				"characters": 470,
				"score": 3.2,
				"completed": true,
				"author": {
					"name": "Author 14",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 9"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "8",
				"title": "Post 7",
				"characters": 70,
				"score": 3.2,
				"completed": false,
				"author": {
					"name": "Author 9",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "33",
				"title": "Post 32",
				"characters": 320,
				"score": 1.5,
				"completed": false,
				"author": {
					"name": "Author 1",
					"friends": [
						{
							"name": "Author 8"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "43",
				"title": "Post 42",
				"characters": 420,
				"score": 4.2,
				"completed": false,
				"author": {
					"name": "Author 13",
					"friends": [
						{
							"name": "Author 15"
						},
						{
							"name": "Author 2"
						},
						{
							"name": "Author 4"
						}
					]
				}
			},
			{
				"id": "56",
				"title": "Post 55",
				"characters": 550,
				"score": 4.2,
				"completed": true,
				"author": {
					"name": "Author 0",
					"friends": [
						{
							"name": "Author 13"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 2"
						}
					]
				}
			},
			{
				"id": "84",
				"title": "Post 83",
				"characters": 830,
				"score": 3,
				"completed": true,
				"author": {
					"name": "Author 7",
					"friends": [
						{
							"name": "Author 0"
						},
						{
							"name": "Author 1"
						},
						{
							"name": "Author 8"
						}
					]
				}
			},
			{
				"id": "30",
				"title": "Post 29",
				"characters": 290,
				"score": 2.4,
				"completed": true,
				"author": {
					"name": "Author 7",
					"friends": [
						{
							"name": "Author 0"
						},
						{
							"name": "Author 1"
						},
						{
							"name": "Author 8"
						}
					]
				}
			},
			{
				"id": "52",
				"title": "Post 51",
				"characters": 510,
				"score": 3.7,
				"completed": true,
				"author": {
					"name": "Author 14",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 9"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "95",
				"title": "Post 94",
				"characters": 940,
				"score": 2.9,
				"completed": true,
				"author": {
					"name": "Author 3",
					"friends": [
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "13",
				"title": "Post 12",
				"characters": 120,
				"score": 4.4,
				"completed": true,
				"author": {
					"name": "Author 12",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "27",
				"title": "Post 26",
				"characters": 260,
				"score": 3.2,
				"completed": false,
				"author": {
					"name": "Author 16",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 6"
						}
					]
				}
			},
			{
				"id": "40",
				"title": "Post 39",
				"characters": 390,
				"score": 1.1,
				"completed": false,
				"author": {
					"name": "Author 0",
					"friends": [
						{
							"name": "Author 13"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 2"
						}
					]
				}
			},
			{
				"id": "64",
				"title": "Post 63",
				"characters": 630,
				"score": 2.3,
				"completed": false,
				"author": {
					"name": "Author 15",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 6"
						},
						{
							"name": "Author 8"
						}
					]
				}
			},
			{
				"id": "89",
				"title": "Post 88",
				"characters": 880,
				"score": 3.4,
				"completed": false,
				"author": {
					"name": "Author 8",
					"friends": [
						{
							"name": "Author 8"
						},
						{
							"name": "Author 17"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "90",
				"title": "Post 89",
				"characters": 890,
				"score": 3,
				"completed": false,
				"author": {
					"name": "Author 0",
					"friends": [
						{
							"name": "Author 13"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 2"
						}
					]
				}
			},
			{
				"id": "11",
				"title": "Post 10",
				"characters": 100,
				"score": 4.2,
				"completed": true,
				"author": {
					"name": "Author 5",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						}
					]
				}
			},
			{
				"id": "15",
				"title": "Post 14",
				"characters": 140,
				"score": 4.6,
				"completed": true,
				"author": {
					"name": "Author 2",
					"friends": [
						{
							"name": "Author 12"
						},
						{
							"name": ""
						},
						{
							"name": "Author 6"
						}
					]
				}
			},
			{
				"id": "69",
				"title": "Post 68",
				"characters": 680,
				"score": 1.7,
				"completed": false,
				"author": {
					"name": "Author 18",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 4"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "73",
				"title": "Post 72",
				"characters": 720,
				"score": 4,
				"completed": false,
				"author": {
					"name": "",
					"friends": null
				}
			},
			{
				"id": "77",
				"title": "Post 76",
				"characters": 760,
				"score": 2.7,
				"completed": true,
				"author": {
					"name": "Author 5",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						}
					]
				}
			},
			{
				"id": "96",
				"title": "Post 95",
				"characters": 950,
				"score": 2.6,
				"completed": true,
				"author": {
					"name": "Author 13",
					"friends": [
						{
							"name": "Author 15"
						},
						{
							"name": "Author 2"
						},
						{
							"name": "Author 4"
						}
					]
				}
			},
			{
				"id": "97",
				"title": "Post 96",
				"characters": 960,
				"score": 3.3,
				"completed": true,
				"author": {
					"name": "Author 3",
					"friends": [
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "4",
				"title": "Post 3",
				"characters": 30,
				"score": 2.6,
				"completed": true,
				"author": {
					"name": "Author 18",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 4"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "18",
				"title": "Post 17",
				"characters": 170,
				"score": 3.7,
				"completed": true,
				"author": {
					"name": "",
					"friends": null
				}
			},
			{
				"id": "42",
				"title": "Post 41",
				"characters": 410,
				"score": 1.4,
				"completed": false,
				"author": {
					"name": "Author 10",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 7"
						}
					]
				}
			},
			{
				"id": "53",
				"title": "Post 52",
				"characters": 520,
				"score": 2.2,
				"completed": true,
				"author": {
					"name": "Author 10",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 7"
						}
					]
				}
			},
			{
				"id": "55",
				"title": "Post 54",
				"characters": 540,
				"score": 2.4,
				"completed": false,
				"author": {
					"name": "Author 8",
					"friends": [
						{
							"name": "Author 8"
						},
						{
							"name": "Author 17"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "78",
				"title": "Post 77",
				"characters": 770,
				"score": 2.4,
				"completed": true,
				"author": {
					"name": "Author 9",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 14"
						}
					]
				}
			},
			{
				"id": "93",
				"title": "Post 92",
				"characters": 920,
				"score": 2.5,
				"completed": false,
				"author": {
					"name": "Author 12",
					"friends": [
						{
							"name": "Author 7"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "2",
				"title": "Post 1",
				"characters": 10,
				"score": 2.4,
				"completed": true,
				"author": {
					"name": "Author 11",
					"friends": [
						{
							"name": "Author 6"
						},
						{
							"name": "Author 5"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "35",
				"title": "Post 34",
				"characters": 340,
				"score": 1.8,
				"completed": false,
				"author": {
					"name": "Author 6",
					"friends": [
						{
							"name": "Author 4"
						},
						{
							"name": "Author 15"
						},
						{
							"name": "Author 1"
						}
					]
				}
			},
			{
				"id": "38",
				"title": "Post 37",
				"characters": 370,
				"score": 4.9,
				"completed": true,
				"author": {
					"name": "Author 16",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 6"
						}
					]
				}
			},
			{
				"id": "39",
				"title": "Post 38",
				"characters": 380,
				"score": 2.2,
				"completed": true,
				"author": {
					"name": "Author 11",
					"friends": [
						{
							"name": "Author 6"
						},
						{
							"name": "Author 5"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "46",
				"title": "Post 45",
				"characters": 450,
				"score": 1.9,
				"completed": true,
				"author": {
					"name": "Author 2",
					"friends": [
						{
							"name": "Author 12"
						},
						{
							"name": ""
						},
						{
							"name": "Author 6"
						}
					]
				}
			},
			{
				"id": "50",
				"title": "Post 49",
				"characters": 490,
				"score": 3.4,
				"completed": false,
				"author": {
					"name": "Author 11",
					"friends": [
						{
							"name": "Author 6"
						},
						{
							"name": "Author 5"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "25",
				"title": "Post 24",
				"characters": 240,
				"score": 3,
				"completed": false,
				"author": {
					"name": "Author 19",
					"friends": [
						{
							"name": "Author 2"
						},
						{
							"name": "Author 17"
						},
						{
							"name": ""
						}
					]
				}
			}
		]
	}
}
```

### query existing post

```gql
query {
  getPost ( postID: "21") {
		title
	}
}
```

response:
```json
{
	"data": {
		"getPost": {
			"title": "Post 20"
		}
	}
}
```

### query nonexisting post

```gql
query {
  getPost ( postID: "001") {
		title
	}
}
```

response:
```json
{
	"errors": [
		{
			"message": "No posts found with id",
			"path": [
				"getPost"
			]
		}
	],
	"data": {
		"getPost": null
	}
}
```

### query filter by characters

```gql
query {
  getPosts (filter: { characters : {gte: 950}}) {
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
				"id": "97",
				"title": "Post 96",
				"characters": 960,
				"score": 4.6,
				"completed": true,
				"author": {
					"name": "Author 2",
					"friends": [
						{
							"name": "Author 3"
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 16"
						}
					]
				}
			},
			{
				"id": "100",
				"title": "Post 99",
				"characters": 990,
				"score": 5,
				"completed": true,
				"author": {
					"name": "Author 9",
					"friends": [
						{
							"name": "Author 11"
						},
						{
							"name": "Author 2"
						},
						{
							"name": "Author 5"
						}
					]
				}
			},
			{
				"id": "98",
				"title": "Post 97",
				"characters": 970,
				"score": 2,
				"completed": true,
				"author": {
					"name": "Author 3",
					"friends": [
						{
							"name": "Author 12"
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 13"
						}
					]
				}
			},
			{
				"id": "96",
				"title": "Post 95",
				"characters": 950,
				"score": 3,
				"completed": false,
				"author": {
					"name": "Author 10",
					"friends": [
						{
							"name": "Author 12"
						},
						{
							"name": "Author 6"
						},
						{
							"name": "Author 4"
						}
					]
				}
			},
			{
				"id": "99",
				"title": "Post 98",
				"characters": 980,
				"score": 3.4,
				"completed": false,
				"author": {
					"name": "Author 12",
					"friends": [
						{
							"name": "Author 15"
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 1"
						}
					]
				}
			}
		]
	}
}
```

### query filter by title

```gql
query {
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
```

response:
```json
{
	"data": {
		"getPosts": [
			{
				"id": "18",
				"title": "Post 17",
				"characters": 170,
				"score": 4,
				"completed": false,
				"author": {
					"name": "Author 15",
					"friends": [
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 7"
						}
					]
				}
			},
			{
				"id": "19",
				"title": "Post 18",
				"characters": 180,
				"score": 4.4,
				"completed": true,
				"author": {
					"name": "Author 7",
					"friends": [
						{
							"name": "Author 15"
						},
						{
							"name": "Author 2"
						},
						{
							"name": "Author 13"
						}
					]
				}
			},
			{
				"id": "62",
				"title": "Post 61",
				"characters": 610,
				"score": 3.1,
				"completed": false,
				"author": {
					"name": "Author 12",
					"friends": [
						{
							"name": "Author 15"
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 1"
						}
					]
				}
			},
			{
				"id": "82",
				"title": "Post 81",
				"characters": 810,
				"score": 3.6,
				"completed": true,
				"author": {
					"name": "Author 13",
					"friends": [
						{
							"name": "Author 2"
						},
						{
							"name": "Author 17"
						},
						{
							"name": "Author 8"
						}
					]
				}
			},
			{
				"id": "17",
				"title": "Post 16",
				"characters": 160,
				"score": 2.5,
				"completed": true,
				"author": {
					"name": "Author 3",
					"friends": [
						{
							"name": "Author 12"
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 13"
						}
					]
				}
			},
			{
				"id": "12",
				"title": "Post 11",
				"characters": 110,
				"score": 2.7,
				"completed": true,
				"author": {
					"name": "Author 17",
					"friends": [
						{
							"name": "Author 2"
						},
						{
							"name": "Author 6"
						},
						{
							"name": "Author 3"
						}
					]
				}
			},
			{
				"id": "11",
				"title": "Post 10",
				"characters": 100,
				"score": 1.2,
				"completed": false,
				"author": {
					"name": "Author 0",
					"friends": [
						{
							"name": "Author 11"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 2"
						}
					]
				}
			},
			{
				"id": "20",
				"title": "Post 19",
				"characters": 190,
				"score": 2,
				"completed": false,
				"author": {
					"name": "Author 2",
					"friends": [
						{
							"name": "Author 3"
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 16"
						}
					]
				}
			},
        ]
    }	
}
```


### query filter not   

```gql
query {
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
				"id": "2",
				"title": "Post 1",
				"characters": 10,
				"score": 1.2,
				"completed": true,
				"author": {
					"name": "Author 15",
					"friends": [
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 7"
						}
					]
				}
			},
			{
				"id": "5",
				"title": "Post 4",
				"characters": 40,
				"score": 4.9,
				"completed": true,
				"author": {
					"name": "Author 6",
					"friends": [
						{
							"name": "Author 10"
						},
						{
							"name": "Author 2"
						},
						{
							"name": "Author 11"
						}
					]
				}
			},
			{
				"id": "4",
				"title": "Post 3",
				"characters": 30,
				"score": 3.1,
				"completed": true,
				"author": {
					"name": "Author 10",
					"friends": [
						{
							"name": "Author 12"
						},
						{
							"name": "Author 6"
						},
						{
							"name": "Author 4"
						}
					]
				}
			},
			{
				"id": "3",
				"title": "Post 2",
				"characters": 20,
				"score": 1.2,
				"completed": false,
				"author": {
					"name": "Author 11",
					"friends": [
						{
							"name": "Author 4"
						},
						{
							"name": "Author 14"
						},
						{
							"name": "Author 12"
						}
					]
				}
			},
			{
				"id": "1",
				"title": "Post 0",
				"characters": 0,
				"score": 3.3,
				"completed": true,
				"author": {
					"name": "Author 8",
					"friends": [
						{
							"name": ""
						},
						{
							"name": "Author 10"
						},
						{
							"name": "Author 12"
						}
					]
				}
			}
		]
	}
}
```


### query filter by and

```gql
query {
  getPosts (
		filter: {
			and: [
			 	{ title: {contains: "1"}},
				{ isComplete: false},
				{ characters : {lt: 200}}
				
			]
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
{
	"data": {
		"getPosts": [
			{
				"id": "11",
				"title": "Post 10",
				"characters": 100,
				"score": 1.2,
				"completed": false,
				"author": {
					"name": "Author 0",
					"friends": [
						{
							"name": "Author 11"
						},
						{
							"name": "Author 16"
						},
						{
							"name": "Author 2"
						}
					]
				}
			},
			{
				"id": "20",
				"title": "Post 19",
				"characters": 190,
				"score": 2,
				"completed": false,
				"author": {
					"name": "Author 2",
					"friends": [
						{
							"name": "Author 3"
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 16"
						}
					]
				}
			},
			{
				"id": "15",
				"title": "Post 14",
				"characters": 140,
				"score": 3.1,
				"completed": false,
				"author": {
					"name": "Author 5",
					"friends": [
						{
							"name": "Author 2"
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 8"
						}
					]
				}
			},
			{
				"id": "16",
				"title": "Post 15",
				"characters": 150,
				"score": 4.6,
				"completed": false,
				"author": {
					"name": "Author 12",
					"friends": [
						{
							"name": "Author 15"
						},
						{
							"name": "Author 12"
						},
						{
							"name": "Author 1"
						}
					]
				}
			},
			{
				"id": "18",
				"title": "Post 17",
				"characters": 170,
				"score": 4,
				"completed": false,
				"author": {
					"name": "Author 15",
					"friends": [
						{
							"name": "Author 3"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 7"
						}
					]
				}
			}
		]
	}
}
```


### query filter by or

```gql
query {
  getPosts (
		filter: {
			or: [
			 	{ title: {contains: "11"}},
				{ characters : {lt: 3}}
				
			]
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
{
	"data": {
		"getPosts": [
			{
				"id": "12",
				"title": "Post 11",
				"characters": 110,
				"score": 4.7,
				"completed": true,
				"author": {
					"name": "Author 16",
					"friends": [
						{
							"name": "Author 4"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 18"
						}
					]
				}
			},
			{
				"id": "1",
				"title": "Post 0",
				"characters": 0,
				"score": 2.7,
				"completed": false,
				"author": {
					"name": "Author 15",
					"friends": [
						{
							"name": "Author 17"
						},
						{
							"name": ""
						},
						{
							"name": "Author 1"
						}
					]
				}
			}
		]
	}
}
```

### query with sorted order

```gql
query {
  getPosts (
		filter: {
			or: [
			 	{ title: {contains: "11"}},
				{ characters : {lt: 3}}
				
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
```
{
	"data": {
		"getPosts": [
			{
				"id": "1",
				"title": "Post 0",
				"characters": 0,
				"score": 2.7,
				"completed": false,
				"author": {
					"name": "Author 15",
					"friends": [
						{
							"name": "Author 17"
						},
						{
							"name": ""
						},
						{
							"name": "Author 1"
						}
					]
				}
			},
			{
				"id": "12",
				"title": "Post 11",
				"characters": 110,
				"score": 4.7,
				"completed": true,
				"author": {
					"name": "Author 16",
					"friends": [
						{
							"name": "Author 4"
						},
						{
							"name": "Author 13"
						},
						{
							"name": "Author 18"
						}
					]
				}
			}
		]
	}
}
```

### query filter by aaggregatePost

```gql
query {
  aggregatePost (
		filter: {
			or: [
			 	{ title: {contains: "11"}},
				{ characters : {lt: 3}}
				
			]
		}
	
	) {
		count
		posts {
			title
			text
			completed
			characters
		}
	}
}
```

response:
```json
{
	"data": {
		"aggregatePost": {
			"count": "2",
			"posts": [
				{
					"title": "Post 11",
					"text": null,
					"completed": true,
					"characters": 110
				},
				{
					"title": "Post 0",
					"text": null,
					"completed": true,
					"characters": 0
				}
			]
		}
	}
}
```

## graphql to database queries

### Create MySQL database
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

## G0-Jet queries


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
	// code with mockdata in resolver r
	// if post, ok := r.Posts[postID]; ok {
	// 	return &post, nil
	// }
	// return nil, fmt.Errorf("No posts found with id")

	var resultPost *model.Post
	resultPost, err := database.GetPost(postID)
	if err != nil {
		return nil, err
	}
	return resultPost, nil

}

// GetPosts is the resolver for the getPosts field.
func (r *queryResolver) GetPosts(ctx context.Context, filter *model.PostFilter, order *model.PostOrder) ([]*model.Post, error) {
	// code with mockdata in resolver r
	// var posts []*model.Post
	// var filteredPosts []*model.Post
	// var sortedPosts []*model.Post
	// var err error
	// //get all posts
	// posts = make([]*model.Post, 0, len(r.Posts))
	// for _, value := range r.Posts {
	// 	posts = append(posts, &value)
	// }
	// //filter posts
	// filteredPosts, err = filterPostHelper(filter, posts)
	// if err != nil {
	// 	return nil, err
	// }
	// //sort filtered posts
	// sortedPosts, err = sortPostHelper(order, filteredPosts)
	// if err != nil {
	// 	return nil, err
	// }
	// return sortedPosts, nil

	resultPosts, err := database.GetPosts(filter, order)
	if err != nil {
		return nil, err
	}
	return resultPosts, nil
}

// AggregatePost is the resolver for the aggregatePost field.
func (r *queryResolver) AggregatePost(ctx context.Context, filter *model.PostFilter) (*model.PostAggregateResult, error) {
	// code with mockdata in resolver r
	// //filter the posts
	// filteredPosts, err := r.GetPosts(ctx, filter, nil)
	// if filteredPosts != nil && err == nil {
	// 	//aggregate the posts
	// 	count := len(filteredPosts)
	// 	countStr := strconv.Itoa(count)
	// 	total := 0
	// 	for _, post := range filteredPosts {
	// 		total += int(*post.Score)
	// 	}
	// 	average := float64(total) / float64(count)
	// 	return &model.PostAggregateResult{
	// 		Count:    &countStr,
	// 		AvgScore: &average,
	// 		Posts:    filteredPosts,
	// 	}, nil
	// }
	// return nil, err

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
		//println("whereClauseSqlExpression=", whereClauseSqlExpression)
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


