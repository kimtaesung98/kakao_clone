# 새 폴더 생성 체크리스트

> 에이전트 또는 개발자가 새 폴더를 생성할 때 이 체크리스트를 순서대로 완료하세요.

---

## Step 1 — 폴더 유형 결정 (템플릿 선택)

먼저 이 폴더가 **소유자(Top)** 인지 **실행자(Sub)** 인지 결정합니다.

| 질문 | Yes | No |
|------|-----|----|
| 이 폴더가 Progress Tracker와 Roadmap을 정의하는 최상위인가? | → `README_TEMPLATE_TOP.md` 사용 | → `README_TEMPLATE_SUB.md` 사용 |

> **소유자(Top)**: 프로젝트 루트, 또는 독립적인 도메인/서비스의 최상위 폴더  
> **실행자(Sub)**: 소유자 아래에 위치하며 DFD로 구현 흐름을 상세화하는 모든 하위 폴더

---

## Step 2 — README_TEMPLATE_TOP.md 사용 시 체크리스트

- [ ] 파일 복사 후 `[프로젝트/모듈명]` 헤더 교체
- [ ] Overview: 시스템 전체 책임 1~2문장 작성
- [ ] DFD Level 0: 외부 액터와 시스템 경계를 단일 버블로 표현
- [ ] Tech Stack: 시스템 전체에서 사용하는 스택 나열
- [ ] Agent Control: 전체 시스템에 적용되는 최상위 제약 정의
- [ ] Sub-folder Map: 하위 폴더와 담당 Feature 매핑 테이블 작성
- [ ] Progress Tracker: 구현할 전체 Feature 목록 초기화 (`⏳ Pending`)
- [ ] Next Roadmap: 첫 번째 작업 항목과 담당 폴더 경로 기재
- [ ] 모든 `[[ ]]` 플레이스홀더가 실제 내용으로 교체되었는가?

---

## Step 3 — README_TEMPLATE_SUB.md 사용 시 체크리스트

- [ ] 파일 복사 후 `[폴더명]` 헤더 교체
- [ ] **상위 작업 연결** 줄: 상위 README.md의 어느 Feature를 담당하는지 명시
- [ ] **Context 섹션**: 상위 DFD 버블명, 입력 데이터, 출력 데이터 채우기
- [ ] **DFD 레벨 결정**: 폴더 깊이에 맞는 레벨 번호와 레벨명 선택
- [ ] **`Decomposed from:` 줄**: 상위 폴더 경로 + 레벨 + 버블명 명시
- [ ] DFD mermaid: Context의 입력/출력과 일치하도록 작성
- [ ] Tech Stack: 상위에서 명시되지 않은 추가 스택만 기재
- [ ] Agent Control: 상위 폴더 경로 명시 + 추가 제약 정의
- [ ] 모든 `[[ ]]` 플레이스홀더가 실제 내용으로 교체되었는가?
- [ ] **상위 DFD와 이 DFD의 입출력이 일치하는가?** (핵심 검증)

---

## Step 4 — 상위 폴더 업데이트

새 하위 폴더를 만든 뒤 **상위 폴더의 README.md도 업데이트**해야 합니다:

- [ ] 상위 `README_TEMPLATE_TOP.md` 기반 파일의 **Sub-folder Map**에 새 폴더 행 추가
- [ ] 상위 DFD에 이 폴더를 나타내는 버블이 없다면 추가
- [ ] 상위 Progress Tracker에 이 폴더가 담당할 Feature가 없다면 추가

---

## 에이전트를 위한 최종 검증 (5가지 질문)

새 폴더의 `README.md`를 작성한 후 아래 질문에 답할 수 있어야 합니다:

| # | 질문 | 확인 섹션 |
|---|------|-----------|
| 1 | 이 폴더는 상위 Progress Tracker의 어느 Feature를 담당하는가? | Context |
| 2 | 상위 DFD의 어느 버블을 이 DFD가 분해하고 있는가? | `Decomposed from:` |
| 3 | 이 폴더로 들어오는 데이터와 나가는 데이터는 무엇인가? | Context + DFD |
| 4 | 이 폴더에서 절대 해서는 안 되는 일은 무엇인가? | Agent Control > 금지 |
| 5 | 이 폴더의 입출력이 상위 DFD의 화살표와 일치하는가? | DFD 일관성 검토 |

이 다섯 가지 질문에 답할 수 없다면 `README.md`가 불완전한 것입니다.
