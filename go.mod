module linker-v3-example

go 1.26.5

require (
	github.com/neteast-software/go-module/acl v0.2.0
	github.com/neteast-software/go-module/acl/session v0.0.0
	github.com/neteast-software/go-module/application v0.2.0
	github.com/neteast-software/go-module/application/http/gin v0.0.0
	github.com/neteast-software/go-module/application/linker v0.2.0
	github.com/neteast-software/go-module/application/store/gorm v0.2.0
	github.com/neteast-software/go-module/audit/operate v0.2.0
	github.com/neteast-software/go-module/audit/operate/http/gin v0.2.0
	github.com/neteast-software/go-module/audit/operate/linker v0.2.0
	github.com/neteast-software/go-module/cache/redis v0.2.0
	github.com/neteast-software/go-module/cache/redis/linker v0.2.0
	github.com/neteast-software/go-module/config v0.2.0
	github.com/neteast-software/go-module/config/env/linker v0.0.0
	github.com/neteast-software/go-module/config/yaml/linker v0.0.0
	github.com/neteast-software/go-module/db/gorm/query v0.0.0
	github.com/neteast-software/go-module/db/gorm/table v0.2.0
	github.com/neteast-software/go-module/db/postgresql v0.2.0
	github.com/neteast-software/go-module/db/postgresql/linker v0.2.0
	github.com/neteast-software/go-module/fault/event v0.2.0
	github.com/neteast-software/go-module/fault/event/linker v0.2.0
	github.com/neteast-software/go-module/fault/notice v0.2.0
	github.com/neteast-software/go-module/fault/notice/linker v0.2.0
	github.com/neteast-software/go-module/http/client v0.0.0
	github.com/neteast-software/go-module/http/client/linker v0.0.0
	github.com/neteast-software/go-module/http/gin v0.2.0
	github.com/neteast-software/go-module/linker/server v0.2.0
	github.com/neteast-software/go-module/mq/consumer v0.2.0
	github.com/neteast-software/go-module/mq/consumer/linker v0.2.0
	github.com/neteast-software/go-module/mq/rocketmq v0.2.0
	github.com/neteast-software/go-module/mq/rocketmq/linker v0.2.0
	github.com/neteast-software/go-module/notify/feishu v0.0.0
	github.com/neteast-software/go-module/notify/feishu/linker v0.0.0
	github.com/neteast-software/go-module/observe/metrics v0.2.0
	github.com/neteast-software/go-module/observe/metrics/linker v0.2.0
	github.com/neteast-software/go-module/observe/metrics/prometheus/linker v0.0.0
	github.com/neteast-software/go-module/observe/metrics/rpc/grpc v0.2.0
	github.com/neteast-software/go-module/observe/tracing v0.2.0
	github.com/neteast-software/go-module/observe/tracing/linker v0.2.0
	github.com/neteast-software/go-module/observe/tracing/mq/consumer v0.2.0
	github.com/neteast-software/go-module/observe/tracing/opentelemetry v0.2.0
	github.com/neteast-software/go-module/observe/tracing/opentelemetry/linker v0.2.0
	github.com/neteast-software/go-module/observe/tracing/rpc/grpc v0.2.0
	github.com/neteast-software/go-module/outbox v0.2.0
	github.com/neteast-software/go-module/registry/service v0.2.0
	github.com/neteast-software/go-module/registry/service/nacos/linker v0.0.0
	github.com/neteast-software/go-module/rpc/grpc v0.2.0
	github.com/neteast-software/go-module/rpc/grpc/linker v0.2.0
	github.com/neteast-software/go-module/rpc/meta v0.2.0
	github.com/neteast-software/go-module/scheduler/cron v0.0.0
	github.com/neteast-software/go-module/scheduler/cron/linker v0.0.0
	github.com/neteast-software/go-module/token v0.2.0
	github.com/neteast-software/go-module/user/account v0.0.0
	github.com/neteast-software/grpc-discovery v0.1.0
	github.com/neteast-software/linker/v3 v3.1.0
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
	gorm.io/gorm v1.31.1
)

require (
	contrib.go.opencensus.io/exporter/ocagent v0.7.0 // indirect
	github.com/alibabacloud-go/alibabacloud-gateway-pop v0.0.6 // indirect
	github.com/alibabacloud-go/alibabacloud-gateway-spi v0.0.5 // indirect
	github.com/alibabacloud-go/darabonba-array v0.1.0 // indirect
	github.com/alibabacloud-go/darabonba-encode-util v0.0.2 // indirect
	github.com/alibabacloud-go/darabonba-map v0.0.2 // indirect
	github.com/alibabacloud-go/darabonba-openapi/v2 v2.0.10 // indirect
	github.com/alibabacloud-go/darabonba-signature-util v0.0.7 // indirect
	github.com/alibabacloud-go/darabonba-string v1.0.2 // indirect
	github.com/alibabacloud-go/debug v1.0.1 // indirect
	github.com/alibabacloud-go/endpoint-util v1.1.0 // indirect
	github.com/alibabacloud-go/kms-20160120/v3 v3.2.3 // indirect
	github.com/alibabacloud-go/openapi-util v0.1.0 // indirect
	github.com/alibabacloud-go/tea v1.2.2 // indirect
	github.com/alibabacloud-go/tea-utils v1.4.4 // indirect
	github.com/alibabacloud-go/tea-utils/v2 v2.0.7 // indirect
	github.com/alibabacloud-go/tea-xml v1.1.3 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1800 // indirect
	github.com/aliyun/alibabacloud-dkms-gcs-go-sdk v0.5.1 // indirect
	github.com/aliyun/alibabacloud-dkms-transfer-go-sdk v0.1.8 // indirect
	github.com/aliyun/aliyun-secretsmanager-client-go v1.1.5 // indirect
	github.com/aliyun/credentials-go v1.4.3 // indirect
	github.com/apache/rocketmq-clients/golang/v5 v5.1.4 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bsm/redislock v0.9.4 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clbanning/mxj/v2 v2.5.5 // indirect
	github.com/dchest/siphash v1.2.3 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/emmansun/gmsm v0.44.0 // indirect
	github.com/gin-gonic/gin v1.11.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.9.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nacos-group/nacos-sdk-go/v2 v2.3.5 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/neteast-software/go-module/audit/operate/store/gorm v0.2.0 // indirect
	github.com/neteast-software/go-module/config/env v0.0.0 // indirect
	github.com/neteast-software/go-module/config/yaml v0.0.0 // indirect
	github.com/neteast-software/go-module/crypto/sm v0.0.0 // indirect
	github.com/neteast-software/go-module/fault v0.2.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/fault/notice v0.2.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/http/gin v0.2.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/http/gin/linker v0.2.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/linker/server v0.2.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/mq/consumer v0.2.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/prometheus v0.0.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/scheduler/cron v0.0.0 // indirect
	github.com/neteast-software/go-module/observe/tracing/http/client v0.0.0 // indirect
	github.com/neteast-software/go-module/observe/tracing/http/gin v0.2.0 // indirect
	github.com/neteast-software/go-module/observe/tracing/http/gin/linker v0.2.0 // indirect
	github.com/neteast-software/go-module/observe/tracing/scheduler/cron v0.0.0 // indirect
	github.com/neteast-software/go-module/provider v0.2.0 // indirect
	github.com/neteast-software/go-module/redact v0.2.0 // indirect
	github.com/neteast-software/go-module/registry/service/linker v0.2.0 // indirect
	github.com/neteast-software/go-module/registry/service/nacos v0.0.0 // indirect
	github.com/neteast-software/go-module/scheduler/cron/store/gorm v0.0.0 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20210501183033-44dafcb38ecc // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/redis/go-redis/v9 v9.16.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/valyala/fastrand v1.1.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.opentelemetry.io/proto/otlp v1.10.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/api v0.230.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260401024825-9d38bb4040a9 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260401024825-9d38bb4040a9 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.6.0 // indirect
)

require (
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic v1.14.2 // indirect
	github.com/bytedance/sonic/loader v0.4.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/gabriel-vasile/mimetype v1.4.11 // indirect
	github.com/gin-contrib/gzip v1.2.5 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.28.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/neteast-software/go-module/graph/console v0.0.0
	github.com/neteast-software/go-module/http/gin/linker v0.2.0
	github.com/neteast-software/go-module/license v0.0.0
	github.com/neteast-software/go-module/license/http/gin v0.0.0
	github.com/neteast-software/go-module/registry/nacos v0.0.0
	github.com/neteast-software/go-module/registry/nacos/linker v0.0.0
	github.com/neteast-software/go-module/security/oauth v0.0.0
	github.com/neteast-software/go-module/security/oauth/http/gin v0.0.0
	github.com/neteast-software/go-module/security/oauth/jwt v0.0.0
	github.com/neteast-software/nacos-kit v0.1.0
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.59.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	golang.org/x/arch v0.22.0 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
)

replace github.com/neteast-software/go-module/http/client => ../modules/http/client

replace github.com/neteast-software/go-module/http/client/linker => ../modules/http/client/linker

replace github.com/neteast-software/go-module/acl/session => ../modules/acl/session

replace github.com/neteast-software/go-module/application/http/gin => ../modules/application/http/gin

replace github.com/neteast-software/go-module/db/gorm/query => ../modules/db/gorm/query

replace github.com/neteast-software/go-module/config/yaml/linker => ../modules/config/yaml/linker

replace github.com/neteast-software/go-module/config/env/linker => ../modules/config/env/linker

replace github.com/neteast-software/go-module/config/env => ../modules/config/env

replace github.com/neteast-software/go-module/crypto/sm => ../modules/crypto/sm

replace github.com/neteast-software/go-module/notify/feishu => ../modules/notify/feishu

replace github.com/neteast-software/go-module/notify/feishu/linker => ../modules/notify/feishu/linker

replace github.com/neteast-software/go-module/observe/metrics/prometheus => ../modules/observe/metrics/prometheus

replace github.com/neteast-software/go-module/observe/metrics/prometheus/linker => ../modules/observe/metrics/prometheus/linker

replace github.com/neteast-software/go-module/observe/metrics/scheduler/cron => ../modules/observe/metrics/scheduler/cron

replace github.com/neteast-software/go-module/observe/tracing/http/client => ../modules/observe/tracing/http/client

replace github.com/neteast-software/go-module/observe/tracing/scheduler/cron => ../modules/observe/tracing/scheduler/cron

replace github.com/neteast-software/go-module/config/yaml => ../modules/config/yaml

replace github.com/neteast-software/go-module/scheduler/cron => ../modules/scheduler/cron

replace github.com/neteast-software/go-module/scheduler/cron/linker => ../modules/scheduler/cron/linker

replace github.com/neteast-software/go-module/scheduler/cron/store/gorm => ../modules/scheduler/cron/store/gorm

replace github.com/neteast-software/go-module/user/account => ../modules/user/account

replace github.com/neteast-software/go-module/registry/nacos/linker => ../modules/registry/nacos/linker

replace github.com/neteast-software/nacos-kit => ../modules/nacos-kit

replace github.com/neteast-software/go-module/registry/nacos => ../modules/registry/nacos

replace github.com/neteast-software/go-module/registry/service/nacos => ../modules/registry/service/nacos

replace github.com/neteast-software/go-module/registry/service/nacos/linker => ../modules/registry/service/nacos/linker

replace github.com/neteast-software/go-module/license => ../modules/license

replace github.com/neteast-software/go-module/license/http/gin => ../modules/license/http/gin

replace github.com/neteast-software/go-module/graph/console => ../modules/graph/console

replace github.com/neteast-software/go-module/security/oauth => ../modules/security/oauth

replace github.com/neteast-software/go-module/security/oauth/http/gin => ../modules/security/oauth/http/gin

replace github.com/neteast-software/go-module/security/oauth/jwt => ../modules/security/oauth/jwt

replace github.com/neteast-software/linker/v3 => ../linker-v3
