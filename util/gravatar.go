package util

import (
  "github.com/valyala/fasthttp"
  "crypto/md5"
  "encoding/hex"
)

func GetGravatar(email string) string {
  hash := md5.New()
  hash.Write([]byte(email))
  emailHash := hex.EncodeToString(hash.Sum(nil))
  url := "https://gravatar.com/avatar/" + emailHash + "?d=404&s=0."

  req := fasthttp.AcquireRequest()
  req.SetRequestURI(url)
  resp := fasthttp.AcquireResponse()
  client.Do(req, resp)

  if resp.StatusCode() != 200 {
    return ""
  }
  return "https://gravatar.com/avatar/" + emailHash + "?d=404&s=128"
}
