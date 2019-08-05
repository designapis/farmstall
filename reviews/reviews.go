package reviews

import (
	"farmstall/problems"
	"farmstall/users"
	"github.com/google/uuid"
)

const BASE_PATH = "/reviews"

type UUID = string
type ReviewMap map[UUID]Review
type Review struct {
	Uuid    UUID        `json:"uuid"`
	Message string      `json:"message"`
	Rating  int         `json:"rating"`
	User    *users.User `json:"userId"`
}

type ReviewFilters struct {
	MaxRating int `json:"maxRating"`
}

type Reviews struct {
	Reviews map[UUID]Review `json:"reviews"`
}

type DeletedReview struct {
	Uuid UUID `json:"uuid"`
}

func NewReviews() *Reviews {
	rs := Reviews{Reviews: ReviewMap{}}
	return &rs
}

func (rs *Reviews) UpdateReview(reviewId UUID, r Review) (*Review, error) {
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

func (rs *Reviews) GetReview(id UUID) (*Review, error) {
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

func (rs *Reviews) DeleteReview(id UUID) error {
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
