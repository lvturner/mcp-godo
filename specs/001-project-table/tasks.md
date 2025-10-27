# Implementation Tasks: Project Table for Todo Items

**Feature**: 001-project-table  
**Branch**: `001-project-table`  
**Date**: 2025-10-25  

Generated from spec.md and plan.md analysis.

---

## Phase 1: Database Schema and Models

### Task 1.1: Create Project Database Schema
**Priority**: P1 | **Parallel**: No  
**Description**: Create database migration for projects table with proper constraints and indexes  
**Files**: `migrations/0003_add_projects_table.sql`, `migrations/0003_add_projects_table.down.sql`  
**Requirements**: FR-001, FR-010  
**Acceptance**: Migration creates projects table with id, name, description, created_at fields and unique constraint on name

### Task 1.2: Update Todo Table Schema  
**Priority**: P1 | **Parallel**: No  
**Description**: Add project_id foreign key to existing todos table with proper indexing  
**Files**: `migrations/0004_update_todos_add_project_id.sql`, `migrations/0004_update_todos_add_project_id.down.sql`  
**Requirements**: FR-006, FR-007  
**Acceptance**: Todos table has optional project_id foreign key with proper indexes

### Task 1.3: Create Project Model Structure
**Priority**: P1 | **Parallel**: No  
**Description**: Define Project struct with validation and database mapping  
**Files**: `pkg/todo/project.go`  
**Requirements**: FR-001, FR-002  
**Acceptance**: Project struct includes all required fields with proper JSON tags and validation

### Task 1.4: Update Todo Model for Project Association
**Priority**: P1 | **Parallel**: No  
**Description**: Modify TodoItem struct to include optional ProjectID field  
**Files**: `pkg/todo/todo.go`  
**Requirements**: FR-006, FR-007  
**Acceptance**: TodoItem struct has optional ProjectID field with proper validation

---

## Phase 2: Database Service Layer

### Task 2.1: Implement Project Service Interface
**Priority**: P1 | **Parallel**: No  
**Description**: Create ProjectService interface with CRUD operations  
**Files**: `pkg/todo/project.go`  
**Requirements**: FR-002, FR-003, FR-004, FR-005  
**Acceptance**: Interface defines CreateProject, GetProject, UpdateProject, DeleteProject, ListProjects methods

### Task 2.2: Implement SQLite Project Service
**Priority**: P1 | **Parallel**: No  
**Description**: Create SQLite implementation of ProjectService with full CRUD operations  
**Files**: `pkg/todo/project_sqlite.go`  
**Requirements**: FR-002, FR-003, FR-004, FR-005, FR-010  
**Acceptance**: All CRUD operations work correctly with proper error handling and unique name validation

### Task 2.3: Implement MySQL Project Service
**Priority**: P1 | **Parallel**: No  
**Description**: Create MySQL implementation of ProjectService with full CRUD operations  
**Files**: `pkg/todo/project_mysql.go`  
**Requirements**: FR-002, FR-003, FR-004, FR-005, FR-010  
**Acceptance**: All CRUD operations work correctly with proper error handling and unique name validation

### Task 2.4: Update Todo Service for Project Operations
**Priority**: P1 | **Parallel**: No  
**Description**: Modify existing todo services to support project assignment and filtering  
**Files**: `pkg/todo/todo_sqlite.go`, `pkg/todo/todo_mysql.go`  
**Requirements**: FR-006, FR-007, FR-008  
**Acceptance**: Todo CRUD operations support optional project_id parameter and project-based filtering

---

## Phase 3: MCP Tool Handlers

### Task 3.1: Create Project MCP Tool Handlers
**Priority**: P1 | **Parallel**: No  
**Description**: Implement MCP tool handlers for project CRUD operations  
**Files**: `pkg/handler/project_handler.go`  
**Requirements**: FR-002, FR-003, FR-004, FR-005  
**Acceptance**: Handlers for create_project, get_project, update_project, delete_project, list_projects work correctly

### Task 3.2: Update Todo MCP Tool Handlers
**Priority**: P1 | **Parallel**: No  
**Description**: Modify existing todo handlers to support project assignment  
**Files**: `pkg/handler/todo_handler.go`  
**Requirements**: FR-006, FR-007  
**Acceptance**: Todo handlers accept optional project_id parameter and handle project assignment correctly

### Task 3.3: Add Project Todo Query Handlers
**Priority**: P2 | **Parallel**: No  
**Description**: Create handlers for querying todos by project  
**Files**: `pkg/handler/project_handler.go`  
**Requirements**: FR-008  
**Acceptance**: Handler for get_project_todos returns all todos for a specific project

---

## Phase 4: Main Application Integration

### Task 4.1: Register Project MCP Tools
**Priority**: P1 | **Parallel**: No  
**Description**: Add project-related MCP tools to the main application  
**Files**: `cmd/mcp-godo/main.go`  
**Requirements**: All functional requirements  
**Acceptance**: All project MCP tools are registered and accessible to LLM agents

### Task 4.2: Update Test Database Setup
**Priority**: P1 | **Parallel**: No  
**Description**: Modify test database setup to include projects table  
**Files**: `sql/setup_test_db.sql`  
**Requirements**: Testing requirements  
**Acceptance**: Test database setup includes projects table and sample data

---

## Phase 5: Testing

### Task 5.1: Create Project Service Unit Tests
**Priority**: P1 | **Parallel**: No  
**Description**: Write comprehensive unit tests for ProjectService implementations  
**Files**: `pkg/todo/project_test.go`  
**Requirements**: All functional requirements  
**Acceptance**: >90% code coverage for project service methods with all edge cases tested

### Task 5.2: Create Project Handler Unit Tests
**Priority**: P1 | **Parallel**: No  
**Description**: Write unit tests for MCP project handlers  
**Files**: `pkg/handler/project_handler_test.go`  
**Requirements**: All functional requirements  
**Acceptance**: All handler methods tested with proper mock setup

### Task 5.3: Create Integration Tests
**Priority**: P2 | **Parallel**: No  
**Description**: Write integration tests for project-todo workflow  
**Files**: `tests/integration/project_todo_test.go`  
**Requirements**: FR-006, FR-007, FR-008  
**Acceptance**: End-to-end tests verify project creation, todo assignment, and project deletion scenarios

### Task 5.4: Test Edge Cases
**Priority**: P2 | **Parallel**: No  
**Description**: Test edge cases from spec: duplicate names, deleted projects, orphaned todos  
**Files**: Various test files  
**Requirements**: Edge cases from spec.md  
**Acceptance**: All edge cases have corresponding test cases with expected behavior verified

---

## Phase 6: Documentation

### Task 6.1: Update API Documentation
**Priority**: P2 | **Parallel**: No  
**Description**: Document new MCP tools and updated existing ones  
**Files**: `specs/001-project-table/contracts/mcp-tools-plan.md`  
**Requirements**: Documentation requirements  
**Acceptance**: All project-related MCP tools documented with examples

### Task 6.2: Create Quick Start Guide
**Priority**: P2 | **Parallel**: No  
**Description**: Write quick start guide for project functionality  
**Files**: `specs/001-project-table/quickstart.md`  
**Requirements**: Documentation requirements  
**Acceptance**: Guide shows how to create projects and assign todos with examples

---

## Success Criteria Verification

### Task 7.1: Performance Testing
**Priority**: P2 | **Parallel**: No  
**Description**: Verify performance meets <200ms p95 response time requirement  
**Files**: Performance test scripts  
**Requirements**: SC-002  
**Acceptance**: All operations complete within performance constraints under load

### Task 7.2: User Experience Testing
**Priority**: P2 | **Parallel**: No  
**Description**: Verify user success rates for project operations  
**Files**: UX test scenarios  
**Requirements**: SC-003, SC-004  
**Acceptance**: 95% success rate for project assignment, 90% for todo creation without project

---

## Completed Tasks (Bug Fixes)

### Bug Fix: Project Deletion Active Todo Handling
**Status**: ✅ Completed | **Date**: 2025-10-27
**Bug**: Project deletion was not properly handling active todos - they should have their project_id set to null when a project is deleted
**Files Modified**:
- [`pkg/todo/project_mariadb.go`](pkg/todo/project_mariadb.go) - Added transaction and active todo update logic
- [`pkg/todo/project_sqlite.go`](pkg/todo/project_sqlite.go) - Added transaction and active todo update logic
- [`pkg/handler/project_handler.go`](pkg/handler/project_handler.go) - Fixed placeholder implementation to actually call project service
- [`tests/unit/project_deletion_test.go`](tests/unit/project_deletion_test.go) - Created comprehensive unit tests
**Solution**: Added database transactions to atomically update active todos (completed_at IS NULL) to set project_id = NULL before deleting the project
**Testing Results**: All tests pass (existing 29 unit tests + 11 contract tests + 7 new project deletion tests)

---

## Task Summary

**Total Tasks**: 18  
**P1 Tasks**: 14 (Critical path)  
**P2 Tasks**: 4 (Enhancements and documentation)  
**Parallel Tasks**: 0 (All tasks have dependencies)

**Estimated Effort**: 3-4 development days  
**Critical Path**: Tasks 1.1 → 1.2 → 1.3 → 1.4 → 2.1 → 2.2 → 2.3 → 2.4 → 3.1 → 3.2 → 4.1 → 5.1 → 5.2

**Dependencies**: 
- Database schema must be complete before service layer
- Service layer must be complete before MCP handlers  
- MCP handlers must be complete before integration testing
- All core functionality must be complete before documentation