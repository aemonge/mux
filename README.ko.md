# mux

**AI CLI 세션 간 전환을 빠르고 직관적으로.**

tmux 세션을 터미널에서 빠르게 탐색하고 관리하는 TUI 도구입니다.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

## 기능

- **실시간 프리뷰** — 선택한 세션의 터미널 출력을 우측 패널에 실시간으로 표시 (500ms 주기 갱신)
- **AI CLI 감지** — `claude`, `codex`, `aider`, `gemini` 등의 AI CLI가 실행 중이면 배지로 표시
- **Git 브랜치 표시** — 각 세션의 현재 브랜치를 표시, worktree는 `⌥⌥`로 구분
- **비용/토큰 추적** — Claude Code 세션의 토큰 사용량과 예상 비용을 실시간 표시 (설정 불필요)
- **상태바 위젯** — `mux status`로 tmux 상태바에 AI 세션 아이콘 표시
- **팝업 오버레이** — AI CLI 실행 중에도 키 하나로 mux를 띄워 세션 전환
- **세션 관리** — TUI 내에서 생성/삭제/이름 변경
- **퀵 필터** — `/` 키로 세션 이름 또는 경로를 실시간 필터링

## 설치

### 인터랙티브 설치 (추천)

```bash
curl -sSL https://raw.githubusercontent.com/lunemis/mux/main/install.sh | bash
```

### Homebrew

```bash
brew install lunemis/tap/mux
```

### 소스에서 빌드

```bash
git clone https://github.com/lunemis/mux.git
cd mux
make install   # /usr/local/bin에 설치
```

### Go install

```bash
go install github.com/lunemis/mux/cmd/mux@latest
```

## 사용법

```bash
mux
```

### 레이아웃

![Screenshot](assets/screenshot.png)

### 팝업 모드 (추천)

tmux 안에서 작업 중일 때, 어떤 프로그램이 실행 중이어도 키 하나로 mux를 오버레이로 띄울 수 있습니다.

```bash
# prefix + m 으로 팝업 열기 (기본값)
mux setup-keybind

# 다른 키로 변경 가능
mux setup-keybind Space

# tmux 설정 리로드
tmux source-file ~/.tmux.conf
```

`Ctrl+b` → `m`으로 mux 팝업이 열리고, 세션 선택 또는 `q`로 닫으면 원래 작업으로 돌아갑니다.

> **참고:** tmux 3.2 이상 필요 (`tmux -V`로 확인)

### 상태바 위젯

tmux 상태바에서 AI 세션 아이콘을 표시:

```bash
# ~/.tmux.conf에 추가
set -g status-right '#(mux status)'
```

AI 세션이 활성화되면 `✦ ◈` 같은 아이콘이 상태바에 표시됩니다.

### 키바인딩

| 키 | 동작 |
|---|---|
| `↑` / `k` | 위로 이동 |
| `↓` / `j` | 아래로 이동 |
| `g` / `G` | 처음 / 마지막으로 이동 |
| `Enter` | 선택한 세션에 attach |
| `n` | 새 세션 생성 |
| `r` | 세션 이름 변경 |
| `x` | 세션 삭제 (확인 후) |
| `/` | 세션 필터링 |
| `Esc` | 필터 초기화 / 모드 취소 |
| `q` / `Ctrl+C` | 종료 |

## 요구사항

- tmux (팝업 모드는 3.2+)
- Linux 또는 macOS

## 기여

[CONTRIBUTING.md](CONTRIBUTING.md)를 참고하세요.

## 라이선스

[MIT](LICENSE)
