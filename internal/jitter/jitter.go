package jitter

import "math/rand/v2"

func GetJitter(jitter int) int {
	if jitter <= 0 {
		return 0
	}

	return rand.IntN(jitter*2+1) - jitter
}
