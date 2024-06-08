import re
from typing import Any, List, Dict

# pydantic.errors.PydanticUserError: Please use `typing_extensions.TypedDict` instead of `typing.TypedDict` on Python < 3.12.
from typing_extensions import TypedDict

from pkg.rpc.server.output_parsers.abstract import OutputParser, CustomBaseModel
from pkg.rpc.server.output_parsers.factory import OutputParserFactory
from pkg.rpc.server.output_parsers.keywords import KeywordsParser
from pkg.rpc.server.output_parsers.summary import SummaryParser
from pkg.rpc.server.prompts.keywords import KeywordsPrompts
from pkg.rpc.server.prompts.summary import SummaryPrompts


class PalaceAddonResultInfo(TypedDict):
    value: Any
    success: bool


class Joined(CustomBaseModel):
    """Data model for joined data models"""

    value: Dict[str, PalaceAddonResultInfo]

    def get_value(self):
        return self.value


class JoinedParser(OutputParser):
    def __init__(self, addons: List[str], verbose=False):
        self.verbose = verbose
        self.addons = addons

    def parse(self, output: str) -> Joined:
        if self.verbose:
            print(f"> Raw output: {output}")

        patterns = []
        parsers: List[OutputParser] = []
        for addon_name in self.addons:
            parser, skip = OutputParserFactory.construct(addon_name)
            if parser is None or skip:
                continue

            parser_instance: OutputParser = parser(verbose=self.verbose)  # type: ignore
            parsers.append(parser_instance)
            patterns.append(parser_instance.pattern)

        pattern = ".*".join(patterns)

        if self.verbose:
            print(f"> Output parser pattern: {pattern}")

        match = re.search(pattern, output, re.DOTALL)

        results: Dict[str, PalaceAddonResultInfo] = {}

        self.success = False
        if match:
            outputs = []
            for parser in parsers:
                output = match.group(parser.group_name)
                outputs.append(output)

            for i, output in enumerate(outputs):
                parser = parsers[i]
                output_model = parser.parse(
                    f"{parser.format_start} {output} {parser.format_end}"
                )

                results[output_model.name] = PalaceAddonResultInfo(
                    value=output_model.value, success=parser.success
                )

        if not self.success:
            pass

        return Joined(value=results)

    @classmethod
    def construct_output(cls, **kwargs) -> str:
        """
        Args:
            summary_input (str): Summary input
            keywords_input (List[str]): Keywords input
        """
        summary_input = kwargs.get("summary_input", "")
        keywords_input = kwargs.get("keywords_input", [])

        summary_output = SummaryParser.construct_output(input=summary_input)
        keywords_output = KeywordsParser.construct_output(input=keywords_input)

        return f"{summary_output}\n{keywords_output}"
