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

		Friends []struct {
			dbModel.Author
		}
	}
}

type GetPostsDest []struct {
	dbModel.Post

	Authors []struct {
		dbModel.Author

		Friends []struct {
			dbModel.Author
		}
	}
}

type GetPostAggregateDest struct {
	Posts []struct {
		dbModel.Post

		Authors []struct {
			dbModel.Author

			Friends []struct {
				dbModel.Author
			}
		}
	}

	Count int32

	AvgScore float32
}

//SQL query to us
// SELECT friends.*
// from Author friends
// inner join Author_Friend_Author afa on afa.friend_author_id = friends.ID
// inner join Author author on author.ID=afa.author_id
// where author.ID =1
// UNION
// SELECT friends.*
// from Author friends
// inner join Author_Friend_Author afa on afa.author_id = friends.ID
// inner join Author author on author.ID=afa.friend_author_id
// where author.ID =1
//output:
//|ID|name|
// |--|----|
// |2|Author 2|
// |3|Author 3|
// |4|Author 4|
// |10|Author 10|

// debug sql: generates correct results. but not this program. i ssuem it is GetPostDest structure. Did not fix it here.
// SELECT `Post`.`ID` AS "Post.ID",
//
//	`Post`.title AS "Post.title",
//	`Post`.characters AS "Post.characters",
//	`Post`.completed AS "Post.completed",
//	`Post`.`datePublished` AS "Post.datePublished",
//	`Post`.score AS "Post.score",
//	`Post`.text AS "Post.text",
//	`Author`.`ID` AS "Author.ID",
//	`Author`.name AS "Author.name",
//	friends.`Friend.ID` AS "Friend.ID",
//	friends.`Friend.name` AS "Friend.name"
//
// FROM posts.`Post`
//
//	INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
//	INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
//	INNER JOIN (
//	     (
//	          SELECT `Friend`.`ID` AS "Friend.ID",
//	               `Friend`.name AS "Friend.name"
//	          FROM posts.`Author` AS `Friend`
//	               INNER JOIN posts.`Author_Friend_Author` ON (`Author_Friend_Author`.friend_author_id = `Friend`.`ID`)
//	               INNER JOIN posts.`Author` ON (`Author_Friend_Author`.author_id = `Author`.`ID`)
//	     )
//	     UNION
//	     (
//	          SELECT `Friend`.`ID` AS "Friend.ID",
//	               `Friend`.name AS "Friend.name"
//	          FROM posts.`Author` AS `Friend`
//	               INNER JOIN posts.`Author_Friend_Author` ON (`Author_Friend_Author`.author_id = `Friend`.`ID`)
//	               INNER JOIN posts.`Author` ON (`Author_Friend_Author`.friend_author_id = `Author`.`ID`)
//	     )
//	) AS friends ON (`Friend.ID` = `Author`.`ID`)
//
// WHERE `Post`.`ID` = 1;
func GetPost(postID string) (*graphModel.Post, error) {
	//get postId as int64
	postIdVal, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	//query
	FRIEND := dbTable.Author.AS("Friend")

	friendsSubQuery :=
		jetMysql.SELECT(
			dbTable.Author.ID.AS("Author.ID"),
			FRIEND.ID.AS("Friend.ID"),
			FRIEND.Name.AS("Friend.name"),
		).FROM(FRIEND.
			INNER_JOIN(dbTable.AuthorFriendAuthor, dbTable.AuthorFriendAuthor.FriendAuthorID.EQ(FRIEND.ID)).
			INNER_JOIN(dbTable.Author, dbTable.AuthorFriendAuthor.AuthorID.EQ(dbTable.Author.ID)),
		).UNION(
			jetMysql.SELECT(
				dbTable.Author.ID, FRIEND.AllColumns,
			).FROM(FRIEND.
				INNER_JOIN(dbTable.AuthorFriendAuthor, dbTable.AuthorFriendAuthor.AuthorID.EQ(FRIEND.ID)).
				INNER_JOIN(dbTable.Author, dbTable.AuthorFriendAuthor.FriendAuthorID.EQ(dbTable.Author.ID)),
			),
		).AsTable("friends")

	//https://github.com/go-jet/jet/wiki/Subquery
	//select columns from friendsSubQuery to use in main
	authorId := dbTable.Author.ID.From(friendsSubQuery)
	friendId := jetMysql.IntegerColumn("Friend.ID").From(friendsSubQuery)
	friendName := jetMysql.StringColumn("Friend.name").From(friendsSubQuery)

	var projectionList jetMysql.ProjectionList
	projectionList = append(projectionList, dbTable.Post.ID)
	projectionList = append(projectionList, dbTable.Post.Title)
	projectionList = append(projectionList, dbTable.Post.Characters)
	projectionList = append(projectionList, dbTable.Post.Completed)
	projectionList = append(projectionList, dbTable.Post.DatePublished)
	projectionList = append(projectionList, dbTable.Post.Score)
	projectionList = append(projectionList, dbTable.Post.Text)
	projectionList = append(projectionList, dbTable.Author.AllColumns)
	projectionList = append(projectionList, friendId)
	projectionList = append(projectionList, friendName)

	query :=
		jetMysql.
			SELECT(projectionList).
			FROM(
				dbTable.Post.
					INNER_JOIN(dbTable.AuthorPost, dbTable.Post.ID.EQ(dbTable.AuthorPost.PostID)).
					INNER_JOIN(dbTable.Author, dbTable.Author.ID.EQ(dbTable.AuthorPost.AuthorID)).
					INNER_JOIN(friendsSubQuery, authorId.EQ(dbTable.Author.ID)),
			).
			WHERE(
				dbTable.Post.ID.EQ(jetMysql.Int(postIdVal)),
			)

	printStatementInfo(query)

	//Store result into desired destination:
	var dest GetPostDest

	err = query.Query(Db, &dest)
	if err != nil {
		return nil, err
	}

	jsonSave("./internal/db/out/resultPost.json", dest)

	datePub := dest.DatePublished.Format("2006-01-02")
	id := strconv.FormatInt(int64(dest.ID), 10)

	var authors []*graphModel.Author
	if len(dest.Authors) > 0 {
		for _, authorVal := range dest.Authors {

			//get friends
			var friends []*model.Author
			for _, friendVal := range authorVal.Friends {
				friendId := strconv.Itoa(int(friendVal.ID))
				friendAuthor := &graphModel.Author{
					ID:   friendId,
					Name: friendVal.Name,
				}
				friends = append(friends, friendAuthor)
			}

			author := &graphModel.Author{
				ID:      strconv.FormatInt(int64(authorVal.ID), 10),
				Name:    authorVal.Name,
				Friends: friends,
			}
			authors = append(authors, author)

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
		Authors:       authors,
	}
	return resultPost, nil
}

func GetPosts(filter *graphModel.PostFilter, order *graphModel.PostOrder) ([]*model.Post, error) {
	//see https://github.com/go-jet/jet/wiki/SELECT#table-aliasing if you have relation in the same table like author-firend-author, manager-employee to write sql

	var orderByClause jetMysql.OrderByClause
	var whereClauseSqlExpression jetMysql.BoolExpression
	var err error

	//https://github.com/go-jet/jet/wiki/Subquery
	FRIEND := dbTable.Author.AS("Friend")

	friendsSubQuery :=
		jetMysql.SELECT(
			dbTable.Author.ID, FRIEND.AllColumns,
		).FROM(FRIEND.
			INNER_JOIN(dbTable.AuthorFriendAuthor, dbTable.AuthorFriendAuthor.FriendAuthorID.EQ(FRIEND.ID)).
			INNER_JOIN(dbTable.Author, dbTable.AuthorFriendAuthor.AuthorID.EQ(dbTable.Author.ID)),
		).UNION(
			jetMysql.SELECT(
				dbTable.Author.ID, FRIEND.AllColumns,
			).FROM(FRIEND.
				INNER_JOIN(dbTable.AuthorFriendAuthor, dbTable.AuthorFriendAuthor.AuthorID.EQ(FRIEND.ID)).
				INNER_JOIN(dbTable.Author, dbTable.AuthorFriendAuthor.FriendAuthorID.EQ(dbTable.Author.ID)),
			),
		).AsTable("friends")

	authorId := dbTable.Author.ID.From(friendsSubQuery)

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
	projectionList = append(projectionList, friendsSubQuery.AllColumns())

	query :=
		jetMysql.
			SELECT(projectionList).
			FROM(
				dbTable.Post.
					INNER_JOIN(dbTable.AuthorPost, dbTable.Post.ID.EQ(dbTable.AuthorPost.PostID)).
					INNER_JOIN(dbTable.Author, dbTable.Author.ID.EQ(dbTable.AuthorPost.AuthorID)).
					INNER_JOIN(friendsSubQuery, authorId.EQ(dbTable.Author.ID)),
			)

	printStatementInfo(query)
	//	Debug sql: This provides correct results in dbeaver, not here though
	//
	// 	SELECT `Post`.`ID` AS "Post.ID",
	// 	`Post`.title AS "Post.title",
	// 	`Post`.characters AS "Post.characters",
	// 	`Post`.completed AS "Post.completed",
	// 	`Post`.`datePublished` AS "Post.datePublished",
	// 	`Post`.score AS "Post.score",
	// 	`Post`.text AS "Post.text",
	// 	`Author`.`ID` AS "Author.ID",
	// 	`Author`.name AS "Author.name",
	// 	friends.`Author.ID` AS "Author.ID",
	// 	friends.`Friend.ID` AS "Friend.ID",
	// 	friends.`Friend.name` AS "Friend.name"
	// FROM posts.`Post`
	// 	INNER JOIN posts.`Author_Post` ON (`Post`.`ID` = `Author_Post`.post_id)
	// 	INNER JOIN posts.`Author` ON (`Author`.`ID` = `Author_Post`.author_id)
	// 	INNER JOIN (
	// 		 (
	// 			  SELECT `Author`.`ID` AS "Author.ID",
	// 				   `Friend`.`ID` AS "Friend.ID",
	// 				   `Friend`.name AS "Friend.name"
	// 			  FROM posts.`Author` AS `Friend`
	// 				   INNER JOIN posts.`Author_Friend_Author` ON (`Author_Friend_Author`.friend_author_id = `Friend`.`ID`)
	// 				   INNER JOIN posts.`Author` ON (`Author_Friend_Author`.author_id = `Author`.`ID`)
	// 		 )
	// 		 UNION
	// 		 (
	// 			  SELECT `Author`.`ID` AS "Author.ID",
	// 				   `Friend`.`ID` AS "Friend.ID",
	// 				   `Friend`.name AS "Friend.name"
	// 			  FROM posts.`Author` AS `Friend`
	// 				   INNER JOIN posts.`Author_Friend_Author` ON (`Author_Friend_Author`.author_id = `Friend`.`ID`)
	// 				   INNER JOIN posts.`Author` ON (`Author_Friend_Author`.friend_author_id = `Author`.`ID`)
	// 		 )
	// 	) AS friends ON (friends.`Author.ID` = `Author`.`ID`)
	// WHERE `Post`.`ID` = 1;

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

		var authors []*graphModel.Author
		if len(post.Authors) > 0 {
			if len(post.Authors) > 0 {
				for _, val := range post.Authors {
					author := &graphModel.Author{
						ID:   strconv.FormatInt(int64(val.ID), 10),
						Name: val.Name,
					}
					authors = append(authors, author)
				}
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
			Authors:       authors,
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
		return nil, fmt.Errorf("unsupported sort field")
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
