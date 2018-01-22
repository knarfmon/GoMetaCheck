package customer

import (
	"errors"
	"github.com/knarfmon/GoMetaCheck/101-WebMetaCheck/config"
	"net/http"
	"strconv"
)


type Customer struct {
	Id      	int
	Name    	string
	Archive 	bool
}

type Site struct {
	Id         int
	CustomerId int
	Name       string
	Url        string
	Archive    bool
}





func AllCustomers()([]Customer,error) {

	rows, err := config.DB.Query("SELECT id,name,archive FROM customer ORDER BY name")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	css := make([]Customer, 0)

	for rows.Next() {

		cs := Customer{}
		err := rows.Scan(&cs.Id,&cs.Name,&cs.Archive) // order matters, everything in select statement

		if err != nil {
			return nil,err
		}

		css = append(css, cs)


	}
	if err = rows.Err(); err != nil {
		return nil,err
	}
	return css, nil
}


func PutCustomer(r *http.Request) (Customer, error) {
	// get form values
	cs := Customer{}
	cs.Name = r.FormValue("name")


	// validate form values
	if cs.Name == "" {
		return cs, errors.New("400. Bad request. All fields must be complete.")
	}



	// insert values
	_, err := config.DB.Exec("INSERT INTO customer (name) VALUES (?)", cs.Name)
	if err != nil {
		return cs, errors.New("500. Internal Server Error." + err.Error())
	}
	return cs, nil
}

func OneCustomer(r *http.Request) (Customer, error) {
	cs := Customer{}
	id := r.FormValue("id")
	if id == "" {
		return cs, errors.New("400. Bad Request.")
	}

	row := config.DB.QueryRow("SELECT id,name,archive FROM customer WHERE id = ?", id)

	err := row.Scan(&cs.Id, &cs.Name, &cs.Archive)
	if err != nil {
		return cs, err
	}

	return cs, nil
}

func UpdateCustomer(r *http.Request) (Customer, error) {
	// get form values
	cs := Customer{}
	cs.Name = r.FormValue("name")
	strId := r.FormValue("id")

	if cs.Name == ""  {
		return cs, errors.New("400. Bad Request. Fields can't be empty.")
	}

	newId, err := strconv.Atoi(strId)
	if err != nil {
		return cs, errors.New("406. Not Acceptable. Id not of correct type")
	}
	cs.Id = newId

	// insert values
	_, err = config.DB.Exec("UPDATE customer SET name = ? WHERE id=?;", cs.Name, cs.Id)
	if err != nil {
		return cs, err
	}
	return cs, nil
}