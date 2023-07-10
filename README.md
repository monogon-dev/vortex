# 🌪️ Vortex

⚠️ _This repository is an **experiment and a work-in-progress**.
It is not an official Monogon product and/or ready for production use (yet?) - check back later!_

Vortex is a proof-of-concept Loki API implementation on top of ClickHouse. It translates LogQL queries to ClickHouse SQL, allowing the use of Loki frontends like Grafana with the performance and flexibility of ClickHouse.

Vortex ingests log streams from Kafka using native ClickHouse functionality. It comes with optimized schemas for the ElasticSearch Beats format, allowing the use with Filebeat, Metricbeat and others.  

---

_(C) 2023 Monogon SE. This software is provided "as-is" and without any express or implied warranties, including, without limitation, the implied warranties of merchantability and fitness for a particular purpose._
