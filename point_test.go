package geo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"testing"
)

// Tests that a call to NewPoint should return a pointer to a Point with the specified values assigned correctly.
func TestNewPoint(t *testing.T) {
	p := NewPoint(40.5, 120.5)

	if p.lat != 40.5 {
		t.Errorf("Expected to be able to specify 40.5 as the lat value of a new point, but got %f instead", p.lat)
	}

	if p.lng != 120.5 {
		t.Errorf("Expected to be able to specify 120.5 as the lng value of a new point, but got %f instead", p.lng)
	}
}

// Tests that calling GetLat() after creating a new point returns the expected lat value.
func TestLat(t *testing.T) {
	p := NewPoint(40.5, 120.5)

	lat := p.Lat()

	if lat != 40.5 {
		t.Errorf("Expected a call to GetLat() to return the same lat value as was set before, but got %f instead", lat)
	}
}

// Tests that calling GetLng() after creating a new point returns the expected lng value.
func TestLng(t *testing.T) {
	p := NewPoint(40.5, 120.5)

	lng := p.Lng()

	if lng != 120.5 {
		t.Errorf("Expected a call to GetLng() to return the same lat value as was set before, but got %f instead", lng)
	}
}

// Ensures that a point can be marhalled into JSON
func TestMarshalJSON(t *testing.T) {
	p := NewPoint(40.7486, -73.9864)
	res, err := json.Marshal(p)

	if err != nil {
		log.Print(err)
		t.Error("Should not encounter an error when attempting to Marshal a Point to JSON")
	}

	if string(res) != `{"lat":40.7486,"lng":-73.9864}` {
		t.Error("Point should correctly Marshal to JSON")
	}
}

// Ensures that a point can be unmarhalled from JSON
func TestUnmarshalJSON(t *testing.T) {
	data := []byte(`{"lat":40.7486,"lng":-73.9864}`)
	p := &Point{}
	err := p.UnmarshalJSON(data)

	if err != nil {
		t.Errorf("Should not encounter an error when attempting to Unmarshal a Point from JSON")
	}

	if p.lat != 40.7486 || p.lng != -73.9864 {
		t.Errorf("Point has mismatched data after Unmarshalling from JSON")
	}
}

// Ensure that a point can be marshalled into slice of binaries
func TestMarshalBinary(t *testing.T) {
	lat, long := 40.7486, -73.9864
	p := NewPoint(lat, long)
	actual, err := p.MarshalBinary()
	if err != nil {
		t.Error("Should not encounter an error when attempting to Marshal a Point to binary", err)
	}

	expected, err := coordinatesToBytes(lat, long)
	if err != nil {
		t.Error("Unable to convert coordinates to bytes slice.", err)
	}

	if !bytes.Equal(actual, expected) {
		t.Errorf("Point should correctly Marshal to Binary.\nExpected %v\nBut got %v", expected, actual)
	}
}

// Ensure that a point can be unmarshalled from a slice of binaries
func TestUnmarshalBinary(t *testing.T) {
	lat, long := 40.7486, -73.9864
	coordinates, err := coordinatesToBytes(lat, long)
	if err != nil {
		t.Error("Unable to convert coordinates to bytes slice.", err)
	}

	actual := Point{}
	err = actual.UnmarshalBinary(coordinates)
	if err != nil {
		t.Error("Should not encounter an error when attempting to Unmarshal a Point from binary", err)
	}

	expected := NewPoint(lat, long)
	if !assertPointsEqual(actual, expected, 4) {
		t.Errorf("Point should correctly Marshal to Binary.\nExpected %+v\nBut got %+v", expected, actual)
	}
}

func coordinatesToBytes(lat, long float64) ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, lat); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, long); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Asserts true when the latitude and longtitude of p1 and p2 are equal up to a certain number of decimal places.
// Precision is used to define that number of decimal places.
func assertPointsEqual(p1, p2 Point, precision int) bool {
	roundedLat1, roundedLng1 := int(p1.lat*float64(precision))/precision, int(p1.lng*float64(precision))/precision
	roundedLat2, roundedLng2 := int(p2.lat*float64(precision))/precision, int(p2.lng*float64(precision))/precision
	return roundedLat1 == roundedLat2 && roundedLng1 == roundedLng2
}
