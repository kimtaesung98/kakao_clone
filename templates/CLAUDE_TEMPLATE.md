# Blueprint Framework — Project Rules

> 이 파일은 `CLAUDE.md`로 복사해 사용합니다.  
> Claude Code는 이 파일을 **매 세션 자동으로 읽습니다** — 항상 최신 상태로 유지하세요.

---

## 세션 시작 프로토콜 (매번 실행)

```
AGENT_BOOTSTRAP.md 존재?
    YES → 전체 읽기 → AGENT_STATE.md 생성 → BOOTSTRAP 삭제
    NO  → AGENT_STATE.md 전체 읽기
              ↓
         Resolved Requirements 확인 (중복 작업 방지)
         Current Task / Next Action 확인 (이어받기)
         Active Rules 확인 (이번 세션 제약)
```

---

## 핵심 규칙 (모든 세션에 적용)

**폴더 구조**
- 모든 폴더에 `README.md` 존재 필수
- Progress Tracker · Roadmap은 **최상위 폴더 README.md에만** 존재
- 하위 폴더 README.md는 `Context` + `DFD` + `Agent Control`만 보유

**DFD**
- 루트 = Level 0, 1단계 하위 = Level 1, 이하 순차 증가
- 모든 DFD에 `> Decomposed from: [상위 폴더] Level N — [버블명] 버블` 명시
- 하위 DFD 외부 입출력 = 상위 DFD 해당 버블 화살표와 일치

**작업 순서**
1. 루트 README.md Tracker → 작업 Feature 확인
2. 담당 폴더 README.md → Context · DFD · Agent Control 확인
3. 코드 작업
4. 루트 Tracker 업데이트 → AGENT_STATE.md 갱신

---

## 세션 종료 시 필수 갱신

- [ ] 루트 README.md Progress Tracker 동기화
- [ ] AGENT_STATE.md 세션 로그 기록
- [ ] 구조 변경 시 DFD 연쇄 업데이트 (하위 → 루트)

---

## 상세 문서 위치

| 문서 | 내용 |
|------|------|
| `AGENT_STATE.md` | 현재 진행 상태 · 누적 요구사항 · 세션 로그 |
| `BLUEPRINT.md` | 프레임워크 전체 명세 (세션 중 읽지 않음) |
| `CONTRIBUTING.md` | 상세 프로토콜 · DFD 드리프트 감지 · 팀 규칙 |
