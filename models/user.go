package models

import (
  "github.com/guregu/null"
  "github.com/go-gorp/gorp"

  "strings"
)

type User struct {
  Id     int64     `db:"user_pid" json:"id,omitempty"`
  Name   string    `db:"user_id" json:"name,omitempty"`
  Nick   string    `db:"user_name" json:"nick,omitempty"`
  Pass   []byte    `db:"user_pass" json:"pass,omitempty"`
  Email  string    `db:"user_email" json:"email,omitempty"`
  Avatar string    `db:"user_avatar" json:"avatar,omitempty"`
  Ip     string    `db:"user_ip" json:"ip,omitempty"`
  Rank   int       `db:"user_rank" json:"rank,omitempty"`
}

func (u User) LoggedIn() bool {
  if u.Id > 0 {
    return true
  }
  return false
}

func (u User) RankClass() string {
  switch(u.Rank) {
    case 1:
      return "donator"
    case 2:
      return "tester"
    case 3:
      return "moderator"
    case 4:
      return "administrator"
    case 5:
      return "developer"
  }
  return ""
}
// Used in "deleted by *"
func (u User) RankText() string {
  switch(u.Rank) {
    case 1:
      return "Donator"
    case 2:
      return "Tester"
    case 3:
      return "Moderator"
    case 4:
      return "Administrator"
    case 5:
      return "Developer"
    default:
      return ""
  }
  return ""
}

func (u User) DoAvatarFeeling(feelingId int) string {
  if strings.HasPrefix(u.Avatar, "http") {
    return u.Avatar
  }
  feeling := "_normal"
  switch(feelingId) {
    case 1:
      feeling = "_happy"
    case 2:
      feeling = "_like"
    case 3:
      feeling = "_surprised"
    case 4:
      feeling = "_frustrated"
    case 5:
      feeling = "_puzzled"
  }
  return "https://mii-secure.cdn.nintendo.net/" + u.Avatar + feeling + "_face.png"
}
func (u User) DoAvatar() string {
  if strings.HasPrefix(u.Avatar, "http") {
    return u.Avatar
  }
  return "https://mii-secure.cdn.nintendo.net/" + u.Avatar + "_normal_face.png"
}

type UserProfile struct {
  User
  Date                   *gorp.NullTime `db:"user_date" json:"date,omitempty"`
  ProfileComment         string         `db:"user_profile_comment" json:"profile_comment,omitempty"`
  Country                string         `db:"user_country" json:"country,omitempty"`
  Birthday               *gorp.NullTime `db:"user_birthday" json:"birthday,omitempty"`
  Website                string         `db:"user_website" json:"website,omitempty"`
  Skill                  int            `db:"user_skill" json:"skill,omitempty"`
  Systems                int            `db:"user_systems" json:"systems,omitempty"`
  FavoritePost           *Post          `db:"-" json:"favorite_post,omitempty"`
  FavoritePostId         *null.Int      `db:"user_favorite_post" json:"favorite_post_id,omitempty"`
  FavoritePostType       int            `db:"user_favorite_post_type" json:"favorite_post_type,omitempty"`
  Nnid                   string         `db:"user_nnid" json:"nnid,omitempty"`
  Code                   string         `db:"user_code" json:"code,omitempty"`
  EmailConfirmed         bool           `db:"user_email_confirmed" json:"email_confirmed,omitempty"`
  RelationshipVisibility int            `db:"user_relationship_visibility" json:"relationship_visibility,omitempty"`
  Role                   int            `db:"user_role" json:"role,omitempty"`
}
