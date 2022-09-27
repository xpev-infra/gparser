package gparser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// 计算int类型表达式
func calculateForInt(x, y interface{}, op token.Token) interface{} {
	x, err := castType(x, TypeInt64)
	if err != nil {
		return err
	}
	y, err = castType(y, TypeInt64)
	if err != nil {
		return err
	}
	return calculateForInt64(x, y, op)
}

// 计算int64类型表达式
func calculateForInt64(x, y interface{}, op token.Token) interface{} {
	x, err := castType(x, TypeInt64)
	if err != nil {
		return err
	}
	xInt, xok := x.(int64)
	yInt, yok := y.(int64)
	if !xok || !yok {
		return errors.New(fmt.Sprintf("%v %v %v eval failed", x, op, y))
	}

	// 计算逻辑
	switch op {
	case token.EQL:
		return xInt == yInt
	case token.NEQ:
		return xInt != yInt
	case token.GTR:
		return xInt > yInt
	case token.LSS:
		return xInt < yInt
	case token.GEQ:
		return xInt >= yInt
	case token.LEQ:
		return xInt <= yInt
	case token.ADD:
		return xInt + yInt
	case token.SUB:
		return xInt - yInt
	case token.MUL:
		return xInt * yInt
	case token.QUO:
		if yInt == 0 {
			return 0
		}
		return xInt / yInt
	default:
		return errors.New(fmt.Sprintf("unsupported binary operator: %s", op.String()))
	}
}

// 计算string类型表达式
func calculateForString(x, y interface{}, op token.Token) interface{} {
	x, err := castType(x, TypeString)
	if err != nil {
		return err
	}
	xString, xok := x.(string)
	yString, yok := y.(string)
	if !xok || !yok {
		return errors.New(fmt.Sprintf("%v %v %v eval failed", x, op, y))
	}

	// 计算逻辑
	switch op {
	case token.EQL: // ==
		return xString == yString
	case token.NEQ: // !=
		return xString != yString
	case token.LSS: // <
		return xString < yString
	case token.GTR: // >
		return xString > yString
	case token.LEQ: // <=
		return xString <= yString
	case token.GEQ: // >=
		return xString >= yString
	}
	return errors.New(fmt.Sprintf("unsupported binary operator: %s", op.String()))
}

// 计算bool类型表达式
func calculateForBool(x, y interface{}, op token.Token) interface{} {
	x, err := castType(x, TypeBool)
	if err != nil {
		return err
	}
	xb, xok := x.(bool)
	yb, yok := y.(bool)
	if !xok || !yok {
		return errors.New(fmt.Sprintf("%v %v %v eval failed", x, op, y))
	}

	// 计算逻辑
	switch op {
	case token.LAND:
		return xb && yb
	case token.LOR:
		return xb || yb
	case token.EQL:
		return xb == yb
	case token.NEQ:
		return xb != yb
	}
	return errors.New(fmt.Sprintf("unsupported binary operator: %s", op.String()))
}

// calculateForFunc 计算函数表达式
func calculateForFunc(funcName string, args []ast.Expr, data map[string]interface{}) interface{} {
	// 根据funcName分发逻辑
	handler, ok := funcNameMap[funcName]
	if !ok {
		return errors.New(fmt.Sprintf("%+v func not support", funcName))
	}
	return handler(args, data)
}

// versionCompare
// @Description: 比较 a b 版本号大小
// @param version1 eg 1.0.0
// @param version2 eg 0.
// @return int  a<b -1, a=b 0 , a>b 1
// @return error
func versionCompare(version1 string, version2 string) (int, error) {
	var res int
	ver1Str := strings.Split(version1, ".")
	ver2Str := strings.Split(version2, ".")
	ver1Len := len(ver1Str)
	ver2Len := len(ver2Str)
	verLen := ver1Len
	if len(ver1Str) < len(ver2Str) {
		verLen = ver2Len
	}
	var err error
	for i := 0; i < verLen; i++ {
		var ver1Int, ver2Int int
		if i < ver1Len {
			ver1Int, err = strconv.Atoi(ver1Str[i])
			if err != nil {
				return 0, err
			}
		}
		if i < ver2Len {
			ver2Int, err = strconv.Atoi(ver2Str[i])
			if err != nil {
				return 0, err
			}
		}
		if ver1Int < ver2Int {
			res = -1
			break
		}
		if ver1Int > ver2Int {
			res = 1
			break
		}
	}
	return res, nil
}
