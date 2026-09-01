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
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price       float64 `json:"price"`
	Instructor  string `json:"instructor"`
}

var CourseList []Course

func init() {
	CourseJSON := `[
		{
			"id": 1,
			"name": "Golang",
			"price": 100,
			"instructor": "John"
		},
		{
			"id": 2,
			"name": "Java",
			"price": 2000,
			"instructor": "Doe"
		}
	]`
	err := json.Unmarshal([]byte(CourseJSON), &CourseList)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	hightestID := -1
	for _, course := range CourseList {
		if hightestID < course.Id {
			hightestID = course.Id
		}
	}
	return hightestID + 1
}

func findID(ID int) (*Course, int) {
	for i, course := range CourseList {
		if course.Id == ID {
			return &course, i
		}
	}
	return nil, 0
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "course/")
	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	course, listItemIndex := findID(ID)
	if course == nil {
		http.Error(w, fmt.Sprintf("no course with id %d",ID), http.StatusNotFound)
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
			w.WriteHeader(http.StatusOK)
			w.Write(courseJSON)
		case http.MethodPut:
			var updateCourse Course
			bodyByte, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = json.Unmarshal(bodyByte, &updateCourse)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if updateCourse.Id != ID {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			course = &updateCourse
			CourseList[listItemIndex] = *course
			w.WriteHeader(http.StatusOK)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
	}
}

func coursesHandler(w http.ResponseWriter, r *http.Request) {
	courseJSON, err := json.Marshal(CourseList)

	switch r.Method {
	case http.MethodGet:
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)
	case http.MethodPost:
		var newCourse Course
		BodyByte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(BodyByte, &newCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newCourse.Id = getNextID()
		CourseList = append(CourseList, newCourse)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before handler middle start")
		handler.ServeHTTP(w, r)	
		fmt.Println("after handler middle end")
	})
}

func main() {
	courseItem := http.HandlerFunc(courseHandler);
	courseList := http.HandlerFunc(coursesHandler);
	http.Handle("/course/", middlewareHandler(courseItem))
	http.Handle("/course", middlewareHandler(courseList))
	http.ListenAndServe(":5000", nil)
}