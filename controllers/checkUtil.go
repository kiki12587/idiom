/*
@Time : 2020/11/30 20:12
@Author : HP
@File : CheckUtil
@Software: GoLand
*/
package controllers

func CheckSlice(data []string) bool {
	if len(data) == 0 {
		return true
	} else {
		return false
	}
}

func CheckString(data string) bool {
	if len(data) == 0 {
		return true
	} else {
		return false
	}
}

func CheckMap(data []map[string]string) bool {
	if len(data) == 0 {
		return true
	} else {
		return false
	}
}
