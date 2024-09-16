package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"myapp/data"
	"net/http"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/techarm/celeritas/mailer"
	"github.com/techarm/celeritas/urlsigner"
)

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "login", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := h.Models.Users.GetByEmail(email)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	matches, err := user.PasswordMatches(password)
	if err != nil {
		w.Write([]byte("Error validating password"))
		return
	}

	if !matches {
		w.Write([]byte("Error validating password"))
		return
	}

	// dit the user check remember me?
	if r.Form.Get("remember") == "remember" {
		randomString := h.randomString(12)
		hasher := sha256.New()
		_, err := hasher.Write([]byte(randomString))
		if err != nil {
			h.App.ErrorInternalServerError(w, r)
			return
		}

		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		rm := data.RemeberToken{}
		err = rm.InsertToken(user.ID, sha)
		if err != nil {
			h.App.ErrorInternalServerError(w, r)
			return
		}

		// set a cookie
		cookie := http.Cookie{
			Name:     fmt.Sprintf("%s_remember", h.App.AppName),
			Value:    fmt.Sprintf("%d|%s", user.ID, sha),
			Path:     "/",
			Expires:  time.Now().Add(365 * 24 * 60 * 60 * time.Second),
			HttpOnly: true,
			Domain:   h.App.Session.Cookie.Domain,
			MaxAge:   315350000,
			Secure:   h.App.Session.Cookie.Secure,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)
		// save hash in session
		h.App.Session.Put(r.Context(), "remember_token", sha)
	}

	h.sessionPut(r.Context(), "userID", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	// delete the remember token if exists
	if h.App.Session.Exists(r.Context(), "remember_token") {
		rt := data.RemeberToken{}
		_ = rt.Delete(h.App.Session.GetString(r.Context(), "remember_token"))
	}

	// delete cooike
	newCookie := http.Cookie{
		Name:     fmt.Sprintf("%s_remember", h.App.AppName),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-100 * time.Hour),
		HttpOnly: true,
		Domain:   h.App.Session.Cookie.Domain,
		MaxAge:   -1,
		Secure:   h.App.Session.Cookie.Secure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &newCookie)

	h.sessionRemove(r.Context(), "userID")
	h.sessionRemove(r.Context(), "remember_token")
	h.sessionDestroy(r.Context())
	h.sessionRenew(r.Context())

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (h *Handlers) Forgot(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "forgot", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("Error rendering: ", err)
		h.App.ErrorInternalServerError(w, r)
	}
}

func (h *Handlers) PostForgot(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorBadRequest(w, r)
		return
	}

	// varify thjat supplied email exists
	var u *data.User
	email := r.Form.Get("email")
	u, err = u.GetByEmail(email)
	if err != nil {
		h.App.ErrorBadRequest(w, r)
		return
	}

	// create a link to password reset form
	link := fmt.Sprintf("%s/users/reset-password?email=%s", h.App.Server.URL, email)
	sign := urlsigner.Signer{
		Secret: []byte(h.App.EncryptionKey),
	}

	// sign the link
	signedLink := sign.GenerateTokenFromString(link)
	h.App.InfoLog.Println("Signed link is ", signedLink)

	// email the message
	var data struct {
		Link string
	}
	data.Link = signedLink

	msg := mailer.Message{
		To:       u.Email,
		Subject:  "Password reset",
		Template: "password-reset",
		Data:     data,
		From:     "admin@example.com",
	}

	h.App.Mail.Jobs <- msg
	res := <-h.App.Mail.Result
	if res.Error != nil {
		h.App.ErrorLog.Println("send reset password error", res.Error)
		h.App.ErrorBadRequest(w, r)
		return
	}

	// redirect the user
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (h *Handlers) ResetPasswordForm(w http.ResponseWriter, r *http.Request) {
	// get from values
	email := r.URL.Query().Get("email")
	theURL := r.RequestURI
	testURL := fmt.Sprintf("%s%s", h.App.Server.URL, theURL)

	// validate the url
	signer := urlsigner.Signer{
		Secret: []byte(h.App.EncryptionKey),
	}
	valid := signer.VerifyToken(testURL)
	if !valid {
		h.App.ErrorLog.Println("invalid url: ", testURL)
		h.App.ErrorUnauthorized(w, r)
		return
	}

	// make sure it's not expred
	expired := signer.Expired(testURL, 60)
	if expired {
		h.App.ErrorLog.Println("link expired: ", testURL)
		h.App.ErrorUnauthorized(w, r)
		return
	}

	// display form
	encryptedMail, _ := h.encrypt(email)
	vars := make(jet.VarMap)
	vars.Set("email", encryptedMail)

	err := h.render(w, r, "reset-password", vars, nil)
	if err != nil {
		h.App.ErrorInternalServerError(w, r)
		return
	}
}

func (h *Handlers) PostResetPassword(w http.ResponseWriter, r *http.Request) {
	// parse the form
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorInternalServerError(w, r)
		return
	}

	// get and decrypt the email
	email, err := h.decrypt(r.Form.Get("email"))
	if err != nil {
		h.App.ErrorInternalServerError(w, r)
		return
	}

	// get the user
	var u data.User
	user, err := u.GetByEmail(email)
	if err != nil {
		h.App.ErrorInternalServerError(w, r)
		return
	}

	// reset the password
	err = user.ResetPassword(user.ID, r.Form.Get("password"))
	if err != nil {
		h.App.ErrorInternalServerError(w, r)
		return
	}

	// redirect
	h.App.Session.Put(r.Context(), "flash", "Password reset. You can now log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
