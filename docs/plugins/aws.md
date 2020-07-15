# AWS

> Nota:
>
> Plugin ainda em desenvolvimento

Neste plugin tentamos abstrair funções para facilitar a vida de desenvolvedores que trabalham focados em criar e dar manuteção em features do projeto e não precisam se preocupar tanto com a infraestrutura do mesmo.

O objetivo aqui é tentar diminuir a carga cognitiva dos engenheiros em cima dos produtos da AWS, deixando apenas o acesso as funções usadas no dia a dia de desenvolvimento.

## Configurações do Plugin

Em seu arquivo de configurações você precisara adicionar suas credenciais para ter permissões, exemplo:

```yaml
aws:
  region: us-east-1
  accessKeyId: ${AWS_ACCESS_KEY_ID}
  secretAccessKey: ${AWS_SECRET_ACCESS_KEY}
```

### EC2

Na listagem das instâncias dos EC2 você pode ocultar algumas que não deseja que os times tenham acesso, exemplo:

```yaml
aws:
  ...
  ec2Config:
    skipList:
      - 'EKSWorkerAutoScalingGroupSpot'
```

### Permissionamento

> No momento a única funcionalidade do AWS plugin que afeta algo na AWS é o `start` e `stop` das instâncias EC2, essa função é recomenda apenas em clusters no ambiente de desenvolvimento.

Exemplo de como adicionar a permissão:

```yaml
aws:
  ...
  permission:
    ec2:
      start: ["org*team"]
      stop: ["org*team"]
```

> `org*team`: Organização e o time do Github com a permissão de executar o `start` ou `stop` das instâncias EC2.
