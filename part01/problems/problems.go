package problems

import (
	"fmt"
)

type ProblemJson struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	Status     int    `json:"status"`
	Instance   string `json:"instance"`
	isAbsolute bool
}

func (pj ProblemJson) Error() string {
	return fmt.Sprintf("%d\t%s\t%s\t%s\t%s", pj.Status, pj.Instance, pj.Type, pj.Title, pj.Detail)
}

func NotFound(pj ProblemJson) *ProblemJson {
	return &ProblemJson{
		Type:     "/not-found",
		Title:    "Resource not found",
		Status:   404,
		Detail:   pj.Detail,
		Instance: pj.Instance,
	}
}

func InvalidCreds(pj ProblemJson) *ProblemJson {
	return &ProblemJson{
		Type:     "/invalid-credentials",
		Title:    "Invalid credentials provided",
		Status:   403,
		Detail:   pj.Detail,
		Instance: pj.Instance,
	}
}

func InvalidRequest(pj ProblemJson) *ProblemJson {
	return &ProblemJson{
		Type:     "/invalid-request",
		Title:    "Invalid request",
		Status:   400,
		Detail:   pj.Detail,
		Instance: pj.Instance,
	}
}

func InvalidBody(pj ProblemJson) *ProblemJson {
	return &ProblemJson{
		Type:     "/invalid-request-body",
		Title:    "Invalid body provided in request",
		Status:   400,
		Detail:   pj.Detail,
		Instance: pj.Instance,
	}
}

func FailedToParseJson(pj ProblemJson) *ProblemJson {
	return &ProblemJson{
		Type:     "/failed-to-parse-json",
		Title:    "Failed to parse the JSON",
		Status:   400,
		Detail:   pj.Detail,
		Instance: pj.Instance,
	}
}

func UpdateNonExisting(pj ProblemJson) *ProblemJson {
	return &ProblemJson{
		Type:     "/update-non-existing",
		Title:    "Refusing to update a non-existing resource. Create one first",
		Status:   400,
		Detail:   pj.Detail,
		Instance: pj.Instance,
	}
}

func CreateAlreadyExists(pj ProblemJson) *ProblemJson {
	return &ProblemJson{
		Type:     "/create-already-exists",
		Title:    "Failed to create resource, it already exists.",
		Status:   409,
		Detail:   pj.Detail,
		Instance: pj.Instance,
	}
}

func Absolutify(pj *ProblemJson, probBase string, apiBase string) ProblemJson {
	pj.Type = probBase + pj.Type

	if pj.Instance != "" {
		pj.Instance = apiBase + pj.Instance
	}

	return *pj
}
