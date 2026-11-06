package guessgame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistanceKm(t *testing.T) {
	tests := []struct {
		name                   string
		lat1, lon1, lat2, lon2 float64
		want                   float64
		delta                  float64
	}{
		{"same point", 48.8566, 2.3522, 48.8566, 2.3522, 0, 1e-9},
		{"Paris to London", 48.8566, 2.3522, 51.5074, -0.1278, 344, 2},
		{"across the antimeridian", 0, 179.9, 0, -179.9, 22.2, 0.1},
		{"pole to pole", 90, 0, -90, 0, 20015, 1},
		{"Tokyo to Sydney", 35.6762, 139.6503, -33.8688, 151.2093, 7823, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, DistanceKm(tt.lat1, tt.lon1, tt.lat2, tt.lon2), tt.delta)
			assert.InDelta(t, tt.want, DistanceKm(tt.lat2, tt.lon2, tt.lat1, tt.lon1), tt.delta, "symmetric")
		})
	}
}

func TestMonthsOff(t *testing.T) {
	assert.Equal(t, 0, MonthsOff(2020, 5, 2020, 5))
	assert.Equal(t, 1, MonthsOff(2019, 12, 2020, 1))
	assert.Equal(t, 1, MonthsOff(2020, 1, 2019, 12))
	assert.Equal(t, 25, MonthsOff(2018, 3, 2020, 4))
}

func TestLocationScore(t *testing.T) {
	region := presetScales[DifficultyRegion]

	assert.Equal(t, MaxPartScore, LocationScore(0, region))
	assert.Equal(t, MaxPartScore, LocationScore(0.02, region), "within GPS noise")
	assert.Equal(t, 4678, LocationScore(10, region))
	assert.Equal(t, 1839, LocationScore(150, region))
	assert.Equal(t, 0, LocationScore(20000, region))

	previous := MaxPartScore
	for d := 0.0; d < 3000; d += 7.3 {
		score := LocationScore(d, region)
		assert.LessOrEqual(t, score, previous, "falls monotonically at %v km", d)
		assert.GreaterOrEqual(t, score, 0)
		previous = score
	}
}

func TestTimeScore(t *testing.T) {
	scales := Scales{TimeMonths: 12}

	assert.Equal(t, MaxPartScore, TimeScore(0, scales))
	assert.Equal(t, 3894, TimeScore(3, scales))
	assert.Equal(t, 1839, TimeScore(12, scales))
	assert.Equal(t, 34, TimeScore(60, scales))

	assert.Equal(t, MaxPartScore, TimeScore(0, Scales{}), "zero scale only rewards exact answers")
	assert.Equal(t, 0, TimeScore(1, Scales{}))
}

func TestScalesFor(t *testing.T) {
	assert.Equal(t, Scales{DistanceKm: 15, TimeMonths: 12}, ScalesFor(DifficultyCity, nil))
	assert.Equal(t, Scales{DistanceKm: 1500, TimeMonths: 12}, ScalesFor(DifficultyWorld, nil))
	assert.Equal(t, AutoScales(nil), ScalesFor(DifficultyAuto, nil))
	assert.Equal(t, AutoScales(nil), ScalesFor("unknown", nil))
}

func TestAutoScales(t *testing.T) {
	t.Run("empty pool", func(t *testing.T) {
		assert.Equal(t, Scales{DistanceKm: 1500, TimeMonths: 12}, AutoScales(nil))
	})

	t.Run("one city over two years", func(t *testing.T) {
		pool := []Point{
			{Latitude: 52.52, Longitude: 13.40, Year: 2020, Month: 1},
			{Latitude: 52.53, Longitude: 13.41, Year: 2020, Month: 6},
			{Latitude: 52.51, Longitude: 13.39, Year: 2021, Month: 1},
			{Latitude: 52.50, Longitude: 13.45, Year: 2022, Month: 1},
		}
		scales := AutoScales(pool)
		assert.Equal(t, 5.0, scales.DistanceKm, "clamped to the 5 km floor")
		assert.InDelta(t, 3.0, scales.TimeMonths, 1e-9, "24 months / 10, clamped to 3")
	})

	t.Run("world over twenty years", func(t *testing.T) {
		pool := []Point{
			{Latitude: 48.85, Longitude: 2.35, Year: 2000, Month: 1},
			{Latitude: 40.71, Longitude: -74.00, Year: 2005, Month: 1},
			{Latitude: 35.68, Longitude: 139.65, Year: 2010, Month: 1},
			{Latitude: -33.87, Longitude: 151.21, Year: 2015, Month: 1},
			{Latitude: 52.52, Longitude: 13.40, Year: 2020, Month: 1},
		}
		scales := AutoScales(pool)
		assert.Equal(t, 1500.0, scales.DistanceKm, "clamped to the 1,500 km ceiling")
		assert.Equal(t, 24.0, scales.TimeMonths)
	})
}
