package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var config *Config

// Config API responser settings
type Config struct {
	ContentType string
}

func (c *Config) getContentType() string {
	if c.ContentType == "" {
		c.ContentType = "application/json"
	}
	return c.ContentType
}

// Setup API general settings
func Setup(c *Config) {
	config = c
}

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
	ID        string      `json:"id,omitempty"`
	Status    int         `json:"-"`
	Code      int         `json:"code,omitempty"`
	Title     string      `json:"title,omitempty"`
	Detail    interface{} `json:"detail,omitempty"`
	Links     Links       `json:"links,omitempty"`
	PrevError error       `json:"-"`
}

// Error text
func (e *Error) Error() string {
	return e.Title
}

// Resp function builds API response
func Resp(c *gin.Context, r *Response) {
	if r.HTTPStatus == 0 {
		r.HTTPStatus = http.StatusOK
	}
	c.JSON(r.HTTPStatus, r)
}

// Err sends error response
func Err(c *gin.Context, code int, msg interface{}, prevErr error) {
	r := Response{}
	if code != 0 {
		r.HTTPStatus = code
	} else {
		r.HTTPStatus = http.StatusInternalServerError
	}
	err := Error{
		ID:        uuid.NewV1().String(),
		Code:      code,
		Title:     http.StatusText(code),
		Detail:    msg,
		PrevError: prevErr,
	}
	r.AddError(err)
	if prevErr != nil {
		log.Println(err)
	}
	c.AbortWithStatusJSON(code, r)
}

// CheckContentType middleware
func CheckContentType(c *gin.Context) {
	act := strings.ToLower(config.getContentType())
	if strings.ToLower(c.Request.Method) == "get" {
		if c.GetHeader("Accept") != "*/*" && strings.ToLower(c.GetHeader("Accept")) != act {
			Err(c, http.StatusNotAcceptable, http.StatusText(http.StatusNotAcceptable), nil)
			return
		}
	} else {
		if strings.ToLower(c.ContentType()) != act {
			Err(c, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType), nil)
			return
		}
		if c.GetHeader("Accept") != "*/*" && strings.ToLower(c.GetHeader("Accept")) != act {
			Err(c, http.StatusNotAcceptable, http.StatusText(http.StatusNotAcceptable), nil)
			return
		}
	}
	c.Next()
}
