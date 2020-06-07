package translator

import (
	"github.com/hexya-erp/hexya/src/server"
)

const MODULE_NAME string = "translator"

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PreInit:  func() {},
		PostInit: func() {},
	})
}
