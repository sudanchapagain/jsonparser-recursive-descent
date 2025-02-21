# JSON parser

```ebnf
JSON    → Object | Array
Object  → "{" PairList? "}"
PairList → Pair ("," Pair)*
Pair    → STRING ":" Value
Array   → "[" ValueList? "]"
ValueList → Value ("," Value)*
Value   → STRING | NUMBER | BOOLEAN | NULL | Object | Array
```
