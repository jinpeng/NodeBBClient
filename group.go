package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type GroupPayload struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Createtime int64 `json:"createtime"`
	UserTitle string `json:"userTitle"`
	UserTitleEnabled int8 `json:"userTitleEnabled"`
	Description string `json:"description"`
	MemberCount int64 `json:"memberCount"`
	Hidden int8 `json:"hidden"`
	System int8 `json:"system"`
	Private int8 `json:"private"`
	DisableJoinRequests int8 `json:"disableJoinRequests"`
	OwnerUid int64 `json:"ownerUid"`
}

type GroupResponse struct {
	Code string `json:"code"`
	Payload GroupPayload `json:"payload"`
}

func createGroup(client *http.Client, name string, ownerUid int64) (string, string) {
	endpoint := SERVER + "/groups"
	v := url.Values{}
	v.Set("name", name)
	v.Set("_uid", "1")
	if ownerUid > 0 {
		v.Set("ownerUid", string(ownerUid))
	}
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(v.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+MASTERTOKEN)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var target GroupResponse;
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		log.Fatal(err)
	}
	if target.Code != "ok" {
		log.Fatal("Failed to create group: ", target.Code)
	}
	log.Printf("[Group] created: slug: %s, name: %s", target.Payload.Slug, target.Payload.Name)
	return target.Payload.Slug, target.Payload.Name
}

func deleteGroup(client *http.Client, groupSlug string)  {
	endpoint := SERVER + "/groups/" + groupSlug + "?_uid=1"
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+MASTERTOKEN)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Printf("[Group] deleted: slug: %s", groupSlug)
	}
}

func addUserToGroup(client *http.Client, groupSlug string, userId int64) {
	endpoint := SERVER + "/groups/" + groupSlug + "/membership/" + strconv.FormatInt(userId, 10) + "?_uid=1"
	req, err := http.NewRequest("PUT", endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+MASTERTOKEN)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Printf("[Group] added user(uid): %d to group(slug): %s", userId, groupSlug)
	}
}
