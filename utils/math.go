package utils

import (
	"math"
	"math/rand"
	"time"
)

//CreateUniqRandArray 创建唯一随机数 [0,max)数组，返回一个数组序列，没有重复数
func CreateUniqRandArray(max int, length int) []int {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	set := make(map[int]struct{})

	nums := make([]int, 0, length)
	for {
		num := rnd.Intn(max)
		if _, ok := set[num]; !ok {
			set[num] = struct{}{}
			nums = append(nums, num)
		}
		if len(nums) == length {
			return nums
		}
	}

}

//CreateRandArray 创建随机数 [0,max)数组，返回随机数序列，可能会有重复
func CreateRandArray(max int, length int) []int {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	nums := make([]int, 0, length)
	for {
		num := rnd.Intn(max)
		nums = append(nums, num)
		if len(nums) == length {
			return nums
		}
	}

}

//Round int32取整 四舍五入
func Round(value float64) int32 {
	return int32(value + 0.5)
}

//Round64 int64取整 四舍五入
func Round64(value float64) int64 {
	return int64(value + 0.5)
}

//RoundFloat64 float64取整 四舍五入
func RoundFloat64(value float64) float64 {
	return float64(int64(value + 0.5))
}

//InvSqrt 平方根倒数速算法
func InvSqrt(x float32) float32 {
	var xhalf float32 = 0.5 * x // get bits for floating VALUE
	i := math.Float32bits(x)    // gives initial guess y0
	i = 0x5f375a86 - (i >> 1)   // convert bits BACK to float
	x = math.Float32frombits(i) // Newton step, repeating increases accuracy
	x = x * (1.5 - xhalf*x*x)
	x = x * (1.5 - xhalf*x*x)
	x = x * (1.5 - xhalf*x*x)
	return 1 / x
}

//Abs int取绝对值
func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

//Abs32 int32取绝对值
func Abs32(i int32) int32 {
	if i < 0 {
		return -i
	}
	return i
}

//IMax 取2数最大
func IMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//I32Max 取2数最大
func I32Max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

//I64Max 取2数最大
func I64Max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

//IMin 取2数最小
func IMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//I32Min 取2数最小
func I32Min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

//I64Min 取2数最小
func I64Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

//IClamp 取中间值
func IClamp(value, min, max int) int {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}

//I32Clamp 取中间值
func I32Clamp(value, min, max int32) int32 {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}

//I64Clamp 取中间值
func I64Clamp(value, min, max int64) int64 {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}

//RandomInt 随机范围取值[from, to]
func RandomInt(from, to int) int {
	return rand.Intn(to-from+1) + from
}

//RandomInt32 随机范围取值[from, to]
func RandomInt32(from, to int32) int32 {
	var i = rand.Intn(int(to)-int(from)+1) + int(from)
	return int32(i)
}

//Shuffle 随机乱序,洗牌
func Shuffle(list []interface{}) {
	var c = len(list)
	if c < 2 {
		return
	}

	for i := 0; i < c-1; i++ {
		var j = rand.Intn(c-i) + i //这里rand需要包含i自身，否则不均匀
		if i != j {
			list[i], list[j] = list[j], list[i]
		}
	}
}

//ShuffleI32 随机乱序
func ShuffleI32(list []int32) {
	var c = len(list)
	if c < 2 {
		return
	}

	for i := 0; i < c-1; i++ {
		var j = rand.Intn(c-i) + i //这里rand需要包含i自身，否则不均匀
		if i != j {
			list[i], list[j] = list[j], list[i]
		}
	}
}

//ShuffleI 随机乱序
func ShuffleI(list []int) {
	var c = len(list)
	if c < 2 {
		return
	}

	for i := 0; i < c-1; i++ {
		var j = rand.Intn(c-i) + i //这里rand需要包含i自身，否则不均匀
		if i != j {
			list[i], list[j] = list[j], list[i]
		}
	}
}

//ShuffleN 随机抽取n张
func ShuffleN(list []interface{}, randCount int) []interface{} {
	var c = len(list)
	if c < 2 {
		return list
	}

	var ct = IMin(c-1, randCount)
	for i := 0; i < ct; i++ {
		var j = rand.Intn(c-i) + i //这里rand需要包含i自身，否则不均匀
		if i != j {
			list[i], list[j] = list[j], list[i]
		}
	}

	return list[:ct]
}

//ShuffleNI32 随机抽取n张
func ShuffleNI32(list []int32, randCount int) []int32 {
	var c = len(list)
	if c < 2 {
		return list
	}

	var ct = IMin(c-1, randCount)
	for i := 0; i < ct; i++ {
		var j = rand.Intn(c-i) + i //这里rand需要包含i自身，否则不均匀
		if i != j {
			list[i], list[j] = list[j], list[i]
		}
	}
	return list[:ct]
}

//ShuffleR 随机乱序,洗牌，反向，效果一样
func ShuffleR(list []interface{}) {
	var c = len(list)
	if c < 2 {
		return
	}

	for i := c - 1; i >= 1; i-- {
		var j = rand.Int() % (i + 1)
		list[i], list[j] = list[j], list[i]
	}
}

//SumI32 数组求和
func SumI32(list []int32) int32 {
	var sum int32
	for _, v := range list {
		sum += v
	}
	return sum
}

//SumMatrixColI32 矩阵列求和
func SumMatrixColI32(mat [][]int32, col int) int32 {
	var list []int32
	for index := 0; index < len(mat); index++ {
		list = append(list, mat[index][col])
	}
	return SumI32(list)
}

//StaticRand 固定种子伪随机
func StaticRand(seedrare, min, max int) int {
	var seed = float64(seedrare)
	seed = seed*2045 + 1
	seed = float64(int(seed) % 1048576)
	var dis = float64(max - min)
	var ret = int(min) + int(math.Floor(seed)*dis/1048576)
	return ret
}

//RandomMultiWeight 随机多个权重
func RandomMultiWeight(weightMapping map[int32]int32, count int) []int32 {
	results := []int32{}
	for i := 0; i < count; i++ {
		result := RandomWeight(weightMapping)
		if result == 0 {
			return results
		}
		results = append(results, result)
		delete(weightMapping, result)
	}
	return results
}

//RandomWeight 随机权重
func RandomWeight(weightMapping map[int32]int32) int32 {
	var totalWeight int32 = 0
	for _, weight := range weightMapping {
		totalWeight += weight
	}
	if totalWeight <= 0 {
		return 0
	}
	randomValue := rand.Int31n(totalWeight) + 1
	var currentValue int32 = 0
	for id, weight := range weightMapping {
		currentValue += weight
		if currentValue >= randomValue {
			return id
		}
	}
	return 0
}
