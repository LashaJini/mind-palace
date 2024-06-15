import pkg.rpc.server.gen.Palace_pb2 as pbPalace
from pkg.rpc.server.addons.abstract import Addon
from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.output_parsers.default import Default


class DefaultAddon(Addon):
    def apply(
        self,
        input: str,
        llm: CustomLlamaCPP,
        verbose=False,
        **kwargs,
    ):
        """default, identity addon"""
        return pbPalace.AddonResult(
            data={Default.name: pbPalace.AddonResultInfo(success=True, value=[input])}
        )
