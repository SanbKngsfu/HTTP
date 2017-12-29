package route

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var sessionObj *session

func init() {
	sessionObj = &session{}
}

type Guid struct {
}

//md5加密
func (this *Guid) md5(v string) string {
	md5Obj := md5.New()
	md5Obj.Write([]byte(v))
	char := md5Obj.Sum(nil)
	return hex.EncodeToString(char)
}

//生成唯一标识
func (this *Guid) Guid() string {
	c := make([]byte, 32)
	//ReadFull从r精确地读取len(buf)字节数据填充进c
	if _, err := io.ReadFull(rand.Reader, c); err != nil {
		return ""
	}
	return this.md5(base64.URLEncoding.EncodeToString(c))
}

type session struct {
	SessionId   string
	SessionName string
	SessionPath interface{}
	Maxlifetime int64
	Gclifetime  int64
}

type SessionConf struct {
	Sessionpath    interface{}
	Sessionname    string
	Cookielifetime int64
	Maxlifetime    int64
	Gclifetime     int64
	Secure         bool
}

//开启会话(Session)
func (this *session) Start(w http.ResponseWriter, r *http.Request, SessionConf SessionConf) {
	obj := &Guid{}
	this.SessionName = SessionConf.Sessionname
	if SessionConf.Maxlifetime == 0 {
		this.Maxlifetime = 3600
	} else {
		this.Maxlifetime = SessionConf.Maxlifetime
	}
	if SessionConf.Gclifetime == 0 {
		this.Gclifetime = 3600
	} else {
		this.Gclifetime = SessionConf.Gclifetime
	}
	//未设置Session存储路径处理
	if SessionConf.Sessionpath == nil {
		//Getwd返回一个对应当前工作目录的根路径。
		path, _ := os.Getwd()
		path += "\\session\\temp\\"
		//Stat返回描述文件f的FileInfo类型值
		if _, err := os.Stat(path); err != nil {
			//ModePerm FileMode = 0777 覆盖所有Unix权限位（用于通过&获取类型位）
			os.Mkdir(path, os.ModePerm)
			this.SessionPath = path
		} else {
			this.SessionPath = path
		}
	}
	//设置Cookie Cookie返回请求中名为name的cookie
	cookie, err := r.Cookie(this.SessionName)
	if err != nil {
		//按条件设置Cookie
		this.SessionId = obj.Guid()
		cookieConf := &http.Cookie{
			Name:     this.SessionName,
			Value:    this.SessionId,
			HttpOnly: true,
			MaxAge:   0}
		if SessionConf.Cookielifetime != 0 {
			cookieConf.Expires = time.Unix(time.Now().Unix()+SessionConf.Cookielifetime, 0)
		}
		if SessionConf.Secure == true {
			cookieConf.Secure = true
		}
		http.SetCookie(w, cookieConf)
	} else {
		this.SessionId = cookie.Value
	}
}

//获取SessionID
func (this *session) GetSessionId() string {
	return this.SessionId
}

//获取Session名称
func (this *session) GetSessionName() string {
	return this.SessionName
}

//设置Session
func (this *session) Set(name, value interface{}) error {
	SessionFile := this.SessionPath.(string) + "sess_" + this.SessionId + ".txt"
	MapValue := make(map[interface{}]interface{})
	//判断当前Session文件是否存在,不存在则创建
	if _, err := os.Stat(SessionFile); err != nil {
		//OpenFile是一个更一般性的文件打开函数，大多数调用者都应用Open或Create代替本函数。
		//它会使用指定的选项（如O_RDONLY等）、指定的模式（如0666等）打开指定名称的文件。
		//如果操作成功，返回的文件对象可用于I/O
		file, _ := os.OpenFile(SessionFile, os.O_CREATE, 0777)
		file.Close()
	}
	//O_RDONLY: open the file read-only.
	f, _ := os.OpenFile(SessionFile, os.O_RDONLY, 0777)
	//函数返回一个从f读取数据的*Decoder，如果f不满足io.ByteReader接口，则会包装f为bufio.Reader。
	d := gob.NewDecoder(f)
	//Decode从输入流读取下一个值并将该值存入e
	//如果e是nil，将丢弃该值；否则e必须是可接收该值的类型的指针
	d.Decode(&MapValue)
	defer func() {
		os.Chtimes(SessionFile, time.Now(), time.Now())
	}()
	f.Close()
	fop, _ := os.OpenFile(SessionFile, os.O_WRONLY|os.O_TRUNC, 0777)
	MapValue[name] = value
	//进行Gob编码
	//NewEncoder返回一个将编码后数据写入w的*Encoder。
	e := gob.NewEncoder(fop)
	//写入到当前Session文件
	//Encode方法将e编码后发送，并且会保证所有的类型信息都先发送。
	e.Encode(MapValue)
	return nil
}

//获取单一Session
func (this *session) Get(name interface{}) (interface{}, error) {
	SessionFile := this.SessionPath.(string) + "sess_" + this.SessionId + ".txt"
	MapValue := make(map[interface{}]interface{})
	//判断当前Session文件是否存在,不存在则创建
	if _, err := os.Stat(SessionFile); err != nil {
		file, _ := os.OpenFile(SessionFile, os.O_CREATE, 0777)
		file.Close()
	}
	//打开当前Session文件
	f, _ := os.OpenFile(SessionFile, os.O_RDONLY, 0777)
	//进行Gob解码
	d := gob.NewDecoder(f)
	d.Decode(&MapValue)
	defer func() {
		os.Chtimes(SessionFile, time.Now(), time.Now())
	}()
	defer f.Close()
	//判断需要获取的Session是否存在
	if MapValue[name] != nil {
		fmt.Println("======================")
		fmt.Println(name.(string) + ":" + MapValue[name].(string))
		fmt.Println("======================")
		return MapValue[name], nil
	} else {
		return "", errors.New("会话" + name.(string) + "不存在")
	}
}

//获取全部Session
func (this *session) GetAll() (map[interface{}]interface{}, error) {
	SessionFile := this.SessionPath.(string) + "sess_" + this.SessionId + ".txt"
	MapValue := make(map[interface{}]interface{})
	//判断当前Session文件是否存在,不存在则创建
	if _, err := os.Stat(SessionFile); err != nil {
		file, _ := os.OpenFile(SessionFile, os.O_CREATE, 0777)
		file.Close()
	}
	//打开当前Session文件
	f, _ := os.OpenFile(SessionFile, os.O_RDONLY, 0777)
	//进行Gob解码
	d := gob.NewDecoder(f)
	d.Decode(&MapValue)
	defer func() {
		os.Chtimes(SessionFile, time.Now(), time.Now())
	}()
	defer f.Close()
	fmt.Println("----------------------")
	fmt.Println(MapValue)
	fmt.Println("----------------------")
	return MapValue, nil
}

//删除Session
func (this *session) Del(name interface{}) error {
	SessionFile := this.SessionPath.(string) + "sess_" + this.SessionId + ".txt"
	MapValue := make(map[interface{}]interface{})
	//判断当前Session文件是否存在,不存在则创建
	if _, err := os.Stat(SessionFile); err != nil {
		file, _ := os.OpenFile(SessionFile, os.O_CREATE, 0777)
		file.Close()
	}
	//打开当前Session文件
	f, _ := os.OpenFile(SessionFile, os.O_RDONLY, 0777)
	d := gob.NewDecoder(f)
	//进行Gob解码
	d.Decode(&MapValue)
	f.Close()
	defer func() {
		os.Chtimes(SessionFile, time.Now(), time.Now())
	}()
	//判断需要获取的Session是否存在,存在则删除
	if MapValue[name] != nil {
		delete(MapValue, name)
		f, _ := os.OpenFile(SessionFile, os.O_WRONLY, 0777)
		e := gob.NewEncoder(f)
		e.Encode(MapValue)
		return nil
	} else {
		return errors.New("会话" + name.(string) + "不存在")
	}
}

//销毁Session
func (this *session) Destroy() error {
	SessionFile := this.SessionPath.(string) + "sess_" + this.SessionId + ".txt"
	if _, err := os.Stat(SessionFile); err != nil {
		return errors.New("Session文件不存在!")
	} else {
		//删除全部Session
		os.Remove(SessionFile)
		return nil
	}
}

//定时垃圾回收
func (this *session) Gc() {
	//设置一个断续器,每隔this.Gclifetime秒后，进行垃圾回收
	t1 := time.NewTicker(time.Duration(this.Gclifetime) * time.Second)
	//开启一个协程处理
	go func() {
		for _ = range t1.C {
			//删除所有this.Gclifetime的过期Session文件
			filepath.Walk(this.SessionPath.(string), func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					if info.ModTime().Unix() < (time.Now().Unix() - this.Maxlifetime) {
						os.Remove(path)
					}
				}
				return nil
			})
		}
	}()
}

//开启会话(Session)
func Start(w http.ResponseWriter, r *http.Request, SessionConf SessionConf) {
	sessionObj.Start(w, r, SessionConf)
}

//获取SessionID
func GetSessionID() string {
	return sessionObj.GetSessionId()
}

//获取Session名称
func GetSessionName() string {
	return sessionObj.GetSessionName()
}

//设置Session
func Set(name, value interface{}) error {
	return sessionObj.Set(name, value)
}

//获取单一session
func Get(name interface{}) (interface{}, error) {
	return sessionObj.Get(name)
}

//删除Session
func Del(name interface{}) error {
	return sessionObj.Del(name)
}

//销毁Session
func Destroy() error {
	return sessionObj.Destroy()
}

//获取全部Session
func GetAll() (map[interface{}]interface{}, error) {
	return sessionObj.GetAll()
}

//垃圾回收
func Gc() {
	sessionObj.Gc()
}
