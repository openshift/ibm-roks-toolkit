package metrics_pusher

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type MetricsPusher struct {
	Log             logr.Logger
	SourceURL       string
	DestinationPath string
	Kubeconfig      string
	Job             string
	Frequency       time.Duration

	restClient rest.Interface
	pushPath   []string
}

const (
	requestTimeout = 2 * time.Minute
)

func (p *MetricsPusher) Run() error {
	var err error
	p.restClient, err = p.getRESTClient()
	if err != nil {
		p.Log.Error(err, "Failed to get REST client for target")
		return err
	}
	p.pushPath = p.getPushPath()
	p.Log.Info("Polling for metrics", "from", p.SourceURL, "to", p.DestinationPath)
	wait.Forever(p.pushMetrics, p.Frequency)
	return nil
}

func (p *MetricsPusher) pushMetrics() {
	resp, err := http.Get(p.SourceURL)
	if err != nil {
		p.Log.Error(err, "failed to fetch metrics from source")
		return
	}
	if resp.StatusCode != http.StatusOK {
		p.Log.Error(fmt.Errorf("Status: %s (%d)", resp.Status, resp.StatusCode), "Unexpected status from source URL")
		return
	}
	metricsBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.Log.Error(err, "Failed to read metrics body from source URL")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	p.Log.Info("Pushing metrics")
	result := p.restClient.Post().AbsPath(p.pushPath...).Body(metricsBody).Do(ctx)
	if result.Error() != nil {
		p.Log.Error(result.Error(), "failed to post metrics")
		return
	}
}

func (p *MetricsPusher) getPushPath() []string {
	parts := strings.Split(p.DestinationPath, "/")
	parts = append(parts, "job", p.Job)
	return parts
}

func (p *MetricsPusher) getRESTClient() (rest.Interface, error) {
	kubeconfig := clientcmd.GetConfigFromFileOrDie(p.Kubeconfig)
	clientConfig := clientcmd.NewDefaultClientConfig(*kubeconfig, &clientcmd.ConfigOverrides{})
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfigOrDie(restConfig).CoreV1().RESTClient(), nil
}
