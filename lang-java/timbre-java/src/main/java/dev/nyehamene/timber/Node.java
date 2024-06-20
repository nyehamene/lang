package dev.nyehamene.timber;

import java.util.Objects;

public record Node(NodeType type, String value) {

    public Node {
        Objects.requireNonNull(type);
    }

    public static Node New(Symbol type) {
        return new Node(type, String.valueOf(type.character));
    }

    public static Node String(String value) {
        return new Node(Value.STRING, value);
    }

    public static Node Name(String value) {
        return new Node(Value.NAME, value);
    }

    public static Node EOF() {
        return New(Symbol.EOF);
    }

    @Override
    public String toString() {
        return STR."Node(type: \{type}; value: '\{value}')";
    }
}

sealed interface NodeType {}

enum Symbol implements NodeType {
    EQUALS('='),
    EOF('\0'),
    LEFT_BRACE('{'),
    LEFT_BRACKET('['),
    LEFT_PAR('('),
    PERIOD('.'),
    PIPE('|'),
    RIGHT_BRACE('}'),
    RIGHT_BRACKET(']'),
    RIGHT_PAR(')'),
    ;

    public final char character;

    Symbol(char c) {
        character = c;
    }
}

enum Value implements NodeType {
    NAME, STRING;
}
