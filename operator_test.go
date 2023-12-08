package go_k8s

import (
	"context"
	"flag"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

var K8sClient client.Client
var setupLog = ctrl.Log.WithName("setup")
var scheme = runtime.NewScheme()

func TestMain(m *testing.M) {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
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
	Job.K8sClient = K8sClient
	go func() {
		err = manager.Start(context.Background())
		if err != nil {
			setupLog.Error(err, "unable to start manager")
			os.Exit(1)
		}
	}()
	m.Run()
}
func Test_Run(t *testing.T) {
	err := Job.Run("job-container-4", "default", "liburdi/go-k8s-job-container:0.0.4", "", []string{""})
	if err != nil {
		t.Log(err)
	}

}
