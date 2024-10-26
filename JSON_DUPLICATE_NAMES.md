# Does JSON allow duplicate names in an Object?

This is a controversion subject that is discussed from time to time.

## TLDR

The JSON Specification does not forbid duplicate names in an object. This implies permission.

The JavaScript language that JSON is based on does forbid duplicate keys in objects. This implies that it is forbidden.

And there is my crazy take: JSON allows duplicate names but all duplicates but the last are to be ignored.

This document is my justification of my crazy take.

## Data Type Background

JSON is not a "Data Type". Nor is it a "collection" in the Computer Science meaning of the word. It is meant to be a serialization of a data type that can contain values that have names.

There are two data types in common usage for values-that-have-names.

### Ordered Associative Array: a sequence of key/value pairs.

When an `associative array` is ordered, this data types allows and encourages duplicate names. Probably the most common used of an associative array on the Internet is the result of an HTML Form. For example,

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

Because of the ordering, you can ask an ordered associative array to "get the second actor" and get a response of "Moe".

Most major languages support associative arrays, though sometimes only through third party libraries or as a sequence of tuples.

Note: there is quite a bit of Python documentation that describes `dict`s as associative arrays (and they are a form of them) and then declare that "all associative arrays" require unique keys. I am shocked by this. For a correct reference, see *"Programming Language Pragmatics (Third Edition)"* chapter **7 - Data Types**. Or pretty much any computer science book.

### Mappings: an unordered list of values mapped to unique keys

These are also sometimes called:
 - `dictionary`
 - `map`
 - `record`
 - `struct`
 - `object`

And, if all of the values are restricted to the same type, you could include `hash table` and `keyed list`.

Nearly all high-level languages support mappings natively.

## Looking at the JSON spec

The JSON spec does not meet either of the definitions mentioned above.

It is not an `ordered associative array` because the key/value pairs are not explicitly ordered.

It is not a `mapping` because they key/value pairs are not explicitly required to be unique.

In fact, the actual wording in the spec...

    `A collection of name/value pairs. In various languages, this is realized as an object, record, struct, dictionary, hash table, keyed list, or associative array.`

mentions collection data types that are NOT neccesarily compatible with each other, such as `dictionary` and `hash table`.

I would argue that the JSON spec, by itself, cannot answer our question.

## Looking at common JSON usage

In every language I've seen so far, JSON has been treated as a mapping when it is deserialized into an internal collection.

Details:

| language   | data type used | language name for type | response to duplicates names |
| ---------- | -------------- | ---------------------- | ---------------------------- |
| JavaScript | Mapping | object | ... |
| Python | Mapping | dict | ... |

(If your language is not in this table, feel free to send me a repo PR. Languages are in English alphabetical order; except that
Javascript is at the top for historical context.)

## Summary

The JSON spec does not require "ordering" OR "uniqueness" to the key/value pairs. However any collection type that does not enforce one of those traits is not very usable

So, my take:

**The SPEC implies that JSON should behave like a JavaScript object and ignore all but the last duplicate key if duplicate keys are found.**

So,

```json
{ "actor": "Moe", "actor": "Larry" }
```

is legitimate JSON. But it should be read/deserialized as:

```json
{ "actor": "Larry" }
```

## Trivia

### Side Effect of This Summary

Because all but the last duplicate key should be ignored, it is critical that any JSON deserialization library keep track of order during deserialization. The "something" that the data is stored in, however, does not need to keep track of order.

I imagine most libraries handle this correctly already, but I bring this up for the sake of completeness.

### Can you create a collection / data type that DOES match JSON spec

Basically, can you create a collection type that matches the open-ended JSON spec? In theory, yes. Perhaps it would be called an `associative set` and would be defined as "an unordered set of name/value pairs". With it you could query some very limited data. Theoretical example in Python:

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

### How does Flask handles Form POST input?

Objviously, you can't use a `Dict` or `OrderedDict` to store those values. So, how is it handled in Flask (a python web server)?

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

Short answer: everything is always a list. Since the HTML form content is flat and does not support lists intrinically, this does not conflict with anything.
