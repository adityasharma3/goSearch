package search

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/adityasharma3/goSearch/cmd/search/interfaces"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	var searchParams interfaces.SearchQueryParams

	criteria := c.Param("criteria")
	value := c.Param("value")

	if criteria == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Criteria or Value cannot be empty"})
		log.Fatal("Search criteria or Value cannot be empty")
	}

	if criteria != "Exact" || criteria != "Contains" {
		log.Fatal("Such a criteria does not exist")
	}

	if err := c.ShouldBindQuery(&searchParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		log.Fatal(`Search Params do not match with validations`)
	}

	jwtToken := c.GetHeader("Authorization")
	decodedToken, err := decodeJWT(jwtToken)
	// log.Print(tokenData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err})
		log.Fatal(`JWT token could not be decoded`)
	}

	response := getSearchQueryResults()

}

func getSearchQueryResults(payload interfaces.ISearchPayload) {
	// aggregations := make(map[string]any)
	var mustConditions []map[string]interface{}

	if payload.Role == "Student" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"status":    "published",
				"isDeleted": false,
				"orgId":     payload.OrgId,
			},
		})
	}

	if payload.CourseOfferingId != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"courseOfferingId": payload.CourseOfferingId,
			},
		})
	}

	if payload.ActivityTypeKey != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"activityType": payload.ActivityTypeKey,
			},
		})
	}

	// matchCase := make(map[string]interface{})
	var matchCase []map[string]interface{}

	if payload.Criteria == "Contains" {
		processedValue := strings.Split(payload.Value, " ")
		for i, word := range processedValue {
			processedValue[i] = "*" + word + "*"
		}

		matchCase = append(matchCase, map[string]interface{}{
			"query_string": map[string]interface{}{
				"default_field": "title",
				"query":         strings.Join(processedValue, " "),
			},
		}, map[string]interface{}{
			"query_string": map[string]interface{}{
				"default_field": "description",
				"query":         strings.Join(processedValue, " "),
			},
		})
	} else if payload.Criteria == "Exact" {
		matchCase = append(matchCase, map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"description": payload.Value,
			},
		}, map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"title": payload.Value,
			},
		})
	}

	// response := // connect to elasticsearch here
}

func decodeJWT(token string) (*interfaces.Token, error) {
	token = token[len("Bearer"):]
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invald token expected 3 parts, got %d", len(parts))
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding token payload: %v", err)
	}

	var claims *interfaces.Token
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return nil, fmt.Errorf("error parsing token payload: %v", err)
	}

	return claims, nil
}
