Estudo de caso - Circuit Breaker
---

Este projeto tem como objetivo o estudo do Pattern Circuit Breaker.
A iniciativa permite criar uma rede Docker com diversos microserviços interconectados contendo diferentes configurações
de Circuit Breaker implementados, e com tempo de respostas diferentes.


## Instalando

O projeto utiliza as linguagens de programação Golang e Python3 e a tecnologia Docker. 
Para executar os testes basta ter instalado Python3 e Docker.

> Mas é aconselhável ter as três tecnologias devidamente configuradas em sua máquina caso deseje desenvolver.

Utilizando devidamente um *virtual environment manager* para projetos Python3 (virtuaenv, pyenv, etc). 
Execute no terminal:

```sh
$ pip install -r case-study/requirements/requirements.txt
```

## Testando

#### Configuração padrão

A arquitetura configurada por padrão está descrita a seguir.

![circuit-breaker-case-Com Circuit Breaker](https://user-images.githubusercontent.com/1136326/97424275-3d528280-18ef-11eb-9f84-a60693e41485.png)

A quantidade de instâncias e seu fluxo de requisições podem ser modificados alterando os arquivos *docker-compose.yaml* e *resources/config.*.yaml*
1. **docker-compose.yaml**: instâncias docker serão levantadas em uma rede específica
2. **resources/config.*.yaml**: Configuração de cada instância docker informando quais instâncias deverá ser requisitada.

Para rodar a aplicação é necessário apenas o Python3 e o Docker.Portanto acesse a pasta do projeto, abra 3 terminais e execute os comandos sequenciamente usando o ambiente com os devidos **requisitos instalados**.

##### Terminal 1
```sh
$ make run
```

##### Terminal 2
```sh
$ cd case-study
$ locust -f locustfile.py -u 20 -r 100 --headless
```

##### Terminal 2
```sh
$ docker stats
```

