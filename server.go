package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"farmstall/openapi"
	"farmstall/reviews"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	Reviews *reviews.Reviews
}

// main
func main() {

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	server := Server{
		Reviews: reviews.NewReviews(),
	}

	server.initDummyData()
	m := mux.NewRouter()

	// API
	api := m.PathPrefix("/v1").Subrouter()
	api.HandleFunc("/reviews", server.getReviews()).Methods(http.MethodGet)
	api.HandleFunc("/reviews", server.addReview()).Methods(http.MethodPost)
	api.HandleFunc("/reviews/{reviewId}", server.getReview()).Methods(http.MethodGet)
	api.HandleFunc("/reviews/{reviewId}", server.deleteReview()).Methods(http.MethodDelete)
	api.HandleFunc("/reviews/{reviewId}", server.updateReview()).Methods(http.MethodPut)

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
	})

	// Wrap in CORS
	handler := c.Handler(m)

	// Create a rate limiter, 1 per second ( 3600 per hour )

	rate, _ := limiter.NewRateFromFormatted("36-M")
	store := memory.NewStore()
	rater := stdlib.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))

	handler = rater.Handler(handler)

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
			ErrorResponse(400, err.Error())(w, r)
			return
		}

		reviewRes, err := ctx.Reviews.UpdateReview(reviewId, review)
		if err != nil {
			ErrorResponse(400, err.Error())(w, r)
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
			ErrorResponse(404, err.Error())(w, r)
		} else {
			// Empty response
			w.WriteHeader(204)
			w.Write(nil)
		}
	}

}

func (ctx *Server) addReview() MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var review reviews.Review
		err := decoder.Decode(&review)
		if err != nil {
			ErrorResponse(400, err.Error())(w, r)
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
			if i <= 0 || i > 5 {
				ErrorResponse(404, "Query parameter 'maxRating' should only be a whole number between 1 and 5 inclusive")(w, r)
				return
			}
			filters := reviews.ReviewFilters{
				MaxRating: i,
			}
			reviewList, _ = ctx.Reviews.GetReviewsFiltered(filters)
		} else {
			reviewList, _ = ctx.Reviews.GetReviews()
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
			ErrorResponse(404, err.Error())(w, r)
		} else {
			writeJson(200, review)(w, r)
		}
	}

}

// Bunch of HTTP stuffs...
type MiddlewareFn func(http.ResponseWriter, *http.Request)

type ErrorResponseS struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func ErrorResponse(status int, msg string) MiddlewareFn {
	return func(w http.ResponseWriter, r *http.Request) {
		errRes := ErrorResponseS{Message: msg, Status: status}
		fmt.Printf("Response error: (%d) %s", status, msg)
		writeJson(status, errRes)(w, r)
	}
}

func response404(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(404, "Resource not found")(w, r)
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
