# giraffe-quail
An implementation of GraphQL, written in idiomatic Go, to the spec.

## The idea
Every year in November I have friends try NaNoWriMo (National Novel Writing Month). It's cool to watch and I have often thought to try it, although I've usually had various deadlines looming over me. This year, I feel like I can carve off the time -- but I'm not really feeling like writing a novel.

So I'm embracing and extending it. Let's call it GoCoWriMo -- Go Code Writing Month.

## The project
I've been experimenting lately with GraphQL. I really like it -- to me, it's that thing I'd want to build on top of a RESTful API but can never justify because I don't have the resources to spend to build something that big for the kinds of projects I've had. It's a nice design, but it's overkill for any one project. Only an organization on the scale of a Facebook (with a really broad array of applications that could use it) can justify the engineering time to design and build it the first time.

When I first looked at it, I looked at [the graphql-go project](https://github.com/graphql-go/graphql). It's an interesting start, but:

* it's now way out of date with respect to the graphql spec
* the author appears to have lost interest in it (hasn't been updated lately)
* it says that it's a straight port of the original Node/JavaScript code, which means it's got some unusual structure.

Then I started playing with the Node/JS project, and at first it felt great. But as my implementation got more complex, it made me remember why I abandoned Node for Go in the first place -- complex Node code becomes very hard to reason about.

The right way to do this is to write a good parser first, to the spec (which was written as a parser spec). And not all that many people know how. I am not a programming language guru, but I have written several parsers for a variety of uses during my career.

So...I have a month, I have a spec, I really believe that a good GraphQL server platform in Go is needed and would be valuable to the community, and I believe I'm qualified to write it. So here we go.

## The plan so far:

- [ ] Tooling
    - [x] create the project and the repository
    - [x] create a sublime text syntax file for peggo
- [ ] write a parser
    - [x] choose a parser (Pigeon/PEG)
    - [ ] define all the tokens
    - [ ] write the parser (transforming to what - an AST?)

## The worklog

* Nov 03 - (2 hr)
    * Adapted the PEGjs Sublime Text settings file to PEGgo. Now I have a usable editor.
    * Settled on a format for peggo files that I think is pretty acceptable.
    * Started writing the grammar from the spec. Got a good way into it, tried to build it, and I'm running into what I think is an unclosed bracket or something. Haven't found it yet, too tired to continue.

* Nov 02 (2.5 hr) - Looked into PEG. I wrote a a calculator (not by copying anyone else's - wanted to try it myself). It's nice because it doesn't need a lexer (unlike yacc). It's also easier to write the grammar. I'm going to run with it, but I'm not sure yet about a couple of things:
    * How best to handle syntax errors (writing them into the grammar is the best way to give good error handling, but it really bloats the grammar).
    * How to format the grammar. The nature of go means that a lot of interface{} and []interface{} are getting passed around, so it's not nice to embed much code in the grammar itself. Indentation and formatting (as with yacc) can get kind of ugly. I played with nested embedded functions which might work out.
    * What the grammar should generate. It might be best to have it build an AST because then I can get all the interface{} elements out of the way.

* Nov 01 (1 hr) - Set up project, write first-pass README, begin planning
