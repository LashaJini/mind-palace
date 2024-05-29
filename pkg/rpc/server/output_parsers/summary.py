from typing import ClassVar, List
from pydantic import ConfigDict

from pkg.rpc.server.output_parsers.abstract import OutputParser, CustomBaseModel
from pkg.rpc.server.prompts.summary import SummaryPrompts


class Summary(CustomBaseModel):
    """Data model for a summary."""

    model_config = ConfigDict(from_attributes=True)

    value: List[str]
    name: ClassVar[str] = "summary"


class SummaryParser(OutputParser):
    def __init__(self, verbose: bool = False):
        format_start = SummaryPrompts.format_start
        format_end = SummaryPrompts.format_end
        group_name = "summary"

        super().__init__(
            format_start=format_start,
            format_end=format_end,
            group_name=group_name,
            pattern=rf"{format_start}\s*(?P<{group_name}>.*?)\s*{format_end}",
            verbose=verbose,
        )

    def parse(self, output: str) -> Summary:
        if self.verbose:
            print(f"> Raw output: {output}")

        summary_part = output.split(self.format_start, 1)[1]
        summary = summary_part.split(self.format_end, 1)[0].strip()

        return Summary(value=[summary])
