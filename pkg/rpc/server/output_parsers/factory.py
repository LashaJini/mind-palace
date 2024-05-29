from typing import NamedTuple, Union

from pkg.rpc.server.output_parsers.abstract import OutputParser
from pkg.rpc.server.output_parsers.default import DefaultParser
from pkg.rpc.server.output_parsers.summary import SummaryParser
from pkg.rpc.server.output_parsers.keywords import KeywordsParser


class OutputModelAndParser(NamedTuple):
    parser: Union[OutputParser, None]
    skip: bool


class OutputParserFactory:
    @classmethod
    def construct(cls, name: str) -> OutputModelAndParser:
        if name not in registry:
            print(f"Unknown parser {name}")

        value = registry[name]
        parser = value.get("parser")
        skip = value.get("skip", True)

        return OutputModelAndParser(parser=parser, skip=skip)


registry = {
    "mind-palace-default": {
        "parser": DefaultParser,
        "skip": True,
    },
    "mind-palace-resource-summary": {
        "parser": SummaryParser,
        "skip": False,
    },
    "mind-palace-resource-keywords": {
        "parser": KeywordsParser,
        "skip": False,
    },
}
