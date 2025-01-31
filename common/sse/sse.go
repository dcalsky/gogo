package sse

import (
	"github.com/dcalsky/gogo/base"
	"io"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/protocol/http1/resp"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/network"
)

var (
	hId                 = []byte("id:")
	headerData          = []byte("data:")
	headerEvent         = []byte("event:")
	headerRetry         = []byte("retry:")
	breakingLine        = []byte("\n")
	doubleBreakingLines = []byte("\n\n")
)

const (
	ContentType  = "text/event-stream"
	noCache      = "no-cache"
	cacheControl = "Cache-Control"
	LastEventID  = "Last-Event-ID"
)

var fieldReplacer = strings.NewReplacer(
	"\n", "\\n",
	"\r", "\\r")

var dataReplacer = strings.NewReplacer(
	"\n", "\ndata:",
	"\r", "\\r")

type Event struct {
	Event string
	ID    string
	Retry uint64
	Data  []byte
}

// GetLastEventID retrieve Last-Event-ID header if present.
func GetLastEventID(c *app.RequestContext) string {
	return c.Request.Header.Get(LastEventID)
}

type HertzStream struct {
	w network.ExtWriter
}

// NewHertzStream creates a new stream for publishing Event.
func NewHertzStream(c *app.RequestContext) *HertzStream {
	c.Response.Header.SetContentType(ContentType)
	if c.Response.Header.Get(cacheControl) == "" {
		c.Response.Header.Set(cacheControl, noCache)
	}

	writer := resp.NewChunkedBodyWriter(&c.Response, c.GetWriter())
	c.Response.HijackWriter(writer)
	return &HertzStream{
		writer,
	}
}

// Publish push an event to client. If error is returned, you need to stop 'publish'.
func (c *HertzStream) Publish(event *Event) error {
	err := Encode(c.w, event)
	if err != nil {
		return err
	}
	return c.w.Flush()
}

func Encode(w io.Writer, e *Event) (err error) {
	if e.ID != "" {
		err = writeID(w, e.ID)
		if err != nil {
			return
		}
	}
	if e.Event != "" {
		err = writeEvent(w, e.Event)
		if err != nil {
			return
		}
	}
	if e.Retry > 0 {
		err = writeRetry(w, e.Retry)
		if err != nil {
			return
		}
	}
	err = writeData(w, e.Data)
	if err != nil {
		return
	}
	return nil
}

func writeID(w io.Writer, id string) (err error) {
	_, err = w.Write(hId)
	if err != nil {
		return
	}
	//_, err = fieldReplacer.WriteString(w, id)
	//if err != nil {
	//	return
	//}
	_, err = w.Write(base.S2b(id))
	if err != nil {
		return
	}
	_, err = w.Write(breakingLine)
	if err != nil {
		return
	}

	return
}

func writeEvent(w io.Writer, event string) (err error) {
	_, err = w.Write(headerEvent)
	if err != nil {
		return
	}
	//_, err = fieldReplacer.WriteString(w, event)
	//if err != nil {
	//	return
	//}
	_, err = w.Write(base.S2b(event))
	if err != nil {
		return
	}
	_, err = w.Write(breakingLine)
	if err != nil {
		return
	}
	return
}

func writeRetry(w io.Writer, retry uint64) (err error) {
	_, err = w.Write(headerRetry)
	if err != nil {
		return
	}
	_, err = w.Write(base.S2b(strconv.FormatUint(retry, 10)))
	if err != nil {
		return
	}
	_, err = w.Write(breakingLine)
	if err != nil {
		return
	}
	return
}

func writeData(w io.Writer, data []byte) (err error) {
	_, err = w.Write(headerData)
	if err != nil {
		return
	}

	_, err = dataReplacer.WriteString(w, base.B2s(data))
	if err != nil {
		return
	}

	_, err = w.Write(doubleBreakingLines)
	if err != nil {
		return
	}

	return nil
}
