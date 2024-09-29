package searchController

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adityasharma3/goSearch/cmd/search/interfaces"
	elasticsearch "github.com/adityasharma3/goSearch/cmd/search/searchclient"
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

	if criteria != "exact" && criteria != "contains" {
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

	params := interfaces.ISearchPayload{
		Value:            value,
		Criteria:         criteria,
		Role:             decodedToken.Role,
		OrgId:            decodedToken.OrgId,
		CourseOfferingId: searchParams.CourseOfferingId,
		ActivityTypeKey:  searchParams.ActivityType,
		Offset:           5,
		Limit:            10,
	}

	response, err := getSearchQueryResults(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, response)
}

func getSearchQueryResults(payload interfaces.ISearchPayload) ([]map[string]interface{}, error) {
	aggregations := make(map[string]any)
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

	ESClient := elasticsearch.GetElasticClient()
	dbPrefix := os.Getenv("DB_PREFIX")

	queryBody := map[string]interface{}{
		"from": 0,
		"size": 10,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": append(mustConditions, map[string]interface{}{
					"bool": map[string]interface{}{
						"should": matchCase,
					},
				}),
			},
		},
		"aggs": aggregations,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(queryBody); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	result, err := ESClient.Search(
		ESClient.Search.WithContext(context.Background()),
		ESClient.Search.WithIndex(dbPrefix+"-activitymetadata-index,"+dbPrefix+"-contentmetadata-index,"+dbPrefix+"-selftasks-index"),
		ESClient.Search.WithBody(&buf),
		ESClient.Search.WithTrackTotalHits(true),
		ESClient.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error getting response from Elasticsearch: %s", err)
	}

	defer result.Body.Close()

	var searchResponse interfaces.SearchResponse
	if err := json.NewDecoder(result.Body).Decode(&searchResponse); err != nil {
		log.Fatalf("error parsing elastic search response body %s", err)
	}

	var sources []map[string]interface{}
	for _, hit := range searchResponse.Hits.Hits {
		sources = append(sources, hit.Source)
	}

	return sources, nil
}

func decodeJWT(token string) (*interfaces.Token, error) {
	token = token[len("Bearer"):]
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token expected 3 parts, got %d", len(parts))
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
