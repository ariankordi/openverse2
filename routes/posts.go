package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/gobwas/ws"
  "github.com/gobwas/ws/wsutil"

  "github.com/go-gorp/gorp"

  "github.com/patrickmn/go-cache"

  "github.com/guregu/null"

  "encoding/json"
  "database/sql"
  "strconv"
  "strings"

  "openverse/models"
  "openverse/util"
  c "openverse/config"
)

func PostAddYeah(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to Yeah a post."})
    return
  }
  userId := user.(models.User).Id
  postId := ctx.Param("id")
  var post models.Post
  err := c.Dbmap.SelectOne(&post, "select post_id, post_by, post_community, post_date from posts where post_id = ? and post_status < 2 limit 1", postId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The post could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that post. " + err.Error()})
    return
  }
  if post.UserId == userId {
    ctx.JSON(403, gin.H{"error": "You cannot give a Yeah to your own post."})
    return
  }
  // check for yeah already given
  myYeahs, _ := c.Dbmap.SelectInt("select count(yeah_id) from post_yeahs where yeah_post = ? and yeah_by = ?", post.Id, userId)
  if myYeahs > 0 {
    ctx.JSON(400, gin.H{"error": "You have already given a Yeah to this post."})
    return
  }
  // Now send the notification!
  err = models.SendNotification(userId, post.UserId, models.NotifPostYeah, post.Id, int64(0))
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to send a notification. " + err.Error()})
    return
  }
  _, err = c.Dbmap.Exec("insert into post_yeahs(yeah_by, yeah_post) values(?, ?)", userId, post.Id)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to give a Yeah to that post. " + err.Error()})
    return
  }
  // Update cache
  yeahCountCached, found := c.CacheStore.Get("yeahcount" + strconv.FormatInt(post.Id, 32))
  if found {
    post.YeahCount = yeahCountCached.(int64)
    c.CacheStore.Set("yeahcount" + strconv.FormatInt(post.Id, 32), post.YeahCount + 1, cache.DefaultExpiration)
  }
  // Update my yeah
  _, found = c.CacheStore.Get("yeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(post.Id, 32))
  if found {
    c.CacheStore.Set("yeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(post.Id, 32), int64(1), cache.DefaultExpiration)
  }
  if(models.ProperSince(post.Date) < 3600000000000) {
    communityStreamId := strconv.FormatInt(post.CommunityId, 32)
    models.GetCommunityBroadcast(communityStreamId).Submit(models.ChannelMessage{
      Topic: post.Id,
      Type: models.MessageYeahPost,
      UserId: userId,
    })
  }
  postStreamId := strconv.FormatInt(post.Id, 32)
  models.GetPostBroadcast(postStreamId).Submit(user)
  ctx.Status(201)
}

func PostDeleteYeah(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to Yeah a post."})
    return
  }
  userId := user.(models.User).Id
  postId := ctx.Param("id")
  var post models.Post
  err := c.Dbmap.SelectOne(&post, "select post_id, post_by, post_community, post_date from posts where post_id = ? and post_status < 2 limit 1", postId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The post could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that post. " + err.Error()})
    return
  }
  if post.UserId == userId {
    ctx.JSON(403, gin.H{"error": "You cannot Unyeah your own post."})
    return
  }
  _, err = c.Dbmap.Exec("delete from post_yeahs where yeah_post = ? and yeah_by = ?", post.Id, userId)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to Unyeah that post. " + err.Error()})
    return
  }
  // Update cache
  yeahCountCached, found := c.CacheStore.Get("yeahcount" + strconv.FormatInt(post.Id, 32))
  if found {
    post.YeahCount = yeahCountCached.(int64)
    c.CacheStore.Set("yeahcount" + strconv.FormatInt(post.Id, 32), post.YeahCount - 1, cache.DefaultExpiration)
  }
  // Update my yeah
  _, found = c.CacheStore.Get("yeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(post.Id, 32))
  if found {
    c.CacheStore.Set("yeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(post.Id, 32), int64(0), cache.DefaultExpiration)
  }
  if(models.ProperSince(post.Date) < 3600000000000) {
    communityStreamId := strconv.FormatInt(post.CommunityId, 32)
    models.GetCommunityBroadcast(communityStreamId).Submit(models.ChannelMessage{
      Topic: post.Id,
      Type: models.MessageUnyeahPost,
      UserId: userId,
    })
  }
  postStreamId := strconv.FormatInt(post.Id, 32)
  models.GetPostBroadcast(postStreamId).Submit(userId)
  ctx.Status(202)
}

func PostView(ctx *gin.Context) {
  postId := ctx.Param("id")
  var post models.Post
  err := c.Dbmap.SelectOne(&post, "select post_id, post_community, post_by, post_feeling_id, post_content, post_screenshot, post_drawing, post_url, post_is_spoiler, post_date, post_status, post_edited from posts where post_id = ? limit 1", postId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The post could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that post. " + err.Error()})
    return
  }
  if post.Status > 1 {
    switch(post.Status) {
      case 2:
        ctx.JSON(404, gin.H{"error": "Deleted by poster."})
      default:
        userDummyRank := strings.ToLower(models.User{Rank: post.Status}.RankText())
        ctx.JSON(404, gin.H{"error": "Deleted by " + userDummyRank + "."})
    }
    return
  }
  user, loggedIn := ctx.Get("user")

  communityCached, found := c.CacheStore.Get("community" + strconv.FormatInt(post.CommunityId, 32))
  if found {
    communityType := communityCached.(models.Community)
    post.Community = &communityType
  } else {
    c.Dbmap.SelectOne(&post.Community, "select community_id, community_name, community_icon from communities where community_id = ? limit 1", post.CommunityId)
    // Not setting the cache this time since it's incomplete
  }

  userCached, found := c.CacheStore.Get("creator" + strconv.FormatInt(post.UserId, 32))
  if found {
    post.User = userCached.(models.User)
  } else {
    c.Dbmap.SelectOne(&post.User, "select user_id, user_name, user_rank, user_avatar from users where user_pid = ? limit 1", post.UserId)
    c.CacheStore.Set("creator" + strconv.FormatInt(post.UserId, 32), post.User, cache.DefaultExpiration)
  }

  yeahCountCached, found := c.CacheStore.Get("yeahcount" + strconv.FormatInt(post.Id, 32))
  if found {
    post.YeahCount = yeahCountCached.(int64)
  } else {
    post.YeahCount, _ = c.Dbmap.SelectInt("select count(yeah_id) from post_yeahs where yeah_post = ?", post.Id)
    c.CacheStore.Set("yeahcount" + strconv.FormatInt(post.Id, 32), post.YeahCount, cache.DefaultExpiration)
  }
  if loggedIn {
    // has yeah cache
    myYeahs := int64(0)
    myYeahsCached, found := c.CacheStore.Get("yeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(post.Id, 32))
    if found {
      myYeahs = myYeahsCached.(int64)
    } else {
      myYeahs, _ = c.Dbmap.SelectInt("select count(yeah_id) from post_yeahs where yeah_post = ? and yeah_by = ?", post.Id, user.(models.User).Id)
      c.CacheStore.Set("yeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(post.Id, 32), myYeahs, cache.DefaultExpiration)
    }
    post.HasYeah = myYeahs > 0
    post.CanYeah = (post.UserId != user.(models.User).Id)
  } else {
    post.HasYeah = false
    post.CanYeah = false
  }
  // comment count cache
  commentCountCache, found := c.CacheStore.Get("commentcount" + strconv.FormatInt(post.Id, 32))
  if found {
    post.CommentCount = commentCountCache.(int64)
  } else {
    post.CommentCount, _ = c.Dbmap.SelectInt("select count(reply_id) from replies where reply_to = ? and reply_status < 2", post.Id)
    c.CacheStore.Set("commentcount" + strconv.FormatInt(post.Id, 32), post.CommentCount, cache.DefaultExpiration)
  }

  var yeahs []models.PostYeah
  _, err = c.Dbmap.Select(&yeahs, "select yeah_by from post_yeahs where yeah_post = ? order by yeah_id desc limit 50", post.Id)
  if err != nil && err != sql.ErrNoRows {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to get the Yeahs for that post. " + err.Error()})
    return
  }
  for i, yeah := range yeahs {
    userCached, found := c.CacheStore.Get("creator" + strconv.FormatInt(yeah.UserId, 32))
    if found {
      yeahs[i].User = userCached.(models.User)
    } else {
      c.Dbmap.SelectOne(&yeahs[i].User, "select user_id, user_rank, user_avatar from users where user_pid = ? limit 1", yeah.UserId)
      // Not setting the cache this time since it's incomplete
    }
  }

  commentOffset := int64(0)
  if(post.CommentCount > 19) {
    commentOffset = post.CommentCount - 20
  }

  var comments []models.Comment
  _, err = c.Dbmap.Select(&comments, "select reply_id, reply_content, reply_feeling_id, reply_screenshot, reply_drawing, reply_url, reply_date, reply_edited, reply_is_spoiler, reply_by from replies where reply_to = ? and reply_status < 2 order by reply_date asc limit 20 offset ?", post.Id, commentOffset)
  if err != nil && err != sql.ErrNoRows {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to get the comments for that post. " + err.Error()})
    return
  }
  for i, comment := range comments {
    userCached, found = c.CacheStore.Get("creator" + strconv.FormatInt(comment.UserId, 32))
    if found {
      comments[i].User = userCached.(models.User)
    } else {
      c.Dbmap.SelectOne(&comments[i].User, "select user_id, user_name, user_rank, user_avatar from users where user_pid = ? limit 1", comment.UserId)
      c.CacheStore.Set("creator" + strconv.FormatInt(comment.UserId, 32), comments[i].User, cache.DefaultExpiration)
    }
    yeahCountCached, found := c.CacheStore.Get("commentyeahcount" + strconv.FormatInt(comment.Id, 32))
    if found {
      comments[i].YeahCount = yeahCountCached.(int64)
    } else {
      comments[i].YeahCount, _ = c.Dbmap.SelectInt("select count(ryeah_id) from reply_yeahs where ryeah_reply = ?", comment.Id)
      c.CacheStore.Set("commentyeahcount" + strconv.FormatInt(comment.Id, 32), comments[i].YeahCount, cache.DefaultExpiration)
    }
    if loggedIn {
      // has yeah cache
      myYeahs := int64(0)
      myYeahsCached, found := c.CacheStore.Get("commentyeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(comment.Id, 32))
      if found {
        myYeahs = myYeahsCached.(int64)
      } else {
        myYeahs, _ := c.Dbmap.SelectInt("select count(ryeah_id) from reply_yeahs where ryeah_reply = ? and ryeah_by = ?", comment.Id, user.(models.User).Id)
        c.CacheStore.Set("commentyeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(comment.Id, 32), myYeahs, cache.DefaultExpiration)
      }
      comments[i].HasYeah = myYeahs > 0
      comments[i].CanYeah = (comment.UserId != user.(models.User).Id)
    } else {
      comments[i].HasYeah = false
      comments[i].CanYeah = false
    }
  }

  ctx.JSON(200, gin.H{"post": post, "yeahs": yeahs, "comments": comments})
}

func PostStream(ctx *gin.Context) {
  postId := ctx.Param("id")
  var post models.Post
  err := c.Dbmap.SelectOne(&post, "select post_id from posts where post_id = ? limit 1", postId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The post could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that post. " + err.Error()})
    return
  }
  // The stream ID is base32
  postStreamId := strconv.FormatInt(post.Id, 32)
  listener := models.OpenPostListener(postStreamId)
	defer models.ClosePostListener(postStreamId, listener)

  conn, _, _, err := ws.UpgradeHTTP(ctx.Request, ctx.Writer)
  if err != nil {
    return
  }
  defer conn.Close()
  for {
    recieved := <-listener
    if recieved != nil {
      message := ""
      if user, ok := recieved.(models.User); ok {
        encodedUser, _ := json.Marshal(user)
        message = "yeah:" + string(encodedUser)
      } else if userId, ok := recieved.(int64); ok {
        encodedUserId := strconv.FormatInt(userId, 10)
        message = "unyeah:" + string(encodedUserId)
      } else if comment, ok := recieved.(models.Comment); ok {
        encodedComment, _ := json.Marshal(comment)
        message = "comment:" + string(encodedComment)
      } else if post, ok := recieved.(models.Post); ok {
        encodedPost, _ := json.Marshal(post)
        message = "edit:" + string(encodedPost)
      } else if recievedmessage, ok := recieved.(models.ChannelMessage); ok {
        encodedMessage, _ := json.Marshal(recievedmessage)
        switch(recievedmessage.Type) {
          case models.MessageYeahPost:
            message = "commentyeah:" + string(encodedMessage)
          case models.MessageUnyeahPost:
            message = "commentunyeah:" + string(encodedMessage)
          case models.MessageEditPost:
            message = "commentedit:" + string(encodedMessage)
        }
      }
      wsutil.WriteServerMessage(conn, ws.OpText, []byte(message))
    }
  }
}

func CommentAddYeah(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to Yeah a reply."})
    return
  }
  userId := user.(models.User).Id
  commentId := ctx.Param("id")
  var comment models.Comment
  err := c.Dbmap.SelectOne(&comment, "select reply_id, reply_by, reply_to from replies where reply_id = ? and reply_status < 2 limit 1", commentId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The reply could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that reply. " + err.Error()})
    return
  }
  if comment.UserId == userId {
    ctx.JSON(403, gin.H{"error": "You cannot give a Yeah to your own reply."})
    return
  }
  // check for yeah already given
  myYeahs, _ := c.Dbmap.SelectInt("select count(ryeah_id) from reply_yeahs where ryeah_reply = ? and ryeah_by = ?", comment.Id, userId)
  if myYeahs > 0 {
    ctx.JSON(400, gin.H{"error": "You have already given a Yeah to this reply."})
    return
  }
  // Now send a notification!
  err = models.SendNotification(userId, comment.UserId, models.NotifCommentYeah, comment.Id, comment.PostId)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to send a notification. " + err.Error()})
    return
  }
  _, err = c.Dbmap.Exec("insert into reply_yeahs(ryeah_by, ryeah_reply) values(?, ?)", userId, comment.Id)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to give a Yeah to that reply. " + err.Error()})
    return
  }
  // Update cache
  yeahCountCached, found := c.CacheStore.Get("commentyeahcount" + strconv.FormatInt(comment.Id, 32))
  if found {
    comment.YeahCount = yeahCountCached.(int64)
    c.CacheStore.Set("commentyeahcount" + strconv.FormatInt(comment.Id, 32), comment.YeahCount + 1, cache.DefaultExpiration)
  }
  // Update my yeah
  _, found = c.CacheStore.Get("commentyeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(comment.Id, 32))
  if found {
    c.CacheStore.Set("commentyeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(comment.Id, 32), int64(1), cache.DefaultExpiration)
  }
  postStreamId := strconv.FormatInt(comment.PostId, 32)
  models.GetPostBroadcast(postStreamId).Submit(models.ChannelMessage{
    Topic: comment.Id,
    Type: models.MessageYeahPost,
    UserId: userId,
  })
  ctx.Status(201)
}

func CommentDeleteYeah(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to Yeah a reply."})
    return
  }
  userId := user.(models.User).Id
  commentId := ctx.Param("id")
  var comment models.Comment
  err := c.Dbmap.SelectOne(&comment, "select reply_id, reply_by, reply_to from replies where reply_id = ? and reply_status < 2 limit 1", commentId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The reply could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that reply. " + err.Error()})
    return
  }
  if comment.UserId == userId {
    ctx.JSON(403, gin.H{"error": "You cannot Unyeah your own reply."})
    return
  }
  _, err = c.Dbmap.Exec("delete from reply_yeahs where ryeah_reply = ? and ryeah_by = ?", comment.Id, userId)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to Unyeah that reply. " + err.Error()})
    return
  }
  // Update cache
  yeahCountCached, found := c.CacheStore.Get("commentyeahcount" + strconv.FormatInt(comment.Id, 32))
  if found {
    comment.YeahCount = yeahCountCached.(int64)
    c.CacheStore.Set("commentyeahcount" + strconv.FormatInt(comment.Id, 32), comment.YeahCount - 1, cache.DefaultExpiration)
  }
  // Update my yeah
  _, found = c.CacheStore.Get("commentyeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(comment.Id, 32))
  if found {
    c.CacheStore.Set("commentyeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(comment.Id, 32), int64(0), cache.DefaultExpiration)
  }
  postStreamId := strconv.FormatInt(comment.PostId, 32)
  models.GetPostBroadcast(postStreamId).Submit(models.ChannelMessage{
    Topic: comment.Id,
    Type: models.MessageUnyeahPost,
    UserId: userId,
  })
  ctx.Status(202)
}

func PostAddComment(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to reply to a post."})
    return
  }
  userId := user.(models.User).Id
  postId := ctx.Param("id")
  var post models.Post
  err := c.Dbmap.SelectOne(&post, "select post_id, post_by, post_community, post_date from posts where post_id = ? and post_status < 2 limit 1", postId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The post could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that post. " + err.Error()})
    return
  }
  body := ctx.PostForm("body")
  feelingId := ctx.PostForm("feeling_id")
  isSpoiler := ctx.PostForm("is_spoiler")
  screenshot := ctx.PostForm("screenshot")
  if body == "" && screenshot == "" {
    ctx.JSON(400, gin.H{"error": "Your reply is empty."})
    return
  }
  if len([]rune(body)) > 5000 {
    ctx.JSON(400, gin.H{"error": "Your reply is too long."})
    return
  }
  isSpoilerBool := false
  if isSpoiler == "true" {
    isSpoilerBool = true
  }
  feelingIdInt, _ := strconv.Atoi(feelingId)

  uploadedScreenshot := ""
  if screenshot != "" {
    uploadedScreenshot, err = util.OptimizeAndUpload(screenshot)
    if err != nil {
      ctx.JSON(500, gin.H{"error": "An error occurred while trying to upload your image. " + err.Error()})
      return
    }
  }

  // Send notification to all commenters if you're the creator
  if userId == post.UserId {
    var users []int64
    _, err = c.Dbmap.Select(&users, "select reply_by from replies where reply_to = ? and reply_by != ? and reply_status < 2 group by reply_by", post.Id, userId)
    if err != nil {
      ctx.JSON(500, gin.H{"error": "An error occurred while trying to send a notification. " + err.Error()})
      return
    }
    for _, targetId := range users {
      err = models.SendNotification(userId, targetId, models.NotifCommentOtherPost, post.Id, int64(0))
      if err != nil {
        ctx.JSON(500, gin.H{"error": "An error occurred while trying to send a notification. " + err.Error()})
        return
      }
    }
  } else {
    err = models.SendNotification(userId, post.UserId, models.NotifCommentMyPost, post.Id, int64(0))
    if err != nil {
      ctx.JSON(500, gin.H{"error": "An error occurred while trying to send a notification. " + err.Error()})
      return
    }
  }

  insert, err := c.Dbmap.Exec("insert into replies(reply_content, reply_feeling_id, reply_is_spoiler, reply_screenshot, reply_by, reply_to) values(?, ?, ?, ?, ?, ?)", body, feelingIdInt, isSpoilerBool, uploadedScreenshot, userId, post.Id)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while inserting your data. Please try again later. " + err.Error()})
    return
  }
  lastId, err := insert.LastInsertId()
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while inserting your data. Please try again later. " + err.Error()})
    return
  }
  // Update cache
  commentCountCached, found := c.CacheStore.Get("commentcount" + strconv.FormatInt(post.Id, 32))
  if found {
    post.CommentCount = commentCountCached.(int64)
    c.CacheStore.Set("commentcount" + strconv.FormatInt(post.Id, 32), post.CommentCount + 1, cache.DefaultExpiration)
  }
  nullScreenshot := null.NewString(uploadedScreenshot, true)
  newComment := models.Comment{Id: lastId, Content: body, Feeling: feelingIdInt, Spoiler: isSpoilerBool, Screenshot: &nullScreenshot, UserId: userId, User: user.(models.User), Date: &gorp.NullTime{Time: models.ProperNow()}}
  if(models.ProperSince(post.Date) < 3600000000000) {
    message := models.ChannelMessage{
      Topic: post.Id,
      Type: models.MessageAddComment,
    }
    if !isSpoilerBool && userId != post.UserId {
      message = models.ChannelMessage{
        Topic: post.Id,
        Type: models.MessageAddComment,
        Comment: &newComment,
      }

    }
    communityStreamId := strconv.FormatInt(post.CommunityId, 32)
    models.GetCommunityBroadcast(communityStreamId).Submit(message)
  }
  postStreamId := strconv.FormatInt(post.Id, 32)
  models.GetPostBroadcast(postStreamId).Submit(newComment)
  // Update recent comment cache if comment can be recent comment
  if !isSpoilerBool && userId != post.UserId {
    c.CacheStore.Set("recentcomment" + strconv.FormatInt(post.Id, 32), newComment, cache.DefaultExpiration)
  }
  c.CacheStore.Set("creator" + strconv.FormatInt(userId, 32), user.(models.User), cache.DefaultExpiration)
  c.CacheStore.Set("commentyeahcount" + strconv.FormatInt(lastId, 32), int64(0), cache.DefaultExpiration)
  c.CacheStore.Set("commentyeah" + strconv.FormatInt(userId, 32) + "-" + strconv.FormatInt(lastId, 32), int64(0), cache.DefaultExpiration)
  ctx.JSON(201, newComment)
}

func CommentView(ctx *gin.Context) {
  commentId := ctx.Param("id")
  var comment models.Comment
  err := c.Dbmap.SelectOne(&comment, "select reply_id, reply_to, reply_content, reply_feeling_id, reply_screenshot, reply_drawing, reply_url, reply_date, reply_edited, reply_is_spoiler, reply_by, reply_status from replies where reply_id = ? limit 1", commentId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The reply could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that reply. " + err.Error()})
    return
  }
  if comment.Status > 1 {
    switch(comment.Status) {
      case 2:
        ctx.JSON(404, gin.H{"error": "Deleted by poster."})
      default:
        userDummyRank := strings.ToLower(models.User{Rank: comment.Status}.RankText())
        ctx.JSON(404, gin.H{"error": "Deleted by " + userDummyRank + "."})
    }
    return
  }

  userCached, found := c.CacheStore.Get("creator" + strconv.FormatInt(comment.UserId, 32))
  if found {
    comment.User = userCached.(models.User)
  } else {
    c.Dbmap.SelectOne(&comment.User, "select user_id, user_name, user_rank, user_avatar from users where user_pid = ? limit 1", comment.UserId)
    c.CacheStore.Set("creator" + strconv.FormatInt(comment.UserId, 32), comment.User, cache.DefaultExpiration)
  }
  var post models.Post
  c.Dbmap.SelectOne(&post, "select post_id, post_community, post_feeling_id, post_content, post_drawing, post_by from posts where post_id = ? limit 1", comment.PostId)
  communityCached, found := c.CacheStore.Get("community" + strconv.FormatInt(post.CommunityId, 32))
  if found {
    communityType := communityCached.(models.Community)
    post.Community = &communityType
  } else {
    c.Dbmap.SelectOne(&post.Community, "select community_id, community_name, community_icon from communities where community_id = ? limit 1", post.CommunityId)
    // Not setting the cache this time since it's incomplete
  }
  userCached, found = c.CacheStore.Get("creator" + strconv.FormatInt(post.UserId, 32))
  if found {
    post.User = userCached.(models.User)
  } else {
    c.Dbmap.SelectOne(&post.User, "select user_id, user_name, user_rank, user_avatar from users where user_pid = ? limit 1", post.UserId)
    c.CacheStore.Set("creator" + strconv.FormatInt(post.UserId, 32), post.User, cache.DefaultExpiration)
  }
  user, loggedIn := ctx.Get("user")
  yeahCountCached, found := c.CacheStore.Get("commentyeahcount" + strconv.FormatInt(comment.Id, 32))
  if found {
    comment.YeahCount = yeahCountCached.(int64)
  } else {
    comment.YeahCount, _ = c.Dbmap.SelectInt("select count(ryeah_id) from reply_yeahs where ryeah_reply = ?", comment.Id)
    c.CacheStore.Set("commentyeahcount" + strconv.FormatInt(comment.Id, 32), comment.YeahCount, cache.DefaultExpiration)
  }
  if loggedIn {
    // has yeah cache
    myYeahs := int64(0)
    myYeahsCached, found := c.CacheStore.Get("commentyeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(comment.Id, 32))
    if found {
      myYeahs = myYeahsCached.(int64)
    } else {
      myYeahs, _ = c.Dbmap.SelectInt("select count(ryeah_id) from reply_yeahs where ryeah_reply = ? and ryeah_by = ?", comment.Id, user.(models.User).Id)
      c.CacheStore.Set("commentyeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(comment.Id, 32), myYeahs, cache.DefaultExpiration)
    }
    comment.HasYeah = myYeahs > 0
    comment.CanYeah = (comment.UserId != user.(models.User).Id)
  } else {
    comment.HasYeah = false
    comment.CanYeah = false
  }

  var yeahs []models.CommentYeah
  _, err = c.Dbmap.Select(&yeahs, "select ryeah_by from reply_yeahs where ryeah_reply = ? order by ryeah_id desc limit 50", comment.Id)
  if err != nil && err != sql.ErrNoRows {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to get the Yeahs for that post. " + err.Error()})
    return
  }
  for i, yeah := range yeahs {
    c.Dbmap.SelectOne(&yeahs[i].User, "select user_id, user_rank, user_avatar from users where user_pid = ? limit 1", yeah.UserId)
  }

  ctx.JSON(200, gin.H{"comment": comment, "post": post, "yeahs": yeahs})
}

func PostDelete(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to delete a post."})
    return
  }
  meUser := user.(models.User)
  postId := ctx.Param("id")
  var post models.Post
  err := c.Dbmap.SelectOne(&post, "select post_id, post_by from posts where post_id = ? and post_status < 2 limit 1", postId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The post could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that post. " + err.Error()})
    return
  }
  if post.UserId != meUser.Id && meUser.Rank < 3 {
    ctx.JSON(403, gin.H{"error": "You don't have permission to delete this post."})
    return
  }
  postStatus := 2
  if post.UserId != meUser.Id {
    postStatus = meUser.Rank
  }
  // TODO: Remove caches for this post
  _, err = c.Dbmap.Exec("update posts set post_status = ?, post_edited = now() where post_id = ?", postStatus, post.Id)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to delete that post. " + err.Error()})
    return
  }
  ctx.Status(202)
}

func CommentDelete(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to delete a reply."})
    return
  }
  meUser := user.(models.User)
  commentId := ctx.Param("id")
  var comment models.Comment
  err := c.Dbmap.SelectOne(&comment, "select reply_id, reply_by, reply_to from replies where reply_id = ? and reply_status < 2 limit 1", commentId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The reply could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that reply. " + err.Error()})
    return
  }
  if comment.UserId != meUser.Id && meUser.Rank < 3 {
    ctx.JSON(403, gin.H{"error": "You don't have permission to delete this reply."})
    return
  }
  commentStatus := 2
  if comment.UserId != meUser.Id {
    commentStatus = meUser.Rank
  }
  _, err = c.Dbmap.Exec("update replies set reply_status = ?, reply_edited = now() where reply_id = ?", commentStatus, comment.Id)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to delete that reply. " + err.Error()})
    return
  }
  // Update cache
  commentCountCached, found := c.CacheStore.Get("commentcount" + strconv.FormatInt(comment.PostId, 32))
  if found {
    commentCount := commentCountCached.(int64)
    c.CacheStore.Set("commentcount" + strconv.FormatInt(comment.PostId, 32), commentCount - 1, cache.DefaultExpiration)
  }
  c.CacheStore.Set("recentcomment" + strconv.FormatInt(comment.PostId, 32), models.Comment{}, cache.DefaultExpiration)
  ctx.Status(202)
}

func PostEdit(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to delete a post."})
    return
  }
  postId := ctx.Param("id")
  var postBefores models.PostBefores
  err := c.Dbmap.SelectOne(&postBefores, "select post_id, post_by, post_community, post_screenshot, post_content, post_content_before, post_screenshot_before, post_edited, post_date from posts where post_id = ? and post_status < 2 limit 1", postId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The post could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that post. " + err.Error()})
    return
  }
  post := postBefores.Post
  if post.UserId != user.(models.User).Id {
    ctx.JSON(403, gin.H{"error": "You don't have permission to edit this post."})
    return
  }

  body := ctx.PostForm("body")
  feelingId := ctx.PostForm("feeling_id")
  isSpoiler := ctx.PostForm("is_spoiler")
  screenshot := ctx.PostForm("screenshot")
  if body == "" && screenshot == "" {
    ctx.JSON(400, gin.H{"error": "Your post is empty."})
    return
  }
  if len([]rune(body)) > 5000 {
    ctx.JSON(400, gin.H{"error": "Your post is too long."})
    return
  }
  isSpoilerBool := false
  if isSpoiler == "true" {
    isSpoilerBool = true
  }
  feelingIdInt, _ := strconv.Atoi(feelingId)

  uploadedScreenshot := ""
  if screenshot != "" {
    uploadedScreenshot, err = util.OptimizeAndUpload(screenshot)
    if err != nil {
      ctx.JSON(500, gin.H{"error": "An error occurred while trying to upload your image. " + err.Error()})
      return
    }
  } else {
    // Screenshot will be the only one retained if it isn't specified
    uploadedScreenshot = post.Screenshot.String
  }

  var contentBefores []string
  var screenshotBefores []string
  json.Unmarshal(postBefores.ContentBefores, &contentBefores)
  json.Unmarshal(postBefores.ScreenshotBefores, &screenshotBefores)
  if post.Content != body {
    contentBefores = append(contentBefores, post.Content)
  } else {
    contentBefores = []string{}
  }
  if uploadedScreenshot != post.Screenshot.String {
    screenshotBefores = append(screenshotBefores, post.Screenshot.String)
  } else {
    screenshotBefores = []string{}
  }
  contentBeforesJson, _ := json.Marshal(contentBefores)
  screenshotBeforesJson, _ := json.Marshal(screenshotBefores)

  // TODO: Re set caches for this post
  _, err = c.Dbmap.Exec("update posts set post_content = ?, post_feeling_id = ?, post_is_spoiler = ?, post_screenshot = ?, post_content_before = ?, post_screenshot_before = ?, post_edited = now() where post_id = ?", body, feelingIdInt, isSpoilerBool, uploadedScreenshot, contentBeforesJson, screenshotBeforesJson, post.Id)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to edit that post. " + err.Error()})
    return
  }
  if(models.ProperSince(post.Date) < 3600000000000) {
    communityStreamId := strconv.FormatInt(post.CommunityId, 32)
    nullScreenshot := null.NewString(uploadedScreenshot, true)
    models.GetCommunityBroadcast(communityStreamId).Submit(models.ChannelMessage{
      Topic: post.Id,
      Type: models.MessageEditPost,
      Post: &models.Post{
        Content: body,
        Feeling: feelingIdInt,
        Spoiler: isSpoilerBool,
        Screenshot: &nullScreenshot,
      },
    })
  }
  ctx.Status(202)
}

func CommentEdit(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to delete a reply."})
    return
  }
  commentId := ctx.Param("id")
  var commentBefores models.CommentBefores
  err := c.Dbmap.SelectOne(&commentBefores, "select reply_id, reply_by, reply_to, reply_screenshot, reply_content, reply_content_before, reply_screenshot_before, reply_edited from replies where reply_id = ? and reply_status < 2 limit 1", commentId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The reply could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that reply. " + err.Error()})
    return
  }
  comment := commentBefores.Comment
  if comment.UserId != user.(models.User).Id {
    ctx.JSON(403, gin.H{"error": "You don't have permission to edit this reply."})
    return
  }

  body := ctx.PostForm("body")
  feelingId := ctx.PostForm("feeling_id")
  isSpoiler := ctx.PostForm("is_spoiler")
  screenshot := ctx.PostForm("screenshot")
  if body == "" && screenshot == "" {
    ctx.JSON(400, gin.H{"error": "Your reply is empty."})
    return
  }
  if len([]rune(body)) > 5000 {
    ctx.JSON(400, gin.H{"error": "Your reply is too long."})
    return
  }
  isSpoilerBool := false
  if isSpoiler == "true" {
    isSpoilerBool = true
  }
  feelingIdInt, _ := strconv.Atoi(feelingId)

  uploadedScreenshot := ""
  if screenshot != "" {
    uploadedScreenshot, err = util.OptimizeAndUpload(screenshot)
    if err != nil {
      ctx.JSON(500, gin.H{"error": "An error occurred while trying to upload your image. " + err.Error()})
      return
    }
  } else {
    // Screenshot will be the only one retained if it isn't specified
    uploadedScreenshot = comment.Screenshot.String
  }

  var contentBefores []string
  var screenshotBefores []string
  json.Unmarshal(commentBefores.ContentBefores, &contentBefores)
  json.Unmarshal(commentBefores.ScreenshotBefores, &screenshotBefores)
  if comment.Content != body {
    contentBefores = append(contentBefores, comment.Content)
  } else {
    contentBefores = []string{}
  }
  if uploadedScreenshot != comment.Screenshot.String {
    screenshotBefores = append(screenshotBefores, comment.Screenshot.String)
  } else {
    screenshotBefores = []string{}
  }
  contentBeforesJson, _ := json.Marshal(contentBefores)
  screenshotBeforesJson, _ := json.Marshal(screenshotBefores)

  _, err = c.Dbmap.Exec("update replies set reply_content = ?, reply_feeling_id = ?, reply_is_spoiler = ?, reply_screenshot = ?, reply_content_before = ?, reply_screenshot_before = ?, reply_edited = now() where reply_id = ?", body, feelingIdInt, isSpoilerBool, uploadedScreenshot, contentBeforesJson, screenshotBeforesJson, comment.Id)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to edit that reply. " + err.Error()})
    return
  }
  postStreamId := strconv.FormatInt(comment.PostId, 32)
  nullScreenshot := null.NewString(uploadedScreenshot, true)
  models.GetPostBroadcast(postStreamId).Submit(models.ChannelMessage{
    Topic: comment.Id,
    Type: models.MessageEditPost,
    Comment: &models.Comment{
      Content: body,
      Feeling: feelingIdInt,
      Spoiler: isSpoilerBool,
      Screenshot: &nullScreenshot,
    },
  })
  ctx.Status(202)
}

func PostGetComments(ctx *gin.Context) {
  postId := ctx.Param("id")
  var post models.Post
  err := c.Dbmap.SelectOne(&post, "select post_id from posts where post_id = ? and post_status < 2 limit 1", postId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The post could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that post. " + err.Error()})
    return
  }
  user, loggedIn := ctx.Get("user")

  offset := "0"
  limit := "20"
  timeOffset := ctx.Request.URL.Query().Get("time")
  offsetParam := ctx.Request.URL.Query().Get("offset")
  limitParam := ctx.Request.URL.Query().Get("limit")
  if offsetParam != "" {
    offsetParamInt, _ := strconv.Atoi(offsetParam)
    if offsetParamInt > 0 {
      offset = offsetParam
    }
  }
  if limitParam != "" {
    limitParamInt, _ := strconv.Atoi(limitParam)
    if limitParamInt > 0 {
      if limitParamInt > 200 {
        ctx.JSON(400, gin.H{"error": "Max limit is 200."})
        return
      }
      limit = limitParam
    }
  }
  var comments []models.Comment
  if timeOffset != "" {
    _, err = c.Dbmap.Select(&comments, "select reply_id, reply_content, reply_feeling_id, reply_screenshot, reply_drawing, reply_url, reply_date, reply_edited, reply_is_spoiler, reply_by from replies where reply_to = ? and reply_status < 2 and post_date < ? order by reply_date asc limit ? offset ?", post.Id, timeOffset, limit, offset)
  } else {
    _, err = c.Dbmap.Select(&comments, "select reply_id, reply_content, reply_feeling_id, reply_screenshot, reply_drawing, reply_url, reply_date, reply_edited, reply_is_spoiler, reply_by from replies where reply_to = ? and reply_status < 2 order by reply_date asc limit ? offset ?", post.Id, limit, offset)
  }
  if err != nil && err != sql.ErrNoRows {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to get the comments for that post. " + err.Error()})
    return
  }
  for i, comment := range comments {
    userCached, found := c.CacheStore.Get("creator" + strconv.FormatInt(comment.UserId, 32))
    if found {
      comments[i].User = userCached.(models.User)
    } else {
      c.Dbmap.SelectOne(&comments[i].User, "select user_id, user_name, user_rank, user_avatar from users where user_pid = ? limit 1", comment.UserId)
      c.CacheStore.Set("creator" + strconv.FormatInt(comment.UserId, 32), comments[i].User, cache.DefaultExpiration)
    }
    yeahCountCached, found := c.CacheStore.Get("commentyeahcount" + strconv.FormatInt(comment.Id, 32))
    if found {
      comments[i].YeahCount = yeahCountCached.(int64)
    } else {
      comments[i].YeahCount, _ = c.Dbmap.SelectInt("select count(ryeah_id) from reply_yeahs where ryeah_reply = ?", comment.Id)
      c.CacheStore.Set("commentyeahcount" + strconv.FormatInt(comment.Id, 32), comments[i].YeahCount, cache.DefaultExpiration)
    }
    if loggedIn {
      // has yeah cache
      myYeahs := int64(0)
      myYeahsCached, found := c.CacheStore.Get("commentyeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(comment.Id, 32))
      if found {
        myYeahs = myYeahsCached.(int64)
      } else {
        myYeahs, _ := c.Dbmap.SelectInt("select count(ryeah_id) from reply_yeahs where ryeah_reply = ? and ryeah_by = ?", comment.Id, user.(models.User).Id)
        c.CacheStore.Set("commentyeah" + strconv.FormatInt(user.(models.User).Id, 32) + "-" + strconv.FormatInt(comment.Id, 32), myYeahs, cache.DefaultExpiration)
      }
      comments[i].HasYeah = myYeahs > 0
      comments[i].CanYeah = (comment.UserId != user.(models.User).Id)
    } else {
      comments[i].HasYeah = false
      comments[i].CanYeah = false
    }
  }

  ctx.JSON(200, comments)
}
