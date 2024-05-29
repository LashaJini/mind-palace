import uuid
from llama_index.core.program import LLMTextCompletionProgram

import gen.Palace_pb2 as pbPalace

from pkg.rpc.server.addons.abstract import Addon
from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.vdb import Milvus
from pkg.rpc.server.prompts.summary import SummaryPrompts
from pkg.rpc.server.output_parsers.summary import Summary, SummaryParser


class SummaryAddon(Addon):
    def apply(
        self,
        id: str,
        input: str,
        llm: CustomLlamaCPP,
        client: Milvus,
        verbose=False,
        **kwargs,
    ):
        """input -> generate summary -> insert embeddings -> return summary"""
        program = LLMTextCompletionProgram.from_defaults(
            llm=llm,
            output_parser=SummaryParser(verbose=verbose),
            output_cls=Summary,  # type:ignore
            prompt_template_str=SummaryPrompts().standalone_template(verbose=verbose),
            verbose=verbose,
        )

        result = program(context_str=input).dict().get("value")

        id = str(uuid.uuid4())
        client.insert({"id": id, "input": input})

        return pbPalace.AddonResult(
            id=id, data={"output": pbPalace.Strings(value=result)}
        )
