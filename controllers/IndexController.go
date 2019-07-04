// Author: James Mallon <jamesmallondev@gmail.com>
// controllers package -
package controllers

import (
	. "GoAuthorization/libs/layout"
	"net/http"
)

// Struct type indexController -
type indexController struct {
	LayoutHelper
}

// IndexController function -
func IndexController() *indexController {
	return &indexController{}
}

// Index method -
func (this *indexController) Index(w http.ResponseWriter, r *http.Request) {
	this.Render(w,
		map[string]interface{}{"PageTitle": "Index"},
		"layout.gohtml", "index/index.gohtml")
}

// About method -
func (this *indexController) About(w http.ResponseWriter, r *http.Request) {
	this.Render(w,
		map[string]interface{}{"PageTitle": "Index"},
		"layout.gohtml", "index/about.gohtml")
}

// ContactUs method -
func (this *indexController) ContactUs(w http.ResponseWriter, r *http.Request) {
	this.Render(w,
		map[string]interface{}{"PageTitle": "Index"},
		"layout.gohtml", "index/contact.gohtml")
}
