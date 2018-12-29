package main

/*
q 请求翻译query
from 翻译源语言
to 译文语言
appid APPID
salt 随机数
sign 签名
*/
import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//Result 用于接收百度返回的结果
type Result struct {
	From         string        `json:"from"`
	To           string        `json:"to"`
	TransResults []TransResult `json:"trans_result"`
}

//TransResult 存放有原文与译文
type TransResult struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

//ErrorResponse 调用百度出错时返回值
type ErrorResponse struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

var (
	//DefalutSrc 默认源语言
	DefalutSrc = "auto"

	//DefalutDec 默认目标语言
	DefalutDec = "zh"

	//私钥
	prikey = "MeBGpG7KPRLCyFs_GQPG"

	//appid
	appid = "20181028000226330"

	//languageList 此键值对用于存储百度支持的语言
	languageList = map[string]string{
		"auto": "自动",
		"zh":   "中文",
		"en":   "英语",
		"yue":  "粤语",
		"wyw":  "文言文",
		"jp":   "日语",
		"kor":  "韩语",
		"fra":  "法语",
		"spa":  "西班牙语",
		"th":   "泰语",
		"ara":  "阿拉伯语",
		"ru":   "俄语",
		"pt":   "葡萄牙语",
		"de":   "德语",
		"it":   "意大利语",
		"el":   "希腊语",
		"nl":   "荷兰语",
		"pl":   "波兰语",
		"bul":  "保加利亚语",
		"est":  "爱沙尼亚语",
		"dan":  "丹麦语",
		"fin":  "芬兰语",
		"cs":   "捷克语",
		"rom":  "罗马尼亚语",
		"slo":  "斯洛文尼亚语",
		"swe":  "瑞典语",
		"hu":   "匈牙利语",
		"cht":  "繁体中文",
		"vie":  "越南语",
	}
)

func main() {
	fmt.Println("当前翻译环境：" + DefalutSrc + " --> " + DefalutDec)
	reader := bufio.NewScanner(os.Stdin)
	//循环读
	for {
		reader.Scan()
		input := reader.Text()

		//判断是否是命令
		if input[0] == '!' && input[1] != ' ' {
			err := Cli(input)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		//检查输入
		query, from, to := checkInput(input)
		result, err := translator(query, from, to)
		fmt.Println(printResult(result, err))
	}
}

//打印输出
func printResult(result interface{}, err error) string {
	var resultStr string
	//服务器错误
	if err != nil {
		errorResponse := result.(ErrorResponse)
		resultStr = "Error Code:" + errorResponse.ErrorCode + "\tError Msg: " + errorResponse.ErrorMsg + "\n"
	} else {
		//翻译成功
		resultRight := result.(Result)
		for _, v := range resultRight.TransResults {
			resultStr += v.Dst
		}
	}

	return resultStr
}

//checkInput 检查用户输入，
//TODO 用于互译功能
func checkInput(input string) (string, string, string) {
	var (
		to    string
		query string
	)
	from := DefalutSrc
	to = DefalutDec

	for _, v := range input {
		if unicode.Is(unicode.Scripts["Han"], v) && DefalutDec == "zh" {
			to = "en"
		}
		break
	}

	query = input
	return query, from, to
}

//调用此函数之前用户的输入已经经过判断
func translator(words string, from string, to string) (interface{}, error) {
	// var query string
	var result Result
	var errorRespon ErrorResponse

	//URL编码
	urlquery := url.QueryEscape(words)

	//随机数
	randnum := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000)
	salt := strconv.FormatInt(randnum, 10)
	str := appid + words + salt + prikey

	//sign
	var sign string
	h := md5.New()
	h.Write([]byte(str))
	sign = hex.EncodeToString(h.Sum(nil))
	url := "http://api.fanyi.baidu.com/api/trans/vip/translate?q=" + urlquery + "&from=" + from + "&to=" + to + "&appid=" + appid + "&salt=" + salt + "&sign=" + sign
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("链接出错")
	}
	defer resp.Body.Close()

	jsonresult := make([]byte, 1024)
	reader := bufio.NewScanner(resp.Body)
	if reader.Scan() {
		jsonresult = reader.Bytes()
	}
	json.Unmarshal(jsonresult, &result)
	if result.TransResults == nil {
		json.Unmarshal(jsonresult, &errorRespon)
		return errorRespon, errors.New("get response error")
	}
	return result, nil
}

//Cli 处理命令相关的逻辑
func Cli(input string) (err error) {
	commands := strings.Fields(input)
	switch commands[0] {
	//帮助命令
	case "!help":
		if len(commands) != 1 {
			return errors.New("命令错误，可用 !help 显示帮助")
		}
		fmt.Println(`
		!help				显示帮助
		!lang [srclang] [declang]	更改默认[翻译语言设置]
		!exit				退出系统
		!setPrikey			设置私钥
		!setAppID			设置appid
		`)
		//退出系统
	case "!exit":
		os.Exit(0)
		//更改默认翻译语言
	case "!lang":
		switch len(commands) {
		case 1:
			for key, val := range languageList {
				fmt.Println(key + "\t" + val)
			}
			fmt.Println("当前翻译环境：" + DefalutSrc + " --> " + DefalutDec)
		case 2:
			//判断目标语言是否在语言列表中
			if languageList[commands[1]] == "" {
				err := errors.New("目标语言不受支持或不识别")
				return err
			}
			DefalutDec = commands[1]
		case 3:
			if languageList[commands[1]] == "" || languageList[commands[2]] == "" {
				err := errors.New("目标语言或源语言不受支持或不识别")
				return err
			}
			DefalutDec = commands[1]
			DefalutSrc = commands[2]
		}
	case "!setPrikey":
		fmt.Println("当前私钥：" + prikey)
		fmt.Printf("请输入新的私钥：")
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		input := reader.Text()
		prikey = input
	case "!setAppID":
		fmt.Println("当前appid：" + appid)
		fmt.Printf("请输入新的appid：")
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		input := reader.Text()
		appid = input
	}

	return nil
}
