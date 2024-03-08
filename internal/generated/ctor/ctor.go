package ctor

/* DO NOT EDIT. This file generated due to '@ctor' annotation*/

import (
	a "github.com/fjrid/parking/internal/app/infra"
	b "github.com/fjrid/parking/internal/app/repo/postgresql/parking"
	c "github.com/fjrid/parking/internal/app/repo/postgresql/user"
	d "github.com/fjrid/parking/internal/app/service/auth"
	e "github.com/fjrid/parking/internal/app/service/parking"
	f "github.com/fjrid/parking/internal/app/service/user"
	"github.com/typical-go/typical-go/pkg/typapp"
)

func init() {
	typapp.Provide("", a.NewCacheStore)
	typapp.Provide("", a.NewDatabases)
	typapp.Provide("", a.NewEcho)
	typapp.Provide("", b.NewParkingRepository)
	typapp.Provide("", c.NewUserRepository)
	typapp.Provide("", d.NewAuthSvc)
	typapp.Provide("", e.NewParkingSvc)
	typapp.Provide("", f.NewUserSvc)
}
