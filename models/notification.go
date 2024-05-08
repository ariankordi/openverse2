package models

import (
  "github.com/guregu/null"
  "github.com/go-gorp/gorp"

  "database/sql"
  "encoding/json"
  "strconv"

  c "openverse/config"
)

type Notification struct {
  Id       int64          `db:"notif_id" json:"id,omitempty"`
  Type     int            `db:"notif_type" json:"type,omitempty"`
  //UserTo   User           `db:"-" json:"user_to,omitempty"`
  //UserToId int64          `db:"notif_to" json:"user_to_id,omitempty"`
  //User     User           `db:"-" json:"user,omitempty"`
  UserId   int64          `db:"notif_by" json:"user_id,omitempty"`
  Users    []User         `db:"-" json:"users,omitempty"`
  // JSON, other users who have merged with this notification
  Others   *null.String `db:"notif_by_others" json:"-,omitempty"`
  Topic    int64          `db:"notif_topic" json:"topic,omitempty"`
  Post     *Post          `db:"-" json:"post,omitempty"`
  Comment  *Comment       `db:"-" json:"comment,omitempty"`
  Date     *gorp.NullTime  `db:"notif_date" json:"date,omitempty"`
  Read     bool           `db:"notif_read" json:"read,omitempty"`
}

type ChannelMessage struct {
  Topic   int64    `json:"topic,omitempty"`
  Type    int      `json:"-,omitempty"`
  Comment *Comment `json:"comment,omitempty"`
  Post    *Post    `json:"post,omitempty"`
  UserId  int64    `json:"user_id,omitempty"`
}

type NotificationHint struct {
  Topic      int64 `json:"topic,omitempty"`
  OtherTopic int64 `json:"other_topic,omitempty"`
  Type       int   `json:"type,omitempty"`
}

const (
  // Notification types
  NotifPostYeah = 0
  NotifCommentYeah = 1
  NotifCommentMyPost = 2
  NotifCommentOtherPost = 3
  NotifFollow = 4
  // Channel message types
  MessageYeahPost = 0
  MessageUnyeahPost = 1
  MessageAddComment = 2
  MessageEditPost = 3
)

// Helper function to send a notification since it's complicated
// Type can be the consts defined above, topic is the post/comment/whatever being mentioned
func SendNotification(source int64, target int64, notifType int, topic int64, otherTopic int64) error {
  // Push to stream
  notificationStreamId := strconv.FormatInt(target, 32)
  GetNotificationBroadcast(notificationStreamId).Submit(NotificationHint{
    Topic: topic,
    OtherTopic: otherTopic,
    Type: notifType,
  })
  // Look for a notification to merge with, or to re-notify
  var notificationMergeSearch Notification
  err := c.Dbmap.SelectOne(&notificationMergeSearch, "select notif_id, notif_by, notif_by_others from notifications where notif_type = ? and notif_topic = ? order by notif_id desc limit 1", notifType, topic)
  if err == nil {
    // get user IDs for this notification
    var notifUserIds []int64
    if notificationMergeSearch.Others != nil {
      json.Unmarshal([]byte(notificationMergeSearch.Others.ValueOrZero()), &notifUserIds)
    }
    // if the user is in merges with the notification, then just update it
    sourceInNotification := false
    for _, id := range notifUserIds {
      if id == source {
        sourceInNotification = true
      }
    }
    if sourceInNotification || source == notificationMergeSearch.UserId {
      // TODO: some sort of spam detection?
      _, err = c.Dbmap.Exec("update notifications set notif_read = 0, notif_date = now() where notif_id = ?", notificationMergeSearch.Id)
      return err
    }
    // If source is not in the notification, update merges
    notifUserIds = append(notifUserIds, source)
    notifUserIdsJson, _ := json.Marshal(notifUserIds)
    _, err = c.Dbmap.Exec("update notifications set notif_by_others = ?, notif_read = 0, notif_date = now() where notif_id = ?", notifUserIdsJson, notificationMergeSearch.Id)
    return err
  }
  if err != sql.ErrNoRows {
    return err
  }

  _, err = c.Dbmap.Exec("insert into notifications(notif_type, notif_to, notif_by, notif_topic) values(?, ?, ?, ?)", notifType, target, source, topic)
  return err
}
