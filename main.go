package main

import (
	"http/route"
	"log"
	"net/http"
)

func main() {
	//注册路由函数
	//HandleFunc注册一个处理器函数handler和对应的模式pattern（注册到DefaultServeMux）
	http.HandleFunc("/check", route.CheckParameters)
	http.HandleFunc("/setCookie", route.SetCookieHandler)
	http.HandleFunc("/getCookie", route.GetCookieHandler)
	http.HandleFunc("/upload", route.UploadHandler)
	http.HandleFunc("/nupload", route.NuploadHandler)
	http.HandleFunc("/view", route.ViewHandler)

	//开始监听 处理请求 返回响应
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
