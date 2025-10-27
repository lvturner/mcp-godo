---
description: "Task list for category feature implementation"
---

# Tasks: Category Feature for Todo Items

**Input**: Design documents from `/specs/001-create-category-feature/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: The examples below include test tasks. Tests are OPTIONAL - only include them if explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `src/`, `tests/` at repository root
- **Web app**: `backend/src/`, `frontend/src/`
- **Mobile**: `api/src/`, `ios/src/` or `android/src/`
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create database migrations for categories table
- [x] T002 Create database migrations for category_id foreign key in todos table
- [x] T003 [P] Update existing todo models to support optional category_id

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Implement Category entity and interfaces in pkg/todo/category.go
- [x] T005 [P] Implement Category SQLite repository in pkg/todo/category_sqlite.go
- [x] T006 [P] Implement Category MariaDB repository in pkg/todo/category_mariadb.go
- [x] T007 Update Todo entity to support optional category relationship
- [x] T008 Create base category service interface

**Checkpoint**: ✅ Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Create and Manage Categories (Priority: P1) 🎯 MVP

**Goal**: Allow users to create and manage categories for organizing todo items

**Independent Test**: Can be fully tested by creating a category, viewing it, and verifying it appears in the category list.

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T009 [P] [US1] Unit test for category creation in tests/unit/category_test.go
- [x] T010 [P] [US1] Unit test for category validation in tests/unit/category_validation_test.go
- [x] T011 [P] [US1] Contract test for create_category MCP method in tests/contract/create_category_test.go
- [x] T012 [P] [US1] Contract test for get_categories MCP method in tests/contract/get_categories_test.go

### Implementation for User Story 1

- [x] T013 [P] [US1] Implement create_category service method in pkg/todo/category.go
- [x] T014 [P] [US1] Implement get_categories service method in pkg/todo/category.go
- [x] T015 [P] [US1] Implement update_category service method in pkg/todo/category.go
- [x] T016 [P] [US1] Implement delete_category service method in pkg/todo/category.go
- [x] T017 [US1] Implement create_category MCP handler in pkg/handler/category_handler.go
- [x] T018 [US1] Implement get_categories MCP handler in pkg/handler/category_handler.go
- [x] T019 [US1] Implement update_category MCP handler in pkg/handler/category_handler.go
- [x] T020 [US1] Implement delete_category MCP handler in pkg/handler/category_handler.go
- [x] T021 [US1] Add validation for unique category names
- [x] T022 [US1] Add error handling for category operations

**Checkpoint**: ✅ User Story 1 is fully functional and testable independently

---

## Phase 4: User Story 2 - Assign Todos to Categories (Priority: P2)

**Goal**: Allow users to assign todo items to categories for better organization

**Independent Test**: Can be fully tested by creating a todo item and assigning it to a category, then verifying it appears in the category's todo list.

### Tests for User Story 2

- [x] T023 [P] [US2] Unit test for todo-category assignment in tests/unit/todo_category_test.go
- [x] T024 [P] [US2] Contract test for assign_todo_to_category MCP method in tests/contract/assign_todo_test.go
- [x] T025 [P] [US2] Contract test for get_category_todos MCP method in tests/contract/get_category_todos_test.go

### Implementation for User Story 2

- [x] T026 [P] [US2] Implement assign_todo_to_category service method in pkg/todo/todo.go
- [x] T027 [P] [US2] Implement get_todos_by_category service method in pkg/todo/todo.go
- [x] T028 [US2] Implement remove_todo_from_category service method in pkg/todo/todo.go
- [x] T029 [US2] Implement assign_todo_to_category MCP handler in pkg/handler/category_handler.go
- [x] T030 [US2] Implement get_category_todos MCP handler in pkg/handler/category_handler.go
- [x] T031 [US2] Implement remove_todo_from_category MCP handler in pkg/handler/category_handler.go
- [x] T032 [US2] Update existing todo creation to support optional category assignment

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Create Todos Without Categories (Priority: P3)

**Goal**: Allow users to create todo items without specifying a category for flexibility

**Independent Test**: Can be fully tested by creating a todo item without selecting a category, then verifying it appears in a default "uncategorized" or "no category" view.

### Tests for User Story 3

- [x] T033 [P] [US3] Unit test for todo creation without category in tests/unit/todo_no_category_test.go
- [x] T034 [P] [US3] Contract test for get_uncategorized_todos MCP method in tests/contract/get_uncategorized_todos_test.go

### Implementation for User Story 3

- [x] T035 [P] [US3] Update create_todo service method to handle optional category_id
- [x] T036 [P] [US3] Implement get_uncategorized_todos service method in pkg/todo/todo.go
- [x] T037 [US3] Update create_todo MCP handler to support optional category parameter
- [x] T038 [US3] Implement get_uncategorized_todos MCP handler in pkg/handler/category_handler.go
- [x] T039 [US3] Add validation to ensure todos can exist without categories

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T040 [P] Documentation updates in docs/
- [x] T041 Code cleanup and refactoring
- [x] T042 Performance optimization for category queries
- [x] T043 [P] Additional unit tests for edge cases in tests/unit/
- [x] T049 [P] Comprehensive unit tests for assign_todo_to_category fix in tests/unit/assign_todo_to_category_test.go
- [x] T050 Fix assign_todo_to_category handler bug - ensure todo service is called for database updates
- [x] T051 Run complete test suite verification - all assign_todo_to_category tests pass (13/13)
- [x] T044 Security hardening for category operations
- [x] T045 Run quickstart.md validation
- [x] T046 Add logging for category operations
- [x] T047 Handle category deletion edge cases (what happens to associated todos)
- [x] T048 Expose category MCP tools in cmd/mcp-godo/main.go

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 for category existence
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - No dependencies on other stories

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Models before services
- Services before endpoints
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 3
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence