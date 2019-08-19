package school

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)


//Adding new struct to demonstrate 1 to 1 relationship between Instructor and Course
//Changes made to instructors csv by adding courseId to link to course table.
type Instructor struct {
	Id int
	Name string
}
type Course struct{
	Id int
	Name string
	Instructor_Id int  //Foreign Key (FK) to establish link to Instructor.
}

//To use after looping through slices
type InstructorCourse struct{
	InstructorName string
	CourseName string
}

//Place structs here


//**************************************************************************************************


// will parse csv file, output a slice of Instructor
func prs(filePath string) ([]Instructor) {
	src, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer src.Close()

	rdr := csv.NewReader(src)
	rows, err := rdr.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	//x denotes for me the slice. Create this outside of the loop to hold all the instructor types.
	xInstructor := []Instructor{}

	for i, row := range rows {
		//Throws away the first row
		if i == 0 {
			continue //then continues
		}


		id, _ := strconv.Atoi(row[0])
		name := row[1]


		//Two different ways: I am creating the Instructor then adding to slice
		/*ins := Instructor{
			Id: id,
			Name: name,
		}
		xInstructor = append(xInstructor, ins)	*/


		//or add type Instructor to the xInstructor all at once, a slice
		xInstructor = append(xInstructor, Instructor{
			Id: id,
			Name: name,
		})
	}
	//returns the slice of Instructor
	return xInstructor

}

// will parse csv file, output a slice of Instructor
func prsCourse(filePath string) ([]Course) {
	src, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer src.Close()

	rdr := csv.NewReader(src)
	rows, err := rdr.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	//x denotes for me the slice. Create this outside of the loop to hold all the instructor types.
	xCourse := []Course{}

	for i, row := range rows {  //i is the iteration count starting at 0, use it to throw away first row
		//Throws away the first row
		//row is the value, in this case its the entire row which has 3 elements seperted by commas.
		if i == 0 {
			continue //then continues
		}

		//This "strconv.Atoi" converts Ascii To Int thus AtoI, elements in csv are Ascii
		id, _ := strconv.Atoi(row[0]) //row[0] is the first element, the id
		// throwing away the error by using _ after the comma
		name := row[1] //no conversion needed since it is already a string, Ascii.
		instructorId,_ := strconv.Atoi(row[2])

		//Two different ways: I am creating the Instructor then adding to slice
		/*ins := Instructor{
			Id: id,
			Name: name,
		}
		xInstructor = append(xInstructor, ins)	*/



		xCourse = append(xCourse, Course{
			Id: id,
			Name: name,
			Instructor_Id: instructorId,
			//New field in this type. Each course will be taught by one instructor
		})
	}
	//returns the slice of Instructor
	return xCourse

}

	func MatchInstructorToCourse(xInstructor []Instructor,xCourse []Course)([]InstructorCourse){
	//bringing in both slices and naming them so as to use them next.
	//result will be a slice of type InstructorCourse
	//common to iterate through slices
	//i name the results outer for outer loop, inner for the inner loop, easy to get confuse

	//instatiate this slice before ranging, if in loop it would just make a new one each loop
	//need to add to this with each loop
		xInstructorCourse := []InstructorCourse{}


		for _, outer := range xInstructor{
			//outer will be of type Instructor
			//will take one Instructor named outer
			for _, inner := range xCourse{
			//inner will be of type Course
			//loop through all of these to find match of Instructor Id's

			if outer.Id == inner.Instructor_Id {
			//once match found, stuff names into InstuctorCourse type then stuff into slice for return.
				xInstructorCourse = append(xInstructorCourse,InstructorCourse{
					InstructorName: outer.Name,
					CourseName: inner.Name,//remember the comma after last field in type
				})

			}}
		}

		return xInstructorCourse
	}