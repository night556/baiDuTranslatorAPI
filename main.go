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

func main() {
	var (
		input string
		from  = "auto"
		to    string
		query string
	)
	for {
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		input = reader.Text()
		if input == "exit" {
			return
		}
		inputs := strings.Split(input, "/")
		to = "zh"
		switch len(inputs) {
		case 1:
			for _, v := range inputs[0] {
				if unicode.Is(unicode.Scripts["Han"], v) {
					fmt.Println(v)
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

		query = inputs[0]
		translator(query, from, to)
	}
}

func translator(words string, from string, to string) {
	// var query string
	var result Result

	//URL编码
	urlquery := url.QueryEscape(words)

	prikey := "百度提供的私钥"
	appid := "百度提供的appid"
	//随机数
	randnum := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000)
	salt := strconv.FormatInt(randnum, 10)
	str := appid + words + salt + prikey

	//sign
	var sign string
	h := md5.New()
	h.Write([]byte(str))
	sign = hex.EncodeToString(h.Sum(nil))
	url := "http://api.fanyi.baidu.com/api/trans/vip/translate?q=" + urlquery + "&from=" + from + "&to=" + to + "&appid=百度提供的apiid&salt=" + salt + "&sign=" + sign
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
