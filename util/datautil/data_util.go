package datautil

import "zyj.com/golang-study/util/varutil"

// MergeConfig 聚合配置
type JoinType string

const (
	InnerJoin JoinType = "InnerJoin"
	LeftJoin  JoinType = "LeftJoin"
	RightJoin JoinType = "RightJoin"
	OuterJoin JoinType = "OuterJoin"
)

// MergeSlicesGeneric 通用的切片聚合函数
func MergeSlicesGeneric[T any, U any, K comparable](
	left []T,
	right []U,
	leftKeyFunc func(T) K,
	rightKeyFunc func(U) K,
	joinType JoinType,
) []map[string]interface{} {

	// 创建右切片的映射
	rightMap := make(map[K]U)
	for _, item := range right {
		key := rightKeyFunc(item)
		rightMap[key] = item
	}

	var result []map[string]interface{}
	matchedKeys := make(map[K]bool)

	// 处理左切片
	for _, leftItem := range left {
		leftKey := leftKeyFunc(leftItem)
		if rightItem, exists := rightMap[leftKey]; exists {
			// 内连接匹配
			result = append(result, map[string]interface{}{
				"key":   leftKey,
				"left":  leftItem,
				"right": rightItem,
				"match": true,
			})
			matchedKeys[leftKey] = true
		} else if joinType == LeftJoin || joinType == OuterJoin {
			// 左连接或外连接
			result = append(result, map[string]interface{}{
				"key":   leftKey,
				"left":  leftItem,
				"right": nil,
				"match": false,
			})
		}
	}
	if joinType == InnerJoin {
		return result
	}
	// 处理右切片中未匹配的元素（右连接或外连接）
	if joinType == RightJoin || joinType == OuterJoin {
		for key, rightItem := range rightMap {
			if !matchedKeys[key] {
				result = append(result, map[string]interface{}{
					"key":   key,
					"left":  nil,
					"right": rightItem,
					"match": false,
				})
			}
		}
	}

	return result
}
func LeftJoinData[L any, R any, O any, K comparable](
	left []L,
	right []R,
	leftKeyFunc func(L) K,
	rightKeyFunc func(R) K,
	getSet func(L, R) O,
) []O {
	// 创建右切片的映射
	rightMap := make(map[K]R)
	for _, item := range right {
		key := rightKeyFunc(item)
		rightMap[key] = item
	}

	var outs = make([]O, 0, len(left))
	// 处理左切片
	for _, leftItem := range left {
		leftKey := leftKeyFunc(leftItem)
		if rightItem, exists := rightMap[leftKey]; exists {
			outs = append(outs, getSet(leftItem, rightItem))
		}
	}

	return outs
}

func JoinData[L any, R any, O any, K comparable](
	left []L,
	right []R,
	leftKeyFunc func(L) K,
	rightKeyFunc func(R) K,
) []O {
	// 创建右切片的映射
	rightMap := make(map[K]R)
	for _, item := range right {
		key := rightKeyFunc(item)
		rightMap[key] = item
	}
	var outs = make([]O, 0, len(left))
	// 处理左切片
	for _, leftItem := range left {
		leftKey := leftKeyFunc(leftItem)
		if rightItem, exists := rightMap[leftKey]; exists {
			var out O
			varutil.CopyTo(leftItem, &out)
			varutil.CopyTo(rightItem, &out)
			outs = append(outs, out)
		}
	}
	return outs
}

func MapKeyObject[T any, K comparable](
	Objs []T,
	leftKeyFunc func(T) K,
) map[K]T {
	rightMap := make(map[K]T, len(Objs))
	for _, item := range Objs {
		key := leftKeyFunc(item)
		rightMap[key] = item
	}
	return rightMap
}

func MapKeyObjectPtr[T any, K comparable](
	Objs []T,
	leftKeyFunc func(T) K,
) map[K]*T {
	rightMap := make(map[K]*T, len(Objs))
	for _, item := range Objs {
		key := leftKeyFunc(item)
		rightMap[key] = &item
	}
	return rightMap
}
