package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Response structure
type Response struct {
	HTTPStatus int         `json:"-"`
	Links      Links       `json:"links,omitempty"`
	Errors     []Error     `json:"errors,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Messages   []string    `json:"messages,omitempty"`
}

// SetData sets response data for BaseResponse
func (r *Response) SetData(data interface{}) {
	r.Data = data
}

// SetMeta sets response meta
func (r *Response) SetMeta(meta interface{}) {
	r.Meta = meta
}

// AddError adds error into errors array
func (r *Response) AddError(err Error) {
	r.Errors = append(r.Errors, err)
}

// AddMessage adds message into messages array
func (r *Response) AddMessage(msg string) {
	r.Messages = append(r.Messages, msg)
}

// AddLink adds link into links array
func (r *Response) AddLink(title, link string) {
	if r.Links == nil {
		r.Links = Links{}
	}
	r.Links[title] = link
}

// Links type
type Links map[string]string

// Data type
type Data map[string]interface{}

// Meta type
type Meta map[string]interface{}

// Error structure
type Error struct {
	ID     string      `json:"id,omitempty"`
	Status int         `json:"-"`
	Code   int         `json:"code,omitempty"`
	Title  string      `json:"title,omitempty"`
	Detail interface{} `json:"detail,omitempty"`
	Links  Links       `json:"links,omitempty"`
}

// Error text
func (e *Error) Error() string {
	return e.Title
}

// Response function builds API response
func Response(c *gin.Context, r *Response) {
	if r.HTTPStatus == 0 {
		r.HTTPStatus = http.StatusOK
	}
	c.JSON(r.HTTPStatus, r)
}

// Error sends error response
func Error(c *gin.Context, code int, msg interface{}) {
	r := Response{}
	if code != 0 {
		r.HTTPStatus = code
	} else {
		r.HTTPStatus = http.StatusInternalServerError
	}
	r.AddError(Error{
		ID:     uuid.Must(uuid.NewV1()).String(),
		Code:   code,
		Title:  http.StatusText(code),
		Detail: msg,
	})
	c.AbortWithStatusJSON(code, r)
}

func CheckContentType(c *gin.Context) {
	act := strings.ToLower(os.Getenv(allowedContentType))
	if strings.ToLower(c.Request.Method) == "get" {
		if c.GetHeader("Accept") != "*/*" && strings.ToLower(c.GetHeader("Accept")) != act {
			apiError(c, http.StatusNotAcceptable, http.StatusText(http.StatusNotAcceptable))
			return
		}
	} else {
		if strings.ToLower(c.ContentType()) != act {
			apiError(c, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}
		if c.GetHeader("Accept") != "*/*" && strings.ToLower(c.GetHeader("Accept")) != act {
			apiError(c, http.StatusNotAcceptable, http.StatusText(http.StatusNotAcceptable))
			return
		}
	}
	c.Next()
}
