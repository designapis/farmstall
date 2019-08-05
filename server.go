package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"farmstall/openapi"
	"farmstall/problems"
	"farmstall/reviews"
	"farmstall/users"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	Reviews *reviews.Reviews
	Users   *users.Users
}

// Set from ENV variable during startup
var PROBS_URL string
var BASE_URL string

const BASE_PATH string = "/v1"

// main
func main() {

	PORT := os.Getenv("PORT")
	FQDN := os.Getenv("FQDN")

	if PORT == "" {
		PORT = "8080"
	}

	if FQDN == "" {
		FQDN = "https://farmstall.ponelat.com"
	}

	// Set global
	PROBS_URL = FQDN + "/probs"
	BASE_URL = FQDN + BASE_PATH

	server := Server{
		Reviews: reviews.NewReviews(),
		Users:   users.NewUsers(),
	}

	server.initDummyData()
	m := mux.NewRouter()

	// API
	api := m.PathPrefix(BASE_PATH).Subrouter()
	api.Use(server.validateRequestMiddleware)

	api.HandleFunc("/reviews", server.getReviews()).Methods(http.MethodGet)
	api.HandleFunc("/reviews", server.addReview()).Methods(http.MethodPost)
	api.HandleFunc("/reviews/{reviewId}", server.getReview()).Methods(http.MethodGet)
	api.HandleFunc("/reviews/{reviewId}", server.deleteReview()).Methods(http.MethodDelete)
	api.HandleFunc("/reviews/{reviewId}", server.updateReview()).Methods(http.MethodPut)

	api.HandleFunc("/users", server.addUser()).Methods(http.MethodPost)
	api.HandleFunc("/tokens", server.createToken()).Methods(http.MethodPost)

	// UI
	m.HandleFunc("/health", server.health())

	// OpenAPI
	m.HandleFunc("/openapi.yaml", openapi.Openapi(openapi.YAML))
	m.HandleFunc("/openapi.json", openapi.Openapi(openapi.JSON))
	m.HandleFunc("/openapi", openapi.Openapi(openapi.ANY))

	// CORS for friendlinesss
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowOriginFunc:  func(origin string) bool { return true },
	})

	// Wrap in CORS
	handler := c.Handler(m)

	// Create a rate limiter, 1 per second ( 3600 per hour )
	rate, _ := limiter.NewRateFromFormatted("36-M")
	store := memory.NewStore()
	rater := stdlib.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))
	handler = rater.Handler(handler)

	// Final handler
	http.Handle("/", handler)

	fmt.Printf("Listening on :%s\n", PORT)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		panic(err)
	}
}

func (ctx *Server) initDummyData() {
	ctx.Reviews.AddReview(reviews.Review{
		Message: "Was awesome!",
		Rating:  5,
	})
	ctx.Reviews.AddReview(reviews.Review{
		Message: "Was okay.",
		Rating:  3,
	})
	ctx.Reviews.AddReview(reviews.Review{
		Message: "Was terrible.",
		Rating:  1,
	})
}

// Validate the incoming request against our schema(s)
func (ctx *Server) validateRequestMiddleware(next http.Handler) http.Handler {
	return http.StripPrefix("/v1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Find operation
		router := openapi3filter.NewRouter().WithSwaggerFromFile("openapi.yaml")
		route, pathParams, errOp := router.FindRoute(r.Method, r.URL)

		if errOp != nil {
			log.Fatalf("Operation not found for %s %s. Error: %s", r.Method, r.URL, errOp)
		}

		// Validate request against operation
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
		}

		something := context.TODO()

		if err := openapi3filter.ValidateRequest(something, requestValidationInput); err != nil {
			switch errVal := err.(type) {
			case *openapi3filter.RequestError:
				ErrorResponse(problems.InvalidBody(problems.ProblemJson{
					Detail: errVal.Reason,
				}))(w, r)
				return
			case *openapi3filter.SecurityRequirementsError:
				// Allow this for now ( optional securities appear to be an issue )
				log.Printf("errVal %s", errVal)
				break
			default:
				ErrorResponse(problems.InvalidRequest(problems.ProblemJson{}))(w, r)
				return
			}
		}

		// All good, carry on...
		next.ServeHTTP(w, r)
	}))
}

func (ctx *Server) health() MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJson(200, map[string]bool{"healthy": true})(w, r)
	}
}

func (ctx *Server) updateReview() MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		vars := mux.Vars(r)
		reviewId := vars["reviewId"]

		decoder := json.NewDecoder(r.Body)
		var review reviews.Review
		err = decoder.Decode(&review)
		if err != nil {
			ErrorResponse(problems.FailedToParseJson(problems.ProblemJson{
				Detail: err.Error(),
			}))(w, r)
			return
		}

		reviewRes, err := ctx.Reviews.UpdateReview(reviewId, review)
		if err != nil {
			ErrorResponse(err.(*problems.ProblemJson))(w, r)
		} else {
			// Empty response
			writeJson(200, reviewRes)(w, r)
		}
	}
}

func (ctx *Server) deleteReview() MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		reviewId := vars["reviewId"]
		err := ctx.Reviews.DeleteReview(reviewId)
		if err != nil {
			ErrorResponse(err.(*problems.ProblemJson))(w, r)
		} else {
			// Empty response
			w.WriteHeader(204)
			w.Write(nil)
		}
	}

}

func (ctx *Server) addUser() MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var user users.NewUser
		err := decoder.Decode(&user)
		if err != nil {
			ErrorResponse(problems.FailedToParseJson(problems.ProblemJson{
				Detail: err.Error(),
			}))(w, r)
			return
		}
		res := ctx.Users.AddUser(user)
		writeJson(201, res)(w, r)
	}
}

func (ctx *Server) createToken() MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var user users.UserLogin
		err := decoder.Decode(&user)
		if err != nil {
			ErrorResponse(problems.FailedToParseJson(problems.ProblemJson{
				Detail: err.Error(),
			}))(w, r)
			return
		}

		token, tokenErr := ctx.Users.CreateToken(user)
		if tokenErr != nil {
			ErrorResponse(tokenErr.(*problems.ProblemJson))(w, r)
			return
		}

		tokenRes := users.TokenResponse{
			Token: token,
		}
		writeJson(201, tokenRes)(w, r)
	}

}

func (ctx *Server) addReview() MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var review reviews.Review
		err := decoder.Decode(&review)
		if err != nil {
			ErrorResponse(problems.FailedToParseJson(problems.ProblemJson{
				Detail: err.Error(),
			}))(w, r)
			return
		}
		res, _ := ctx.Reviews.AddReview(review)
		writeJson(201, res)(w, r)
	}

}

func (ctx *Server) getReviews() MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		maxRating := query.Get("maxRating")
		log.Printf("\nmaxRating: %s\n", maxRating)
		var reviewList *[]reviews.Review
		if maxRating != "" {
			i, _ := strconv.Atoi(maxRating)

			// if i <= 0 || i > 5 {
			// 	ErrorResponse(404, "Query parameter 'maxRating' should only be a whole number between 1 and 5 inclusive")(w, r)
			// 	return
			// }

			filters := reviews.ReviewFilters{
				MaxRating: i,
			}
			reviewList = ctx.Reviews.GetReviewsFiltered(filters)
		} else {
			reviewList = ctx.Reviews.GetReviews()
		}
		writeJson(200, reviewList)(w, r)
	}
}

func (ctx *Server) getReview() MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		reviewId := vars["reviewId"]
		review, err := ctx.Reviews.GetReview(reviewId)
		if err != nil {
			ErrorResponse(err.(*problems.ProblemJson))(w, r)
		} else {
			writeJson(200, review)(w, r)
		}
	}

}

// Bunch of HTTP stuffs...
type MiddlewareFn func(http.ResponseWriter, *http.Request)

func ErrorResponse(prob *problems.ProblemJson) MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(prob.Error())
		prob := problems.Absolutify(prob, PROBS_URL, BASE_URL)
		writeJson(prob.Status, prob)(w, r)
	}
}

func writeJson(status int, msg interface{}) MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("%s", msg)
		w.Header().Add("Content-type", "application/json")
		msgBytes, _ := json.Marshal(msg)
		w.WriteHeader(status)
		w.Write([]byte(msgBytes))
	}
}
