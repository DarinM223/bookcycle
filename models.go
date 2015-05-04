package main

type User struct {
	Id        int `sql:"AUTO_INCREMENT"`
	Firstname string
	Lastname  string
	Username  string  `sql:"not null"`
	Rating    float64 `sql:"not null" sql:"default:0"`
	Email     string  `sql:"not null"`
	Phone     int
}

type Book struct {
	Id        int    `sql:"AUTO_INCREMENT"`
	Title     string `sql:"not null"`
	Author    string `sql:"not null"`
	Class     string `sql:"not null"`
	Professor string `sql:"not null"`
	Version   string `sql:"not null"`
	Price     string `sql:"not null"`
	Condition string `sql:"not null"`
	UserId    int    `sql:"index"`
}
