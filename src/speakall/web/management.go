package web

import (
	"log"
	"net/http"
	"speakall/db"
)

func databaseHandler(w http.ResponseWriter, r *http.Request) {

	user := getLoginUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tc := make(map[string]interface{})
	sql := ""
	if r.Method == "POST" {
		sql = r.FormValue("SQL")
		rows, err := db.Query(sql)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, elm := range rows.Records {
			for _, val := range elm {
				log.Println(val)
				log.Println(*val)
				log.Println(&val)
			}
		}

		tc["Columns"] = rows.Columns
		tc["Records"] = rows.Records
	}

	tc["User"] = user
	tc["SQL"] = sql

	setTemplates(w, tc, "menu.tmpl", "database.tmpl")
	return
}
