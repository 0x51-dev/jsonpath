# JSONPath

JSONPath defines a string syntax for selecting and extracting JSON (RFC 8259) values from within a given JSON value.

## Overview

A brief overview of JSONPath syntax.

| Syntax Element    | Description                                                                                                                            |
|-------------------|----------------------------------------------------------------------------------------------------------------------------------------|
| `$`               | [root node identifier](https://www.rfc-editor.org/rfc/rfc9535.html#root-identifier)                                                    |
| `@`               | [current node identifier](https://www.rfc-editor.org/rfc/rfc9535.html#filter-selector) (valid only within filter selectors)            |
| `[<selectors>]`   | [child segment](https://www.rfc-editor.org/rfc/rfc9535.html#child-segment): selects zero or more children of a node                    |
| `.name`           | shorthand for `['name']`                                                                                                               |
| `.*`              | shorthand for `[*]`                                                                                                                    |
| `..[<selectors>]` | [descendant segment](https://www.rfc-editor.org/rfc/rfc9535.html#descendant-segment): selects zero or more descendants of a node       |
| `..name`          | shorthand for `..['name']`                                                                                                             |
| `..*`             | shorthand for `..[*]`                                                                                                                  |
| `'name'`          | [name selector](https://www.rfc-editor.org/rfc/rfc9535.html#name-selector): selects a named child of an object                         |
| `*`               | [wildcard selector](https://www.rfc-editor.org/rfc/rfc9535.html#wildcard-selector): selects all children of a node                     |
| `3`               | [index selector](https://www.rfc-editor.org/rfc/rfc9535.html#index-selector): selects an indexed child of an array (from 0)            |
| `0:100:5`         | [array slice selector](https://www.rfc-editor.org/rfc/rfc9535.html#slice): start:end:step for arrays                                   |
| `?<logical-expr>` | [filter selector](https://www.rfc-editor.org/rfc/rfc9535.html#filter-selector): selects particular children using a logical expression |
| `length(@.foo)`   | [function extension](https://www.rfc-editor.org/rfc/rfc9535.html#fnex): invokes a function in a filter expression                      |

## References

- Gössner, S., Ed., Normington, G., Ed., and C. Bormann, Ed., "JSONPath: Query Expressions for JSON", RFC 9535, DOI
  10.17487/RFC9535, February 2024, <https://www.rfc-editor.org/info/rfc9535>.
