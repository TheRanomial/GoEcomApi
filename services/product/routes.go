package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TheRanomial/GoEcomApi/services/auth"
	"github.com/TheRanomial/GoEcomApi/types"
	"github.com/TheRanomial/GoEcomApi/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore,userStore types.UserStore) *Handler{
	return &Handler{store:store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router){

	router.HandleFunc("/products",h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/product/{productId}",h.handleGetProductById).Methods(http.MethodGet)

	//middleware
	router.HandleFunc("/products", auth.WithJWTAuth(h.handleCreateProduct, h.userStore)).Methods("POST")
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {

	products,err:=h.store.GetProducts()

	if err!=nil{
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no product available"))
		return
	}
	utils.WriteJSON(w,http.StatusOK,products)

}

func (h *Handler) handleGetProductById(w http.ResponseWriter, r *http.Request){
	vars:=mux.Vars(r)

	str,ok:=vars["productID"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product id"))
		return
	}

	productId,err:=strconv.Atoi(str)

	if err!=nil{
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product id"))
		return
	}

	user,err:=h.store.GetProductById(productId)

	if err!=nil{
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no product with such id"))
		return
	}
	utils.WriteJSON(w,http.StatusOK,user)
}


func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.CreateProductPayload
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.store.CreateProduct(product)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, product)
}
