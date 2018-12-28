package main

import "testing"

//1.功能测试  以Test开头
//2.压力测试  以Benchmark开头
//3.测试代码覆盖率的测试

//translator 只用于性能测试，用户输入从main函数判断

//功能测试：main函数的用户输入校验测试
func TestMain(t *testing.T) {
	query, from, to, err := checkInput("nihao/zh/wyw")
	if query != "nihao" {
		t.Error("翻译字符串错误")
	}
	if from != "wyw" {
		t.Error("源语言解析错误")
	}
	if to != "zh" {
		t.Error("目标语言解析错误")
	}
	if err != nil {
		t.Error("判断出错")
	}
}
func TestCheckInput1(t *testing.T) {
	query, from, to, err := checkInput("nihao")
	if query != "nihao" {
		t.Error("翻译字符串错误")
	}
	if from != "auto" {
		t.Error("源语言解析错误")
	}
	if to != "zh" {
		t.Error("目标语言解析错误")
	}
	if err != nil {
		t.Error("判断出错")
	}

}

//
func TestCheckInput3(t *testing.T) {
	query, from, to, err := checkInput("nihao/zh/en")
	if query != "nihao" {
		t.Error("翻译字符串错误")
	}
	if from != "en" {
		t.Error("源语言解析错误")
	}
	if to != "zh" {
		t.Error("目标语言解析错误")
	}
	if err != nil {
		t.Error("判断出错")
	}
}

func TestCheckInput4(t *testing.T) {
	query, from, to, err := checkInput("你猜/ab")
	if query != "你猜" {
		t.Error("翻译字符串错误")
	}
	if from != "auto" {
		t.Error("源语言解析错误")
	}
	if to != "en" {
		t.Error("目标语言解析错误")
	}
	if err != nil {
		t.Error("判断出错")
	}
}
func TestCheckInput5(t *testing.T) {
	query, from, to, err := checkInput("你猜")
	if query != "你猜" {
		t.Error("翻译字符串错误")
	}
	if from != "auto" {
		t.Error("源语言解析错误")
	}
	if to != "en" {
		t.Error("目标语言解析错误")
	}
	if err != nil {
		t.Error("判断出错")
	}
}
