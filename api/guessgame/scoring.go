// Package guessgame holds the database-free logic of the 4D-guesser game:
// scoring, difficulty scales, photo sampling and round ids.
//
// Fork-specific: not part of upstream Photoview.
package guessgame

import (
	"math"
	"slices"
)

const (
	// MaxPartScore is the most a round can score for place, and separately for time.
	MaxPartScore = 5000
	// MaxRoundScore is the most a single round can score.
	MaxRoundScore = 2 * MaxPartScore

	earthRadiusKm = 6371.0

	// exactDistanceKm absorbs GPS noise: a pin this close counts as exact.
	exactDistanceKm = 0.025
)

// Scales are the decay constants of a game: D₀ for distance and T₀ for time.
type Scales struct {
	DistanceKm float64
	TimeMonths float64
}

// Difficulty names a preset of Scales. DifficultyAuto derives them from the photo pool.
type Difficulty string

const (
	DifficultyAuto   Difficulty = "AUTO"
	DifficultyCity   Difficulty = "CITY"
	DifficultyRegion Difficulty = "REGION"
	DifficultyWorld  Difficulty = "WORLD"
)

// defaultTimeScaleMonths is T₀ for every difficulty. With the bell-shaped time score
// a year off keeps 90% (4,500 points), 3 years off about 1,900 and 5 years off about 360.
const defaultTimeScaleMonths = 37

var presetScales = map[Difficulty]Scales{
	DifficultyCity:   {DistanceKm: 15, TimeMonths: defaultTimeScaleMonths},
	DifficultyRegion: {DistanceKm: 150, TimeMonths: defaultTimeScaleMonths},
	DifficultyWorld:  {DistanceKm: 1500, TimeMonths: defaultTimeScaleMonths},
}

// Point is a photo location with the year and month it was taken.
type Point struct {
	Latitude  float64
	Longitude float64
	Year      int
	Month     int
}

// DistanceKm returns the great-circle distance between two coordinates, using the haversine formula.
func DistanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	dPhi := (lat2 - lat1) * math.Pi / 180
	dLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Pow(math.Sin(dPhi/2), 2) + math.Cos(phi1)*math.Cos(phi2)*math.Pow(math.Sin(dLambda/2), 2)
	return 2 * earthRadiusKm * math.Asin(math.Min(1, math.Sqrt(a)))
}

// MonthsOff returns the gap between two year/month pairs in whole months.
func MonthsOff(year1, month1, year2, month2 int) int {
	diff := 12*(year1-year2) + (month1 - month2)
	if diff < 0 {
		return -diff
	}
	return diff
}

// LocationScore scores a guess that is distanceKm away from the answer.
func LocationScore(distanceKm float64, scales Scales) int {
	if distanceKm <= exactDistanceKm {
		return MaxPartScore
	}
	return decay(distanceKm, scales.DistanceKm)
}

// TimeScore scores a guess that is monthsOff months away from the answer, as
// 5000·e^(−(monthsOff/T₀)²): close guesses lose little, far ones fall off quickly.
func TimeScore(monthsOff int, scales Scales) int {
	if scales.TimeMonths <= 0 {
		return decay(float64(monthsOff), 0)
	}
	ratio := float64(monthsOff) / scales.TimeMonths
	return int(math.Round(MaxPartScore * math.Exp(-ratio*ratio)))
}

func decay(value, scale float64) int {
	if scale <= 0 {
		if value == 0 {
			return MaxPartScore
		}
		return 0
	}
	return int(math.Round(MaxPartScore * math.Exp(-value/scale)))
}

// ScalesFor returns the scales of a difficulty. For DifficultyAuto, and for unknown
// difficulties, they are derived from the photo pool.
func ScalesFor(difficulty Difficulty, pool []Point) Scales {
	if scales, ok := presetScales[difficulty]; ok {
		return scales
	}
	return AutoScales(pool)
}

// AutoScales derives the distance scale from the spread of a photo pool:
// D₀ is a quarter of the 90th-percentile distance from the pool's median point, clamped to 5-1,500 km.
// T₀ is the same for every pool, defaultTimeScaleMonths.
func AutoScales(pool []Point) Scales {
	scales := Scales{DistanceKm: 1500, TimeMonths: defaultTimeScaleMonths}
	if len(pool) == 0 {
		return scales
	}

	lats := make([]float64, len(pool))
	lons := make([]float64, len(pool))
	for i, p := range pool {
		lats[i] = p.Latitude
		lons[i] = p.Longitude
	}
	medianLat, medianLon := median(lats), median(lons)

	distances := make([]float64, len(pool))
	for i, p := range pool {
		distances[i] = DistanceKm(medianLat, medianLon, p.Latitude, p.Longitude)
	}
	slices.Sort(distances)
	p90 := distances[(len(distances)-1)*9/10]

	scales.DistanceKm = clamp(p90/4, 5, 1500)
	return scales
}

func median(values []float64) float64 {
	sorted := slices.Clone(values)
	slices.Sort(sorted)
	return sorted[len(sorted)/2]
}

func clamp(value, low, high float64) float64 {
	return math.Max(low, math.Min(value, high))
}
