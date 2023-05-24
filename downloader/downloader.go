package downloader

type Client struct {
	contentLength      int64
	Link               string
	Parts              int
	Auto               bool
	downloadRangePairs [][2]int64
}

func NewClient(link string, parts int, auto bool) *Client {
	return &Client{
		Link:               link,
		Parts:              parts,
		Auto:               auto,
		downloadRangePairs: nil,
	}
}

func (c *Client) SetLink(link string) {
	c.Link = link
}

func (c *Client) SetParts(parts int) {
	c.Parts = parts
}
