import gen.Palace_pb2 as pbPalace
import gen.Palace_pb2_grpc as grpcPalace
from sentence_transformers import SentenceTransformer
from pkg.addons.example import Apply as exampleApply
import grpc
import time
import os
from dotenv import load_dotenv
from concurrent import futures

load_dotenv()

PYTHON_GRPC_SERVER_PORT = os.getenv("PYTHON_GRPC_SERVER_PORT", 50051)
_ONE_DAY_IN_SECONDS = 60 * 60 * 24

transformer = SentenceTransformer("all-MiniLM-L6-v2")


def convert_data(data):
    converted_data = []
    for item in data:
        new_item = {
            "name": item["name"],
            "vector": transformer.encode(item["input"]),
        }
        converted_data.append(new_item)
    return converted_data


addonDict = {"mind-palace-resource-summary": exampleApply}


class MindPalaceService:
    def Add(self, request, context):
        file = request.file
        with open(file) as f:
            input = f.read()

        for step in request.steps:
            if step in addonDict:
                print(addonDict[step](input))

        data = [{"name": "original", "input": input}]
        ins = convert_data(data)

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
