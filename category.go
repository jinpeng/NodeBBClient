package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type CategoryPayload struct {
	Cid int64 `json:"cid"`
	Name string `json:"name"`
	Description string `json:"description"`
	DescriptionParsed string `json:"descriptionParsed"`
	Icon string `json:"icon"`
	BgColor string `json:"bgColor"`
	Color string `json:"color"`
	Slug string `json:"slug"`
	ParentCid int64 `json:"parentCid"`
	TopicCount int64 `json:"topic_count"`
	PostCount int64 `json:"post_count"`
	Disabled int8 `json:"disabled"`
	Order int64 `json:"order"`
	Link string `json:"link"`
	NumRecentReplies int64 `json:"numRecentReplies"`
	Class string `json:"class"`
	ImageClass string `json:"imageClass"`
	IsSection int8 `json:"isSection"`
}

type CategoryResponse struct {
	Code string `json:"code"`
	Payload CategoryPayload `json:"payload"`
}

func createCategory(client *http.Client, name string, parentCid int64) (int64, string) {
	endpoint := SERVER + "/categories"
	v := url.Values{}
	v.Set("name", name)
	v.Set("_uid", "1")
	if parentCid > 0 {
		v.Set("parentCid", string(parentCid))
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

	var target CategoryResponse;
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[Categroy] created: cid: %d, name: %s", target.Payload.Cid, target.Payload.Name)
	return target.Payload.Cid, target.Payload.Name
}

func deleteCategory(client *http.Client, categoryId int64)  {
	endpoint := SERVER + "/categories/" + strconv.FormatInt(categoryId, 10) + "?_uid=1"
	//log.Printf("About to delete category: %d, endpoint: %s", categoryId, endpoint)
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
		log.Printf("[Categroy] deleted: cid: %d", categoryId)
	}
}

func grantCategoryPrivileges(client *http.Client, categoryId int64, groupSlug string) {
	endpoint := SERVER + "/categories/" + strconv.FormatInt(categoryId, 10) + "/privileges"
	v := url.Values{}
	v.Add("privileges", "find")
	v.Add("privileges", "read")
	v.Add("privileges", "topics:read")
	v.Add("privileges", "topics:create")
	v.Add("privileges", "topics:reply")
	v.Add("privileges", "topics:tag")
	v.Add("privileges", "posts:edit")
	v.Add("privileges", "posts:history")
	v.Add("privileges", "posts:delete")
	v.Add("privileges", "posts:upvote")
	v.Add("privileges", "posts:downvote")
	v.Add("privileges", "topics:delete")
	v.Add("groups", groupSlug)
	v.Set("_uid", "1")
	req, err := http.NewRequest("PUT", endpoint, strings.NewReader(v.Encode()))
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
		log.Printf("[Categroy] granted group(slug): %s to category(cid): %d", groupSlug, categoryId)
	}
}
