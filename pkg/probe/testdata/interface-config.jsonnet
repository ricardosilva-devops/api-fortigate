// Resposta de api/v2/cmdb/system/interface.
//
// Amostra sintetica de um FortiGate para testes, reduzida aos casos que importam:
//   wan     - fisica, com IP e MAC ficticios, role wan
//   lan1    - fisica, com IP, role lan, alias preenchido
//   LAB_VLAN - vlan: MAC zerado (herda da interface-pai)
//   VPN_x   - tunel: IP 0.0.0.0 (sem endereco)
//   modem   - ADMINISTRATIVAMENTE DESABILITADA, o caso que o endpoint
//             /monitor/system/interface nem devolve
[
  {
    results: [
      {
        name: 'wan',
        alias: '(Example ISP)',
        description: 'Example uplink',
        status: 'up',
        role: 'wan',
        type: 'physical',
        mode: 'static',
        ip: '198.51.100.10 255.255.255.248',
        macaddr: '02:00:00:00:10:01',
        vdom: 'root',
      },
      {
        name: 'lan1',
        alias: '(LAN)',
        description: '',
        status: 'up',
        role: 'lan',
        type: 'physical',
        mode: 'static',
        ip: '192.0.2.1 255.255.255.0',
        macaddr: '02:00:00:00:10:02',
        vdom: 'root',
      },
      {
        name: 'LAB_VLAN',
        alias: '',
        description: '',
        status: 'up',
        role: 'lan',
        type: 'vlan',
        mode: 'static',
        ip: '203.0.113.1 255.255.255.0',
        macaddr: '00:00:00:00:00:00',
        vdom: 'root',
      },
      {
        name: 'VPN_EXAMPLE',
        alias: '',
        description: '',
        status: 'up',
        role: 'undefined',
        type: 'tunnel',
        mode: 'static',
        ip: '0.0.0.0 0.0.0.0',
        macaddr: '00:00:00:00:00:00',
        vdom: 'root',
      },
      {
        name: 'modem',
        alias: '',
        description: '',
        status: 'down',
        role: 'undefined',
        type: 'physical',
        mode: 'static',
        ip: '0.0.0.0 0.0.0.0',
        macaddr: '00:00:00:00:00:00',
        vdom: 'root',
      },
    ],
    vdom: 'root',
  },
]
