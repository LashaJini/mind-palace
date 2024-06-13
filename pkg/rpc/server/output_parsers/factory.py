from typing import Type

from pkg.rpc.server.output_parsers.abstract import OutputParser
from pkg.rpc.server.addons.scheme import addons as registry


class OutputParserFactory:
    @classmethod
    def construct(cls, name: str) -> Type[OutputParser]:
        if name not in registry:
            print(f"Unknown parser {name}")

        value = registry[name]
        parser = value.get("parser")

        return parser
