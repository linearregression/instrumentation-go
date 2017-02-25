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

// GaugeBoolViewDesc defines an bool gauge view.
type GaugeBoolViewDesc struct {
	*ViewDescCommon
}

func (gd *GaugeBoolViewDesc) createAggregator(t time.Time) (aggregator, error) {
	return newGaugeAggregatorBool(), nil
}

func (gd *GaugeBoolViewDesc) retrieveView(now time.Time) (*View, error) {
	gav, err := gd.retrieveAggreationView(now)
	if err != nil {
		return nil, err
	}
	return &View{
		ViewDesc: gd,
		ViewAgg:  gav,
	}, nil
}

func (gd *GaugeBoolViewDesc) viewDesc() *ViewDescCommon {
	return gd.ViewDescCommon
}

func (gd *GaugeBoolViewDesc) isValid() error {
	return nil
}

func (gd *GaugeBoolViewDesc) retrieveAggreationView(t time.Time) (*GaugeBoolAggView, error) {
	var aggs []*GaugeBoolAgg

	for sig, a := range gd.signatures {
		tags, err := tagging.TagsFromValuesSignature([]byte(sig), gd.TagKeys)
		if err != nil {
			return nil, fmt.Errorf("malformed signature '%v'. %v", sig, err)
		}
		aggregator, ok := a.(*gaugeAggregatorBool)
		if !ok {
			return nil, fmt.Errorf("unexpected aggregator type. got %T, want stats.gaugeAggregatorBool", a)
		}
		ga := &GaugeBoolAgg{
			GaugeBoolStats: aggregator.retrieveCollected(),
			Tags:           tags,
		}
		aggs = append(aggs, ga)
	}

	return &GaugeBoolAggView{
		Descriptor:   gd,
		Aggregations: aggs,
	}, nil
}

func (gd *GaugeBoolViewDesc) stringWithIndent(tabs string) string {
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

func (gd *GaugeBoolViewDesc) String() string {
	return gd.stringWithIndent("")
}

// GaugeBoolAggView is the set of collected GaugeBoolAgg associated with
// ViewDesc.
type GaugeBoolAggView struct {
	Descriptor   *GaugeBoolViewDesc
	Aggregations []*GaugeBoolAgg
}

func (gv *GaugeBoolAggView) stringWithIndent(tabs string) string {
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

func (gv *GaugeBoolAggView) String() string {
	return gv.stringWithIndent("")
}

// A GaugeBoolAgg is a statistical summary of measures associated with a
// unique tag set.
type GaugeBoolAgg struct {
	*GaugeBoolStats
	Tags []tagging.Tag
}

func (ga *GaugeBoolAgg) stringWithIndent(tabs string) string {
	if ga == nil {
		return "nil"
	}

	tabs2 := tabs + "  "
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%T {\n", ga)
	fmt.Fprintf(&buf, "%v  Stats: %v,\n", tabs, ga.GaugeBoolStats.stringWithIndent(tabs2))
	fmt.Fprintf(&buf, "%v  Tags: %v,\n", tabs, ga.Tags)
	fmt.Fprintf(&buf, "%v}", tabs)
	return buf.String()
}

func (ga *GaugeBoolAgg) String() string {
	return ga.stringWithIndent("")
}
