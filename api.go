package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenAddress string
	Store         Storage
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apifunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string
}

func makeHttpHandler(f apifunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJson(w, 400, err.Error())
		}
	}
}

func newAPIServer(listen string, store Storage) *APIServer {
	return &APIServer{
		ListenAddress: listen,
		Store:         store,
	}
}

func (s *APIServer) run() {
	r := mux.NewRouter()

	r.HandleFunc("/account", makeHttpHandler(s.handleAccount))
	r.HandleFunc("/account/{id}", makeHttpHandler(s.handleAccountByID))
	r.HandleFunc("/transfer", makeHttpHandler(s.handleTransfer))
	r.HandleFunc("/login", makeHttpHandler(s.handleLogin))
	http.ListenAndServe(s.ListenAddress, r)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	loginRequest := LoginRequest{}

	err := json.NewDecoder(r.Body).Decode(loginRequest)
	if err != nil {
		return fmt.Errorf("wrong")
	}

	accInt, err := strconv.Atoi(loginRequest.Username)
	if err != nil {
		return fmt.Errorf("ow")
	}
	acc, err := s.Store.GetAccount(accInt)

	token, err := generateJWT()
	if err != nil {
		return fmt.Errorf("oops again")
	}

	resp := LoginResponse{
		Token:  token,
		UserID: string(acc.Number),
	}
	return writeJson(w, http.StatusOK, resp)
}

func (s *APIServer) handleAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		strId := mux.Vars(r)["id"]
		id, err := strconv.Atoi(strId)
		if err != nil {
			return fmt.Errorf("invalid id %s", strId)
		}
		account, err := s.Store.GetAccount(id)
		if err != nil {
			return err
		}
		return writeJson(w, http.StatusOK, account)
	}
	return fmt.Errorf("invalid method type %s", r.Method)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not supported %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.Store.GetAllAccounts()
	if err != nil {
		return fmt.Errorf("nothing found %+v", err)
	}
	return writeJson(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accountapi := &AccountAPI{}
	err := json.NewDecoder(r.Body).Decode(accountapi)
	if err != nil {
		return err
	}
	newAccount := makeAccount(accountapi.Firstname, accountapi.Lastname)
	err = s.Store.CreateAccount(newAccount)
	if err != nil {
		return err
	}
	writeJson(w, http.StatusOK, accountapi)
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		transferData := &TransferRequest{}
		err := json.NewDecoder(r.Body).Decode(transferData)
		if err != nil {
			return err
		}
		return writeJson(w, http.StatusOK, transferData)
	}
	return nil
}

var secret = []byte("somesecretherer")

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["auth"] = true
	claims["user"] = "just"

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	// secret := os.Getenv("JWT_SECRET")
	secret := "JWT_SECRET"

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}
