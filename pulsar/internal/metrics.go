// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package internal

import (
	"github.com/prometheus/client_golang/prometheus"
)

var topicLabelNames = []string{"pulsar_tenant", "pulsar_namespace", "topic"}

type Metrics struct {
	messagesPublished *prometheus.CounterVec
	bytesPublished    *prometheus.CounterVec
	messagesPending   *prometheus.GaugeVec
	bytesPending      *prometheus.GaugeVec
	publishErrors     *prometheus.CounterVec
	publishLatency    *prometheus.HistogramVec
	publishRPCLatency *prometheus.HistogramVec

	messagesReceived   *prometheus.CounterVec
	bytesReceived      *prometheus.CounterVec
	prefetchedMessages *prometheus.GaugeVec
	prefetchedBytes    *prometheus.GaugeVec
	acksCounter        *prometheus.CounterVec
	nacksCounter       *prometheus.CounterVec
	dlqCounter         *prometheus.CounterVec
	processingTime     *prometheus.HistogramVec

	producersOpened     *prometheus.CounterVec
	producersClosed     *prometheus.CounterVec
	producersPartitions *prometheus.GaugeVec
	consumersOpened     *prometheus.CounterVec
	consumersClosed     *prometheus.CounterVec
	consumersPartitions *prometheus.GaugeVec
	readersOpened       *prometheus.CounterVec
	readersClosed       *prometheus.CounterVec

	// Metrics that are not labeled with topic, are immediately available
	ConnectionsOpened                     prometheus.Counter
	ConnectionsClosed                     prometheus.Counter
	ConnectionsEstablishmentErrors        prometheus.Counter
	ConnectionsHandshakeErrors            prometheus.Counter
	LookupRequestsCount                   prometheus.Counter
	PartitionedTopicMetadataRequestsCount prometheus.Counter
	RPCRequestCount                       prometheus.Counter
}

type TopicMetrics struct {
	MessagesPublished        prometheus.Counter
	BytesPublished           prometheus.Counter
	MessagesPending          prometheus.Gauge
	BytesPending             prometheus.Gauge
	PublishErrorsTimeout     prometheus.Counter
	PublishErrorsMsgTooLarge prometheus.Counter
	PublishLatency           prometheus.Observer
	PublishRPCLatency        prometheus.Observer

	MessagesReceived   prometheus.Counter
	BytesReceived      prometheus.Counter
	PrefetchedMessages prometheus.Gauge
	PrefetchedBytes    prometheus.Gauge
	AcksCounter        prometheus.Counter
	NacksCounter       prometheus.Counter
	DlqCounter         prometheus.Counter
	ProcessingTime     prometheus.Observer

	ProducersOpened     prometheus.Counter
	ProducersClosed     prometheus.Counter
	ProducersPartitions prometheus.Gauge
	ConsumersOpened     prometheus.Counter
	ConsumersClosed     prometheus.Counter
	ConsumersPartitions prometheus.Gauge
	ReadersOpened       prometheus.Counter
	ReadersClosed       prometheus.Counter
}

func NewMetricsProvider(userDefinedLabels map[string]string) *Metrics {
	constLabels := map[string]string{
		"client": "go",
	}
	for k, v := range userDefinedLabels {
		constLabels[k] = v
	}

	metrics := &Metrics{
		messagesPublished: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_messages_published",
			Help:        "Counter of messages published by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		bytesPublished: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_bytes_published",
			Help:        "Counter of messages published by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		messagesPending: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:        "pulsar_client_producer_pending_messages",
			Help:        "Counter of messages pending to be published by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		bytesPending: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:        "pulsar_client_producer_pending_bytes",
			Help:        "Counter of bytes pending to be published by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		publishErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_producer_errors",
			Help:        "Counter of publish errors",
			ConstLabels: constLabels,
		}, append(topicLabelNames, "error")),

		publishLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        "pulsar_client_producer_latency_seconds",
			Help:        "Publish latency experienced by the client",
			ConstLabels: constLabels,
			Buckets:     []float64{.0005, .001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		}, topicLabelNames),

		publishRPCLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        "pulsar_client_producer_rpc_latency_seconds",
			Help:        "Publish RPC latency experienced internally by the client when sending data to receiving an ack",
			ConstLabels: constLabels,
			Buckets:     []float64{.0005, .001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		}, topicLabelNames),

		producersOpened: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_producers_opened",
			Help:        "Counter of producers created by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		producersClosed: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_producers_closed",
			Help:        "Counter of producers closed by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		producersPartitions: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:        "pulsar_client_producers_partitions_active",
			Help:        "Counter of individual partitions the producers are currently active",
			ConstLabels: constLabels,
		}, topicLabelNames),

		consumersOpened: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_consumers_opened",
			Help:        "Counter of consumers created by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		consumersClosed: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_consumers_closed",
			Help:        "Counter of consumers closed by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		consumersPartitions: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:        "pulsar_client_consumers_partitions_active",
			Help:        "Counter of individual partitions the consumers are currently active",
			ConstLabels: constLabels,
		}, topicLabelNames),

		messagesReceived: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_messages_received",
			Help:        "Counter of messages received by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		bytesReceived: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_bytes_received",
			Help:        "Counter of bytes received by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		prefetchedMessages: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:        "pulsar_client_consumer_prefetched_messages",
			Help:        "Number of messages currently sitting in the consumer pre-fetch queue",
			ConstLabels: constLabels,
		}, topicLabelNames),

		prefetchedBytes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:        "pulsar_client_consumer_prefetched_bytes",
			Help:        "Total number of bytes currently sitting in the consumer pre-fetch queue",
			ConstLabels: constLabels,
		}, topicLabelNames),

		acksCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_consumer_acks",
			Help:        "Counter of messages acked by client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		nacksCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_consumer_nacks",
			Help:        "Counter of messages nacked by client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		dlqCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_consumer_dlq_messages",
			Help:        "Counter of messages sent to Dead letter queue",
			ConstLabels: constLabels,
		}, topicLabelNames),

		processingTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        "pulsar_client_consumer_processing_time_seconds",
			Help:        "Time it takes for application to process messages",
			Buckets:     []float64{.0005, .001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			ConstLabels: constLabels,
		}, topicLabelNames),

		readersOpened: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_readers_opened",
			Help:        "Counter of readers created by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		readersClosed: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "pulsar_client_readers_closed",
			Help:        "Counter of readers closed by the client",
			ConstLabels: constLabels,
		}, topicLabelNames),

		ConnectionsOpened: prometheus.NewCounter(prometheus.CounterOpts{
			Name:        "pulsar_client_connections_opened",
			Help:        "Counter of connections created by the client",
			ConstLabels: constLabels,
		}),

		ConnectionsClosed: prometheus.NewCounter(prometheus.CounterOpts{
			Name:        "pulsar_client_connections_closed",
			Help:        "Counter of connections closed by the client",
			ConstLabels: constLabels,
		}),

		ConnectionsEstablishmentErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name:        "pulsar_client_connections_establishment_errors",
			Help:        "Counter of errors in connections establishment",
			ConstLabels: constLabels,
		}),

		ConnectionsHandshakeErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name:        "pulsar_client_connections_handshake_errors",
			Help:        "Counter of errors in connections handshake (eg: authz)",
			ConstLabels: constLabels,
		}),

		LookupRequestsCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name:        "pulsar_client_lookup_count",
			Help:        "Counter of lookup requests made by the client",
			ConstLabels: constLabels,
		}),

		PartitionedTopicMetadataRequestsCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name:        "pulsar_client_partitioned_topic_metadata_count",
			Help:        "Counter of partitioned_topic_metadata requests made by the client",
			ConstLabels: constLabels,
		}),

		RPCRequestCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name:        "pulsar_client_rpc_count",
			Help:        "Counter of RPC requests made by the client",
			ConstLabels: constLabels,
		}),
	}

	err := prometheus.DefaultRegisterer.Register(metrics.messagesPublished)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.messagesPublished = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.bytesPublished)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.bytesPublished = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.messagesPending)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.messagesPending = are.ExistingCollector.(*prometheus.GaugeVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.bytesPending)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.bytesPending = are.ExistingCollector.(*prometheus.GaugeVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.publishErrors)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.publishErrors = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.publishLatency)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.publishLatency = are.ExistingCollector.(*prometheus.HistogramVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.publishRPCLatency)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.publishRPCLatency = are.ExistingCollector.(*prometheus.HistogramVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.messagesReceived)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.messagesReceived = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.bytesReceived)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.bytesReceived = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.prefetchedMessages)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.prefetchedMessages = are.ExistingCollector.(*prometheus.GaugeVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.prefetchedBytes)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.prefetchedBytes = are.ExistingCollector.(*prometheus.GaugeVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.acksCounter)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.acksCounter = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.nacksCounter)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.nacksCounter = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.dlqCounter)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.dlqCounter = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.processingTime)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.processingTime = are.ExistingCollector.(*prometheus.HistogramVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.producersOpened)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.producersOpened = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.producersClosed)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.producersClosed = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.producersPartitions)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.producersPartitions = are.ExistingCollector.(*prometheus.GaugeVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.consumersOpened)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.consumersOpened = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.consumersClosed)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.consumersClosed = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.consumersPartitions)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.consumersPartitions = are.ExistingCollector.(*prometheus.GaugeVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.readersOpened)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.readersOpened = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.readersClosed)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.readersClosed = are.ExistingCollector.(*prometheus.CounterVec)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.ConnectionsOpened)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.ConnectionsOpened = are.ExistingCollector.(prometheus.Counter)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.ConnectionsClosed)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.ConnectionsClosed = are.ExistingCollector.(prometheus.Counter)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.ConnectionsEstablishmentErrors)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.ConnectionsEstablishmentErrors = are.ExistingCollector.(prometheus.Counter)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.ConnectionsHandshakeErrors)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.ConnectionsHandshakeErrors = are.ExistingCollector.(prometheus.Counter)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.LookupRequestsCount)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.LookupRequestsCount = are.ExistingCollector.(prometheus.Counter)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.PartitionedTopicMetadataRequestsCount)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.PartitionedTopicMetadataRequestsCount = are.ExistingCollector.(prometheus.Counter)
		}
	}
	err = prometheus.DefaultRegisterer.Register(metrics.RPCRequestCount)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			metrics.RPCRequestCount = are.ExistingCollector.(prometheus.Counter)
		}
	}
	return metrics
}

func (mp *Metrics) GetTopicMetrics(t string) *TopicMetrics {
	tn, _ := ParseTopicName(t)
	topic := TopicNameWithoutPartitionPart(tn)
	labels := map[string]string{
		"pulsar_tenant":    tn.Tenant,
		"pulsar_namespace": tn.Namespace,
		"topic":            topic,
	}

	tm := &TopicMetrics{
		MessagesPublished:        mp.messagesPublished.With(labels),
		BytesPublished:           mp.bytesPublished.With(labels),
		MessagesPending:          mp.messagesPending.With(labels),
		BytesPending:             mp.bytesPending.With(labels),
		PublishErrorsTimeout:     mp.publishErrors.With(mergeMaps(labels, map[string]string{"error": "timeout"})),
		PublishErrorsMsgTooLarge: mp.publishErrors.With(mergeMaps(labels, map[string]string{"error": "msg_too_large"})),
		PublishLatency:           mp.publishLatency.With(labels),
		PublishRPCLatency:        mp.publishRPCLatency.With(labels),

		MessagesReceived:   mp.messagesReceived.With(labels),
		BytesReceived:      mp.bytesReceived.With(labels),
		PrefetchedMessages: mp.prefetchedMessages.With(labels),
		PrefetchedBytes:    mp.prefetchedBytes.With(labels),
		AcksCounter:        mp.acksCounter.With(labels),
		NacksCounter:       mp.nacksCounter.With(labels),
		DlqCounter:         mp.dlqCounter.With(labels),
		ProcessingTime:     mp.processingTime.With(labels),

		ProducersOpened:     mp.producersOpened.With(labels),
		ProducersClosed:     mp.producersClosed.With(labels),
		ProducersPartitions: mp.producersPartitions.With(labels),
		ConsumersOpened:     mp.consumersOpened.With(labels),
		ConsumersClosed:     mp.consumersClosed.With(labels),
		ConsumersPartitions: mp.consumersPartitions.With(labels),
		ReadersOpened:       mp.readersOpened.With(labels),
		ReadersClosed:       mp.readersClosed.With(labels),
	}

	return tm
}

func mergeMaps(a, b map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range a {
		res[k] = v
	}
	for k, v := range b {
		res[k] = v
	}

	return res
}
