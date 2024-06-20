module.exports = grammar({
  name: "okra",

  extras: ($) => [$.comment, /\s/],

  word: ($) => $.identifier,

  rules: {
    program: ($) => repeat($._toplevel_definition),

    _toplevel_definition: ($) =>
      choice(
        alias($.namespace_statement, $.namespace),
        alias($.import_statement, $.import),
        alias($.test_definition, $.test),
        alias($.interface_definition, $.interface),
        alias($.record_definition, $.record),
        alias($.enum_definition, $.enum),
        alias($.function_definition, $.function),
        $._statements,
        $.block,
        $._literals,
      ),

    namespace_statement: ($) => {
      return seq("namespace", field("name", $.identifier), ";");
    },

    import_statement: ($) => {
      let name = choice(
        seq($.identifier, repeat1(seq(".", $.identifier))),
        $.identifier,
      );
      return seq("import", field("name", name), ";");
    },

    test_definition: ($) => {
      let type = $.identifier;
      let func = alias($._test_function, $.function);
      let target = field("target", choice(type, func));
      let desc = field("desc", alias($.string_literal, $.description));
      let body = field("body", alias($._test_block, $.block));
      // body
      let empty = seq("test", optional(target), optional(desc), "end");
      let notempty = seq("test", optional(target), optional(desc), body);
      return choice(empty, notempty);
    },

    _test_block: ($) => seq("do", repeat(choice($._statements)), "end"),

    _test_function: ($) => {
      let name = field("name", $.identifier);
      let params = alias($._test_parameters, $.parameters);
      let ret = field("return", $._type);
      return seq(name, params, optional(ret));
    },

    _test_parameters: ($) => {
      let type = alias($._test_parameter, $.parameter);
      let types = seq("(", optional(type), repeat(seq(",", type)), ")");
      return types;
    },

    _test_parameter: ($) => {
      let type = field("type", $._type);
      return type;
    },

    interface_definition: ($) => {
      let name = field("name", $.identifier);
      let body = field(
        "body",
        choice(
          alias($.interface_block, $.block),
          alias($.interface_case_block, $.block),
        ),
      );
      let start = "interface";
      let empty = seq(start, name, "end");
      let notempty = seq(start, name, body);
      return choice(empty, notempty);
    },

    interface_block: ($) => {
      let func = alias($.interface_function, $.function);
      let funcs = repeat(func);
      return seq("do", funcs, "end");
    },

    interface_case_block: ($) => {
      let func = alias($.interface_function, $.case_function);
      let funcs = repeat1(func);
      return seq("case", funcs, "end");
    },

    interface_function: ($) => {
      let name = field("name", $.identifier);
      let ret = field("return", choice($.void, $.identifier));
      let params = alias($._interfacefn_parameters, $.parameters);
      return seq("fun", name, params, optional(ret), "end");
    },

    _interfacefn_parameters: ($) => {
      let param = alias($._interfacefn_parameter, $.parameter);
      let params = seq(param, repeat(seq(",", param)));
      return seq("(", optional(params), ")");
    },

    _interfacefn_parameter: ($) => {
      let name = field("name", $.identifier);
      let type = field("type", $._type);
      return seq(optional(name), type);
    },

    record_definition: ($) => {
      let name = field("name", $.identifier);
      let body = field("body", alias($.record_block, $.block));
      let start = "record";
      let empty = seq(start, name, "end");
      let noempty = seq(start, name, body);
      return choice(empty, noempty);
    },

    record_block: ($) => {
      let component = alias($._record_component, $.component);
      let components = repeat(component);
      return seq("do", components, "end");
    },

    _record_component: ($) => {
      return seq(field("name", $.identifier), field("type", $._type), ";");
    },

    enum_definition: ($) => {
      let name = field("name", $.identifier);
      let body = field("body", alias($._enum_block, $.block));
      return seq("enum", name, body);
    },

    _enum_block: ($) => {
      let constant = alias($.enum_record, $.record);
      let constants = repeat(constant);
      return seq("do", constants, "end");
    },

    enum_record: ($) => {
      let name = field("name", $.identifier);
      let body = alias($._enum_record_block, $.block);
      let start = "record";
      let empty = seq(start, name, "end");
      let notempty = seq(start, name, body);
      return choice(empty, notempty);
    },

    _enum_record_block: ($) => {
      let component = alias($._enum_component, $.component);
      let components = repeat(seq(component, ";"));
      return seq("do", components, "end");
    },

    _enum_component: ($) => {
      let name = field("name", $.identifier);
      let type = field("type", $._type);
      let value = field("value", seq("=", $._literals));
      return seq(name, optional(type), value);
    },

    function_definition: ($) => {
      let name = field("name", $.identifier);
      let ret = field("return", $._type);
      let body = field("body", $.block);
      let params = alias($._function_parameters, $.parameters);
      return seq("fun", name, params, optional(ret), body);
    },

    _function_parameters: ($) => {
      let param = alias($._function_parameter, $.parameter);
      let params = seq(param, repeat(seq(",", param)));
      return seq("(", optional(params), ")");
    },

    _function_parameter: ($) => {
      let name = field("name", $.identifier);
      let type = field("type", $._type);
      return seq(name, optional(type));
    },

    block: ($) => seq("do", repeat($._statements), "end"),

    _statements: ($) =>
      choice(
        alias($.let_statement, $.let),
        seq(alias($.call_expression, $.call), ";"),
      ),

    let_statement: ($) => {
      let name = field("name", $.identifier);
      let value = field("value", seq("=", $._expression));
      let type = field("type", $._type);
      return seq("let", name, optional(type), optional(value), ";");
    },

    _function_type: ($) => {
      let parameters = alias($._function_type_parameters, $.parameters);
      let ret = field("return", $._type);
      return seq("fn", parameters, optional(ret));
    },

    _function_type_parameters: ($) => {
      let parameter = alias($._function_type_parameter, $.parameter);
      let more = repeat(seq(",", parameter));
      return seq("(", optional(seq(parameter, more)), ")");
    },

    _function_type_parameter: ($) => {
      let parameter = field("type", $.identifier);
      return parameter;
    },

    _expression: ($) => {
      return choice(
        alias($.let_expression, $.let_in),
        alias($._primary_expression, ""),
      );
    },

    _primary_expression: ($) =>
      choice(
        alias($.call_expression, $.call),
        $.member,
        $._literals,
        $.identifier,
      ),

    let_expression: ($) => {
      let bindings = alias($._let_expression_bindings, $.bindings);
      let body = field("body", alias($._let_expression_block, $.block));
      let empty = seq("let", bindings, "in", "end");
      let notempty = seq("let", bindings, body);
      return choice(empty, notempty);
    },

    _let_expression_bindings: ($) => {
      let binding = alias($._let_expression_binding, $.binding);
      let bindings = seq(binding, repeat(seq(",", binding)));
      return bindings;
    },

    _let_expression_binding: ($) => {
      let name = field("name", $.identifier);
      let type = field("type", $._type);
      let value = field("value", choice($._literals, $.identifier));
      return seq(name, optional(type), "=", value);
    },

    _let_expression_block: ($) => {
      return seq("in", repeat1($._statements), "end");
    },

    call_expression: ($) => {
      let name = field("fun", $.identifier);
      let arguments = alias($._call_arguments, $.arguments);
      return seq(name, arguments);
    },

    _call_arguments: ($) => {
      let argument = alias($._call_argument, $.argument);
      let more = repeat(seq(",", argument));
      return seq("(", optional(seq(argument, more)), ")");
    },

    _call_argument: ($) => {
      let argument = $._expression;
      return argument;
    },

    member: ($) => {
      let object = field("object", $._primary_expression);
      let field0 = field("field", $.identifier);
      let method = field("method", alias($.call_expression, $.call));
      return seq(object, ".", choice(method, field0));
    },

    _type: ($) =>
      choice(
        alias($._function_type, $.function),
        alias($._namespaced_type, $.namespaced_type),
        $.void,
        $.identifier,
      ),

    _namespaced_type: ($) => {
      let name = field("namespace", $.identifier);
      let type = field("type", $._type);
      return seq(name, ".", type);
    },

    void: () => "void",

    identifier: () => /[a-zA-Z_]\w*/,

    _literals: ($) =>
      choice(
        $.string_literal,
        $.integer_literal,
        $.nil_literal,
        $.true_literal,
        $.false_literal,
        $.this_literal,
      ),

    string_literal: () => token(seq('"', /[^"\n]*/, '"')),

    integer_literal: () => choice(token(/0/), token(/[1-9][0-9_]*/)),

    nil_literal: () => "nil",

    true_literal: () => "true",

    false_literal: () => "false",

    this_literal: () => "this",

    comment: () => token(repeat1(seq("---", /[^\n]*/, "\n"))),
  },
});
