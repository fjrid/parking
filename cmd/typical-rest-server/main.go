package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fjrid/parking/internal/app"
	"github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/envkit"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typgo"

	// Important to enable dependency injectino
	_ "github.com/fjrid/parking/internal/generated/ctor"
	_ "github.com/fjrid/parking/internal/generated/envcfg"
)

func main() {
	// Print application header
	fmt.Printf("Start %s (%s) at %s\n",
		typgo.ProjectName, typgo.ProjectVersion, time.Now().Format(time.RFC3339))

	// Read dotenv file
	if dotenv := os.Getenv("CONFIG"); dotenv != "" {
		fmt.Printf("Set ENV from '%s'", dotenv)
		m, _ := envkit.ReadFile(dotenv)
		envkit.Setenv(m)
	}

	if err := typapp.StartApp(app.Start, app.Shutdown); err != nil {
		logrus.Fatal(err.Error())
	}

}
