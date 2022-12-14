// Package generate provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package generate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Backend defines model for backend.
type Backend struct {
	HealthCheckNodePort *int     `json:"health_check_node_port,omitempty"`
	Server              []string `json:"server"`
}

// Config defines model for config.
type Config struct {
	Frontends []Frontend `json:"frontends"`
}

// Error defines model for error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Frontend defines model for frontend.
type Frontend struct {
	Backend Backend `json:"backend"`
	Port    int     `json:"port"`
}

// Loadbalancer defines model for loadbalancer.
type Loadbalancer struct {
	Config Config  `json:"config"`
	Name   *string `json:"name,omitempty"`
	Status *Status `json:"status,omitempty"`
}

// Status defines model for status.
type Status struct {
	Hostname *string `json:"hostname,omitempty"`
	Ip       *string `json:"ip,omitempty"`
}

// Name defines model for name.
type Name = string

// CreateLoadBalancerJSONBody defines parameters for CreateLoadBalancer.
type CreateLoadBalancerJSONBody = Loadbalancer

// CreateLoadBalancerJSONRequestBody defines body for CreateLoadBalancer for application/json ContentType.
type CreateLoadBalancerJSONRequestBody = CreateLoadBalancerJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetHealth request
	GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetLoadbalancers request
	GetLoadbalancers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteLoadBalancer request
	DeleteLoadBalancer(ctx context.Context, name Name, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetLoadbalancer request
	GetLoadbalancer(ctx context.Context, name Name, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateLoadBalancer request with any body
	CreateLoadBalancerWithBody(ctx context.Context, name Name, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateLoadBalancer(ctx context.Context, name Name, body CreateLoadBalancerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHealthRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetLoadbalancers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetLoadbalancersRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteLoadBalancer(ctx context.Context, name Name, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteLoadBalancerRequest(c.Server, name)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetLoadbalancer(ctx context.Context, name Name, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetLoadbalancerRequest(c.Server, name)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateLoadBalancerWithBody(ctx context.Context, name Name, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateLoadBalancerRequestWithBody(c.Server, name, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateLoadBalancer(ctx context.Context, name Name, body CreateLoadBalancerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateLoadBalancerRequest(c.Server, name, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetHealthRequest generates requests for GetHealth
func NewGetHealthRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/healthz")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetLoadbalancersRequest generates requests for GetLoadbalancers
func NewGetLoadbalancersRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/loadbalancers")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewDeleteLoadBalancerRequest generates requests for DeleteLoadBalancer
func NewDeleteLoadBalancerRequest(server string, name Name) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "name", runtime.ParamLocationPath, name)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/loadbalancers/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetLoadbalancerRequest generates requests for GetLoadbalancer
func NewGetLoadbalancerRequest(server string, name Name) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "name", runtime.ParamLocationPath, name)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/loadbalancers/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateLoadBalancerRequest calls the generic CreateLoadBalancer builder with application/json body
func NewCreateLoadBalancerRequest(server string, name Name, body CreateLoadBalancerJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateLoadBalancerRequestWithBody(server, name, "application/json", bodyReader)
}

// NewCreateLoadBalancerRequestWithBody generates requests for CreateLoadBalancer with any type of body
func NewCreateLoadBalancerRequestWithBody(server string, name Name, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "name", runtime.ParamLocationPath, name)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/loadbalancers/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetHealth request
	GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error)

	// GetLoadbalancers request
	GetLoadbalancersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetLoadbalancersResponse, error)

	// DeleteLoadBalancer request
	DeleteLoadBalancerWithResponse(ctx context.Context, name Name, reqEditors ...RequestEditorFn) (*DeleteLoadBalancerResponse, error)

	// GetLoadbalancer request
	GetLoadbalancerWithResponse(ctx context.Context, name Name, reqEditors ...RequestEditorFn) (*GetLoadbalancerResponse, error)

	// CreateLoadBalancer request with any body
	CreateLoadBalancerWithBodyWithResponse(ctx context.Context, name Name, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateLoadBalancerResponse, error)

	CreateLoadBalancerWithResponse(ctx context.Context, name Name, body CreateLoadBalancerJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateLoadBalancerResponse, error)
}

type GetHealthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetHealthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetHealthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetLoadbalancersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Loadbalancer
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetLoadbalancersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetLoadbalancersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteLoadBalancerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r DeleteLoadBalancerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteLoadBalancerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetLoadbalancerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Loadbalancer
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetLoadbalancerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetLoadbalancerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateLoadBalancerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Loadbalancer
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r CreateLoadBalancerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateLoadBalancerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetHealthWithResponse request returning *GetHealthResponse
func (c *ClientWithResponses) GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error) {
	rsp, err := c.GetHealth(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHealthResponse(rsp)
}

// GetLoadbalancersWithResponse request returning *GetLoadbalancersResponse
func (c *ClientWithResponses) GetLoadbalancersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetLoadbalancersResponse, error) {
	rsp, err := c.GetLoadbalancers(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetLoadbalancersResponse(rsp)
}

// DeleteLoadBalancerWithResponse request returning *DeleteLoadBalancerResponse
func (c *ClientWithResponses) DeleteLoadBalancerWithResponse(ctx context.Context, name Name, reqEditors ...RequestEditorFn) (*DeleteLoadBalancerResponse, error) {
	rsp, err := c.DeleteLoadBalancer(ctx, name, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteLoadBalancerResponse(rsp)
}

// GetLoadbalancerWithResponse request returning *GetLoadbalancerResponse
func (c *ClientWithResponses) GetLoadbalancerWithResponse(ctx context.Context, name Name, reqEditors ...RequestEditorFn) (*GetLoadbalancerResponse, error) {
	rsp, err := c.GetLoadbalancer(ctx, name, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetLoadbalancerResponse(rsp)
}

// CreateLoadBalancerWithBodyWithResponse request with arbitrary body returning *CreateLoadBalancerResponse
func (c *ClientWithResponses) CreateLoadBalancerWithBodyWithResponse(ctx context.Context, name Name, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateLoadBalancerResponse, error) {
	rsp, err := c.CreateLoadBalancerWithBody(ctx, name, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateLoadBalancerResponse(rsp)
}

func (c *ClientWithResponses) CreateLoadBalancerWithResponse(ctx context.Context, name Name, body CreateLoadBalancerJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateLoadBalancerResponse, error) {
	rsp, err := c.CreateLoadBalancer(ctx, name, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateLoadBalancerResponse(rsp)
}

// ParseGetHealthResponse parses an HTTP response from a GetHealthWithResponse call
func ParseGetHealthResponse(rsp *http.Response) (*GetHealthResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetHealthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseGetLoadbalancersResponse parses an HTTP response from a GetLoadbalancersWithResponse call
func ParseGetLoadbalancersResponse(rsp *http.Response) (*GetLoadbalancersResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetLoadbalancersResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Loadbalancer
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseDeleteLoadBalancerResponse parses an HTTP response from a DeleteLoadBalancerWithResponse call
func ParseDeleteLoadBalancerResponse(rsp *http.Response) (*DeleteLoadBalancerResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteLoadBalancerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseGetLoadbalancerResponse parses an HTTP response from a GetLoadbalancerWithResponse call
func ParseGetLoadbalancerResponse(rsp *http.Response) (*GetLoadbalancerResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetLoadbalancerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Loadbalancer
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseCreateLoadBalancerResponse parses an HTTP response from a CreateLoadBalancerWithResponse call
func ParseCreateLoadBalancerResponse(rsp *http.Response) (*CreateLoadBalancerResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateLoadBalancerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest Loadbalancer
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /healthz)
	GetHealth(ctx echo.Context) error

	// (GET /loadbalancers)
	GetLoadbalancers(ctx echo.Context) error

	// (DELETE /loadbalancers/{name})
	DeleteLoadBalancer(ctx echo.Context, name Name) error

	// (GET /loadbalancers/{name})
	GetLoadbalancer(ctx echo.Context, name Name) error

	// (PUT /loadbalancers/{name})
	CreateLoadBalancer(ctx echo.Context, name Name) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetHealth converts echo context to params.
func (w *ServerInterfaceWrapper) GetHealth(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetHealth(ctx)
	return err
}

// GetLoadbalancers converts echo context to params.
func (w *ServerInterfaceWrapper) GetLoadbalancers(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetLoadbalancers(ctx)
	return err
}

// DeleteLoadBalancer converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteLoadBalancer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "name" -------------
	var name Name

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, ctx.Param("name"), &name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteLoadBalancer(ctx, name)
	return err
}

// GetLoadbalancer converts echo context to params.
func (w *ServerInterfaceWrapper) GetLoadbalancer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "name" -------------
	var name Name

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, ctx.Param("name"), &name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetLoadbalancer(ctx, name)
	return err
}

// CreateLoadBalancer converts echo context to params.
func (w *ServerInterfaceWrapper) CreateLoadBalancer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "name" -------------
	var name Name

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, ctx.Param("name"), &name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateLoadBalancer(ctx, name)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/healthz", wrapper.GetHealth)
	router.GET(baseURL+"/loadbalancers", wrapper.GetLoadbalancers)
	router.DELETE(baseURL+"/loadbalancers/:name", wrapper.DeleteLoadBalancer)
	router.GET(baseURL+"/loadbalancers/:name", wrapper.GetLoadbalancer)
	router.PUT(baseURL+"/loadbalancers/:name", wrapper.CreateLoadBalancer)

}
