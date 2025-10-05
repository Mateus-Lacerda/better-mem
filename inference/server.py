from __future__ import print_function

import grpc
from grpc import aio
# import message_pb2
from prediction_pb2 import (
    PredictionRequest, PredictionResponse, EmbedRequest, EmbedResponse
)
from prediction_pb2_grpc import PredictionServicer, add_PredictionServicer_to_server
from label import predict, embed, create_inference_pool


class Server(PredictionServicer):
    def __init__(self, inference_pool):
        super().__init__()
        self.inference_pool = inference_pool

    async def Predict(self, request: PredictionRequest, context):
        print("[PREDICT] Request received")
        prediction = await predict(request.message, self.inference_pool)
        context.set_code(grpc.StatusCode.OK)
        context.set_details("Prediction successful")
        if request.return_embedding:
            embedding = await embed(request.message, self.inference_pool)
        else:
            embedding = []
        return PredictionResponse(label=prediction, embedding=embedding)

    async def Embed(self, request: EmbedRequest, context):
        print("[EMBED] Request received")
        embedding = await embed(request.message, self.inference_pool)
        context.set_code(grpc.StatusCode.OK)
        context.set_details("Tokenization successful")
        return EmbedResponse(embedding=embedding)


async def serve():
    print("TODO: Add a relashionship inference, that predicts if two messages are related as contradiction, entailment or neutral")
    inference_pool = create_inference_pool(
        model_path="models/prediction/model.onnx",
        tokenizer_path="models/prediction/tokenizer.json",
        embedder_model_path="models/embedding/model.onnx",
        embedder_tokenizer_path="models/embedding/tokenizer.json",
        pool_size=100
    )
    port = 50051
    server = aio.server()
    add_PredictionServicer_to_server(Server(inference_pool), server)
    server.add_insecure_port(f"[::]:{port}")
    await server.start()
    print(f"Server started, listening on {port}")
    await server.wait_for_termination()
