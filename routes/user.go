package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/gobwas/ws"
  "github.com/gobwas/ws/wsutil"

  "github.com/patrickmn/go-cache"

  "encoding/json"
  "database/sql"
  "strconv"
  "regexp"

  "openverse/models"
  c "openverse/config"
)

func GetNotifications(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to view your notifications."})
    return
  }
  userId := user.(models.User).Id
  var notifications []models.Notification
  _, err := c.Dbmap.Select(&notifications, "select notif_by, notif_by_others, notif_type, notif_topic, notif_date, notif_read from notifications where notif_to = ? order by notif_date desc limit 64", userId)
  if err != nil && err != sql.ErrNoRows {
    ctx.JSON(500, gin.H{"error": "An error occurred while grabbing your notifications. " + err.Error()})
    return
  }
  _, err = c.Dbmap.Exec("update notifications set notif_read = 1 where notif_to = ?", userId)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while reading your notifications. " + err.Error()})
    return
  }
  // Update notif read on stream
  notificationStreamId := strconv.FormatInt(user.(models.User).Id, 32)
  models.GetNotificationBroadcast(notificationStreamId).Submit(0)

  for i, notification := range notifications {
    var notifUserIds []int64
    if notification.Others != nil {
      json.Unmarshal([]byte(notification.Others.ValueOrZero()), &notifUserIds)
    }
    // prepend UserId to notifUserIds
    notifUserIds = append([]int64{notification.UserId}, notifUserIds...)
    var notifUsers []models.User
    for _, id := range notifUserIds {
      var notifUser models.User
      userCached, found := c.CacheStore.Get("creator" + strconv.FormatInt(id, 32))
      if found {
        notifUser = userCached.(models.User)
      } else {
        c.Dbmap.SelectOne(&notifUser, "select user_id, user_name, user_rank, user_avatar from users where user_pid = ? limit 1", id)
        c.CacheStore.Set("creator" + strconv.FormatInt(id, 32), notifUser, cache.DefaultExpiration)
      }
      notifUsers = append(notifUsers, notifUser)
    }
    notifications[i].Users = notifUsers
    // assign the Topic for the specific type (post, comment...)
    switch(notification.Type) {
      case models.NotifPostYeah:
        // gave your post () a yeah
        var post models.Post
        c.Dbmap.SelectOne(&post, "select post_id, post_feeling_id, post_content, post_by from posts where post_id = ? limit 1", notification.Topic)
        notifications[i].Post = &post
      case models.NotifCommentYeah:
        // gave your comment () a yeah
        var comment models.Comment
        c.Dbmap.SelectOne(&comment, "select reply_id, reply_feeling_id, reply_content, reply_by from replies where reply_id = ? limit 1", notification.Topic)
        notifications[i].Comment = &comment
      case models.NotifCommentMyPost:
        // commented on your post ()
        var post models.Post
        c.Dbmap.SelectOne(&post, "select post_id, post_feeling_id, post_content, post_by from posts where post_id = ? limit 1", notification.Topic)
        notifications[i].Post = &post
      case models.NotifCommentOtherPost:
        // commented on someone's post ()
        var post models.Post
        c.Dbmap.SelectOne(&post, "select post_id, post_feeling_id, post_content, post_by from posts where post_id = ? limit 1", notification.Topic)
        notifications[i].Post = &post
    }
    // do merges here once that's implemented
  }

  ctx.JSON(200, notifications)
}

func GetNotificationsStream(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to view your notifications."})
    return
  }
  // The stream ID is base32
  notificationStreamId := strconv.FormatInt(user.(models.User).Id, 32)
  listener := models.OpenNotificationListener(notificationStreamId)
  defer models.CloseNotificationListener(notificationStreamId, listener)

  conn, _, _, err := ws.UpgradeHTTP(ctx.Request, ctx.Writer)
  if err != nil {
    return
  }
  defer conn.Close()
  for {
    recieved := <-listener
    if recieved != nil {
      if _, ok := recieved.(int); ok {
        message := strconv.Itoa(recieved.(int))
        wsutil.WriteServerMessage(conn, ws.OpText, []byte(message))
      }
      message, _ := json.Marshal(recieved)
      wsutil.WriteServerMessage(conn, ws.OpText, message)
    }
  }
}

func GetProfileSettings(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to edit your profile."})
    return
  }
  var profile models.UserProfile
  err := c.Dbmap.SelectOne(&profile, "select user_profile_comment, user_country, user_birthday, user_website, user_favorite_post, user_skill, user_nnid, user_relationship_visibility from users where user_pid = ? limit 1", user.(models.User).Id)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while searching for your profile. " + err.Error()})
    return
  }
  // Get user stuff, favorite post
  if profile.FavoritePostId != nil {
    c.Dbmap.SelectOne(&profile.FavoritePost, "select post_screenshot from posts where post_id = ? limit 1", profile.FavoritePostId)
  }

  ctx.JSON(200, profile)
}

func EditProfileSettings(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to edit your profile."})
    return
  }
  meUser := user.(models.User)
  /*var profile UserProfile
  err := c.Dbmap.SelectOne(&profile, "select user_profile_comment, user_country, user_birthday, user_website, user_favorite_post, user_nnid, user_relationship_visibility from users where user_pid = ? limit 1", meUser.Id)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while searching for your profile. " + err.Error()})
    return
  }
  */
  Name := ctx.PostForm("name")
  Nick := ctx.PostForm("nick")
  Avatar := ctx.PostForm("avatar")
  if Name == "" || Nick == "" || Avatar == "" {
    ctx.JSON(400, gin.H{"error": "The Username/Nickname/Avatar fields cannot be blank."})
    return
  }
  nameSyntaxValid, err := regexp.MatchString("^[^/]{2,32}$", Name)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "Something went wrong while verifying your account. " + err.Error()})
    return
  }
  if !nameSyntaxValid {
    ctx.JSON(400, gin.H{"error": "The username you entered contains an invalid character (forward slash), or is too long or short."})
    return
  }
  if Nick == "" {
    ctx.JSON(400, gin.H{"error": "The Nickname field cannot be blank."})
    return
  }
  if Nick != "" && len([]rune(Nick)) > 64 {
    ctx.JSON(400, gin.H{"error": "The nickname you entered is too long."})
    return
  }
  meUser.Name = Name
  meUser.Nick = Nick
  meUser.Avatar = Avatar
  _, err = c.Dbmap.Exec("update users set user_id = ?, user_name = ?, user_avatar = ? where user_pid = ?", Name, Nick, Avatar, meUser.Id)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while saving your profile. " + err.Error()})
    return
  }
  // Update cache
  c.CacheStore.Set("creator" + strconv.FormatInt(meUser.Id, 32), meUser, cache.DefaultExpiration)
}

func UserView(ctx *gin.Context) {
  userName := ctx.Param("id")
  var user models.UserProfile
  err := c.Dbmap.SelectOne(&user, "select user_pid, user_id, user_name, user_rank, user_avatar, user_date, user_profile_comment, user_country, user_birthday, user_website, user_skill, user_systems, user_favorite_post, user_favorite_post_type, user_nnid, user_relationship_visibility from users where user_id = ? limit 1", userName)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The user could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that user. " + err.Error()})
    return
  }
  // Get user stuff, favorite post
  if user.FavoritePostId != nil {
    c.Dbmap.SelectOne(&user.FavoritePost, "select post_screenshot from posts where post_id = ? limit 1", user.FavoritePostId)
  }
  //meUserObj, loggedIn := ctx.Get("user")
  // TODO: literally everything
  ctx.JSON(200, user)
}
