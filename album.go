package ampache

type AlbumResults struct {
	Count  int
	Albums []Album `xml:"album"`
}

type Album struct {
	ID            string     `xml:"id, attr"`
	Name          string     `xml:"name"`
	Artist        ArtistName `xml:"artist"`
	Year          string     `xml:"year"`
	Tracks        int        `xml:"tracks"`
	Disk          int        `xml:"disk"`
	Tags          []Tag      `xml:"tag"`
	Art           string     `xml:"art"`
	PreciseRating int        `xml:"preciserating"`
	Rating        float32    `xml:"rating"`
}

type AlbumTitle struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",chardata"`
}
