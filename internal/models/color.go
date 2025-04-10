package models

type Color struct {
	ID   int
	Name string `meta:"required=true;max_len=16"`
}
