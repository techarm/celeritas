package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"net/http"
	"strings"
	"time"

	up "github.com/upper/db/v4"
)

type Token struct {
	ID        int       `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Email     string    `db:"email" json:"email"`
	PlainText string    `db:"token" json:"token"`
	Hash      []byte    `db:"token_hash" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Expires   time.Time `db:"expiry" json:"expiry"`
}

func (t *Token) Table() string {
	return "tokens"
}

func (t *Token) GetUserForToken(tokenString string) (*User, error) {
	var token Token
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": tokenString})
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	var user User
	collection = upper.Collection(user.Table())
	res = collection.Find(up.Cond{"id": token.UserID})
	err = res.One(&user)
	if err != nil {
		return nil, err
	}

	user.Token = token
	return &user, nil
}

func (t *Token) GetTokensForUser(id int) ([]*Token, error) {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"user_id": id})

	var tokens []*Token
	err := res.All(&tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (t *Token) Get(id int) (*Token, error) {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"id": id})

	var token Token
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *Token) GetByToken(plainText string) (*Token, error) {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": plainText})

	var token Token
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *Token) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)

	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) DeleteByToken(plainText string) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": plainText})

	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) Insert(token Token, u User) error {
	collection := upper.Collection(t.Table())

	// delete existing tokne
	res := collection.Find(up.Cond{"user_id": u.ID})
	err := res.Delete()
	if err != nil {
		return err
	}

	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	token.FirstName = u.FirstName
	token.Email = u.Email

	_, err = collection.Insert(token)
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) GenerateToken(userID int, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID:  t.UserID,
		Expires: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

func (t *Token) AuthenticateToken(r *http.Request) (*User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorizaiton header received")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorizaiton header received")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("token wrong size")
	}

	t, err := t.GetByToken(token)
	if err != nil {
		return nil, errors.New("no matching token found")
	}

	if t.Expires.Before(time.Now()) {
		return nil, errors.New("expired token")
	}

	user, err := t.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}

func (t *Token) ValidToken(token string) (bool, error) {
	user, err := t.GetUserForToken(token)
	if err != nil {
		return false, errors.New("no matching user found")
	}

	if user.Token.PlainText == "" {
		return false, errors.New("no matching token found")
	}

	if user.Token.Expires.Before(time.Now()) {
		return false, errors.New("expired token")
	}

	return true, nil
}
