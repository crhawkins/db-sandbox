package models

type Stock struct {
	ID    int
	Count int
	Price float64

	Car Car
}
