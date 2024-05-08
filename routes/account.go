package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"

  "github.com/patrickmn/go-cache"

  "golang.org/x/crypto/bcrypt"

  "database/sql"
  "strconv"
  "regexp"

  "openverse/util"
  "openverse/models"
  c "openverse/config"
)

func UserSession(ctx *gin.Context) {
  session := sessions.Default(ctx)
  ctx.Set("session", session)

  userId := session.Get("user")
  if userId != nil {
    var user models.User
    userCached, found := c.CacheStore.Get("creator" + strconv.FormatInt(userId.(int64), 32))
    if found {
      user = userCached.(models.User)
    } else {
      err := c.Dbmap.SelectOne(&user, "select user_id, user_name, user_rank, user_avatar from users where user_pid = ? limit 1", userId)
      if err != nil {
        return
      }
      c.CacheStore.Set("creator" + strconv.FormatInt(userId.(int64), 32), user, cache.DefaultExpiration)
    }
    user.Id = userId.(int64)
    ctx.Set("user", user)
  }

  ctx.Next()
}

func Login(ctx *gin.Context) {
  Name := ctx.PostForm("name")
  Pass := ctx.PostForm("pass")
  if Name == "" || Pass == "" {
    ctx.JSON(400, gin.H{"error": "The Username/Password fields cannot be blank."})
    return
  }
  var user models.User
  err := c.Dbmap.SelectOne(&user, "select user_pid, user_pass from users where user_id = ? limit 1", Name)
  if err != nil {
    if err == sql.ErrNoRows {
      ctx.JSON(404, gin.H{"error": "That user could not be found."})
      return
    }
    ctx.JSON(500, gin.H{"error": "Something went wrong while signing in. " + err.Error()})
    return
  }
  err = bcrypt.CompareHashAndPassword(user.Pass, []byte(Pass))
  if err != nil {
    if err == bcrypt.ErrMismatchedHashAndPassword {
      ctx.JSON(401, gin.H{"error": "Invalid password."})
      return
    }
    ctx.JSON(500, gin.H{"error": "Something went wrong while verifying your password. " + err.Error()})
    return
  }
  session := sessions.Default(ctx)
  session.Set("user", user.Id)
  session.Save()
  ctx.JSON(202, gin.H{"id": user.Id})
}

func Signup(ctx *gin.Context) {
  Name := ctx.PostForm("name")
  Nick := ctx.PostForm("nick")
  Nnid := ctx.PostForm("nnid")
  Email := ctx.PostForm("email")
  Avatar := ctx.PostForm("avatar")
  Pass := ctx.PostForm("pass")
  PassAgain := ctx.PostForm("pass_again")
  if Name == "" || Pass == "" || PassAgain == "" {
    ctx.JSON(400, gin.H{"error": "The Username/Password/Confirm Password fields cannot be blank."})
    return
  }
  if Pass != PassAgain {
    ctx.JSON(400, gin.H{"error": "The passwords you entered do not match."})
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
  if Nick == "" && Nnid == "" {
    ctx.JSON(400, gin.H{"error": "The Nickname field cannot be blank."})
    return
  }
  if Nick != "" && len([]rune(Nick)) > 64 {
    ctx.JSON(400, gin.H{"error": "The nickname you entered is too long."})
    return
  }
  if Nnid != "" {
    nnidSyntaxValid, err := regexp.MatchString("^[A-Za-z0-9-._]{6,16}$", Nnid)
    if err != nil {
      ctx.JSON(500, gin.H{"error": "Something went wrong while verifying your account. " + err.Error()})
      return
    }
    if !nnidSyntaxValid {
      ctx.JSON(400, gin.H{"error": "The Nintendo Network ID you entered is invalid."})
      return
    }
  }
  if Email != "" {
    emailSyntaxValid, err := regexp.MatchString("^[A-Za-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$", Email)
    if err != nil {
      ctx.JSON(500, gin.H{"error": "Something went wrong while verifying your account. " + err.Error()})
      return
    }
    if !emailSyntaxValid {
      ctx.JSON(400, gin.H{"error": "The email address you entered is invalid."})
      return
    }
  }
  // End of form validation
  nameExists, err := c.Dbmap.SelectInt("select count(user_pid) from users where user_id = ?", Name)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "A database error occurred while trying to verify your account. " + err.Error()})
    return
  }
  if nameExists > 0 {
    ctx.JSON(400, gin.H{"error": "The username you entered is already taken."})
    return
  }
  if Email != "" {
    emailExists, err := c.Dbmap.SelectInt("select count(user_pid) from users where user_email = ?", Email)
    if err != nil {
      ctx.JSON(500, gin.H{"error": "A database error occurred while trying to verify your account. " + err.Error()})
      return
    }
    if emailExists > 0 {
      ctx.JSON(400, gin.H{"error": "The email you entered is already registered to another account."})
      return
    }
  }
  if Nnid != "" {
    nnidExists, err := c.Dbmap.SelectInt("select count(user_pid) from users where user_nnid = ?", Nnid)
    if err != nil {
      ctx.JSON(500, gin.H{"error": "A database error occurred while trying to verify your account. " + err.Error()})
      return
    }
    if nnidExists > 0 {
      ctx.JSON(400, gin.H{"error": "The Nintendo Network ID you entered is already registered to another account."})
      return
    }
  }
  // End of form checks
  finalAvatar := ""
  if Avatar != "" {
    finalAvatar, err = util.OptimizeAndUpload128(Avatar)
    if err != nil {
      ctx.JSON(500, gin.H{"error": "An error occurred while trying to upload your image. " + err.Error()})
      return
    }
  } else {
    if Email != "" {
      gravatarAvatar := util.GetGravatar(Email)
      if gravatarAvatar != "" {
        finalAvatar = gravatarAvatar
      }
    }
  }
  if Nnid != "" {
    getNnid, err := util.GetNNIDInfo(Nnid)
    if err != nil {
      if err.Error() == "NNID doesn't exist" {
        ctx.JSON(400, gin.H{"error": "The Nintendo Network ID you entered does not exist."})
        return
      }
      ctx.JSON(500, gin.H{"error": "Something went wrong while getting your Nintendo Network ID. " + err.Error()})
      return
    }
    if finalAvatar == "" {
      finalAvatar = getNnid.Hash
    }
    Nnid = getNnid.Nnid
    if Nick == "" {
      Nick = getNnid.Name
    }
  }

  newPass, err := bcrypt.GenerateFromPassword([]byte(Pass), 13)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "Something went wrong while creating a password. " + err.Error()})
    return
  }
  insert, err := c.Dbmap.Exec("insert into users(user_id, user_name, user_pass, user_email, user_avatar, user_ip, user_nnid) values(?, ?, ?, ?, ?, ?, ?)", Name, Nick, string(newPass), Email, finalAvatar, ctx.ClientIP(), Nnid)
  if err != nil {
    ctx.JSON(500, gin.H{"error": "A database error occurred while trying to create your account. " + err.Error()})
    return
  }
  userId, err := insert.LastInsertId()
  if err != nil {
    ctx.JSON(500, gin.H{"error": "A database error occurred while trying to create your account. " + err.Error()})
    return
  }
  session := sessions.Default(ctx)
  session.Set("user", userId)
  session.Save()
  ctx.Status(201)
}

func Logout(ctx *gin.Context) {
  session := sessions.Default(ctx)
  session.Set("user", nil)
  session.Save()
  ctx.Status(202)
}
