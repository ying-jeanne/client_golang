// Copyright 2019 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package v1_test provides examples making requests to Prometheus using the
// Golang client.
package v2_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/common/config"

	"github.com/prometheus/client_golang/api"
	v2 "github.com/prometheus/client_golang/exp/api/client/v2"
)

const DemoPrometheusURL = "https://demo.prometheus.io:443"

func ExampleAPI_query() {
	client, err := api.NewClient(api.Config{
		Address: DemoPrometheusURL,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v2api := v2.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, annotations, err := v2api.Query(ctx, "up", time.Now(), v2.WithTimeout(5*time.Second))
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(annotations.Warnings) > 0 {
		fmt.Printf("Warnings: %v\n", annotations.Warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func ExampleAPI_queryRange() {
	client, err := api.NewClient(api.Config{
		Address: DemoPrometheusURL,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v2api := v2.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v2.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
	result, annotations, err := v2api.QueryRange(ctx, "rate(prometheus_tsdb_head_samples_appended_total[5m])", r, v2.WithTimeout(5*time.Second))
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(annotations.Warnings) > 0 {
		fmt.Printf("Warnings: %v\n", annotations.Warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

type userAgentRoundTripper struct {
	name string
	rt   http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface.
func (u userAgentRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.UserAgent() == "" {
		// The specification of http.RoundTripper says that it shouldn't mutate
		// the request so make a copy of req.Header since this is all that is
		// modified.
		r2 := new(http.Request)
		*r2 = *r
		r2.Header = make(http.Header)
		for k, s := range r.Header {
			r2.Header[k] = s
		}
		r2.Header.Set("User-Agent", u.name)
		r = r2
	}
	return u.rt.RoundTrip(r)
}

func ExampleAPI_queryRangeWithUserAgent() {
	client, err := api.NewClient(api.Config{
		Address:      DemoPrometheusURL,
		RoundTripper: userAgentRoundTripper{name: "Client-Golang", rt: api.DefaultRoundTripper},
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v2api := v2.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v2.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
	result, annotations, err := v2api.QueryRange(ctx, "rate(prometheus_tsdb_head_samples_appended_total[5m])", r)
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(annotations.Warnings) > 0 {
		fmt.Printf("Warnings: %v\n", annotations.Warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func ExampleAPI_queryRangeWithBasicAuth() {
	client, err := api.NewClient(api.Config{
		Address: DemoPrometheusURL,
		// We can use amazing github.com/prometheus/common/config helper!
		RoundTripper: config.NewBasicAuthRoundTripper(
			config.NewInlineSecret("me"),
			config.NewInlineSecret("definitely_me"),
			api.DefaultRoundTripper,
		),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v2api := v2.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v2.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
	result, annotations, err := v2api.QueryRange(ctx, "rate(prometheus_tsdb_head_samples_appended_total[5m])", r)
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(annotations.Warnings) > 0 {
		fmt.Printf("Warnings: %v\n", annotations.Warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func ExampleAPI_queryRangeWithAuthBearerToken() {
	client, err := api.NewClient(api.Config{
		Address: DemoPrometheusURL,
		// We can use amazing github.com/prometheus/common/config helper!
		RoundTripper: config.NewAuthorizationCredentialsRoundTripper(
			"Bearer",
			config.NewInlineSecret("secret_token"),
			api.DefaultRoundTripper,
		),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v2api := v2.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v2.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
	result, annotations, err := v2api.QueryRange(ctx, "rate(prometheus_tsdb_head_samples_appended_total[5m])", r)
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(annotations.Warnings) > 0 {
		fmt.Printf("Warnings: %v\n", annotations.Warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func ExampleAPI_queryRangeWithAuthBearerTokenHeadersRoundTripper() {
	client, err := api.NewClient(api.Config{
		Address: DemoPrometheusURL,
		// We can use amazing github.com/prometheus/common/config helper!
		RoundTripper: config.NewHeadersRoundTripper(
			&config.Headers{
				Headers: map[string]config.Header{
					"Authorization": {
						Values: []string{"Bearer secret"},
					},
				},
			},
			api.DefaultRoundTripper,
		),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v2api := v2.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v2.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
	result, annotations, err := v2api.QueryRange(ctx, "rate(prometheus_tsdb_head_samples_appended_total[5m])", r)
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(annotations.Warnings) > 0 {
		fmt.Printf("Warnings: %v\n", annotations.Warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func ExampleAPI_series() {
	client, err := api.NewClient(api.Config{
		Address: DemoPrometheusURL,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v2api := v2.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	lbls, annotations, err := v2api.Series(ctx, []string{
		"{__name__=~\"scrape_.+\",job=\"node\"}",
		"{__name__=~\"scrape_.+\",job=\"prometheus\"}",
	}, time.Now().Add(-time.Hour), time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(annotations.Warnings) > 0 {
		fmt.Printf("Warnings: %v\n", annotations.Warnings)
	}
	fmt.Println("Result:")
	for _, lbl := range lbls {
		fmt.Println(lbl)
	}
}
