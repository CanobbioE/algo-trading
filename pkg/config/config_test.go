package config_test

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/CanobbioE/algo-trading/pkg/config"
	"github.com/CanobbioE/algo-trading/pkg/monitor"
	"github.com/CanobbioE/algo-trading/pkg/strategies"
	"github.com/CanobbioE/algo-trading/pkg/utilities"
)

func substringErrMatcher(substr string) func(err error) (string, bool) {
	return func(err error) (string, bool) {
		if strings.Contains(err.Error(), substr) {
			return "", true
		}
		return fmt.Sprintf(
			"expected error to contain \"%q\", instead got: \"%v\"",
			substr, err), false
	}
}

func TestConfig_UnmarshalJSON(t *testing.T) {
	type output struct {
		wantErr    bool
		errMatcher func(error) (string, bool)
		cfg        *config.Config
	}
	type testCase struct {
		name string
		data []byte
		want *output
	}

	for _, tc := range []testCase{
		{
			name: "succeeds with valid config",
			data: utilities.MustReturn(os.ReadFile("./testdata/full-config.json")),
			want: &output{
				cfg: &config.Config{
					StockUniverse: []string{"GME"},
					Strategies: []*strategies.StrategyWeight{
						{
							Strategy: &strategies.BreakoutStrategy{},
							Weight:   1.8,
						},
						{
							Strategy: &strategies.VWAPStrategy{},
							Weight:   1.0,
						},
						{
							Strategy: &strategies.MeanReversionStrategy{},
							Weight:   0.6,
						},
						{
							Strategy: &strategies.BollingerBandSqueezeStrategy{},
							Weight:   1.2,
						},
						{
							Strategy: &strategies.MomentumStrategy{},
							Weight:   1.8,
						},
						{
							Strategy: &strategies.MACDStrategy{},
							Weight:   3.0,
						},
					},
					Thresholds: &strategies.Thresholds{
						AtrPeriod:         3,
						LowATRThreshold:   0.07,
						HighATRThreshold:  0.4,
						LowLookback:       8,
						HighLookback:      3,
						VolumeThreshold:   1.0,
						Deviation:         0.03,
						Squeeze:           0.07,
						MinMomentumReturn: 0.03,
					},
					MACDParams: &strategies.MACDParams{
						FastPeriod:      6,
						SlowPeriod:      13,
						SignalPeriod:    5,
						TriggerDistance: 0.005,
					},
					LookBack:             3,
					MomentumLookBack:     8,
					BollingerCoefficient: 1.8,
					Filters: &monitor.ScanFilters{
						MinConfidence:    0.4,
						MinWeightedScore: 0.4,
						MaxRisk:          monitor.RiskHigh,
						MinOpportunity:   monitor.OpportunityHigh,
						MinVolume:        1000,
						RequiredSignals:  1,
					},
				},
			},
		},
		{
			name: "succeeds with minimal config",
			data: utilities.MustReturn(os.ReadFile("./testdata/minimal-config.json")),
			want: &output{
				cfg: &config.Config{
					StockUniverse: []string{"GME"},
					Strategies: []*strategies.StrategyWeight{
						{
							Strategy: &strategies.BreakoutStrategy{},
							Weight:   1.8,
						},
					},
					Thresholds: &strategies.Thresholds{
						AtrPeriod:         3,
						LowATRThreshold:   0.07,
						HighATRThreshold:  0.4,
						LowLookback:       8,
						HighLookback:      3,
						VolumeThreshold:   1.0,
						Deviation:         0.03,
						Squeeze:           0.07,
						MinMomentumReturn: 0.03,
					},
				},
			},
		},
		{
			name: "fails with invalid json",
			data: []byte(`{ "broken":`),
			want: &output{
				wantErr:    true,
				errMatcher: substringErrMatcher("failed to parse json raw"),
			},
		},
		{
			name: "fails with missing threshold",
			data: []byte(`{}`),
			want: &output{
				wantErr:    true,
				errMatcher: substringErrMatcher("no thresholds specified"),
			},
		},
		{
			name: "fails with missing strategies",
			data: []byte(`{"Thresholds": {}}`),
			want: &output{
				wantErr:    true,
				errMatcher: substringErrMatcher("at least one strategy must be specified"),
			},
		},
		{
			name: "fails with unknown strategies",
			data: []byte(`{"Thresholds": {}, "Strategies": [{"strategy": "what's this?'"}]}`),
			want: &output{
				wantErr:    true,
				errMatcher: substringErrMatcher("unknown strategy"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := &config.Config{}
			err := got.UnmarshalJSON(tc.data)

			switch {
			case err == nil && !tc.want.wantErr:
				// continue as usual
			case err != nil && !tc.want.wantErr:
				t.Fatalf("expected no error, instead got: %v", err)
			case err == nil && tc.want.wantErr:
				t.Fatalf("expected error, instead got: %v", err)
			case tc.want.errMatcher != nil:
				if msg, ok := tc.want.errMatcher(err); !ok {
					t.Fatal(msg)
				}
				return
			default:
				t.Fatalf("unexpected error, instead got: %v", err)
			}

			compareStrategies := cmp.Comparer(func(x, y strategies.Strategy) bool {
				return reflect.TypeOf(x) == reflect.TypeOf(y)
			})

			isStrategyField := func(path cmp.Path) bool {
				ps := path.String()
				return ps == "Strategies.Strategy" || path.Last().Type() == reflect.TypeOf((*strategies.Strategy)(nil)).Elem()
			}

			diff := cmp.Diff(tc.want.cfg, got,
				cmp.AllowUnexported(
					config.Config{},
					monitor.ScanFilters{},
					strategies.Thresholds{},
					strategies.MACDParams{},
					strategies.BreakoutStrategy{},
					strategies.VWAPStrategy{},
					strategies.MeanReversionStrategy{},
					strategies.BollingerBandSqueezeStrategy{},
					strategies.MomentumStrategy{},
				), cmp.FilterPath(isStrategyField, compareStrategies))

			if diff != "" {
				t.Errorf("unmarshal json result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
