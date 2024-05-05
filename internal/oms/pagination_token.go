package oms

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chrisrob11/oms/internal/oms/db"
	"github.com/gin-gonic/gin"
)

const TokenQueryParamName = "$token"

// PaginationToken holds information for pagination.
type PaginationToken struct {
	StartID int `json:"start_id"`
	Size    int `json:"size"`
}

// EncodeToken encodes the PaginationToken to base64.
func EncodeToken(token PaginationToken) string {
	data, _ := json.Marshal(token)
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeToken decodes the base64-encoded token to PaginationToken.
func DecodeToken(encodedToken string) (PaginationToken, error) {
	data, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return PaginationToken{}, err
	}

	var token PaginationToken
	err = json.Unmarshal(data, &token)

	return token, err
}

type QueryPagingInfo struct {
	Limit      int
	StartingID *int
}

func extractLimit(c *gin.Context) (*int32, bool) {
	limitStr := c.Query("$limit")
	if limitStr != "" {
		limitInt64, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return nil, true
		}

		if limitInt64 > 500 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit cannot be higher than 500"})
			return nil, true
		}

		limitInt := int32(limitInt64)

		return &limitInt, false
	}

	return nil, false
}

func extractTokenFromQuery(c *gin.Context) (paging *PaginationToken, hasError bool) {
	token := c.Query(TokenQueryParamName)
	params := db.ListCampaignLineItemsParams{
		Limit: 100,
	}

	if token != "" {
		pagingToken, err := DecodeToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return nil, true
		}

		params.ID = int32(pagingToken.StartID)
		params.Limit = int32(pagingToken.Size)

		return &PaginationToken{StartID: pagingToken.StartID, Size: pagingToken.Size}, false
	}

	return nil, false
}
