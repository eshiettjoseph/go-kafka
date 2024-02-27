## Weather Data Collection with Golang and Apache Kafka
A webscraper set up as a producer that scrapers weather data hourly from https://wttr.in publishes to Kafka. A consumer subscribes to this data and stores it in a PostgreSQL database.
**TODO:**
- Deploy to EKS
- Setup GitOps with ArgoCD
- Setup IPFS for data storage
- CICD pipeline to build and deploy images to Dockerhub
- Loadbalancing
- Monitoring, Logging, tracing with OpenTelemetry
