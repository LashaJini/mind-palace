import re
from typing import ClassVar

from pydantic import PrivateAttr

import pkg.rpc.server.gen.Palace_pb2 as pbPalace
from pkg.rpc.server import addon_names, logger
from pkg.rpc.server.output_parsers.abstract import (
    OutputParser,
    CustomBaseModel,
)
from pkg.rpc.server.prompts.summary import SummaryPrompts


class Summary(CustomBaseModel):
    name: ClassVar[str] = addon_names.summary
    _summary: str = PrivateAttr("")

    class Config:
        arbitrary_types_allowed = True

    def __init__(self, summary: str = "", **kwargs):
        super().__init__(**kwargs)

        self._summary = summary

    @property
    def summary(self) -> str:
        return self._summary

    def to_addon_result(self) -> pbPalace.AddonResult:
        response = pbPalace.AddonResponse(
            summary_response=pbPalace.SummaryResponse(summary=self._summary),
            success=self.success,
        )
        return pbPalace.AddonResult(map={Summary.name: response})


class SummaryParser(OutputParser):
    def __init__(self, verbose: bool = False, **kwargs):
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
            **kwargs,
        )

    def parse(self, output: str) -> Summary:
        if self.verbose:
            logger.log.debug(f"> Raw output: {output}")

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
            logger.log.debug(f"> Extracted summary: {summary}")

        return Summary(summary=summary, success=self.success)

    @classmethod
    def construct_output(cls, **kwargs) -> str:
        """
        Args:
            input (str): Input text.
        """
        input = kwargs.get("input", "")
        return f"{SummaryPrompts.format_start} {input} {SummaryPrompts.format_end}"
