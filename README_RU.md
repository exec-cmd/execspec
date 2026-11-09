<div align="center">

# execspec

**Разработка по спецификациям для работы с ИИ-агентами, в приватной папке `spec/`.**

[![Go 1.27+](https://img.shields.io/badge/go-1.27%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green)](LICENSE)
![Local first](https://img.shields.io/badge/local--first-без%20LLM%20и%20сети-informational)

[English](README.md) · Русский

</div>

> [!WARNING]
> Проект находится в ранней разработке и может сильно меняться: команды, флаги, формат файлов в `spec/` и навык для агентов могут меняться между версиями без обратной совместимости. Закрепляйте версию и читайте release notes перед обновлением.

`xspec` — небольшая CLI-утилита, которая держит намерения рядом с кодом. Изменения (changes) и обновления (updates) получают номера и лежат обычным Markdown в папке `spec/`, которую git игнорирует; рядом короткий бриф проекта `SPEC.md`. Вы и ваш ИИ-агент читаете и пишете одни и те же файлы. CLI создаёт их из шаблонов, закрывает `open → done` и отказывает, пока остаются открытые вопросы или невыполненные задачи. В комплекте навык для агентов: он учит Claude Code, Codex и другие агенты тому же процессу.

## Возможности

- **Приватно по умолчанию.** В `spec/` лежат намерения: планы, мнения, открытые вопросы. `xspec init` убирает папку из git.
- **Номера.** Один счётчик на оба типа: `xspec done 7`, `--with 7,9`.
- **Проверяемые требования.** Когда поведению нужна точная формулировка, добавляйте требования в нотации EARS (`WHEN … THE SYSTEM SHALL …`) и задачи, которые на них ссылаются (`(R1)`). У бага три поведения: текущее, ожидаемое и *неизменное*.
- **Агент работает, вы отвечаете.** Агент показывает короткий план в чате и начинает работу. Он останавливается только на открытой метке `[?: вопрос]` или в зоне *Ask first* из `SPEC.md`. Неотмеченные задачи и упавшие команды `Verify:` не дают закрыть.
- **Спеки нужного размера.** Задача в одно предложение обходится без файла. *Change* — одна работа, *update* — пакет; оба начинаются с Goal / Tasks / Verify и растут по мере нужды.
- **Бриф, который пишет агент.** В `SPEC.md` — Context, Audience, Principles, Design, Stack и Ask first. Агент составляет черновик по вашему репозиторию и нескольким вопросам.
- **Всегда понятно, что дальше.** `xspec next`, `xspec questions` и генерируемый `ROADMAP.md` с историей.
- **Чистая история.** `done` записывает в файл коммит и переносит файл в `spec/archive/`; `--decision` сохраняет одну строку с причиной в `notes/decisions.md`.
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
xspec init                                  # spec/, бриф SPEC.md, строка в .gitignore
xspec new change "Login form"               # spec/changes/001-login-form.md
xspec new change "Crash on save" --bug      # spec/changes/002-crash-on-save.md
xspec new update v0.2 --with 1,2            # spec/updates/003-v0-2.md (--dir добавит design.md)
xspec questions                             # что ждёт вашего ответа
xspec next                                  # следующая задача
xspec done 3                                # откажет, если задачи не отмечены или Verify упал
xspec check                                 # код выхода 1 при любой проблеме
```

### Как выглядит работа

```text
$ xspec list
NUM  TYPE    STATUS  [?]  TASKS  TITLE
001  change  open    0    1/2    Login form
002  change  open    0    0/2    Crash on save
003  update  open    0    0/1    v0.2

$ xspec done 3 --decision "sessions over JWT: simpler to revoke"
verify: go test ./...
003 done   v0.2 -> spec/archive/updates/003-v0-2.md
001 done   Login form -> spec/archive/changes/001-login-form.md
002 done   Crash on save -> spec/archive/changes/002-crash-on-save.md
```

Так выглядит новый файл change. Агент его заполняет, показывает короткий план и пишет код:

```markdown
status: open
# 001 Login form

## Goal
[?: what changes, and for whom]

## Tasks
- [ ] T1 [?: task]

Verify: [?: a command that must pass before done, e.g. go test ./...; or delete this line]
```

## Документация

Документация на английском.

| Страница | Что там |
| --- | --- |
| [Usage](docs/user/usage.md) | Все команды и флаги |
| [Workflow](docs/user/workflow.md) | Структура `spec/`, типы, бриф и рабочий цикл |
| [Architecture](docs/dev/architecture.md) | Внутреннее устройство для контрибьюторов и агентов; начните с [AGENTS.md](AGENTS.md) |

## Разработка

```sh
make          # fmt, vet, test, build → build/bin/xspec
go test ./...
```

## Лицензия

MIT. См. [LICENSE](LICENSE).
