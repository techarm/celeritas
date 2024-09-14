package handlers

import (
	"fmt"
	"myapp/data"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

// FormvalHandlers comment goes here
func (h *Handlers) FormvalHandlers(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)
	validator := h.App.Validator(nil)
	vars.Set("validator", validator)
	vars.Set("user", data.User{})

	err := h.App.Render.Page(w, r, "form", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) SubmitForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	validator := h.App.Validator(nil)
	validator.Required(r, "first_name", "last_name", "email")
	validator.Check(len(r.Form.Get("first_name")) > 1, "first_name", "Must be at least two characters")
	validator.Check(len(r.Form.Get("last_name")) > 1, "last_name", "Must be at least two characters")

	var user data.User
	user.FirstName = r.Form.Get("first_name")
	user.LastName = r.Form.Get("last_name")
	user.Email = r.Form.Get("email")
	user.Validate(validator)

	if !validator.Valid() {
		vars := make(jet.VarMap)
		vars.Set("validator", validator)

		vars.Set("user", user)
		if err := h.App.Render.Page(w, r, "form", vars, nil); err != nil {
			h.App.ErrorLog.Println(err)
			return
		}

		return
	}

	fmt.Fprintf(w, "valide data")
}
