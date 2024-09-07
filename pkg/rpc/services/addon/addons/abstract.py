from abc import ABC, abstractmethod
from typing import Any, Optional, Self

import pkg.rpc.gen.Palace_pb2 as pbPalace

from pkg.rpc.services.llm.llm import CustomLlamaCPP
from pkg.rpc.services.addon.output_parsers.abstract import CustomBaseModel, OutputParser


class Addon(ABC):
    _parser: Any
    _output_model: Any
    _result: pbPalace.AddonResult | None
    _prompt_variables: dict

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

    def result(self) -> pbPalace.AddonResult | None:
        return self._result

    @property
    @abstractmethod
    def parser(self) -> OutputParser:
        raise NotImplementedError

    @property
    @abstractmethod
    def output_model(self) -> CustomBaseModel:
        raise NotImplementedError
