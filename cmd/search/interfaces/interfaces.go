package interfaces

import "github.com/golang-jwt/jwt"

type SearchQueryParams struct {
	CourseOfferingId string `form:"courseOfferingId"`
	ActivityType     string `form:"activityType"`
	Offset           int    `form:"offset"`
	Limit            int    `form:"limit"`
}

type Token struct {
	UserId    string `json:"custom:userId"`
	OrgId     string `json:"custom:orgId"`
	Role      string `json:"custom:role"`
	ContactId string `json:"custom:contactId"`
	jwt.StandardClaims
}

type ISearchPayload struct {
	Value            string `json:"value"`
	Criteria         string `json:"criteria"`
	Role             string `json:"role"`
	OrgId            string `json:"orgId"`
	CourseOfferingId string `json:"courseOfferingId"`
	ActivityTypeKey  string `json:"activityTypeKey"`
	Offset           int    `json:"offset"`
	Limit            int    `json:"limit"`
}

type Hit struct {
	Source map[string]interface{} `json:"_source"`
}

type SearchResponse struct {
	Hits struct {
		Hits []Hit `json:"hits"`
	} `json:"hits"`
}

type ESResponse struct {
	Aggregations Aggregations `json:"aggregations"`
}

type Aggregations struct {
	ActivityType ActivityType `json:"activity_type"`
}

type ActivityType struct {
	Buckets []Bucket `json:"buckets"`
}

type Bucket struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
	Hits     Hits   `json:"activity_type"` // Assuming activity_type contains hits
}

type Hits struct {
	Hits []Hit `json:"hits"`
}

type ActivityData struct {
	// Define your activity data fields here
	Name string `json:"name"`
}

type Result struct {
	Type  string         `json:"type"`
	Count int            `json:"count"`
	Data  []ActivityData `json:"data"`
}
