# Vortex
| ⚠️         **PRE-RELEASE**: This is a work in progress - please watch this repo for news. |
|-------------------------------------------------------------------------------------------|

🌪️ Vortex is a powerful proof of concept tool designed to emulate the Loki query API and seamlessly convert LogQL queries to ClickHouse SQL.

## ✨ Features
Vortex offers the following key features:

- LogQL to ClickHouse SQL Conversion: By seamlessly converting LogQL queries to ClickHouse SQL, Vortex enables you to leverage the power of ClickHouse for log analysis, enhancing performance and scalability.
- Filebeat Integration: Vortex integrates seamlessly with Filebeat, a reliable log shipper, enabling easy log ingestion and efficient transfer of logs to Kafka.
- Kafka Integration: Vortex utilizes Kafka as the intermediary transport layer, facilitating smooth and real-time log streaming for further processing.
- JSON Parsing with Materialized Views: Vortex leverages ClickHouse's Materialized Views to parse JSON logs efficiently, extracting meaningful information such as timestamp, message, and labels.
- Flexible Log Schema: Vortex's log schema comprises three columns: timestamp, message, and a Map(String, String) for labels, offering a comprehensive representation of log data while ensuring flexibility and extensibility.

## 🚀 Getting Started 
To get started with Vortex, follow these simple steps:

- Install and configure Filebeat to send logs to Kafka. Make sure the logs are in a compatible JSON format.
- Set up Kafka and create a topic where Filebeat can send the logs.
- Install and configure ClickHouse, ensuring that it is properly connected to Kafka.
- Start Vortex and ensure that it is properly configured to interact with ClickHouse and Kafka.

---
(C) 2023 Monogon SE.

This software is provided "as-is" and without any express or implied warranties, including, without limitation, the implied warranties of merchantability and fitness for a particular purpose.