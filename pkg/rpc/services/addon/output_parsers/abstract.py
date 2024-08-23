from abc import ABC, abstractmethod
from typing import ClassVar
from llama_index.core.output_parsers.base import ChainableOutputParser
from pydantic import BaseModel

import pkg.rpc.gen.Palace_pb2 as pbPalace


class CustomBaseModel(BaseModel, ABC):
    name: ClassVar[str] = ""
    success: bool = False

    def __init__(self, **kwargs):
        super().__init__(**kwargs)

    @abstractmethod
    def to_addon_result(self) -> pbPalace.AddonResult:
        raise NotImplementedError

    def __iter__(self):
        return super().__iter__()


class OutputParser(ChainableOutputParser, ABC):
    def __init__(
        self,
        format_start: str = "",
        format_end: str = "",
        group_name: str = "",
        pattern: str = "",
        skip: bool = False,
        verbose: bool = False,
        **kwargs,
    ):
        self.success = False
        self.skip = skip
        self.verbose = verbose
        self.format_start = format_start
        self.format_end = format_end
        self.group_name = group_name
        self.pattern = pattern
        self.kwargs = kwargs

    @abstractmethod
    def parse(self, output: str) -> CustomBaseModel:
        raise NotImplementedError

    @classmethod
    @abstractmethod
    def construct_output(cls, **kwargs) -> str:
        """For testing purposes."""
        raise NotImplementedError
