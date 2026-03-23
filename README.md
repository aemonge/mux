# mux

tmux 세션을 터미널에서 빠르게 탐색하고 관리하는 TUI 도구입니다.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)

## 기능

- **세션 목록** — 활성/비활성 세션을 최근 활동 순으로 정렬하여 표시
- **실시간 프리뷰** — 선택한 세션의 터미널 출력을 우측 패널에 실시간으로 표시 (500ms 주기 갱신)
- **AI CLI 감지** — 세션에서 `claude`, `codex`, `aider`, `gemini` 등의 AI CLI가 실행 중이면 배지로 표시
- **세션 생성/삭제/이름 변경** — TUI 내에서 모든 세션 관리 작업 가능
- **퀵 필터** — `/` 키로 세션 이름 또는 경로를 실시간 필터링
- **즉시 연결** — `Enter` 키로 선택한 세션에 바로 attach (tmux 내부에서는 `switch-client` 사용)

## 설치

```bash
git clone https://github.ecodesamsung.com/euteum-park/mux.git
cd mux
go build -o mux .
```

빌드된 바이너리를 PATH에 추가:

```bash
mv mux /usr/local/bin/
```

## 사용법

```bash
mux
```

### 팝업 모드 (추천)

tmux 안에서 작업 중일 때, 어떤 프로그램(claude, codex 등)이 실행 중이어도 키 하나로 mux를 오버레이로 띄울 수 있습니다.

**설정:**

```bash
# prefix + m 으로 팝업 열기 (기본값)
mux setup-keybind

# 다른 키로 변경 가능
mux setup-keybind Space

# tmux 설정 리로드
tmux source-file ~/.tmux.conf
```

**사용:**

`Ctrl+b` → `m` (또는 설정한 키)으로 mux 팝업이 열리고, 세션 선택 또는 `q`로 닫으면 원래 작업으로 돌아갑니다.

수동으로 팝업을 열 수도 있습니다:

```bash
mux popup
```

> **참고:** tmux 3.2 이상 필요 (`tmux -V`로 확인)

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

## 레이아웃

```
⚡ tmux sessions (3)
┌─────────────────┐┌──────────────────────────────────────┐
│ ● my-project    ││ [ my-project ]  ~/dev/project  ✦ claude│
│   dev-server    ││ ─────────────────────────────────────  │
│   dotfiles      ││ ...터미널 출력 미리보기...               │
└─────────────────┘└──────────────────────────────────────┘
↑↓/jk navigate  •  enter attach  •  n new  •  x kill  •  r rename  •  / filter  •  q quit
```

## 요구사항

- Go 1.21+
- tmux (팝업 모드는 3.2+)

## 의존성

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI 프레임워크
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI 컴포넌트
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — 스타일링
