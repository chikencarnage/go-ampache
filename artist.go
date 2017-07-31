package ampache

type ArtistName struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",chardata"`
}
