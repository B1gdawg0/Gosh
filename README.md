# GO# (Gosharp) or Gosh

## Look, I love Go. But sometimes... ugh.

You know that feeling? You're cruising through Go, loving the speed, loving the simplicity. But man, I love OOP. I love Java, C#, and them. Go's got no actual classes, no try-catch, weird syntax sometimes.

And you're just sitting there like... "why tho?"

So I built Gosh for the love of the game. Life's too short to choose between performance and developer happiness.

---

## What's the deal?

**Go is FAST.** Like, real real fast. But it's missing some love:
- No OOP (classes? inheritance? nah)
- No try-catch (good luck with that `if err != nil` ... millions of times)
- No enums (const spam go brrrr)
- No event-driven patterns
- No dependency injection
- No... you get it

**Gosh is the "yes, and" to Go's "no, but"**

Write OOP-style code that transpiles to pure Go. You get modern syntax with Go's legendary performance baked in. And, This might be the **fastest OOP language ever** right? Because under the hood, it's just Go. We're not interpreting. We're not adding runtime overhead. We're transpiling to native Go code and letting the Go compiler do its magic. Java? Slower. C#? Slower. TypeScript? Don't even start. They're all battling JIT compilers and runtime VMs. Gosh just... becomes Go. That's it. That's the trick.

---

## How does this transpiler work?

```
Your .gosh file
    ↓
Lexer breaks it into tokens (char by char)
    ↓
Parser builds an AST (abstract syntax tree)
    ↓
Transpiler converts AST → pure Go code
    ↓
Go compiler goes brrrr
    ↓
Done
```

No VM. No interpreter. No runtime BS. Just straight Go performance with better syntax.

---

## Show me the goods

### Install & Run
```bash
gosh script.gosh              # Run it
gosh script.gosh --debug=true # See the Go code it generates
```

### Syntax (it's dead simple)

**Imports** - throw 'em at the top:
```gosh
import fmt math strings;
```

**Variables** - type come first:
```gosh
int speed = 9999;
bool isAwesome = true;
string vibe = "immaculate";
```

**Math** - normal math:
```gosh
int result = (10 + 20) * 3;
int chaos = ((a + b) * c) / (d - e);
```

**Native Go escape hatch** - need raw Go? Use `::` (syntax was actually added for debugging, but let's keep it):
```gosh
int x = 42;

:: fmt.Println("The answer is", x)
:: if x > 40 {
::     fmt.Println("We're in the big leagues now")
:: }
```

---

## Full Example

```gosh
import fmt math;

int radius = 10;
float area = math.Pi * float(radius * radius);

:: fmt.Printf("Circle area: %.2f\n", area)
:: fmt.Println("This ran at Go speed, btw")
```

**What you write** ↑

**What runs** ↓

```go
package main
import ("fmt" "math")

func main() {
    radius := 10
    area := math.Pi * float64(radius * radius)
    fmt.Printf("Circle area: %.2f\n", area)
    fmt.Println("This ran at Go speed, btw")
}
```

---

## For you, yes YOU reading this

If you're an early adopter of this project: good luck with debugging and breaking prod. It's gonna be a pain in your ass. But man, in the end, if all this works, you'll be part of the **fastest OOP language ever created**.

If you wanna add something, contribute something, or just vibe - let's talk.

**My contact:** 
Email: lerdphipat.k@gmail.com

---

**Built by a dev who refused to choose between Go or my GOAT C#😘**

*Let's cook* 👨‍🍳