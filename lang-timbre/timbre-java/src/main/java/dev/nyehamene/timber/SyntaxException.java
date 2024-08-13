package dev.nyehamene.timber;

public class SyntaxException extends Exception {

    private SyntaxException(String message) {
        super(message);
    }

    public static SyntaxException InvalidCharacter(char invalid) {
        return new SyntaxException(String.valueOf(invalid));
    }

    public static SyntaxException IllegalCharacterEscape(char illegal) {
        return new SyntaxException(STR."\\\{illegal}");
    }

    public static SyntaxException UnterminatedString(String unterminated) {
        return new SyntaxException(STR."\"\{unterminated}");
    }
}
