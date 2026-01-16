// Package supabase provides database operations using PostgREST API.
package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// FilterOperator represents a PostgREST filter operator.
type FilterOperator string

const (
	OpEq    FilterOperator = "eq"
	OpNeq   FilterOperator = "neq"
	OpGt    FilterOperator = "gt"
	OpGte   FilterOperator = "gte"
	OpLt    FilterOperator = "lt"
	OpLte   FilterOperator = "lte"
	OpLike  FilterOperator = "like"
	OpILike FilterOperator = "ilike"
	OpIn    FilterOperator = "in"
	OpIs    FilterOperator = "is"
)

// Filter represents a query filter.
type Filter struct {
	Column   string
	Operator FilterOperator
	Value    string
}

// QueryBuilder provides a fluent API for building database queries.
type QueryBuilder struct {
	client    *Client
	table     string
	columns   []string
	filters   []Filter
	orderBy   string
	orderAsc  bool
	limit     int
	offset    int
	single    bool
	userToken string
}

// From starts a new query on the specified table.
func (c *Client) From(table string) *QueryBuilder {
	return &QueryBuilder{
		client:  c,
		table:   table,
		columns: []string{"*"},
	}
}

// Select specifies which columns to return.
func (q *QueryBuilder) Select(columns ...string) *QueryBuilder {
	if len(columns) > 0 {
		q.columns = columns
	}
	return q
}

// Filter adds a filter condition to the query.
func (q *QueryBuilder) Filter(column string, operator FilterOperator, value string) *QueryBuilder {
	q.filters = append(q.filters, Filter{
		Column:   column,
		Operator: operator,
		Value:    value,
	})
	return q
}

// Eq is a shorthand for Filter with OpEq.
func (q *QueryBuilder) Eq(column, value string) *QueryBuilder {
	return q.Filter(column, OpEq, value)
}

// Order specifies the column to order by.
func (q *QueryBuilder) Order(column string, ascending bool) *QueryBuilder {
	q.orderBy = column
	q.orderAsc = ascending
	return q
}

// Limit sets the maximum number of rows to return.
func (q *QueryBuilder) Limit(n int) *QueryBuilder {
	q.limit = n
	return q
}

// Offset sets the number of rows to skip.
func (q *QueryBuilder) Offset(n int) *QueryBuilder {
	q.offset = n
	return q
}

// Single expects a single result (returns error if not exactly one row).
func (q *QueryBuilder) Single() *QueryBuilder {
	q.single = true
	q.limit = 1
	return q
}

// WithToken sets the user's access token for RLS policies.
func (q *QueryBuilder) WithToken(token string) *QueryBuilder {
	q.userToken = token
	return q
}

// Execute runs the query and unmarshals the result into dest.
func (q *QueryBuilder) Execute(dest interface{}) error {
	reqURL := q.buildURL()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}

	q.setHeaders(req)

	if q.single {
		req.Header.Set("Accept", "application/vnd.pgrst.object+json")
	}

	resp, err := q.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("query failed: %s", string(body))
	}

	return json.Unmarshal(body, dest)
}

func (q *QueryBuilder) buildURL() string {
	baseURL := fmt.Sprintf("%s/rest/v1/%s", q.client.baseURL, q.table)

	params := url.Values{}
	params.Set("select", strings.Join(q.columns, ","))

	for _, f := range q.filters {
		params.Add(f.Column, fmt.Sprintf("%s.%s", f.Operator, f.Value))
	}

	if q.orderBy != "" {
		order := q.orderBy
		if !q.orderAsc {
			order += ".desc"
		}
		params.Set("order", order)
	}

	if q.limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", q.limit))
	}

	if q.offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", q.offset))
	}

	return baseURL + "?" + params.Encode()
}

func (q *QueryBuilder) setHeaders(req *http.Request) {
	req.Header.Set("apikey", q.client.apiKey)
	req.Header.Set("Content-Type", "application/json")

	if q.userToken != "" {
		req.Header.Set("Authorization", "Bearer "+q.userToken)
	} else {
		req.Header.Set("Authorization", "Bearer "+q.client.apiKey)
	}
}

// Insert adds one or more rows to the table.
func (c *Client) Insert(table string, data interface{}, userToken string) error {
	return c.mutate("POST", table, data, nil, userToken, false)
}

// InsertReturning adds rows and returns the inserted data.
func (c *Client) InsertReturning(table string, data interface{}, result interface{}, userToken string) error {
	return c.mutateReturning("POST", table, data, nil, userToken, result)
}

// Update modifies rows matching the filters.
func (c *Client) Update(table string, data interface{}, filters []Filter, userToken string) error {
	return c.mutate("PATCH", table, data, filters, userToken, false)
}

// UpdateReturning modifies rows and returns the updated data.
func (c *Client) UpdateReturning(table string, data interface{}, filters []Filter, result interface{}, userToken string) error {
	return c.mutateReturning("PATCH", table, data, filters, userToken, result)
}

// Delete removes rows matching the filters.
func (c *Client) Delete(table string, filters []Filter, userToken string) error {
	return c.mutate("DELETE", table, nil, filters, userToken, false)
}

func (c *Client) mutate(method, table string, data interface{}, filters []Filter, userToken string, returning bool) error {
	reqURL := c.buildMutateURL(table, filters)

	var body []byte
	var err error
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	c.setMutateHeaders(req, userToken, returning)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s failed: %s", method, string(respBody))
	}

	return nil
}

func (c *Client) mutateReturning(method, table string, data interface{}, filters []Filter, userToken string, result interface{}) error {
	reqURL := c.buildMutateURL(table, filters)

	var body []byte
	var err error
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	c.setMutateHeaders(req, userToken, true)
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s failed: %s", method, string(respBody))
	}

	return json.Unmarshal(respBody, result)
}

func (c *Client) buildMutateURL(table string, filters []Filter) string {
	baseURL := fmt.Sprintf("%s/rest/v1/%s", c.baseURL, table)

	if len(filters) == 0 {
		return baseURL
	}

	params := url.Values{}
	for _, f := range filters {
		params.Add(f.Column, fmt.Sprintf("%s.%s", f.Operator, f.Value))
	}

	return baseURL + "?" + params.Encode()
}

func (c *Client) setMutateHeaders(req *http.Request, userToken string, returning bool) {
	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	if userToken != "" {
		req.Header.Set("Authorization", "Bearer "+userToken)
	} else {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	if returning {
		req.Header.Set("Prefer", "return=representation")
	}
}
