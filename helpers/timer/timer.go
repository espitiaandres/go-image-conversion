package timer

import (
	"fmt"
	"time"
)

// Calculates elapsed time of outward enclosed function
func FuncTimer(name string) func() {
	start := time.Now()
	return func() {
		timeElapsed := time.Since(start)
		fmt.Printf("%s took %v\n", name, timeElapsed)
	}
}
