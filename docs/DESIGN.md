# DESIGN

This documentation basically describes the architecture of the application.

## Components and Collaborations

- **User**: basic user that has username and password, with limited permission
  to application.
- **Student**: a kind of user who can join courses, submit homework and view
  commited homework.
- **Teacher**: a kind of user who can setup courses, publish homework, view
  and download commited homework.
- **Admininstrator**: has full access to the system.
- **Course**: container for teachers, students and assignments.
- **Assignment**: attached to courses. There may be multiple assignments in a
  course. But there cannot be more than one uploaded file for each person.

## Key Points

### Registration
Everyone (including teachers and students) should use real identity to register
to the system.

### Assignments
When receiving assignments, submitted homework can be automatically renamed to a
regularized form.

## Implementation


