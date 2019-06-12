package reviews

import (
	"farmstall/users"
	"fmt"

	"github.com/google/uuid"
)

const EDENIEDBYCODE = "e4033"
const ENOTFOUND = "e4040"

type UUID = string
type ReviewMap map[UUID]Review
type Review struct {
	Uuid    UUID        `json:"uuid"`
	Message string      `json:"message"`
	Rating  int         `json:"rating"`
	User    *users.User `json:"userId"`
}

type ReviewError struct {
	Msg  string `json:"error"`
	Uuid UUID   `json:"uuid"`
	Code string `json:"code"`
}

func (e ReviewError) Error() string {
	return fmt.Sprintf("%s: %s %s", e.Uuid, e.Msg, e.Code)
}

type Reviews struct {
	Reviews map[UUID]Review `json:"reviews"`
}

type DeletedReview struct {
	Uuid UUID `json:"uuid"`
}

func (rs *Reviews) UpdateReview(reviewId UUID, r Review) (*Review, error) {
	if _, ok := rs.Reviews[reviewId]; !ok {
		return nil, &ReviewError{
			Code: EDENIEDBYCODE,
			Msg:  "Review does not exist. Won't update. Please create new one.",
			Uuid: reviewId,
		}
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
		reErr := &ReviewError{Uuid: id, Msg: "Review not found", Code: ENOTFOUND}
		return nil, reErr
	}
	return &review, nil
}

func (rs *Reviews) DeleteReview(id UUID) error {
	if _, ok := rs.Reviews[id]; !ok {
		reErr := ReviewError{Uuid: id, Msg: "Review not found", Code: ENOTFOUND}
		return reErr
	}
	delete(rs.Reviews, id)
	return nil
}

func (rs *Reviews) GetReviews() (*[]Review, error) {
	v := make([]Review, 0, len(rs.Reviews))
	for _, value := range rs.Reviews {
		v = append(v, value)
	}
	return &v, nil
}

func NewReviews() *Reviews {
	rs := Reviews{Reviews: ReviewMap{}}
	return &rs
}
