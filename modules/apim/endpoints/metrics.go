// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package endpoints

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/discover"
	"github.com/erda-project/erda/pkg/httpserver"
	"github.com/erda-project/erda/pkg/strutil"
)

func ProxyMetrics(ctx context.Context, r *http.Request, vars map[string]string) error {

	// proxy
	r.URL.Scheme = "http"
	r.Host = discover.Monitor()
	r.URL.Host = discover.Monitor()
	r.URL.Path = strings.Replace(r.URL.Path, "/api/apim/metrics", "/api/metrics", 1)

	return nil
}

func InternalReverseHandler(handler func(context.Context, *http.Request, map[string]string) error) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			logrus.Debugf("start %s %s", r.Method, r.URL.String())

			handleRequest(r)

			err := handler(context.Background(), r, mux.Vars(r))
			if err != nil {
				logrus.Errorf("failed to handle request: %s (%v)", r.URL.String(), err)
				return
			}
		},
		FlushInterval: -1,
	}
}

func handleRequest(r *http.Request) {
	// base64 decode request body if declared in header
	if strutil.Equal(r.Header.Get(httpserver.Base64EncodedRequestBody), "true", true) {
		r.Body = ioutil.NopCloser(base64.NewDecoder(base64.StdEncoding, r.Body))
	}
}
