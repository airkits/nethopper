package utils

// PowerCalc return size & power
func PowerCalc(size int32) (int32, uint8) {

	if size <= 0 {
		return 0, 0
	}
	power := uint8(0)
	value := size
	for {
		if value <= 1 {
			break
		}
		value >>= 1
		power++
	}
	if size&(size-1) == 0 { //is power of 2
		return 1 << power, power
	}
	power++
	return 1 << power, power
}
