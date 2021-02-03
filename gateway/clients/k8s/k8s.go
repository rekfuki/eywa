package k8s

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corelister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	// Required for auth to init
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const (
	faasNamespace   = "faas"
	faasIDLabel     = "function"
	faasSecretMount = "/var/faas/secrets"

	faasMinReplicasIDLabel = "faas.replicas.min"
	faasMaxReplicasIDLabel = "faas.replicas.max"
	faasScaleFactorIDLabel = "faas.scale.factor"

	updatedAtLabel = "updated_at"

	defaultMinReplicas   = 0
	defaultMaxReplicas   = 100
	defaultScalingFactor = 20

	limitRangeName = "resources-min-max"

	dockerPullSecret = "image-pull-secret"
)

// Config represents the configuration of k8s client
type Config struct {
	InCluster           bool
	CacheExpiryDuration time.Duration
	LimitCPUMin         string
	LimitMemMin         string
	LimitCPUMax         string
	LimitMemMax         string
}

// Client represents the k8s client
type Client struct {
	clientset      *kubernetes.Clientset
	endpointLister corelister.EndpointsNamespaceLister
	limitRange     ResourceLimits
	cache          *cache.Cache
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

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	defaultResync := time.Second * 1

	kubeInformerOpt := kubeinformers.WithNamespace(faasNamespace)
	kubeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(clientset, defaultResync, kubeInformerOpt)
	go startFactory(kubeInformerFactory)

	endpointsInformer := kubeInformerFactory.Core().V1().Endpoints()
	endpointsLister := endpointsInformer.Lister()

	return &Client{
		clientset:      clientset,
		endpointLister: endpointsLister.Endpoints(faasNamespace),
		cache:          cache.New(conf.CacheExpiryDuration, conf.CacheExpiryDuration),
		limitRange: ResourceLimits{
			MinCPU: conf.LimitCPUMin,
			MaxCPU: conf.LimitCPUMax,
			MinMem: conf.LimitMemMin,
			MaxMem: conf.LimitMemMax,
		},
	}, nil
}

// GetLimits returns imposed resource limits under the namespace where functions are running
func (c *Client) GetLimits() *ResourceLimits {
	return &c.limitRange
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
