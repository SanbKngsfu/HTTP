package route

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

const (
	UPLOAD_DIR = "./uploads"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("tmpl/upload.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}
	if r.Method == "POST" {
		//FormFile返回以key为键查询r.MultipartForm字段得到结果中的第一个文件和它的信息。
		//如果必要，本函数会隐式调用ParseMultipartForm和ParseForm。查询失败会返回ErrMissingFile错误。
		f, h, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filename := h.Filename
		defer f.Close()
		//Create采用模式0666（任何人都可读写，不可执行）创建一个名为name的文件，如果文件已存在会截断它（为空文件）。
		//如果成功，返回的文件对象可用于I/O；对应的文件描述符具有O_RDWR模式。如果出错，错误底层类型是*PathError。

		t, err := os.Create(UPLOAD_DIR + "/" + filename)
		fmt.Println("==============")
		fmt.Println(t)
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/view?id="+filename,
			http.StatusFound)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageId
	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image")
	//ServeFile回复请求name指定的文件或者目录的内容。
	http.ServeFile(w, r, imagePath)
}
func isExists(path string) bool {
	//Stat返回一个描述name指定的文件对象的FileInfo。如果指定的文件对象是一个符号链接，
	//返回的FileInfo描述该符号链接指向的文件的信息，本函数会尝试跳转该链接。
	//如果出错，返回的错误值为*PathError类型。
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	//返回一个布尔值说明该错误是否表示一个文件或目录已经存在。ErrExist和一些系统调用错误会使它返回真。
	return os.IsExist(err)
}
