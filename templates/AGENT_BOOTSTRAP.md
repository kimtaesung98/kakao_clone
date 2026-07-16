# AGENT BOOTSTRAP — 최초 1회 읽고 삭제

> ⚠️ **이 파일은 일회용입니다.**  
> 에이전트는 이 파일을 읽은 즉시 아래 순서를 따르고, **마지막에 반드시 이 파일을 삭제**해야 합니다.  
> 이후 세션에서는 `AGENT_STATE.md`만 읽으면 됩니다.

---

## 에이전트 실행 순서 (이 파일을 열었을 때)

```
1. 이 파일 전체를 읽는다
2. 루트 README.md → Progress Tracker 현재 상태 확인
3. CLAUDE.md가 없으면 templates/CLAUDE_TEMPLATE.md를 복사해 루트에 CLAUDE.md 생성
4. AGENT_STATE.md 파일을 생성한다  ← 아래 §5 초기화 지시 참고
5. 이 파일(AGENT_BOOTSTRAP.md)을 삭제한다
6. 사용자에게 보고: "초기화 완료. CLAUDE.md + AGENT_STATE.md 생성됨. 이후 세션은 자동 로딩됩니다"
```

---

## §1. 이 프로젝트의 Blueprint Framework 적용 구조

```
[[ 프로젝트 루트 경로 ]]/
├── README.md                  ← Level 0 DFD + Progress Tracker (소유자)
├── AGENT_STATE.md             ← 에이전트 상태 파일 (이 파일 삭제 후 생성)
├── CONTRIBUTING.md            ← 에이전트 행동 강령
└── [[ 주요 폴더들 ]]/
    └── README.md              ← Level N DFD + Context (실행자, Tracker 없음)
```

**소유자 폴더** (Progress Tracker 보유): `[[ 경로 나열 ]]`  
**실행자 폴더** (DFD만 보유): `[[ 경로 나열 ]]`

---

## §2. DFD 레벨링 핵심 규칙 (요약)

| 폴더 깊이 | DFD 레벨 | 헤더 표기 |
|-----------|----------|-----------|
| 루트 | Level 0 — Context | 시스템 전체, 외부 액터만 |
| 1단계 하위 | Level 1 — Main Processes | 루트 버블 분해 |
| 2단계 하위 | Level 2 — Detailed Processes | Level 1 버블 분해 |
| 3단계+ | Level 3+ — Function Level | 함수 수준 |

- 각 DFD는 `> Decomposed from: [상위 폴더] Level N — [버블명] 버블` 명시 필수
- 하위 DFD 외부 입출력 = 상위 DFD 해당 버블 화살표와 반드시 일치
- 새 외부 의존성 추가 시 Level N → N-1 → 0 순으로 상위 DFD 연쇄 업데이트

---

## §3. 폴더 소유권 모델 (요약)

```
최상위(소유자) README.md         하위(실행자) README.md
─────────────────────────────    ─────────────────────────────
✅ DFD Level 0                   ✅ Context (상위 Tracker 항목 참조)
✅ Agent Control (전역)           ✅ DFD Level N (버블 분해)
✅ Sub-folder Map                 ✅ Agent Control (상속 + 추가)
✅ Progress Tracker  ← 유일       ❌ Progress Tracker 없음
✅ Next Roadmap      ← 유일       ❌ Next Roadmap 없음
```

**에이전트 작업 순서:**  
① 루트 README.md → Tracker에서 Feature 확인  
② 담당 하위 폴더 README.md → Context + DFD로 구현 흐름 파악  
③ Agent Control 금지 규칙 확인 → 코드 작업  
④ 루트 README.md Tracker 업데이트 (🔄 → ✅)  
⑤ AGENT_STATE.md 업데이트

---

## §4. 이 프로젝트 전역 Agent Control (요약)

> 이 섹션을 프로젝트 루트 README.md의 Agent Control에서 복사해 채우세요.

**허용:**
- [[ 전역 허용 규칙 ]]

**금지:**
- [[ 전역 금지 규칙 ]]

**필수:**
- 작업 완료마다 루트 README.md Progress Tracker 업데이트
- 구조 변경 시 해당 레벨부터 루트까지 DFD 연쇄 업데이트
- 작업 완료마다 AGENT_STATE.md 세션 로그 기록

---

## §5. AGENT_STATE.md 초기화 지시

이 섹션을 읽고 `AGENT_STATE.md`를 생성하세요.  
아래 템플릿에서 `[[ ]]` 부분을 현재 프로젝트 상태로 채워 생성합니다.

```markdown
# Agent State

> 이 파일은 에이전트가 매 세션 시작 시 읽고, 작업 완료 후 갱신합니다.
> BLUEPRINT.md와 AGENT_BOOTSTRAP.md는 읽지 않습니다 — 이 파일이 전체 맥락을 제공합니다.

---

## 누적 사용자 요구사항 (Resolved Requirements)

> 사용자가 요청하고 반영 완료된 항목을 기록합니다.
> 에이전트는 이 목록을 보고 이미 처리된 요구사항을 재처리하지 않습니다.

| # | 요구사항 | 반영 날짜 | 변경된 파일 |
|---|----------|-----------|-------------|
| - | (초기화 — 아직 없음) | - | - |

---

## 현재 작업 상태 (Progress Snapshot)

> 루트 README.md Progress Tracker에서 복사. 작업 완료마다 갱신.

| Feature | Status | Owner Folder |
|---------|--------|--------------|
| [[ Feature A ]] | ⏳ Pending | `/[[폴더]]` |

---

## 다음 할 일 (Next Action)

> 루트 README.md Next Roadmap에서 복사. 완료 항목은 삭제.

1. [[ 다음 작업 항목 ]]

---

## 핵심 규칙 스냅샷 (Active Rules)

> 전역 Agent Control 핵심만 3~7개. 상세는 각 폴더 README.md 참조.

- [[ 핵심 규칙 1 ]]
- [[ 핵심 규칙 2 ]]

---

## DFD 드리프트 상태 (Drift Status)

> 구조적 변경 후 에이전트가 업데이트합니다.  
> ⚠️ 드리프트 의심 항목은 다음 작업 전 반드시 검증합니다.

| 폴더 | DFD 레벨 | 마지막 검증일 | 상태 |
|------|----------|--------------|------|
| `/` | Level 0 | [[ 오늘 날짜 ]] | ✅ 최신 |

---

## 세션 로그 (Session Log)

| 날짜 | 수행 작업 | 변경 파일 | 다음 세션 인계 사항 |
|------|-----------|-----------|---------------------|
| [[ 오늘 날짜 ]] | 초기화 (BOOTSTRAP 읽고 STATE 생성) | AGENT_STATE.md 생성, AGENT_BOOTSTRAP.md 삭제 | - |
```

---

## 이 파일 삭제 확인

AGENT_STATE.md 생성이 완료되었으면 이 파일을 삭제하세요.  
삭제 명령: `rm AGENT_BOOTSTRAP.md` (또는 에디터에서 파일 삭제)

> 삭제 후 이 파일은 git에서도 제거됩니다.  
> 이후 세션에서는 절대 이 파일을 다시 만들지 않습니다.
