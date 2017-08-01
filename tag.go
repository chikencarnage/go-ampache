package ampache

type TagResult struct {
	Tag []Tag `xml:"tag"`
}

type Tag struct {
	ID       string `xml:"id,attr"`
	Name     string `xml:"name"`
	Albums   int    `xml:"albums"`
	Artists  int    `xml:"artists"`
	Songs    int    `xml:"songs"`
	Video    int    `xml:"video"`
	Playlist int    `xml:"playlist"`
	Stream   int    `xml:"stream"`
}

type TagInfo struct {
	ID    string `xml:"id,attr"`
	Count string `xml:"count,attr"`
	Value string `xml:",chardata"`
}
