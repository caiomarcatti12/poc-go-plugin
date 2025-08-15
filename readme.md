# POC — Sistema de Plugins Nativos em Go com `-buildmode=plugin`

Esta POC demonstra como criar e carregar **plugins nativos em Go** usando o recurso `-buildmode=plugin`.
A aplicação (`host`) carrega um plugin externo (`greeter`) em tempo de execução, valida compatibilidade via **ABI** e consome a interface definida em um **módulo de contrato** (`pluginapi`).

---

## Objetivo

* Separar o **contrato** (interfaces e metadados) da implementação.
* Compilar plugins de forma **independente** do host.
* Demonstrar **carregamento dinâmico** de `.so` em Go.
* Garantir **compatibilidade binária** via verificação de ABI.
* Facilitar a extensão da aplicação sem recompilar o host.

---

## Arquitetura do Projeto

O projeto é dividido em três módulos independentes:

```
poc-go-plugin/
├── pluginapi/   → Contrato compartilhado (interfaces, constantes, tipos)
├── greeter/     → Plugin de exemplo (implementa a interface Greeter)
├── host/        → Aplicação principal (carrega e usa plugins)
└── bin/         → Saída dos plugins compilados (.so)
```

### **1. pluginapi** (Contrato)

* Define:

  * `Greeter` (interface)
  * `Info` (struct com metadados)
  * `ABI` (constante para controle de compatibilidade)
* **Importado** tanto pelo host quanto pelos plugins.
* **Versão única** que deve ser idêntica nos dois lados para evitar erros de tipo.

### **2. greeter** (Plugin)

* Implementa `Greeter` retornando uma saudação.
* Exporta:

  * `Plugin` (instância que implementa a interface)
  * `ABI` (constante com versão binária)
  * `Info` (metadados do plugin)
* Compilado como `.so` com `go build -buildmode=plugin`.

### **3. host** (Aplicação principal)

* Recebe via flag o caminho do `.so` do plugin.
* Usa `plugin.Open` para carregar o `.so` em runtime.
* Valida:

  * **ABI** do plugin vs. ABI esperado.
  * Tipo do símbolo `Plugin` (garante que implementa `Greeter`).
* Executa `Greet("Mundo")` e imprime no terminal.

---

## ⚙️ Como funciona

1. **Contrato compartilhado (`pluginapi`)**
   Define a API pública que tanto o host quanto os plugins conhecem.

2. **Plugin (`greeter`)**
   Compilado separadamente como `.so` com:

   ```bash
   go build -buildmode=plugin -o ../bin/greeter_v1.so
   ```

3. **Host (`host`)**
   Carrega o `.so`:

   ```go
   p, _ := plugin.Open("../bin/greeter_v1.so")
   sym, _ := p.Lookup("Plugin")
   greeter := sym.(pluginapi.Greeter)
   fmt.Println(greeter.Greet("Mundo"))
   ```

4. **Verificação de ABI**
   Antes de usar o plugin, o host confere se `ABI` no plugin é igual ao `ABI` esperado.

5. **Execução**
   O plugin roda no mesmo processo que o host (sem IPC), permitindo chamadas diretas.

---

## 🚀 Como rodar

### 1. Compilar o plugin

```bash
cd greeter
go build -buildmode=plugin -o ../bin/greeter_v1.so
```

### 2. Rodar o host

```bash
cd host
go run . --plugin ../bin/greeter_v1.so
```

**Saída esperada:**

```
2025/08/15 Plugin carregado: greeter v1.0.0 (ABI=1) — Plugin de exemplo que dá oi
Olá, Mundo 👋
```

---

## 📌 Considerações

* **Compatibilidade binária**
  Plugins precisam ser compilados com a **mesma versão do Go** e **mesmo pacote de contrato** (`pluginapi`), caso contrário falharão no `plugin.Open` ou no `type assertion`.

* **Portabilidade**
  O recurso `-buildmode=plugin` não é suportado no Windows (funciona em Linux/macOS).

* **Segurança**
  O plugin roda no mesmo processo, então **não há isolamento** — bugs ou panics no plugin afetam o host.

* **Quando usar**
  Ideal quando:

  * Você controla o ambiente (SO e versão do Go).
  * Precisa de performance (sem custo de IPC).
  * Quer estender o software sem recompilar o host.

