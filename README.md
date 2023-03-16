# Log Benchmarks

> Benchmarks on using logging, tracing and without them.

```bash
# setup jaeger
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.43
```

| Case                        | Data Size | Total loop | Average elapsed time |
| --------------------------- | :-------: | ---------: | -------------------: |
| Without Logging and Tracing |   small   |      15217 |          77804 ns/op |
| Without Logging and Tracing |    big    |       3790 |         267474 ns/op |
| Tracing                     |   small   |      13303 |          90404 ns/op |
| Tracing                     |    big    |       3592 |         308865 ns/op |
| With Logging                |   small   |       5738 |         240336 ns/op |
| With Logging                |    big    |       1150 |        1689071 ns/op |
