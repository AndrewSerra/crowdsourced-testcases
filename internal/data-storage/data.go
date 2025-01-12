/*
 * Created on Mon Dec 16 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package datastorage

import (
	"strings"
	"time"
)

type profileData struct {
	Username string
	Courses  []Course `json:"courses"`
}

func NewProfileData(profileName string, courses ...Course) *profileData {
	return &profileData{
		Username: profileName,
		Courses:  courses,
	}
}

func (p *profileData) AddCourse(c Course) {
	p.Courses = append(p.Courses, c)
}

func (p *profileData) RemoveCourse(c Course) {
	for i, course := range p.Courses {
		if course.Id == c.Id {
			p.Courses = append(p.Courses[:i], p.Courses[i+1:]...)
		}
	}
}

func (p profileData) GetCourseByName(name string) *Course {
	for _, c := range p.Courses {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

func (p profileData) GetCourseById(id int) *Course {
	for _, c := range p.Courses {
		if c.Id == id {
			return &c
		}
	}
	return nil
}

type Course struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Assignments []Assignment `json:"assignments"`
}

func NewCourse(name string, assignments ...Assignment) *Course {
	return &Course{
		Name:        name,
		Assignments: assignments,
	}
}

func (c *Course) AddAssigment(a Assignment) {
	c.Assignments = append(c.Assignments, a)
}

func (c *Course) RemoveAssigment(a Assignment) {
	for i, assignment := range c.Assignments {
		if assignment.Id == a.Id {
			c.Assignments = append(c.Assignments[:i], c.Assignments[i+1:]...)
		}
	}
}

func (c Course) PushlishAssignment(id int) {
	for _, a := range c.Assignments {
		if a.Id == id {
			a.Publish()
		}
	}
}

func (c Course) GetAssignmentByName(name string) *Assignment {
	for _, a := range c.Assignments {
		if a.Name == name {
			return &a
		}
	}
	return nil
}

func (c Course) GetAssignmentById(id int) *Assignment {
	for _, a := range c.Assignments {
		if a.Id == id {
			return &a
		}
	}
	return nil
}

func (c Course) String() string {
	return c.Name
}

type customTime struct {
	time.Time
}

func (ct *customTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(time.DateTime, s)
	return
}

type Assignment struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Start     customTime `json:"start_date"`
	End       customTime `json:"end_date"`
	Published bool       `json:"is_published"`
	IsOpen    bool       `json:"is_open"`
}

func NewAssignment(name string, start customTime, end customTime, published bool, open bool) *Assignment {
	return &Assignment{
		Name:      name,
		Start:     start,
		End:       end,
		Published: published,
		IsOpen:    open,
	}
}

func (a *Assignment) Publish() {
	a.Published = true
}

func (a *Assignment) Unpublish() {
	a.Published = false
}

func (a *Assignment) Open() {
	a.IsOpen = true
}

func (a *Assignment) Close() {
	a.IsOpen = false
}

func (a Assignment) String() string {
	return a.Name
}
