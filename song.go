package ampache

type SongResult struct {
	Songs []Song `xml:"song"`
}

type Song struct {
	ID            string     `xml:"id,attr"`
	Title         string     `xml:"title"`
	Artist        ArtistName `xml:"artist"`
	Album         AlbumTitle `xml:"album"`
	Tags          []Tag      `xml:"tag"`
	Track         int        `xml:"track"`
	Time          int        `xml:"time"`
	Url           string     `xml:"url"`
	Site          uint64     `xml:"site"`
	Art           string     `xml:"art"`
	PreciseRating int        `xml:"preciserating"`
	Rating        float32    `xml:"rating"`
}
