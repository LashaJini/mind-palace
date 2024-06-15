from llama_index.core import PromptTemplate
from llama_index.core.program import LLMTextCompletionProgram

import pkg.rpc.server.gen.Palace_pb2 as pbPalace
from pkg.rpc.server.addons.abstract import Addon
from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.prompts.summary import SummaryPrompts
from pkg.rpc.server.output_parsers.summary import Summary, SummaryParser


class SummaryAddon(Addon):
    def apply(
        self,
        input: str,
        llm: CustomLlamaCPP,
        verbose=False,
        **kwargs,
    ):
        """input -> generate summary -> insert embeddings -> return summary"""
        prompt = SummaryPrompts().prompt(text=input, verbose=verbose, **kwargs)
        parser = SummaryParser(verbose=verbose)
        program = LLMTextCompletionProgram(
            llm=llm,
            output_parser=parser,
            output_cls=Summary,  # type:ignore
            prompt=PromptTemplate(prompt),
            verbose=verbose,
        )

        value = program(context_str=input, verbose=verbose).dict().get("value")

        return pbPalace.AddonResult(
            data={
                Summary.name: pbPalace.AddonResultInfo(
                    success=parser.success, value=value
                )
            },
        )
