package graph

// always place any additoonal code needed for schema.resolvers.go in a seperate file i.e. here to stop being overwritten every time
// command go run github.com/99designs/gqlgen generate is run

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/dmrhimali/go-graphql-filtering/graph/model"
)

func sortPostHelper(order *model.PostOrder, postsToSort []*model.Post) ([]*model.Post, error) {
	if order == nil {
		return postsToSort, nil
	}

	if order.Field != nil {

		if order.Order == nil || *order.Order == model.SortOrderAsc { //ASC
			if *order.Field == model.SortableFieldCharacters {
				sort.Slice(postsToSort, func(i, j int) bool {
					return postsToSort[i].Characters < postsToSort[j].Characters
				})
			}
			if *order.Field == model.SortableFieldDatePublished {
				sort.Slice(postsToSort, func(i, j int) bool {
					return *postsToSort[i].DatePublished < *postsToSort[j].DatePublished
				})
			}
			if *order.Field == model.SortableFieldDatePublished {
				sort.Slice(postsToSort, func(i, j int) bool {
					return postsToSort[i].Title < postsToSort[j].Title
				})
			}
		} else { //DESC
			if *order.Field == model.SortableFieldCharacters {
				sort.Slice(postsToSort, func(i, j int) bool {
					return postsToSort[i].Characters > postsToSort[j].Characters
				})
			}
			if *order.Field == model.SortableFieldDatePublished {
				sort.Slice(postsToSort, func(i, j int) bool {
					return *postsToSort[i].DatePublished > *postsToSort[j].DatePublished
				})
			}
			if *order.Field == model.SortableFieldDatePublished {
				sort.Slice(postsToSort, func(i, j int) bool {
					return postsToSort[i].Title > postsToSort[j].Title
				})
			}
		}
		return postsToSort, nil
	}
	return nil, fmt.Errorf("no order field specified")

}
func filterPostHelper(filter *model.PostFilter, postsToFilter []*model.Post) ([]*model.Post, error) {
	var filteredPosts []*model.Post
	var err error

	filteredPosts = postsToFilter

	if filter == nil || IsEmptyStruct(filter) {
		return postsToFilter, nil
	}

	if filter.Title != nil {
		filteredPosts, err = filterByTitle(filteredPosts, filter)
	}

	if filter.Characters != nil {
		filteredPosts, err = filterByCharacters(filteredPosts, filter)
	}

	if filter.IsComplete != nil {
		filteredPosts, err = filterByIsComplete(filteredPosts, filter)
	}

	if filter.Not != nil {
		filteredPosts, err = filterByNot(filteredPosts, filter.Not)
	}

	if filter.And != nil {
		filteredPosts, err = filterByAnd(filteredPosts, filter.And)
	}

	if filter.Or != nil {
		filteredPosts, err = filterByOr(filteredPosts, filter.Or)
	}

	return filteredPosts, err

}
func filterByIsComplete(postsToFilter []*model.Post, filter *model.PostFilter) ([]*model.Post, error) {
	var filteredPosts []*model.Post
	if filter != nil {
		filteredPosts = make([]*model.Post, 0)
		for _, post := range postsToFilter {
			if post.Completed == *filter.IsComplete {
				filteredPosts = append(filteredPosts, post)
			}
		}
		return filteredPosts, nil
	}
	return nil, fmt.Errorf("filter iscomplete not specified")
}
func filterByCharacters(postsToFilter []*model.Post, filter *model.PostFilter) ([]*model.Post, error) {
	var filteredPosts []*model.Post

	if filter != nil {
		filteredPosts = make([]*model.Post, 0)
		for _, post := range postsToFilter {
			if filter.Characters != nil {
				if filter.Characters.Equals != nil && filter.Characters.Equals == &post.Characters {
					filteredPosts = append(filteredPosts, post)
				} else if filter.Characters.Gt != nil && *filter.Characters.Gt < post.Characters {
					filteredPosts = append(filteredPosts, post)
				} else if filter.Characters.Gte != nil && *filter.Characters.Gte <= post.Characters {
					filteredPosts = append(filteredPosts, post)
				} else if filter.Characters.Lt != nil && *filter.Characters.Lt > post.Characters {
					filteredPosts = append(filteredPosts, post)
				} else if filter.Characters.Lte != nil && *filter.Characters.Lte >= post.Characters {
					filteredPosts = append(filteredPosts, post)
				}
			}
		}
		return filteredPosts, nil
	}
	return nil, fmt.Errorf("filter character not specified")
}
func filterByTitle(postsToFilter []*model.Post, filter *model.PostFilter) ([]*model.Post, error) {
	var filteredPosts []*model.Post
	if filter != nil {
		filteredPosts = make([]*model.Post, 0)
		for _, post := range postsToFilter {
			if filter.Title != nil {
				if filter.Title.Contains != nil && strings.Contains(post.Title, *filter.Title.Contains) {
					filteredPosts = append(filteredPosts, post)
				} else if filter.Title.Equals != nil && *filter.Title.Equals == post.Title {
					filteredPosts = append(filteredPosts, post)
				}
			}
		}
		return filteredPosts, nil
	}
	return nil, fmt.Errorf("filter title not specified")
}
func filterByNot(postsToFilter []*model.Post, filter *model.PostFilter) ([]*model.Post, error) {
	//invert filter
	var filteredOpposite []*model.Post
	var filtered []*model.Post
	var err error

	filteredOpposite, err = filterPostHelper(filter, postsToFilter)
	filtered = difference(postsToFilter, filteredOpposite)
	return filtered, err
}

func filterByOr(postsToFilter []*model.Post, filters []*model.PostFilter) ([]*model.Post, error) {
	var filtered []*model.Post
	var err error

	for _, fieldFilter := range filters {
		var fieildfiltered []*model.Post
		fieildfiltered, err = filterPostHelper(fieldFilter, postsToFilter)
		filtered = union(filtered, fieildfiltered)
	}

	return filtered, err
}

// and.every(subFilter => applyFilter(user, subFilter));
func filterByAnd(postsToFilter []*model.Post, filters []*model.PostFilter) ([]*model.Post, error) {
	var filtered []*model.Post
	var err error

	filtered = postsToFilter

	for _, fieldFilter := range filters {
		filtered, err = filterPostHelper(fieldFilter, filtered)
	}

	return filtered, err
}

func difference(list1, list2 []*model.Post) []*model.Post {

	diff := make([]*model.Post, 0)
	for _, el1 := range list1 {
		found := false
		for _, el2 := range list2 {
			if el1 == el2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, el1)
		}
	}
	return diff
}

func union(slice1, slice2 []*model.Post) []*model.Post {
	set := make(map[*model.Post]bool)
	for _, v := range slice1 {
		set[v] = true
	}
	for _, v := range slice2 {
		set[v] = true
	}

	result := make([]*model.Post, 0)
	for k := range set {
		result = append(result, k)
	}

	return result
}

func intersection(slice1, slice2 []*model.Post) []*model.Post {
	set := make(map[*model.Post]bool)
	for _, v := range slice1 {
		set[v] = true
	}

	result := make([]*model.Post, 0)
	for _, v := range slice2 {
		if set[v] {
			result = append(result, v)
		}
	}

	return result
}

func IsEmptyStruct(s interface{}) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			return false
		}
	}

	return true
}

// for debugging printing structs
func printStructInJson(title string, val interface{}) {
	b, err := json.Marshal(val)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(title, ":\n", string(b))
}
