// Copyright 2017 Google Inc.
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
//

package stats

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/instrumentation-go/stats/tagging"
)

// GaugeStringViewDesc defines an string gauge view.
type GaugeStringViewDesc struct {
	*ViewDescCommon
}

func (gd *GaugeStringViewDesc) createAggregator(t time.Time) (aggregator, error) {
	return newGaugeAggregatorString(), nil
}

func (gd *GaugeStringViewDesc) retrieveView(now time.Time) (*View, error) {
	gav, err := gd.retrieveAggreationView(now)
	if err != nil {
		return nil, err
	}
	return &View{
		ViewDesc: gd,
		ViewAgg:  gav,
	}, nil
}

func (gd *GaugeStringViewDesc) viewDesc() *ViewDescCommon {
	return gd.ViewDescCommon
}

func (gd *GaugeStringViewDesc) isValid() error {
	return nil
}

func (gd *GaugeStringViewDesc) retrieveAggreationView(t time.Time) (*GaugeStringAggView, error) {
	var aggs []*GaugeStringAgg

	for sig, a := range gd.signatures {
		tags, err := tagging.TagsFromValuesSignature([]byte(sig), gd.TagKeys)
		if err != nil {
			return nil, fmt.Errorf("malformed signature '%v'. %v", sig, err)
		}
		aggregator, ok := a.(*gaugeAggregatorString)
		if !ok {
			return nil, fmt.Errorf("unexpected aggregator type. got %T, want stats.gaugeAggregatorString", a)
		}
		ga := &GaugeStringAgg{
			GaugeStringStats: aggregator.retrieveCollected(),
			Tags:             tags,
		}
		aggs = append(aggs, ga)
	}

	return &GaugeStringAggView{
		Descriptor:   gd,
		Aggregations: aggs,
	}, nil
}

func (gd *GaugeStringViewDesc) stringWithIndent(tabs string) string {
	if gd == nil {
		return "nil"
	}
	vd := gd.ViewDescCommon
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%T {\n", gd)
	fmt.Fprintf(&buf, "%v  Name: %v,\n", tabs, vd.Name)
	fmt.Fprintf(&buf, "%v  Description: %v,\n", tabs, vd.Description)
	fmt.Fprintf(&buf, "%v  MeasureDescName: %v,\n", tabs, vd.MeasureDescName)
	fmt.Fprintf(&buf, "%v  TagKeys: %v,\n", tabs, vd.TagKeys)
	fmt.Fprintf(&buf, "%v}", tabs)
	return buf.String()
}

func (gd *GaugeStringViewDesc) String() string {
	return gd.stringWithIndent("")
}

// GaugeStringAggView is the set of collected GaugeStringAgg associated with
// ViewDesc.
type GaugeStringAggView struct {
	Descriptor   *GaugeStringViewDesc
	Aggregations []*GaugeStringAgg
}

func (gv *GaugeStringAggView) stringWithIndent(tabs string) string {
	if gv == nil {
		return "nil"
	}

	tabs2 := tabs + "    "
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%T {\n", gv)
	fmt.Fprintf(&buf, "%v  Aggregations:\n", tabs)
	for _, agg := range gv.Aggregations {
		fmt.Fprintf(&buf, "%v%v,\n", tabs2, agg.stringWithIndent(tabs2))
	}
	fmt.Fprintf(&buf, "%v}", tabs)
	return buf.String()
}

func (gv *GaugeStringAggView) String() string {
	return gv.stringWithIndent("")
}

// A GaugeStringAgg is a statistical summary of measures associated with a
// unique tag set.
type GaugeStringAgg struct {
	*GaugeStringStats
	Tags []tagging.Tag
}

func (ga *GaugeStringAgg) stringWithIndent(tabs string) string {
	if ga == nil {
		return "nil"
	}

	tabs2 := tabs + "  "
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%T {\n", ga)
	fmt.Fprintf(&buf, "%v  Stats: %v,\n", tabs, ga.GaugeStringStats.stringWithIndent(tabs2))
	fmt.Fprintf(&buf, "%v  Tags: %v,\n", tabs, ga.Tags)
	fmt.Fprintf(&buf, "%v}", tabs)
	return buf.String()
}

func (ga *GaugeStringAgg) String() string {
	return ga.stringWithIndent("")
}