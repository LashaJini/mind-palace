from abc import ABC, abstractmethod
from typing import Any, ClassVar, List
from llama_index.core.output_parsers.base import ChainableOutputParser
from pydantic import BaseModel


class CustomBaseModel(BaseModel, ABC):
    name: ClassVar[str] = ""
    value: Any

    def __init__(self, /, **data) -> None:
        super().__init__(**data)

    @abstractmethod
    def get_value(self) -> Any:
        pass


class OutputParser(ChainableOutputParser, ABC):
    def __init__(
        self,
        format_start: str,
        format_end: str,
        group_name: str,
        pattern: str,
        skip: bool = False,
        verbose: bool = False,
    ):
        self.success = False
        self.skip = skip
        self.verbose = verbose
        self.format_start = format_start
        self.format_end = format_end
        self.group_name = group_name
        self.pattern = pattern

    @abstractmethod
    def parse(self, output: str) -> CustomBaseModel:
        raise NotImplementedError

    @classmethod
    @abstractmethod
    def construct_output(cls, **kwargs) -> str:
        """For testing purposes."""
        raise NotImplementedError
