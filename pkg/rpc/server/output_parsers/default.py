from typing import ClassVar

import pkg.rpc.server.gen.Palace_pb2 as pbPalace
from pkg.rpc.server import addon_names
from pkg.rpc.server.output_parsers.abstract import (
    OutputParser,
    CustomBaseModel,
)


class Default(CustomBaseModel):
    name: ClassVar[str] = addon_names.default
    _default: str = ""

    class Config:
        arbitrary_types_allowed = True

    def __init__(self, default: str = "", **kwargs):
        super().__init__(**kwargs)

        self._default = default

    @property
    def default(self) -> str:
        return self._default

    @default.setter
    def default(self, default: str):
        self._default = default

    def to_addon_result(self) -> pbPalace.AddonResult:
        response = pbPalace.AddonResponse(
            default_response=pbPalace.DefaultResponse(default=self._default),
            success=True,
        )

        return pbPalace.AddonResult(map={Default.name: response})


class DefaultParser(OutputParser):
    _default: str = ""

    def __init__(self, verbose: bool = False, **kwargs):
        format_start = "DEFAULT"
        format_end = "DEFAULT_END"
        group_name = "default"
        pattern = ""
        skip = True

        super().__init__(
            format_start=format_start,
            format_end=format_end,
            group_name=group_name,
            pattern=pattern,
            skip=skip,
            verbose=verbose,
            **kwargs,
        )

        self._default: str = kwargs.get("default", "")

    def parse(self, output: str) -> Default:
        self.success = True
        return Default(default=self._default, success=self.success)

    @classmethod
    def construct_output(cls, **kwargs) -> str:
        """For testing purposes."""
        return ""

    @property
    def default(self) -> str:
        return self._default

    @default.setter
    def default(self, default: str):
        self._default = default
