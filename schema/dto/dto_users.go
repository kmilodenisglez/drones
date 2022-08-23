package dto

// User struct
type User struct {
	Username   string `json:"username"`
	Passphrase string `json:"passphrase"`
	Name       string `json:"name"`
}

var Users = []User{{
	Passphrase: "0b14d501a594442a01c6859541bcb3e8164d183d32937b851835442f69d5c94e", // password1
	Username:   "richard.sargon@meinermail.com",
	Name:       "Richard Sargon",
}, {
	Passphrase: "6cf615d5bcaac778352a8f1f3360d23f02f34ec182e259897fd6ce485d7870d4", // password1
	Username:   "tom.carter@meinermail.com",
	Name:       "Tom Carter",
}}