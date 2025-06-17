package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

type ApiClient struct {
	apiUrl    *url.URL
	authToken string
	client    HTTPClient
	UserAgent string

	FeatureTags    *FeatureTagsService
	FeatureToggles *FeatureTogglesService
	Projects       *ProjectsService
	FeatureTypes   *FeatureTypesService
	Strategies     *StrategiesService
	Users          *UsersService
}

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	userAgent = "go-unleash-api/api/" + LibraryVersion
)

func (c *ApiClient) setApiUrl(urlStr string) error {
	if urlStr == "" {
		return ErrApiUrlCannotBeEmpty
	}

	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.apiUrl, err = url.Parse(urlStr)
	return err
}

func (c *ApiClient) setAuthToken(authToken string) error {
	if authToken == "" {
		return ErrTokenAuthCannotBeEmpty
	}
	c.authToken = authToken

	return nil
}

func NewClient(httpClient HTTPClient, apiUrl string, authToken string) (*ApiClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}
	c := &ApiClient{}
	if err := c.setApiUrl(apiUrl); err != nil {
		return nil, err
	}
	if err := c.setAuthToken(authToken); err != nil {
		return nil, err
	}

	c.client = httpClient
	c.UserAgent = userAgent
	c.FeatureTags = &FeatureTagsService{client: c}
	c.FeatureToggles = &FeatureTogglesService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.FeatureTypes = &FeatureTypesService{client: c}
	c.Strategies = &StrategiesService{client: c}
	c.Users = &UsersService{client: c}

	return c, nil
}

func (c *ApiClient) newRequest(path string, method string, opt interface{}) (*http.Request, error) {
	var u = *c.apiUrl
	u.Opaque = c.apiUrl.Path + path

	if opt != nil && isStruct(opt) {
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req := &http.Request{
		Method:     method,
		URL:        &u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
	req.Header.Set("User-Agent", userAgent)

	if method == "POST" || method == "PUT" {
		bodyBytes, err := json.Marshal(opt)
		if err != nil {
			return nil, err
		}
		bodyReader := bytes.NewReader(bodyBytes)

		u.RawQuery = ""
		req.Body = ioutil.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.authToken)

	return req, nil
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

type Response struct {
	*http.Response
}

func (c *ApiClient) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil && response.StatusCode != http.StatusNoContent {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}

func isStruct(s interface{}) bool {
	return reflect.ValueOf(s).Kind() == reflect.Struct
}
