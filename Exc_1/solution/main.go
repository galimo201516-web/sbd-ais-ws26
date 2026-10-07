package main

import "fmt"

// Student is a struct: a custom type that groups related fields together.
type Student struct {
	Name     string
	Program  string
	Semester int
}

// Print is a METHOD with a RECEIVER of type Student.
// "(s Student)" ties this function to the Student type,
// so any Student value can call it as s.Print().
func (s Student) Print() {
	fmt.Println("Hello, World!")
	fmt.Println("Name:    ", s.Name)
	fmt.Println("Program: ", s.Program)
	fmt.Println("Semester:", s.Semester)
	fmt.Printf("Whole struct: %+v\n", s)
}

func main() {
	s := Student{
		Name:     "Ghali",
		Program:  "Artificial Intelligence Systems",
		Semester: 3,
	}
	s.Print()
}