# giraffe-quail
An implementation of GraphQL, written in idiomatic Go, to the spec.

## The idea
Every year in November I have friends try NaNoWriMo (National Novel Writing Month). It's cool to watch and I have often thought to try it, although I've usually had various deadlines looming over me. This year, I feel like I can carve off the time -- but I'm not really feeling like writing a novel.

So I'm embracing and extending it. Let's call it GoCoWriMo -- Go Code Writing Month.

## The project
I've been experimenting lately with GraphQL. I really like it -- to me, it's that thing I'd want to build on top of a RESTful API but can never justify because I don't have the resources to spend to design build something that big for the kinds of projects I've had. It's a nice design, but it's overkill for any one project. Only an organization on the scale of a Facebook could justify the engineering time to design and build it the first time.

When I first looked at it, I looked at [the graphql-go project](https://github.com/graphql-go/graphql). It's an interesting start, but it's now way out of date with respect to the graphql spec, the author appears to have lost interest in it, and it also says that it's a straight port of the original Node/JavaScript code, which means it's got some unusual structure.

Then I started playing with the Node/JS project, and at first it felt great. But as my implementation got more complex, it made me remember why I abandoned Node for Go in the first place -- complex Node code becomes very hard to reason about.

The right way to do this is to write a good parser first, to the spec (which was written as a parser spec). And not all that many people know how. I am not a programming language guru, but I have written several parsers for a variety of uses during my career.

So...I have a month, I have a spec, I really believe that a good GraphQL server platform in Go is needed and would be valuable to the community, and I believe I'm qualified to write it. So here we go.

## The plan so far:

- [x] create the project and the repository
- [ ] write a parser
    - [ ] choose a parser (yacc? Pigeon/PEG?)
    - [ ] define all the tokens
    - [ ] write a good lexer
    - [ ] write the parser (transforming to what - an AST?)

