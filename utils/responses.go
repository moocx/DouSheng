package utils

import "time"

type VideoResponse struct {
	StatusCode int            `json:"status_code"`
	StatusMsg  string         `json:"status_msg"`
	NextTime   int64          `json:"next_time"`
	VideoList  []VideoResItem `json:"video_list"`
}

type VideoResItem struct {
	ID            uint   `json:"id"`
	Author        Author `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount uint   `json:"favorite_count"`
	CommentCount  uint   `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

type Author struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	FollowCount    int    `json:"follow_count"`
	FollowerCount  int    `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	Avatar         string `json:"avatar"`
	Background     string `json:"background_image"`
	Signature      string `json:"signature"`
	TotalFavorited string `json:"total_favorited"`
	WorkCount      int    `json:"work_count"`
	FavoriteCount  int    `json:"favorite_count"`
}

type CommentListResponse struct {
	StatusCode  int              `json:"status_code"`
	StatusMsg   string           `json:"status_msg"`
	CommentList []CommentResItem `json:"comment_list"`
}

type CommentResItem struct {
	ID         uint      `json:"id"`
	User       Author    `json:"user"`
	Content    string    `json:"content"`
	CreateDate time.Time `json:"create_date"`
}

type CommentResponse struct {
	StatusCode int             `json:"status_code"`
	StatusMsg  string          `json:"status_msg"`
	Comment    *CommentResItem `json:"comment,omitempty"`
}
