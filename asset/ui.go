// +build dev

package asset

import (
	"net/http"

	"github.com/shurcooL/httpfs/union"
)

var Assets http.FileSystem = union.New(map[string]http.FileSystem{
	"/check":    http.Dir("./check"),
	"/settings": http.Dir("./settings"),
})
