from typing import List
from llama_index.core import PromptTemplate
from llama_index.core.program import LLMTextCompletionProgram

import pkg.rpc.server.gen.Palace_pb2 as pbPalace
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
        prompt = JoinedPrompts().prompt(text=input, verbose=verbose, **kwargs)
        parser = JoinedParser(verbose=verbose, addons=self.addons)
        program = LLMTextCompletionProgram(
            llm=llm,
            output_parser=parser,
            output_cls=Joined,  # type:ignore
            prompt=PromptTemplate(prompt),
            verbose=verbose,
        )

        results = (
            program(context_str=input, verbose=verbose, **kwargs).dict().get("value")
        )

        data = {}
        if results is not None:
            for key, addon_result_info in results.items():
                data[key] = pbPalace.AddonResultInfo(
                    value=addon_result_info.get("value"),
                    success=addon_result_info.get("success"),
                )

        return pbPalace.AddonResult(
            id="",
            data=data,
        )
