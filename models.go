package main

type Product struct {
	ID       int    `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Quantity int
	Price    float64 `gorm:"type:numeric(10,7)"`
}
