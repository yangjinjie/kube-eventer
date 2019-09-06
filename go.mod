module github.com/yangjinjie/kube-eventer

go 1.12

require (
	github.com/yangjinjie/kube-eventer v1.0.0
	github.com/Shopify/sarama v1.22.1
	github.com/denverdino/aliyungo v0.0.0-20190410085603-611ead8a6fed
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.1
	github.com/google/cadvisor v0.33.1
	github.com/influxdata/influxdb v1.7.7
	github.com/pborman/uuid v1.2.0
	github.com/prometheus/client_golang v1.0.0
	github.com/riemann/riemann-go-client v0.4.0
	github.com/smartystreets/go-aws-auth v0.0.0-20180515143844-0c1422d1fdb9
	github.com/stretchr/testify v1.3.0
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980
	gopkg.in/olivere/elastic.v3 v3.0.75
	gopkg.in/olivere/elastic.v5 v5.0.81
	k8s.io/api v0.0.0-20190627205229-acea843d18eb
	k8s.io/apimachinery v0.0.0-20190627205106-bc5732d141a8
	k8s.io/apiserver v0.0.0-20190606205144-71ebb8303503
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v0.3.1
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190606204050-af9c91bd2759
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190606205144-71ebb8303503
)
