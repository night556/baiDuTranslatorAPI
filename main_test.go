package main

import (
	"errors"
	"testing"
)

//1.功能测试  以Test开头
//2.压力测试  以Benchmark开头
//3.测试代码覆盖率的测试
//translator 只用于性能测试，用户输入从main函数判断
//功能测试：main函数的用户输入校验测试

//TestCheckInput1 测试英文
func TestCheckInput1(t *testing.T) {
	query, from, to := checkInput("hello")
	if query != "hello" {
		t.Error("翻译字符串错误")
	}
	if from != "auto" {
		t.Error("源语言解析错误")
	}
	if to != "zh" {
		t.Error("目标语言解析错误")
	}
}

//TestCheckInput2 测试自动转换，中英互译
func TestCheckInput2(t *testing.T) {
	query, from, to := checkInput("你好")
	if query != "你好" {
		t.Error("翻译字符串错误")
	}
	if from != "auto" {
		t.Error("源语言解析错误")
	}
	if to != "en" {
		t.Error("目标语言解析错误")
	}
}

//TestTranslator1 检测正常的翻译
func TestTranslator1(t *testing.T) {
	result, err := translator("hello", "auto", "zh")
	resultRight, ok := result.(Result)
	if err != nil {
		t.Error("非正常error")
	}
	if !ok {
		t.Error("类型断言失败")
	}
	if resultRight.From != "en" || resultRight.To != "zh" || resultRight.TransResults == nil {
		t.Error("翻译逻辑失败")
	}
}

//TestTranslator2 检测正常的翻译
func TestTranslator2(t *testing.T) {
	result, err := translator("hello", "", "zh")
	resultRight, ok := result.(ErrorResponse)
	if err == nil {
		t.Error("非正常error")
	}
	if !ok {
		t.Error("类型断言失败")
	}
	if resultRight.ErrorCode != "54000" || resultRight.ErrorMsg != "PARAM_FROM_TO_OR_Q_EMPTY" {

		t.Error("翻译逻辑失败")
	}
}

func TestPrintResult(t *testing.T) {
	result := Result{From: "en", To: "zh", TransResults: nil}
	str := printResult(result, nil)
	if str != "" {
		t.Error("打印翻译成功结果失败")
	}
	errmsg := ErrorResponse{ErrorCode: "54000", ErrorMsg: "ERROR "}
	err := errors.New("nicai")
	str = printResult(errmsg, err)
	if str != "Error Code:"+errmsg.ErrorCode+"\tError Msg: "+errmsg.ErrorMsg+"\n" {
		t.Error("打印翻译错误字符串失败")
	}
}

//TestCli1 检测!help命令
func TestCli1(t *testing.T) {
	err := Cli("!help")
	if err != nil {
		t.Error("!help 命令出错")
	}
	err = Cli("!help ok")
	if err == nil {
		t.Error("!help 命令参数判断出错")
	}
}

//TestCli2 检测!lang命令
func TestCli2(t *testing.T) {
	err := Cli("!lang")
	if err != nil {
		t.Error("!lang 命令出错")
	}
	err = Cli("!lang fra")
	if err != nil {
		t.Error("!lang 命令参数判断出错")
	}
	err = Cli("!lang fra zh")
	if err != nil {
		t.Error("!lang 命令参数判断出错")
	}
	err = Cli("!lang fr")
	if err == nil {
		t.Error("!lang 命令参数判断出错")
	}
	err = Cli("!lang fra dd")
	if err == nil {
		t.Error("!lang 命令参数判断出错")
	}
}

//TestCli3 检测!setPrikey与!setAppID命令
func TestCli3(t *testing.T) {
	err := Cli("!setPrikey")
	if err != nil {
		t.Error("!setPrikey 命令出错")
	}
	err = Cli("!setAppID")
	if err != nil {
		t.Error("!setAppID 命令出错")
	}
}
