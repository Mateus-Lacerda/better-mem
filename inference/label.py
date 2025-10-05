import logging
from typing import AsyncGenerator, Sequence
from functools import partial

import asyncio
import onnxruntime as ort
from tokenizers import Tokenizer
from numpy import ndarray, argmax, mean

from multiprocessing import Pool
from multiprocessing.pool import Pool as PoolType

_session = None
_tokenizer = None
_embedder_tokenizer = None
_embedder_model = None

logger = logging.getLogger(__name__)


def start_session(model_path: str) -> None:
    global _session
    _session = ort.InferenceSession(
        model_path,
        sess_options=ort.SessionOptions(),
        providers=["CPUExecutionProvider"],
    )


def start_tokenizer(tokenizer_path: str) -> None:
    global _tokenizer
    _tokenizer = Tokenizer.from_file(tokenizer_path)


def start_embedder(
    embedder_model_path: str,
    embedder_tokenizer_path: str,
) -> None:
    global _embedder_model
    global _embedder_tokenizer
    _embedder_model = ort.InferenceSession(
        embedder_model_path,
        sess_options=ort.SessionOptions(),
        providers=ort.get_available_providers(),
    )
    _embedder_tokenizer = Tokenizer.from_file(embedder_tokenizer_path)



def _init_worker(
    model_path: str,
    tokenizer_path: str,
    embedder_model_path: str,
    embedder_tokenizer_path: str,
) -> None:
    """Initialize both session and tokenizer in worker process."""
    start_session(model_path)
    start_tokenizer(tokenizer_path)
    start_embedder(embedder_model_path, embedder_tokenizer_path)


def create_inference_pool(
    model_path: str,
    tokenizer_path: str,
    embedder_model_path: str,
    embedder_tokenizer_path: str,
    pool_size: int = 1,
) -> PoolType:
    """Create a pool with both session and tokenizer initialized in each worker."""
    initializer = partial(
        _init_worker,
        model_path,
        tokenizer_path,
        embedder_model_path,
        embedder_tokenizer_path,
    )
    pool = Pool(pool_size, initializer=initializer)
    logger.info(
        "Created inference pool with %d workers (model: %s, tokenizer: %s, embedder_model: %s, embedder_tokenizer: %s)",
        pool_size, model_path, tokenizer_path, embedder_model_path, embedder_tokenizer_path
    )
    return pool


def _predict_in_worker(input_data: str) -> int:
    """Function to run prediction entirely within worker process."""
    global _session

    try:
        if _tokenizer is None:
            raise Exception("Tokenizer not initialized in worker")
        if _session is None:
            raise Exception("Session not initialized in worker")

        tokens = _tokenizer.encode(input_data).ids

        # Tokenize input
        attention_mask = get_attention_mask(tokens)

        # Prepare input for ONNX model
        model_input = {
            "input_ids": [tokens],
            "attention_mask": [attention_mask],
        }

        # Run inference
        result = _session.run(None, model_input)

        if not isinstance(result[0], ndarray):
            return 0

        first_hidden_state = result[0].tolist()[0]
        label = argmax(first_hidden_state).item()

        return label

    except Exception as e:
        logger.error(f"Error in worker prediction: {e}")
        return 0


def _embed_in_worker(input_data: str) -> list:
    """Function to run embedding entirely within worker process."""
    global _embedder_model
    global _embedder_tokenizer
    try:
        if _embedder_model is None:
            raise Exception("Embedder model not initialized in worker")
        if _embedder_tokenizer is None:
            raise Exception("Embedder tokenizer not initialized in worker")
        
        # Tokenize input
        encoded = _embedder_tokenizer.encode(input_data)
        input_ids = encoded.ids
        attention_mask = encoded.attention_mask
        token_type_ids = encoded.type_ids
        
        # Prepare input for ONNX model
        model_input = {
            "input_ids": [input_ids],
            "attention_mask": [attention_mask],
            "token_type_ids": [token_type_ids],
        }
        
        # Run inference
        result = _embedder_model.run(None, model_input)
        if len(result) == 0:
            logger.error(f"Error in worker embedding: No result returned")
            return []

        if not isinstance(result[0], ndarray):
            def custom_mean(x):
                return [sum(x) / len(x) for x in x]
            embedding = custom_mean(result[0])
            return embedding

        embedding = mean(result[0], axis=1).tolist()[0]
        return embedding
    
    except Exception as e:
        logger.error(f"Error in worker embedding: {e}")
        return []


async def yield_input_ids(input_ids: Sequence) -> AsyncGenerator:
    for token in input_ids:
        yield token


def get_attention_mask(input_ids: Sequence, pad_token_id: int = 0) -> list:
    return [1 if token != pad_token_id else 0 for token in input_ids]


async def predict(
    input_data: str,
    inference_pool: PoolType,
) -> int:
    """Run prediction using worker pool."""
    try:
        # Use session_pool to run prediction in worker process
        # (tokenizer_pool is not needed since both are initialized in same workers)
        loop = asyncio.get_event_loop()
        result = await loop.run_in_executor(None, inference_pool.apply, _predict_in_worker, (input_data,))
        return result
    except Exception as e:
        logger.error(f"Error predicting: {e}")
        return 0


async def embed(
    input_data: str,
    inference_pool: PoolType,
) -> list:
    """Run prediction using worker pool."""
    try:
        # Use session_pool to run prediction in worker process
        # (tokenizer_pool is not needed since both are initialized in same workers)
        loop = asyncio.get_event_loop()
        result = await loop.run_in_executor(None, inference_pool.apply, _embed_in_worker, (input_data,))
        return result
    except Exception as e:
        logger.error(f"Error embedding: {e}")
        return []
