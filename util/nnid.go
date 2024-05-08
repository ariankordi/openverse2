package util

import (
  "github.com/valyala/fasthttp"

  "encoding/json"

  "errors"
  "strings"
)

var client = &fasthttp.Client{}

type MappedID struct {
	InID string `json:"in_id" db:"user_id"`
	OutID string `json:"out_id" db:"pid"`
}
type MappedIDsResponse struct {
	MappedIDs []MappedID `json:"mapped_ids"`
}
// ID is omitted from this
type Mii struct {
	Name string `json:"name"`
	Data string `json:"data"`
	UserID string `json:"user_id"`
	Images []MiiImage `json:"images"`
}
type MiiImage struct {
	Type string `json:"type"`
	//URL string `json:"url"`
	CachedURL string `json:"cached_url"`
}
type MiiResponse struct {
	Miis []Mii `json:"miis"`
}

type MiiInfo struct {
  Nnid string `json:"nnid"`
  Name string `json:"name"`
  Hash string `json:"hash"`
  Data string `json:"data"`
}

func doRequest(url string) []byte {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Nintendo-Client-ID", "a2efa818a34fa16b8afbc8a74eba3eda")
	req.Header.Add("X-Nintendo-Client-Secret", "c91cdb5658bd4954ade78533a339cf9a")

	resp := fasthttp.AcquireResponse()
	client.Do(req, resp)

	return resp.Body()
}
func GetNNID(nnid string) (Mii, error) {
	endpoint := "https://accountws.nintendo.net/v1/api/admin/mapped_ids?input_type=user_id&output_type=pid&input=" + nnid
	response := doRequest(endpoint)
	var MappedIDs = new(MappedIDsResponse)
	json.Unmarshal(response, &MappedIDs)
	if MappedIDs.MappedIDs[0].OutID == "" {
		return Mii{}, errors.New("NNID doesn't exist")
	}
	endpoint = "https://accountws.nintendo.net/v1/api/miis?pids=" + MappedIDs.MappedIDs[0].OutID
  response = doRequest(endpoint)
	if string(response) == "" {
		return Mii{}, errors.New("Mii doesn't exist")
	}
	var Miis = new(MiiResponse)
	json.Unmarshal(response, &Miis)
	return Miis.Miis[0], nil
}

// Shortcut to get a Mii hash and names from GetNNID
func GetNNIDInfo(nnid string) (MiiInfo, error) {
  nnidAll, err := GetNNID(nnid)
  if err != nil {
    return MiiInfo{}, err
  }
  fullUrl := nnidAll.Images[0].CachedURL
  fullUrl1 := strings.Split(fullUrl, "/")[3]
  miiHash := strings.Split(fullUrl1, "_")[0]
  return MiiInfo{
    Name: nnidAll.Name,
    Data: nnidAll.Data,
    Nnid: nnidAll.UserID,
    Hash: miiHash,
  }, nil
}
