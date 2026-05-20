# Archive Notice

**Archived**: May 2026  
**Status**: No longer actively developed

## O que foi este projeto

DashOps nasceu de uma dor recorrente: em cada empresa, a mesma reinvenção da roda — portais internos fragmentados, ferramentas de infra sem conexão, desenvolvedores alternando entre dezenas de abas para ter contexto do próprio serviço.

A ideia era um hub centralizado com arquitetura de plugins que conectasse Kubernetes, AWS, Service Catalog, Observability (Loki, Prometheus, Tempo) em uma interface unificada, alimentada por IA.

## O que foi construído

- Arquitetura hexagonal completa em Go (8 módulos migrados)
- Sistema de plugins dinâmico com Dependency Injection
- Integrações reais com Loki, Prometheus, Tempo, AlertManager
- Frontend React modular com módulos de K8s, AWS, Observability, Service Catalog
- Stack de desenvolvimento local com Kind + observability stack completa
- Website de documentação em https://dash-ops.github.io/
- Design system com mockups interativos via Figma submodules

## Por que foi arquivado

O ecossistema mudou mais rápido do que o projeto conseguiu acompanhar.

O MCP (Model Context Protocol) resolveu o layer de integração de uma forma que nenhum portal web conseguia — hoje um desenvolvedor consegue plugar um LLM no contexto da própria infraestrutura em horas. Ferramentas como K8sGPT, HolmesGPT, os MCP servers da AWS e Red Hat para Kubernetes, e agregadores como MetaMCP preencheram os gaps que o DashOps se propunha a resolver.

O espaço de Internal Developer Platforms também consolidou em torno de Backstage (CNCF) e Port.io, tornando difícil justificar manutenção de uma alternativa sem time dedicado.

## O que ficou de aprendizado

- Arquitetura hexagonal em Go com DI aplicada a um projeto real
- Plugin system com carregamento dinâmico de módulos
- Integração com o ecossistema de observabilidade (Loki, Tempo, Prometheus)
- Design de APIs com separação clara entre ports, adapters e controllers
- Reconhecer quando um mercado se move e encerrar um ciclo com clareza
