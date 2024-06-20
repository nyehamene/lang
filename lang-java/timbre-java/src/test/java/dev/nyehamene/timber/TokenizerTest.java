package dev.nyehamene.timber;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatNoException;
import static org.junit.jupiter.api.Assertions.assertAll;

import org.assertj.core.api.Assertions;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.EnumSource;

public class TokenizerTest {

    @Test
    void emptySource() {
        whitespace("");
    }

    @Test
    void whitespace() {
        whitespace("  ");
    }

    void whitespace(String source) {
        assertThatNoException().isThrownBy(() -> {
            // Given
            var tokenizer = Tokenizer.New(source);
            // When
            var node = tokenizer.advance();
            // Then
            var expected = Node.EOF();
            assertThat(node).isEqualTo(expected);
        });
    }

    @Nested
    class TokenizeSpecialSymbols {

        @Test
        void tokenizeColon() {
            tokenizeChar('=', Symbol.EQUALS);
        }

        @Test
        void tokenizeParenthesis() {
            tokenizeChar('(', Symbol.LEFT_PAR);
        }

        @Test
        void tokenizeRightPar() {
            tokenizeChar(')', Symbol.RIGHT_PAR);
        }

        @Test
        void tokenizeLeftBrace() {
            tokenizeChar('{', Symbol.LEFT_BRACE);
        }

        @Test
        void tokenizeRightBrace() {
            tokenizeChar('}', Symbol.RIGHT_BRACE);
        }

        @Test
        void tokenizeLeftBracket() {
            tokenizeChar('[', Symbol.LEFT_BRACKET);
        }

        @Test
        void tokenizeRightBracket() {
            tokenizeChar(']', Symbol.RIGHT_BRACKET);
        }

        @Test
        void tokenizeSemiColon() {
            tokenizeChar('.', Symbol.PERIOD);
        }

        @Test
        void tokenizePipe() {
            tokenizeChar('|', Symbol.PIPE);
        }

        @Test
        void leadingWhitespace() {
            assertAll(
                    () -> tokenizeChar(" =", Symbol.EQUALS),
                    () -> tokenizeChar(" (", Symbol.LEFT_PAR),
                    () -> tokenizeChar(" )", Symbol.RIGHT_PAR),
                    () -> tokenizeChar(" {", Symbol.LEFT_BRACE),
                    () -> tokenizeChar(" }", Symbol.RIGHT_BRACE),
                    () -> tokenizeChar(" [", Symbol.LEFT_BRACKET),
                    () -> tokenizeChar(" ]", Symbol.RIGHT_BRACKET),
                    () -> tokenizeChar(" .", Symbol.PERIOD),
                    () -> tokenizeChar(" |", Symbol.PIPE));
        }

        @Test
        @Disabled("""
                The tokenizer advances by one token so never reaches
                the trailing whitespace when the test is executed
                """)
        void trailingWhitespace() {
            assertAll(
                    () -> tokenizeChar("= ", Symbol.EQUALS),
                    () -> tokenizeChar("( ", Symbol.LEFT_PAR),
                    () -> tokenizeChar(") ", Symbol.RIGHT_PAR),
                    () -> tokenizeChar("{ ", Symbol.LEFT_BRACE),
                    () -> tokenizeChar("} ", Symbol.RIGHT_BRACE),
                    () -> tokenizeChar("[ ", Symbol.LEFT_BRACKET),
                    () -> tokenizeChar("] ", Symbol.RIGHT_BRACKET),
                    () -> tokenizeChar(". ", Symbol.PERIOD),
                    () -> tokenizeChar("| ", Symbol.PIPE));
        }

        @Test
        @Disabled("The tokenizer never reaches the trailing whitespace")
        void bothLeadingAndTralingWhitespace() {
            assertAll(
                    () -> tokenizeChar(" = ", Symbol.EQUALS),
                    () -> tokenizeChar(" ( ", Symbol.LEFT_PAR),
                    () -> tokenizeChar(" ) ", Symbol.RIGHT_PAR),
                    () -> tokenizeChar(" { ", Symbol.LEFT_BRACE),
                    () -> tokenizeChar(" } ", Symbol.RIGHT_BRACE),
                    () -> tokenizeChar(" [ ", Symbol.LEFT_BRACKET),
                    () -> tokenizeChar(" ] ", Symbol.RIGHT_BRACKET),
                    () -> tokenizeChar(" . ", Symbol.PERIOD),
                    () -> tokenizeChar(" | ", Symbol.PIPE));
        }

        private void tokenizeChar(char character, Symbol match) {
            tokenizeChar(String.valueOf(character), match);
        }

        private void tokenizeChar(String source, Symbol match) {
            assertThatNoException().isThrownBy(() -> {
                // Given
                var tokenizer = Tokenizer.New(source);
                // When
                var node = tokenizer.advance();
                // Then
                var expected = Node.New(match);
                assertThat(node).isEqualTo(expected);
            });
        }
    }

    @Nested
    class TokenizeDoubleQuoteStrings {

        @Test
        void emptyString() {
            tokenizeString("\"\"", "");
        }

        @Test
        void nonemptyString() {
            tokenizeString("\"timber\"", "timber");
        }

        @ParameterizedTest
        @EnumSource(Symbol.class)
        void containingSpecialSymbols(Symbol symbol) {
            var source = String.valueOf(symbol.character);
            tokenizeString(STR."\"\{source}\"", source);
        }

        @Test
        void escapedDoubleQuote() {
            var source = "\"\\\"\"";
            tokenizeString(source, "\"");
        }

        @Test
        void escaptedForwardSlash() {
            var source = "\"\\\\\"";
            tokenizeString(source, "\\");
        }

        @Test
        void unterminatedString() {
            failToTokenize("\"");
            failToTokenize("\"timber");
        }

        @Test
        void leadingWhitespace() {
            assertAll(
                    () -> tokenizeString(" \"timber\"", "timber"),
                    () -> tokenizeString(" \"TIMBER\"", "TIMBER"),
                    () -> tokenizeString(" \"timber_land\"", "timber_land"));
        }

        @Test
        @Disabled("""
                The tokenizer advances by one token so never reaches
                the trailing whitespace when the test is executed
                """)
        void trailingWhitespace() {
            assertAll(
                    () -> tokenizeString("\"timber\" ", "timber"),
                    () -> tokenizeString("\"TIMBER\" ", "TIMBER"),
                    () -> tokenizeString("\"timber_land\" ", "timber_land"));
        }

        @Test
        @Disabled("The tokenizer never reaches the trailing whitespace")
        void bothLeadingAndTralingWhitespace() {
            assertAll(
                    () -> tokenizeString(" \"timber\" ", "timber"),
                    () -> tokenizeString(" \"TIMBER\" ", "TIMBER"),
                    () -> tokenizeString(" \"timber_land\" ", "timber_land"));
        }

        private void failToTokenize(String source) {
            // Given
            var tokenizer = Tokenizer.New(source);
            // Then
            Assertions.assertThatExceptionOfType(SyntaxException.class)
                    .isThrownBy(() -> /* When */ tokenizer.advance());
        }

        private void tokenizeString(String source, String match) {
            assertThatNoException().isThrownBy(() -> {
                // Given
                var tokenizer = Tokenizer.New(source);
                // When
                var node = tokenizer.advance();
                // Then
                var expected = Node.String(match);
                assertThat(node).isEqualTo(expected);
            });
        }
    }

    @Nested
    class TokenizeName {

        @Test
        void lowerCase() {
            tokenizeName("timber");
        }

        @Test
        void upperCase() {
            tokenizeName("TIMBER");
        }

        @Test
        void leadingUnderScore() {
            tokenizeName("_timber");
        }

        @Test
        void trailingUnderScore() {
            tokenizeName("timber_");
        }

        @Test
        void snakeCase() {
            tokenizeName("timber_land");
        }

        @Test
        void leadingWhitespace() {
            assertAll(
                    () -> tokenizeName(" timber", "timber"),
                    () -> tokenizeName(" TIMBER", "TIMBER"),
                    () -> tokenizeName(" timber_land", "timber_land"));
        }

        @Test
        @Disabled("""
                The tokenizer advances by one token so never reaches
                the trailing whitespace when the test is executed
                """)
        void trailingWhitespace() {
            assertAll(
                    () -> tokenizeName("timber ", "timber"),
                    () -> tokenizeName("TIMBER ", "TIMBER"),
                    () -> tokenizeName("timber_land ", "timber_land"));
        }

        @Test
        @Disabled("The tokenizer never reaches the trailing whitespace")
        void bothLeadingAndTralingWhitespace() {
            assertAll(
                    () -> tokenizeName(" timber ", "timber"),
                    () -> tokenizeName(" TIMBER ", "TIMBER"),
                    () -> tokenizeName(" timber_land ", "timber_land"));
        }

        private void tokenizeName(String source) {
            tokenizeName(source, source);
        }

        private void tokenizeName(String source, String expectedValue) {
            assertThatNoException().isThrownBy(() -> {
                // Given
                var tokenizer = Tokenizer.New(source);
                // When
                var node = tokenizer.advance();
                // Then
                var expected = Node.Name(expectedValue);
                assertThat(node).isEqualTo(expected);
            });
        }
    }
}
