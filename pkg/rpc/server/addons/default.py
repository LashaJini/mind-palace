import pkg.rpc.server.gen.Palace_pb2 as pbPalace
from pkg.rpc.server.addons.abstract import Addon
from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.vdb import Milvus


class DefaultAddon(Addon):
    def apply(
        self,
        id: str,
        input: str,
        llm: CustomLlamaCPP,
        client: Milvus,
        verbose=False,
        **kwargs,
    ):
        """input -> insert embeddings"""
        client.insert({"id": id, "input": input})

        return pbPalace.AddonResult(id="", data=None)
