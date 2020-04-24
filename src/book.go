package main

type Book struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Author           string `json:"author"`
	AuthorAvatar     string `json:"author_avatar"`
	Abbreviation     string `json:"abbreviation"`
	Cover            string `json:"cover"`
	Finished         bool   `json:"finished"`
	TotalReads       int    `json:"total_reads"`
	TotalChars       int    `json:"total_chars"`
	LastUpdateTime   string `json:"last_update_time"`
	Class            string `json:"class"`
	TotalSearches    int    `json:"total_searches"`
	TotalVotes       int    `json:"total_votes"`
	LastChapterTitle string `json:"last_chapter_title"`
	LastChapterUrl   string `json:"last_chapter_url"`
	WithVIPChapter   bool   `json:"with_vip_chapter"`
	Gender           string `json:"gender"`
	Score            int    `json:"score"`
	RWords           string `json:"rwords,omitempty"`
	RUser            string `json:"ruser,omitempty"`
}

type Chapter struct {
	NativeId int    `json:"native_id"`
	Id       string `json:"id"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	Vip      bool   `json:"vip"`
}
