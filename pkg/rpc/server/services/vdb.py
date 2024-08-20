from pkg.rpc.server.vdb import Milvus, MilvusInsertData
import pkg.rpc.server.gen.SharedTypes_pb2 as pbShared


class VDBService:
    def __init__(self, client: Milvus):
        self.client = client

    def Insert(self, request, context):
        self.client.insert(
            user=request.user,
            data=MilvusInsertData(ids=request.ids, inputs=request.inputs),
        )
        return pbShared.Empty()

    def Ping(self, request, context):
        if not self.client.ping():
            raise ConnectionError
        return pbShared.Empty()

    def Drop(self, request, context):
        self.client.drop()
        return pbShared.Empty()
