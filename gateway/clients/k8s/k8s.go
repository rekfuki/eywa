package k8s

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	resourcev1 "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	defaultMinReplicas   = 0
	defaultMaxReplicas   = 100
	defaultScalingFactor = 20

	limitRangeName = "resources-min-max"

	updatedAtLabel = "updated_at"
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
	limitRange     *v1.LimitRange
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

	limitRange, err := setupLimits(clientset, conf)
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
		limitRange:     limitRange,
		cache:          cache.New(conf.CacheExpiryDuration, conf.CacheExpiryDuration),
	}, nil
}

// GetLimits returns imposed resource limits under the namespace where functions are running
func (c *Client) GetLimits() *ResourceLimits {
	return &ResourceLimits{
		MinCPU: c.limitRange.Spec.Limits[0].Min.Cpu().String(),
		MaxCPU: c.limitRange.Spec.Limits[0].Max.Cpu().String(),
		MinMem: c.limitRange.Spec.Limits[0].Min.Memory().String(),
		MaxMem: c.limitRange.Spec.Limits[0].Max.Memory().String(),
	}
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

func setupLimits(clientset *kubernetes.Clientset, conf *Config) (*v1.LimitRange, error) {
	context := context.TODO()
	err := clientset.CoreV1().LimitRanges(faasNamespace).Delete(context, limitRangeName, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
	}

	maxCPU, err := resourcev1.ParseQuantity(conf.LimitCPUMax)
	if err != nil {
		return nil, err
	}

	minCPU, err := resourcev1.ParseQuantity(conf.LimitCPUMin)
	if err != nil {
		return nil, err
	}

	maxMem, err := resourcev1.ParseQuantity(conf.LimitMemMax)
	if err != nil {
		return nil, err
	}

	minMem, err := resourcev1.ParseQuantity(conf.LimitMemMin)
	if err != nil {
		return nil, err
	}

	limitRange := &v1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name:      limitRangeName,
			Namespace: faasNamespace,
		},
		Spec: v1.LimitRangeSpec{
			Limits: []v1.LimitRangeItem{
				{
					Type: v1.LimitTypeContainer,
					Max: v1.ResourceList{
						v1.ResourceCPU:    maxCPU,
						v1.ResourceMemory: maxMem,
					},
					Min: v1.ResourceList{
						v1.ResourceCPU:    minCPU,
						v1.ResourceMemory: minMem,
					},
				},
			},
		},
	}

	limitRange, err = clientset.CoreV1().LimitRanges(faasNamespace).Create(context, limitRange, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return limitRange, nil
}
