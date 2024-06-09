import re
from typing import ClassVar, List
from pydantic import ConfigDict

from pkg.rpc.server import addon_names
from pkg.rpc.server.output_parsers.abstract import OutputParser, CustomBaseModel
from pkg.rpc.server.prompts.summary import SummaryPrompts


class Summary(CustomBaseModel):
    """Data model for a summary."""

    model_config = ConfigDict(from_attributes=True)

    value: List[str]
    name: ClassVar[str] = addon_names.summary

    def get_value(self) -> str:
        return self.value[0]


class SummaryParser(OutputParser):
    def __init__(self, verbose: bool = False):
        format_start = SummaryPrompts.format_start
        format_end = SummaryPrompts.format_end
        group_name = "summary"
        pattern = rf"{format_start}\s*(?P<{group_name}>.*?)\s*{format_end}"

        super().__init__(
            format_start=format_start,
            format_end=format_end,
            group_name=group_name,
            pattern=pattern,
            verbose=verbose,
        )

    def parse(self, output: str) -> Summary:
        if self.verbose:
            print(f"> Raw output: {output}")

        self.success = False
        summary = ""
        match = re.search(self.pattern, output, re.DOTALL)
        if match:
            summary = match.group(self.group_name)
            self.success = len(summary) > 0

        if not self.success:
            e = (
                f"Summary output format should be: {self.format_start} <{self.group_name}> {self.format_end}\n",
                f"Got: {output}",
            )
            raise ValueError(e)

        if self.verbose:
            print(f"> Extracted summary: {summary}")

        return Summary(value=[summary])

    @classmethod
    def construct_output(cls, **kwargs) -> str:
        """
        Args:
            input (str): Input text.
        """
        input = kwargs.get("input", "")
        return f"{SummaryPrompts.format_start} {input} {SummaryPrompts.format_end}"
