package utils

import "strconv"

func IsValidOrderId(orderId string) bool {
	number, err := strconv.ParseInt(orderId, 10, 64)
	if err != nil {
		return false
	}
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int64) int64 {
	var luhn int64

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
