package example

//go:generate ../enum-gen

// enum: RED=#FF0000|GREEN=#00FF00|BLUE=#0000FF
type Color string

// enum: TESLA|VOLVO |MERCEDES
type Car int
