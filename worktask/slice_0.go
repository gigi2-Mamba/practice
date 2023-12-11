package main

import "fmt"

/*
根据下标，删除切片元素
要求1： 实现删除
要求2： 实现高性能删除，避免声明赋值和扩容
要求3： 实现泛型
要求4： 缩容
*/

// 根据具体下标做切割，尽量避免赋值，纯切割可以吗  应该不行

func Delete(arr []int, index int) []int {
	newarr := make([]int, 0, len(arr)-1)

	newarr = append(newarr, arr[:index]...)
	newarr = append(newarr, arr[index+1:]...)

	return newarr

}

func DeleteGeneric[T any](vals []T, idx int) []T {
	newArr := make([]T, 0, len(vals)-1)
	newArr = append(newArr, vals[:idx]...)
	newArr = append(newArr, vals[idx+1:]...)

	return newArr
}

func main() {
	a := []int{1, 2, 3, 4, 5}
	ints := Delete(a, 2)
	fmt.Println(ints)
	a1 := []string{"a", "b", "c"}
	a1s := DeleteGeneric(a1, 1)
	fmt.Println(a1s)
}
