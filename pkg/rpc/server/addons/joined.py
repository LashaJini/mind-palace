from typing import List
from llama_index.core.program import LLMTextCompletionProgram

import gen.Palace_pb2 as pbPalace

from pkg.rpc.server.addons.abstract import Addon
from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.vdb import Milvus
from pkg.rpc.server.prompts.joined import JoinedPrompts
from pkg.rpc.server.output_parsers.joined import Joined, JoinedParser


class JoinedAddons(Addon):
    def __init__(self, addons: List[str]):
        self.addons = addons

    def apply(
        self,
        id: str,
        input: str,
        llm: CustomLlamaCPP,
        client: Milvus,
        verbose=False,
        **kwargs,
    ):
        program = LLMTextCompletionProgram.from_defaults(
            llm=llm,
            output_parser=JoinedParser(verbose=verbose, addons=self.addons),
            output_cls=Joined,  # type:ignore
            prompt_template_str=JoinedPrompts().standalone_template(),
            verbose=verbose,
        )
        results = (
            program(context_str=input, verbose=verbose, **kwargs).dict().get("value")
        )

        data = {}
        if results is not None:
            for key, value in results.items():
                data[key] = pbPalace.Strings(value=value)

        return pbPalace.AddonResult(
            id="",
            data=data,
        )
