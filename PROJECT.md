# api-fortigate — notas de implementação

Este repositório é baseado em [prometheus-community/fortigate_exporter](https://github.com/prometheus-community/fortigate_exporter), versão de referência v1.25.0. O código derivado mantém os avisos de atribuição e a licença Apache-2.0 em `LICENSE` e `NOTICE`.

## Objetivo

Coletar métricas da API REST do FortiOS e exportá-las no formato Prometheus, incluindo informações de configuração e estado das interfaces de rede. Os exemplos, endereços, nomes de equipamentos e tokens apresentados neste repositório são fictícios e não identificam nenhuma infraestrutura operacional.

## Funcionalidades acrescentadas

- `fortigate_interface_info`: metadados da interface obtidos de `/api/v2/monitor/system/interface/select`, incluindo IP, máscara, MAC e duplex.
- `fortigate_interface_config_info`: nome, alias, descrição, tipo, função, modo, status, IP e MAC via `/api/v2/cmdb/system/interface`.
- `fortigate_interface_admin_up`: indica se a interface está administrativamente habilitada.
- `fortigate_setup_rate`: taxa de criação de novas sessões por protocolo.
- `fortigate_offloaded_sessions`: sessões tratadas em hardware, quando disponíveis.
- `fortigate_log_rate`: taxa de registros por destino, quando disponível.

As informações descritivas ficam em métricas de metadados, evitando adicionar atributos voláteis, como IP de interface WAN, às séries de contadores de tráfego. Endereços `0.0.0.0` e MACs zerados são normalizados como campos vazios quando não representam um valor configurado.

## Dados sintéticos para testes

O arquivo `pkg/probe/testdata/interface-config.jsonnet` contém exclusivamente uma topologia ilustrativa de cinco interfaces: `wan`, `lan1`, `LAB_VLAN`, `VPN_EXAMPLE` e `modem`. Usa as faixas reservadas para documentação RFC 5737 (`192.0.2.0/24`, `198.51.100.0/24` e `203.0.113.0/24`) e MACs fictícios localmente administrados.

O teste `pkg/probe/system_interface_config_test.go` usa os mesmos identificadores e endereços para verificar o formato das métricas. `fortigate-key.example.yaml` contém somente um token fictício de exemplo; copie-o para `fortigate-key.yaml` no ambiente local antes de configurar uma credencial operacional. Nunca substitua exemplos versionados por credenciais, endereços, nomes de clientes ou respostas obtidas de dispositivos reais.

## Validação

```bash
gofmt -l ./pkg/
go vet ./...
go test ./...
```

Para uso real, mantenha os tokens fora do controle de versão, em um arquivo local protegido e ignorado pelo Git, conforme o `.gitignore`.

## Origem e licença

Conserve os avisos legais dos autores originais e identifique adaptações próprias ao redistribuir o projeto. Este software não é um produto oficial da Fortinet.
