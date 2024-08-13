package example;

import java.io.IOException;
import dev.nyehamene.timber.Tokenizer;
import dev.nyehamene.timber.Node;
import dev.nyehamene.timber.SyntaxException;

class Example {

    static final String EBNF = """
        grammar = rules .
        rules = rule rules | rule .
        rule = NAME "=" alternatives .
        alternatives =
          alternative "|" alternatives |
          alternative
          .
        alternative = terms .
        terms = term terms | term.
        term = NAME | STRING | group .
        group =
          "(" alternatives ")" |
          "{" alternatives "}" |
          "[" alternatives "]"
          .
          """;

    public static void main(String[] args) throws IOException, SyntaxException {
        tokenizer(EBNF);
    }

    private static void tokenizer(String grammar) throws IOException, SyntaxException {
        var tokenizer = Tokenizer.New(grammar);
        Node tokenized = null;
        do {
            tokenized = tokenizer.advance();
            System.out.println(tokenized);
        } while (!tokenized.equals(Node.EOF()));
    }
}