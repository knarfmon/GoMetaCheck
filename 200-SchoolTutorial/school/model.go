package school

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)


// A new type called Instructor, a blueprint, for compiler to determine size req.
// Remember the fields are capitalized because they will be moved from package to package.
type Instructor struct {
	Id int
	Name string
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

		//date, _ := time.Parse("2006-01-02", row[0])  strconv.Atoi(r.FormValue("id"))
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