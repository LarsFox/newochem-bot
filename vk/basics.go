package vk

const (
	newochemID = -80512191

	vkMethod = "https://api.vk.com/method/%s?access_token=%s&v=%s"
	vkCount  = 100
)

// GetWallResponse is a response for api.vk.com/wall.get
type GetWallResponse struct {
	Count int            `json:"count"`
	Items []*getWallItem `json:"items"`
}

type getWallItem struct {
	ID          int                  `json:"id"`
	Date        int                  `json:"date"`
	Text        string               `json:"text"`
	Attachments []*getWallAttachment `json:"attachments"`
	Likes       []*counter           `json:"likes"`
	Reposts     []*counter           `json:"reposts"`
	Views       []*counter           `json:"views"`
}

type getWallAttachment struct {
	Type string          `json:"type"`
	Link *attachmentLink `json:"link"`
}

type attachmentLink struct {
	URL     string           `json:"url"`
	Title   string           `json:"title"`
	Caption string           `json:"caption"`
	Photo   *attachmentPhoto `json:"photo"`
}

type attachmentPhoto struct {
	ID        int    `json:"id"`
	AlbumID   int    `json:"album_id"`
	Photo1280 string `json:"photo_1280"` // need more?
}

type counter struct {
	Count int `json:"count"`
}
