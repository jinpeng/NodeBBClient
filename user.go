package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type UserPayload struct {
	Uid int64 `json:"uid"`
}

type UserResponse struct {
	Code string `json:"code"`
	Payload UserPayload `json:"payload"`
}

func createUser(client *http.Client, username string, password string, email string) (int64, string) {
	endpoint := SERVER + "/users"
	v := url.Values{}
	v.Set("username", username)
	v.Set("password", password)
	v.Set("email", email)
	v.Set("_uid", "1")
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

	var target UserResponse;
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[User] created: uid: %d, username: %s", target.Payload.Uid, username)
	return target.Payload.Uid, username
}

func deleteUser(client *http.Client, userId int64)  {
	endpoint := SERVER + "/users/" + strconv.FormatInt(userId, 10) + "?_uid=1"
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
		log.Printf("[User] deleted: uid: %d", userId)
	}
}