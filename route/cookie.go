package route

import (
	"fmt"
	"net/http"
	"time"
)

func SetCookieHandler(w http.ResponseWriter, r *http.Request) {
	COOKIE_MAX_MAX_AGE := time.Hour * 24 / time.Second // 单位：秒。
	maxAge := int(COOKIE_MAX_MAX_AGE)
	uid := "10"

	uid_cookie := &http.Cookie{
		Name:     "uid",
		Value:    uid,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   maxAge,
	}

	//SetCookie在w的头域中添加Set-Cookie头，该HTTP头的值为cookie。
	http.SetCookie(w, uid_cookie)

	c1 := http.Cookie{
		Name:     "first_cookie",
		Value:    "vanyar",
		HttpOnly: true,
	}

	c2 := http.Cookie{
		Name:     "second_cookie",
		Value:    "noldor",
		HttpOnly: true,
	}
	//Set添加键值对到Header，如键已存在则会用只有新值一个元素的切片取代旧值切片。
	w.Header().Set("Set-Cookie", c1.String()) //String返回该cookie的序列化结果。
	//Add添加键值对到Header，如键已存在则会将新的值附加到旧值切片后面。
	w.Header().Add("Set-Cookie", c2.String())

}

func GetCookieHandler(w http.ResponseWriter, r *http.Request) {
	h := r.Header["Cookie"]
	fmt.Fprintln(w, h)

	c1, err := r.Cookie("first_cookie")
	if err != nil {
		fmt.Fprintln(w, "Cannot get cookie")
	}
	cs := r.Cookies()
	fmt.Fprintln(w, c1)
	fmt.Fprintln(w, cs)

}
