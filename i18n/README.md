# Sistema de Internacionalização (i18n) - MGC CLI

Este diretório contém o sistema de internacionalização da MGC CLI, permitindo que a interface seja exibida em múltiplos idiomas.

## Estrutura

```
i18n/
├── i18n.go                    # Sistema principal de i18n
├── translations/              # Arquivos de tradução
│   ├── pt-BR.json            # Português (Brasil)
│   └── en-US.json            # Inglês (EUA)
└── README.md                 # Esta documentação
```

## Como Funciona

### 1. Detecção Automática de Idioma

O sistema detecta automaticamente o idioma preferido do usuário na seguinte ordem:

1. **Variável de ambiente `CLI_LANG`** (ex: `export CLI_LANG=pt-BR`)
2. **Variável de ambiente `LANG`** (ex: `export LANG=pt_BR.UTF-8`)
3. **Variável de ambiente `LC_ALL`** (ex: `export LC_ALL=pt_BR.UTF-8`)
4. **Idioma padrão** (pt-BR)

### 2. Flag de Linha de Comando

Você pode forçar um idioma específico usando a flag `--lang`:

```bash
mgc --lang=en-US help
mgc --lang=pt-BR i18n list
```

### 3. Comando i18n

A CLI inclui um comando dedicado para gerenciar idiomas:

```bash
# Listar idiomas disponíveis
mgc i18n list

# Definir idioma atual (temporário)
mgc i18n set en-US

# Mostrar informações sobre um idioma
mgc i18n info pt-BR

# Mostrar idioma atual
mgc i18n current
```

## Adicionando Novos Idiomas

### 1. Criar Arquivo de Tradução

Crie um novo arquivo JSON em `i18n/translations/` seguindo o padrão:

```json
{
  "code": "es-ES",
  "name": "Spanish (Spain)",
  "native_name": "Español (España)",
  "translations": {
    "cli.short_description": "MGC CLI - Herramienta de línea de comandos para Magalu Cloud",
    "cli.long_description": "MGC CLI es una herramienta poderosa para gestionar recursos en Magalu Cloud.",
    "cli.usage": "Uso",
    "cli.available_commands": "Comandos Disponibles",
    // ... mais traduções
  }
}
```

### 2. Estrutura de Chaves de Tradução

Use uma hierarquia de chaves para organizar as traduções:

- `cli.*` - Textos da interface principal
- `i18n.*` - Textos do sistema de i18n
- `commands.*` - Textos específicos de comandos
- `errors.*` - Mensagens de erro

### 3. Usando Traduções no Código

```go
import "github.com/magaluCloud/mgccli/i18n"

func myFunction() {
    manager := i18n.GetInstance()
    
    // Tradução simples
    message := manager.T("cli.short_description")
    
    // Tradução com parâmetros
    message := manager.T("cli.welcome_user", "João")
    
    // Verificar idioma atual
    currentLang := manager.GetLanguage()
}
```

## Formatos de Código de Idioma

Use códigos de idioma no formato ISO 639-1 + ISO 3166-1:

- `pt-BR` - Português (Brasil)
- `en-US` - Inglês (EUA)
- `es-ES` - Espanhol (Espanha)
- `fr-FR` - Francês (França)
- `de-DE` - Alemão (Alemanha)

## Variáveis de Ambiente

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `CLI_LANG` | Idioma preferido da CLI | `export CLI_LANG=pt-BR` |
| `LANG` | Idioma do sistema | `export LANG=pt_BR.UTF-8` |
| `LC_ALL` | Localização completa | `export LC_ALL=pt_BR.UTF-8` |

## Fallback

Se uma tradução não for encontrada:

1. O sistema usa a chave original como fallback
2. Se o idioma não for suportado, usa o idioma padrão (pt-BR)
3. Se nenhum arquivo de tradução estiver disponível, a CLI funciona normalmente

## Desenvolvimento

### Adicionando Novas Traduções

1. Adicione a chave em todos os arquivos de tradução
2. Use a função `manager.T("chave")` no código
3. Teste com diferentes idiomas

### Testando

```bash
# Testar com português
export CLI_LANG=pt-BR
./gfcli --help

# Testar com inglês
export CLI_LANG=en-US
./gfcli --help

# Testar com flag
./gfcli --lang=en-US --help
```

## Boas Práticas

1. **Mantenha consistência** nas chaves de tradução
2. **Use hierarquia** para organizar traduções
3. **Teste todos os idiomas** ao adicionar novas funcionalidades
4. **Documente** novas chaves de tradução
5. **Mantenha fallbacks** para casos de erro 