package number

import "math"

// IFloorDiv 整除 向负无穷取整
func IFloorDiv(a, b int64) int64 {
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	}
	// 向负无穷取整
	return a/b - 1
}

// FFloorDiv 浮点数整除 向负无穷取整
func FFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

// IMod 取余 利用整除实现
func IMod(a, b int64) int64 {
	return a - IFloorDiv(a, b)*b
}

// FMod 浮点数取余 利用整除实现
func FMod(a, b float64) float64 {
	return a - FFloorDiv(a, b)*b
}

// ShiftLeft 按位左移
func ShiftLeft(a, n int64) int64 {
	if n >= 0 {
		return a << uint64(n)
	}
	return ShiftRight(a, -n)
}

// ShiftRight 按位右移
func ShiftRight(a, n int64) int64 {
	if n >= 0 {
		return int64(uint64(a) >> uint64(n))
	}
	return ShiftLeft(a, -n)
}

// FloatToInteger 将浮点数转换为整数
func FloatToInteger(f float64) (int64, bool) {
	i := int64(f)
	return i, float64(i) == f
}
