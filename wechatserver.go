package main

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	TOKEN    = "dreamcastleshanghai"
	Text     = "text"
	Location = "location"
	Image    = "image"
	Link     = "link"
	Event    = "event"
	Music    = "music"
	News     = "news"
)

type msgBase struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
}

type Request struct {
	XMLName                xml.Name `xml:"xml"`
	msgBase                         // base struct
	Location_X, Location_Y float32
	Scale                  int
	Label                  string
	PicUrl                 string
	MsgId                  int
}

type Response struct {
	XMLName xml.Name `xml:"xml"`
	msgBase
	ArticleCount int     `xml:",omitempty"`
	Articles     []*item `xml:"Articles>item,omitempty"`
	FuncFlag     int
}

type item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	PicUrl      string
	Url         string
}

func weixinEvent(w http.ResponseWriter, r *http.Request) {
	if weixinCheckSignature(w, r) == false {
		fmt.Println("auth failed, attached?")
		return
	}

	fmt.Println("auth success, parse POST")

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(string(body))
	var wreq *Request
	if wreq, err = DecodeRequest(body); err != nil {
		log.Fatal(err)
		return
	}

	wresp, err := dealwith(wreq)
	if err != nil {
		log.Fatal(err)
		return
	}

	data, err := wresp.Encode()
	if err != nil {
		fmt.Printf("error:%v\n", err)
		return
	}

	fmt.Println(string(data))
	fmt.Fprintf(w, string(data))
	return
}

func dealwith(req *Request) (resp *Response, err error) {
	resp = NewResponse()
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = Text

	if req.MsgType == Event {
		if req.Content == "subscribe" {
			resp.Content = "欢迎关注微信公众号 dreamcastleshanghai, 上海梦堡信息科技有限公司。"
			return resp, nil
		}
	}

	if req.MsgType == Text {
		if strings.Trim(strings.ToLower(req.Content), " ") == "help" {
			resp.Content = "欢迎关注微信公众号 dreamcastleshanghai, 上海梦堡信息科技有限公司。"
			return resp, nil
		}
		resp.Content = "将尽快回复您."
	} else if req.MsgType == Image {

		var a item
		a.Description = "description"
		a.Title = "title"
		a.PicUrl = "http://news.07073.com/uploads/090929/9_132710_1_lit.jpg"
		a.Url = "http://news.07073.com/pingce/320407.html"

		resp.MsgType = News
		resp.ArticleCount = 1
		resp.Articles = append(resp.Articles, &a)
		resp.FuncFlag = 1
	} else {
		resp.Content = "暂时还不支持其他的类型"
	}
	return resp, nil
}

func weixinAuth(w http.ResponseWriter, r *http.Request) {
	if weixinCheckSignature(w, r) == true {
		fmt.Println("auth check true.")
		var echostr string = strings.Join(r.Form["echostr"], "")
		fmt.Fprintf(w, echostr)
	}
	fmt.Println("auth check false.")
}

func weixinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("GET begin...")
		weixinAuth(w, r)
		fmt.Println("GET END...")
	} else {
		fmt.Println("POST begin...")
		weixinEvent(w, r)
		fmt.Println("POST END...")
	}
}

func main() {
	http.HandleFunc("/check", weixinHandler)
	//http.HandleFunc("/", action)
	port := "80"
	println("Listening on port ", port, "...")
	err := http.ListenAndServe(":"+port, nil) //设置监听的端口

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func str2sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func weixinCheckSignature(w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	fmt.Println(r.Form)

	var signature string = strings.Join(r.Form["signature"], "")
	var timestamp string = strings.Join(r.Form["timestamp"], "")
	var nonce string = strings.Join(r.Form["nonce"], "")

	fmt.Println("signature : ", signature)
	fmt.Println("timestamp : ", timestamp)
	fmt.Println("nonce : ", nonce)

	tmps := []string{TOKEN, timestamp, nonce}

	fmt.Println(tmps)
	sort.Strings(tmps)
	fmt.Println(tmps)

	tmpStr := tmps[0] + tmps[1] + tmps[2]
	fmt.Println("tmpStr : ", tmpStr)

	tmp := str2sha1(tmpStr)
	fmt.Println("tmp : ", tmp)

	if tmp == signature {
		return true
	}
	return false
}

func DecodeRequest(data []byte) (req *Request, err error) {
	req = &Request{}
	if err = xml.Unmarshal(data, req); err != nil {
		return
	}
	req.CreateTime *= time.Second
	return
}

func NewResponse() (resp *Response) {
	resp = &Response{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	return
}

func (resp Response) Encode() (data []byte, err error) {
	resp.CreateTime = time.Duration(time.Now().Unix())
	data, err = xml.Marshal(resp)
	return
}
