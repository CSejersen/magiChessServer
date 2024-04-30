import asyncio
import websockets
import time

async def connect_to_server():
    # Connect to your WebSocket server
    uri = "ws://localhost:3000/ws"
    async with websockets.connect(uri) as websocket:

        print("Sending identification msg to server")
        await websocket.send("input")
        ack = await websocket.recv()
        print(f"from server {ack}")

        while(True):
            msg = input("input move: ")
            await websocket.send("e2")
            print(f"Sending msg: {msg} to server")

            response = await websocket.recv()
            print("From server:", response)

# Run the asyncio event loop
asyncio.run(connect_to_server())

