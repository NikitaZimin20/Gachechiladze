package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"wildberries/pkg/product"
	"wildberries/pkg/session"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

type ProductHandler struct {
	Tmpl        *template.Template
	Logger      *zap.SugaredLogger
	Sessions    *session.SessionsManager
	ProductRepo product.ProductRepo
	Cloudinary  *cloudinary.Cloudinary
}

func (h *ProductHandler) Index(w http.ResponseWriter, r *http.Request) {
	orderBy := r.URL.Query().Get("order_by")
	sess, err := session.SessionFromContext(r.Context())
	bsk := &product.Basket{}
	if err == nil {
		h.ProductRepo.AddBasket(sess.UserID)
		bsk, err = h.ProductRepo.GetBasketByID(sess.UserID)
		if err != nil {
			http.Error(w, `DB err`, http.StatusInternalServerError)
			return
		}
		for _, prd := range bsk.Products {
			sess.AddPurchase(prd.ID)
		}
	}

	elems, err := h.ProductRepo.GetAll(orderBy)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Products   []*product.Product
		Session    *session.Session
		TotalCount uint32
	}{
		Products:   elems,
		Session:    sess,
		TotalCount: bsk.TotalCount,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := h.Tmpl.ExecuteTemplate(w, "about.html", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) Privacy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := h.Tmpl.ExecuteTemplate(w, "privacy.html", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) Product(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	bsk := &product.Basket{}
	if err == nil {
		h.ProductRepo.AddBasket(sess.UserID)
		bsk, err = h.ProductRepo.GetBasketByID(sess.UserID)
		if err != nil {
			http.Error(w, `DB err`, http.StatusInternalServerError)
			return
		}
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Bad Id", http.StatusBadGateway)
		return
	}

	prod, err := h.ProductRepo.GetByID(uint32(id))
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	prod.Views++

	_, err = h.ProductRepo.Update(prod)
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	prds, err := h.ProductRepo.GetRelated(prod.Type, 5)
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	fmt.Print("ebal", prds)

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "product.html", struct {
		Product    *product.Product
		Related    []*product.Product
		Session    *session.Session
		TotalCount uint32
	}{
		Product:    prod,
		Related:    prds,
		Session:    sess,
		TotalCount: bsk.TotalCount,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) AddProductToBasket(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "DB Error", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Bad Id", http.StatusBadGateway)
		return
	}

	prod, err := h.ProductRepo.GetByID(uint32(id))
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	prodId, err := h.ProductRepo.AddProductToBasket(sess.UserID, prod)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	sess.AddPurchase(prodId)
	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]uint32{
		"updated": prodId,
	})
	w.Write(respJSON)
}

func (h *ProductHandler) DeleteProductFromBasket(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "sess Error", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Bad Id", http.StatusBadGateway)
		return
	}

	_, err = h.ProductRepo.DeleteProductFromBasket(sess.UserID, uint32(id))
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	sess.DeletePurchase(uint32(id))
	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]uint32{
		"updated": uint32(id),
	})
	w.Write(respJSON)
}

func (h *ProductHandler) History(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "sess Error", http.StatusBadRequest)
		return
	}

	landInf, err := h.ProductRepo.GetOrders(sess.UserID)
	if err != nil {
		http.Error(w, "bd Error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "history.html", struct {
		Landings []*product.LandingInfo
	}{
		Landings: landInf,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) Basket(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	bsk := &product.Basket{}
	if err == nil {
		h.ProductRepo.AddBasket(sess.UserID)
		bsk, err = h.ProductRepo.GetBasketByID(sess.UserID)
		if err != nil {
			http.Error(w, `DB err`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "basket.html", bsk)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "sess Error", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Bad Id", http.StatusBadGateway)
		return
	}

	_, err = h.ProductRepo.Delete(uint32(id))
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	sess.DeletePurchase(uint32(id))
	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]uint32{
		"updated": uint32(id),
	})
	w.Write(respJSON)
}

func (h *ProductHandler) AddProductForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := h.Tmpl.ExecuteTemplate(w, "createproduct.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `Bad id`, http.StatusBadRequest)
		return
	}

	r.ParseForm()
	product := new(product.Product)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(product, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}
	product.ID = uint32(id)

	ok, err := h.ProductRepo.Update(product)
	if err != nil {
		http.Error(w, `db error`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("update: %v %v", product, ok)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *ProductHandler) RegisterOrder(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "sess Error", http.StatusBadRequest)
		return
	}
	bsk, err := h.ProductRepo.GetBasketByID(sess.UserID)
	if err != nil {
		http.Error(w, "DB Error", http.StatusBadRequest)
		return
	}
	id, err := h.ProductRepo.RegisterOrder(sess.UserID, bsk.Products)
	if err != nil {
		http.Error(w, "DB Error", http.StatusBadRequest)
		return
	}
	fmt.Print("kekw", id)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	product := new(product.Product)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(product, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}
	fmt.Print("??????????????", product)

	// Get handler for filename, size and headers
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	resp, err := h.Cloudinary.Upload.Upload(context.TODO(), file, uploader.UploadParams{})
	if err != nil {
		http.Error(w, `upload err`, http.StatusInternalServerError)
		return
	}

	product.ImageURL = resp.URL

	defer file.Close()

	lastID, err := h.ProductRepo.Add(product)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Insert with id LastInsertId: %v", lastID)
	http.Redirect(w, r, "/", http.StatusFound)
}
