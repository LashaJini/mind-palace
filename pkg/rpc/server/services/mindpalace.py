from typing import List

from llama_index.core import SimpleDirectoryReader

from pkg.rpc.server.addons.factory import AddonFactory
from pkg.rpc.server.addons.joined import JoinedAddons
from pkg.rpc.server.prompts.joined import JoinedPrompts
from pkg.rpc.server.prompts.abstract import Prompts
from pkg.rpc.server.prompts.factory import PromptsFactory
from pkg.rpc.server.vdb import Milvus, MilvusInsertData
from pkg.rpc.server.llm import CustomLlamaCPP

import pkg.rpc.server.gen.Palace_pb2 as pbPalace

host = "localhost"
port = 19530
db_name = "user1_mind_palace"
collection_name = "llamatest"
client = Milvus(
    host=host,
    port=port,
    db_name=db_name,
    collection_name=collection_name,
)


class MindPalaceService:
    def __init__(self, llm: CustomLlamaCPP, verbose=False):
        self.llm = llm
        self.verbose = verbose

    def ApplyAddon(self, request: pbPalace.JoinedAddons, context):
        file = request.file
        documents = SimpleDirectoryReader(input_files=[file]).load_data()

        input = documents[0].text

        addons = request.addons
        if addons.joined:
            instructions = []
            formats = []
            names = []
            for name in addons.names:
                prompt = PromptsFactory.construct(name)
                joinable = prompt.joinable_template()

                addon_instructions = joinable.instructions
                addon_format = joinable.format

                instructions.append(addon_instructions)
                formats.append(addon_format)
                names.append(name)

            return JoinedAddons(names).apply(
                input=input,
                llm=self.llm,
                client=client,
                instructions=", ".join([s for s in instructions if s]),
                format="\n".join([s for s in formats if s]),
                verbose=self.verbose,
            )

        name = addons.names[0]
        addon = AddonFactory.construct(name)
        return addon.apply(
            input=input,
            llm=self.llm,
            client=client,
            verbose=self.verbose,
            max_keywords=10,
        )

    def JoinAddons(self, request, context):
        file = request.file
        documents = SimpleDirectoryReader(input_files=[file]).load_data()
        input_text = documents[0].text

        input_text_token_count = self.llm.token_size(input_text)
        sys_prompt_token_count = Prompts.system_prompt_token_count(self.llm)
        joined_prompt_token_count = JoinedPrompts().standalone_template_token_count(
            self.llm
        )

        addons_tokens = []
        for step in request.steps:
            prompt = PromptsFactory.construct(name=step)
            token_count = prompt.joinable_template_token_count(llm=self.llm)
            addons_tokens.append({"name": step, "token": token_count})

        addons_tokens.sort(key=lambda x: x["token"])

        available_tokens = self.llm.calculate_available_tokens(
            input_text_token_count,
            sys_prompt_token_count,
            joined_prompt_token_count,
            self.verbose,
        )
        addons: List[pbPalace.JoinedAddon] = []
        batch_addons = []
        for addon in addons_tokens:
            if addon.get("token") < available_tokens:
                available_tokens -= addon.get("token")
                batch_addons.append(addon.get("name"))
                continue

            available_tokens = self.llm.calculate_available_tokens(
                input_text_token_count,
                sys_prompt_token_count,
                joined_prompt_token_count,
                self.verbose,
            )
            addons.append(
                pbPalace.JoinedAddon(names=batch_addons, joined=len(batch_addons) > 1)
            )
            batch_addons = []

            available_tokens -= addon.get("token")
            batch_addons.append(addon.get("name"))

        if len(batch_addons) > 0:
            addons.append(
                pbPalace.JoinedAddon(names=batch_addons, joined=len(batch_addons) > 1)
            )

        joinedAddons = pbPalace.JoinedAddonsResponse(addons=addons)

        return joinedAddons

    def VDBInsert(self, request, context):
        client.insert(MilvusInsertData(id=request.id, input=request.input))
        return pbPalace.Empty()
