# Dash-OPS - WIP

![DashOps](https://github.com/dash-ops/dash-ops/workflows/DashOps/badge.svg)

Dash-ops is under construction, the goal is to create a simple and permissioned interface with small actions for developers who use some features of structuring tools such as Kubernetes, AWS, GitHub...

We want to remove the cognitive burden of many engineers who just want to focus on the development of their features and leave the responsibility of managing the structuring part of their system to teams focused as teams of SREs.

## Docs

Access the [document directory here](/docs).

## Running local

Create a configuration file in the project's root directory called `dash-ops.yaml`, example:

```yaml
port: 8080
origin: http://localhost:8080
headers:
  - 'Content-Type'
  - 'Authorization'
front: app
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
    urlLoginSuccess: 'http://localhost:8080'
    orgPermission: 'dash-ops'
    scopes:
      - user
      - repo
      - read:org
kubernetes:
  - name: 'Kubernetes Dev'
    kubeconfig: ${HOME}/.kube/config
    context: 'dev'
  - name: 'Kubernetes Prod'
    kubeconfig: ${HOME}/.kube/config
    context: 'prod'
aws:
  region: us-east-1
  accessKeyId: ${AWS_ACCESS_KEY_ID}
  secretAccessKey: ${AWS_SECRET_ACCESS_KEY}
  ec2Config:
    skipList:
      - 'EKSWorkerAutoScalingGroupSpot'
```

This project has backend in GoLang and frontend with React.
We recommend using the last generated docker image to test locally:

```sh
docker run --rm \
  -v $(pwd)/dash-ops.yaml:/dash-ops.yaml \
  -v /home/my-user/.kube/config:/.kube/config \
  -e GITHUB_CLIENT_ID=666 \
  -e GITHUB_CLIENT_SECRET=666xpto \
  -e AWS_ACCESS_KEY_ID=999 \
  -e AWS_SECRET_ACCESS_KEY=999xpto \
  -p 8080:8080 \
  -it dashops/dash-ops
```

## Running on a kubernetes cluster with helm:

We use our own helm chart, for more information go here, or follow the instructions below to execute:

### Create a `values.yaml`

File with your settings and secrets, example:

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
    - name: 'Kubernetes InCluster'
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
