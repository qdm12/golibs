package random

type Rand interface {
	// Int returns a pseudo random int number.
	Int() int
	// Intn returns a pseudo random int number in the range [0..maxN).
	Intn(maxN int) int
}
