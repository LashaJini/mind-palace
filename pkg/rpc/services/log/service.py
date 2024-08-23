from pkg.rpc.services.log.log import Logger
import pkg.rpc.gen.Log_pb2 as pbLog
import pkg.rpc.gen.SharedTypes_pb2 as pbShared


class LogService:
    def __init__(self, log: Logger):
        self.log = log

        self.callers = {
            "info": self.log.info,
            "db_info": self.log.db_info,
            "tx_info": self.log.tx_info,
            "debug": self.log.debug,
            "warning": self.log.warning,
            "exception": self.log.exception,
            "error": self.log.error,
        }
        self.types = self.callers.keys()

    def Message(self, request: pbLog.LogRequest, context):
        extra = {
            "caller_filename": request.filename,
            "line": request.line,
            "service_name": request.service_name,
            "id": request.id,
        }

        if request.type not in self.types:
            raise ValueError(f"Unknown type {request.type}")

        self.callers[request.type](msg=request.msg, extra=extra)
        return pbShared.Empty()
