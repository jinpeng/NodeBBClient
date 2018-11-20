package main

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	SERVER = "http://localhost:4567/api/v2"
	MASTERTOKEN = "22dc106b-86ce-49f3-ae83-7ab5c2bbee7c"
)

func createSchoolStructure(client *http.Client, schoolIndex int) {
	schoolIndexStr := strconv.Itoa(schoolIndex)
	schoolEmailSuffix := "@u" + schoolIndexStr + ".edu"
	// Create school admin
	adminUsername := "adm-" + schoolIndexStr + "-1";
	adminEmail := adminUsername + schoolEmailSuffix
	admUid, _ := createUser(client, adminUsername, "123456", adminEmail)
	// Create school group, owner is school admin
	schoolGroup := "group-" + schoolIndexStr
	schoolGroupSlug, _ := createGroup(client, schoolGroup, admUid)
	// Create school category
	schoolCategory := "school-" + schoolIndexStr
	schoolCategoryId, _ := createCategory(client, schoolCategory, 0)
	// Grant school group to school category
	grantCategoryPrivileges(client, schoolCategoryId, schoolGroupSlug)

	// Create 10 teachers, 100 class groups, owner is teacher
	// Create 1000 students, add to 100 class groups
	for idx:=1; idx<=10; idx++ {
		teacherUsername := "tch-" + schoolIndexStr + "-" + strconv.Itoa(idx)
		teacherEmail := teacherUsername + schoolEmailSuffix
		tchId, _ := createUser(client, teacherUsername, "123456", teacherEmail)
		for index:=1; index<=10; index++ {
			classGroup := "group-" + schoolIndexStr + "-" + strconv.Itoa((idx-1)*10+index)
			groupSlug, _ := createGroup(client, classGroup, tchId)
			// Create class category
			classCategory := "class-" + schoolIndexStr + "-" + strconv.Itoa((idx-1)*10+index)
			classCategoryId, _ := createCategory(client, classCategory, schoolCategoryId)
			// Grant class group to class category
			grantCategoryPrivileges(client, classCategoryId, groupSlug)

			studentUsernamePrefix := "stu-" + schoolIndexStr + "-" + strconv.Itoa((idx-1)*10+index) + "-"
			for idxStu:=1; idxStu<=10; idxStu++ {
				studentUsername := studentUsernamePrefix + strconv.Itoa(idxStu)
				studentEmail := studentUsername + schoolEmailSuffix
				stuId, _ := createUser(client, studentUsername, "123456", studentEmail)
				addUserToGroup(client, groupSlug, stuId)
			}
		}
	}
}

func main() {
	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 30,
		Transport: netTransport,
	}

	//for idx:=1; idx<200; idx++ {
	//	deleteCategory(netClient, int64(idx))
	//}
	//
	//for idx:=1; idx<200; idx++ {
	//	deleteGroup(netClient, "group-" + strconv.Itoa(idx))
	//}
	//
	//for idx:=2; idx<200; idx++ {
	//	deleteUser(netClient, int64(idx))
	//}

	createSchoolStructure(netClient, 1)
}
