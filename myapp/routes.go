package main

import (
	"fmt"
	"myapp/data"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/techarm/celeritas/mailer"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.get("/", a.Handlers.Home)
	a.get("/go-page", a.Handlers.GoPage)
	a.get("/jet-page", a.Handlers.JetPage)
	a.get("/sessions", a.Handlers.SessionPage)

	a.get("/users/login", a.Handlers.UserLogin)
	a.post("/users/login", a.Handlers.PostUserLogin)
	a.get("/users/logout", a.Handlers.Logout)
	a.get("/form", a.Handlers.FormvalHandlers)
	a.post("/form", a.Handlers.SubmitForm)

	a.get("/json", a.Handlers.JSON)
	a.get("/xml", a.Handlers.XML)
	a.get("/download-file", a.Handlers.DownloadFile)
	a.get("/crypto", a.Handlers.TestCrypto)

	a.get("/cache", a.Handlers.ShowCachePage)
	a.post("/api/save-in-cache", a.Handlers.SaveInCache)
	a.post("/api/get-from-cache", a.Handlers.GetFromCache)
	a.post("/api/delete-from-cache", a.Handlers.DeleteFromCache)
	a.post("/api/empty-cache", a.Handlers.EmptyCache)

	a.get("/test-mail", func(w http.ResponseWriter, r *http.Request) {
		msg := mailer.Message{
			From:        "test@example.com",
			To:          "you@there.com",
			Subject:     "Test Subject - send using channel",
			Template:    "test",
			Attachments: nil,
			Data:        nil,
		}

		a.App.Mail.Jobs <- msg
		res := <-a.App.Mail.Result
		if res.Error != nil {
			a.App.ErrorLog.Println(res.Error)
		}

		err := a.App.Mail.SendSMTPMessage(msg)
		if err != nil {
			a.App.ErrorLog.Println(err)
		}

		fmt.Fprintf(w, "Send mail success.")
	})

	a.get("/create-user", func(w http.ResponseWriter, r *http.Request) {
		firstName := a.App.RandomString(5)
		lastName := a.App.RandomString(5)
		u := data.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     fmt.Sprintf("%s_%s@gmail.com", firstName, lastName),
			Active:    1,
			Password:  "password",
		}

		id, err := a.Models.Users.Insert(u)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "%d: %s", id, u.FirstName)
	})

	a.get("/get-all-users", func(w http.ResponseWriter, r *http.Request) {
		users, err := a.Models.Users.GetAll()
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		var sb strings.Builder
		for _, x := range users {
			sb.WriteString(fmt.Sprintf("%d: %s, %s, %s\n", x.ID, x.FirstName, x.LastName, x.Email))
		}
		fmt.Fprint(w, sb.String())
	})

	a.get("/get-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "%d: %s, %s", id, u.FirstName, u.Email)
	})

	a.get("/update-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		u.LastName = a.App.RandomString(10)
		u.LastName = ""
		validator := a.App.Validator(nil)
		u.Validate(validator)

		if !validator.Valid() {
			fmt.Fprint(w, "validation error")
			return
		}

		err = u.Update(*u)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "update user lastname: %s", u.LastName)
	})

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
