package dev.nyehamene.timber;

public class Tokenizer {
    private final char[] chunks;
    private int current;

    private Tokenizer(char[] chunks) {
        this.chunks = chunks;
    }

    public static Tokenizer New(String source) {
        return new Tokenizer(source.toCharArray());
    }

    public Node advance() throws SyntaxException {
        while(notAtEnd() && Character.isWhitespace(peek())) {
            current++;
        }
        if (isAtEnd()) {
            return Node.EOF();
        }
        char currentChar = chunks[current++];
        return switch(currentChar) {
            case '=' -> Node.New(Symbol.EQUALS);
            case '(' -> Node.New(Symbol.LEFT_PAR);
            case ')' -> Node.New(Symbol.RIGHT_PAR);
            case '{' -> Node.New(Symbol.LEFT_BRACE);
            case '}' -> Node.New(Symbol.RIGHT_BRACE);
            case '[' -> Node.New(Symbol.LEFT_BRACKET);
            case ']' -> Node.New(Symbol.RIGHT_BRACKET);
            case '.' -> Node.New(Symbol.PERIOD);
            case '|' -> Node.New(Symbol.PIPE);
            case '"' -> tokenizeString();
            default -> {
                if (isAlpha(currentChar)) {
                    yield tokenizeName(currentChar);
                }
                throw SyntaxException.InvalidCharacter(currentChar);
            }
        };
    }

    private Node tokenizeName(char start) {
        var builder = new StringBuilder();
        builder.append(start);
        while(isAlpha()) {
            builder.append(advance0());
        }
        return Node.Name(builder.toString());
    }

    private Node tokenizeString() throws SyntaxException {
        var builder = new StringBuilder();
        while(notAtEnd() && peek() != '"') {
            char c = advance0();
            if(c == '\\') {
                if (peek() == '\\') {
                    builder.append(advance0());
                } else if(peek() == '"') {
                    builder.append(advance0());
                } else {
                    // reject other escapes. test tbd.
                    throw SyntaxException.IllegalCharacterEscape(peek());
                }
            } else {
                builder.append(c);
            }
        }

        if (isAtEnd()) {
            throw SyntaxException.UnterminatedString(builder.toString());
        }
        // consume closing '"'
        current++;
        return Node.String(builder.toString());
    }

    private char advance0() {
        return chunks[current++];
    }

    private char peek() {
        return chunks[current];
    }

    private boolean isAtEnd() {
        return !(current < chunks.length);
    }

    private boolean notAtEnd() {
        return !isAtEnd();
    }

    private boolean isAlpha() {
        return notAtEnd() && isAlpha(peek());
    }

    private boolean isAlpha(char current) {
        return (current >= 'A' && current <= 'Z') ||
               (current >= 'a' && current <= 'z') ||
               current == '_';
    }
}
