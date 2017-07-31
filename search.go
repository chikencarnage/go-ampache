package ampache

type Search struct {
	Filter      string
	Exact       bool
	SinceAdd    string
	SinceUpdate string
	Include     []string
}

func NewSearch(filter string) *Search {
	return &Search{Filter: filter}
}
