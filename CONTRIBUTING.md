# Contributing — 에이전트 행동 강령

> 이 파일은 **Claude Code를 포함한 모든 AI 에이전트**가 이 레포지토리에서 작업할 때 반드시 따라야 하는 규칙입니다.

---

## 세션 시작 시 — 어떤 파일을 읽어야 하는가?

```
AGENT_BOOTSTRAP.md가 존재하는가?
        │
       YES → §1 최초 초기화 프로토콜 실행
        │
        NO → AGENT_STATE.md가 존재하는가?
                    │
                   YES → §2 일반 세션 프로토콜 실행
                    │
                    NO → 사용자에게 문의 (비정상 상태)
```

---

## §1. 최초 초기화 프로토콜 (AGENT_BOOTSTRAP.md 존재 시)

> 이 프로토콜은 프로젝트 생애 **딱 한 번**만 실행됩니다.

1. `AGENT_BOOTSTRAP.md` 전체를 읽는다
2. 루트 `README.md`의 Progress Tracker와 Next Roadmap을 읽는다
3. 루트에 `CLAUDE.md`가 없으면 `templates/CLAUDE_TEMPLATE.md`를 복사해 생성한다
4. `AGENT_BOOTSTRAP.md §5`의 지시에 따라 `AGENT_STATE.md`를 생성한다
5. `AGENT_BOOTSTRAP.md`를 **삭제**한다
6. 사용자에게 보고: `"초기화 완료. CLAUDE.md + AGENT_STATE.md 생성. 이후 세션은 CLAUDE.md가 자동 로딩됩니다."`

**이후 모든 세션은 §2 프로토콜만 따릅니다. BLUEPRINT.md는 읽지 않습니다.**

---

## §2. 일반 세션 프로토콜 (매 세션, AGENT_STATE.md 존재 시)

### 세션 시작 (읽기)

1. `AGENT_STATE.md` 전체를 읽는다
   - **누적 사용자 요구사항** → 이미 처리된 것 확인, 중복 방지
   - **현재 작업 상태** → 이어받을 Feature 파악
   - **다음 할 일** → 이번 세션의 작업 범위 확인
   - **핵심 규칙 스냅샷** → 이번 세션에 적용할 제약 확인
   - **세션 로그** → 직전 세션 인계 사항 확인

2. 사용자의 이번 세션 요청을 받는다

3. 요청이 **누적 사용자 요구사항**에 이미 있으면 → 완료 여부 확인 후 불필요한 재작업 방지

### 작업 중

4. 작업 대상 폴더의 `README.md`를 읽는다 (DFD + Agent Control 확인)
5. Agent Control **금지 규칙** 위반 여부 확인 후 코드 작성
6. 구조 변경 시 → 해당 레벨부터 루트까지 DFD 연쇄 업데이트

### 세션 종료 (갱신)

7. 루트 `README.md` **Progress Tracker** 업데이트 (🔄 시작, ✅ 완료)
8. 루트 `README.md` **Next Roadmap** 갱신
9. `AGENT_STATE.md`를 아래 순서로 갱신한다:

```
a. 누적 사용자 요구사항  → 이번 세션에서 처리한 새 요구사항 행 추가
b. 현재 작업 상태        → 루트 Tracker 내용과 동기화
c. 다음 할 일            → 루트 Roadmap 내용과 동기화
d. 핵심 규칙 스냅샷      → 새로 추가/변경된 규칙 반영
e. 세션 로그             → 오늘 날짜, 수행 작업, 변경 파일, 인계 사항 기록
                           (10개 초과 시 가장 오래된 행 삭제)
```

---

## §3. AGENT_STATE.md 갱신 상세 규칙

### 누적 사용자 요구사항 — 언제 추가하는가

- 사용자가 새로운 구조적 요청을 하고 이를 반영 완료했을 때
- 새로운 Agent Control 규칙이 추가되었을 때
- 폴더 구조가 변경되었을 때

> 사소한 버그 수정이나 단순 코드 변경은 추가하지 않습니다.  
> "이 프로젝트의 방향과 규칙에 영향을 주는 결정"만 기록합니다.

### Progress Snapshot — 동기화 기준

- 루트 `README.md` Tracker와 **항상 동일**해야 합니다
- 세션 종료 전 두 파일을 비교해 불일치가 없는지 확인합니다

### 세션 로그 — 기록 형식

```markdown
| 2026-05-18 | DFD 레벨링 적용, 폴더 소유권 모델 정의 | BLUEPRINT.md, templates/* | 하위 폴더 DFD 연쇄 업데이트 미완료 |
```

---

## §4. 절대 금지 행동

| 금지 행동 | 이유 |
|-----------|------|
| 매 세션 `BLUEPRINT.md` 전체 읽기 | 불필요한 토큰 낭비 |
| `AGENT_BOOTSTRAP.md` 삭제 없이 세션 종료 | 다음 세션이 중복 초기화 실행 |
| `AGENT_STATE.md` 미갱신 세션 종료 | 맥락 단절, 다음 세션이 처음부터 재파악 |
| `AGENT_STATE.md`에 코드 내용 기록 | 파일 비대화, 토큰 낭비 |
| 이미 처리된 요구사항 재처리 | 누적 사용자 요구사항 미확인에서 발생 |
| `Progress Snapshot`과 루트 Tracker 불일치 유지 | 맥락 오염 |

---

## §5. 커밋 메시지 규칙

```
<type>(<scope>): <description>
```

| Type | 사용 시점 |
|------|-----------|
| `feat` | 새 기능 추가 |
| `fix` | 버그 수정 |
| `docs` | 문서 수정 (README, DFD, AGENT_STATE 등) |
| `refactor` | 기능 변경 없는 구조 개선 |
| `chore` | 빌드, 의존성 등 기타 |

**AGENT_STATE.md 갱신은 기능 커밋에 포함하거나 `docs` 커밋으로 분리합니다.**

예: `feat(auth): add JWT refresh endpoint + update agent state`

---

## §6. DFD 드리프트 감지 프로토콜

DFD는 시간이 지나면 실제 코드와 멀어집니다. 아래 트리거마다 검증을 실행합니다.

### 검증 트리거

- 5회 이상 커밋 후 해당 폴더 DFD 미갱신 상태
- 새 외부 라이브러리 또는 서비스 추가
- 폴더 간 데이터 전달 방식 변경
- `AGENT_STATE.md` Drift Status에 `⚠️ 드리프트 의심` 표시된 항목 존재

### 드리프트 감지 프롬프트

의심되는 폴더에 아래 프롬프트를 실행하세요:

```
[폴더 경로]의 실제 코드를 읽고,
현재 README.md의 DFD와 비교해서 불일치 항목을 찾아줘.
불일치가 있으면 DFD를 코드 기준으로 수정하고,
상위 폴더 DFD에도 영향이 있으면 함께 업데이트해줘.
```

### 드리프트 해결 후

- 해당 폴더 DFD 업데이트
- `AGENT_STATE.md` Drift Status 테이블의 마지막 검증일 갱신
- 상위 DFD에도 영향이 있었으면 연쇄 업데이트 확인

---

## §7. 멀티 에이전트 / 팀 협업 규칙

여러 사람 또는 여러 세션이 동시에 작업할 때의 충돌 방지 규칙입니다.

### AGENT_STATE.md 충돌 방지

`AGENT_STATE.md`는 **한 번에 한 세션만 수정**합니다.

```
작업 시작 전: git pull로 최신 STATE 동기화
작업 종료 후: STATE 갱신 → 즉시 커밋 → push
```

동시 수정이 발생해 충돌이 생기면:
1. 두 버전의 세션 로그를 모두 보존 (합치기)
2. Progress Snapshot은 루트 README.md Tracker 기준으로 재동기화
3. 누적 요구사항은 양쪽 항목을 모두 유지

### Progress Tracker 충돌 방지

같은 Feature를 두 세션이 동시에 수정하는 것을 막으려면,  
작업 시작 시 Feature 상태를 `🔄 In Progress`로 먼저 커밋해 점유를 선언합니다.

```
# 작업 선점 커밋
git add README.md && git commit -m "chore: claim [Feature명] In Progress"
git push
```

다른 세션은 `🔄 In Progress` 상태를 보고 해당 Feature 작업을 건너뜁니다.

### 폴더 소유 원칙

복수의 에이전트/개발자가 협업할 때 **한 폴더는 한 세션이 담당**합니다.  
동일 폴더를 동시에 수정해야 할 경우 git branch를 분리하고 PR로 통합합니다.
