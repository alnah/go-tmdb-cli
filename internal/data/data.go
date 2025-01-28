package data

import (
	"sort"
	"time"
)

type (
	Movies []Movie
	Movie  struct {
		Id            int     `json:"id"`
		Date          string  `json:"release_date"`
		OriginalTitle string  `json:"original_title"`
		Title         string  `json:"title"`
		Average       float64 `json:"vote_average"`
		Votes         int     `json:"vote_count"`
	}
)

type OrderError struct{}

func (e OrderError) Error() string {
	return `Sorting order must be "asc" or "desc".`
}

func (m Movies) Sort() (Movies, error) {
	sortOrder := []struct {
		sortFunc func(order string) (Movies, error)
		order    string
	}{
		{m.SortByOriginalTitle, "asc"},
		{m.SortByTitle, "asc"},
		{m.SortByDate, "asc"},
		{m.SortByVotes, "desc"},
		{m.SortByAverage, "desc"},
	}

	var err error
	for _, s := range sortOrder {
		m, err = s.sortFunc(s.order)
		if err != nil {
			return m, err
		}
	}

	if err != nil {
		return m, err
	}

	return m, nil
}

func (m Movies) SortByDate(order string) (Movies, error) {
	return m.sort(order, func(i int, j int) bool {
		iDate, _ := time.Parse(time.DateOnly, m[i].Date)
		jDate, _ := time.Parse(time.DateOnly, m[j].Date)

		return iDate.Before(jDate)
	})
}

func (m Movies) SortByOriginalTitle(order string) (Movies, error) {
	return m.sort(order, func(i int, j int) bool {
		return m[i].OriginalTitle < m[j].OriginalTitle
	})
}

func (m Movies) SortByTitle(order string) (Movies, error) {
	return m.sort(order, func(i int, j int) bool {
		return m[i].Title < m[j].Title
	})
}

func (m Movies) SortByAverage(order string) (Movies, error) {
	return m.sort(order, func(i int, j int) bool {
		return m[i].Average < m[j].Average
	})
}

func (m Movies) SortByVotes(order string) (Movies, error) {
	return m.sort(order, func(i int, j int) bool {
		return m[i].Votes < m[j].Votes
	})
}

func (m Movies) sort(order string, compare func(i, j int) bool) (Movies, error) {
	if err := checkOrder(order); err != nil {
		return m, err
	}

	sort.Slice(m, func(i, j int) bool {
		if order == "asc" {
			return compare(i, j)
		}

		return !compare(i, j)
	})

	return m, nil
}

func checkOrder(order string) error {
	if order != "asc" && order != "desc" {
		return &OrderError{}
	}

	return nil
}
