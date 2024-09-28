package interfaces

import "github.com/golang-jwt/jwt"

type SearchQueryParams struct {
	CourseOfferingId string `json:"courseOfferingId"`
	ActivityType     string `json:"activityType"`
	Offset           int    `json:"offset"`
	Limit            int    `json:"limit"`
}

type Token struct {
	UserId    string `json:"userId"`
	OrgId     string `json:"orgId"`
	Role      string `json:"role"`
	ContactId string `json:"contactId"`
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
