# Kubernetes

> Nota:
>
> Plugin ainda em desenvolvimento

Neste plugin tentamos abstrair funções para facilitar a vida de desenvolvedores que trabalham focados em criar e dar manuteção em features do projeto e não precisam se preocupar tanto com a infraestrutara do mesmo.

O objetivo aqui é tentar diminuir a carga cognitiva dos engenheiros em cima do kubernetes, deixando apenas acesso as informações como visualição do que esta rodando no momento e formas de consultar logs com facilidade.

## Configurações do Plugin

Em seu arquivo de configuração você poderá adicionar seus clusters do kubernetes, exemplo:

```yaml
kubernetes:
  - name: 'Kubernetes Dev'
    kubeconfig: ~/.kube/config
    context: 'dev'
  - name: 'Kubernetes Prod'
    kubeconfig: ~/.kube/config
    context: 'prod'
```

Caso você esteja rodando no cluster o `dash-ops` e adicionou configurações de `ClusterRole` no `rbac`, você simplesmente pode rodar o plugin kuberentes sem o `kubeconfig`, seria a configuração `inCluster`, exemplo:

```yaml
kubernetes:
  - name: 'Kubernetes Dev'
    kubeconfig:
```

Neste caso o plugin vai ter permissões de acessar diretamente a API do kuberentes que ele esta rodando sem um kubeconfig.

> Se você esta usando o nosso template helm o `ClusterRole` já vem pre configurado (https://github.com/dash-ops/helm-charts)

### Permissionamento

> No momento a única funcionalidade do Kubernetes plugin que afeta algo no kuberentes é o `scale` dos `deployments`, essa função é recomenda apenas em clusters no ambiente de desenvolvimento.

Exemplo de como adicionar a permissão:

```yaml
kubernetes:
  - name: 'Kubernetes Dev'
    kubeconfig: ~/.kube/config
    context: 'dev'
    permission:
      deployments:
        start: ['org*team']
        stop: ['org*team']
```

> `org*team`: Organização e o time do Github com a permissão de executar o `scale` para `1` ou `0` do `deployment` no kuberentes.
