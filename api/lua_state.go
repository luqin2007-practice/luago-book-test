package api

type LuaType = int

type ArithOp = int
type CompareOp = int

type GoFunction func(state LuaState) int

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
	// LuaUpvalueIndex 将注册表伪索引转换为 Upvalue 索引
	LuaUpvalueIndex(i int) int

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
	// IsGoFunction 判断给定 index 是否为 Go 函数
	IsGoFunction(index int) bool
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
	// ToGoFunction 将给定 index 值转换为 Go 函数，返回转换是否成功
	ToGoFunction(index int) GoFunction

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
	// PushGoFunction 将一个 Go 函数入栈
	PushGoFunction(f GoFunction, n int)
	// PushGlobalTable 将全局变量表推入栈顶
	PushGlobalTable()

	/* 运算函数 */

	// Arith 按位运算、算术运算
	//   从栈顶弹出运算数并将结果压栈
	Arith(op ArithOp)
	// Compare 比较运算
	//   比较栈中两个位置的变量，不改变栈状态
	Compare(index1, index2 int, op CompareOp) bool
	// Len 长度计算，获取列表长度 #
	//   从栈中指定索引获取值，取值的长度并压栈
	Len(index int)
	// Concat 字符串拼接 ..
	//   从栈顶弹出 n 个值，拼接成一个字符串后压栈
	Concat(n int)

	/* 表相关 */

	// NewTable 创建无法预估大小的表，并放入栈顶
	//   等价于 CreateTable(0,0)
	NewTable()
	// CreateTable 创建表，并放入栈顶
	//   nArr: 预估列表部分长度
	//   nRec: 预估记录部分长度
	CreateTable(nArr, nRec int)
	// GetTable 获取表元素
	//   将栈顶元素作为键，从 index 位置的表中获取数据并放入栈顶
	//   返回表元素的类型
	GetTable(index int) LuaType
	// GetField 获取表元素
	//   将给定字符串，从 index 位置的表中获取数据并放入栈顶
	//   返回表元素的类型
	GetField(index int, k string) LuaType
	// GetI 获取表元素
	//   将给定字符串，从 index 位置的表中获取数据并放入栈顶
	//   返回表元素的类型
	GetI(index int, i int64) LuaType
	// GetGlobal 使用字符串从全局变量表中取值
	//   比 PushGlobalTable + GetField 效率高一点
	GetGlobal(name string) LuaType
	// SetTable 将栈顶两个元素作为 k v 存入 index 位置的表中
	SetTable(index int)
	// SetField 将栈顶元素和给定字符串键存入 index 位置的表中
	SetField(index int, k string)
	// SetI 将栈顶元素和给定数字键存入 index 位置的表中
	SetI(index int, i int64)
	// SetGlobal 将栈顶元素存入全局变量表
	//   比 PushGlobalTable + SetField 效率高一点
	SetGlobal(name string)
	// Register 在全局变量表中添加 Go 函数
	Register(name string, f GoFunction)

	/* 函数调用 */

	// Load 加载编译过的二进制脚本或 Lua 脚本
	//   mode 表示加载的数据类型，b 表示二进制，t 表示 lua 脚本，bt 表示二者皆可
	//   返回值 0 表示加载成功
	Load(chunk []byte, chunkName string, mode string) int
	// Call 调用函数应先把函数推入栈顶，然后将参数依次推入栈顶
	//   nArgs 实际传入的参数数量
	//   nResults 实际需要的参数数量
	Call(nArgs, nResults int)

	/* 元表 */

	// GetMetatable 检查 index 处值是否有元表，若有则入栈，否则不改变栈状态
	//   返回值：值是否包含元表
	GetMetatable(index int) bool
	// SetMetatable 将栈顶 Table 弹出，将其设置为 index 处值的元表
	SetMetatable(index int)
	// RawLen 等效 Len，忽略元方法
	RawLen(index int) uint
	// RawEqual 等效 Equal，忽略元方法
	RawEqual(index1, index2 int) bool
	// RawGet 等效 GetTable，忽略元方法
	RawGet(index int) LuaType
	// RawSet 等效 SetTable，忽略元方法
	RawSet(index int)
	// RawGetI 等效 GetI，忽略元方法
	RawGetI(index int, i int64) LuaType
	// RawSetI 等效 SetI，忽略元方法
	RawSetI(index int, i int64)
}
