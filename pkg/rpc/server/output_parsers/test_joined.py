from typing import List
import pytest

from pkg.rpc.server import addon_names
from pkg.rpc.server.output_parsers.joined import JoinedParser
from pkg.rpc.server.output_parsers.keywords import KeywordsParser
from pkg.rpc.server.output_parsers.summary import SummaryParser


verbose = True


def validate_keywords(
    parser: JoinedParser,
    output: str,
    expected_keywords: List[List[str]],
    expected_success=True,
):
    success = False
    keywords = []

    keywords_result = (
        parser.parse(output).to_addon_result().map.get(addon_names.keywords)
    )

    if keywords_result is not None:
        keywords_chunks = keywords_result.keywords_response.list
        keywords = [keyword_chunk.keywords for keyword_chunk in keywords_chunks]
        success = keywords_result.success

    assert keywords == expected_keywords
    assert success
    assert success == expected_success


def validate_summary(
    parser: JoinedParser, output: str, expected_summary: str, expected_success=True
):
    success = False
    summary = ""

    summary_result = parser.parse(output).to_addon_result().map.get(addon_names.summary)

    if summary_result is not None:
        summary = summary_result.summary_response.summary
        success = summary_result.success

    assert summary == expected_summary
    assert success == expected_success


def validate_everything(
    parser,
    output: str,
    expected_keywords: List[List[str]],
    expected_summary: str,
):
    validate_summary(
        parser=parser,
        output=output,
        expected_summary=expected_summary,
        expected_success=True,
    )
    validate_keywords(
        parser=parser,
        output=output,
        expected_keywords=expected_keywords,
        expected_success=True,
    )


class TestJoinedParser:
    @pytest.mark.parametrize(
        "input_keywords, chunks, small_summary, summary_with_new_lines",
        [
            (
                ["key1", "key 2", "key 3", "key4"],
                [""],
                "small summary",
                "summary summary\nsummary\nsummary",
            )
        ],
    )
    def test_valid_output_returns_joined(
        self,
        input_keywords: List[str],
        chunks: List[str],
        small_summary: str,
        summary_with_new_lines: str,
    ):
        keywords_parser = KeywordsParser()
        keywords_parser.chunks = chunks
        summary_parser = SummaryParser()

        parser = JoinedParser(
            parsers=[keywords_parser],
            verbose=verbose,
        )
        validate_keywords(
            parser=parser,
            output=KeywordsParser.construct_output(input=input_keywords),
            expected_keywords=[input_keywords],
        )

        parser = JoinedParser(
            parsers=[summary_parser],
            verbose=verbose,
        )
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
            parsers=[summary_parser, keywords_parser],
            verbose=verbose,
        )
        validate_everything(
            parser=parser,
            output=JoinedParser.construct_output(
                summary_input=summary_with_new_lines, keywords_input=input_keywords
            ),
            expected_keywords=[input_keywords],
            expected_summary=summary_with_new_lines,
        )
