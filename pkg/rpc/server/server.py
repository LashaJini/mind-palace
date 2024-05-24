from typing import List
from llama_index.core import PromptTemplate, SimpleDirectoryReader
import grpc
import time
import os
from dotenv import load_dotenv
from concurrent import futures

from llama_index.core.program import LLMTextCompletionProgram
from llama_index.llms.llama_cpp.llama_utils import DEFAULT_SYSTEM_PROMPT
from llama_index.core.prompts.default_prompts import (
    DEFAULT_SUMMARY_PROMPT_TMPL,
    DEFAULT_KEYWORD_EXTRACT_TEMPLATE_TMPL,
)

from pkg.rpc.server.vdb import InputDataDict, Milvus
from pkg.rpc.server.llm import Summary, Album, llm

import gen.Palace_pb2 as pbPalace
import gen.Palace_pb2_grpc as grpcPalace

load_dotenv()

PYTHON_GRPC_SERVER_PORT = os.getenv("PYTHON_GRPC_SERVER_PORT", 50051)
_ONE_DAY_IN_SECONDS = 60 * 60 * 24

prompt = "The tower is 324 metres (1,063 ft) tall, about the same height as an 81-storey building, and the tallest structure in Paris. Its base is square, measuring 125 metres (410 ft) on each side. During its construction, the Eiffel Tower surpassed the Washington Monument to become the tallest man-made structure in the world, a title it held for 41 years until the Chrysler Building in New York City was finished in 1930. It was the first structure to reach a height of 300 metres. Due to the addition of a broadcasting aerial at the top of the tower in 1957, it is now taller than the Chrysler Building by 5.2 metres (17 ft). Excluding transmitters, the Eiffel Tower is the second tallest free-standing structure in France after the Millau Viaduct."
addonDict = {
    "mind-palace-resource-summary": llm.gen_summary,
    "mind-palace-resource-keywords": llm.gen_keywords,
}

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

# documents = SimpleDirectoryReader("/home/jini/examples/").load_data()
# data: List[InputDataDict] = [
#     {"id": "abcdefg", "input": documents[0].text},
#     {"id": "abcdefgh", "input": documents[1].text},
# ]
# client.insert(data)

# grpc input file
# input -> simple document
# gen summary
# gen tags
# gen embeddings for original document + summary -> insert embeddings
# summary, tags -> grpc


class MindPalaceService:
    def Add(self, request, context):
        file = request.file
        document = SimpleDirectoryReader(input_files=[file]).load_data()
        original = document[0].text

        summary = llm.gen_summary(original)
        keywords = llm.gen_keywords(original)
        print(summary)
        print(keywords)

        # for step in request.steps:
        #     if step in addonDict:
        #         print(addonDict[step](input))

        data = [{"name": "original", "input": input}]
        ins = data

        return pbPalace.Vectors(vectors=ins)


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
