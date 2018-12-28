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

type Result struct {
	From         string        `json:"from"`
	To           string        `json:"to"`
	TransResults []TransResult `json:"trans_result"`
}
type TransResult struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

var languageList = []string{"zh", "en", "yue", "wyw"}
var flag = false

func main() {

	//循环读
	for {
		flag = false
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		input := reader.Text()
		if input == "exit" {
			return
		}
		query, from, to, err := checkInput(input)
		if err != nil {
			fmt.Println("输入有误，无法翻译！！！\n" + err.Error())
			continue
		}
		translator(query, from, to)
	}
}

func checkInput(input string) (string, string, string, error) {
	var (
		from  = "auto"
		to    string
		query string
	)
	inputs := strings.Split(input, "/")

	to = "zh"
	switch len(inputs) {
	case 1:
		for _, v := range inputs[0] {
			if unicode.Is(unicode.Scripts["Han"], v) {
				to = "en"
			}
			break
		}
	case 2:
		to = inputs[1]

	case 3:
		from = inputs[2]
		to = inputs[1]
	}
	for _, dirlang := range languageList {
		if to == dirlang {
			break
		}
		err := errors.New("目标语言错误")
		return "", "", "", err
	}
	query = inputs[0]
	return query, from, to, nil
}

//调用此函数之前用户的输入已经经过判断
func translator(words string, from string, to string) {
	// var query string
	var result Result

	//URL编码
	urlquery := url.QueryEscape(words)

	prikey := "MeBGpG7KPRLCyFs_GQPG"
	appid := "20181028000226330"
	//随机数
	randnum := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000)
	salt := strconv.FormatInt(randnum, 10)
	str := appid + words + salt + prikey

	//sign
	var sign string
	h := md5.New()
	h.Write([]byte(str))
	sign = hex.EncodeToString(h.Sum(nil))
	url := "http://api.fanyi.baidu.com/api/trans/vip/translate?q=" + urlquery + "&from=" + from + "&to=" + to + "&appid=20181028000226330&salt=" + salt + "&sign=" + sign
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	jsonresult := make([]byte, 1024)
	reader := bufio.NewScanner(resp.Body)
	if reader.Scan() {
		jsonresult = reader.Bytes()
	}
	err = json.Unmarshal(jsonresult, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range result.TransResults {
		fmt.Println(v.Dst)
	}

}
