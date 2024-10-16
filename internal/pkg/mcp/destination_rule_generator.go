// Copyright Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mcp

import (
	"github.com/jewertow/federation/internal/pkg/istio"
	"github.com/jewertow/federation/internal/pkg/xds"
	"github.com/jewertow/federation/internal/pkg/xds/adss"
	"google.golang.org/protobuf/types/known/anypb"
	istiocfg "istio.io/istio/pkg/config"
)

var _ adss.RequestHandler = (*DestinationRuleResourceGenerator)(nil)

type DestinationRuleResourceGenerator struct {
	cf *istio.ConfigFactory
}

func NewDestinationRuleResourceGenerator(cf *istio.ConfigFactory) *DestinationRuleResourceGenerator {
	return &DestinationRuleResourceGenerator{cf: cf}
}

func (v *DestinationRuleResourceGenerator) GetTypeUrl() string {
	return xds.DestinationRuleTypeUrl
}

func (v *DestinationRuleResourceGenerator) GenerateResponse() ([]*anypb.Any, error) {
	dr := v.cf.GetDestinationRules()
	if dr == nil {
		return nil, nil
	}
	return serialize(&istiocfg.Config{
		Meta: istiocfg.Meta{
			Name:      dr.Name,
			Namespace: dr.Namespace,
		},
		Spec: &dr.Spec,
	})
}
