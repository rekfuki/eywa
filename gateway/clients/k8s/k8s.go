package k8s

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corelister "k8s.io/client-go/listers/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	faasNamespace   = "faas"
	faasIDLabel     = "function"
	faasSecretMount = "/var/faas/secrets"
)

type Config struct {
	InCluster bool
}

type Client struct {
	clientset      *kubernetes.Clientset
	endpointLister corelister.EndpointsNamespaceLister
}

type functionLookup struct {
	endpointLister corelister.EndpointsLister
	listers        map[string]corelister.EndpointsNamespaceLister

	lock sync.RWMutex
}

// Setup initialises a k8s client
func Setup(conf *Config) (*Client, error) {
	var config *rest.Config
	var err error
	if conf.InCluster {
		config, err = rest.InClusterConfig()
	} else {
		home := homedir.HomeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	defaultResync := time.Second * 1

	kubeInformerOpt := kubeinformers.WithNamespace(faasNamespace)
	kubeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(clientset, defaultResync, kubeInformerOpt)
	go startFactory(kubeInformerFactory)

	endpointsInformer := kubeInformerFactory.Core().V1().Endpoints()
	lister := endpointsInformer.Lister()

	return &Client{
		clientset:      clientset,
		endpointLister: lister.Endpoints(faasNamespace),
	}, nil
}

func (c *Client) Resolve(fnName string) (string, error) {
	if strings.Contains(fnName, ".") {
		fnName = strings.TrimSuffix(fnName, "."+faasNamespace)
	}

	svc, err := c.endpointLister.Get(fnName)
	if err != nil {
		return "", fmt.Errorf("Error listing \"%s.%s\": %s", fnName, faasNamespace, err)
	}

	if len(svc.Subsets) == 0 {
		return "", fmt.Errorf("No subsets available for \"%s.%s\"", fnName, faasNamespace)
	}

	all := len(svc.Subsets[0].Addresses)
	if len(svc.Subsets[0].Addresses) == 0 {
		return "", fmt.Errorf("No addresses in subset for \"%s.%s\"", fnName, faasNamespace)
	}

	target := rand.Intn(all)

	serviceIP := svc.Subsets[0].Addresses[target].IP

	return fmt.Sprintf("http://%s:%d", serviceIP, 8080), nil
}

func startFactory(f kubeinformers.SharedInformerFactory) {
	stop := make(chan struct{})
	defer close(stop)

	go f.Start(stop)

	for {
		<-stop
		log.Errorf("K8s informer factory stopped...")
		os.Exit(0)
	}
}
