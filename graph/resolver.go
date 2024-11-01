package graph

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/dmrhimali/go-graphql-filtering/graph/model"
	"golang.org/x/exp/rand"
)

type Resolver struct {
	Authors map[string]model.Author
	Posts   map[string]model.Post
}

func NewResolver() Config {
	const nAuthors = 20
	const nPosts = 100
	r := Resolver{}
	r.Authors = make(map[string]model.Author, nAuthors)
	r.Posts = make(map[string]model.Post, nPosts)

	rand.Seed(uint64(time.Now().UnixNano()))

	for i := 0; i < nAuthors; i++ {
		authorId := strconv.Itoa(i + 1)
		mockAuthor := model.Author{
			ID:   authorId,
			Name: fmt.Sprintf("Author %d", i),
		}
		r.Authors[authorId] = mockAuthor

	}

	//set friends
	for i := 0; i < nAuthors; i++ {
		authorId := strconv.Itoa(i + 1)

		var friends []*model.Author
		friends = make([]*model.Author, 0)
		randIndices := pickRandomIndices(nAuthors, 3, i)

		for _, friendAuthorId := range randIndices {
			friendAuthor := r.Authors[strconv.Itoa(friendAuthorId)]
			friends = append(friends, &friendAuthor)
		}
		author := r.Authors[authorId]
		author.Friends = friends

		r.Authors[authorId] = author

	}

	for i := 0; i < nPosts; i++ {
		postId := strconv.Itoa(i + 1)
		assignedAuthorId := strconv.Itoa(RandIndex(len(r.Authors)))
		author := r.Authors[assignedAuthorId]
		mockPost := model.Post{
			ID:            postId,
			Title:         fmt.Sprintf("Post %d", i),
			Characters:    i * 10,
			Completed:     RandBool(),
			Score:         RandFloat(),
			DatePublished: RandDate(),
			Author:        &author,
		}
		r.Posts[postId] = mockPost

		author.Posts = append(author.Posts, &mockPost)
		r.Authors[assignedAuthorId] = author

	}

	return Config{
		Resolvers: &r,
	}

}

func RandBool() bool {
	rand.Seed(uint64(time.Now().UnixNano()))
	return rand.Intn(2) == 1
}

func RandIndex(nItems int) int {
	// Seed the random number generator
	rand.Seed(uint64(time.Now().UnixNano()))
	// Generate a random index within the range of the list
	randomIndex := rand.Intn(nItems)
	return randomIndex
}

// return a date of yyyy-mm-dd in string format
func RandDate() *string {
	// Seed the random number generator
	rand.Seed(uint64(time.Now().UnixNano()))

	// Generate a random time between two dates
	min := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	max := time.Date(2025, 12, 31, 23, 59, 59, 999999999, time.UTC)

	// Calculate the difference in nanoseconds
	delta := max.UnixNano() - min.UnixNano()

	// Generate a random duration within the range
	randomDuration := rand.Int63n(delta)

	// Add the random duration to the minimum date
	randomTime := min.Add(time.Duration(randomDuration))
	timeString := strconv.Quote(randomTime.Format("2006-01-02"))
	return &timeString
}

func RandFloat() *float64 {
	// Seed the random number generator
	rand.Seed(uint64(time.Now().UnixNano()))

	// Define the range
	min := 1.0
	max := 5.0

	// Generate a random float within the specified range
	randFloatInRange := min + rand.Float64()*(max-min)

	precision := 1 //1 decimal point
	randFloatInRageWithPrecision := fixPrecision(randFloatInRange, precision)
	return &randFloatInRageWithPrecision
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func fixPrecision(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func pickRandomIndices(length int, count int, excludeIndex int) []int {
	rand.Seed(uint64(time.Now().UnixNano())) // Seed the random number generator

	if count > length {
		count = length // Ensure we don't pick more indices than the slice has
	}

	indices := make([]int, count)
	i := 0
	for i < count {
		// Generate a random index within the slice's length
		randomIndex := rand.Intn(length)
		if randomIndex != excludeIndex {
			indices[i] = randomIndex
			i++
		}
	}

	return indices
}
