# RJSON: Repeatable JavaScript Object Notation

"Repeatable JSON", or RJSON, is a strict formatting of a standard JSON document. It is called "repeatable" because it is formatted in strict manner such that if two documents contain the same content, the two resulting JSON document are identically serialized.

## Goals

In order of importance:

1. It is compatible with the JSON spec. In the end, any JSON parser can read a RJSON document as simply JSON. The "repeatable" part is in the strict manner of the formatting.
2. An object with the same content will ALWAYS create exactly the same serialization of the document. Two document files containing the same data will create the same hash code number.
3. It should be human readable.
4. The JSON should be line-oriented for easy tracking of differences in utilities such as GIT.

Easy human *writability* is NOT a goal. It should be fairly easy for a person to read RJSON; but writing is very strict. Enough so that having a human write anything non-trivial in RJSON is going to be frustrating.

The goals are to be met by a set of rules that:

* state the standard for layout and indentation, and
* more precisely state the definitions of some elements than the original spec.

## Target Libraries

1. Javascript
2. Go
3. Python

It would be good to have support in other languages, but this author is concentrating strictly on those three.

Because RJSON is about the expression of standard JSON, the libraries need only influence the serialization. There is no need for handling deserialization.

However, it is my intention to also write a linter in Go at some point.

## The Rules

### `FORMATTING 1` - For objects and array entries, indentation is two spaces and ONLY two spaces.

Good:

```json
{
  "a": 3.2,
  "b": [],
  "c": false
}
```

Bad RJSON (even though valid JSON):

```json
{"a": 3.2,"b": [],"c": false}
```

*Indentation is expected before each name.*

```json
{
    "a": 3.2,
    "b": [],
    "c": false
}
```

*The indentation is 4 spaces, but should be 2 spaces.*

### `FORMATTING 2` - For object key/value pair lines: indentation, field name (key), colon, exactly one space, start of value. 

Good:

```json
{
  "a": 3.2,
  "b": [],
  "c": false
}
```

Bad RJSON (even though valid JSON):

```json
{
  "a":3.2,
  "b":[],
  "c":false
}
```

*A single space is expected after each colon.*

```json
{
  "a" : 3.2,
  "b" : [],
  "c" : false
}
```

*There cannot be spacing before the name and the colon.*

### `FORMATTING 3` - For object and arrays entries, the end of a value ends the line (via LINEFEED or COMMA LINEFEED) 

Good:

```json
{
  "a": 3.2,
  "b": [],
  "c": false,
  "d": [
    1,
    2
  ]
}
```

Bad RJSON (even though valid JSON):

```json
{
  "a": 3.2 ,
  "b": [] ,
  "c": false ,
  "d": [
    1 ,
    2
  ]
}
```

*There should not be an extra space between the value and the comma.*

```json
{
  "a": 3.2
  , "b": []
  , "c": false
  , "d": [
    1
    ,2
  ]
}
```

*The comma comes before the LINEFEED or indentation.*

Not shown: adding white space at end of line also not allowed. But this is difficult to demonstrate using markdown.

### `FORMATTING 4` - All whitespace is unicode 'SPACE' U+0020 and all lines are terminated with a single 'LINE FEED' unicode U+000A (`\n`)

The JSON spec allows for `carriage returns` and `tabs`, however RJSON never uses those in it's expression.

### `FORMATTING 5` - An empty object is represented by `{}` and an empty array by `[]`

Good:

```json
{
  "foo": {},
  "bar": []
}
```

Bad RJSON (even though valid JSON):

```json
{
  "foo": { },
  "bar": [ ]
}
```

*A space between left and right brackets is not allowed.*

```json
{
  "foo": {
  },
  "bar": [
  ]
}
```

*If empty, an object or array cannot span lines.*

### `FORMATTING 6` - Object field names (keys) start on the same line as the start of the value.

Good:

```json
{
  "zed": [
    1
  ],
  "zigzigzigzigzig": "some really long text",
  "zag": {
    "a": 99
  }
}
```

Bad RJSON (even though valid JSON):

```json
{
  "zed": [ 1 ],
  "zigzigzigzigzig": "some really long text",
  "zzz": { "a": 99 }
}
```

*Even if small, non-empty arrays and objects span multiple lines and are indented.*

```json
{
  "zed": [
    1
  ],
  "zigzigzigzigzig": 
    "some really long text",
  "zzz": {
    "a": 99
  }
}
```

*Event if a name/value combination is long, the value must start on the same line as the name. In the case of a string value, the whole name/value pair will be on one line.*

### `FORMATTING 6` - All documents end with LINEFEED

The last character in a document or the general expresson of RJSON should be a single LINEFEED character following the closure of the document's root array or object.

As an example, the RJSON document

```json
{
  "foo": "bar"
}
```

contains 3 lines and a LINEFEED character at the end of each line. So, a total of 3 LINEFEED characters.

### `FORMATTING 7` - A JSON document's root element must either be an object or array

The JSON spec sort-of kind-of implies this with the words "JSON is built on two structures...".

But RJSON is more explicit: the top-level structure is either an object or an array.

Because of RJSON strict formatting, a reader need only look at the first character in a document to detect which type of root element is contained in the document. If the first character is a curly brace `{` then the document contains an object. If the first character is a left square bracket `[` then it is an array.

Good:

```json
{
  "sample": 42
}
```

```json
{}
```

```json
[
  42
]
```

```json
[]
```

```json
[
  {
    "id": 293,
    "name": "box"
  },
  {
    "id": 429,
    "name": "cylinder"
  }
]
```

Bad RJSON (and possibly bad JSON, but we are not sure):

```json
42
```

*This is a bare Number. RJSON does not allow this but JSON might or might not.*

```json
"hello world"
```

*This is a bare string. RJSON does not allow this but JSON might or might not.*

```json
true
```

*This is a bare boolean `true`. RJSON does not allow this but JSON might or might not.*

### `NUMBERS 1` - Numbers use a capital E if using an exponent

Good:

```json
{
  "number": 2.345E20
}
```

Bad RJSON (even though valid JSON):

```json
{
  "number": 2.345e20
}
```

*Do not use a lower-case 'e' for the exponent marker.*

### `NUMBERS 2` - Numbers do not use the plus symbol if using a positive exponent

Good:

```json
{
  "number": 2.345E20,
  "small": 4.56E-10
}
```

Bad RJSON (even though valid JSON):

```json
{
  "number": 2.345E+20,
  "small": 4.56E-10
}
```

*The +20 for an exponent is allowed with JSON but not RJSON*

Bad RJSON **and** bad JSON:

```json
  "number": +2.345E20
```

*The JSON spec already forbids using a plus symbol for positive numbers.*

### `NUMBERS 3` - If reading/writing RJSON documents, Number precision must be preserved.

Decimal floating-point numbers have implied precision and it should not be discarded. For example:

```json
{
  "a": 1,
  "b": 1.000,
  "c": 3.1415E20,
  "d": 3.14150000E20,
}
```

The values for "a", "b", "c", and "d" are different. "a" is exactly 1 but "b" is 1 to three decimal points. They may both have the same scalar value, but the precision (or *significance*) differs.

### `NUMBERS 4` - Numbers are in decimal not binary

The JSON spec explicitly defines a "number" as involving the Base10 digits 0 to 9. It does not support Binary (Base2) or Hex (Base16 aka Base2^4). For integers, this is not critical. But for floating numbers, a binary floating point number cannot store the full set of numbers expressed by decimal floating point numbers. For example, one-fifth (1/5) is "0.2" in decimal. But in binary, one-fifth is a repeating number and can only be approximated. This is very similar to how one-third is a repeating number in decimal (0.3333...) and can only be stored as an approximation in decimal.

(Trivia: one-third CAN be precisely be stored in Sumerian sexagesimal as ` .𒌋𒌋` (note the space before the period to denote nothing/zero.) That is the numbering system used by the extinct civilization of Ur. The numbering is Base60. However, one-seventh is a repeating number in sexagesimal. No numbering system can prevent all divisor flaws; but that is a mathematical discussion well outside the scope of this document.)

This gives way to "subtle" rounding errors when converting between the decimal to binary.

For many applications, this is a subtle distinction that does not mean much and is handled by good rounding algorithms. But if, for example, a RJSON document is used in a financial application, the numbers are technically decimal not binary should ideally be treated as such.

Not all programming languages support decimal numbers. Some helpful references:

| language    | ref |
| ----------- | --------------------- |
| Python      | there is a standard library called 'decimal' |
| JavaScript  | https://www.npmjs.com/package/js-big-decimal |
| Go          | https://pkg.go.dev/github.com/shopspring/decimal |

### `STRINGS 1` - In a string, the `\U` escapement is not to be used.

The actual UTF-8 codepoints are to be inserted into the string.

Good:

```json
{
  "greeting": "hello 鲍勃"
}
```

Bad RJSON (even though valid JSON):

```json
{
  "greeting": "hello \u9c8d\u52c3"
}
```

*You must choose to use the actual codepoints, not the escaped equivalents.*

```json
{
  "greetin\u0067": "hello 鲍勃"
}
```

*This rule applies to the name string as well.*

### `DATA TYPES 1` - Arrays are actually "lists" not "arrays"

In Computer Science standard terminology, an "array" is a fixed list of items of the same type. However, in JSON, an "array" is neither fixed nor is it required to be of the same type. 

It is called an "array" for historic reasons: the JavaScript language also calls them arrays despite not being actual arrays in JavaScript either. The reason for this in JavaScript is not known to us. Perhaps the lists are stored using an array-of-pointers in the underlying interpreter; but that would be an hidden implementation detail not a language spec detail.

Here are examples:

Two valid JSON documents with differening sizes for the field "aaa":

```json
{
  "aaa": [
    1,
    2,
    3
  ]
}
```

```json
{
  "aaa": [
    1,
    2,
    3,
    4
  ]
}
```

A valid JSON document with mixed types in the field "aaa":

```json
{
  "aaa": [
    1,
    null,
    "hello",
    false
  ]
}
```

### `DATA TYPES 2` - `null` represents "not known" or "unknown"

In RJSON, `null` very specifically represents the database meaning commonly found in SQL specifications. It means that a value is not known. Or, to be explicit: it does NOT mean empty or does-not-exist.

As a demonstration:

```json
{
  "a": 0,
  "c": null
}
```

This above document shows that:

- "a" is an empty integer. 
- "b" does not exist. (aka "missing" or "undefined")
- "c" is not known.

These concepts are expressed differently in different languages:

| concept         | JSON              | Python                | JavaScript    | Go                 |
| --------------- | ----------------- | --------------------- | ------------- | ------------------ |
| empty number    | `0`               | `Decimal("0")` or `0` | `0`           | `decimal.NewFromInt(3)` or `0` |
| empty list      | `[]`              | `[]`                  | `[]`          | `make([]string, 0)` |
| empty object    | `{}`              | `{}`                  | `{}`          | `map[string]interface{}{}` |
| empty string    | `""`              | `""`                  | `''` or `""`  | `""` |
| non-existance   | (field missing)   |                       | `undefined`   |  |
| null / unknown  | `null`            | `None`                | `null`        |  |

Python does not have a positive means of noting non-existance. However `key in obj` can detect existance. And `del obj[key]` can remove existance.

Go does not have a positive means of notating unknown vs non-existance. Often `nil` will be used to indicate either of those concepts, but `nil` really means "undefined pointer". Extra effort must be taken during marshalling/unmarshalling. Also see 'omitempty' handling.

Sadly, some languages, such as C#, have standard libraries that treat unknown and non-existance as the same thing leading to possible security vulnerabilities. (ref) Writing a RJSON library in that language will be ... challenging.

C supports `NIL` but that is distinctly a pointer with no reference. Like *many* things in a mid-level language like C, supporting null and/or non-existence is left as a exersize to the programmer.

**TLDR:** So, if a RJSON library or program reads a RJSON document and writes that RJSON document back out:

* If a field *does not exist* when it was read, and no change is made to that field, it is not written out.
* If a field is `null` when it was read, and no change is made to that field, it should remain `null` when it is written.

### `DATA TYPES 3` - An Object's field names may not repeat in the same level of that object

Good:

```json
{
  "id": "Joe"
}
```

```json
{
  "id": "Joe",
  "pet": {
    "id": "Mittens",
    "type": "rabbit"
  },
  "addr": {
    "detail": "1234 Main St.",
    "id": "home"
  }
}
```

*Though the name `id` shows up in three places, they are not at the same level in the same object, so that is okay.*

Bad RJSON (even though valid JSON):

```json
{
  "id": "Joe",
  "id": "Joe"
}
```

*Duplicate keys ("id") not allowed in RJSON*

```json
{
  "id": "Joe",
  "id": "Larry"
}
```

*Duplicate keys ("id") not allowed in RJSON*

```json
{
  "id": "Joe",
  "pet": {
    "id": "Mittens",
    "type": "rabbit",
    "type": "domestic"
  }
}
```

*Duplicate keys ("type") in subtending "pet" object not allowed in RJSON*

It *might* seem like this rule violates goal #1: "Compatible with the JSON spec." However, it is this author's contention that it does not. If curious, I've discussed this subject in FAR TOO MUCH DETAIL in another document: [JSON_DUPLICATE_NAMES](JSON_DUPLICATE_NAMES.md)

### `DATA TYPES 3` - Use a repeatable order for the fields of an object

The items in a JSON object are not placed in a particular order per the JSON spec.

But, to meet the goals of RJSON the expression of those fields needs to be non arbitrary.

For example, if application A has a user named "Joe" who is age 32 and writes that information to a file "a_joe.json":

```json
{
  "name": "Joe",
  "age": 32
}
```

And another application has the same user who is age 32 with the name "Joe" and write that information to a file "b_joe.json":

```json
{
  "age": 32,
  "name": "Joe"
}
```

The content of these files, for exact same data, are different. A hash of `a_joe.json` will generate a different checksum than the hash of `b_joe.json`, implying different data. This is precisely what we do not want to happen.

Another example: a utility reads file "abc.json" in a git repo, makes no changes whatsoever, and writes a new "abc.json". Will the repo detect a change? With just JSON it might because the order of any object fields is not predictable. With RJSON, the fields are always ordered the same way, so the file should not have any changes.

**SORT is by rune codepoint**

Object items are sorted in the order of their 32-bit equivalent UTF-8 code points (runes) list of their names (keys).

That is, if each name is treated as a list of utf-8 characters (aka runes), and each utf-8 character is treated as it's normalized 32-bit integer, then the sorting is based on the difference of 32-bit equivalants in those lists.

If two names contain the same shared code points, then the longer name is last and the shorter name is first.

Good:

```json
{
  "a": 4948,
  "aa": 4949,
  "b": 223,
  "c": false,
  "d": "jerry"
}
```

Bad RJSON (even though valid JSON):

```json
{
  "aa": 4949,
  "b": 223,
  "a": 4948,
  "c": false,
  "d": "jerry"
}
```

*The "a" name should come before all the other fields because it is "lowest" in the sorting order*

Please note that this is NOT the same thing as collation (aka "alphabetical order"). Collation cannot be used because it is contextual to the reader's locality/language. See the COLLATION spec for Unicode for many examples. Plus the COLLATION spec can (and does) change over time; which is not a desirable trait for a repeatable spec like RJSON.

Good:

```json
{
  "Aa": "dkd",
  "Bb": "vud",
  "aa": "joe",
  "bb": "ivw"
}
```

Bad RJSON (even though valid JSON):

```json
{
  "Aa": "dkd",
  "aa": "joe",
  "Bb": "vud",
  "bb": "ivw"
}
```

*This doc object is bad because it is in alphabetical order per English rules; and is not sorted by unicode code point values.*

## Downsides to RJSON

No specification is perfect. RJSON has definite flaws:

### Object fields are likely not in logical order

Good (but unfortunately not in a logical order):

```json
{
  "addr1": "123 Main Str",
  "age": 32,
  "city": "Springfield",
  "family_name": "Smith",
  "given_name": "Joe",
  "zip": "65102"
}
```

Bad RJSON containing the same data:

```json
{
  "given_name": "Joe",
  "family_name": "Smith",
  "age": 32,
  "addr1": "123 Main Str",
  "city": "Springfield",
  "zip": "65102"
}
```

RJSON is about the expression of JSON. Internally, an application could read a RJSON document and order the fields in any way that it wants, including logically or by category. It is only when writing/serializing/expressing the JSON document as RJSON that code-point equivalent field ordering is applied.

### RJSON is not as small as compacted JSON

Or, to show it visibly. The RJSON doc:

```json
{
  "a": 4948,
  "b": 223,
  "c": false,
  "d": "jerry"
}
```

is 58 bytes long. This includes the final LINEFEED. It is larger than minified JSON with the same content:

```json
{"a":4948,"b":223,"c":false,"d":"jerry"}
```

which is only 40 bytes long.

For larger objects with greater indentation, the descrepancy can get much worse. These larger sizes will likely have an impact on performance.

None-the-less RJSON is never minimized as performance is not one of the goals.

### In `git` diff, appending to an array or object shows up as multiple changed lines, not one.

With the original doc:

```json
{
  "aaa": [
    4,
    5,
    99,
    3
  ]
}
```

Adding `8` to "aaa" will show
* a changes on line 6 (adding a comma), `    3,`
* a new line on line 7. `    8`

This is due to JSON spec not allowing a trailing commas on the last value of an array.

In fact, appending to an empty list shows three changes. For example:

```json
{
  "bbb": []
}
```

Adding `8` to "bbb" becomes:

```json
{
  "bbb": [
    8
  ]
}
```

So, one line changed and two lines inserted: three changes. When the RJSON spec was being put together, we considered showing empty lists like this:

```json
{
  "ccc": [
  ]
}
```

However, our experience with larger documents showed that indented empty lists made things _much_ harder to read for humans. The empty sets of indented brackets created a lot of visual whitepace "noise". And, in the end, human readability was a greater goal than git tracking.
