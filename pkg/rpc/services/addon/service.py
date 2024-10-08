from typing import List
import grpc
from llama_index.core import SimpleDirectoryReader

from pkg.rpc.services.addon.addons.factory import AddonFactory
from pkg.rpc.services.addon.addons.joined import JoinedAddons
from pkg.rpc.services.llm.llm import CustomLlamaCPP
import pkg.rpc.services.addon.prompts.joined as joined_prompts
from pkg.rpc.services.addon.prompts.abstract import Prompts
from pkg.rpc.services.addon.prompts.factory import PromptsFactory

import pkg.rpc.gen.Palace_pb2 as pbPalace
import pkg.rpc.gen.SharedTypes_pb2 as pbShared
from pkg.rpc.log.client import LogGrpcClient
from pkg.rpc.loggers.vdb import log


class AddonService:
    def __init__(self, llm: CustomLlamaCPP, log: LogGrpcClient, verbose=False):
        self.llm = llm
        self.verbose = verbose
        self.log = log

    def ApplyAddon(self, request: pbPalace.JoinedAddons, context: grpc.ServicerContext):
        if context.is_active():
            try:
                file = request.file
                documents = SimpleDirectoryReader(input_files=[file]).load_data()

                input = documents[0].text

                addons = request.addons
                if addons.joined:
                    joined_addons = JoinedAddons([name for name in addons.names])
                    prepared_input = joined_addons.prepare_input(user_input=input)

                    instructions = []
                    formats: List[str] = []
                    for name in addons.names:
                        prompt = PromptsFactory.construct(name)
                        joinable = prompt.joinable_template(
                            **joined_addons.prompt_variables
                        )

                        addon_instructions = joinable.instructions
                        addon_format = joinable.format

                        instructions.append(addon_instructions)
                        formats.append(addon_format)

                    result = (
                        prepared_input.apply(
                            llm=self.llm,
                            instructions=", ".join([s for s in instructions if s]),
                            format=".\n".join(["- " + s for s in formats if s]),
                            verbose=self.verbose,
                        )
                        .finalize(verbose=self.verbose)
                        .result()
                    )
                    self.log.debug(f"Result {result}")
                    return result

                # TODO: abstract away
                name = addons.names[0]
                result = (
                    AddonFactory.construct(name)
                    .prepare_input(user_input=input)
                    .apply(
                        llm=self.llm,
                        verbose=self.verbose,
                    )
                    .finalize(verbose=self.verbose)
                    .result()
                )
                self.log.debug(f"Result {result}")
                return result
            except Exception as e:
                self.log.exception(str(e))
        log.warning("context is not active. Skipping addon 'ApplyAddon'")

    def JoinAddons(self, request, context: grpc.ServicerContext):
        if context.is_active():
            try:
                file = request.file
                documents = SimpleDirectoryReader(input_files=[file]).load_data()
                input_text = documents[0].text

                input_text_token_count = self.llm.token_size(input_text)
                sys_prompt_token_count = self.llm.token_size(Prompts.sys_prompt)
                joined_prompt_token_count = self.llm.token_size(
                    joined_prompts.DEFAULT_JOINED_TMPL
                )
                token_decrements = [
                    input_text_token_count,
                    sys_prompt_token_count,
                    joined_prompt_token_count,
                ]

                addons_tokens = []
                for step in request.steps:
                    prompt = PromptsFactory.construct(name=step)
                    format_token_count = prompt.joinable_template_token_count(
                        llm=self.llm
                    )
                    addons_tokens.append({"name": step, "token": format_token_count})

                # TODO: algo for most addons joined together
                addons_tokens.sort(key=lambda x: x["token"])

                self.log.debug(f"addons + approx tokens required: {addons_tokens}")

                available_tokens = self.llm.calculate_available_tokens(token_decrements)
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
                    self.log.debug(
                        f"{'Joined addons' if addon_cluster.joined else 'Single addon'}: {', '.join(addon_cluster.names) if addon_cluster.joined else addon_cluster.names[0]}"
                    )
                    for addon_cluster in addons
                ]

                joinedAddons = pbPalace.JoinedAddonsResponse(addons=addons)

                return joinedAddons
            except Exception as e:
                self.log.exception(str(e))
        log.warning("context is not active. Skipping addon 'JoinAddons'")

    def Ping(self, request, context):
        return pbShared.Empty()
