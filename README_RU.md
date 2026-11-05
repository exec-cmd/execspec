<div align="center">

# execspec

**Разработка по спецификациям для работы с ИИ-агентами, в приватной папке `spec/`.**

[![Go 1.27+](https://img.shields.io/badge/go-1.27%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green)](LICENSE)
![Local first](https://img.shields.io/badge/local--first-без%20LLM%20и%20сети-informational)

[English](README.md) · Русский

</div>

`xspec` — небольшая CLI-утилита, которая держит намерения рядом с кодом. Фичи, баги, обновления, фиксы и планы получают номера и лежат обычным Markdown в папке `spec/`, которую git игнорирует. Вы и ваш ИИ-агент читаете и пишете одни и те же файлы. CLI создаёт их из шаблонов, ведёт по статусам `open → doing → done` и не пускает дальше, пока остаются открытые вопросы. В комплекте навык (skill), который учит этому процессу Claude Code, Codex и других агентов.

## Возможности

- **Приватно по умолчанию.** В `spec/` лежат намерения: планы, мнения, открытые вопросы. `xspec init` убирает папку из git.
- **Номера.** Один счётчик на все типы: `xspec start 7`, `xspec done 7`, `--with 7,9`.
- **Проверяемые требования.** У фич требования записаны в нотации EARS (`WHEN … THE SYSTEM SHALL …`), а каждая задача ссылается на своё требование (`(R1)`). У багов три поведения: текущее, ожидаемое и *то, что не должно сломаться*.
- **Вы одобряете, агент строит.** Метки `[?: вопрос]` не дают одобрить спеку. `xspec approve` запускаете только вы, и любая правка после него сбрасывает одобрение. Неотмеченные задачи и упавшие команды `Verify:` не дают сделать `done`.
- **Спеки нужного размера.** Kickoff на весь проект, `--depth light` для мелкой работы, один update с частями, когда задач много сразу.
- **Всегда понятно, что дальше.** `xspec next`, `xspec questions` и генерируемый `ROADMAP.md` с историей.
- **Пакеты.** *update* собирает фичи и issues, *fix* — issues и мелкие баги. Когда пакет закрывается, закрывается и всё, что в него входит.
- **Чистая история.** `done` записывает в файл коммит и переносит файл в `spec/archive/`.
- **Навык для агентов.** `xspec install global` ставит его для Claude Code (`~/.claude/skills`) и других агентов (`~/.agents/skills`).
- **Без LLM внутри.** Ключи и сеть не нужны, всё в одном статическом бинарнике.

## Установка

Скачайте бинарник для своей системы со страницы [Releases](https://github.com/exec-cmd/execspec/releases), переименуйте в `xspec` и положите в `PATH`:

```sh
chmod +x xspec-linux-amd64 && mv xspec-linux-amd64 ~/.local/bin/xspec
```

Или через Go 1.27+:

```sh
go install github.com/exec-cmd/execspec/src/xspec@latest
# или из исходников:
make build    # → build/bin/xspec
```

Затем поставьте навык для агентов:

```sh
xspec install global     # ~/.claude/skills/execspec и ~/.agents/skills/execspec
```

## Быстрый старт

```sh
xspec init                                  # spec/, kickoff, строка в .gitignore
xspec new feature "Login form"              # spec/features/002-login-form.md
xspec new issue "Crash on save"             # spec/issues/003-crash-on-save.md
xspec new update v0.2 --with 2,3            # spec/updates/004-v0-2/{update,design}.md
xspec questions                             # что ждёт вашего ответа
xspec approve 4                             # только вы; откажет, пока есть метки [?: …]
xspec start 4
xspec next                                  # следующая задача
xspec done 4                                # откажет, если задачи не отмечены или Verify упал
xspec check                                 # код выхода 1 при любой проблеме
```

### Как выглядит работа

```text
$ xspec list
NUM  TYPE     STATUS  [?]  TASKS  TITLE
002  feature  open    0    0/2    Login form
003  issue    open    0    0/2    Crash on save
004  update   open    0    0/1    v0.2

$ xspec done 4
verify: go test ./...
004 done   v0.2 -> spec/archive/updates/004-v0-2
002 done   Login form -> spec/archive/features/002-login-form.md
003 done   Crash on save -> spec/archive/issues/003-crash-on-save.md
```

Так выглядит новый файл фичи. Агент его заполняет, вы одобряете, и только после этого пишется код:

```markdown
status: open
# 002 Login form

## Why
[?: the problem, and who gets what when it is solved]

## Requirements
- R1: WHEN [?: condition] THE SYSTEM SHALL [?: behavior]

## Design
[?: the approach and the code it touches]

## Tasks
- [ ] T1 [?: task] (R1)

Verify: [?: a command that must pass before done, e.g. go test ./...; or delete this line]
```

## Документация

Документация на английском.

| Страница | Что там |
| --- | --- |
| [Usage](docs/user/usage.md) | Все команды и флаги |
| [Workflow](docs/user/workflow.md) | Структура `spec/`, типы, шаблоны и рабочий цикл |
| [Architecture](docs/dev/architecture.md) | Внутреннее устройство для контрибьюторов и агентов; начните с [AGENTS.md](AGENTS.md) |

## Разработка

```sh
make          # fmt, vet, test, build → build/bin/xspec
go test ./...
```

## Лицензия

MIT. См. [LICENSE](LICENSE).
