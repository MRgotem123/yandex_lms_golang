package main

type Shape interface {
	Area() float64
}

func CalculateArea(s Shape) float64 {
	return s.Area()
}

type Rectangle struct {
	Width  float64
	Height float64
}
type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}
func (c Circle) Area() float64 {
	return c.Radius * c.Radius * 3.1415926
}
