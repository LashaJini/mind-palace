from abc import ABC, abstractmethod
from typing import Any, Optional, Self

from pkg.rpc.server import logger
import pkg.rpc.server.gen.Palace_pb2 as pbPalace

from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.output_parsers.abstract import CustomBaseModel, OutputParser


class Addon(ABC):
    _parser: Any
    _output_model: Any
    _result: pbPalace.AddonResult | None

    def __init__(self, **kwargs):
        self._result = None

    @abstractmethod
    def input(self, verbose=False) -> str:
        raise NotImplementedError

    @abstractmethod
    def apply(
        self,
        llm: CustomLlamaCPP,
        verbose=False,
        **kwargs,
    ) -> Self:
        raise NotImplementedError

    @abstractmethod
    def prepare_input(
        self,
        user_input: str,
    ) -> Self:
        raise NotImplementedError

    @abstractmethod
    def finalize(
        self, result: Optional[pbPalace.AddonResult] = None, verbose=False
    ) -> Self:
        raise NotImplementedError

    def result(self, verbose=False) -> pbPalace.AddonResult | None:
        logger.log.debug(f"> RESULT {self._result}")
        return self._result

    @property
    @abstractmethod
    def parser(self) -> OutputParser:
        raise NotImplementedError

    @property
    @abstractmethod
    def output_model(self) -> CustomBaseModel:
        raise NotImplementedError
