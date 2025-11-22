# Gosh (Go#)

Gosh (Go#) is an experimental, OOP-inspired programming language built on top of Go. 
It keeps Go’s performance and simplicity while introducing higher-level syntax, 
object-oriented concepts, and an easier developer experience (DX).

---

Key Features:

- OOP-inspired syntax – define variables, expressions, and logic in a more intuitive way.
- Seamless Go integration – write native Go code inline using `::` for ultimate flexibility.
- Automatic transpilation – `.gosh` files are transpiled into Go code and executed automatically.
- Performance-first – leverages Go’s speed and efficiency under the hood.
- Easy import handling – simple `import` syntax that maps directly to Go packages.

---

Concept Overview:

Gosh’s pipeline follows a clear, structured flow:

+------------------+
|     .gosh file    |
+------------------+
          |
          v
+------------------+
| Read bytes from  |
|     file         |
+------------------+
          |
          v
+------------------+
| Convert to string |
|     content       |
+------------------+
          |
          v
+------------------+
|   Create Lexer    |
+------------------+
          |
          v
+----------------------------+
| Generate Tokens character  |
|       by character         |
+----------------------------+
          |
          v
+----------------------------+
| Parse Tokens into AST      |
|        expressions         |
+----------------------------+
          |
          v
+----------------------------+
| Transpile AST to Go code   |
+----------------------------+
          |
          v
+----------------------------+
| Run Go code automatically  |
+----------------------------+
          |
          v
+----------------------------+
|      Execution complete    |
+----------------------------+


---

Syntax:

1. Imports:

import fmt; math;

- All imports must be at the **top of the file**.
- Automatically handled in transpilation.

2. Variable declaration:

[type] [variable_name] = [RHS];

Examples:

int a = 10;
bool b = true;
string message = "Hello Gosh!";

3. Native Go code:

Use `::` to write native Go code directly:

:: fmt.Println(a, b, message)

- Allows you to execute Go code inline without transpilation interference.

4. Expressions:

- Supports numeric operations with proper operator precedence:
  
int c = (10 * (a + b));

- Parentheses are supported to control operation order.

---

Usage:

gosh <filename.gosh> [--debug=true|false]

- `--debug=true` prints the generated Go code without running it.
- `--debug=false` (or omit) executes the transpiled code immediately.

---

Gosh is designed to combine Go’s performance with a more flexible, high-level syntax, 
making it easier to write expressive code without sacrificing speed.
