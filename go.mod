module linker-v3-example

go 1.26.5

require (
	github.com/neteast-software/go-module/acl v0.3.0
	github.com/neteast-software/go-module/acl/session v0.4.0
	github.com/neteast-software/go-module/application v0.2.0
	github.com/neteast-software/go-module/application/http/gin v0.4.0
	github.com/neteast-software/go-module/application/linker v0.4.0
	github.com/neteast-software/go-module/application/store/gorm v0.4.0
	github.com/neteast-software/go-module/audit/operate v0.2.0
	github.com/neteast-software/go-module/audit/operate/http/gin v0.3.1
	github.com/neteast-software/go-module/audit/operate/linker v0.4.0
	github.com/neteast-software/go-module/cache/redis v0.3.0
	github.com/neteast-software/go-module/cache/redis/linker v0.4.0
	github.com/neteast-software/go-module/config v0.3.0
	github.com/neteast-software/go-module/config/env/linker v0.4.0
	github.com/neteast-software/go-module/config/yaml/linker v0.4.0
	github.com/neteast-software/go-module/db/gorm/model v0.1.0
	github.com/neteast-software/go-module/db/gorm/query v0.4.0
	github.com/neteast-software/go-module/db/gorm/table v0.3.1
	github.com/neteast-software/go-module/db/postgresql v0.4.0
	github.com/neteast-software/go-module/db/postgresql/linker v0.4.0
	github.com/neteast-software/go-module/fault/event v0.3.0
	github.com/neteast-software/go-module/fault/event/linker v0.4.0
	github.com/neteast-software/go-module/fault/notice v0.2.0
	github.com/neteast-software/go-module/fault/notice/linker v0.4.0
	github.com/neteast-software/go-module/http/client v0.4.0
	github.com/neteast-software/go-module/http/client/linker v0.4.0
	github.com/neteast-software/go-module/http/gateway v0.1.0
	github.com/neteast-software/go-module/http/gateway/declaration v0.1.0
	github.com/neteast-software/go-module/http/gateway/linker v0.2.0
	github.com/neteast-software/go-module/http/gin v0.5.0
	github.com/neteast-software/go-module/linker/server v0.4.1
	github.com/neteast-software/go-module/mq/consumer v0.2.0
	github.com/neteast-software/go-module/mq/consumer/linker v0.4.1
	github.com/neteast-software/go-module/mq/rocketmq v0.4.0
	github.com/neteast-software/go-module/mq/rocketmq/linker v0.4.0
	github.com/neteast-software/go-module/notify/feishu v0.1.0
	github.com/neteast-software/go-module/notify/feishu/linker v0.4.0
	github.com/neteast-software/go-module/observe/metrics v0.3.0
	github.com/neteast-software/go-module/observe/metrics/http/gateway v0.1.0
	github.com/neteast-software/go-module/observe/metrics/linker v0.4.0
	github.com/neteast-software/go-module/observe/metrics/prometheus/linker v0.4.0
	github.com/neteast-software/go-module/observe/metrics/rpc/grpc v0.3.1
	github.com/neteast-software/go-module/observe/tracing v0.3.0
	github.com/neteast-software/go-module/observe/tracing/http/gateway v0.1.0
	github.com/neteast-software/go-module/observe/tracing/linker v0.4.0
	github.com/neteast-software/go-module/observe/tracing/mq/consumer v0.4.1
	github.com/neteast-software/go-module/observe/tracing/opentelemetry v0.4.0
	github.com/neteast-software/go-module/observe/tracing/opentelemetry/linker v0.4.1
	github.com/neteast-software/go-module/observe/tracing/rpc/grpc v0.3.1
	github.com/neteast-software/go-module/outbox v0.2.0
	github.com/neteast-software/go-module/registry/discovery/nacos/linker v0.1.0
	github.com/neteast-software/go-module/registry/service v0.3.0
	github.com/neteast-software/go-module/registry/service/nacos/linker v0.4.0
	github.com/neteast-software/go-module/rpc/grpc v0.3.1
	github.com/neteast-software/go-module/rpc/grpc/linker v0.4.0
	github.com/neteast-software/go-module/rpc/meta v0.3.1
	github.com/neteast-software/go-module/scheduler/cron v0.4.1
	github.com/neteast-software/go-module/scheduler/cron/linker v0.4.1
	github.com/neteast-software/go-module/token v0.2.0
	github.com/neteast-software/go-module/user/account v0.4.0
	github.com/neteast-software/go-module/worker/periodic v0.1.0
	github.com/neteast-software/go-module/worker/periodic/linker v0.4.0
	github.com/neteast-software/linker/v3 v3.7.0
	google.golang.org/grpc v1.82.1
	google.golang.org/protobuf v1.36.12-0.20260120151049-f2248ac996af
	gorm.io/gorm v1.31.1
)

require (
	github.com/dgryski/go-jump v0.0.0-20211018200510-ba001c3ffce0 // indirect
	github.com/dgryski/go-mpchash v0.0.0-20200819201138-7382f34c4cd1 // indirect
	github.com/dimfeld/httppath v0.0.0-20170720192232-ee938bf73598 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/pprof v0.0.0-20260402051712-545e8a4df936 // indirect
	github.com/instana/go-sensor v1.73.5 // indirect
	github.com/lightstep/lightstep-tracer-common/golang/gogo v0.0.0-20210210170715-a8dfcb80d3a7 // indirect
	github.com/lightstep/lightstep-tracer-go v0.26.0 // indirect
	github.com/looplab/fsm v1.0.3 // indirect
	github.com/neteast-software/go-module/audit/operate/store/gorm v0.4.0 // indirect
	github.com/neteast-software/go-module/config/env v0.2.0 // indirect
	github.com/neteast-software/go-module/config/yaml v0.2.0 // indirect
	github.com/neteast-software/go-module/crypto/sm v0.1.0 // indirect
	github.com/neteast-software/go-module/fault v0.3.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/fault/notice v0.2.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/http/gin v0.2.1 // indirect
	github.com/neteast-software/go-module/observe/metrics/http/gin/linker v0.4.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/linker/server v0.4.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/mq/consumer v0.4.1 // indirect
	github.com/neteast-software/go-module/observe/metrics/prometheus v0.1.0 // indirect
	github.com/neteast-software/go-module/observe/metrics/scheduler/cron v0.1.0 // indirect
	github.com/neteast-software/go-module/observe/tracing/http/client v0.4.0 // indirect
	github.com/neteast-software/go-module/observe/tracing/http/gin v0.2.1 // indirect
	github.com/neteast-software/go-module/observe/tracing/http/gin/linker v0.4.0 // indirect
	github.com/neteast-software/go-module/observe/tracing/scheduler/cron v0.1.0 // indirect
	github.com/neteast-software/go-module/provider v0.3.0 // indirect
	github.com/neteast-software/go-module/redact v0.3.0 // indirect
	github.com/neteast-software/go-module/registry/discovery v0.1.0 // indirect
	github.com/neteast-software/go-module/registry/discovery/linker v0.1.0 // indirect
	github.com/neteast-software/go-module/registry/discovery/nacos v0.1.0 // indirect
	github.com/neteast-software/go-module/registry/service/linker v0.4.0 // indirect
	github.com/neteast-software/go-module/registry/service/nacos v0.4.0 // indirect
	github.com/neteast-software/go-module/scheduler/cron/store/gorm v0.4.0 // indirect
	github.com/neteast-software/grpc-discovery v0.1.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opentracing/basictracer-go v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20250401214520-65e299d6c5c9 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/sony/gobreaker v1.0.0 // indirect
	github.com/szuecs/rate-limit-buffer v0.9.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/valkey-io/valkey-go v1.0.76 // indirect
	github.com/valkey-io/valkey-go/valkeyhook v1.0.76 // indirect
	github.com/valkey-io/valkey-go/valkeyotel v1.0.76 // indirect
	github.com/yookoala/gofast v0.8.0 // indirect
	github.com/zalando/skipper v0.27.41 // indirect
	go.opentelemetry.io/otel/bridge/opentracing v1.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.44.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.44.0 // indirect
	go4.org/netipx v0.0.0-20231129151722-fdeea329fbba // indirect
	golang.org/x/exp v0.0.0-20251209150349-8475f28825e9 // indirect
	golang.org/x/tools/godoc v0.1.0-deprecated // indirect
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
	github.com/buger/jsonparser v1.2.0 // indirect
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
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.9.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.19.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nacos-group/nacos-sdk-go/v2 v2.3.5 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/orcaman/concurrent-map v0.0.0-20210501183033-44dafcb38ecc // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.24.1 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.1 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	github.com/redis/go-redis/v9 v9.21.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/valyala/fastrand v1.1.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/sdk v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	go.opentelemetry.io/proto/otlp v1.10.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	google.golang.org/api v0.256.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	gopkg.in/ini.v1 v1.67.1 // indirect
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
	github.com/neteast-software/go-module/graph/console v0.4.0
	github.com/neteast-software/go-module/http/gin/linker v0.4.0
	github.com/neteast-software/go-module/license v0.4.0
	github.com/neteast-software/go-module/license/http/gin v0.4.0
	github.com/neteast-software/go-module/registry/nacos v0.1.1
	github.com/neteast-software/go-module/registry/nacos/linker v0.4.0
	github.com/neteast-software/go-module/security/oauth v0.1.0
	github.com/neteast-software/go-module/security/oauth/http/gin v0.1.1
	github.com/neteast-software/go-module/security/oauth/jwt v0.1.0
	github.com/neteast-software/nacos-kit v0.1.0
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.59.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	golang.org/x/arch v0.22.0 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)
