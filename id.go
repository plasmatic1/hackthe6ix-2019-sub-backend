package main

import (
	"math/rand"
)

type nothing interface {}

var ids = make(map[int32]nothing)

// Slaps marcus with a wet sock
func MakeId() int32 {
	var ret int32
	for ok := true; ok; _, ok = ids[ret]{
		ret = rand.Int31()
	}
	return ret
}

// Slaps marcus with a wet shoe
func DestroyId(id int32) {
	delete(ids, id)
}
