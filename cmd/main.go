package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	go_k8s "github.com/liburdi/go-k8s"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var K8sClient client.Client
var setupLog = ctrl.Log.WithName("setup")
var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func main() {
	var name *string = flag.String("name", "default name", "help message for name")
	var image *string = flag.String("image", "default image", "help message for image")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	logger := zap.New(zap.UseFlagOptions(&opts))
	ctrl.SetLogger(logger)

	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	K8sClient = manager.GetClient()
	go_k8s.Job.K8sClient = K8sClient

	logger.Info(fmt.Sprintf("param %s %s", *image, *name))

	go func() {
		time.Sleep(10 * time.Second)
		err = go_k8s.Job.Run(*name, "go-k8s-operator", *image, "", []string{""})
		if err != nil {
			logger.Info(fmt.Sprintf("param %s %s %s", *image, *name, err.Error()))
		}
	}()

	err = manager.Start(context.Background())
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
}
