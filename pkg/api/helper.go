package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ncarlier/webhookd/pkg/helper"
)

// HTTPParamsToShellVars convert URL values to shell vars.
func HTTPParamsToShellVars[T url.Values | http.Header](params T) []string {
	var result []string
	for k, v := range params {
		var buf bytes.Buffer
		value, err := url.QueryUnescape(strings.Join(v, ","))
		if err != nil {
			continue
		}
		buf.WriteString(helper.ToSnake(k))
		buf.WriteString("=")
		buf.WriteString(value)
		result = append(result, buf.String())
	}
	return result
}

func nextRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
