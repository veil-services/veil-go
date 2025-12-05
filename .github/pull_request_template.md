## 📝 Description

## 🏷️ Type of Change

- [ ] **feat**: Nova funcionalidade (detectores, opções de configuração)
- [ ] **fix**: Correção de bug
- [ ] **perf**: Melhoria de performance (latência ou alocação de memória)
- [ ] **refactor**: Alteração de código sem mudança de funcionalidade externa
- [ ] **test**: Adição ou correção de testes
- [ ] **docs**: Atualização de documentação
- [ ] **chore**: Tarefas de build, CI ou manutenção

## ⚡ Performance & Safety

- [ ] **Zero-Allocation Target**: Esta mudança mantém ou melhora a alocação de memória? (Verifique com `go test -bench=. -benchmem`)
- [ ] **Thread-Safety**: O código é seguro para uso concorrente? (Verifique com `go test -race ./...`)
- [ ] **Algoritmo**: Se você adicionou um detector, evitou o uso de `regexp` em caminhos críticos (hot paths)?

**Benchmark Results (Opcional mas recomendado para PRs de `perf`):**
```text
(Cole a saída do benchmark aqui se relevante)
```

## 🧪 Testing Checklist

- [ ] **Unit Tests**: Adicionei testes cobrindo casos de sucesso e erro.
- [ ] **Corpus Tests**: Se adicionei um novo detector, atualizei o testdata/corpus.json com casos de Falso/Verdadeiro Positivo.
- [ ] **Regression**: Executei `go test ./...` e todos os testes passaram.
- [ ] **Race Detector**: Executei `go test -race ./...` sem erros.

## ✅ Final Checklist

- [ ] O título do PR segue o formato Conventional Commits.
- [ ] O código segue o estilo do projeto (formatação, nomes de variáveis).
- [ ] Atualizei a documentação (README ou GoDocs) se necessário.
- [ ] Concordo que minha contribuição será licenciada sob a licença MIT do projeto.
