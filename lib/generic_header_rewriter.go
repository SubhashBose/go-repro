package lib

import (
	"net/http"
	"strings"
)

type GenericHeaderRewriter struct{}

func (*GenericHeaderRewriter) RewriteSpecifiedHeaders(keys []string, headers http.Header, ctx RequestContext) (rewritten bool) {
	for _, key := range keys {

		if value := headers.Get(key); value != "" {
			for _, mapping := range ctx.HostMappings() {
				if strings.Contains(value, mapping.remote) {
					value = strings.Replace(value, mapping.remote, mapping.host, -1)
					rewritten = true
				}
			}

			headers.Set(key, value)
		}
	}

	return
}

func (*GenericHeaderRewriter) RewriteSpecifiedIncomingHeaders(keys []string, headers http.Header, ctx RequestContext) (rewritten bool) {
	for _, key := range keys {

		if value := headers.Get(key); value != "" {
			for _, mapping := range ctx.HostMappings() {
				if strings.Contains(value, mapping.host) {
					value = strings.Replace(value, mapping.host, mapping.remote, -1)
					rewritten = true
				}
			}

			headers.Set(key, value)
		}
	}

	return
}
