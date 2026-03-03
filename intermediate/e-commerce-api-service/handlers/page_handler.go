// This handler serves HTML pages using Go's html/template.
// Each page method prepares data, then renders through the shared layout.
//
// Key pattern: every page handler builds a PageData struct, which the
// layout template uses for the nav, and the child template uses for content.

package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/middleware"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/services"
)

// PageData is the shared struct every template receives.
// The layout reads User/CartCount/Flash; each page reads its own fields.
type PageData struct {
	User      *models.User
	CartCount int
	Flash     string
	FlashType string // "success" or "error"

	// Page-specific data — only one is populated per request
	Products []models.Product
	Product  *models.Product
	Cart     *models.CartResponse
	Orders   []models.Order
	Order    *models.Order
	Checkout *models.CheckoutResponse
}

type PageHandler struct {
	templates  map[string]*template.Template
	productSvc *services.ProductService
	cartSvc    *services.CartService
	orderSvc   *services.OrderService
	authSvc    *services.AuthService
}

// FuncMap provides helper functions available in every template.
var tmplFuncs = template.FuncMap{
	"formatPrice": func(cents int64) string {
		return "$" + strconv.FormatFloat(float64(cents)/100.0, 'f', 2, 64)
	},
	"multiply": func(a int64, b int) int64 {
		return a * int64(b)
	},
	"productGradient": func(id int64) string {
		gradients := []string{
			"linear-gradient(135deg, #f5f0e8, #e8dcc8)",
			"linear-gradient(135deg, #e8efe5, #c8d8c0)",
			"linear-gradient(135deg, #f0e8f5, #d8c8e8)",
			"linear-gradient(135deg, #e8f0f5, #c8d8e8)",
			"linear-gradient(135deg, #f5ede8, #e8d4c8)",
			"linear-gradient(135deg, #e8f5f0, #c8e8d8)",
			"linear-gradient(135deg, #f5e8ea, #e8c8cc)",
			"linear-gradient(135deg, #eef0e8, #dae0c8)",
		}
		return gradients[(id-1)%int64(len(gradients))]
	},
	"productIcon": func(id int64) string {
		icons := []string{"◎", "⏀", "◈", "▣", "◉", "⬡", "◇", "△"}
		return icons[(id-1)%int64(len(icons))]
	},
	"stockBadge": func(stock int) string {
		if stock == 0 {
			return "out-of-stock"
		}
		if stock <= 5 {
			return "low-stock"
		}
		return "in-stock"
	},
	"seq": func(n int) []int {
		s := make([]int, n)
		for i := range s {
			s[i] = i + 1
		}
		return s
	},
}

func NewPageHandler(
	tmplDir string,
	productSvc *services.ProductService,
	cartSvc *services.CartService,
	orderSvc *services.OrderService,
	authSvc *services.AuthService,
) *PageHandler {
	h := &PageHandler{
		templates:  make(map[string]*template.Template),
		productSvc: productSvc,
		cartSvc:    cartSvc,
		orderSvc:   orderSvc,
		authSvc:    authSvc,
	}

	// Parse each page template together with the shared layout.
	// This gives every page access to the layout's nav, CSS, and scripts.
	pages := []string{
		"home", "product", "cart", "checkout",
		"orders", "auth", "admin",
	}
	layoutFile := filepath.Join(tmplDir, "layout.html")
	for _, page := range pages {
		pageFile := filepath.Join(tmplDir, page+".html")
		t := template.Must(
			template.New("layout.html").Funcs(tmplFuncs).ParseFiles(layoutFile, pageFile),
		)
		h.templates[page] = t
	}

	return h
}

// render executes the named template with the given data.
func (h *PageHandler) render(w http.ResponseWriter, name string, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates[name].ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Printf("[TEMPLATE] Error rendering %s: %v", name, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// pageData builds the shared PageData with nav info from the request context.
func (h *PageHandler) pageData(r *http.Request) PageData {
	pd := PageData{}

	// Try to extract user from context (set by auth middleware)
	if uid, ok := r.Context().Value(middleware.UserIdKey).(int64); ok {
		role, _ := r.Context().Value(middleware.RoleKey).(string)
		pd.User = &models.User{ID: uid, Role: role}

		// Load cart count for the nav badge
		if cart, err := h.cartSvc.GetCart(uid); err == nil && cart != nil {
			pd.CartCount = len(cart.Items)
		}
	}

	// Flash messages via query params (set by redirect after POST)
	pd.Flash = r.URL.Query().Get("flash")
	pd.FlashType = r.URL.Query().Get("flash_type")
	if pd.FlashType == "" {
		pd.FlashType = "success"
	}

	return pd
}

// ── Page Handlers

// POST /auth/cookie
// Bridges API auth to page sessions: takes Authorization: Bearer <token>
// validates it, and sets the HttpOnly cookie so server-rendered pages work.
func (h *PageHandler) SetAuthCookie(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		http.Error(w, `{"error":"missing or invalid authorization header"}`, http.StatusUnauthorized)
		return
	}
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	if _, err := h.authSvc.ValidateToken(tokenStr); err != nil {
		http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   259200, // 72 hours
	})
	w.WriteHeader(http.StatusNoContent)
}

// GET /
func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	q := models.ProductQuery{
		Name:  r.URL.Query().Get("q"),
		Page:  1,
		Limit: 20,
	}
	if p, _ := strconv.Atoi(r.URL.Query().Get("page")); p > 0 {
		q.Page = p
	}
	products, err := h.productSvc.List(q)
	if err != nil {
		log.Printf("[PAGE] Error loading products: %v", err)
	}
	pd.Products = products
	h.render(w, "home", pd)
}

// GET /products/{id}
func (h *PageHandler) ProductDetail(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Redirect(w, r, "/?flash=Product+not+found&flash_type=error", http.StatusSeeOther)
		return
	}
	product, err := h.productSvc.GetById(id)
	if err != nil {
		http.Redirect(w, r, "/?flash=Product+not+found&flash_type=error", http.StatusSeeOther)
		return
	}
	pd.Product = product
	h.render(w, "product", pd)
}

// GET /cart
func (h *PageHandler) Cart(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	h.render(w, "cart", pd)
}

// POST /cart/add
func (h *PageHandler) CartAdd(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	if pd.User == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	productId, _ := strconv.ParseInt(r.FormValue("product_id"), 10, 64)
	qty, _ := strconv.Atoi(r.FormValue("quantity"))
	if qty <= 0 {
		qty = 1
	}

	err := h.cartSvc.AddItem(pd.User.ID, models.AddToCartRequest{
		ProductId: productId,
		Quantity:  qty,
	})
	if err != nil {
		http.Redirect(w, r, "/cart?flash="+err.Error()+"&flash_type=error", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/cart?flash=Item+added+to+cart", http.StatusSeeOther)
}

// POST /cart/update
func (h *PageHandler) CartUpdate(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	if pd.User == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	productID, _ := strconv.ParseInt(r.FormValue("product_id"), 10, 64)
	qty, _ := strconv.Atoi(r.FormValue("quantity"))

	h.cartSvc.UpdateItem(pd.User.ID, productID, models.UpdateCartItemRequest{Quantity: qty})
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

// POST /cart/remove
func (h *PageHandler) CartRemove(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	if pd.User == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	productId, _ := strconv.ParseInt(r.FormValue("product_id"), 10, 64)
	h.cartSvc.RemoveItem(pd.User.ID, productId)
	http.Redirect(w, r, "/cart?flash=Item+removed", http.StatusSeeOther)
}

// POST /checkout
func (h *PageHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	if pd.User == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}
	resp, err := h.orderSvc.Checkout(pd.User.ID)
	if err != nil {
		http.Redirect(w, r, "/cart?flash="+err.Error()+"&flash_type=error", http.StatusSeeOther)
		return
	}
	pd.Checkout = resp
	h.render(w, "checkout", pd)
}

// GET /orders
func (h *PageHandler) Orders(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	h.render(w, "orders", pd)
}

// GET /auth
func (h *PageHandler) AuthPage(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	h.render(w, "auth", pd)
}

// POST /auth/login
func (h *PageHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	resp, err := h.authSvc.Login(models.UserLoginRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})
	log.Printf("[PAGE] Login response: %v", resp)
	if err != nil {
		http.Redirect(w, r, "/auth?flash=Invalid+credentials&flash_type=error", http.StatusSeeOther)
		return
	}
	// Set JWT as a cookie for server-side rendered pages
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   259200, // 72 hours (matches JWT expiry)
	})
	http.Redirect(w, r, "/?flash=Welcome+back!", http.StatusSeeOther)
}

// POST /auth/register
func (h *PageHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	resp, err := h.authSvc.Register(models.UserRegisterRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		http.Redirect(w, r, "/auth?flash="+err.Error()+"&flash_type=error", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   259200,
	})
	http.Redirect(w, r, "/?flash=Account+created!", http.StatusSeeOther)
}

// GET /logout
func (h *PageHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/?flash=Signed+out", http.StatusSeeOther)
}

// GET /admin
func (h *PageHandler) AdminPage(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	products, _ := h.productSvc.List(models.ProductQuery{Limit: 100})
	pd.Products = products
	h.render(w, "admin", pd)
}

// POST /admin/products/create
func (h *PageHandler) AdminCreateProduct(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	if pd.User == nil || pd.User.Role != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	_, err := h.productSvc.Create(models.ProductCreateRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       int64(price * 100), // dollars to cents
		Stock:       stock,
		ImageURL:    r.FormValue("image_url"),
	})
	if err != nil {
		http.Redirect(w, r, "/admin?flash="+err.Error()+"&flash_type=error", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/admin?flash=Product+created", http.StatusSeeOther)
}

// POST /admin/products/delete
func (h *PageHandler) AdminDeleteProduct(w http.ResponseWriter, r *http.Request) {
	pd := h.pageData(r)
	if pd.User == nil || pd.User.Role != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	id, _ := strconv.ParseInt(r.FormValue("product_id"), 10, 64)
	h.productSvc.Delete(id)
	http.Redirect(w, r, "/admin?flash=Product+deleted", http.StatusSeeOther)
}
