// Copyright The Prometheus Authors
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

package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSystemInterfaceConfig(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/cmdb/system/interface", "testdata/interface-config.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemInterfaceConfig, c, r) {
		t.Errorf("probeSystemInterfaceConfig() returned non-success")
	}

	// O que este teste protege:
	//   - "modem" aparece com admin_up=0. Ela nao existe no endpoint
	//     /monitor/system/interface, entao so esta probe a enxerga
	//   - tunel com ip "0.0.0.0 0.0.0.0" vira ip/mask vazios, nao "0.0.0.0"
	//   - vlan com macaddr zerado vira mac vazio, nao "00:00:00:00:00:00"
	//   - description vazia e exportada como label vazio, sem sumir a serie
	em := `
	# HELP fortigate_interface_admin_up Whether the interface is administratively enabled (1) or disabled (0)
	# TYPE fortigate_interface_admin_up gauge
	fortigate_interface_admin_up{alias="",name="LAB_VLAN",role="lan",vdom="root"} 1
	fortigate_interface_admin_up{alias="",name="VPN_EXAMPLE",role="undefined",vdom="root"} 1
	fortigate_interface_admin_up{alias="",name="modem",role="undefined",vdom="root"} 0
	fortigate_interface_admin_up{alias="(LAN)",name="lan1",role="lan",vdom="root"} 1
	fortigate_interface_admin_up{alias="(Example ISP)",name="wan",role="wan",vdom="root"} 1
	# HELP fortigate_interface_config_info Interface configuration metadata (always 1)
	# TYPE fortigate_interface_config_info gauge
	fortigate_interface_config_info{alias="",description="",ip="203.0.113.1",mac="",mask="255.255.255.0",mode="static",name="LAB_VLAN",role="lan",status="up",type="vlan",vdom="root"} 1
	fortigate_interface_config_info{alias="",description="",ip="",mac="",mask="",mode="static",name="VPN_EXAMPLE",role="undefined",status="up",type="tunnel",vdom="root"} 1
	fortigate_interface_config_info{alias="",description="",ip="",mac="",mask="",mode="static",name="modem",role="undefined",status="down",type="physical",vdom="root"} 1
	fortigate_interface_config_info{alias="(LAN)",description="",ip="192.0.2.1",mac="02:00:00:00:10:02",mask="255.255.255.0",mode="static",name="lan1",role="lan",status="up",type="physical",vdom="root"} 1
	fortigate_interface_config_info{alias="(Example ISP)",description="Example uplink",ip="198.51.100.10",mac="02:00:00:00:10:01",mask="255.255.255.248",mode="static",name="wan",role="wan",status="up",type="physical",vdom="root"} 1
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
