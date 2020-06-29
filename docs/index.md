# Dash-OPS

In this first moment we will keep the docs in Portuguese (BR), until we have more contributors.

> Nota:
>
> Projeto em desenvolvimento

## API

Nossa API foi toda criada em GoLang, nela temos muitas melhorias a fazer como performance, qualidade de código e testes.

Separamos o projeto em pequenos pacotes com objetivo de algum dia separar eles em bibliotecas e plugins habilitados via configurações.

### Estrutura

| Diretório ou Arquivo | Descrição                                                                     |
| -------------------- | ----------------------------------------------------------------------------- |
| `pkg`                | Pacotes                                                                       |
| `pkg/aws`            | Futuro plugin que integra com o SDK da AWS                                    |
| `pkg/commons`        | Futura biblioteca usada para abstrair funções comuns entre os plugins e a API |
| `pkg/config`         | Pacote responsável por gerenciar as configurações do projeto                  |
| `pkg/kubernetes`     | Futuro plugin que integra com o SDK do Kubernetes                             |
| `pkg/oauth2`         | Futuro plugin que atualmente integra apenas com o SDK oauth2 do Github        |
| `pkg/spa`            | Pacote responsável por expor o frontend do projeto                            |
| `main.go`            | Arquivo de inicialização do projeto                                           |
| `dash-ops.yaml`      | Arquivo de configurações do projeto                                           |

### Rodando

Para rodar o projeto localmente é necessario criar o arquivo de configurações `dash-ops.yaml` na raiz do projeto, exemplo:

```yaml
port: 8080
origin: http://localhost:3000
headers:
  - 'Content-Type'
  - 'Authorization'
front: front/build
plugins:
  - 'OAuth2'
  - 'Kubernetes'
  - 'AWS'
oauth2:
  - provider: github
    clientId: ${GITHUB_CLIENT_ID}
    clientSecret: ${GITHUB_CLIENT_SECRET}
    authURL: 'https://github.com/login/oauth/authorize'
    tokenURL: 'https://github.com/login/oauth/access_token'
    urlLoginSuccess: 'http://localhost:3000'
    scopes:
      - user
      - repo
      - read:org
kubernetes:
  kubeconfig: ~/.kube/config
aws:
  region: us-east-1
  accessKeyId: ${AWS_ACCESS_KEY_ID}
  secretAccessKey: ${AWS_SECRET_ACCESS_KEY}
  ec2Config:
    skiplist:
      - 'EKSWorkerAutoScalingGroupSpot'
```

Caso deseje personalizar as configurações olhe as docs de cada plugin, em execute o seguinte comando:

```sh
go run main.go
```

### Build

Para rodar o build do projeto localmente, basta executar o seguinte comando:

```sh
CGO_ENABLED=0 go build -ldflags="-s -w" -o dash-ops .
```

## Frontend

Nosso frontend foi todo criado com React, nele tentamos simplicar o uso das APIs e futuros plugins.

Toda a estrutura do projeto é pensando em futuros desaclopamentos semelhante a API backend.

### Estrutura

Dentro da estrutura do React temos a pasta `src`, nela se encontra todo o nosso código:

| Diretório ou Arquivo | Descrição                                                              |
| -------------------- | ---------------------------------------------------------------------- |
| `src/components`     | Apenas componentes da estrutura(Layout) do frontend                    |
| `src/helpers`        | Futura biblioteca de utilitarios do projeto                            |
| `src/modules`        | Estrutura das páginas separadas baseado nos futuros plugins do backend |
| `src/pages`          | Páginas genéricas do projeto                                           |
| `src/App.js`         | Arquivo de inicialização do projeto front                              |

### Rodando

Para rodar o front do projeto localmente, basta executar os seguintes comandos:

```sh
cd front
yarn install
yarn start
```

### Build

Para rodar o build do front localmente, basta executar os seguintes comandos:

```sh
cd front
yarn
yarn build
```
