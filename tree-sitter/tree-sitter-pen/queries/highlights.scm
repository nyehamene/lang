[
  "do"
  "end"
  "return"
  "let"
] @keyword

[
  "fun"
  "example"
  "test"
] @keyword.function

[
  (true_literal)
  (false_literal)
] @constant.builtin.boolean

(void_type) @type.builtin

(nil_literal) @constant.builtin.nil
(numeric_literal) @constant.builtin.numeric

(string_literal) @string.quoted.double

(comment) @comment

(prefix_documentation) @comment.block.documentation.prefix
(postfix_documentation) @comment.block.documentation.postfix

(function
    name: (identifier) @function)

(example
    function: (identifier) @function)

(example
    desc: (string_literal) @string.quoted.double)

