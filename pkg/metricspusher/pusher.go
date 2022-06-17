package metricspusher

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"

	knet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type MetricsPusher struct {
	Log             logr.Logger
	SourceURL       string
	SourcePath      string
	DestinationPath string
	Kubeconfig      string
	Clientca        string
	Frequency       time.Duration

	restClient rest.Interface
	pushPath   []string
	readPath   []string
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
	p.splitPaths()
	p.Log.Info("Polling for metrics", "from", p.SourceURL, "to", p.DestinationPath)
	wait.Forever(p.pushMetrics, p.Frequency)
	return nil
}

func (p *MetricsPusher) pushMetrics() {
	var metricsBody []byte
	if len(p.SourcePath) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		var err error
		metricsBody, err = p.restClient.Get().AbsPath(p.readPath...).Do(ctx).Raw()
		if err != nil {
			p.Log.Error(err, "failed to fetch metrics from source path")
			return
		}
	} else {
		var resp *http.Response
		var err error
		if p.Clientca != "" {
			var caCert []byte
			caCert, err = ioutil.ReadFile(p.Clientca)
			if err != nil {
				p.Log.Error(err, "Unable to read CA cert for fetching metrics: Please provide valid CA cert or path")
				return
			}

			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			tlsConfig := tls.Config{
				RootCAs:    caCertPool,
				MinVersion: tls.VersionTLS12,
			}

			tr := knet.SetTransportDefaults(&http.Transport{
				TLSHandshakeTimeout: 10 * time.Second,
				TLSClientConfig:     &tlsConfig,
			})

			client := &http.Client{Transport: tr}
			resp, err = client.Get(p.SourceURL)
		} else {
			resp, err = http.Get(p.SourceURL)
		}

		if err != nil {
			p.Log.Error(err, "failed to fetch metrics from source URL")
			return
		}
		if resp.StatusCode != http.StatusOK {
			p.Log.Error(fmt.Errorf("Status: %s (%d)", resp.Status, resp.StatusCode), "Unexpected status from source URL")
			return
		}

		metricsBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			p.Log.Error(err, "Failed to read metrics body from source URL")
			return
		}
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

func (p *MetricsPusher) splitPaths() {
	p.pushPath = strings.Split(p.DestinationPath, "/")
	p.readPath = strings.Split(p.SourcePath, "/")

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
