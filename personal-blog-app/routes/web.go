package routes

import (
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/personal-blog-app/handlers"
	"github.com/jaygaha/roadmap-go-projects/personal-blog-app/middleware"
)

var middlewareList = []func(http.HandlerFunc) http.HandlerFunc{
	middleware.Authenticate,
}

func WebRoutes() {
	// app routes
	http.HandleFunc("/", handlers.WelcomeHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	// show article
	http.HandleFunc("/article/", handlers.BlogShowHandler) // /blogs/1234567890 - show blog

	// protected routes
	adminHandler := http.HandlerFunc(handlers.AdminDashboardHandler)
	for _, middlewareFunc := range middlewareList {
		adminHandler = middlewareFunc(adminHandler)
	}
	http.Handle("/admin", adminHandler)

	// logout
	logoutHandler := http.HandlerFunc(handlers.LogoutHandler)
	for _, middlewareFunc := range middlewareList {
		logoutHandler = middlewareFunc(logoutHandler)
	}
	http.Handle("/logout", logoutHandler)

	// blogs
	newBlogHandler := http.HandlerFunc(handlers.BlogNewHandler)
	for _, middlewareFunc := range middlewareList {
		newBlogHandler = middlewareFunc(newBlogHandler)
	}
	http.Handle("/blogs/new", newBlogHandler)

	// Save blog
	saveBlogHandler := http.HandlerFunc(handlers.BlogPostHandler)
	for _, middlewareFunc := range middlewareList {
		saveBlogHandler = middlewareFunc(saveBlogHandler)
	}
	http.Handle("/blogs/submit", saveBlogHandler)

	editBlogHandler := http.HandlerFunc(handlers.BlogEditHandler)
	for _, middlewareFunc := range middlewareList {
		editBlogHandler = middlewareFunc(editBlogHandler)
	}
	http.Handle("/edit/", editBlogHandler)

	deleteBlogHandler := http.HandlerFunc(handlers.BlogDeleteHandler)
	for _, middlewareFunc := range middlewareList {
		deleteBlogHandler = middlewareFunc(deleteBlogHandler)
	}
	http.Handle("/blogs/delete", deleteBlogHandler)
}
