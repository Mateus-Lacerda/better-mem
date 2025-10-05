from __future__ import print_function

from contextlib import contextmanager
import logging
import time
from argparse import ArgumentParser
from concurrent.futures import ThreadPoolExecutor
from random import choice

import grpc
from prediction_pb2 import PredictionRequest
from prediction_pb2_grpc import PredictionStub


@contextmanager
def get_stub():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = PredictionStub(channel)
        yield stub


def run_with_stub(stub, message="Hello, World!", debug=True):
    response = stub.Predict(PredictionRequest(message=message))
    if debug:
        print(f"Client received: {list(response.embedding)}")
        print("Label: " + str(response.label))
    return response

def embed(stub, message="Hello, World!", debug=True):
    response = stub.Embed(PredictionRequest(message=message))
    if debug:
        print(f"Client received: {list(response.embedding)}")

    return response


def run(message="Hello, World!", debug=True):
    if debug:
        print("Will try to connect to server")
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = PredictionStub(channel)
        response = stub.Predict(PredictionRequest(message=message))
    if debug:
        print(f"Client received: {list(response.embedding)}")
        print("Label: " + str(response.label))
    return response


def create_parser():
    parser = ArgumentParser()
    parser.add_argument("--stress", "-s", action="store_true")
    parser.add_argument("--present_tests", "-p", action="store_true")
    parser.add_argument("--message", "-m", type=str)
    parser.add_argument("--test", "-t", type=str)
    parser.add_argument("--test_embedding", "-e", action="store_true")
    return parser


if __name__ == '__main__':
    parser = create_parser()
    args = parser.parse_args()
    if not any(vars(args).values()):
        parser.print_help()
        exit(1)
    sentences = [
        {"sentence": "Hello, how are you?", "label": 0},
        {"sentence": "My name is John.", "label": 2},
        {"sentence": "I am a student.", "label": 2},
        {"sentence": "I love zoos!", "label": 2},
        {"sentence": "I saw a giraffe today.", "label": 0},
        {"sentence": "What's the weather like today?", "label": 0},
        {"sentence": "What are your plans for the weekend?", "label": 0},
        {"sentence": "My mother's birthday is next week.", "label": 2},
        {"sentence": "I am going to the zoo tomorrow.", "label": 1},
        {"sentence": "I am going to the zoo next week.", "label": 1},
        {"sentence": "I am going to the zoo next month.", "label": 1},
    ]
    results = {
        "correct": 0,
        "incorrect": 0,
    }
    with get_stub() as stub:
        def run_with_analytics(sentence, expected_label, count):
            start_time = time.time()
            response = run_with_stub(stub, sentence, False)
            end_time = time.time()
            print(f"Time taken: {end_time - start_time}")
            print(f"Expected label: {expected_label}")
            print(f"Count: {count}")
            print(f"Label: {response.label}")
            if response.label == expected_label:
                results["correct"] += 1
            else:
                results["incorrect"] += 1
            return response
                
        if args.stress:
            logging.basicConfig()
            num_tasks = 100
            start_time = time.time()
            futures = []
            count = 0
            with ThreadPoolExecutor(max_workers=500) as executor:
                for _ in range(num_tasks):
                    count += 1
                    sentence = choice(sentences)
                    futures.append(executor.submit(run_with_analytics, sentence["sentence"], sentence["label"], count))


            for future in futures:
                future.result()

            end_time = time.time()

            print(f"Time taken: {end_time - start_time}")
            print(f"Correct: {results['correct']}")
            print(f"Incorrect: {results['incorrect']}")
            print(f"Accuracy: {results['correct'] / (results['correct'] + results['incorrect'])}")
            print(f"Average time: {(end_time - start_time) / num_tasks}")

        if args.present_tests:
            for sentence in sentences:
                result = run_with_analytics(sentence["sentence"], sentence["label"], 0)
                print("--------------------------------------------------")
                print(sentence["sentence"])
                print(f"Expected label: {sentence['label']}")
                print(f"Label: {result.label}")
                print(f"Embedding: {list(result.embedding)}")


        if args.message:
            run(args.message)

        if args.test:
            while True:
                message = input("Enter message (`exit` to exit): ")
                if message == "exit":
                    break
                try:
                    run(message)
                except Exception as e:
                    print(e)
        if args.test_embedding:
            while True:
                message = input("Enter message (`exit` to exit): ")
                if message == "exit":
                    break
                try:
                    embed(stub, message)
                except Exception as e:
                    print(f"Error: {e}")
