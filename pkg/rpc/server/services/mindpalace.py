from typing import List

from llama_index.core import SimpleDirectoryReader

from pkg.rpc.server import logger
from pkg.rpc.server.addons.factory import AddonFactory
from pkg.rpc.server.addons.joined import JoinedAddons
from pkg.rpc.server.config import ServerConfig
from pkg.rpc.server.prompts.joined import JoinedPrompts
from pkg.rpc.server.prompts.abstract import Prompts
from pkg.rpc.server.prompts.factory import PromptsFactory
from pkg.rpc.server.vdb import Milvus, MilvusInsertData
from pkg.rpc.server.llm import CustomLlamaCPP

import pkg.rpc.server.gen.Palace_pb2 as pbPalace


class MindPalaceService:
    def __init__(
        self, llm: CustomLlamaCPP, client: Milvus, server_config: ServerConfig
    ):
        self.llm = llm
        self.server_config = server_config

        self.client = client

    def ApplyAddon(self, request: pbPalace.JoinedAddons, context):
        try:
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

                return (
                    JoinedAddons(names)
                    .prepare_input(user_input=input)
                    .apply(
                        llm=self.llm,
                        instructions=", ".join([s for s in instructions if s]),
                        format="\n".join([s for s in formats if s]),
                        verbose=self.server_config.verbose,
                    )
                    .finalize(verbose=self.server_config.verbose)
                    .result(verbose=self.server_config.verbose)
                )

            # TODO: abstract away
            name = addons.names[0]
            return (
                AddonFactory.construct(name)
                .prepare_input(user_input=input)
                .apply(
                    llm=self.llm,
                    verbose=self.server_config.verbose,
                    max_keywords=10,
                )
                .finalize(verbose=self.server_config.verbose)
                .result(verbose=self.server_config.verbose)
            )
        except Exception as e:
            logger.log.exception(e)

    def JoinAddons(self, request, context):
        try:
            file = request.file
            documents = SimpleDirectoryReader(input_files=[file]).load_data()
            input_text = documents[0].text

            input_text_token_count = self.llm.token_size(input_text)
            sys_prompt_token_count = Prompts.system_prompt_token_count(self.llm)  # ok
            joined_prompt_token_count = JoinedPrompts().standalone_template_token_count(
                self.llm
            )

            addons_tokens = []
            for step in request.steps:
                prompt = PromptsFactory.construct(name=step)
                format_token_count = prompt.joinable_template_token_count(llm=self.llm)
                addons_tokens.append({"name": step, "token": format_token_count})

            # TODO: algo for most addons joined together
            addons_tokens.sort(key=lambda x: x["token"])

            logger.log.debug(f"> addons + approx tokens required: {addons_tokens}")

            available_tokens = self.llm.calculate_available_tokens(
                input_text_token_count,
                sys_prompt_token_count,
                joined_prompt_token_count,
                self.server_config.verbose,
            )
            tokens_left = available_tokens

            addons: List[pbPalace.JoinedAddon] = []
            batch_addons = []

            def add_batch_to_addons():
                if len(batch_addons) > 0:
                    addons.append(
                        pbPalace.JoinedAddon(
                            names=batch_addons, joined=len(batch_addons) > 1
                        )
                    )

            for addon in addons_tokens:
                token_cost = addon.get("token")
                addon_name = addon.get("name")

                if token_cost < tokens_left:
                    tokens_left -= token_cost
                    batch_addons.append(addon_name)
                else:
                    add_batch_to_addons()

                    # reset for the next batch
                    batch_addons = [addon_name]
                    tokens_left = available_tokens - token_cost

            # drain remaining addons in batch
            add_batch_to_addons()

            [
                logger.log.debug(
                    f"> {'Joined addons' if addon_cluster.joined else 'Single addon'}: {', '.join(addon_cluster.names) if addon_cluster.joined else addon_cluster.names[0]}"
                )
                for addon_cluster in addons
            ]

            joinedAddons = pbPalace.JoinedAddonsResponse(addons=addons)

            return joinedAddons
        except Exception as e:
            logger.log.exception(e)

    def VDBInsert(self, request, context):
        self.client.insert(
            user=request.user,
            data=MilvusInsertData(ids=request.ids, inputs=request.inputs),
        )
        return pbPalace.Empty()

    def Ping(self, request, context):
        return pbPalace.Empty()

    def VDBPing(self, request, context):
        if not self.client.ping():
            raise ConnectionError
        return pbPalace.Empty()

    def VDBDrop(self, request, context):
        self.client.drop()
        return pbPalace.Empty()

    def SetConfig(self, request: pbPalace.Config, context):
        if request.map is not None:
            self.server_config.update(**request.map)

        return pbPalace.Empty()
