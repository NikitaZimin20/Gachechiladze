package main

import (
	"html/template"
	"net/http"

	"wildberries/middleware"
	"wildberries/pkg/handlers"
	"wildberries/pkg/product"
	"wildberries/pkg/session"
	"wildberries/pkg/user"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	templates := template.Must(template.ParseGlob("./static/html/*"))
	userRepo := user.NewMemoryRepo()
	productRepo := product.NewMemoryRepo()
	sm := session.NewSessionsManager()
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()
	cld, _ := cloudinary.NewFromParams("davauqkbe", "637661274245857", "xf9AYi4mcK2yLIfi3ERiTFVqeq4")

	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		Sessions: sm,
		UserRepo: userRepo,
		Logger:   logger,
	}

	productHandler := &handlers.ProductHandler{
		Tmpl:        templates,
		Sessions:    sm,
		ProductRepo: productRepo,
		Logger:      logger,
		Cloudinary:  cld,
	}

	staticHandler := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("./static")),
	)

	r := mux.NewRouter()

	r.HandleFunc("/", productHandler.Index).Methods("GET")
	r.HandleFunc("/about", productHandler.About).Methods("GET")
	r.HandleFunc("/privacy", productHandler.Privacy).Methods("GET")
	r.HandleFunc("/history", productHandler.History).Methods("GET")

	r.HandleFunc("/products/new", productHandler.AddProductForm).Methods("GET")
	r.HandleFunc("/products/new", productHandler.AddProduct).Methods("POST")
	r.HandleFunc("/products/{id}", productHandler.Product).Methods("PUT")
	r.HandleFunc("/products/{id}", productHandler.Product).Methods("GET")
	r.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	r.HandleFunc("/basket/{id}", productHandler.AddProductToBasket).Methods("GET")
	r.HandleFunc("/basket/{id}", productHandler.DeleteProductFromBasket).Methods("DELETE")
	r.HandleFunc("/basket", productHandler.Basket).Methods("GET")
	r.HandleFunc("/register_order", productHandler.RegisterOrder).Methods("GET")

	r.HandleFunc("/register", userHandler.Register).Methods("GET")
	r.HandleFunc("/login", userHandler.Login).Methods("GET")
	r.HandleFunc("/logout", userHandler.Logout).Methods("GET")

	r.HandleFunc("/sign_up", userHandler.SignUp).Methods("POST")
	r.HandleFunc("/sign_in", userHandler.SignIn).Methods("POST")

	r.PathPrefix("/static/").Handler(staticHandler)

	mux := middleware.Auth(sm, r)
	mux = middleware.AccessLog(logger, mux)
	mux = middleware.Panic(mux)

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	http.ListenAndServe(addr, mux)
}
