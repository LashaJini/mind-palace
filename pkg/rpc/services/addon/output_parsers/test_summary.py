import pytest

from pkg.rpc.services.addon.output_parsers.summary import SummaryParser


class TestSummaryParser:
    @pytest.fixture(scope="class")
    def parser(self):
        return SummaryParser()

    @pytest.mark.parametrize(
        "invalid_input",
        [
            "SUMMARY: oi",
            "oi SUMMARY_END",
            "SUMMARY: SUMMARY_END",
            "oi",
            "SUMMARY_END",
            "SUMMARY: ",
        ],
    )
    def test_invalid_summary_output_raises_errors(
        self, parser: SummaryParser, invalid_input
    ):
        with pytest.raises(ValueError):
            parser.parse(invalid_input)

        assert not parser.success

    def test_valid_output_returns_summary(self, parser: SummaryParser):
        small_summary = "small summary"

        output = SummaryParser.construct_output(input=small_summary)
        assert parser.parse(output).summary == small_summary
        assert parser.success

        summary_with_new_lines = "summary summary\nsummary\nsummary"

        output = SummaryParser.construct_output(input=summary_with_new_lines)
        assert parser.parse(output).summary == summary_with_new_lines
        assert parser.success
