module.exports = grammar({
  name: "penny",

  extras: ($) => [$.comment, /\s/, /\n/],

  word: ($) => $.identifier,

  // inline: ($) => [$._expression],

  // TODO: remove conflict for expression and function_call
  conflicts: ($) => [
    [$._expression, $.function_call],
    // [$._expression, $.block_function_call],
    [$.block_function_call, $.function_call],
    // [$._expression, $.function_call, $.block_function_call],
  ],

  rules: {
    program: ($) => {
      return repeat($._declaration);
    },

    _declaration: ($) =>
      choice($.function_declaration, $.example_declaration, $.test_declaration),

    function_declaration: ($) => {
      return seq(
        optional($.prefix_documentation),
        $.function,
        optional($.postfix_documentation),
      );
    },

    function: ($) =>
      seq(
        "fun",
        field("name", $.identifier),
        optional(repeat(field("parameter", $.identifier))),
        "=",
        field("return", choice($.void_type, $.identifier)),
        "do",
        optional(field("body", $.block)),
        "end",
      ),

    example_declaration: ($) =>
      seq(
        optional($.prefix_documentation),
        $.example,
        optional($.postfix_documentation),
      ),

    example: ($) =>
      seq(
        "example",
        optional(
          seq(
            field("function", $.identifier),
            repeat(field("parameter_type", $.identifier)),
          ),
        ),
        optional(field("desc", $.string_literal)),
        "do",
        optional(field("body", $.block)),
        "end",
      ),

    test_declaration: ($) =>
      seq(
        optional($.prefix_documentation),
        $.test,
        optional($.postfix_documentation),
      ),

    test: ($) =>
      seq(
        "test",
        optional(
          seq(
            field("function", $.identifier),
            repeat(field("parameter_type", $.identifier)),
          ),
        ),
        optional(field("desc", $.string_literal)),
        "do",
        optional(field("body", $.block)),
        "end",
      ),

    block: ($) => repeat1($._statements),

    void_type: () => "void",

    _statements: ($) =>
      choice(
        $.let_statement,
        $.return_statement,
        alias($._expression_statement, $.statement),
      ),

    let_statement: ($) =>
      seq(
        "let",
        field("name", $.identifier),
        "=",
        field("value", $._expression),
        ";",
      ),

    return_statement: ($) => seq("return", $._expression, ";"),

    block_function_call: ($) =>
      seq(
        field("name", choice($._literals, $.identifier)),
        repeat(
          field(
            "argument",
            // choice(
            //   $.block_function_call,
            //   $.function_call,
            //   $.try_expression,
            //   $.field_access,
            //   $.identifier,
            //   $._literals,
            // ),
            $._expression,
          ),
        ),
        "do",
        optional($.block),
        "end",
      ),

    _expression_statement: ($) =>
      choice($.block_function_call, seq($._expression, ";")),

    _expression: ($) =>
      choice(
        $.try_expression,
        $.field_access,
        $.function_call,
        $.identifier,
        $._literals,
      ),

    try_expression: ($) => prec.left(seq("try", $.function_call)),

    function_call: ($) =>
      prec.left(
        seq(
          choice(field("name", $._literals, $.identifier)),
          field("argument", repeat($._expression)),
        ),
      ),

    field_access: ($) => {
      let start = choice($._literals, $.identifier);
      return seq(field("object", start), ".", $.identifier);
    },

    identifier: () => /[a-zA-Z_]\w*/,

    string_literal: () => seq('"', /[^"]*/, '"'),

    numeric_literal: () => /[1-9][0-9_]+/,

    true_literal: () => "true",

    false_literal: () => "false",

    nil_literal: () => "nil",

    _literals: ($) =>
      choice(
        $.string_literal,
        $.numeric_literal,
        $.true_literal,
        $.false_literal,
        $.nil_literal,
      ),

    prefix_documentation: () => token(repeat1(seq("--*", /[^\n]*/, "\n"))),

    postfix_documentation: () => token(repeat1(seq("--+", /[^\n]*/, "\n"))),

    comment: () => token(repeat1(seq("---", /[^\n]*/, "\n"))),
  },
});
