package handlerfunction

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/iharshyadav/backend/internal/storage"
	"github.com/iharshyadav/backend/internal/types"
	"github.com/iharshyadav/backend/internal/utils/response"
)

func Create(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var createUser types.CreateUser

		err := json.NewDecoder(r.Body).Decode(&createUser)

		if errors.Is(err , io.EOF){
            response.WriteJson(w,http.StatusBadRequest,response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		slog.Info("creating a user")

		if err != nil {
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(createUser); err != nil {
			validationeErr := err.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidationError(validationeErr))
			return
		}

		lastId , err := storage.CreateUserInterface(
			createUser.Name,
			createUser.Email,
			createUser.Age,
		)

		slog.Info("user created successfully",slog.String("userId",fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w,http.StatusInternalServerError,err)
			return
		}

		response.WriteJson(w,http.StatusCreated,map[string]int64{"id":lastId})
	}
}


func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.PathValue("id")
		slog.Info("getting a user", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)


		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		user, err := storage.GetUserById(intId)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, user)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")
		students, err := storage.GetUsers()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}