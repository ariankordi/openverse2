package util

import (
  "github.com/valyala/fasthttp"
  "github.com/h2non/bimg"

  "mime/multipart"
  "encoding/base64"
  "encoding/json"

  "strings"
  "strconv"
  "bytes"
  "io"
  "io/ioutil"
  "errors"
)

func UploadToKek(image []byte, ext string) (string, error) {
  imageEncoded := "data:image/" + ext + ";base64," + base64.StdEncoding.EncodeToString(image)
  req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://u.kek.gg/v1/upload-to-kek")
	req.Header.SetContentType("text/plain")
	req.Header.SetMethod("POST")

	req.SetBodyString(imageEncoded)
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	if err != nil {
		return "", err
	}

  if resp.StatusCode() > 300 {
    return "", errors.New("http response code " + strconv.Itoa(resp.StatusCode()) + ", " + string(resp.Body()))
  }

	body := string(resp.Body())
  body = strings.Replace(body, "?", "", 1)
  body, err = strconv.Unquote(body)
  if err != nil {
    return "", err
  }
	return body, nil
}

func UploadToCatgirl(image []byte, ext string) (string, error) {
  /*imageSlice := strings.SplitAfter(image, ";base64,")
  if len(imageSlice) < 1 {
    return "", errors.New("no image mime type on given string")
  }
  fileExt := strings.Split(strings.Split(image, ";base64,")[0], "data:image/")[1]
  base64Decoded := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imageSlice[1]))
  */
  imageReader := bytes.NewReader(image)
  body := &bytes.Buffer{}
  writer := multipart.NewWriter(body)
  part, err := writer.CreateFormFile("sharex", "." + ext)
  if err != nil {
      return "", err
  }
  _, err = io.Copy(part, imageReader)
  if err != nil {
      return "", err
  }
  _ = writer.WriteField("secret", "8ufLPswPPSeMbMiF37fehoNJTLlbXRbtK1nhxE6aPhqR78LV7W8XdonOvsRp9lS9")
  err = writer.Close()
  if err != nil {
      return "", err
  }
  bodyRequest, err := ioutil.ReadAll(body)
  if err != nil {
      return "", err
  }

  req := fasthttp.AcquireRequest()
  req.SetRequestURI("https://catgirl.host/upload")
  req.Header.SetMethod("POST")
  req.Header.SetContentType(writer.FormDataContentType())

  req.SetBody(bodyRequest)
	resp := fasthttp.AcquireResponse()
	err = client.Do(req, resp)
	if err != nil {
		return "", err
	}

  if resp.StatusCode() > 300 {
    return "", errors.New("http response code " + strconv.Itoa(resp.StatusCode()) + ", " + string(resp.Body()))
  }

	bodyResp := string(resp.Body())
  bodyResp = strings.Replace(bodyResp, "i.", "u.", 1)
  bodyResp = bodyResp + "." + ext
  if err != nil {
    return "", err
  }
	return bodyResp, nil
}

func UploadToLolisafe(image []byte, ext string) (string, error) {
  imageReader := bytes.NewReader(image)
  body := &bytes.Buffer{}
  writer := multipart.NewWriter(body)
  part, err := writer.CreateFormFile("files[]", " ." + ext)
  if err != nil {
      return "", err
  }
  _, err = io.Copy(part, imageReader)
  if err != nil {
      return "", err
  }
  err = writer.Close()
  if err != nil {
      return "", err
  }
  bodyRequest, err := ioutil.ReadAll(body)
  if err != nil {
      return "", err
  }

  req := fasthttp.AcquireRequest()
  // Other lolisafe servers like lolisareinthe.club, dmca.gripe also work
  req.SetRequestURI("https://safe.moe/api/upload")
  req.Header.SetMethod("POST")
  req.Header.SetContentType(writer.FormDataContentType())

  req.SetBody(bodyRequest)
	resp := fasthttp.AcquireResponse()
	err = client.Do(req, resp)
	if err != nil {
		return "", err
	}

  if resp.StatusCode() > 300 {
    return "", errors.New("http response code " + strconv.Itoa(resp.StatusCode()) + ", " + string(resp.Body()))
  }

  var bodyDecoded map[string]interface{}
	err = json.Unmarshal(resp.Body(), &bodyDecoded)
  if err != nil {
    return "", err
  }
  imageUrl := bodyDecoded["files"].([]interface{})[0].(map[string]interface{})["url"].(string)
	return imageUrl, nil
}

func OptimizeAndUpload(image string) (string, error) {
  imageSlice := strings.SplitAfter(image, ";base64,")
  if len(imageSlice) < 1 {
    return "", errors.New("no image mime type on given string")
  }
  fileExt := strings.Split(strings.Split(image, ";base64,")[0], "data:image/")[1]
  base64Decoded, err := base64.StdEncoding.DecodeString(imageSlice[1])
  if err != nil {
    return "", err
  }
  newImage := base64Decoded
  if fileExt != "gif" {
    options := bimg.Options{
      Quality: 85,
      Compression: 9,
      StripMetadata: true,
    }
    newImage, err = bimg.NewImage(base64Decoded).Process(options)
    if err != nil {
      return "", err
    }
  }

  kek, err := UploadToLolisafe(newImage, fileExt)
  return kek, err
}

func OptimizeAndUpload128(image string) (string, error) {
  imageSlice := strings.SplitAfter(image, ";base64,")
  if len(imageSlice) < 1 {
    return "", errors.New("no image mime type on given string")
  }
  fileExt := strings.Split(strings.Split(image, ";base64,")[0], "data:image/")[1]
  base64Decoded, err := base64.StdEncoding.DecodeString(imageSlice[1])
  if err != nil {
    return "", err
  }
  newImage := base64Decoded
  if fileExt == "gif" {
    return "", errors.New("Sorry, but you cannot use a GIF as an avatar.")
  }
  options := bimg.Options{
    Quality: 85,
    Compression: 9,
    StripMetadata: true,
  }
  newImage, err = bimg.NewImage(base64Decoded).Process(options)
  if err != nil {
    return "", err
  }
  newImage, err = bimg.NewImage(newImage).Resize(128, 128)
  if err != nil {
    return "", err
  }

  kek, err := UploadToLolisafe(newImage, fileExt)
  return kek, err
}
