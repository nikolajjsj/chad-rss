package server

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWT struct {
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}
