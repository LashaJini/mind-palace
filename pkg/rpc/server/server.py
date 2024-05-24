from llama_index.core import SimpleDirectoryReader
import grpc
import time
import os
from dotenv import load_dotenv
from concurrent import futures

from pkg.rpc.server.addons import AddonsDict
from pkg.rpc.server.vdb import Milvus

import gen.Palace_pb2 as pbPalace
import gen.Palace_pb2_grpc as grpcPalace

load_dotenv()

PYTHON_GRPC_SERVER_PORT = os.getenv("PYTHON_GRPC_SERVER_PORT", 50051)
_ONE_DAY_IN_SECONDS = 60 * 60 * 24

prompt = "The tower is 324 metres (1,063 ft) tall, about the same height as an 81-storey building, and the tallest structure in Paris. Its base is square, measuring 125 metres (410 ft) on each side. During its construction, the Eiffel Tower surpassed the Washington Monument to become the tallest man-made structure in the world, a title it held for 41 years until the Chrysler Building in New York City was finished in 1930. It was the first structure to reach a height of 300 metres. Due to the addition of a broadcasting aerial at the top of the tower in 1957, it is now taller than the Chrysler Building by 5.2 metres (17 ft). Excluding transmitters, the Eiffel Tower is the second tallest free-standing structure in France after the Millau Viaduct."

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
    def ApplyAddon(self, request, context):
        file = request.file
        documents = SimpleDirectoryReader(input_files=[file]).load_data()

        # TODO: AddMany
        if len(documents) > 0:
            pass

        document = documents[0]
        input = document.text

        result = None
        try:
            if request.step in AddonsDict.keys():
                addonApply = AddonsDict[request.step]
                result = addonApply(original_id=request.id, input=input, client=client)
        except Exception as e:
            print(e)

        return result


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
