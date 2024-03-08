package main

import (
	"github.com/fjrid/parking/pkg/typcfg"
	"github.com/fjrid/parking/pkg/typdb"
	"github.com/fjrid/parking/pkg/typdocker"
	"github.com/fjrid/parking/pkg/typredis"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typgen"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
)

var descriptor = typgo.Descriptor{
	ProjectName:    "typical-rest-server",
	ProjectVersion: "0.9.21",
	Environment:    typgo.DotEnv(".env"),

	Tasks: []typgo.Tasker{
		// generate
		&typgen.Generator{
			Processor: typgen.Processors{
				&typapp.CtorAnnot{},
				&typdb.DBRepoAnnot{},
				&typcfg.EnvconfigAnnot{GenDotEnv: ".env", GenDoc: "USAGE.md"},
			},
		},
		// test
		&typgo.GoTest{
			Includes: []string{"internal/app/**", "pkg/**"},
		},
		// compile
		&typgo.GoBuild{},
		// run
		&typgo.RunBinary{
			Before: typgo.TaskNames{"build"},
		},
		// mock
		&typmock.GoMock{},
		// docker
		&typdocker.DockerTool{
			ComposeFiles: typdocker.ComposeFiles("deploy"),
			EnvFile:      ".env",
		},
		// pg
		&typdb.PostgresTool{Name: "pg"},
		// cache
		&typredis.RedisTool{Name: "cache"},
		// setup
		&typgo.Task{
			Name:  "setup",
			Usage: "setup dependency",
			Action: typgo.TaskNames{
				"pg drop", "pg create", "pg migrate", "pg seed",
			},
		},
		// release
		&typrls.ReleaseProject{
			Before: typgo.TaskNames{"test", "build"},
			// Releaser:  &typrls.CrossCompiler{Targets: []typrls.Target{"darwin/amd64", "linux/amd64"}},
			Publisher: &typrls.Github{Owner: "typical-go", Repo: "typical-rest-server"},
		},
	},
}

func main() {
	typgo.Start(&descriptor)
}
