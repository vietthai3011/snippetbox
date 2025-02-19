package main

import (
	"net/http"
)

func (app *application) routers() *http.ServeMux {

	/*
		1. http.StripPrefix("/static", ...) sẽ loại bỏ phần /static khỏi đường dẫn, vì vậy yêu cầu sẽ chuyển thành images/logo.png.
		2. http.FileServer(http.Dir("./ui/static/")) sẽ dùng đường dẫn đã được sửa đổi (images/logo.png) để tìm file trong thư mục ./ui/static/.
		3. Go sẽ trả về file logo.png trong thư mục ./ui/static/images/ cho người dùng.
	*/
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", http.HandlerFunc(app.home))
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	return mux
}
