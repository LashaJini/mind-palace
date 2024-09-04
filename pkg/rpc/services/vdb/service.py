import grpc
from pkg.rpc.services.vdb.vdb import Milvus, MilvusInsertData
import pkg.rpc.gen.SharedTypes_pb2 as pbShared
from pkg.rpc.loggers.vdb import log


class VDBService:
    def __init__(self, client: Milvus):
        self.client = client

    def Insert(self, request, context: grpc.ServicerContext):
        if context.is_active():
            self.client.insert(
                user=request.user,
                data=MilvusInsertData(ids=request.ids, inputs=request.inputs),
            )
            return pbShared.Empty()
        log.warning("context is not active. Skipping vdb 'Insert'")

    def Ping(self, request, context):
        if context.is_active():
            if not self.client.ping():
                raise ConnectionError
            return pbShared.Empty()
        log.warning("context is not active. Skipping vdb 'Ping'")

    def Drop(self, request, context: grpc.ServicerContext):
        if context.is_active():
            self.client.drop()
            return pbShared.Empty()
        log.warning("context is not active. Skipping vdb 'Drop'")
