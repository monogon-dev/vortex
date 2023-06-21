CREATE OR REPLACE TABLE filebeat_kafka
(
   `line` String
) ENGINE = Kafka()
      SETTINGS
          kafka_broker_list = 'kafka:9092',
          kafka_topic_list = 'filebeat',
          kafka_group_name = 'clickhouse',
          kafka_format = 'JSONAsString',
          kafka_max_block_size = 1000000,
          kafka_num_consumers = 1,
          kafka_thread_per_consumer = 1,
          input_format_import_nested_json = 1;

CREATE OR REPLACE TABLE filebeat
(
    `@timestamp` DateTime64(3) CODEC(Delta(8), ZSTD(1)),
    `message` String CODEC(ZSTD(1)),
    `labels` Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    INDEX idx_fields_key mapKeys(`labels`) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_fields_value mapValues(`labels`) TYPE bloom_filter(0.01) GRANULARITY 1
)
    ENGINE = MergeTree
        PARTITION BY toDate(`@timestamp`)
        ORDER BY (toUnixTimestamp(`@timestamp`))
        SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;

DROP VIEW filebeat_view;
-- generate based on config?
CREATE MATERIALIZED VIEW filebeat_view TO filebeat
AS
SELECT `@timestamp`,
       `message`,
       mapFilter((k,v) -> v != '', mapConcat(
               mapFromArrays(
                       ['logfile', 'hostname', 'input', 'agent', 'tags'],
                       [logfile, hostname, input, agent, tags]
                   ),
               `labels`
           )) as `labels`
FROM (
         SELECT parseDateTime64BestEffort(JSON_VALUE(line, '$."@timestamp"'))       AS `@timestamp`,
                JSON_VALUE(line, '$.message')          AS `message`,
                JSON_VALUE(line, '$.log.file.path')    AS `logfile`,
                JSON_VALUE(line, '$.host.name')        AS `hostname`,
                JSON_VALUE(line, '$.input.type')       AS `input`,
                JSON_VALUE(line, '$.agent.type')       AS `agent`,
                arrayStringConcat(JSONExtract(JSON_VALUE(line, '$.tags'), 'Array(String)'),',') AS `tags`,
                JSONExtract(JSON_VALUE(line, '$.labels'), 'Map(String, String)') AS `labels`
                 FROM filebeat_kafka
                 SETTINGS
                 function_json_value_return_type_allow_complex= true
         );



SELECT *
from filebeat;



TRUNCATE filebeat;



WITH '{
  "@timestamp": "2023-06-21T14:49:16.828Z",
  "@metadata": {
    "beat": "filebeat",
    "type": "_doc",
    "version": "8.8.1"
  },
  "message": "2023.06.21 14:49:13.398571 [ 272 ] {} <Trace> MergedBlockOutputStream: filled checksums 202306_57_57_0 (state Temporary)",
  "input": {
    "type": "filestream"
  },
  "agent": {
    "id": "e742104a-5f35-4685-a054-34a1ecf61836",
    "name": "e66f13c97673",
    "type": "filebeat",
    "version": "8.8.1",
    "ephemeral_id": "81377a64-f212-4980-b316-f0df2206d42f"
  },
  "labels": {
    "application": "foobar"
  },
  "tags": [
    "production",
    "env2"
  ],
  "ecs": {
    "version": "8.0.0"
  },
  "host": {
    "name": "e66f13c97673"
  },
  "log": {
    "offset": 385774,
    "file": {
      "path": "/var/log/clickhouse.log"
    }
  }
}' AS inputData
SELECT `@timestamp`,
       `message`,
       mapFilter((k,v) -> v != '', mapConcat(
               mapFromArrays(
                       ['logfile', 'hostname', 'input', 'agent', 'tags'],
                       [logfile, hostname, input, agent, tags]
                   ),
               `labels`
           )) as `labels`
FROM (
         SELECT JSON_VALUE(inputData, '$."@timestamp"')       AS `@timestamp`,
                JSON_VALUE(inputData, '$.message')          AS `message`,
                JSON_VALUE(inputData, '$.log.file.path')    AS `logfile`,
                JSON_VALUE(inputData, '$.host.name')        AS `hostname`,
                JSON_VALUE(inputData, '$.input.type')       AS `input`,
                JSON_VALUE(inputData, '$.agent.type')       AS `agent`,
                arrayStringConcat(JSONExtract(JSON_VALUE(inputData, '$.tags'), 'Array(String)'),',') AS `tags`,
                JSONExtract(JSON_VALUE(inputData, '$.labels'), 'Map(String, String)') AS `labels`
            SETTINGS
                 function_json_value_return_type_allow_complex= true
         );