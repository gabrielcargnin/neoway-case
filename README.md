### Configurações iniciais ###

- Setar as variáveis de ambiente no arquivo .env na raiz do projeto;
- Rodar os sqls do arquivo app.sql na pasta postgres;

### Rodar localmente ###

Executar os comandos:
```
go build %GOPATH%/src/neoway-case/consumption-service
go run %GOPATH%/src/neoway-case/consumption-service

```

### Rodar localmente no Idea GoLand ###

Ir em Run/Debug Configurations, adicionar nova configuração de Go Build e preencher o diretório com o caminho para o pacote consumption-service.
Confirmar que a opção "Run after build" está marcada.
    
### Rodar no docker ###

Para rodar no docker basta rodar 
```
docker-compose up -d --build
```



