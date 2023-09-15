package main

import (
	"fmt"
	"github.com/Southclaws/fault"
	apihttp "github.com/Southclaws/fault/examples/api/http"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

var ExistingUserID = "123"
var ErrorUserID = "999"
var logger *slog.Logger

func init() {
	logger = slog.Default()
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if _, ok := v.(error); ok {

			// We change the response to not reveal the actual error message,
			// instead we can transform the message something more friendly or mapped
			// to some code / language, etc.
			render.DefaultResponder(w, r, render.M{"status": "error"})
			return
		}

		render.DefaultResponder(w, r, v)
	}
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == ExistingUserID {
		user := User{
			ID:   userID,
			Name: "Bob",
		}
		render.JSON(w, r, user)
		return
	}

	if userID == ErrorUserID {
		err := fault.New("db error: connection lost")
		err = fault.Wrap(err,
			fmsg.WithDesc("Could not get user", "An error occured while getting the user. Try again later"),
			ftag.With(ftag.Internal),
		)
		apihttp.RespondWithError(logger, err, w, r)
		return
	}

	err := fault.New(fmt.Sprintf("db error: user id[%s] not found", userID))
	err = fault.Wrap(err,
		fmsg.WithDesc("User not found", "Cannot find the requested user"),
		ftag.With(ftag.NotFound))
	apihttp.RespondWithError(logger, err, w, r)
}

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(apihttp.DecorateRequestMetadata)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(apihttp.LoggerRequest(logger))

	r.Route("/users/{userID}", func(r chi.Router) {
		r.Use(apihttp.PathVariableAsFCtx("userID", "user_id"))
		r.Get("/", GetUser)
	})

	fmt.Printf("Listening on :3333 ...\n")
	http.ListenAndServe(":3333", r)
}
