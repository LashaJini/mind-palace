from typing import List

from pkg.rpc.server import logger
import pkg.rpc.server.gen.Palace_pb2 as pbPalace
from pkg.rpc.server.output_parsers.abstract import (
    OutputParser,
    CustomBaseModel,
)
from pkg.rpc.server.output_parsers.keywords import KeywordsParser
from pkg.rpc.server.output_parsers.summary import SummaryParser


class Joined(CustomBaseModel):
    _addon_results: List[pbPalace.AddonResult] = []

    class Config:
        arbitrary_types_allowed = True

    def __init__(self, addon_results: List[pbPalace.AddonResult] = [], **kwargs):
        super().__init__(**kwargs)

        self._addon_results = addon_results

    def to_addon_result(self) -> pbPalace.AddonResult:
        result = pbPalace.AddonResult()
        for addon_result in self._addon_results:
            result.MergeFrom(addon_result)

        return result


class JoinedParser(OutputParser):
    def __init__(self, parsers: List[OutputParser], **kwargs):
        super().__init__(**kwargs)

        self._parsers = parsers

    def parse(self, output: str) -> Joined:
        if self.verbose:
            logger.log.debug(f"> Raw output:\n{output}")

        addon_results: List[pbPalace.AddonResult] = []

        for parser in self._parsers:
            output_model = parser.parse(output)
            addon_result = output_model.to_addon_result()

            addon_results.append(addon_result)

        return Joined(addon_results=addon_results, success=True)

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
