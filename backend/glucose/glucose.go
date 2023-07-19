package glucose

import "math"

const mmolformula float32 = 0.0555555555555556

func MmolToMg(mmol float32) int {
	res := mmol / mmolformula
	return int(res + 0.5)
}

func MgToMmol(mg int) float32 {
	res := float32(mg) * mmolformula
	return float32(math.Round(float64(res)/0.05) * 0.05)
}
