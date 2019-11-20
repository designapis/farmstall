package reviews

import (
	"encoding/json"
	"farmstall/problems"
	"farmstall/utils"
	"github.com/google/uuid"
	_ "log"
)

const BASE_PATH = "/reviews"

type ReviewMap map[string]Review

type Review struct {
	Message string `json:"message"`
	Rating  int    `json:"rating"`
	Uuid    string `json:"uuid"`
	UserID  string `json:"-"`
}
type ReviewFilters struct {
	MaxRating int `json:"maxRating"`
}

type Reviews struct {
	Reviews map[string]Review `json:"reviews"`
}

type DeletedReview struct {
	Uuid string `json:"uuid"`
}

func NewReviews() *Reviews {
	rs := Reviews{Reviews: ReviewMap{}}
	return &rs
}

func (rs *Reviews) UpdateReview(reviewId string, r Review) (*Review, error) {
	if _, ok := rs.Reviews[reviewId]; !ok {
		return nil, problems.UpdateNonExisting(problems.ProblemJson{
			Instance: BASE_PATH + "/" + reviewId,
		})
	}

	r.Uuid = reviewId
	rs.Reviews[reviewId] = r
	return &r, nil
}

func (rs *Reviews) AddReview(r Review) (*Review, error) {
	uuidVal := uuid.New().String()
	r.Uuid = uuidVal
	rs.Reviews[uuidVal] = r
	return &r, nil
}

func (rs *Reviews) GetReview(id string) (*Review, error) {
	var review Review
	var ok bool
	review, ok = rs.Reviews[id]
	if !ok {
		return nil, problems.NotFound(problems.ProblemJson{
			Instance: BASE_PATH + "/" + id,
		})
	}
	return &review, nil
}

func (rs *Reviews) DeleteReview(id string) error {
	if _, ok := rs.Reviews[id]; !ok {
		return problems.NotFound(problems.ProblemJson{
			Instance: BASE_PATH + "/" + id,
		})
	}
	delete(rs.Reviews, id)
	return nil
}

func (rs *Reviews) GetReviews() *[]Review {
	v := make([]Review, 0, len(rs.Reviews))
	for _, value := range rs.Reviews {
		v = append(v, value)
	}
	return &v
}

func (rs *Reviews) GetReviewsFiltered(filters ReviewFilters) *[]Review {
	v := make([]Review, 0, len(rs.Reviews))
	for _, value := range rs.Reviews {
		if value.Rating <= filters.MaxRating {
			v = append(v, value)
		}
	}
	return &v
}

// JSON marshal/unmarshal
// Allows for UserID to be null

type ReviewAlias Review
type ReviewJSON struct {
	ReviewAlias
	UserID utils.NullString `json:"userId"`
}

func NewReviewJSON(r Review) ReviewJSON {
	rj := ReviewJSON{}
	rj.ReviewAlias = ReviewAlias(r)
	rj.UserID = utils.NullString(r.UserID)
	return rj
}

func (rj ReviewJSON) toObj() Review {
	r := Review(rj.ReviewAlias)
	r.UserID = string(rj.UserID)
	return r
}

func (r Review) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewReviewJSON(r))
}

func (r *Review) UnmarshalJSON(data []byte) error {
	var rj ReviewJSON
	if err := json.Unmarshal(data, &rj); err != nil {
		return err
	}
	*r = rj.toObj()
	return nil
}
