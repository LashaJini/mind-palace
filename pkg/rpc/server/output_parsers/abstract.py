from abc import ABC, abstractmethod
from typing import Any, ClassVar
from llama_index.core.output_parsers.base import ChainableOutputParser
from pydantic import BaseModel


class CustomBaseModel(BaseModel):
    name: ClassVar[str] = ""
    value: Any

    def __init__(self, /, **data) -> None:
        super().__init__(**data)


class OutputParser(ChainableOutputParser, ABC):
    def __init__(
        self,
        format_start: str,
        format_end: str,
        group_name: str,
        pattern: str,
        verbose: bool = False,
    ):
        self.verbose = verbose
        self.format_start = format_start
        self.format_end = format_end
        self.group_name = group_name
        self.pattern = pattern

    @abstractmethod
    def parse(self, output: str) -> CustomBaseModel:
        raise NotImplementedError
