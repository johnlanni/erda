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

package cmdb

import (
	"net/http"

	"github.com/erda-project/erda/modules/openapi/api/apis"
)

var CMDB_FILE_HEAD = apis.ApiSpec{
	Path:          "/api/files/<uuid>",
	BackendPath:   "/api/files/<uuid>",
	Host:          "cmdb.marathon.l4lb.thisdcos.directory:9093",
	Scheme:        "http",
	Method:        http.MethodHead,
	CheckLogin:    false,
	TryCheckLogin: true,
	CheckToken:    true,
	IsOpenAPI:     true,
	Doc:           "summary: HEAD 请求，文件下载，在 path 中指定具体文件",
}
