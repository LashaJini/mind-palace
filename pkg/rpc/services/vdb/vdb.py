import grpc
from typing import List
from pymilvus import MilvusClient, connections, DataType, db

from pkg.rpc import config
from pkg.rpc.services.llm.llm import EmbeddingModel
from pkg.rpc.log.client import LogGrpcClient
import pkg.rpc.gen.Palace_pb2_grpc as palaceService
import pkg.rpc.gen.Palace_pb2 as pbPalace
import pkg.rpc.gen.VDB_pb2 as pbVDB


class VDBSearchResult:
    id: str
    distance: float
    type: str

    def __init__(self, id: str, distance: float, type: str):
        self.id = id
        self.distance = distance
        self.type = type

    def __repr__(self):
        return f"VDBSearchResult({self.id}, {self.distance})"


class VDBRows:
    def __init__(self, rows: List[pbVDB.InsertRequest.VDBRow]):
        self.rows = [VDBRows.VDBRow(row.id, row.input, row.type) for row in rows]

    class VDBRow:
        def __init__(self, id: str, input: str, type: str):
            self.id = id
            self.input = input
            self.type = type

        def __repr__(self):
            return f"VDBRow({self.id}, {self.input}, {self.type})"

    def __repr__(self):
        return f"VDBRows({self.rows})"


class VDB:
    db_name: str = ""
    ROW_TYPES = ["whole", "chunk"]

    def __init__(self, host: str, port: int, log: LogGrpcClient):
        self.host = host
        self.port = port
        self.db_name = config.VDB_NAME
        self.search_limit = config.VDB_SEARCH_LIMIT

        self.log = log

        channel = grpc.insecure_channel(f"localhost:{config.PALACE_GRPC_SERVER_PORT}")
        self.embedding_model = palaceService.EmbeddingModelStub(channel)

        connections.connect(host=self.host, port=self.port)
        if not self.db_exists():
            self.log.info(f"Database '{config.VDB_NAME}' not found. Creating...")
            db.create_database(config.VDB_NAME)

        self.client = MilvusClient(
            uri=f"http://{self.host}:{self.port}",
            db_name=self.db_name,
        )

    def create_collection(self, collection_name: str) -> None:
        schema = self.client.create_schema(
            auto_id=False,
            enable_dynamic_field=True,
        )

        schema.add_field(
            field_name="id",
            datatype=DataType.VARCHAR,
            is_primary=True,
            max_length=64,
        )
        schema.add_field(
            field_name="vector",
            datatype=DataType.FLOAT_VECTOR,
            dim=EmbeddingModel.dimension,
        )
        schema.add_field(
            field_name="type",
            description=f"input type either {'|'.join(VDB.ROW_TYPES)}",
            datatype=DataType.VARCHAR,
            max_length=5,
        )

        index_params = self.client.prepare_index_params()
        index_params.add_index(field_name="id", index_type="INVERTED")
        index_params.add_index(
            field_name="vector",
            index_type="IVF_FLAT",
            metric_type=EmbeddingModel.metric_type,
            params={"nlist": 1536},
        )

        self.log.info(f"Collection {collection_name} not found. Creating...")
        self.client.create_collection(
            collection_name=collection_name,
            schema=schema,
            index_params=index_params,
        )

    def collection(self, user: str) -> str:
        collection_name = f"{user}_collection"
        if not self.client.has_collection(collection_name):
            self.create_collection(collection_name)

        return collection_name

    def insert(self, user: str, data: VDBRows) -> bool:
        for row in data.rows:
            if row.type not in VDB.ROW_TYPES:
                return False

        embedded_data = self._get_text_embedding(data)
        collection_name = self.collection(user)
        self.log.info(f"Insert into {collection_name}: {data}")
        self.client.insert(collection_name=collection_name, data=embedded_data)

        return True

    def search(self, user: str, text: str) -> List[VDBSearchResult]:
        embeddings = self.embedding_model.CalculateEmbeddings(pbPalace.Text(text=text))

        rows = self.client.search(
            collection_name=self.collection(user),
            data=[list(embeddings.embedding)],
            limit=self.search_limit,
            output_fields=["id", "distance", "type"],
        )

        self.log.debug(f"VDB search result ({len(rows[0])}): {rows}")

        result = []
        for row in rows[0]:
            entity = row.get("entity", {})
            type: str = entity.get("type", VDB.ROW_TYPES[0])

            result.append(
                VDBSearchResult(
                    row.get("id", -1),
                    row.get("distance", -1),
                    type,
                )
            )

        return result

    def _get_text_embedding(self, data: VDBRows):
        embedded_data = []

        for row in data.rows:
            embeddings: pbPalace.Embeddings = self.embedding_model.CalculateEmbeddings(
                pbPalace.Text(text=row.input)
            )

            new_item = {
                "id": row.id,
                "vector": list(embeddings.embedding),
                "type": row.type,
            }
            embedded_data.append(new_item)
        return embedded_data

    def drop(self):
        for collection in self.client.list_collections():
            self.client.drop_collection(collection)
            self.log.info(f"Drop {collection} collection")

        db.drop_database(self.db_name)
        self.log.info(f"Drop {self.db_name} database")

    def db_exists(self):
        return self.db_name in db.list_database()

    def ping(self):
        return self.db_exists()
