from typing import List
import grpc
import time
import os
from dotenv import load_dotenv
from concurrent import futures

from llama_index.core import SimpleDirectoryReader

from pkg.rpc.server.addons.factory import AddonFactory
from pkg.rpc.server.addons.joined import JoinedAddons
from pkg.rpc.server.prompts.joined import JoinedPrompts
from pkg.rpc.server.prompts.abstract import Prompts
from pkg.rpc.server.prompts.factory import PromptsFactory
from pkg.rpc.server.vdb import Milvus, MilvusInsertData
from pkg.rpc.server.llm import CustomLlamaCPP

import pkg.rpc.server.gen.Palace_pb2 as pbPalace
import pkg.rpc.server.gen.Palace_pb2_grpc as grpcPalace

load_dotenv()

PYTHON_GRPC_SERVER_PORT = os.getenv("PYTHON_GRPC_SERVER_PORT", 50051)
_ONE_DAY_IN_SECONDS = 60 * 60 * 24

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


verbose = True

llm = CustomLlamaCPP(
    verbose=verbose,
    generate_kwargs={
        "top_k": 1,  # TODO: config
        "stop": ["<|endoftext|>", "</s>"],  # TODO: wtf
        # "seed": 4294967295,
        # "seed": -1,
    },
    # kwargs to pass to __init__()
    model_kwargs={
        "n_gpu_layers": -1,  # TODO: config
    },
)


class MindPalaceService:
    def ApplyAddon(self, request: pbPalace.JoinedAddons, context):
        file = request.file
        documents = SimpleDirectoryReader(input_files=[file]).load_data()

        # TODO: AddMany
        if len(documents) > 0:
            pass

        document = documents[0]
        input = document.text

        result = None
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

            result = JoinedAddons(names).apply(
                input=input,
                llm=llm,
                client=client,
                instructions=", ".join([s for s in instructions if s]),
                format="\n".join([s for s in formats if s]),
                verbose=verbose,
            )
        else:
            name = addons.names[0]
            addon = AddonFactory.construct(name)
            result = addon.apply(
                input=input,
                llm=llm,
                client=client,
                verbose=verbose,
                max_keywords=10,
            )

        return result

    def JoinAddons(self, request, context):
        file = request.file
        documents = SimpleDirectoryReader(input_files=[file]).load_data()
        input_text = documents[0].text

        input_text_token_count = llm.token_size(input_text)
        sys_prompt_token_count = Prompts.system_prompt_token_count(llm)
        joined_prompt_token_count = JoinedPrompts().standalone_template_token_count(llm)

        addons_tokens = []
        for step in request.steps:
            prompt = PromptsFactory.construct(name=step)
            token_count = prompt.joinable_template_token_count(llm=llm)
            addons_tokens.append({"name": step, "token": token_count})

        addons_tokens.sort(key=lambda x: x["token"])

        available_tokens = llm.calculate_available_tokens(
            input_text_token_count,
            sys_prompt_token_count,
            joined_prompt_token_count,
            verbose,
        )
        addons: List[pbPalace.JoinedAddon] = []
        batch_addons = []
        for addon in addons_tokens:
            if addon.get("token") < available_tokens:
                available_tokens -= addon.get("token")
                batch_addons.append(addon.get("name"))
                continue

            available_tokens = llm.calculate_available_tokens(
                input_text_token_count,
                sys_prompt_token_count,
                joined_prompt_token_count,
                verbose,
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


def server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    grpcPalace.add_PalaceServicer_to_server(MindPalaceService(), server)
    server.add_insecure_port(f"[::]:{PYTHON_GRPC_SERVER_PORT}")
    server.start()
    print("Server started. Listen on port:", PYTHON_GRPC_SERVER_PORT)
    # server.wait_for_termination()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


def main():
    server()


if __name__ == "__main__":
    main()
