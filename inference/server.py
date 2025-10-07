from __future__ import print_function

import grpc
from grpc import aio
# import message_pb2
from prediction_pb2 import (
    PredictionRequest, PredictionResponse, EmbedRequest, EmbedResponse
)
from prediction_pb2_grpc import PredictionServicer, add_PredictionServicer_to_server
from label import InferencePool, predict, embed, create_inference_pool


class Server(PredictionServicer):
    def __init__(self, inference_pool):
        super().__init__()
        self.inference_pool = inference_pool

    async def Predict(self, request: PredictionRequest, context):
        print(f"[PREDICT] Request received: {request.message}", flush=True)
        prediction = await predict(request.message, self.inference_pool)
        print(f"[PREDICT] Prediction result: {prediction}", flush=True)
        context.set_code(grpc.StatusCode.OK)
        context.set_details("Prediction successful")
        if request.return_embedding:
            embedding = await embed(request.message, self.inference_pool)
            print(f"[PREDICT] Embedding length: {len(embedding)}", flush=True)
        else:
            embedding = []
        return PredictionResponse(label=prediction, embedding=embedding)

    async def Embed(self, request: EmbedRequest, context):
        print(f"[EMBED] Request received: {request.message}", flush=True)
        embedding = await embed(request.message, self.inference_pool)
        print(f"[EMBED] Embedding length: {len(embedding)}", flush=True)
        context.set_code(grpc.StatusCode.OK)
        context.set_details("Tokenization successful")
        return EmbedResponse(embedding=embedding)


async def serve(pool_size: int = 3, wait_for_termination: bool = True) -> None | tuple[aio.Server, InferencePool]:
    print("TODO: Add a relashionship inference, that predicts if two messages are related as contradiction, entailment or neutral", flush=True)
    inference_pool = create_inference_pool(
        model_path="models/prediction/model.onnx",
        tokenizer_path="models/prediction/tokenizer.json",
        embedder_model_path="models/embedding/model.onnx",
        embedder_tokenizer_path="models/embedding/tokenizer.json",
        pool_size=pool_size,
    )
    print(f"Created inference pool with pool_size={pool_size}...", flush=True)
    port = 50051
    server = aio.server()
    add_PredictionServicer_to_server(Server(inference_pool), server)
    server.add_insecure_port(f"[::]:{port}")
    await server.start()
    print(f"Server started, listening on {port}", flush=True)

    if wait_for_termination:
        print("Waiting for termination...", flush=True)
        await server.wait_for_termination()
        print("Server terminated", flush=True)
        inference_pool.close()
        print("Inference pool terminated", flush=True)
        print("Server stopped", flush=True)
        return
    else:
        # Return the server so it can be managed externally
        print("Returning server and inference pool...", flush=True)
        return server, inference_pool
