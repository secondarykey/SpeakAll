package web

import (
	"net/http"
	"speakall/db"
	"strings"
)

//category
func categoryHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	pathS := strings.Split(url, "/")
	if len(pathS) > 2 {
		http.Error(w, url+" is Not Found.", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		user := getLoginUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		tc := make(map[string]interface{})
		tc["User"] = user

		key, err := db.GenerateCategoryKey()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tc["CategoryKey"] = key

		templateDir := "templates/"
		setTemplates(w, tc, templateDir+"menu.tmpl", templateDir+"category.tmpl")
	} else {
		name := r.FormValue("name")
		desc := r.FormValue("description")
		key := r.FormValue("key")
		_, err := db.InsertCategory(key, name, desc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//DashboardにURLを教える

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func categoryListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "GETしないで><", http.StatusBadRequest)
		return
	}
	cats, err := db.SelectAllCategory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	setJson(cats, w)
}

func categoryViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "GETしないで><", http.StatusBadRequest)
		return
	}
	url := r.URL.Path
	pathS := strings.Split(url, "/")

	cat, err := db.SelectCategory(pathS[3])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	setJson(cat, w)
}