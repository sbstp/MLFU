package drivers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/publicsuffix"
)

type Body struct {
	ContentType string
	Data        []byte
}

func (b Body) Reader() io.Reader {
	return bytes.NewReader(b.Data)
}

func marshalFieldsMultipart(fields map[string]string) Body {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for key, value := range fields {
		mw.WriteField(key, value)
	}
	mw.Close()
	return Body{
		ContentType: mw.FormDataContentType(),
		Data:        buf.Bytes(),
	}
}

func marshalFieldsURLEncoded(fields map[string]string) Body {
	query := url.Values{}
	for key, value := range fields {
		query.Set(key, value)
	}
	return Body{
		ContentType: "application/x-www-form-urlencoded",
		Data:        []byte(query.Encode()),
	}
}

func newClient() *http.Client {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	return &http.Client{
		Jar:     jar,
		Timeout: time.Second * 15,
	}
}
