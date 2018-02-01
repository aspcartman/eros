package tinder

import "time"

type User struct {
	GroupMatched      bool          `json:"group_matched"`
	DistanceMi        int           `json:"distance_mi"`
	ContentHash       string        `json:"content_hash"`
	CommonFriends     []interface{} `json:"common_friends"`
	CommonLikes       []interface{} `json:"common_likes"`
	CommonFriendCount int           `json:"common_friend_count"`
	CommonLikeCount   int           `json:"common_like_count"`
	ConnectionCount   int           `json:"connection_count"`
	ID                string        `json:"_id"`
	Bio               string        `json:"bio"`
	BirthDate         time.Time     `json:"birth_date"`
	Name              string        `json:"name"`
	PingTime          time.Time     `json:"ping_time"`
	Photos []struct {
		ID  string `json:"id"`
		URL string `json:"url"`
		ProcessedFiles []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"processedFiles"`
		FileName  string `json:"fileName"`
		Extension string `json:"extension"`
	} `json:"photos"`
	Instagram struct {
		LastFetchTime         time.Time `json:"last_fetch_time"`
		CompletedInitialFetch bool      `json:"completed_initial_fetch"`
		Photos []struct {
			Image     string `json:"image"`
			Thumbnail string `json:"thumbnail"`
			Ts        string `json:"ts"`
			Link      string `json:"link"`
		} `json:"photos"`
		MediaCount     int    `json:"media_count"`
		ProfilePicture string `json:"profile_picture"`
		Username       string `json:"username"`
	} `json:"instagram"`
	Jobs    []interface{} `json:"jobs"`
	Schools []interface{} `json:"schools"`
	Teaser struct {
		String string `json:"string"`
	} `json:"teaser"`
	Teasers []struct {
		Type   string `json:"type"`
		String string `json:"string"`
	} `json:"teasers"`
	Gender        int    `json:"gender"`
	BirthDateInfo string `json:"birth_date_info"`
	SNumber       int    `json:"s_number"`
}
