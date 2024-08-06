package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TheRanomial/GoEcomApi/configs"
	"github.com/TheRanomial/GoEcomApi/services/auth"
	"github.com/TheRanomial/GoEcomApi/types"
	"github.com/TheRanomial/GoEcomApi/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router){
	router.HandleFunc("/login",h.handleLogin).Methods("POST")
	router.HandleFunc("/register",h.handleRegister).Methods("POST")

	router.HandleFunc("/users/{userID}", auth.WithJWTAuth(h.handleGETUser, h.store)).Methods(http.MethodGet)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request){
	var user types.LoginUserPayload
	
	if err:=utils.ParseJSON(r,&user);err!=nil{
		utils.WriteError(w,http.StatusBadRequest,err)
	}

	if err:=utils.Validate.Struct(user);err!=nil{
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u,err:=h.store.GetUserByEmail(user.Email)
	if err!=nil{
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no user with this email or password"))
		return
	}

	val:=CheckPasswordHash(user.Password,u.Password)

	if !val{
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("password is incorrect"))
		return
	}

	secret := []byte(configs.Envs.JWTSecret)
	token,err:=auth.CreateJWT(secret,u.ID)

	if err!=nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w,http.StatusOK,map[string]string{"token":token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request){
	var user types.RegisterUserPayload

	if err:=utils.ParseJSON(r,&user);err!=nil{
		utils.WriteError(w,http.StatusBadRequest,err)
		return
	}

	if err:=utils.Validate.Struct(user);err!=nil{
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	//user exists
	_,err:=h.store.GetUserByEmail(user.Email)
	if err==nil{
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already exists"))
		return
	}

	//hash password
	hash,err:=HashPassword(user.Password)

	if err!=nil{
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("can't hash given password"))
		return
	}

	err=h.store.CreateUser(types.User{
		FirstName:user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Password: hash,
	})

	if err!=nil{
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("can't create a new user"))
		return
	}

	utils.WriteJSON(w,http.StatusCreated,user)
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func (h *Handler) handleGETUser(w http.ResponseWriter, r *http.Request){
	vars:=mux.Vars(r)

	str,ok:=vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	userId,err:=strconv.Atoi(str)

	if err!=nil{
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	user,err:=h.store.GetUserByID(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w,http.StatusOK,user)
}