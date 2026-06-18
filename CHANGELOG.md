# Changelog

This changelog covers the work in this branch after `efb2e9344772ea1b46893f39a21d08cebb63ea62`.

## 2026-05-22
- Fixed product query logic by using `MAX` for product names and adjusting grouping.
- Fixed attribute name-row inserts so `attribute_id = 0` is handled correctly.
- Merged the database migration and model-fix work into the branch.

## 2026-05-25
- Added support for `MOD` and `POW` operators in the formula engine.
- Added function-call parsing and evaluation for `MIN`, `MAX`, and `ROUND`.
- Updated lexer, parser, AST, and evaluator logic to support the new formula syntax.

## 2026-05-26 to 2026-05-27
- Fixed mixed integer/float type conversion and grouped-expression parsing.
- Added delete endpoints and frontend delete actions.
- Improved product details behavior and category/product selection handling.

## 2026-05-28 to 2026-05-29
- Added inline documentation across store, model, category, product, and formula code.
- Updated CI workflows to build backend and frontend on feature branches.

## 2026-06-03 to 2026-06-05
- Bumped Go to 1.25 and cleaned up formatting.
- Added endpoint tests plus integration test setup and teardown.
- Added Windows and Bash test runner scripts and refined PostgreSQL test configuration.

## 2026-06-17
- Replaced standard logging with structured logging across main, HTTP, store, migration, and service layers.