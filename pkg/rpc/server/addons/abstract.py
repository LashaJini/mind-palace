from abc import ABC, abstractmethod

import pkg.rpc.server.gen.Palace_pb2 as pbPalace

from pkg.rpc.server.llm import CustomLlamaCPP


class Addon(ABC):
    @abstractmethod
    def apply(
        self,
        input: str,
        llm: CustomLlamaCPP,
        verbose=False,
        **kwargs,
    ) -> pbPalace.AddonResult:
        raise NotImplementedError
