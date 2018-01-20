package customer

import (
	//"errors"

	"github.com/knarfmon/GoMetaCheck/101-WebMetaCheck/config"
	"net/http"
	//"strconv"
)

type customer struct {
	Id      int
	Name    string
	Archive bool

}

func AllCustomers()([]customer,error) {

	rows, err := DB.Query("SELECT * FROM customer")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return nil, err
	}

	defer rows.Close()

	css := make([]customer, 0)

	for rows.Next() {

		cs := customer{}
		err := rows.Scan(&cs.Id,&cs.Name,&cs.Archive) // order matters, everything in select statement

		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return nil,err
		}

		css = append(css, cs)


	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)

		return nil,err
	}
	return css, nil
}