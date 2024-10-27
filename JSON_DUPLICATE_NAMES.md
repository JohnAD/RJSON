# Does JSON allow duplicate names in an Object?

This is a mildly controversial subject that is discussed from time to time.

## TLDR

The JSON specifications do not forbid duplicate names (keys) in an object. This implies permission.

The JavaScript language that JSON is based on forbids duplicate keys in objects. This implies that it is forbidden.

And there is my crazy take: JSON allows duplicate names but all duplicates but the last one are to be quietly ignored.

This document is my justification of my crazy take.

## Data Type Background

JSON is not a "Data Type". Nor is it a "collection" in the Computer Science meaning of the word. It is meant to be a serialization of a data type that can contain values that have names.

There are two data types in common usage for values-that-have-names.

### Ordered Associative Array: a sequence of key/value pairs.

When an `associative array` is ordered, the data types allows and encourages duplicate names. Probably the most common use of an ordered associative array on the Internet is the result of an HTML Form. For example,

```html
<html>
 <body>
  <form method="POST">
   <input type="text" name="actor">
   <input type="text" name="actor">
   <input type="text" name="actor">
   <input type="text" name="year">
  </form>
 </body>
<html>
```

Generates something like the following HTTP body after the header:

```text
actor=Larry
actor=Moe
actor=Curly
year=1944
```

Because of the ordering, you can ask an ordered associative array to "get the second 'actor'" and receive a response of "Moe".

Most major languages support ordered associative arrays, though sometimes only through third party libraries or as a sequence of tuples.

Note: there is quite a few of Python web pages that describes `dict` as an associative array (and it is a *subtype*). They then also declare that "**all** associative arrays" require unique keys. **I am shocked by this.** For a correct reference, see *"Programming Language Pragmatics (Third Edition)"* by Michael Scott; see chapter "7 - Data Types". Or pretty much any computer science book on data types.

### Mapping: an unordered set of values mapped by unique keys

A `mapping` is also sometimes called:

 - `dictionary`
 - `map`
 - `record`
 - `struct`
 - `object`

And, if all of the values are restricted to the same type, you could include `hash table`.

Nearly all high-level languages support mappings natively.

## Looking at the JSON spec

The [JSON spec](https://www.json.org/json-en.html) does not meet either of the definitions mentioned above.

It is not an `ordered associative array` because the key/value pairs are not explicitly ordered.

It is not a `mapping` because they key/value pairs are not explicitly required to be unique.

In fact, the actual wording in the spec...

    `A collection of name/value pairs. In various languages, this is realized as an object, record, struct, dictionary, hash table, keyed list, or associative array.`

mentions various data types that are NOT neccesarily compatible with each other, such as `dictionary` and `keyed list`.

I would argue that the JSON spec, by itself, cannot answer our question.

In fact, the corresponding RFC for JSON, [RFC-8259](https://www.rfc-editor.org/rfc/rfc8259#section-4) , says:

     The names within an object SHOULD be unique.

In the grammar of RFCs, SHOULD means "not required strongly encouraged". Again, this implies technical permission.

## Looking at common JSON usage

In every language I've seen so far, JSON has been treated as a mapping when it is deserialized into an internal collection.

Details:

| language   | data type used | language name for type | response to duplicates names |
| ---------- | -------------- | ---------------------- | ---------------------------- |
| JavaScript | Mapping | `object` | the object contains the last duplicate name's value |
| Go | Mapping | `map[string]interface{}` using `json.Unmarshal` | the map contains the last duplicate name's value |
| Python | Mapping | `dict` using `json.load` or `json.loads` | the dict contains the last duplicate name's value |

(If your language is not in this table, feel free to send me a repo PR. Languages are in English alphabetical order; except that
Javascript is at the top for historical context.)

But, so far, languages in general do the "ignore the early duplicates" behavior.

## Summary

The JSON spec does not require "ordering" OR "uniqueness" to the key/value pairs. Common practice is to store JSON to a Mapping and silently ignore all but the last duplicate name.

So, my take:

**JSON libraries should behave like a JavaScript and purposely ignore all but the last duplicate key if duplicate keys are found.**

That is all I'm saying. I've used a lot of words to say a simple thing. Sorry about that.

So,

```json
{ "actor": "Moe", "actor": "Larry" }
```

is legitimate JSON and should be deserialized as:

```json
{ "actor": "Larry" }
```

## Trivia

### Side Effect of This Summary

Because all but the last duplicate key should be ignored, it is critical that any JSON deserialization library keep track of order during deserialization. But, the "something" that the data is stored in, however, does not need to keep track of order.

I imagine most libraries handle this correctly already, but I bring this up for the sake of completeness.

### Can you create a collection / data type that DOES match JSON spec

Could someone create a collection type that matches the open-ended JSON spec? In theory, yes. Perhaps it would be called an `associative set` and would be defined as "an unordered set of name/value pairs". With it you could query some very limited data. Theoretical example in Python:

```python
from associativeset import AS

x = AS()
x["actor"] = "Larry"
x["actor"] = "Moe"
x["actor"] = "Curly"

assert x["actor"].set_count == 3
assert "actor" in x
assert "other" not in x
assert "Moe" in x["actor"]
assert "Curly" in x["actor"]
assert "Shemp" not in x["actor"]
assert {"Larry", "Curly", "Moe"} == x["actor"]
assert {"Curly", "Moe", "Larry"} == x["actor"]  # order DOES NOT MATTER in a set
```

I'm not convinced writing such a library would be useful. But, in theory, it could be done.

### How does Flask handle Form POST input?

Obviously, you can't use a `dict` or `OrderedDict` to store those values. So, how is it handled in Flask (a python web server)?

It is stored as a specialized class of the WTF libraries `Form` called `FlaskForm`. This class builds what they call a "multi dict" which gets an ordered list of values for any one key applied against a schema. So, using the example of:

```text
actor=Larry
actor=Moe
actor=Curly
year=1944
```

You would see in the `route` handler:

```python
assert request.form["actor"] == ["Larry", "Moe", "Curly"]
assert request.form["year"] == ["1944"]
```

Short answer: everything is always a list. Since the HTML form content is flat and does not support lists intrinsically, this does not conflict with anything.
