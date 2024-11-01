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