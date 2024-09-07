import grpc
from pkg.rpc.services.vdb.vdb import VDB, VDBRows
import pkg.rpc.gen.SharedTypes_pb2 as pbShared
import pkg.rpc.gen.VDB_pb2 as pbVDB
from pkg.rpc.loggers.vdb import log


class VDBService:
    def __init__(self, client: VDB):
        self.client = client

    def Insert(self, request: pbVDB.InsertRequest, context: grpc.ServicerContext):
        if context.is_active():
            success = self.client.insert(
                user=request.user,
                data=VDBRows(list(request.rows)),
            )

            if not success:
                context.abort(
                    code=grpc.StatusCode.INVALID_ARGUMENT,
                    details="could not insert into vdb",
                )

            return pbShared.Empty()
        log.warning("context is not active. Skipping vdb 'Insert'")

    def Search(self, request: pbVDB.SearchRequest, context: grpc.ServicerContext):
        if context.is_active():
            result = self.client.search(
                user=request.user,
                text=request.text,
            )
            return pbVDB.SearchResponse(
                rows=[
                    pbVDB.SearchResponse.VDBRow(
                        id=x.id, distance=x.distance, type=x.type
                    )
                    for x in result
                ]
            )
        log.warning("context is not active. Skipping vdb 'Search'")

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
