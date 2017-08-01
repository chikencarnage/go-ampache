package ampache

type ArtistResults struct {
	Count   int
	Artists []Artist `xml:"artist"`
}

type Artist struct {
	ID            string  `xml:"id,attr"`
	Name          string  `xml:"name"`
	Albums        int     `xml:"albums"`
	Songs         int     `xml:"songs"`
	Tags          []Tag   `xml:"tag"`
	PreciseRating int     `xml:"preciserating"`
	Rating        float32 `xml:"rating"`
}

type ArtistName struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",chardata"`
}

func (s *Search) SearchArtists() (*ArtistResults, error) {
	return &ArtistResults{}, nil
}

func (s *Search) SearchArtist() (*ArtistResults, error) {
	return &ArtistResults{}, nil
}

func (s *Search) SearchArtistAlbums() (*AlbumResults, error) {
	return &AlbumResults{}, nil
}

func (s *Search) SearchArtistSongs() (*SongResults, error) {
	return &SongResults{}, nil
}
