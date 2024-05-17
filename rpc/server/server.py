import gen.Palace_pb2 as pbPalace
import gen.Palace_pb2_grpc as grpcPalace
import grpc
import time
from concurrent import futures

PYTHON_GRPC_PORT = 50052
_ONE_DAY_IN_SECONDS = 60 * 60 * 24


class AddService:
    def Add(self, request, context):
        print(pbPalace.Memory)
        return pbPalace.Status(code=1)


def server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    grpcPalace.add_PalaceServicer_to_server(AddService(), server)
    server.add_insecure_port(f"[::]:{PYTHON_GRPC_PORT}")
    server.start()
    print("Server started. Listen on port:", PYTHON_GRPC_PORT)
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
