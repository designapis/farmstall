FarmStall API
-------------

*This is a contrived API* As part of the ongoing Swagger In Action book by Manning.

The API is hosted at https://farmstall.ponelat.com/v1

## FarmStall API

The FarmStall service helps the Farmer's Market by encouraging feedback on the experience from the patrons that visit.

The patrons can create reviews, and view the reviews of others. They can do so anonymously or with a user account.

Reviews are messages ( in markdown format ), with a corresponding rating ( 1 to 5 inclusive ) that helps broadly categorize the feedback into shades of positive/negative. Where a rating of 5 is the most postive type of review.

## Swagger in Action TODOs
- [x] Add maxRating to GET /reviews
- [x] Add rate limiting

- [x] Add user store
- [x] Add token store ( assoc with user )
- [x] Add problem+json for error responses

- [x] Add auth to POST /reviews and GET /reviews/{reviewId}

- [ ] Prevent creating a user with the same username 

