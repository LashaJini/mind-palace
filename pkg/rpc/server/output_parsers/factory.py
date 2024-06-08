from typing import NamedTuple, Type

from pkg.rpc.server.output_parsers.abstract import OutputParser
from pkg.rpc.server.addons.scheme import addons as registry


class OutputModelAndParser(NamedTuple):
    parser: Type[OutputParser]
    skip: bool


class OutputParserFactory:
    @classmethod
    def construct(cls, name: str) -> OutputModelAndParser:
        if name not in registry:
            print(f"Unknown parser {name}")

        value = registry[name]
        parser = value.get("parser")
        skip = value.get("skip")

        return OutputModelAndParser(parser=parser, skip=skip)
