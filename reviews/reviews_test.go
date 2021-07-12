package reviews

import (
	"testing"

	"encoding/json"
	"github.com/google/uuid"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestGetReviewsEmpty(t *testing.T) {
	reviews := NewReviews()
	res := reviews.GetReviews()
	assert.Assert(t, is.Len(*res, 0), "should return empty list of reviews")
}

func TestAddOneReview(t *testing.T) {
	reviews := NewReviews()
	_, err := reviews.AddReview(Review{
		Message: "good",
		Rating:  5,
	})
	assert.NilError(t, err, "should have no errors")
	assert.Assert(t, is.Len(reviews.Reviews, 1), "should update the list of reviews")
}

func TestAddThenGetOneReview(t *testing.T) {
	reviews := NewReviews()
	addedReview, err := reviews.AddReview(Review{
		Message: "good",
		Rating:  5,
	})

	gottenReview, err := reviews.GetReview(addedReview.Uuid)

	assert.NilError(t, err, "should have no errors")
	assert.Assert(t, is.DeepEqual(addedReview, gottenReview), "should match the review in memory")
}

func TestAddTwoThenGetAllReviews(t *testing.T) {
	reviews := NewReviews()
	reviews.AddReview(Review{
		Message: "good",
		Rating:  5,
	})
	reviews.AddReview(Review{
		Message: "poor",
		Rating:  1,
	})

	allReviews := reviews.GetReviews()
	assert.Assert(t, is.Len(*allReviews, 2), "should match the number of reviews added")
}

func TestDeleteOnReview(t *testing.T) {
	reviews := NewReviews()
	addedReview, _ := reviews.AddReview(Review{
		Message: "good",
		Rating:  5,
	})

	err := reviews.DeleteReview(addedReview.Uuid)
	assert.NilError(t, err, "should have no errors")
	assert.Assert(t, is.Len(reviews.Reviews, 0), "should have no reviews")
}

func TestUpdateReview(t *testing.T) {
	reviews := NewReviews()
	oriReview, _ := reviews.AddReview(Review{
		Message: "poor",
		Rating:  1,
	})

	newReview := Review{
		Message: "good",
		Rating:  5,
	}

	_, err := reviews.UpdateReview(oriReview.Uuid, newReview)
	assert.NilError(t, err, "should have no errors")

	updatedReview, err := reviews.GetReview(oriReview.Uuid)
	assert.NilError(t, err, "should have no errors")

	assert.Assert(t, is.Equal(updatedReview.Message, "good"), "should match the updated review's Message")
	assert.Assert(t, is.Equal(updatedReview.Rating, 5), "should match the updated review's Rating")
	assert.Assert(t, is.Len(reviews.Reviews, 1), "should not add any more reviews")
}

func TestUpdateReviewInvalidUuid(t *testing.T) {
	reviews := NewReviews()
	newReview := Review{
		Message: "good",
		Rating:  5,
		Uuid:    uuid.New().String(),
	}
	_, err := reviews.UpdateReview(newReview.Uuid, newReview)
	assert.ErrorContains(t, err, "Refusing to update a non-existing resource", "should return an error")
}

func TestGetAllReviews(t *testing.T) {
	reviews := NewReviews()
	reviews.AddReview(Review{
		Message: "good",
		Rating:  5,
	})
	reviews.AddReview(Review{
		Message: "average",
		Rating:  3,
	})
	reviews.AddReview(Review{
		Message: "poor",
		Rating:  1,
	})

	allReviews := reviews.GetReviews()
	assert.Assert(t, is.Len(*allReviews, 3), "should equal the number of reviews added")
}

func TestGetReviewsByMaxRatingBelowOrEquals(t *testing.T) {
	reviews := NewReviews()
	reviews.AddReview(Review{
		Message: "good",
		Rating:  5,
	})
	reviews.AddReview(Review{
		Message: "average",
		Rating:  3,
	})
	reviews.AddReview(Review{
		Message: "poor",
		Rating:  1,
	})

	filters := ReviewFilters{
		MaxRating: 3,
	}

	allReviews := reviews.GetReviewsFiltered(filters)

	assert.Assert(t, is.Len(*allReviews, 2), "should equal two, for the reviews with ratings 1 and 3")
}

func TestGetReviewsByMaxRatingEquals(t *testing.T) {
	reviews := NewReviews()
	reviews.AddReview(Review{
		Message: "good",
		Rating:  5,
	})
	reviews.AddReview(Review{
		Message: "average",
		Rating:  3,
	})
	reviews.AddReview(Review{
		Message: "poor",
		Rating:  1,
	})

	filters := ReviewFilters{
		MaxRating: 1,
	}

	allReviews := reviews.GetReviewsFiltered(filters)

	assert.Assert(t, is.Len(*allReviews, 1), "should equal one, for the review with rating 1")
}

func TestAnonymousReviewHasNullForUserID(t *testing.T) {

	reviews := NewReviews()
	review, _ := reviews.AddReview(Review{
		Message: "good",
		Rating:  5,
	})

	jsonBytes, _ := json.Marshal(review)

	assert.Assert(t, is.Contains(string(jsonBytes), `"userId":null`), "should equal null when serialized in JSON")
}
