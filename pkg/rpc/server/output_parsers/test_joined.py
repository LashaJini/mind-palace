import pytest

from pkg.rpc.server.output_parsers.joined import JoinedParser
from pkg.rpc.server.output_parsers.keywords import KeywordsParser
from pkg.rpc.server.output_parsers.summary import SummaryParser


verbose = True


def validate_keywords(parser, output, expected_keywords, expected_success=True):
    success = False
    keywords = []

    parser_result = parser.parse(output).get_value()
    keywords_result = parser_result.get("keywords")

    if keywords_result is not None:
        keywords = keywords_result.get("value")
        success = keywords_result.get("success")

    assert keywords == expected_keywords
    assert success
    assert success == expected_success


def validate_summary(parser, output, expected_summary, expected_success=True):
    success = False
    summary = ""

    parser_result = parser.parse(output).get_value()
    summary_result = parser_result.get("summary")

    if summary_result is not None:
        summary = summary_result.get("value")[0]
        success = summary_result.get("success")

    assert summary == expected_summary
    assert success == expected_success


def validate_everything(parser, output, **kwargs):
    validate_summary(
        parser,
        output,
        kwargs.get("expected_summary"),
        kwargs.get("expected_summary_success", True),
    )
    validate_keywords(
        parser,
        output,
        kwargs.get("expected_keywords"),
        kwargs.get("expected_keywords_success", True),
    )


class TestJoinedParser:
    @pytest.mark.parametrize(
        "input_keywords, small_summary, summary_with_new_lines",
        [
            (
                ["key1", "key 2", "key 3", "key4"],
                "small summary",
                "summary summary\nsummary\nsummary",
            )
        ],
    )
    def test_valid_output_returns_joined(
        self, input_keywords, small_summary, summary_with_new_lines
    ):
        parser = JoinedParser(verbose=verbose, addons=["mind-palace-resource-keywords"])
        validate_keywords(
            parser=parser,
            output=KeywordsParser.construct_output(input=input_keywords),
            expected_keywords=input_keywords,
        )

        parser = JoinedParser(verbose=verbose, addons=["mind-palace-resource-summary"])
        validate_summary(
            parser=parser,
            output=SummaryParser.construct_output(input=small_summary),
            expected_summary=small_summary,
        )

        validate_summary(
            parser=parser,
            output=SummaryParser.construct_output(input=summary_with_new_lines),
            expected_summary=summary_with_new_lines,
        )

        parser = JoinedParser(
            verbose=verbose,
            addons=["mind-palace-resource-summary", "mind-palace-resource-keywords"],
        )
        validate_everything(
            parser=parser,
            output=JoinedParser.construct_output(
                summary_input=summary_with_new_lines, keywords_input=input_keywords
            ),
            expected_keywords=input_keywords,
            expected_summary=summary_with_new_lines,
        )
