package large

import (
	"log"
	"math"
	"testing"
)

func TestLarge(t *testing.T) {
	a := math.MaxInt64
	log.Println(a, a+1)
	b := math.MaxFloat64
	log.Println(b, b+1)
}
