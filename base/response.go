package base

import (
	"github.com/cloudwego/hertz/pkg/app"
)

type ErrorObj struct {
	Code      string
	Message   string
	MessageCn string
	Detail    string `json:",omitempty"`
}

type ResponseMetaData struct {
	RequestId string
	Error     *ErrorObj `json:",omitempty"`
}

type ErrorResponse struct {
	Meta ResponseMetaData
}

type CommonResponse struct {
	Meta ResponseMetaData
	Data any `json:",omitempty"`
}

func RespondJson[T any](args Args, c *app.RequestContext, data T) {
	c.Header("X-Request-Id", args.RequestId)
	c.JSON(200, CommonResponse{
		Meta: ResponseMetaData{RequestId: args.RequestId},
		Data: data,
	})
	c.Abort()
}

func RespondError(args Args, c *app.RequestContext, except Exception) {
	c.Header("X-Request-Id", args.RequestId)
	c.JSON(except.StatusCode, ErrorResponse{
		Meta: ResponseMetaData{
			RequestId: args.RequestId,
			Error: &ErrorObj{
				Code:      except.Code,
				Message:   except.Message,
				MessageCn: except.MessageCn,
				Detail:    except.RawError,
			},
		},
	})
	c.Abort()
}
