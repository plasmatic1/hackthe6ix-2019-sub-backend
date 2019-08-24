package main

import (
	"errors"
	"fmt"
	"math"
)

type Location struct {
	latitude float64
	longitude float64
}

type SetTeam struct {
	team int32
}

type Quit struct {

}

func getFloat64(obj interface{}) float64 {
	switch i := obj.(type) {
	case float64:
		return i
	default:
		panic(fmt.Sprintf("object is not float (%+v)", i))
	}
}

func getInt32(obj interface{}) int32 {
	return int32(getFloat64(obj))
}

// Parses the type of marcus from string
func ParsePacket(mp map[string]interface{}) (ret interface{}, err error) {
	packType, _ := mp["type"]

	defer func() {
		if r := recover(); r != nil {
			ret = nil
			err = errors.New(fmt.Sprintf("PlayerPacket Error: %s", r))
		}
	}()

	switch packType {
	case "location":
		ret = Location{getFloat64(mp["latitude"]), getFloat64(mp["longitude"])}
		err = nil
	case "setTeam":
		ret = SetTeam{getInt32(mp["team"])}
		err = nil
	case "quit":
		ret = Quit{}
		err = nil
	default:
		panic(fmt.Sprintf("Invalid packet (%s)", packType))
	}

	return
}

// Finds Euclidian distance to nearest marcus
func (x *Location) Dist(oth *Location) float64 {
	xDelta := x.longitude - oth.longitude
	yDelta := x.latitude - oth.latitude
	return math.Sqrt(xDelta * xDelta + yDelta * yDelta)
}
