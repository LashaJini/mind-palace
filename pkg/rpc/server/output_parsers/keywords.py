import re
from typing import ClassVar, List
from pydantic import ConfigDict

from pkg.rpc.server.output_parsers.abstract import OutputParser, CustomBaseModel
from pkg.rpc.server.prompts.keywords import KeywordsPrompts


class Keywords(CustomBaseModel):
    """Data model for a keywords."""

    model_config = ConfigDict(from_attributes=True)

    value: List[str]
    name: ClassVar[str] = "keywords"

    def get_value(self) -> List[str]:
        return self.value


class KeywordsParser(OutputParser):
    def __init__(self, verbose: bool = False):
        format_start = KeywordsPrompts.format_start
        format_end = KeywordsPrompts.format_end
        group_name = "keywords"
        pattern = rf"{format_start}\s*(?P<{group_name}>.*?)\s*{format_end}"

        super().__init__(
            format_start=format_start,
            format_end=format_end,
            group_name=group_name,
            pattern=pattern,
            verbose=verbose,
        )

    def parse(self, output: str) -> Keywords:
        if self.verbose:
            print(f"> Raw output: {output}")

        self.success = False
        keywords = []
        match = re.search(self.pattern, output, re.DOTALL)
        if match:
            keywords = match.group(self.group_name)
            keywords = [
                keyword.strip().lower() for keyword in keywords.split(",") if keyword
            ]
            self.success = len(keywords) > 0

        if not self.success:
            e = (
                f"Keywords output format should be: {self.format_start} <{self.group_name}> {self.format_end}\n",
                f"Got: {output}",
            )
            raise ValueError(e)

        if self.verbose:
            print(f"> Extracted keywords: {keywords}")

        return Keywords(value=keywords)
