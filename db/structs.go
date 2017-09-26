package db

// Article is a VK page or Zen page
type Article struct {
	ID         int    `json:"id"`
	URL        string `json:"url"`
	Title      string `json:"title"`
	VKWallID   int    `json:"vk_wall_id"`
	Categories []int  `json:"categories"`
}

// Category of an article
type Category struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	VKPublicID int    `json:"vk_public_id"` // is shown for everyone
	VKStatsID  int    `json:"vk_stats_id"`  // is used for creating articles
}

// Job is a single mention of a user in an article
type Job struct {
	ID        int    `json:"id"` // TODO: check if ID is not used when saving job with db.Save
	ArticleID int    `json:"article_id"`
	UserID    int    `json:"user_id"`
	Kind      string `json:"kind"`
}

// State is a current user state
type State struct {
	ID    int
	User  string
	State int
	Input string
}

// StateInput is a struct of all possible inputs. Stored as string
type StateInput struct {
	Article *Article `json:"article"`
	Jobs    []*Job   `json:"jobs"`
}

// Worker is a translator, editor (or anyone else involved)
type Worker struct {
	ID        int
	VKID      int
	ShortName string
	Active    bool
}
