from abc import ABC, abstractmethod

import pkg.rpc.server.gen.Palace_pb2 as pbPalace

from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.vdb import Milvus


class Addon(ABC):
    @abstractmethod
    def apply(
        self,
        id: str,
        input: str,
        llm: CustomLlamaCPP,
        client: Milvus,
        verbose=False,
        **kwargs,
    ) -> pbPalace.AddonResult:
        raise NotImplementedError
