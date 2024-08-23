import grpc
import inspect

from pkg.rpc import config
import pkg.rpc.gen.Log_pb2_grpc as logService
import pkg.rpc.gen.Log_pb2 as pbLog


class LogGrpcClient:
    def __init__(self, service_name: str):
        self.service_name = service_name

        channel = grpc.insecure_channel(f"localhost:{config.LOG_GRPC_SERVER_PORT}")
        self.client = logService.LogStub(channel)

    def _request(self, msg: str, type: str) -> pbLog.LogRequest:
        s = inspect.stack()
        # caller [2] -> log type (info, debug...) [1] -> _request [0]
        nth_call = 2
        filename = s[nth_call].filename
        line = s[nth_call].lineno

        return pbLog.LogRequest(
            msg=msg,
            filename=filename,
            line=line,
            service_name=self.service_name,
            type=type,
        )

    def info(self, msg: str):
        self.client.Message(request=self._request(msg=msg, type="info"))

    def debug(self, msg: str):
        self.client.Message(request=self._request(msg=msg, type="debug"))

    def warning(self, msg: str):
        self.client.Message(request=self._request(msg=msg, type="warning"))

    def exception(self, msg: str):
        self.client.Message(request=self._request(msg=msg, type="exception"))
