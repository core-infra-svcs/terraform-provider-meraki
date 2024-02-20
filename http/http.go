package httpclient

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// HttpClient is a http client wrapper that can be used by services to enable retries.
/*
usage:
// Create a new HttpClient.
retries := 1 // number of times to retry
retryDelay := 3 // Seconds to wait before retrying
connectionTimeout := 10 // Seconds to wait for a request to return
logger := lumberjack.Logger{}
c := httpclient.New(1, 3, 10 logger)
// c can now be used as a service client to retry any GET requests. To enable retry on all methods set that on the object.
c.SetRetryAllMethods()
// c will now retry all methods, use this only when dealing with idempotent services (like POST for FDWS).
*/
type HttpClient struct {
	RetryDelay time.Duration
	MaxRetries int
	client     *http.Client
	retryAll   bool
	apiToken   string
}

// Do is a wrapper method for http client.Do method. It will retry the request
// if the request method is GET up to MaxRetries times.
func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	if c.client == nil {
		c.client = &http.Client{}
	}
	if !c.shouldRetry(req) {
		return c.client.Do(req)
	}
	var err error
	var resp *http.Response
	for i := 0; i <= c.MaxRetries; i++ {
		retry := false
		errLog := ""

		if c.retryAll {
			// Always rewind the request body when non-nil.
			if req.Body != nil {
				body, err := req.GetBody()
				if err != nil {
					return resp, err
				}

				req.Body = io.NopCloser(body)
			}
		}

		resp, err = c.client.Do(req)
		// Make sure we're not trying to reference anything from resp or err if either are nil. Construct the log messages for each only if
		// we can dereference them.
		if resp != nil {
			if resp.StatusCode != http.StatusOK {
				retry = true
			}
			if resp.StatusCode != http.StatusTooManyRequests {
				time.Sleep(c.RetryDelay)
			}
		}
		if err != nil {
			retry = true
			errLog = fmt.Sprintf("- Error: %v", err)
			// Explicitly capture timeout errors
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				errLog += " - request timed out"
			}
		}
		if retry {
			// Clean up the log a little by removing any leading or trailing whitespace. This keeps the log spacing consistent if either resp or err is nil.
			time.Sleep(time.Duration(c.RetryDelay) * time.Second)
			continue
		}
		return resp, nil
	}
	return resp, err
}

// SetRetryAllMethods will set the client to retry all methods not just GET. This is because some services (like FDWS) treat a fetch of data
// as a POST instead of a GET. Be VERY careful when enabling this for a service to make sure that you understand what is going on. This is why it
// is a setter method instead of part of the constructor or a public property. This has to be explicitly set for a service client.
func (c *HttpClient) SetRetryAllMethods() {
	c.retryAll = true
}

func (c *HttpClient) SetTransport(tripper http.RoundTripper) {
	c.client.Transport = tripper
}

func (c *HttpClient) shouldRetry(req *http.Request) bool {
	if c.retryAll {
		return true
	}
	return req.Method == http.MethodGet
}

// New creates a new HttpClient instance.
func New(maxRetries int, retryDelay, connectionTimeout time.Duration) *HttpClient {
	return &HttpClient{
		MaxRetries: maxRetries,
		RetryDelay: retryDelay,
		client: &http.Client{
			Timeout: connectionTimeout * time.Second,
		},
		retryAll: false,
	}
}
