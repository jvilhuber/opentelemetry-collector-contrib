// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memcachedreceiver

//go:generate mdatagen metadata.yaml

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

const (
	typeStr = "memcached"
)

// NewFactory creates a factory for memcached receiver.
func NewFactory() component.ReceiverFactory {
	return receiverhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		receiverhelper.WithMetrics(createMetricsReceiver))
}

func createDefaultConfig() configmodels.Receiver {
	return &config{
		ReceiverSettings: configmodels.ReceiverSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		ScraperControllerSettings: receiverhelper.ScraperControllerSettings{
			CollectionInterval: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
		TCPAddr: confignet.TCPAddr{
			Endpoint: "localhost:11211",
		},
	}
}

func createMetricsReceiver(
	ctx context.Context,
	params component.ReceiverCreateParams,
	rConf configmodels.Receiver,
	consumer consumer.MetricsConsumer,
) (component.MetricsReceiver, error) {
	cfg := rConf.(*config)

	scraper := newMemcachedScraper(params.Logger, cfg)

	return receiverhelper.NewScraperControllerReceiver(
		&cfg.ScraperControllerSettings, consumer,
		receiverhelper.AddResourceMetricsScraper(scraper),
	)
}
