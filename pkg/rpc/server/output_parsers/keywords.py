from typing import ClassVar, List
from pydantic import ConfigDict

from pkg.rpc.server.output_parsers.abstract import OutputParser, CustomBaseModel
from pkg.rpc.server.prompts.keywords import KeywordsPrompts


class Keywords(CustomBaseModel):
    """Data model for a keywords."""

    model_config = ConfigDict(from_attributes=True)

    value: List[str]
    name: ClassVar[str] = "keywords"


class KeywordsParser(OutputParser):
    def __init__(self, verbose: bool = False):
        format_start = KeywordsPrompts.format_start
        format_end = KeywordsPrompts.format_end
        group_name = "keywords"

        super().__init__(
            format_start=format_start,
            format_end=format_end,
            group_name=group_name,
            pattern=rf"{format_start}\s*(?P<{group_name}>.*?)\s*{format_end}",
            verbose=verbose,
        )

    def parse(self, output: str) -> Keywords:
        if self.verbose:
            print(f"> Raw output: {output}")

        keywords_part = output.split(self.format_start, 1)[1]
        keywords = keywords_part.split(self.format_end, 1)[0].strip()
        if self.verbose:
            print(f"> Extracted keywords: {keywords}")
        keywords = [keyword.strip().lower() for keyword in keywords.split(",")]

        return Keywords(value=keywords)
