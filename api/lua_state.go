package api

type LuaType = int

type ArithOp = int
type CompareOp = int

type LuaState interface {

	/* 基础栈操作 */

	// GetTop 获取栈顶索引
	GetTop() int
	// AbsIndex 将相对索引转换为绝对索引
	AbsIndex(index int) int
	// CheckStack 确保 Lua 栈可以存入 n 个元素，若不能存入则扩展栈空间
	CheckStack(n int) bool
	// Pop 弹出栈中 n 个元素
	Pop(n int)
	// Copy 将 from 位置的元素复制到 to
	Copy(from, to int)
	// PushValue 将 index 处的元素复制并压入栈顶
	PushValue(index int)
	// Replace 将栈顶值弹出并写入 index 位置
	Replace(index int)
	// Insert 将栈顶值弹出并插入到 index 位置，其他值依次后移
	Insert(index int)
	// Remove 移除 index 位置元素
	Remove(index int)
	// Rotate 将 [index, top] 区间的值向栈顶方向旋转 n 次，n < 0 则逆向旋转
	Rotate(index, n int)
	// SetTop 设置栈顶索引
	SetTop(index int)

	/* 索引访问函数 */

	// TypeName 获取类型名
	TypeName(tp LuaType) string
	// Type 获取给定 index 值的类型，若下标错返回 LUA_TNONE
	Type(index int) LuaType
	// IsNone 判断给定 index 值是否为 LUA_TNONE
	IsNone(index int) bool
	// IsNil 判断给定 index 值是否为 Nil
	IsNil(index int) bool
	// IsNoneOrNil 判断指定 index 是否有值
	IsNoneOrNil(index int) bool
	// IsBoolean 判断指定 index 是否为布尔，或可被转换为布尔
	IsBoolean(index int) bool
	// IsInteger 判断指定 index 是否为整型，或可被转换为整型
	IsInteger(index int) bool
	// IsNumber 判断给定 index 是否为数字，或可被转换为数字
	IsNumber(index int) bool
	// IsString 判断给定 index 是否为字符串，或可被转换为字符串
	IsString(index int) bool
	// ToBoolean 将给定 index 值转换为布尔值
	ToBoolean(index int) bool
	// ToInteger 将给定 index 值转换为整型
	ToInteger(index int) int64
	// ToIntegerX 将给定 index 值转换为整型，返回转换是否成功
	ToIntegerX(index int) (int64, bool)
	// ToNumber 将给定 index 值转换为数字
	ToNumber(index int) float64
	// ToNumberX 将给定 index 值转换为数字，返回转换是否成功
	ToNumberX(index int) (float64, bool)
	// ToString 将给定 index 值转换为字符串
	ToString(index int) string
	// ToStringX 将给定 index 值转换为数字，返回转换是否成功
	ToStringX(index int) (string, bool)

	/* 入栈函数 */

	// PushNil 将一个 nil 入栈
	PushNil()
	// PushBoolean 将一个 bool 入栈
	PushBoolean(b bool)
	// PushInteger 将一个整型入栈
	PushInteger(n int64)
	// PushNumber 将一个浮点数字入栈
	PushNumber(n float64)
	// PushString 将一个字符串入栈
	PushString(s string)

	/* 运算函数 */

	// Arith 按位运算、算术运算
	// 从栈顶弹出运算数并将结果压栈
	Arith(op ArithOp)
	// Compare 比较运算
	// 比较栈中两个位置的变量，不改变栈状态
	Compare(index1, index2 int, op CompareOp) bool
	// Len 长度计算，获取列表长度 #
	// 从栈中指定索引获取值，取值的长度并压栈
	Len(index int)
	// Concat 字符串拼接 ..
	// 从栈顶弹出 n 个值，拼接成一个字符串后压栈
	Concat(n int)
}
