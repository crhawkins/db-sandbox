package models

type Car struct {
	ID      int
	Company Company
	Model   string `meta:"required=true;max_len=128"`
	Color   Color
}
