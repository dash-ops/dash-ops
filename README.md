# Dash-OPS - WIP

![DashOps](https://github.com/dash-ops/dash-ops/workflows/DashOps/badge.svg)

Dash-ops is under construction, the goal is to create a simple and permissioned interface with small actions for developers who use some features of structuring tools such as Kubernetes, AWS, GitHub...

## Running local

Create a configuration file in the project's root directory called `dash-ops.yaml`, example:

```yaml
port: 8080
origin: http://localhost:3000
headers:
  - 'Content-Type'
  - 'Authorization'
front: front/build
plugins:
  - 'Oauth2'
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
    blacklist:
      - 'EKSWorkerAutoScalingGroupSpot'
```

This project has backend in GoLang and frontend with React.
To run the backend just run:

```sh
go run main.go
```

And for the frontend run:

```sh
cd front
yarn
yarn start
```

## Running on a kubernetes cluster with helm:

Create a `values.yaml` file with your settings and secrets, example:

```yaml
name: dash-ops
image:
  name: dashops/dash-ops
  tag: latest

container:
  limits:
    memory: 100Mi
  requests:
    memory: 100Mi
  httpPort: 8080
  probes:
    readinessDelay: 5
    livenessDelay: 5
    path: /api/health
  env: {}
  secrets:
    GITHUB_CLIENT_ID: 666=
    GITHUB_CLIENT_SECRET: 666xpto==
    AWS_ACCESS_KEY_ID: 999=
    AWS_SECRET_ACCESS_KEY: 999xpto==

ingress:
  enabled: true
  externalDNSTarget: traefik.local.
  hosts:
    - host: dash-ops.dev

configMap: |
  port: 8080
  origin: http://dash-ops.dev
  headers: 
    - "Content-Type"
    - "Authorization"
  front: app
  plugins:
    - "Oauth2"
    - "Kubernetes"
    - "AWS"
  oauth2:
    - provider: github
      clientId: ${GITHUB_CLIENT_ID}
      clientSecret: ${GITHUB_CLIENT_SECRET}
      authURL: "https://github.com/login/oauth/authorize"
      tokenURL: "https://github.com/login/oauth/access_token"
      urlLoginSuccess: "http://dash-ops.dev"
      scopes: 
        - user
        - repo
        - read:org
  kubernetes:
    kubeconfig:
  aws:
    region: us-east-1
    accessKeyId: ${AWS_ACCESS_KEY_ID}
    secretAccessKey: ${AWS_SECRET_ACCESS_KEY}
```

### Install

```sh
helm install --name dash-ops dash-ops/dash-ops --values ./values.yaml
```

### Upgrading

```sh
helm upgrade dash-ops dash-ops/dash-ops --values ./values.yaml
```
