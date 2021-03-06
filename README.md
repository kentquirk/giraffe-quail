# giraffe-quail
An implementation of GraphQL, written in idiomatic Go, to the spec.

## The idea
Every year in November I have friends try NaNoWriMo (National Novel Writing Month). It's cool to watch and I have often thought to try it, although I've usually had various deadlines looming over me. This year, I feel like I can carve off the time -- but I'm not really feeling like writing a novel.

So I'm embracing and extending it. Let's call it GoCoWriMo -- Go Code Writing Month.

## The reality
I got sick in the middle of the month, had Thanksgiving and also a couple of work-related crises -- so getting this done by mid-Nov didn't happen. But it's still underway and looking a lot more like a reasonable reality.

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

- [x] Tooling
    - [x] create the project and the repository
    - [x] create a sublime text syntax file for peggo
- [x] write a parser
    - [x] choose a parser (Pigeon/PEG)
    - [x] define all the tokens
    - [x] smoketest the parser (build a recognizer testbed)
    - [x] write or acquire a bunch of parser tests - probably can steal test cases from the graphql implementation in Node. (I have nowhere near enough tests, but at the moment there are enough to move on to other things.)
    - [x] create a parser test suite
    - [x] fix the problems the grammar tests find
    - [x] start adding error productions and build an error reporting mechanism
- [x] Schema
    - [x] Create a schema definition language parser
    - [x] Test it on some standard schema
    - [x] Make it generate type definitions for the graphql parser
- [x] Data structures
    - [x] design a parser data structure
    - [x] get the parser to produce it
- [ ] Generate a working query
    - [x] Organize for use by clients
    - [ ] Write some test queries and schema supporting them
    - [ ] Begin defining a structure for writing query resolvers
    - [ ] Build a server that can answer a simple query

## The worklog

* Week of Dec 20
    * Add schema from a real app we're building. Nothing complex, but it's in active use for a Node-based GraphQL implementation.

* Week of Dec 13
    * I am now working on this "officially" as part of my job at Achievement Network. Not going to be maintaining this log as well, but should be making more regular progress. One of these days I'll update this readme to be more appropriate to the current status.
    * Meanwhile, I now have a working parser that propagates the discovered queries up to the return value of the parser. So I can parse a schema and have the data structures to support it, then pass that schema on to the parser which can then read a query against that schema. Now I need to use that parsed object to actually do the query.
    * This will require that I have a structure for building query resolvers. That's the next step.
    * Have a simple test app that can resolve a query and return its results. The API is starting to look like something I'd be willing to deploy.
    * Add a document datatype to reflect proper graphql structure. Now getting into the nitty-gritty of the semantics of a query.
    * Implement a sample query to our status endpoint that actually kinda works.
    * Allow the data for directives to flow all the way through to the Operation and Field objects.

* Dec 4 - (5 hrs) - Sunday
    * I realized that ValueRegistry is really a Scope, so renamed it.
    * Now parser is using the type registry. It's not tested yet, and value assignment only code to work if types are identical. Lots more work to do but I want to start testing it.
    * Imported the Star Wars schema from the FB examples; it loaded right away.
    * Created a small test framework and loaded my test file from it. It too worked, once I imported the starwars schema.

* Dec 3 - (4 hrs) - Saturday
    * Now have a type registry and a value registry that I'm pretty happy with
    * Implemented Enum where the enum name is a type and the enum values are value objects containing strings but marked with the enum type. It is passing tests.
    * Implemented FieldSets and Interfaces and Obj
    * Temporary types are resolved after being defined later
    * Implemented Union
    * I tried to pull everything into the root, but we can't have two instances of a PEGGO-generated parser in the same project. So I refactored to pull types and values and the registries out to a separate project. Now both parsers can use the types project.

* Dec 2 - (5 hrs) - Abbreviated work day plus some evening time.
    * Starting building some tests for the schema parser; found a few minor bugs in the .peggo file but now have a full set of tests working, and it parses the example grammar used in an explanatory page. Now I have to have it start building type definitions.
    * Decided my approach to type registry and values was not quite right; rewriting it.

* Dec 1 - (5 hrs) - Well, so much for GoCoWriMo -- I was sick and then life got in the way. Now I'm finally getting back to it. The good news is that I've been able to fit this into my job, so I do have some time to work on it now.
    * Today I got a collection of data types and variables and values working, so I'm starting to build some useful data structures here.
    * Also wrote the first pass of a grammar for parsing the GraphQL type schema language, which is oddly unspecified. I used the examples and followed the same principles as GraphQL itself (example: commas are considered whitespace) and got something that compiles.

* Nov 14 - () - I'm sick

* Nov 13 - (1.5 hr)
    * Sleeping on it was good. I now have a solid plan and a good part of the implementation for values and types. The key was separating them and having a type registry that knows how to create values.

* Nov 12 - (3 hrs)
    * Figuring out how best to handle variables and types. Made good progress on it but realized that types shouldn't be an enum, they should be a service where types get registered, with some types built in.

* Nov 11 - (4 hrs) - no work for Veteran's Day so I get to do this.
    * Move main out of the parser directory, rename the package.
    * Learn more about how to structure things so that the .peggo file is as minimal as possible.
    * Add some basic productions for simple types.
    * Add parserHelpers.go to help keep code out of the .peggo file.
    * Decide that an AST is probably overkill for this problem -- it's not a full language, and there is no need to generate source code or to transform the result
    into some other form (optimization)
    * Start working on the structure of Value objects; implement the simple types. Looks like variables are going to require a little bit of rework.

* Nov 08-09-10 - recover

* Nov 07 - (2 hrs)
    * Search for test files; grab some test queries from the tests for graphql-js. Find some more missing ignorable locations and fix them.
    * Hunt down a problem with ListValues.
    * Learn how to resolve shift-reduce conflicts in Peg.

* Nov 06 - (3.5 hrs) - had to take a couple of days off
    * My problem was a 0 that should have been a '0'. Obtuse error message meant that to find it I had to binary search the file. That sucked. I resolve to do better with my parser.
    * With that solved, I managed to get the entire grammar into PEGgo, and it parses and builds a program that compiles.
    * Built a simple test driver and a simple test file and it failed. Realized that I hadn't dealt with a whole basket of ignorables and added the _ marker to the productions.

* Nov 03 - (2 hr)
    * Adapted the PEGjs Sublime Text settings file to PEGgo. Now I have a usable editor.
    * Settled on a format for peggo files that I think is pretty acceptable.
    * Started writing the grammar from the spec. Got a good way into it, tried to build it, and I'm running into what I think is an unclosed bracket or something. Haven't found it yet, too tired to continue.

* Nov 02 (2.5 hr) - Looked into PEG. I wrote a a calculator (not by copying anyone else's - wanted to try it myself). It's nice because it doesn't need a lexer (unlike yacc). It's also easier to write the grammar. I'm going to run with it, but I'm not sure yet about a couple of things:
    * How best to handle syntax errors (writing them into the grammar is the best way to give good error handling, but it really bloats the grammar).
    * How to format the grammar. The nature of go means that a lot of interface{} and []interface{} are getting passed around, so it's not nice to embed much code in the grammar itself. Indentation and formatting (as with yacc) can get kind of ugly. I played with nested embedded functions which might work out.
    * What the grammar should generate. It might be best to have it build an AST because then I can get all the interface{} elements out of the way.

* Nov 01 (1 hr) - Set up project, write first-pass README, begin planning
