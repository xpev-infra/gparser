package gparser

import (
	"errors"
	"go/ast"
	"go/token"
	"strings"
)

// 注册可执行函数
var funcNameMap = map[string]func(args []ast.Expr, data map[string]interface{}) interface{}{}

func init() {
	funcNameMap = map[string]func(args []ast.Expr, data map[string]interface{}) interface{}{
		"in_array":             inArray,
		"start_with":           startWith,
		"in_organization":      inOrganization,
		"contain_organization": containOrg,
		"compare_version":      compareVersion,
	}
}

func startWith(args []ast.Expr, data map[string]interface{}) interface{} {
	// 规则表达式中的变量
	param := eval(args[0], data)

	exp, ok := args[1].(*ast.BasicLit)
	if !ok {
		return errors.New("func start_with error usage")
	}

	if exp.Kind != token.STRING {
		return errors.New("func start_with only support string")
	}

	exp.Value = strings.Trim(exp.Value, "\"")

	return strings.HasPrefix(param.(string), exp.Value)
}

// inOrganization
// @Title inOrganization
// @Description: 定制方法，判断用户所属的部门是否和规则中的匹配，本来也可以用start_with替代这个功能，但为了功能表达清晰，多加一个方法
// @param args
// @param data
// @return interface{}
func inOrganization(args []ast.Expr, data map[string]interface{}) interface{} {
	param := eval(args[0], data)

	exp, ok := args[1].(*ast.BasicLit)
	if !ok {
		return errors.New("func start_with error usage")
	}

	if exp.Kind != token.STRING {
		return errors.New("func start_with only support string")
	}

	exp.Value = strings.Trim(exp.Value, "\"")

	// 用户所属的部门路径，如   1111111/222222/3333333333 表示 部门3
	inputValues := strings.Split(param.(string), "/")

	// 规则中写的部门路径 如 1111111/222222 表示A中心
	ruleValues := strings.Split(exp.Value, "/")

	// 用户所属的部门路径长度如果都小于规则中的路径，那么一定是不匹配的
	if len(inputValues) < len(ruleValues) {
		return false
	}

	// 如果用户属于规则中的部门，那么规则中的路径前缀和用户的一定是相同的
	// 如 规则中 配置 1111111/222222 表示A中心
	// 用户所在  1111111/222222/3333333333 表示A中心下的部门3
	for i, org := range ruleValues {
		if org != inputValues[i] {
			return false
		}
	}

	return true

}

// inArray 判断变量是否存在在数组中
func inArray(args []ast.Expr, data map[string]interface{}) interface{} {
	// 规则表达式中的变量
	param := eval(args[0], data)
	vRange, ok := args[1].(*ast.CompositeLit)
	if !ok {
		return errors.New("func in_array 2ed params is not a composite lit")
	}

	// 规则表达式中数组里的元素
	eltNodes := make([]interface{}, 0, len(vRange.Elts))
	for _, p := range vRange.Elts {
		elt := eval(p, data)
		eltNodes = append(eltNodes, elt)
	}

	for _, node := range eltNodes {
		switch node.(type) {
		case int64:
			param, err := castType(param, TypeInt64)
			if err != nil {
				return false
			}
			paramInt64, paramOk := param.(int64)
			nodeInt64, nodeOk := node.(int64)
			if !paramOk || !nodeOk {
				return false
			}
			if nodeInt64 == paramInt64 {
				return true
			}
		case string:
			param, err := castType(param, TypeString)
			if err != nil {
				return false
			}
			nodeString, nodeOk := node.(string)
			paramString, paramOk := param.(string)
			if !paramOk || !nodeOk {
				return false
			}
			if nodeString == paramString {
				return true
			}
		}
	}
	return false
}

// containOrg
// @Title containOrg
// @Description: 是否包含组织
// @return args
// @return data
func containOrg(args []ast.Expr, data map[string]interface{}) interface{} {
	// 规则表达式中的变量
	param := eval(args[0], data)
	vRange, ok := args[1].(*ast.CompositeLit)
	if !ok {
		return false
	}

	eltNodes := make([]interface{}, 0, len(vRange.Elts))
	for _, p := range vRange.Elts {
		elt := eval(p, data)
		eltNodes = append(eltNodes, elt)
	}

	param, err := castType(param, TypeString)
	if err != nil {
		return false
	}

	paramString, ok := param.(string)
	if !ok {
		return false
	}

	for _, node := range eltNodes {
		nodeString, ok := node.(string)

		if !ok {
			return false
		}

		// 匹配中某个部门
		if strings.HasPrefix(paramString, nodeString) {
			return true
		}

	}

	return false

}

func compareVersion(args []ast.Expr, data map[string]interface{}) interface{} {
	// 规则表达式中的变量
	param := eval(args[0], data)
	versionLit, ok := args[1].(*ast.BasicLit)
	if !ok {
		return false
	}

	operate, ok := args[2].(*ast.BasicLit)
	if !ok {
		return false
	}

	param, err := castType(param, TypeString)
	if err != nil {
		return false
	}

	version := strings.Trim(versionLit.Value, "\"")

	result, err := versionCompare(param.(string), version)
	if err != nil {
		return err
	}

	switch strings.Trim(operate.Value, "\"") {
	case ">":
		return result == 1
	case ">=":
		return result == 1 || result == 0
	case "=":
		return result == 0
	case "!=":
		return result != 0
	case "<":
		return result == -1
	case "<=":
		return result == -1 || result == 0
	default:
		return false
	}

}
