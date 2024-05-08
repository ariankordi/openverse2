/** Some things in post/reply/user need to be changed to NULL
  * Add reply_drawing to reply
  * Add post_content_before, post_screenshot_before to post, reply_content_before, post_screenshot_before to reply
  * Add notif_by_others as text to notification, default null
  * Add index to notification on notif_to
*/
package main

import (
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/utrack/gin-csrf"

  "github.com/patrickmn/go-cache"

  "openverse/routes"
  "openverse/models"
  "openverse/config"

  "strings"
  "strconv"
  "encoding/json"
  "io/ioutil"
  "bytes"

  "os"
  /*"os/user"
  "net"
  "net/http"*/

  "log"
)


func main() {
  config.Dbmap.TraceOn("[GORP]", log.New(os.Stdout, "", 0))

  router := gin.Default()
  router.Use(sessions.Sessions("openverse_0", config.SessionStore))
  router.Use(csrf.Middleware(csrf.Options{
    Secret: "C`v/$2-!YLBK*Vg>^-.bh?zJ65)>BWuEWFDy),M2ySL%W9~3vFgZ+U</?4HwvUe,",
    ErrorFunc: func(ctx *gin.Context) {
      ctx.JSON(403, gin.H{"error": "CSRF token mismatch"})
      ctx.Abort()
    },
  }))
  router.Use(routes.UserSession)
  router.POST("/app/login", routes.Login)
  router.POST("/app/logout", routes.Logout)
  router.POST("/app/signup", routes.Signup)

  router.GET("/app/titles", routes.TitlesList)
  router.GET("/app/communities", routes.CommunitiesList)
  router.GET("/app/communities/:id", routes.CommunityPosts)
  router.POST("/app/communities/:id/post", routes.CommunityCreatePost)

  router.GET("/app/posts/:id", routes.PostView)
  router.POST("/app/posts/:id/yeah", routes.PostAddYeah)
  router.POST("/app/posts/:id/unyeah", routes.PostDeleteYeah)
  router.POST("/app/posts/:id/reply", routes.PostAddComment)
  router.POST("/app/replies/:id/yeah", routes.CommentAddYeah)
  router.POST("/app/replies/:id/unyeah", routes.CommentDeleteYeah)
  router.GET("/app/replies/:id", routes.CommentView)
  router.POST("/app/posts/:id/delete", routes.PostDelete)
  router.POST("/app/replies/:id/delete", routes.CommentDelete)
  router.POST("/app/posts/:id/edit", routes.PostEdit)
  router.POST("/app/replies/:id/edit", routes.CommentEdit)
  router.GET("/app/posts/:id/replies", routes.PostGetComments)

  router.GET("/app/users/:id", routes.UserView)

  router.GET("/app/news/my_news", routes.GetNotifications)
  router.GET("/app/settings/profile", routes.GetProfileSettings)
  router.POST("/app/settings/profile", routes.EditProfileSettings)

  router.GET("/app/communities/:id/stream", routes.CommunityPostsStream)
  router.GET("/app/posts/:id/stream", routes.PostStream)
  router.GET("/app/news/my_news/stream", routes.GetNotificationsStream)

  // Debug, remove this soon
  router.GET("/app/flush", func(ctx *gin.Context) {
    config.CacheStore.Flush()
    ctx.String(200, "Store flushed")
  })

  router.Static("/static", "./static")
  // Vue index.html
  router.NoRoute(func(c *gin.Context) {
    // Return 404 if the browser isn't looking for an HTML page
    acceptHeader := c.GetHeader("Accept")
    if !strings.Contains(acceptHeader, "html") && !strings.HasPrefix(acceptHeader, "*/*") {
      c.Status(404)
      return
    }
    // Give 404 for /app/
    if strings.HasPrefix(c.Request.URL.Path, "/app/") {
      c.Status(404)
      return
    }
    user, exists := c.Get("user")
    if exists != true {
      user = &models.User{}
    }
    userData, err := json.Marshal(user)
    if err != nil {
      userData = []byte("null")
    }
    csrfToken := []byte(csrf.GetToken(c))

    unreadNotifications := int64(0)
    if meUser, ok := user.(models.User); ok && meUser.Id > 0 {
      unreadNotifications, _ = config.Dbmap.SelectInt("select count(notif_id) from notifications where notif_to = ? and notif_read = 0", meUser.Id)
    }

    // Get cookie for dark preference
    darkPreference := []byte(" disabled")
    darkCookie, _ := c.Request.Cookie("openverse_dark")
    if darkCookie != nil {
      darkPreference = nil
    }

    // index.html cache
    html := []byte{}
    htmlCached, found := config.CacheStore.Get("index")
    if found {
      html = htmlCached.([]byte)
    } else {
      html, err = ioutil.ReadFile("./index.html")
      if err != nil {
        html = []byte(err.Error())
      } else {
        config.CacheStore.Set("index", html, cache.NoExpiration)
      }
    }

    template := bytes.Replace(html, []byte("__USER_DATA"), userData, 1)
    template = bytes.Replace(template, []byte("__CSRF_TOKEN"), csrfToken, 1)
    template = bytes.Replace(template, []byte("__UNREAD_NOTIFS"), []byte(strconv.FormatInt(unreadNotifications, 10)), 1)
    template = bytes.Replace(template, []byte("__DARK_PREFERENCE"), darkPreference, 1)
    //template := mustache.RenderFile("./index.html", user)
    c.Data(200, "text/html", template)
  })

  // Manually running gin, so we can use it with a sock and permissions
  /*os.Remove("./run.sock")
  listener, err := net.Listen("unix", "./run.sock")
	if err != nil {
		log.Fatal(err)
	}
  // Change the group of the sock to nginx so nginx can use it
  nginxUser, err := user.Lookup("nginx")
  if err == nil {
    // This gid is being converted from a string now but it's going to be an int in Go 2
    nginxUserUid, _ := strconv.Atoi(nginxUser.Uid)
    nginxUserGid, _ := strconv.Atoi(nginxUser.Gid)
    os.Chown("./run.sock", nginxUserUid, nginxUserGid)
  } else {
    log.Println(err)
  }
	defer listener.Close()
  // Finally, serve
	log.Fatal(http.Serve(listener, router))*/
	log.Fatal(router.Run(":8082"))
}
