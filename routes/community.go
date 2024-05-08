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
  "time"

  "openverse/util"
  "openverse/models"
  c "openverse/config"
)

func TitlesList(ctx *gin.Context) {
  var generalTitles []models.Title
  var gameTitles []models.Title
  var specialTitles []models.Title
  //_, err := c.Dbmap.Select(&titles, "SELECT titles1.title_id, titles1.title_type, titles1.title_icon, titles1.title_platform FROM titles titles1 LEFT OUTER JOIN titles titles2 ON (titles1.title_type = titles2.title_type AND titles1.title_id < titles2.title_id) GROUP BY titles1.title_id HAVING COUNT(titles1.title_id) < 8 ORDER BY title_type, title_id DESC")
  _, err := c.Dbmap.Select(&generalTitles, "select title_id, title_name, title_type, title_icon, title_platform from titles where title_type = 1 order by title_id desc limit 8")
  if err != nil {
    ctx.JSON(500, gin.H{"error": "The database server is giving errors right now, try again in a moment. " + err.Error()})
    return
  }
  _, err = c.Dbmap.Select(&gameTitles, "select title_id, title_name, title_type, title_icon, title_platform from titles where title_type = 0 order by title_id desc limit 8")
  if err != nil {
    ctx.JSON(500, gin.H{"error": "The database server is giving errors right now, try again in a moment. " + err.Error()})
    return
  }
  _, err = c.Dbmap.Select(&specialTitles, "select title_id, title_name, title_type, title_icon, title_platform from titles where title_type = 2 order by title_id desc limit 8")
  if err != nil {
    ctx.JSON(500, gin.H{"error": "The database server is giving errors right now, try again in a moment. " + err.Error()})
    return
  }
  ctx.JSON(200, gin.H{"general": generalTitles, "game": gameTitles, "special": specialTitles})
}

func CommunitiesList(ctx *gin.Context) {
  var generalCommunities []models.Community
  var gameCommunities []models.Community
  var specialCommunities []models.Community
  cachedGeneral, found := c.CacheStore.Get("general")
  if found {
    generalCommunities = cachedGeneral.([]models.Community)
  } else {
    _, err := c.Dbmap.Select(&generalCommunities, "select community_id, community_name, community_type, community_icon, community_platform from communities where community_type = 1 order by community_id desc limit 12")
    if err != nil {
      ctx.JSON(500, gin.H{"error": "The database server is giving errors right now, try again in a moment. " + err.Error()})
      return
    }
    c.CacheStore.Set("general", generalCommunities, 10 * time.Minute)
  }
  cachedGame, found := c.CacheStore.Get("game")
  if found {
    gameCommunities = cachedGame.([]models.Community)
  } else {
    _, err := c.Dbmap.Select(&gameCommunities, "select community_id, community_name, community_type, community_icon, community_platform from communities where community_type = 0 order by community_id desc limit 12")
    if err != nil {
      ctx.JSON(500, gin.H{"error": "The database server is giving errors right now, try again in a moment. " + err.Error()})
      return
    }
    c.CacheStore.Set("game", gameCommunities, 10 * time.Minute)
  }
  cachedSpecial, found := c.CacheStore.Get("special")
  if found {
    specialCommunities = cachedSpecial.([]models.Community)
  } else {
    _, err := c.Dbmap.Select(&specialCommunities, "select community_id, community_name, community_type, community_icon, community_platform from communities where community_type = 2 order by community_id desc limit 12")
    if err != nil {
      ctx.JSON(500, gin.H{"error": "The database server is giving errors right now, try again in a moment. " + err.Error()})
      return
    }
    c.CacheStore.Set("special", specialCommunities, 10 * time.Minute)
  }
  ctx.JSON(200, gin.H{"general": generalCommunities, "game": gameCommunities, "special": specialCommunities})
}

func CommunityPosts(ctx *gin.Context) {
  communityId := ctx.Param("id")
  var community models.Community
  communityIdInt, err := strconv.ParseInt(communityId, 10, 64)
  if err != nil {
    ctx.JSON(400, gin.H{"error": "You have entered an invalid community ID."})
    return
  }
  communityCached, found := c.CacheStore.Get("community" + strconv.FormatInt(communityIdInt, 32))
  if found {
    community = communityCached.(models.Community)
  } else {
    err := c.Dbmap.SelectOne(&community, "select community_id, community_name, community_description, community_type, community_icon, community_banner, community_platform, community_perms from communities where community_id = ? limit 1", communityId)
    if err != nil {
      if err == sql.ErrNoRows {
        ctx.JSON(404, gin.H{"error": "The community could not be found."})
        return
      }
      ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that community. " + err.Error()})
      return
    }
    c.CacheStore.Set("community" + strconv.FormatInt(communityIdInt, 32), community, cache.DefaultExpiration)
  }
  // TODO: Community permissions
  user, loggedIn := ctx.Get("user")
  offset := "0"
  limit := "50"
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
      if limitParamInt > 500 {
        ctx.JSON(400, gin.H{"error": "Max limit is 500."})
        return
      }
      limit = limitParam
    }
  }
  var posts []models.Post
  if timeOffset != "" {
    _, err = c.Dbmap.Select(&posts, "select post_id, post_feeling_id, post_content, post_screenshot, post_drawing, post_url, post_is_spoiler, post_date, post_edited, post_by from posts where post_community = ? and post_status < 2 and post_date < ? order by post_id desc limit ? offset ?", community.Id, timeOffset, limit, offset)
  } else {
    _, err = c.Dbmap.Select(&posts, "select post_id, post_feeling_id, post_content, post_screenshot, post_drawing, post_url, post_is_spoiler, post_date, post_edited, post_by from posts where post_community = ? and post_status < 2 order by post_id desc limit ? offset ?", community.Id, limit, offset)
  }
  if err != nil && err != sql.ErrNoRows {
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to grab the posts for that community. " + err.Error()})
    return
  }
  for i, post := range posts {
    // creator cache
    userCached, found := c.CacheStore.Get("creator" + strconv.FormatInt(post.UserId, 32))
    if found {
      posts[i].User = userCached.(models.User)
    } else {
      c.Dbmap.SelectOne(&posts[i].User, "select user_id, user_name, user_rank, user_avatar from users where user_pid = ? limit 1", post.UserId)
      c.CacheStore.Set("creator" + strconv.FormatInt(post.UserId, 32), posts[i].User, cache.DefaultExpiration)
    }
    // yeah count cache
    yeahCountCached, found := c.CacheStore.Get("yeahcount" + strconv.FormatInt(post.Id, 32))
    if found {
      posts[i].YeahCount = yeahCountCached.(int64)
    } else {
      posts[i].YeahCount, _ = c.Dbmap.SelectInt("select count(yeah_id) from post_yeahs where yeah_post = ?", post.Id)
      c.CacheStore.Set("yeahcount" + strconv.FormatInt(post.Id, 32), posts[i].YeahCount, cache.DefaultExpiration)
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
      posts[i].HasYeah = myYeahs > 0
      posts[i].CanYeah = (post.UserId != user.(models.User).Id)
    } else {
      posts[i].HasYeah = false
      posts[i].CanYeah = false
    }
    // comment count cache
    commentCountCache, found := c.CacheStore.Get("commentcount" + strconv.FormatInt(post.Id, 32))
    if found {
      posts[i].CommentCount = commentCountCache.(int64)
    } else {
      posts[i].CommentCount, _ = c.Dbmap.SelectInt("select count(reply_id) from replies where reply_to = ? and reply_status < 2", post.Id)
      c.CacheStore.Set("commentcount" + strconv.FormatInt(post.Id, 32), posts[i].CommentCount, cache.DefaultExpiration)
    }
    if posts[i].CommentCount > 0 {
      var recentComment models.Comment
      // recent comment cache
      recentCommentCache, found := c.CacheStore.Get("recentcomment" + strconv.FormatInt(post.Id, 32))
      if found {
        recentComment = recentCommentCache.(models.Comment)
      } else {
        c.Dbmap.SelectOne(&recentComment, "select reply_id, reply_feeling_id, reply_content, reply_drawing, reply_screenshot, reply_by, reply_date, reply_edited from replies where reply_to = ? and reply_by != ? and reply_status < 2 and reply_is_spoiler = 0 order by reply_id desc limit 1", post.Id, post.UserId)
        c.CacheStore.Set("recentcomment" + strconv.FormatInt(post.Id, 32), recentComment, cache.DefaultExpiration)
      }
      if recentComment.Id > 0 {
        // recent comment creator cache
        recentCommentUserCached, found := c.CacheStore.Get("creator" + strconv.FormatInt(recentComment.UserId, 32))
        if found {
          recentComment.User = recentCommentUserCached.(models.User)
        } else {
          c.Dbmap.SelectOne(&recentComment.User, "select user_id, user_name, user_rank, user_avatar from users where user_pid = ? limit 1", recentComment.UserId)
          c.CacheStore.Set("creator" + strconv.FormatInt(recentComment.UserId, 32), recentComment.User, cache.DefaultExpiration)
        }
        posts[i].RecentComment = &recentComment
      }
    }
  }
  ctx.JSON(200, gin.H{"community": community, "posts": posts})
}

func CommunityPostsStream(ctx *gin.Context) {
  communityId := ctx.Param("id")
  var community models.Community
  err := c.Dbmap.SelectOne(&community, "select community_id from communities where community_id = ? limit 1", communityId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The community could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that community. " + err.Error()})
    return
  }
  // The stream ID is base32
  communityStreamId := strconv.FormatInt(community.Id, 32)
  listener := models.OpenCommunityListener(communityStreamId)
	defer models.CloseCommunityListener(communityStreamId, listener)

  conn, _, _, err := ws.UpgradeHTTP(ctx.Request, ctx.Writer)
  if err != nil {
    return
  }
  defer conn.Close()
  for {
    recieved := <-listener
    if recieved != nil {
      message := ""
      if post, ok := recieved.(models.Post); ok {
        encodedPost, _ := json.Marshal(post)
        message = "post:" + string(encodedPost)
      } else {
        recievedMessage := recieved.(models.ChannelMessage)
        encodedMessage, _ := json.Marshal(recievedMessage)
        switch(recievedMessage.Type) {
        case models.MessageYeahPost:
            message = "yeah:" + string(encodedMessage)
          case models.MessageUnyeahPost:
            message = "unyeah:" + string(encodedMessage)
          case models.MessageAddComment:
            message = "comment:" + string(encodedMessage)
          case models.MessageEditPost:
            message = "edit:" + string(encodedMessage)
        }
      }
      wsutil.WriteServerMessage(conn, ws.OpText, []byte(message))
    }
  }
}

func CommunityCreatePost(ctx *gin.Context) {
  user, loggedIn := ctx.Get("user")
  if !loggedIn {
    ctx.JSON(500, gin.H{"error": "You must be signed in to make a post."})
    return
  }
  userId := user.(models.User).Id
  communityId := ctx.Param("id")
  var community models.Community
  err := c.Dbmap.SelectOne(&community, "select community_id, community_name, community_description, community_type, community_icon, community_banner, community_platform from communities where community_id = ? limit 1", communityId)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "The community could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "An error occurred while trying to find that community. " + err.Error()})
    return
  }
  if user.(models.User).Rank < community.Permission {
    ctx.JSON(404, gin.H{"error": "You don't have permission to make a post to that community."})
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
  }
  insert, err := c.Dbmap.Exec("insert into posts(post_content, post_feeling_id, post_is_spoiler, post_screenshot, post_by, post_community) values(?, ?, ?, ?, ?, ?)", body, feelingIdInt, isSpoilerBool, uploadedScreenshot, userId, community.Id)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while inserting your data. Please try again later. " + err.Error()})
    return
  }
  lastId, err := insert.LastInsertId()
  if err != nil {
    ctx.JSON(500, gin.H{"error": "An error occurred while inserting your data. Please try again later. " + err.Error()})
    return
  }
  communityStreamId := strconv.FormatInt(community.Id, 32)
  nullScreenshot := null.NewString(uploadedScreenshot, true)
  newPost := models.Post{Id: lastId, Content: body, Feeling: feelingIdInt, Spoiler: isSpoilerBool, Screenshot: &nullScreenshot, UserId: userId, User: user.(models.User), Date: &gorp.NullTime{Time: models.ProperNow()}/*, CommunityId: community.Id, Community: &community*/}
  models.GetCommunityBroadcast(communityStreamId).Submit(newPost)
  c.CacheStore.Set("creator" + strconv.FormatInt(userId, 32), user.(models.User), cache.DefaultExpiration)
  c.CacheStore.Set("yeahcount" + strconv.FormatInt(lastId, 32), int64(0), cache.DefaultExpiration)
  c.CacheStore.Set("yeah" + strconv.FormatInt(userId, 32) + "-" + strconv.FormatInt(lastId, 32), int64(0), cache.DefaultExpiration)
  c.CacheStore.Set("commentcount" + strconv.FormatInt(lastId, 32), int64(0), cache.DefaultExpiration)
  c.CacheStore.Set("recentcomment" + strconv.FormatInt(lastId, 32), models.Comment{}, cache.DefaultExpiration)
  ctx.JSON(201, newPost)
}
