import pytest

from pkg.rpc.server.output_parsers.keywords import KeywordsParser


class TestKeywordsParser:
    @pytest.fixture(scope="class")
    def parser(self):
        return KeywordsParser(verbose=True)

    @pytest.mark.parametrize(
        "invalid_input",
        [
            "KEYWORDS: oi",
            "oi KEYWORDS_END",
            "KEYWORDS: KEYWORDS_END",
            "oi",
            "KEYWORDS_END",
            "KEYWORDS: ",
        ],
    )
    def test_invalid_keywords_output_raises_errors(
        self, parser: KeywordsParser, invalid_input
    ):
        with pytest.raises(ValueError):
            parser.parse(invalid_input)

        assert not parser.success

    def test_valid_output_returns_keywords(self, parser: KeywordsParser):
        keywords = ["key1", " key 2", "key 3", "key4"]
        expected_keywords = ["key1", "key 2", "key 3", "key4"]

        output = KeywordsParser.construct_output(input=keywords)
        assert parser.parse(output).keywords == [expected_keywords]
        assert parser.success
