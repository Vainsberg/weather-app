package formatint

import "math"

func FormatInt(num float64) int {
	return int(math.Ceil(num))
}
