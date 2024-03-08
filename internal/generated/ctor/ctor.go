package ctor

/* DO NOT EDIT. This file generated due to '@ctor' annotation*/

import (
	a "github.com/fjrid/parking/internal/app/infra"
	b "github.com/fjrid/parking/internal/app/repo"
	c "github.com/fjrid/parking/internal/app/service"
	"github.com/typical-go/typical-go/pkg/typapp"
)

func init() {
	typapp.Provide("", a.NewCacheStore)
	typapp.Provide("", a.NewDatabases)
	typapp.Provide("", a.NewEcho)
	typapp.Provide("", b.NewBookRepo)
	typapp.Provide("", c.NewBookSvc)
}
