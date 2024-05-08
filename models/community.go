package models

import (
  "github.com/guregu/null"
  "github.com/go-gorp/gorp"
)

type Title struct {
  Id       int64  `db:"title_id" json:"id,omitempty"`
  Type     int    `db:"title_type" json:"type,omitempty"`
  Name     string `db:"title_name" json:"name,omitempty"`
  Icon     string `db:"title_icon" json:"icon,omitempty"`
  Banner   string `db:"title_banner" json:"banner,omitempty"`
  Platform int    `db:"title_platform" json:"platform,omitempty"`
}

type Community struct {
  Id          int64  `db:"community_id" json:"id,omitempty"`
  TitleId     int64  `db:"community_title" json:"title_id,omitempty"`
  Type        int    `db:"community_type" json:"type,omitempty"`
  Name        string `db:"community_name" json:"name,omitempty"`
  Description string `db:"community_description" json:"description,omitempty"`
  Icon        string `db:"community_icon" json:"icon,omitempty"`
  Banner      string `db:"community_banner" json:"banner,omitempty"`
  Platform    int    `db:"community_platform" json:"platform,omitempty"`
  Permission  int    `db:"community_perms" json:"permission,omitempty"`
}

type Post struct {
  Id             int64           `db:"post_id" json:"id,omitempty"`
  // Making things pointers ensures that interfaces are omitted if they're empty
  Community      *Community      `db:"-" json:"community,omitempty"`
  CommunityId    int64           `db:"post_community" json:"community_id,omitempty"`
  User           User            `db:"-" json:"user,omitempty"`
  UserId         int64           `db:"post_by" json:"user_id,omitempty"`
  Feeling        int             `db:"post_feeling_id" json:"feeling,omitempty"`
  Content        string          `db:"post_content" json:"content,omitempty"`
  Screenshot     *null.String     `db:"post_screenshot" json:"screenshot,omitempty"`
  Drawing        *null.String     `db:"post_drawing" json:"drawing,omitempty"`
  Url            *null.String     `db:"post_url" json:"url,omitempty"`
  Spoiler        bool            `db:"post_is_spoiler" json:"spoiler,omitempty"`
  Date           *gorp.NullTime       `db:"post_date" json:"date,omitempty"`
  Status         int             `db:"post_status" json:"status,omitempty"`
  Edited         *gorp.NullTime       `db:"post_edited" json:"edited,omitempty"`
  //Html       bool            `db:"post_html" json:"html,omitempty"`
  YeahCount      int64           `db:"-" json:"yeah_count,omitempty"`
  CanYeah        bool            `db:"-" json:"can_yeah,omitempty"`
  HasYeah        bool            `db:"-" json:"has_yeah,omitempty"`
  CommentCount   int64           `db:"-" json:"comment_count,omitempty"`
  RecentComment  *Comment        `db:"-" json:"recent_comment,omitempty"`
}

type PostBefores struct {
  Post
  // json.Unmarshal requires a []byte so this makes it easier
  ContentBefores    []byte `db:"post_content_before", json:"-,omitempty"`
  ScreenshotBefores []byte `db:"post_screenshot_before", json:"-,omitempty"`
}

type Comment struct {
  Id           int64           `db:"reply_id" json:"id,omitempty"`
  //Post         Post            `db:"-" json:"post,omitempty"`
  PostId       int64           `db:"reply_to" json:"post_id,omitempty"`
  User         User            `db:"-" json:"user,omitempty"`
  UserId       int64           `db:"reply_by" json:"user_id,omitempty"`
  Feeling      int             `db:"reply_feeling_id" json:"feeling,omitempty"`
  Content      string          `db:"reply_content" json:"content,omitempty"`
  Screenshot   *null.String     `db:"reply_screenshot" json:"screenshot,omitempty"`
  Drawing      *null.String     `db:"reply_drawing" json:"drawing,omitempty"`
  Url          *null.String     `db:"reply_url" json:"url,omitempty"`
  Spoiler      bool            `db:"reply_is_spoiler" json:"spoiler,omitempty"`
  Date         *gorp.NullTime       `db:"reply_date" json:"date,omitempty"`
  Status       int             `db:"reply_status" json:"status,omitempty"`
  Edited       *gorp.NullTime       `db:"reply_edited" json:"edited,omitempty"`
  //Html         bool            `db:"reply_html" json:"html,omitempty"`
  YeahCount    int64           `db:"-" json:"yeah_count,omitempty"`
  CanYeah      bool            `db:"-" json:"can_yeah,omitempty"`
  HasYeah      bool            `db:"-" json:"has_yeah,omitempty"`
}

type CommentBefores struct {
  Comment
  // json.Unmarshal requires a []byte so this makes it easier
  ContentBefores    []byte `db:"reply_content_before", json:"-,omitempty"`
  ScreenshotBefores []byte `db:"reply_screenshot_before", json:"-,omitempty"`
}

type PostYeah struct {
  Id     int64         `db:"yeah_id" json:"id,omitempty"`
  PostId int64         `db:"yeah_post" json:"post_id,omitempty"`
  User   User          `db:"-" json:"user,omitempty"`
  UserId int64         `db:"yeah_by" json:"user_id,omitempty"`
  //Date   *gorp.NullTime `db:"yeah_date" json:"date,omitempty"`
}

type CommentYeah struct {
  Id        int64         `db:"ryeah_id" json:"id,omitempty"`
  CommentId int64         `db:"ryeah_reply" json:"comment_id,omitempty"`
  User      User          `db:"-" json:"user,omitempty"`
  UserId    int64         `db:"ryeah_by" json:"user_id,omitempty"`
  //Date   *gorp.NullTime `db:"yeah_date" json:"date,omitempty"`
}
