package route

import (
	"fmt"
	"net/http"
	"strings"
)

func CheckParameters(w http.ResponseWriter, r *http.Request) {
	//request的结构体中有Header Header这一属性
	//Header中放的就是各种请求头的信息
	//通过观察文档可以发现Header的具体结构 通过遍历即可得到我们想要的请求头信息
	//	Header = map[string][]string{
	//		"Accept-Encoding": {"gzip, deflate"},
	//		"Accept-Language": {"en-us"},
	//		"Connection": {"keep-alive"},
	//	}
	for k, v := range r.Header {
		for _, j := range v {
			fmt.Println(k + ":" + j)
		}
	}

	//(把http请求拆开 放到map里)
	//ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
	r.ParseForm()
	// Form是解析好的表单数据，包括URL字段的query参数和POST或PUT的表单数据。
	// 本字段只有在调用ParseForm后才有效。在客户端，会忽略请求中的本字段而使用Body替代。
	//Form的结构与Header类似 map[string][]string
	fmt.Println(r.Form)
	//Path是请求的路径
	fmt.Println("path", r.URL.Path)
	//如果url参数中有url_long 则将会被取出
	fmt.Println(r.Form["url_long"])
	//遍历参数的key和value
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	Start(w, r, SessionConf{
		Sessionname:    "GOSESSION",
		Sessionpath:    nil,
		Maxlifetime:    3600,
		Cookielifetime: 3600,
		Gclifetime:     3600,
		Secure:         true,
	})
	Set("session1", "SESSION1")
	Set("session2", "SESSION2")
	Set("session3", "SESSION3")
	fmt.Println(Get("session1"))
	Del("session1")
	fmt.Println(GetAll())
	Gc()
	//这个写入到w的是输出到客户端的
	fmt.Fprintf(w, "Hello world!")
}
