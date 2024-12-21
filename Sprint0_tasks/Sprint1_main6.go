package main

import "reflect"

type Vehicle interface {
	CalculateTravelTime(distance float64) float64
	GetType() string
}

type Car struct {
	Speed    float64
	Type     string
	FuelType string
}

func (c Car) CalculateTravelTime(distance float64) float64 {
	return distance / c.Speed
}

func (c Car) GetType() string {
	return c.Type
}

type Motorcycle struct {
	Speed float64
	Type  string
}

func (m Motorcycle) CalculateTravelTime(distance float64) float64 {
	return distance / m.Speed
}

func (m Motorcycle) GetType() string {
	return m.Type
}

func EstimateTravelTime(vehicles []Vehicle, distance float64) map[string]float64 {
	travelTimes := make(map[string]float64)
	for _, vehicle := range vehicles {
		vehicleType := reflect.TypeOf(vehicle).String()
		travelTimes[vehicleType] = vehicle.CalculateTravelTime(distance)
	}
	return travelTimes
}
