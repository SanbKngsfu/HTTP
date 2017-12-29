package route

import (
	"html/template"
	"io"
	"net/http"
	"os"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles("tmpl/uploadn.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	//ExecuteTemplate方法类似Execute，但是使用名为name的t关联的模板产生输出。
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

func NuploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET displays the upload form.
	case "GET":
		display(w, "uploadn", nil)

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		//parse the multipart form in the request
		//ParseMultipartForm将请求的主体作为multipart/form-data解析。
		//请求的整个主体都会被解析，得到的文件记录最多maxMemery字节保存在内存，
		//其余部分保存在硬盘的temp文件里。
		err := r.ParseMultipartForm(100000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get a ref to the parsed multipart form
		m := r.MultipartForm

		//get the *fileheaders
		files := m.File["myfiles"]
		for i, _ := range files {
			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//create destination file making sure the path is writeable.
			dst, err := os.Create("./uploads/" + files[i].Filename)
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//copy the uploaded file to the destination file
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}
		//display success message.
		display(w, "uploadn", "Upload successful.")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
