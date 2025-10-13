package main

import (
	"fmt"
	xerrors "github.com/pkg/errors"
)

func main() {
	/*
		%s,%v 功能一样，输出错误信息，不包含堆栈
		%q 输出的错误信息带引号，不包含堆栈
		%+v 输出错误信息和堆栈
	*/
	e1 := xerrors.New("msg1") // 新生成错误e1，带堆栈信息和msg1
	e2 := xerrors.WithMessage(e1, "msg2")
	e3 := xerrors.WithMessage(e2, "msg3")
	fmt.Printf("%+v\n\n", e1) // 打印底层堆栈信息和msg1
	fmt.Printf("%+v\n\n", e2) // 打印底层堆栈信息和msg1，msg2
	fmt.Printf("%+v\n\n", e3) // 打印底层堆栈信息和msg1，msg2，msg3
	e4 := xerrors.Cause(e3)
	fmt.Printf("msg4 %+v\n\n", e4) // // 打印底层堆栈信息和msg1
}
