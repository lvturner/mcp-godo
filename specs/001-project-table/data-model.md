# Data Model for Project Table Feature

## Entities

### Projects
- **id** (integer, primary key, auto-increment)
- **name** (string, unique, not null)
- **description** (text)
- **created_at** (timestamp, not null)
- **updated_at** (timestamp, not null)

### Todos
- **id** (integer, primary key, auto-increment)
- **title** (string, not null)
- **description** (text)
- **project_id** (integer, nullable)
- **completed_at** (datetime, default NULL)
- **due_date** (datetime, default NULL)
- **created_date** (datetime, not null, default current_timestamp())
- **reference_id** (integer, default NULL)

## Relationships

1. **Projects to Todos** (One-to-Many)
   - One project can have zero or more todos
   - A todo belongs to zero or one project
   - project_id field references Projects.id (foreign key constraint not enforced in current schema)

2. **Optional Association**
   - Todos can exist without being associated with a project
   - Projects can exist without having any todos

## Validation Rules

1. **Projects**
   - Name must be unique
   - Name cannot be empty
   - Created_at and updated_at timestamps are automatically managed

2. **Todos**
   - Title is required
   - Project_id references an existing project (if provided)
   - Completed_at is a datetime field indicating completion
   - Due_date is optional
   - Created_date is automatically managed
   - Reference_id is optional

## State Transitions

1. **Project Lifecycle**
   - Created → Deleted (soft delete or hard delete)

2. **Todo Lifecycle**
   - Created → Completed (completed_at timestamp set)
   - Completed → Uncompleted (completed_at timestamp cleared)

## Database Schema

The database schema follows the existing patterns in the codebase:
- SQLite/MySQL compatible
- Auto-incrementing primary keys
- Timestamp fields for audit trails
- Foreign key constraints for data integrity (project_id references Projects.id)
- Proper indexing for performance