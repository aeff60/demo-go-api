package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Course struct {
	ID         int    `json: "id"`
	Name       string `json: "name"`
	Price      int    `json: "price"`
	Instructor string `json: "instructor"`
}

var courseList []Course

func init() {
	courseJSON := `[ 
		{
			"id":101,
			"name":"Python",
			"price":2590,
			"instructor":"BorntoDev"
		},
		{
			"id":102,
			"name":"JavaScript",
			"price":0,
			"instructor":"BorntoDev"
		},
		{
			"id":103,
			"name":"SQL",
			"price":0,
			"instructor":"BorntoDev"
		}
	]`
	err := json.Unmarshal([]byte(courseJSON), &courseList)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	highestID := -1
	for _, course := range courseList {
		if highestID < course.ID {
			highestID = course.ID
		}
	}
	return highestID + 1
}

func findID(ID int) (*Course, int) {
	for i, course := range courseList {
		if course.ID == ID {
			return &course, i
		}

	}
	return nil, 0

}
func courseHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "course/")
	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	course, listItemIndex := findID(ID)
	if course == nil {
		http.Error(w, fmt.Sprintf("no course with id %d", ID), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		courseJSON, err := json.Marshal(course)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)

	case http.MethodPut:
		var updatedCourse Course
		byteBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(byteBody, &updatedCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedCourse.ID != ID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		course = &updatedCourse
		courseList[listItemIndex] = *course
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

}

func coursesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		courseJSON, err := json.Marshal(courseList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)
	case http.MethodPost:
		var newCourse Course
		bodybytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodybytes, &newCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newCourse.ID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newCourse.ID = getNextID()
		courseList = append(courseList, newCourse)
		w.WriteHeader(http.StatusCreated)
		return
	}

}

// func middlewareHandler(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("before handler middle start")
// 		handler.ServeHTTP(w, r)
// 		fmt.Println("middleware finised\n")

// 	})

// }

func enableCorsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, X-CSRF-Token, Accept-Encoding")
		handler.ServeHTTP(w, r)
	})
}
func main() {

	courseItemHandler := http.HandlerFunc(courseHandler)
	courseListHandler := http.HandlerFunc(coursesHandler)
	http.Handle("/course/", enableCorsMiddleware(courseItemHandler))
	http.Handle("/course", enableCorsMiddleware(courseListHandler))
	http.ListenAndServe(":5000", nil)

}

//noted map
