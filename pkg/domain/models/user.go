package models

type User struct {
	CreatedAt  int64  `json:"created_at"`
	ModifiedAt int64  `json:"modified_at"`
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Email      string `json:"email"`
	Password   string `json:"Password"`
	Verified   bool   `json:"verified"`
}
