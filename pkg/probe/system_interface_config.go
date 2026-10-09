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
	"encoding/json"
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus-community/fortigate_exporter/pkg/fortigatehttpclient"
)

// probeSystemInterfaceConfig le a CONFIGURACAO das interfaces em
// /api/v2/cmdb/system/interface.
//
// POR QUE UMA PROBE SEPARADA da System/Interface:
//
//   - o endpoint /monitor/system/interface/select devolve apenas as interfaces
//     operacionais - neste FortiGate 40F, 6 de 16. O cmdb devolve TODAS,
//     inclusive as administrativamente desabilitadas, que sao justamente as que
//     interessam num painel de inventario;
//   - role, description, type e o status ADMINISTRATIVO so existem no cmdb;
//   - configuracao muda raramente, entao esta probe pode ser desabilitada em
//     instalacoes que nao a queiram, sem afetar a coleta de contadores.
//
// Exporta duas metricas:
//
//	fortigate_interface_config_info  - metadados, valor sempre 1
//	fortigate_interface_admin_up     - 1 se habilitada administrativamente
//
// O status administrativo e diferente do link: uma interface pode estar
// habilitada e sem link (cabo solto), ou desabilitada pelo operador. O
// fortigate_interface_link_up existente cobre o primeiro caso; este cobre o
// segundo.
func probeSystemInterfaceConfig(c fortigatehttpclient.FortiHTTP, _ *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mConfigInfo = prometheus.NewDesc(
			"fortigate_interface_config_info",
			"Interface configuration metadata (always 1)",
			[]string{"vdom", "name", "alias", "description", "role", "type", "mode", "ip", "mask", "mac", "status"}, nil,
		)
		mAdminUp = prometheus.NewDesc(
			"fortigate_interface_admin_up",
			"Whether the interface is administratively enabled (1) or disabled (0)",
			[]string{"vdom", "name", "alias", "role"}, nil,
		)
	)

	type ifConfig struct {
		Name string `json:"name"`
		// alias e description sao texto livre e frequentemente vazios; sao
		// exportados assim mesmo para o painel poder mostrar a lacuna em vez
		// de esconder a interface.
		Alias       string `json:"alias"`
		Description string `json:"description"`
		// "up" ou "down" - status ADMINISTRATIVO, definido pelo operador.
		Status string `json:"status"`
		// lan, wan, dmz, undefined
		Role string `json:"role"`
		// physical, vlan, tunnel, aggregate, loopback...
		Type string `json:"type"`
		// static, dhcp, pppoe
		Mode string `json:"mode"`
		// O cmdb devolve IP e mascara num unico campo separado por espaco
		// ("10.20.20.254 255.255.255.0"), diferente do /monitor, que usa
		// campos distintos. Separamos abaixo.
		IP     string `json:"ip"`
		Macddr string `json:"macaddr"`
		VDOM   string `json:"vdom"`
	}

	type ifConfigResponse struct {
		Results json.RawMessage `json:"results"`
		VDOM    string          `json:"vdom"`
	}

	var resp []ifConfigResponse
	if err := c.Get("api/v2/cmdb/system/interface", "vdom=*", &resp); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range resp {
		var ifs []ifConfig
		if err := json.Unmarshal(r.Results, &ifs); err != nil {
			log.Printf("Error parsing interface config: %v", err)
			continue
		}
		for _, ic := range ifs {
			vdom := ic.VDOM
			if vdom == "" {
				vdom = r.VDOM
			}

			ip, mask := "", ""
			if f := strings.Fields(ic.IP); len(f) == 2 {
				ip, mask = f[0], f[1]
				// 0.0.0.0 0.0.0.0 significa "sem IP", tipico de tunel e de
				// interface em modo dhcp ainda sem lease. Exportar como vazio
				// evita painel mostrando 0.0.0.0 como se fosse endereco.
				if ip == "0.0.0.0" {
					ip, mask = "", ""
				}
			}

			// VLANs e tuneis herdam o MAC da interface-pai e o cmdb devolve
			// tudo zero. Vazio e mais honesto que um MAC falso.
			mac := ic.Macddr
			if mac == "00:00:00:00:00:00" {
				mac = ""
			}

			m = append(m, prometheus.MustNewConstMetric(mConfigInfo, prometheus.GaugeValue, 1,
				vdom, ic.Name, ic.Alias, ic.Description, ic.Role, ic.Type, ic.Mode, ip, mask, mac, ic.Status))

			adminUp := 0.0
			if ic.Status == "up" {
				adminUp = 1.0
			}
			m = append(m, prometheus.MustNewConstMetric(mAdminUp, prometheus.GaugeValue, adminUp,
				vdom, ic.Name, ic.Alias, ic.Role))
		}
	}
	return m, true
}
