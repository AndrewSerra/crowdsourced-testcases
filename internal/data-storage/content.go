/*
 * Created on Mon Dec 16 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package datastorage

import "time"

type profileData struct {
	Username string
	Courses  []*course
}

func NewProfileData(profileName string, courses ...*course) *profileData {
	return &profileData{
		Username: profileName,
		Courses:  courses,
	}
}

func (p *profileData) AddCourse(c *course) {
	p.Courses = append(p.Courses, c)
}

func (p *profileData) RemoveCourse(c *course) {
	for i, course := range p.Courses {
		if course == c {
			p.Courses = append(p.Courses[:i], p.Courses[i+1:]...)
		}
	}
}

func (p profileData) GetCourseByName(name string) *course {
	for _, c := range p.Courses {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (p profileData) GetCourseById(id int) *course {
	for _, c := range p.Courses {
		if c.Id == id {
			return c
		}
	}
	return nil
}

type course struct {
	Id          int
	Name        string
	Assignments []*assignment
}

func NewCourse(name string, assignments ...*assignment) *course {
	return &course{
		Name:        name,
		Assignments: assignments,
	}
}

func (c *course) AddAssigment(a *assignment) {
	c.Assignments = append(c.Assignments, a)
}

func (c *course) RemoveAssigment(a *assignment) {
	for i, assignment := range c.Assignments {
		if assignment == a {
			c.Assignments = append(c.Assignments[:i], c.Assignments[i+1:]...)
		}
	}
}

func (c course) GetAssignmentByName(name string) *assignment {
	for _, a := range c.Assignments {
		if a.Name == name {
			return a
		}
	}
	return nil
}

func (c course) GetAssignmentById(id int) *assignment {
	for _, a := range c.Assignments {
		if a.Id == id {
			return a
		}
	}
	return nil
}

type assignment struct {
	Id        int
	Name      string
	Start     time.Time
	End       time.Time
	Published bool
	Open      bool
}

func NewAssignment(name string, start time.Time, end time.Time, published bool, open bool) *assignment {
	return &assignment{
		Name:      name,
		Start:     start,
		End:       end,
		Published: published,
		Open:      open,
	}
}
